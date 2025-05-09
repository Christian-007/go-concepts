package main

import (
	"fmt"
	"grpc/client"
	pb "grpc/proto/addressbook"
	"grpc/server"
	"log/slog"
	"os"
)


func main() {
	fmt.Println("Hello from ./grpc/main.go")

	person := &pb.Person{
		Id: 1234,
		Name: "John Doe",
		Email: "jdoe@example.com",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "555-4321", Type: pb.PhoneType_PHONE_TYPE_HOME},
		},
	}

	fmt.Println("[./grpc/main.go] Person:", person)

	// The following is for running the grpc client and server
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	grpcServer := server.NewGrpcServer(logger)
	go func() {
		if err := grpcServer.Start(); err != nil {
			logger.Error("Failed to start the gRPC server",
				slog.String("error", err.Error()),
			)
		}
	}()

	// Wait for gRPC server is running properly
	<-grpcServer.Running()

	grpcClient := client.NewGrpcClient(logger)
	if err := grpcClient.Call(); err != nil {
		logger.Error("Failed to make a request from gRPC client",
				slog.String("error", err.Error()),
			)
	}
}