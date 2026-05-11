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

func (s *user) CreateUser(req *model.User) (string, error) {
	hasEmail, err := s.hasEmail(req.Email)
	if err != nil {
		return "", err
	}
	if hasEmail {
		return "", errors.New("Invalid Email")
	}

	hashedPassword, err := helpers.Hash(req.Password)
	if err != nil {
		return "", err
	}
	req.Password = hashedPassword

	if err = s.repository.Create(req); err != nil {
		return "", err
	}
	return req.Id, nil
}

func (s *user) LoginUser(req *api.LoginUserRequest) (string, error) {
	user, err := s.repository.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("Invalid Credentials")
	}

	if user == nil {
		return "", errors.New("Invalid Credentials")
	}

	if err := helpers.VerifyHashed(user.Password, req.Password); err != nil {
		return "", errors.New("Invalid Credentials")
	}

	return user.Id, nil
}

func (s *user) GetUserById(id string) (*model.User, error) {
	return s.repository.GetUserById(id)
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
