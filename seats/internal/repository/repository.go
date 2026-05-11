package repository

import (
	"github.com/ThuraMinThein/seats/internal/model"
	"gorm.io/gorm"
)

type Repository struct {
	database *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		database: db,
	}
}

func (r *Repository) Create(seats []*model.Seat) error {
	return r.database.Create(seats).Error
}

func (r *Repository) FindAll() ([]*model.Seat, error) {
	var seat []*model.Seat
	err := r.database.Order("seat_number ASC").Find(&seat).Error
	return seat, err
}

func (r *Repository) Delete(seatID int64) error {
	return r.database.Delete(&model.Seat{}, seatID).Error
}
