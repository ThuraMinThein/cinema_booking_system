package main

import (
	"github.com/ThuraMinThein/gateway/config"
	"github.com/sirupsen/logrus"
)

func main() {

	config.LoadConfig()

	grpcServer := NewGRPCServer(config.Config.UsersServiceAddress)
	if err := grpcServer.Run(); err != nil {
		logrus.Fatal("User GRPC Service", err)
	}
}
