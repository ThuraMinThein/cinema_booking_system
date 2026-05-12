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
	if err := h.seatService.SetSeats(); err != nil {
		return nil, err
	}

	return &api.SetSeatsResponse{
		Message: "Set seats successfully",
	}, nil
}

func (h *seatGRPCHandler) GetSeats(c context.Context, request *api.GetSeatsRequest) (*api.GetSeatsResponse, error) {
	seats, err := h.seatService.FindAll(request.MovieId)
	if err != nil {
		return nil, err
	}

	var seatList []*api.Seat
	for _, seat := range seats {
		seatList = append(seatList, &api.Seat{
			Id:     int64(seat.ID),
			Name:   seat.SeatNumber,
			Row:    seat.RowNumber,
			Column: seat.ColumnNumber,
			Status: seat.Status,
		})
	}

	return &api.GetSeatsResponse{
		Seats: seatList,
	}, nil
}

func (h *seatGRPCHandler) GetOneSeat(c context.Context, request *api.GetOneSeatRequest) (*api.GetOneSeatResponse, error) {
	seat, err := h.seatService.FindOne(request.SeatId)
	if err != nil {
		return nil, err
	}

	return &api.GetOneSeatResponse{
		Seat: &api.Seat{
			Id:     int64(seat.ID),
			Name:   seat.SeatNumber,
			Row:    seat.RowNumber,
			Column: seat.ColumnNumber,
		},
	}, nil
}

func (h *seatGRPCHandler) DeleteSeat(c context.Context, request *api.DeleteSeatRequest) (*api.DeleteSeatResponse, error) {
	return nil, nil
}
