package handling

import (
	"github.com/sirupsen/logrus"
	"telegram-parser/db"
	"telegram-parser/grpc_client"
	"telegram-parser/mq"
)

type Handle struct {
	DbClient           db.DB
	Rabbit             *mq.Rabbit
	ServiceConnections *grpc_client.ServiceConnections
	HashingInfo        *HashingInfo
}

// NewHandleStruct returns a structure with connections to all key services
func NewHandleStruct(dbClient db.DB, rabbit *mq.Rabbit, serviceConnections *grpc_client.ServiceConnections) *Handle {
	return &Handle{
		DbClient:           dbClient,
		Rabbit:             rabbit,
		ServiceConnections: serviceConnections,
		HashingInfo:        &HashingInfo{},
	}
}

// HandleNewServiceAddresses handles new addresses of parser services.
// Always run in goroutine
func (h *Handle) HandleNewServiceAddresses(serviceAddresses chan string) {
	for newAddr := range serviceAddresses {
		h.ServiceConnections.CreateNewConnToService(newAddr)
	}
}

// HandleMsgsFromMQ processes new messages from telegram
func (h *Handle) HandleMsgsFromMQ() {

	rabbit := h.Rabbit
	dbClient := h.DbClient
	//serviceConnections := h.ServiceConnections

	updates := rabbit.Consume()

	for update := range updates {
		msg := mq.UnmarshalRabbitBody(update.Body)

		//logrus.Infof("RECIEVED NEW MESSAGE FROM RABBIT: %#v \n\n", msg)

		err := dbClient.Insert(msg)
		if err != nil {
			// when the unresolved rabbit timeout expires, the message will be returned to the queue
			logrus.Error(err)
		}

		logrus.Info("Consumer получил сообщение из rabbit'а\n\n")

		//	TODO сообщение отправляется в сервис (с помощью консистентного хеширования) для дальнейшего наблюдения

		//	1. Вычисляем уникальное имя сообщения =>  "chatId:messageId"
		//	2. Вычисляем номер ноды
		//	3. Отправляем

		//	HASHING.GO

		//	передавать serviceConnections в hashing
		//GetServiceAddr()

		//	TODO send ACK true if OK
	}
}
