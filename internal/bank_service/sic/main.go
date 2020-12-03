package main

import (
	"context"
	"fmt"
	"net"
	"tpayment/internal/bank_service/bank_common"

	"google.golang.org/grpc"
)

type Server struct {
}

func (s *Server) BaseTxn(ctx context.Context, in *bank_common.BaseRequest) (*bank_common.BaseReply, error) {

	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50001")
	if err != nil {
		fmt.Println("net.Listen fail->", err.Error())
		return
	}

	s := grpc.NewServer()

	bank_common.RegisterTxnServer(s, &Server{})
	err = s.Serve(lis)
	if err != nil {
		panic(err.Error())
	}
}
