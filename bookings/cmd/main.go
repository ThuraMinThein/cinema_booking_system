package main

import (
	"github.com/ThuraMinThein/bookings/config"
	"github.com/ThuraMinThein/bookings/database"
	"github.com/ThuraMinThein/bookings/grpc_client"
	redis_client "github.com/ThuraMinThein/bookings/pkg/redis"
	"github.com/sirupsen/logrus"
)

func main() {

	config.LoadConfig()

	if err := database.DatabaseInit(); err != nil {
		logrus.Fatal("Database connection failed", err)
	}

	uc, sc := grpc_client.InitGrpcClients()

	if err := redis_client.InitRedis(); err != nil {
		logrus.Fatal("Redis connection failed", err)
	}

	defer func() {
		uc.Close()
		sc.Close()
	}()

	grpcServer := NewGRPCBookingServer(config.Config.BookingsServiceAddress)
	if err := grpcServer.Run(); err != nil {
		logrus.Fatal("User GRPC Service", err)
	}
}
