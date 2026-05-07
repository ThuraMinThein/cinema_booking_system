package service

import (
	"github.com/ThuraMinThein/seats/internal/model"
	"github.com/ThuraMinThein/seats/internal/repository"
)

type service struct {
	repository *repository.Repository
}

func NewService(repo *repository.Repository) *service {
	return &service{
		repository: repo,
	}
}

func (s *service) SetSeats() error {
	seats := []*model.Seat{
		{SeatNumber: "A1", ColumnNumber: "1", RowNumber: "A"},
		{SeatNumber: "A2", ColumnNumber: "2", RowNumber: "A"},
		{SeatNumber: "A3", ColumnNumber: "3", RowNumber: "A"},
		{SeatNumber: "B1", ColumnNumber: "1", RowNumber: "B"},
		{SeatNumber: "B2", ColumnNumber: "2", RowNumber: "B"},
		{SeatNumber: "B3", ColumnNumber: "3", RowNumber: "B"},
	}
	return s.repository.Create(seats)
}

func (s *service) FindAll() ([]*model.Seat, error) {
	return s.repository.FindAll()
}

func (s *service) Delete(seatID int64) error {
	return s.repository.Delete(seatID)
}
