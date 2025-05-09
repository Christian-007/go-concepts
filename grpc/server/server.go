package server

import (
	"context"
	"fmt"
	pb "grpc/proto/addressbook"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

// AddressBookService implementation
type AddressBookServiceServer struct {
	pb.UnimplementedAddressBookServiceServer
}

func (a *AddressBookServiceServer) GetAll(context.Context, *pb.GetAllParams) (*pb.AddressBook, error) {
	fmt.Println("This is a dummy AddressBookServiceServer.GetAll() method")
	return &pb.AddressBook{}, nil
}

// gRPC server code
type GrpcServer struct {
	ready chan struct{}
	logger *slog.Logger
}

func NewGrpcServer(logger *slog.Logger) *GrpcServer {
	return &GrpcServer{
		ready: make(chan struct{}),
		logger: logger,
	}
}

func (s *GrpcServer) Start() error {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAddressBookServiceServer(grpcServer, &AddressBookServiceServer{})

	go func() {
		s.logger.Info("gRPC server listening on :50051")
		close(s.ready)
	}()

	return grpcServer.Serve(listener)
}

func (s *GrpcServer) Running() <-chan struct{} {
	return s.ready
}
