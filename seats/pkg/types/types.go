package types

import "github.com/ThuraMinThein/seats/internal/model"

type SeatService interface {
	SetSeats() error
	FindAll() ([]*model.Seat, error)
	FindOne(seatId int64) (*model.Seat, error)
	Delete(seatID int64) error
}

type SeatRepository interface {
	Create(seats []*model.Seat) error
	FindAll() ([]*model.Seat, error)
	FindOne(seatId int64) (*model.Seat, error)
	Delete(seatID int64) error
}
