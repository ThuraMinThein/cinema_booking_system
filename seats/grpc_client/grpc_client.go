package grpc_client

import (
	pb "github.com/ThuraMinThein/common/api"
	"github.com/ThuraMinThein/seats/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var bc pb.BookingServiceClient

func InitGrpcClients() *grpc.ClientConn {

	bookingConn, err := grpc.NewClient(config.Config.BookingsServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Failed to create gRPC client: %v", err)
	}

	logrus.WithField("port", config.Config.BookingsServiceAddress).Info("Connected to gRPC booking service")
	bc = pb.NewBookingServiceClient(bookingConn)

	return bookingConn
}

func GetBookingClient() pb.BookingServiceClient {
	return bc
}
