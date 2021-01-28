package grpc_client

import (
	"google.golang.org/grpc"
	"telegram-parser/proto"
)

func ParserClientInit() (proto.ParserClient, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":3000", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	parserClient := proto.NewParserClient(conn)

	return parserClient, err
}
