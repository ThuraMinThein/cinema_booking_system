package main

import "context"

type user struct {
	store UsersStore
}

func NewService(store UsersStore) *user {
	return &user{store: store}
}

func (s *user) CreateUser(ctx context.Context) error {
	return s.store.Create(ctx)
}
