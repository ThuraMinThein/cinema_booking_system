package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ThuraMinThein/bookings/internal/model"
	"github.com/ThuraMinThein/bookings/internal/repository"
	"github.com/ThuraMinThein/bookings/pkg/memory_storage"
	"github.com/ThuraMinThein/common/api"
)

type service struct {
	repository         *repository.Repository
	userGrpcConnection api.UserServiceClient
	seatGrpcConnection api.SeatsServiceClient
}

func NewService(repo *repository.Repository, userGrpcConnection api.UserServiceClient, seatGrpcConnection api.SeatsServiceClient) *service {
	return &service{
		repository:         repo,
		userGrpcConnection: userGrpcConnection,
		seatGrpcConnection: seatGrpcConnection,
	}
}

func (s *service) Create(request *api.CreateRequest) error {
	var bookings []*model.Booking

	user, err := s.userGrpcConnection.GetUserById(
		context.Background(),
		&api.GetUserByIdRequest{
			UserId: request.UserId,
		},
	)

	if err != nil {
		return err
	}

	for book := range request.SeatIds {

		seat, err := s.seatGrpcConnection.GetOneSeat(context.Background(), &api.GetOneSeatRequest{SeatId: request.SeatIds[book]})
		if err != nil {
			return err
		}

		isAvailable, message, err := s.IsSeatAvailable(request.MovieId, request.SeatIds[book])

		if err != nil {
			return err
		}

		if !isAvailable {
			return errors.New("Seat is already " + message)
		}

		bookings = append(bookings, &model.Booking{
			UserID:     request.UserId,
			UserName:   user.User.Name,
			MovieID:    request.MovieId,
			SeatID:     request.SeatIds[book],
			SeatNumber: seat.Seat.Name,
			Showtime:   time.Now().Add(2 * 24 * time.Hour),
		})
	}
	return s.repository.Create(bookings)
}

func (s *service) HoldBooking(request *api.HoldBookingRequest) error {

	user, err := s.userGrpcConnection.GetUserById(
		context.Background(),
		&api.GetUserByIdRequest{
			UserId: request.UserId,
		},
	)

	if err != nil {
		return err
	}

	seat, err := s.seatGrpcConnection.GetOneSeat(context.Background(), &api.GetOneSeatRequest{SeatId: request.SeatId})
	if err != nil {
		return err
	}

	isAvailable, message, err := s.IsSeatAvailable(request.MovieId, request.SeatId)
	if err != nil {
		return err
	}

	if !isAvailable {
		return errors.New("Seat is already " + message)
	}

	booking := &model.Booking{
		UserID:     request.UserId,
		UserName:   user.User.Name,
		MovieID:    request.MovieId,
		SeatID:     request.SeatId,
		SeatNumber: seat.Seat.Name,
		Showtime:   time.Now().Add(2 * 24 * time.Hour),
	}

	return memory_storage.SaveBooking(request, booking)
}

func (s *service) FindAll(userId string, movieId int64) ([]model.Booking, error) {
	return s.repository.FindAll(userId, movieId)
}

func (s *service) FindAllBookedSeats(movieId int64) ([]model.Booking, error) {
	dbBookings, err := s.repository.FindAllBookedSeats(movieId)
	if err != nil {
		return nil, err
	}

	redisBookings, err := memory_storage.GetMovieBookings(movieId)
	if err != nil {
		return nil, err
	}

	for i := range dbBookings {
		dbBookings[i].Status = "Booked"
	}

	for i := range redisBookings {
		redisBookings[i].Status = "Held"
	}

	allBookings := append(dbBookings, redisBookings...)

	return allBookings, nil
}

func (s *service) IsSeatAvailable(movieID int64, seatID int64) (bool, string, error) {

	var booking *model.Booking

	key := fmt.Sprintf("booking:%d:%d", movieID, seatID)

	if err := memory_storage.GetData(key, &booking); err != nil {
		return false, "Redis Error", err
	}

	if booking != nil {
		return false, "Held", nil
	}

	booking, err := s.repository.FindByMovieAndSeatID(movieID, seatID)
	if err != nil {
		return false, "DB Error", err
	}

	var message string
	if booking == nil {
		message = "Available"
	} else {
		message = "Booked"
	}

	return booking == nil, message, nil
}

func (s *service) Update(booking *model.Booking) error {
	return s.repository.Update(booking)
}

func (s *service) Cancel(bookingID int64) error {
	return s.repository.Cancel(bookingID)
}
