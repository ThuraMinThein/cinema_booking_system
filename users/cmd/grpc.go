package main

import (
	"net"

	"github.com/ThuraMinThein/users/internal/handler"
	"github.com/ThuraMinThein/users/internal/repository"
	"github.com/ThuraMinThein/users/internal/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type grpcServer struct {
	address string
}

func NewGRPCServer(address string) *grpcServer {
	return &grpcServer{address: address}
}

func (s *grpcServer) Run() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	// register the service
	repository := repository.NewRepository()
	svc := service.NewService(repository)
	handler.NewGRPCUsersService(grpcServer, svc)

	logrus.WithField("port", s.address).Info("User gRPC server is running")

	return grpcServer.Serve(lis)
}
