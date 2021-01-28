package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"telegram-parser/proto"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 3000))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	s := Server{}
	grpcServer := grpc.NewServer()

	proto.RegisterParserServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %s", err)
	}
}

type Server struct {
}

func (s *Server) Ping(ctx context.Context, in *proto.PingRequest) (*proto.PingResponse, error) {
	logrus.Infof("PING recieved %v", in.Msg)

	return &proto.PingResponse{Msg: "Pong"}, nil
}

func (s *Server) NewTrackedMsg(ctx context.Context, in *proto.NewTrackedMsgRequest) (*proto.NewTrackedMsgResponse, error) {
	logrus.Infof("NEWTRACKEDMSG recieved %v", in.Content)
	return nil, nil
}
