package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	proto "user/api/pb"
	"user/business"
)

type UserService struct {
	proto.UnimplementedUserServer
}

// Create 创建用户
func (s *UserService) Create(ctx context.Context, request *proto.SignInRequest) (*proto.UserResponse, error) {
	fmt.Println("创建用户")
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

func (s *UserService) Update(ctx context.Context, request *proto.UpdateRequest) (*proto.UserResponse, error) {
	fmt.Println("更新用户")
	return nil, status.Error(codes.Unimplemented, "服务不可用")
}

func (s *UserService) CheckPassword(ctx context.Context, request *proto.CheckPasswordRequest) (*proto.CheckPasswordResponse, error) {
	fmt.Println("验证用户密码")
	return nil, status.Error(codes.Unimplemented, "服务不可用")
}

// Detail 获取用户
func (s *UserService) Detail(ctx context.Context, request *proto.GetUserRequest) (*proto.UserResponse, error) {
	fmt.Println("获取用户详情: ", request.Id)
	b := business.UserBusiness{Id: request.Id}
	entity, err := b.GetUserById()
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *UserService) List(ctx context.Context, request *proto.SearchRequest) (*proto.UsersResponse, error) {
	fmt.Println("获取用户列表")
	return nil, status.Error(codes.Unimplemented, "服务不可用")
}

func (s *UserService) ListByIds(ctx context.Context, request *proto.ListByIdsRequest) (*proto.UsersResponse, error) {
	fmt.Println("批量获取用户:", request.Ids)
	b := business.UserBusiness{Ids: request.Ids}
	res, err := b.GetUserByIds()
	if err != nil {
		return nil, err
	}
	return res, nil
}
