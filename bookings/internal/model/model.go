package model

type Booking struct {
	ID         int64  `json:"id"`
	SeatID     int64  `json:"seat_id"`
	SeatNumber string `json:"seat_number"`
	UserID     int64  `json:"user_id"`
	UserName   string `json:"user_name"`
	MovieID    int64  `json:"movie_id"`
	ShowTime   string `json:"show_time"`
}
