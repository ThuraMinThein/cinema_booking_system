package main

import (
	"github.com/ThuraMinThein/bookings/config"
	"github.com/ThuraMinThein/bookings/database"
	"github.com/sirupsen/logrus"
)

func main() {

	config.LoadConfig()

	if err := database.DatabaseInit(); err != nil {
		logrus.Fatal("Database connection failed", err)
	}
	grpcServer := NewGRPCBookingServer(config.Config.BookingsServiceAddress)
	if err := grpcServer.Run(); err != nil {
		logrus.Fatal("User GRPC Service", err)
	}
}
