package types

import (
	"context"

	"github.com/ThuraMinThein/common/api"
)

type UsersService interface {
	CreateUser(context.Context, *api.CreateUserRequest) error
	LoginUser(context.Context) error
}

type UsersRepository interface {
	Create(context.Context) error
}
