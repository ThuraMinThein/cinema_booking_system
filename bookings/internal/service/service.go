package service

import (
	"github.com/ThuraMinThein/bookings/internal/model"
	"github.com/ThuraMinThein/bookings/internal/repository"
)

type service struct {
	repository *repository.Repository
}

func NewService(repo *repository.Repository) *service {
	return &service{
		repository: repo,
	}
}

func (s *service) Create(booking *model.Booking) error {
	return s.repository.Create(booking)
}

func (s *service) FindAll(movieID int64) ([]model.Booking, error) {
	return s.repository.FindAll(movieID)
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
