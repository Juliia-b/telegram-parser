package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"telegram-parser/proto"
)

func main() {
	port := parseFlags()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		logrus.Fatal(err)
	}

	s := Server{}
	grpcServer := grpc.NewServer()

	proto.RegisterParserServer(grpcServer, &s)

	logrus.Info("Service is running")

	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %s", err)
	}
}

type Server struct {
}

func (s *Server) Ping(ctx context.Context, in *proto.PingRequest) (*proto.PingResponse, error) {
	logrus.Infof("PING recieved `%v`", in.Msg)
	return &proto.PingResponse{Msg: "Pong"}, nil
}

func (s *Server) AddMsg(ctx context.Context, in *proto.AddMsgRequest) (*proto.AddMsgResponse, error) {
	logrus.Infof("NEWTRACKEDMSG recieved msgkey `%v`", in.MsgKey)
	return &proto.AddMsgResponse{Processed: true}, nil
	//	if *proto.AddMsgResponse == nil
	//	err = rpc error: code = Internal desc = grpc: error while marshaling: proto: Marshal called with nil
}

//    --------------------------------------------------------------------------------
//                                     HELPERS
//    --------------------------------------------------------------------------------

func parseFlags() string {
	port := flag.String("port", "foo", "")
	flag.Parse()

	return *port
}
