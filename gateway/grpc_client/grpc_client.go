package grpc_client

import (
	pb "github.com/ThuraMinThein/common/api"
	"github.com/ThuraMinThein/gateway/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var uc pb.UserServiceClient
var bc pb.BookingServiceClient
var sc pb.SeatsServiceClient

func InitGrpcClients() {

	usersConn, err := grpc.NewClient(config.Config.UsersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Failed to create gRPC client: %v", err)
	}

	defer usersConn.Close()
	logrus.WithField("port", config.Config.UsersServiceAddress).Info("Connected to gRPC users service")

	uc = pb.NewUserServiceClient(usersConn)

	bookingsConn, err := grpc.NewClient(config.Config.BookingsServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Failed to create gRPC client: %v", err)
	}

	defer bookingsConn.Close()
	logrus.WithField("port", config.Config.BookingsServiceAddress).Info("Connected to gRPC bookings service")

	bc = pb.NewBookingServiceClient(bookingsConn)

	seatsConn, err := grpc.NewClient(config.Config.SeatsServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Failed to create gRPC client: %v", err)
	}

	defer seatsConn.Close()
	logrus.WithField("port", config.Config.SeatsServiceAddress).Info("Connected to gRPC seats service")
	sc = pb.NewSeatsServiceClient(seatsConn)

}

func GetUserClient() pb.UserServiceClient {
	return uc
}

func GetBookingClient() pb.BookingServiceClient {
	return bc
}

func GetSeatsClient() pb.SeatsServiceClient {
	return sc
}
