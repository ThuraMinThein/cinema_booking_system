package service

import (
	"context"
	"errors"
	"time"

	"github.com/ThuraMinThein/common/api"
	"github.com/ThuraMinThein/seats/internal/model"
	"github.com/ThuraMinThein/seats/internal/repository"
)

type service struct {
	repository            *repository.Repository
	bookingGrpcConnection api.BookingServiceClient
}

func NewService(repo *repository.Repository, bookingGrpcConnection api.BookingServiceClient) *service {
	return &service{
		repository:            repo,
		bookingGrpcConnection: bookingGrpcConnection,
	}
}

func (s *service) SetSeats() error {
	hasSeats, err := s.repository.FindAll()
	if err != nil {
		return err
	}
	if len(hasSeats) > 0 {
		return errors.New("Seats Already Created")
	}
	seats := []*model.Seat{
		{SeatNumber: "A1", ColumnNumber: "1", RowNumber: "A"},
		{SeatNumber: "A2", ColumnNumber: "2", RowNumber: "A"},
		{SeatNumber: "A3", ColumnNumber: "3", RowNumber: "A"},
		{SeatNumber: "A4", ColumnNumber: "4", RowNumber: "A"},
		{SeatNumber: "A5", ColumnNumber: "5", RowNumber: "A"},
		{SeatNumber: "A6", ColumnNumber: "6", RowNumber: "A"},
		{SeatNumber: "A7", ColumnNumber: "7", RowNumber: "A"},
		{SeatNumber: "A8", ColumnNumber: "8", RowNumber: "A"},
		{SeatNumber: "A9", ColumnNumber: "9", RowNumber: "A"},
		{SeatNumber: "A10", ColumnNumber: "10", RowNumber: "A"},
		{SeatNumber: "B1", ColumnNumber: "1", RowNumber: "B"},
		{SeatNumber: "B2", ColumnNumber: "2", RowNumber: "B"},
		{SeatNumber: "B3", ColumnNumber: "3", RowNumber: "B"},
		{SeatNumber: "B4", ColumnNumber: "4", RowNumber: "B"},
		{SeatNumber: "B5", ColumnNumber: "5", RowNumber: "B"},
		{SeatNumber: "B6", ColumnNumber: "6", RowNumber: "B"},
		{SeatNumber: "B7", ColumnNumber: "7", RowNumber: "B"},
		{SeatNumber: "B8", ColumnNumber: "8", RowNumber: "B"},
		{SeatNumber: "B9", ColumnNumber: "9", RowNumber: "B"},
		{SeatNumber: "B10", ColumnNumber: "10", RowNumber: "B"},
	}
	return s.repository.Create(seats)
}

func (s *service) FindAll(movieId int64) ([]*model.Seat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	seats, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	bookings, err := s.bookingGrpcConnection.FindAllBookedSeats(ctx, &api.FindAllBookedSeatsRequest{
		MovieId: movieId,
	})
	if err != nil {
		return nil, err
	}

	bookedMap := make(map[int64]string)

	for _, booking := range bookings.Bookings {
		bookedMap[booking.SeatId] = booking.Status
	}

	var seatList []*model.Seat
	for _, seat := range seats {

		status := "Available"

		if bookedMap[seat.ID] != "" {
			status = bookedMap[seat.ID]
		}

		seatList = append(seatList, &model.Seat{
			ID:           seat.ID,
			SeatNumber:   seat.SeatNumber,
			ColumnNumber: seat.ColumnNumber,
			RowNumber:    seat.RowNumber,
			Status:       status,
		})
	}
	return seatList, nil

}

func (s *service) FindOne(seatId int64) (*model.Seat, error) {
	return s.repository.FindOne(seatId)
}

func (s *service) Delete(seatID int64) error {
	return s.repository.Delete(seatID)
}
