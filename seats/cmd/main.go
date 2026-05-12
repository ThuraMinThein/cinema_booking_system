package main

import (
	"github.com/ThuraMinThein/seats/config"
	"github.com/ThuraMinThein/seats/database"
	"github.com/ThuraMinThein/seats/grpc_client"
	"github.com/sirupsen/logrus"
)

func main() {

	config.LoadConfig()

	if err := database.DatabaseInit(); err != nil {
		logrus.Fatal("Database connection failed", err)
	}

	bc := grpc_client.InitGrpcClients()

	defer func() {
		bc.Close()
	}()
	grpcServer := NewGRPCSeatsServer(config.Config.SeatsServiceAddress)
	if err := grpcServer.Run(); err != nil {
		logrus.Fatal("User GRPC Service", err)
	}
}
