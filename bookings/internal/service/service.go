package service

import (
	"context"
	"errors"
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

		isAvailable, err := s.IsSeatAvailable(request.MovieId, request.SeatIds[book])

		if err != nil {
			return err
		}

		if !isAvailable {
			return errors.New("Seat is already taken")
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

	isAvailable, err := s.IsSeatAvailable(request.MovieId, request.SeatId)
	if err != nil {
		return err
	}

	if !isAvailable {
		return errors.New("Seat is already taken")
	}

	booking := &model.Booking{
		UserID:     request.UserId,
		UserName:   user.User.Name,
		MovieID:    request.MovieId,
		SeatID:     request.SeatId,
		SeatNumber: seat.Seat.Name,
		Showtime:   time.Now().Add(2 * 24 * time.Hour),
	}

	return memory_storage.SetData(
		"booking:"+string(request.MovieId)+string(request.SeatId),
		booking,
		2*time.Minute,
	)
}

func (s *service) FindAll(userId string, movieId int64) ([]model.Booking, error) {
	return s.repository.FindAll(userId, movieId)
}

func (s *service) IsSeatAvailable(movieID int64, seatID int64) (bool, error) {

	var booking *model.Booking

	if err := memory_storage.GetData("booking:"+string(movieID)+string(seatID), &booking); err != nil {
		return false, err
	}

	if booking != nil {
		return false, nil
	}

	booking, err := s.repository.FindByMovieAndSeatID(movieID, seatID)
	if err != nil {
		return false, err
	}

	return booking == nil, nil
}

func (s *service) Update(booking *model.Booking) error {
	return s.repository.Update(booking)
}

func (s *service) Cancel(bookingID int64) error {
	return s.repository.Cancel(bookingID)
}
