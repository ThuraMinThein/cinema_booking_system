package main

import "context"

type UsersService interface {
	CreateUser(context.Context) error
}

type UsersStore interface {
	Create(context.Context) error
}
