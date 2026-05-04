package repository

import "context"

type store struct {
}

func NewRepository() *store {
	return &store{}
}

func (s *store) Create(context.Context) error {
	return nil
}
