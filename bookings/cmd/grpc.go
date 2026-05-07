package main

import (
	"net"

	"github.com/ThuraMinThein/bookings/database"
	"github.com/ThuraMinThein/bookings/internal/handler"
	"github.com/ThuraMinThein/bookings/internal/repository"
	"github.com/ThuraMinThein/bookings/internal/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type grpcServer struct {
	address string
}

func NewGRPCBookingServer(address string) *grpcServer {
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
	handler.NewGRPCBookingsService(grpcServer, svc)

	logrus.WithField("port", s.address).Info("Bookings gRPC server is running")

	return grpcServer.Serve(lis)
}
