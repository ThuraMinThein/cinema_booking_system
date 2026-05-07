package handler

import (
	"context"

	"github.com/ThuraMinThein/common/api"
	"github.com/ThuraMinThein/seats/pkg/types"
	"google.golang.org/grpc"
)

type seatGRPCHandler struct {
	seatService types.SeatService
	api.UnimplementedSeatsServiceServer
}

func NewGRPCSeatsService(grpc *grpc.Server, seatService types.SeatService) {
	grpcHandler := &seatGRPCHandler{
		seatService: seatService,
	}

	api.RegisterSeatsServiceServer(grpc, grpcHandler)
}

func (h *seatGRPCHandler) SetSeats(c context.Context, request *api.SetSeatsRequest) (*api.SetSeatsResponse, error) {
	return nil, nil
}

func (h *seatGRPCHandler) GetSeats(c context.Context, request *api.GetSeatsRequest) (*api.GetSeatsResponse, error) {
	return nil, nil
}

func (h *seatGRPCHandler) DeleteSeat(c context.Context, request *api.DeleteSeatRequest) (*api.DeleteSeatResponse, error) {
	return nil, nil
}
