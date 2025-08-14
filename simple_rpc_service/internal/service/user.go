package service

import (
	"context"
	"simple_rpc_svc/internal/proto"
)

type UserService struct {
	proto.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	userInfo := &proto.GetUserResponse{
		Id:   1,
		Name: "jack",
		Age:  20,
	}
	return userInfo, nil

}
