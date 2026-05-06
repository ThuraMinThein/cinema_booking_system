package service

import (
	"errors"

	"github.com/ThuraMinThein/common/api"
	"github.com/ThuraMinThein/users/internal/model"
	"github.com/ThuraMinThein/users/pkg/helpers"
	types "github.com/ThuraMinThein/users/pkg/type"
)

type user struct {
	repository types.UsersRepository
}

func NewService(repository types.UsersRepository) *user {
	return &user{repository: repository}
}

func (s *user) CreateUser(req *model.User) error {
	hasEmail, err := s.hasEmail(req.Email)
	if err != nil {
		return err
	}
	if hasEmail {
		return errors.New("Invalid Email")
	}

	hashedPassword, err := helpers.Hash(req.Password)
	if err != nil {
		return err
	}
	req.Password = hashedPassword

	return s.repository.Create(req)
}

func (s *user) LoginUser(req *api.LoginUserRequest) error {
	user, err := s.repository.FindByEmail(req.Email)
	if err != nil {
		return err
	}

	if user.Password != req.Password {
		return err
	}
	return nil
}

// utils
func (s *user) hasEmail(email string) (bool, error) {
	if email != "" {
		user, err := s.repository.FindByEmail(email)
		if err != nil {
			return false, err
		}
		if user != nil {
			return true, nil
		}
	}

	return false, nil
}
