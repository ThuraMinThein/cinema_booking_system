package repository

import (
	"errors"

	"github.com/ThuraMinThein/users/internal/model"
	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *store {
	return &store{
		db: db,
	}
}

func (s *store) Create(req *model.User) error {
	return s.db.Create(req).Error
}

func (s *store) FindByEmail(email string) (*model.User, error) {
	var user *model.User
	result := s.db.First(&user, "email = ?", email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return user, nil
}

func (s *store) GetUserById(id string) (*model.User, error) {
	var user *model.User
	result := s.db.First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return user, nil
}
