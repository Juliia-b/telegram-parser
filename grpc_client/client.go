package grpc_client

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
	"telegram-parser/proto"
	"time"
)

//    --------------------------------------------------------------------------------
//                                    STRUCTS
//    --------------------------------------------------------------------------------

// ServiceConnections contains connections to all replicas of the parser
type ServiceConnections struct {
	Clients map[string]proto.ParserClient // string - is parser service address; proto.ParserClient - connection to service
	//UnavailableServices []string                      // addresses of services to which it was not possible to connect
}

//    --------------------------------------------------------------------------------
//                                     METHODS
//    --------------------------------------------------------------------------------

// NewParserConnections returns new instance of ServiceConnections
func NewParserConnections() *ServiceConnections {
	return &ServiceConnections{Clients: make(map[string]proto.ParserClient)}
}

// ConnectToService creates gRPC client for the transferred node address
func (s *ServiceConnections) ConnectToService(addr string) (ok bool) {
	_, contains := getParserClientFromMap(s, addr)
	if contains {
		// Client already exists
		logrus.Errorf("The address `%v` is already in the nodes list\n", addr)
		return true
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		// Unable to connect to the service
		logrus.Errorf("Failed to connect to the service at %v\n", addr)
		// TODO возможно стоит убрать UnavailableServices ?
		//s.UnavailableServices = append(s.UnavailableServices, addr)
		return false
	}

	parserClient := proto.NewParserClient(conn)

	addParserClientToMap(s, addr, parserClient)

	logrus.Infof("A connection has been established with the 'parser' microservice at %v\n", addr)

	return true
}

// GetConnectionToService returns the gRPC client by the service address
func (s *ServiceConnections) GetConnectionToService(nodeAddr string) (proto.ParserClient, error) {
	client, ok := getParserClientFromMap(s, nodeAddr)
	if !ok {
		return nil, errors.New(fmt.Sprintf("The parser service with the address `%v` does not exist", nodeAddr))
	}

	return client, nil
}

// SendAddMsgRequest sends a request "AddMsg" to the service
func (s *ServiceConnections) SendAddMsgRequest(msgKey string, serviceAddr string) (isNodeFailed bool, panicErr error) {
	isNodeFailed = false
	panicErr = nil

	clientConn, panicErr := s.GetConnectionToService(serviceAddr)
	if panicErr != nil {
		return
	}

	// Add timeout to the request
	ctx, _ := context.WithTimeout(context.TODO(), 2*time.Second)
	_, err := clientConn.AddMsg(ctx, &proto.AddMsgRequest{MsgKey: msgKey})
	if err != nil {
		// No response has been received from the microservice when the timeout expires.
		logrus.Errorf("Error in add msg is %v \n", err)
		isNodeFailed = true
		return
	}

	return
}

//    --------------------------------------------------------------------------------
//                                     HELPERS
//    --------------------------------------------------------------------------------

func addParserClientToMap(s *ServiceConnections, key string, value proto.ParserClient) {
	mu := sync.RWMutex{}
	mu.Lock()
	s.Clients[key] = value
	mu.Unlock()
}

func getParserClientFromMap(s *ServiceConnections, key string) (client proto.ParserClient, exist bool) {
	mu := sync.RWMutex{}
	mu.RLock()
	client, exist = s.Clients[key]
	mu.RUnlock()

	return client, exist
}

//    --------------------------------------------------------------------------------
//                                        EXTRA
//    --------------------------------------------------------------------------------

// ParserClientsInit creates gRPC clients for the transferred node addresses
//func (p *ServiceConnections) ParserClientsInit(conf *flags.Config) {
//	for _, addr := range conf.ParserAddrs {
//		if _, contains := p.Clients[addr]; contains {
//			logrus.Infof("The address `%v` is already in the nodes list\n", addr)
//			continue
//		}
//
//		var conn *grpc.ClientConn
//		conn, err := grpc.Dial(addr, grpc.WithInsecure())
//		if err != nil {
//			logrus.Error(err)
//
//			// TODO пометить адрес как не релевантный, удалить из адресов
//
//
//			continue
//		}
//
//		parserClient := proto.NewParserClient(conn)
//
//		p.Clients[addr] = parserClient
//	}
//}

//func (p *ServiceConnections) NewTrackedMsg(nodeAddr string, msg *db.Message) error {
//	client, ok := p.Clients[nodeAddr]
//	if !ok {
//		return errors.New(fmt.Sprintf("The parser service with the address `%v` does not exist", nodeAddr))
//	}
//
//	ctx := context.TODO()
//	resp, err := client.NewTrackedMsg(ctx, &proto.sendNewTrackedMsg{
//		MessageID: msg.MessageID,
//		ChatId:    msg.ChatID,
//		ChatTitle: msg.ChatTitle,
//		Content:   msg.Content,
//		Date:      msg.Date,
//		Views:     msg.Views,
//		Forwards:  msg.Forwards,
//		Replies:   msg.Replies,
//	})
//
//	if err != nil {
//		return err
//	}
//
//	if !resp.Received {
//		//TODO	Сообщение по каким-то причинам не поступило в обработку. Попробовать снова через время
//	}
//
//	return nil
//}
