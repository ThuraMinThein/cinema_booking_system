package service

import (
	"context"
	"time"

	"github.com/ThuraMinThein/bookings/grpc_client"
	"github.com/ThuraMinThein/bookings/internal/model"
	"github.com/ThuraMinThein/bookings/internal/repository"
	"github.com/ThuraMinThein/common/api"
)

type service struct {
	repository *repository.Repository
}

func NewService(repo *repository.Repository) *service {
	return &service{
		repository: repo,
	}
}

func (s *service) Create(request *api.CreateRequest) error {
	var bookings []*model.Booking
	for book := range request.SeatIds {
		seatConn := grpc_client.GetSeatsClient()
		seat, err := seatConn.GetOneSeat(context.Background(), &api.GetOneSeatRequest{SeatId: request.SeatIds[book]})
		if err != nil {
			return err
		}

		userConn := grpc_client.GetUserClient()
		user, err := userConn.GetUserById(context.Background(), &api.GetUserByIdRequest{UserId: request.UserId})
		if err != nil {
			return err
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

func (s *service) FindAll(userId string, movieId int64) ([]model.Booking, error) {
	return s.repository.FindAll(userId, movieId)
}

func (s *service) IsSeatAvailable(movieID int64, seatID int64) (bool, error) {
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
