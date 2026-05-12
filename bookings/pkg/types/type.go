package types

import (
	"github.com/ThuraMinThein/bookings/internal/model"
	"github.com/ThuraMinThein/common/api"
)

type BookingService interface {
	Create(*api.CreateRequest) error
	HoldBooking(*api.HoldBookingRequest) error
	FindAll(string, int64) ([]model.Booking, error)
	IsSeatAvailable(int64, int64) (bool, string, error)
	Update(*model.Booking) error
	Cancel(int64) error
}

type BookingRepository interface {
	Create(*model.Booking) error
	FindAll(int64) ([]model.Booking, error)
	FindByMovieAndSeatID(int64, int64) (*model.Booking, error)
	Update(*model.Booking) error
	Cancel(int64) error
}
