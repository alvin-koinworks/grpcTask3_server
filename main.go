package main

import (
	"log"
	"net"

	proto "serverGRPC/resources/proto"
	depositServer "serverGRPC/resources/service"

	"google.golang.org/grpc"
)

const port = ":9000"

func main() {
	listen, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatal("Failed to connect")
	}

	serve := grpc.NewServer()

	proto.RegisterDepositServiceServer(serve, &depositServer.DepositService{})

	log.Printf("Server connect at %v", listen.Addr())

	if err := serve.Serve(listen); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
}
