package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	proto "user/api/qvbilam/user/v1"
	"user/business"
	"user/model"
	"user/utils"
)

type UserService struct {
	proto.UnimplementedUserServer
}

// Create 创建用户
func (s *UserService) Create(ctx context.Context, request *proto.UpdateRequest) (*proto.UserResponse, error) {
	fmt.Println("注册用户")

	b := business.UserBusiness{
		Gender:   request.Gender,
		Nickname: request.Nickname,
		Avatar:   request.Avatar,
	}

	entity, err := b.Create()
	if err != nil {
		return nil, err
	}

	return &proto.UserResponse{Id: entity.ID}, nil
}

func (s *UserService) Update(ctx context.Context, request *proto.UpdateRequest) (*emptypb.Empty, error) {
	fmt.Println("更新用户")

	b := business.UserBusiness{
		Id:       request.Id,
		Code:     request.Code,
		Gender:   request.Gender,
		Nickname: request.Nickname,
		Avatar:   request.Avatar,
	}

	err := b.Update()
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) Delete(ctx context.Context, request *proto.UpdateRequest) (*emptypb.Empty, error) {
	fmt.Println("注销用户")

	b := business.UserBusiness{
		Id: request.Id,
	}

	err := b.Delete()
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) CheckPassword(ctx context.Context, request *proto.CheckPasswordRequest) (*proto.CheckPasswordResponse, error) {
	fmt.Println("验证用户密码")

	res := utils.CheckPassword(request.Password, request.EnctypedPassword)
	return &proto.CheckPasswordResponse{Success: res}, nil
}

// Detail 获取用户
func (s *UserService) Detail(ctx context.Context, request *proto.GetUserRequest) (*proto.UserResponse, error) {
	fmt.Println("获取用户详情: ", request.Id)

	b := business.UserBusiness{Id: request.Id}
	entity, err := b.GetById()
	if err != nil {
		return nil, err
	}
	res := userEntityToResponse(entity)
	return res, nil
}

func (s *UserService) Search(ctx context.Context, request *proto.SearchRequest) (*proto.UsersResponse, error) {
	b := business.UserBusiness{
		Keyword:   request.Keyword,
		Sort:      request.Sort,
		IsVisible: &request.IsVisible,
		Page:      request.Page.Page,
		PerPage:   request.Page.PerPage,
	}
	fmt.Printf("搜索用户, 搜索条件: %+v\n", b)

	entities, count := b.Search()
	res := &proto.UsersResponse{}
	res.Total = count
	for _, entity := range *entities {
		res.Users = append(res.Users, userEntityToResponse(&entity))
	}

	return res, nil
	//return nil, status.Error(codes.Unimplemented, "服务不可用")
}

func (s *UserService) List(ctx context.Context, request *proto.SearchRequest) (*proto.UsersResponse, error) {
	fmt.Println("获取用户列表")
	b := business.UserBusiness{Ids: request.Id}

	entities, count := b.GetByIds()

	res := &proto.UsersResponse{}
	res.Total = count
	for _, entity := range *entities {
		res.Users = append(res.Users, userEntityToResponse(&entity))
	}
	return res, nil
}

func (s *UserService) Login(ctx context.Context, request *proto.LoginRequest) (*proto.UserResponse, error) {
	fmt.Println("登陆")
	b := business.AccountBusiness{
		Id:              0,
		UserName:        "",
		Mobile:          "",
		Email:           "",
		Password:        "",
		Ip:              "",
		LoginMethod:     "",
		AccountPlatform: nil,
	}

	b := business.UserBusiness{Mobile: request.Mobile}
	entity, err := b.GetByMobile()
	if err != nil {
		return nil, err
	}
	if check := utils.CheckPassword(request.Password, entity.Password); check == false {
		return nil, status.Error(codes.InvalidArgument, "密码错误")
	}
	return userEntityToResponse(entity), nil
}

// Logout todo 退出登陆
func (s *UserService) Logout(ctx context.Context, request *proto.GetUserRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "服务不可用")
}

func (s *UserService) Auth(ctx context.Context, request *proto.AuthRequest) (*proto.UserResponse, error) {
	b := business.UserBusiness{Id: request.Id}
	user, err := b.GetById()
	if err != nil {
		return nil, err
	}

	return &proto.UserResponse{
		Id:       user.ID,
		Code:     user.Code,
		Nickname: user.Nickname,
	}, nil
}

func userEntityToResponse(user *model.User) *proto.UserResponse {
	return &proto.UserResponse{
		Id:       user.ID,
		Code:     user.Code,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		Level:    nil,
	}
}
