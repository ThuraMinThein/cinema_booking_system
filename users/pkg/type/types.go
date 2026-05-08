package types

import (
	"github.com/ThuraMinThein/common/api"
	"github.com/ThuraMinThein/users/internal/model"
)

type UsersService interface {
	CreateUser(*model.User) error
	LoginUser(*api.LoginUserRequest) error
	GetUserById(string) (*model.User, error)
}

type UsersRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	GetUserById(string) (*model.User, error)
}
