package grpc_client

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"telegram-parser/proto"
)

//    --------------------------------------------------------------------------------
//                                    STRUCTS
//    --------------------------------------------------------------------------------

// ServiceConnections contains connections to all replicas of the parser
type ServiceConnections struct {
	Clients             map[string]proto.ParserClient // string - is parser service address; proto.ParserClient - connection to service
	UnavailableServices []string                      // addresses of services to which it was not possible to connect
}

//    --------------------------------------------------------------------------------
//                                     METHODS
//    --------------------------------------------------------------------------------

// NewParserConnections returns new instance of ServiceConnections
func NewParserConnections() *ServiceConnections {
	return &ServiceConnections{Clients: make(map[string]proto.ParserClient)}
}

// CreateNewConnToService creates gRPC client for the transferred node address
func (p *ServiceConnections) CreateNewConnToService(addr string) {
	if _, contains := p.Clients[addr]; contains {
		// Client already exists
		logrus.Errorf("The address `%v` is already in the nodes list\n", addr)
		return
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		// Unable to connect to the service
		p.UnavailableServices = append(p.UnavailableServices, addr)
		return
	}

	parserClient := proto.NewParserClient(conn)

	p.Clients[addr] = parserClient
}

// GetConnectionToService returns the gRPC client by the service address
func (p *ServiceConnections) GetConnectionToService(nodeAddr string) (proto.ParserClient, error) {
	client, ok := p.Clients[nodeAddr]
	if !ok {
		return nil, errors.New(fmt.Sprintf("The parser service with the address `%v` does not exist", nodeAddr))
	}

	return client, nil
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
//	resp, err := client.NewTrackedMsg(ctx, &proto.NewTrackedMsgRequest{
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
