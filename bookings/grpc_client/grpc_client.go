package grpc_client

import (
	"github.com/ThuraMinThein/bookings/config"
	pb "github.com/ThuraMinThein/common/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var uc pb.UserServiceClient
var sc pb.SeatsServiceClient

func InitGrpcClients() (*grpc.ClientConn, *grpc.ClientConn) {

	usersConn, err := grpc.NewClient(config.Config.UsersServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Failed to create gRPC client: %v", err)
	}

	logrus.WithField("port", config.Config.UsersServiceAddress).Info("Connected to gRPC users service")

	uc = pb.NewUserServiceClient(usersConn)

	seatsConn, err := grpc.NewClient(config.Config.SeatsServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Failed to create gRPC client: %v", err)
	}

	logrus.WithField("port", config.Config.SeatsServiceAddress).Info("Connected to gRPC seats service")
	sc = pb.NewSeatsServiceClient(seatsConn)

	return usersConn, seatsConn
}

func GetUserClient() pb.UserServiceClient {
	return uc
}

func GetSeatsClient() pb.SeatsServiceClient {
	return sc
}
