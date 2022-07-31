package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	proto "user/api/v1"
	"user/business"
)

type UserService struct {
	proto.UnimplementedUserServer
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, request *proto.SignInRequest) (*proto.UserResponse, error) {
	b := business.UserBusiness{
		Mobile:   request.Mobile,
		Gender:   request.Gender,
		Nickname: request.Nickname,
		Password: request.Password,
	}
	entity, err := b.CreateUser()
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *UserService) UpdateUser(ctx context.Context, request *proto.UpdateRequest) (*proto.UserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "服务不可用")
}

func (s *UserService) CheckPassword(ctx context.Context, request *proto.CheckPasswordRequest) (*proto.CheckPasswordResponse, error) {
	return nil, status.Error(codes.Unimplemented, "服务不可用")
}

// GetUser 获取用户
func (s *UserService) GetUser(ctx context.Context, request *proto.GetUserRequest) (*proto.UserResponse, error) {
	b := business.UserBusiness{Id: request.Id}
	entity, err := b.GetUserById()
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *UserService) GetUserByCode(ctx context.Context, request *proto.GetUserRequest) (*proto.UserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "服务不可用")
}
