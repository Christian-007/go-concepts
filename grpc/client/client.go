package client

import (
	"context"
	pb "grpc/proto/addressbook"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	logger *slog.Logger
}

func NewGrpcClient(logger *slog.Logger) *GrpcClient {
	return &GrpcClient{logger: logger}
}

func (c *GrpcClient) Call() error {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	addressBookService := pb.NewAddressBookServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = addressBookService.GetAll(ctx, &pb.GetAllParams{})
	if err != nil {
		return err
	}

	c.logger.Info("Successfully made a request from gRPC Client")
	return nil
}
