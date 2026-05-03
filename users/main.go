package main

import "context"

func main() {
	store := NewStore()
	svc := NewService(store)

	svc.CreateUser(context.Background())
}
