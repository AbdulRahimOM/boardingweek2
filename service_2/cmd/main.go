package main

import (
	"boarding-week2/service_2/config"
	"boarding-week2/service_2/server"
	"log"
	"net"

	pb "boarding-week2/pb"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":"+config.EnvValues.Svc2Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Println("Account Service >>>>>> Listening on port: ", config.EnvValues.Svc2Port)
	}

	svc2Server := server.InitializeServer()
	grpcServer := grpc.NewServer()
	pb.RegisterSvc2Server(grpcServer, svc2Server)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln("failed to serve", err)
	}

}
