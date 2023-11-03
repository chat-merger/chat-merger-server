package main

import (
	"chatmerger/mergerapi"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

func main() {
	println("start")
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("running server failed: %v", err)
	}
	grpcServer := grpc.NewServer()

	s := new(mergerapi.Server)

	mergerapi.RegisterBaseServiceServer(grpcServer, s)

	//for i := 0; i < 100; i++ {
	//	go clientRun(fmt.Sprintf("cl_№%d", i))
	//}
	go clientRun("cl_0")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpc server failed: %v", err)
	}
	println("end")
}

func clientRun(name string) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed create client: %v", err)
	}

	defer conn.Close()

	serviceClient := mergerapi.NewBaseServiceClient(conn)

	messageClient, err := serviceClient.Connect(context.Background())
	if err != nil {
		log.Fatalf("failed CreateMessage: %v", err)
	}

	for {
		var msg = new(mergerapi.Response)
		err = messageClient.RecvMsg(msg)
		if err != nil {
			log.Fatalf("failed RecvMsg: %v", err)
		}
		fmt.Printf("client (%s) > recive msg: \n\t%#v\n", name, msg.Event)
	}
}
