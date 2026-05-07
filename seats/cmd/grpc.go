package main

import (
	"net"

	"github.com/ThuraMinThein/seats/database"
	"github.com/ThuraMinThein/seats/internal/handler"
	"github.com/ThuraMinThein/seats/internal/repository"
	"github.com/ThuraMinThein/seats/internal/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type grpcServer struct {
	address string
}

func NewGRPCSeatsServer(address string) *grpcServer {
	return &grpcServer{
		address: address,
	}
}

func (s *grpcServer) Run() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	// register the service
	repository := repository.NewRepository(database.DB)
	svc := service.NewService(repository)
	handler.NewGRPCSeatsService(grpcServer, svc)

	logrus.WithField("port", s.address).Info("Seats gRPC server is running")

	return grpcServer.Serve(lis)
}
