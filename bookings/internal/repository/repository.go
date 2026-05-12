package repository

import (
	"github.com/ThuraMinThein/bookings/internal/model"
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

func (r *Repository) Create(booking []*model.Booking) error {
	return r.database.Create(booking).Error
}

func (r *Repository) FindAll(userId string, movieId int64) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.database.Where("user_id = ? AND movie_id = ?", userId, movieId).Find(&bookings).Error
	return bookings, err
}

func (r *Repository) FindByMovieAndSeatID(movieID int64, seatID int64) (*model.Booking, error) {
	var booking *model.Booking
	result := r.database.Where("movie_id = ? AND seat_id = ?", movieID, seatID).Find(&booking)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}
	return booking, nil
}

func (r *Repository) Update(booking *model.Booking) error {
	return r.database.Save(booking).Error
}

func (r *Repository) Cancel(bookingID int64) error {
	return r.database.Delete(&model.Booking{}, bookingID).Error
}
