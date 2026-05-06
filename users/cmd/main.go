package main

import (
	"github.com/ThuraMinThein/users/config"
	"github.com/ThuraMinThein/users/database"
	"github.com/sirupsen/logrus"
)

func main() {

	config.LoadConfig()

	if err := database.DatabaseInit(); err != nil {
		logrus.Fatal("Database connection failed", err)
	}

	grpcServer := NewGRPCServer(config.Config.UsersServiceAddress)
	if err := grpcServer.Run(); err != nil {
		logrus.Fatal("User GRPC Service", err)
	}
}
