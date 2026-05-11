package handler

import (
	"context"

	"github.com/ThuraMinThein/bookings/pkg/types"
	"github.com/ThuraMinThein/common/api"
	"google.golang.org/grpc"
)

type bookingGRPCHandler struct {
	bookingService types.BookingService
	api.UnimplementedBookingServiceServer
}

func NewGRPCBookingsService(grpc *grpc.Server, bookingService types.BookingService) {
	grpcHandler := &bookingGRPCHandler{
		bookingService: bookingService,
	}

	api.RegisterBookingServiceServer(grpc, grpcHandler)
}

func (h *bookingGRPCHandler) Create(c context.Context, request *api.CreateRequest) (*api.CreateResponse, error) {
	if err := h.bookingService.Create(request); err != nil {
		return nil, err
	}
	return &api.CreateResponse{
		Status: "success",
	}, nil
}
func (h *bookingGRPCHandler) FindAll(c context.Context, request *api.FindAllRequest) (*api.FindAllResponse, error) {
	bookings, err := h.bookingService.FindAll(request.UserId, request.MovieId)
	if err != nil {
		return nil, err
	}

	bookingResponse := make([]*api.Booking, 0, len(bookings))
	for _, booking := range bookings {
		bookingResponse = append(bookingResponse, &api.Booking{
			UserId:     booking.UserID,
			UserName:   booking.UserName,
			MovieId:    booking.MovieID,
			SeatId:     booking.SeatID,
			SeatNumber: booking.SeatNumber,
		})
	}

	return &api.FindAllResponse{
		Bookings: bookingResponse,
	}, nil
}
func (h *bookingGRPCHandler) IsSeatAvailable(c context.Context, request *api.IsSeatAvailableRequest) (*api.IsSeatAvailableResponse, error) {
	return nil, nil
}
func (h *bookingGRPCHandler) Update(c context.Context, request *api.UpdateRequest) (*api.UpdateResponse, error) {
	return nil, nil
}
func (h *bookingGRPCHandler) Cancel(c context.Context, request *api.CancelRequest) (*api.CancelResponse, error) {
	return nil, nil
}
