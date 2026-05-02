package users

import "context"

func main() {
	store := NewStore()
	svc := NewService(store)

	svc.CreateUser(context.Background())
}
