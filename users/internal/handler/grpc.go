package handler

import (
	"context"

	"github.com/ThuraMinThein/common/api"
	types "github.com/ThuraMinThein/users/pkg/type"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	userService types.UsersService
	api.UnimplementedUserServiceServer
}

func NewGRPCUsersService(grpc *grpc.Server, userService types.UsersService) {
	grpcHandler := &grpcHandler{
		userService: userService,
	}
	api.RegisterUserServiceServer(grpc, grpcHandler)
}

func (h *grpcHandler) CreateUser(ctx context.Context, request *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	err := h.userService.CreateUser(ctx, request)
	if err != nil {
		return nil, err
	}

	return &api.CreateUserResponse{}, nil
}
func (h *grpcHandler) LoginUser(ctx context.Context, request *api.LoginUserRequest) (*api.LoginUserResponse, error) {
	err := h.userService.LoginUser(ctx)
	if err != nil {
		return nil, err
	}

	return &api.LoginUserResponse{}, nil
}
