package model

import "time"

type Seat struct {
	ID           int64
	SeatNumber   string
	ColumnNumber string
	RowNumber    string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
