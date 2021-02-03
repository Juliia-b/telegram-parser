package handling

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"telegram-parser/db"
	"telegram-parser/grpc_client"
	"telegram-parser/mq"
	"time"
)

//    --------------------------------------------------------------------------------
//                                    STRUCTS
//    --------------------------------------------------------------------------------

type Handler struct {
	ServiceConnections *grpc_client.ServiceConnections
	HashingInfo        *HashRing
	DbClient           db.DB
	Rabbit             *mq.Rabbit
}

//    --------------------------------------------------------------------------------
//                                     METHODS
//    --------------------------------------------------------------------------------

// NewHandleStruct returns a structure with connections to all key services
func NewHandleStruct(serviceConnections *grpc_client.ServiceConnections, dbClient db.DB, rabbit *mq.Rabbit) *Handler {
	return &Handler{
		ServiceConnections: serviceConnections,
		HashingInfo:        hashRingInit(),
		DbClient:           dbClient,
		Rabbit:             rabbit,
	}
}

// HandleNewServiceAddresses handles new addresses of parser services.
// Always run in goroutine
func (h *Handler) HandleNewServiceAddresses(serviceAddresses chan string) {
	for nodeAddr := range serviceAddresses {
		ok := h.ServiceConnections.ConnectToService(nodeAddr)
		if ok {
			//  Connection to the service was successful.
			//	Need to create a rb tree
			h.HashingInfo.CreateRBTree(nodeAddr)
			//	Add node address to hash ring
			h.HashingInfo.addNodeToHashRing(nodeAddr)
		}
	}
}

// RunHandlingMsgsFromMQ launches multiple handlers.
// @param handlersCount - number of message handlers to run
func (h *Handler) RunHandlingMsgsFromMQ(handlersCount int) {
	if handlersCount < 10 {
		handlersCount = 10
	}

	for i := 0; i < handlersCount; i++ {
		go h.handleMsgsFromMQ()
	}
}

// handleMsgsFromMQ processes new messages from telegram
func (h *Handler) handleMsgsFromMQ() {
	rabbit := h.Rabbit
	dbClient := h.DbClient

	updates := rabbit.Consume()

	for update := range updates {
		msg := mq.UnmarshalRabbitBody(update.Body)

		logrus.Infof("Получено сообщение с id %v из rabbit \n\n", msg.MessageID)

		// TODO будет ли ошибка при внесении одинаковых зн?
		err := dbClient.Insert(msg)
		if err != nil {
			logrus.Error(err)
			update.Nack(false, false)
			continue
		}

		msgKey := createUniqueKey(msg)

		nodeAddr, ok := h.HashingInfo.calculateNodeAddr(msgKey)
		if !ok {
			logrus.Errorf("Can't get node address from hash ring. Unique msgKey is %v \n", msgKey)
			update.Nack(false, false)
			continue
		}

		isNodeFailed, panicErr := h.ServiceConnections.SendAddMsgRequest(msgKey, nodeAddr)
		if panicErr != nil {
			update.Nack(false, false)
			logrus.Panic(panicErr)
		}

		logrus.Infof("Сообщение %v отправлено в микросервис %v\n", msgKey, nodeAddr)

		if isNodeFailed {
			//	Microservice fell. the message was not delivered. Need to move data to another microservice
			logrus.Errorf("Microservice '%v' fell.\n", nodeAddr)

			h.handleServiceDrop(nodeAddr)

			update.Nack(false, false)
			continue
		}

		h.HashingInfo.AddValToTree(nodeAddr, msgKey)

		// message processed
		update.Ack(false)
	}
}

// handleServiceDrop removes a service from the list of available services.
// Redistributes the load of the fallen service to other services. Starts pinging the node
func (h *Handler) handleServiceDrop(nodeAddr string) {

	h.HashingInfo.DeleteRbTree(nodeAddr)
	h.HashingInfo.removeNodeFromHashRing(nodeAddr)

	// TODO перенаправить сообщения в др сервисы
	h.redirectMsgs(nodeAddr)

	// endless attempt to connect to the service
	go h.tryConnectToServiceAfterFail(nodeAddr)
}

// redirectMsgs distributes all messages of the crashed service to the remaining available services
func (h *Handler) redirectMsgs(failedNodeAddr string) {
	values := h.HashingInfo.GetTreeValues(failedNodeAddr)

	for _, msgKey := range values {
		nodeAddr, ok := h.HashingInfo.calculateNodeAddr(msgKey)
		if !ok {
			logrus.Errorf("Can't get node address from hash ring. Unique msgKey is %v \n", msgKey)
			continue
		}

		isNodeFailed, panicErr := h.ServiceConnections.SendAddMsgRequest(msgKey, nodeAddr)

		//isNodeFailed, panicErr := h.ServiceConnections.SendAddMsgRequest(msgKey, nodeAddr)
		if panicErr != nil {
			logrus.Panic(panicErr)
		}

		if isNodeFailed {
			//	Microservice fell. the message was not delivered. Need to move data to another microservice
			logrus.Errorf("Microservice '%v' fell.\n", nodeAddr)

			h.handleServiceDrop(nodeAddr)
		}
	}
}

// tryConnectToServiceAfterFail tries to reconnect to microservice after it crashes.
// Is in a separate goroutine, pings the service every 1 minute.
// Run in goroutine
func (h *Handler) tryConnectToServiceAfterFail(nodeAddr string) (done bool) {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-ticker.C:
			ok := h.ServiceConnections.ConnectToService(nodeAddr)
			if ok {
				ticker.Stop()

				//  Connection to the service was successful.
				//	Need to create a rb tree
				h.HashingInfo.CreateRBTree(nodeAddr)
				//	Add node address to hash ring
				h.HashingInfo.addNodeToHashRing(nodeAddr)

				// TODO перенаправить в сервис сообщения для парсинга и удалить эти сообщения с других нод
				return ok
			}
		}
	}
}

//faildetection

//    --------------------------------------------------------------------------------
//                                     HELPERS
//    --------------------------------------------------------------------------------

// createUniqueKey creates a unique key by which a search is performed in a consistent hashing circle
func createUniqueKey(msg *db.Message) string {
	return fmt.Sprintf("%v:%v", msg.ChatID, msg.MessageID)
}
