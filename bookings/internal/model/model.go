package model

import "time"

type Booking struct {
	ID         int64     `json:"id"`
	SeatID     int64     `json:"seat_id"`
	SeatNumber string    `json:"seat_number"`
	UserID     string    `json:"user_id"`
	UserName   string    `json:"user_name"`
	MovieID    int64     `json:"movie_id"`
	Showtime   time.Time `json:"showtime"`
	Status     string    `json:"status" gorm:"status" default:"Booked"`
}
