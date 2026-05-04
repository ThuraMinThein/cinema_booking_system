package service

import (
	"context"

	"github.com/ThuraMinThein/common/api"
	types "github.com/ThuraMinThein/users/pkg/type"
)

type user struct {
	store types.UsersRepository
}

func NewService(store types.UsersRepository) *user {
	return &user{store: store}
}

func (s *user) CreateUser(ctx context.Context, user *api.CreateUserRequest) error {
	return s.store.Create(ctx)
}

func (s *user) LoginUser(ctx context.Context) error {
	return s.LoginUser(ctx)
}
