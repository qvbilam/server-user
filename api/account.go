package api

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	proto "user/api/qvbilam/user/v1"
	"user/business"
	"user/enum"
	"user/model"
)

type AccountService struct {
	proto.UnimplementedAccountServer
}

// Create 创建账号
func (s *AccountService) Create(ctx context.Context, request *proto.UpdateAccountRequest) (*proto.AccountResponse, error) {
	b := business.AccountBusiness{
		Mobile:   request.Mobile,
		Email:    request.Email,
		Password: request.Password,
		Ip:       request.Ip,
	}
	// 创建账号
	entity, err := b.Create()
	if err != nil {
		return nil, err
	}
	// 创建默认用户信息
	ub := business.UserBusiness{AccountId: entity.ID}
	if _, err := ub.Create(nil); err != nil {
		return nil, err
	}

	return &proto.AccountResponse{Id: entity.ID}, nil
}

// Update 更新账号
func (s *AccountService) Update(ctx context.Context, request *proto.UpdateAccountRequest) (*emptypb.Empty, error) {
	return nil, nil
}

// LoginPassword 密码登陆
func (s *AccountService) LoginPassword(ctx context.Context, request *proto.LoginPasswordRequest) (*proto.AccountResponse, error) {
	b := business.AccountBusiness{
		UserName:    request.Username,
		Mobile:      request.Mobile,
		Email:       request.Email,
		Password:    request.Password,
		Ip:          request.Ip,
		LoginMethod: request.Method,
	}
	// 账号信息
	entity, err := b.LoginPassword()
	if err != nil {
		return nil, err
	}

	return s.loginResponse(entity)
}

// LoginSmsCode 短信登陆
func (s *AccountService) LoginSmsCode(ctx context.Context, request *proto.LoginMobileRequest) (*proto.AccountResponse, error) {
	b := business.AccountBusiness{
		Mobile:      request.Mobile,
		Ip:          request.Ip,
		LoginMethod: enum.LoginMethodSms,
	}
	entity, err := b.LoginMobileCode()
	if err != nil {
		return nil, err
	}

	return s.loginResponse(entity)
}

// LoginPlatform 平台登陆
func (s *AccountService) LoginPlatform(ctx context.Context, request *proto.LoginPlatformRequest) (*proto.AccountResponse, error) {
	b := business.AccountBusiness{
		AccountPlatform: &business.AccountPlatform{
			Type: request.Type,
		},
	}
	entity, err := b.LoginPlatform(request.Code)
	if err != nil {
		return nil, err
	}
	return s.loginResponse(entity)
}

// BindPlatform 绑定平台账号
func (s *AccountService) BindPlatform(ctx context.Context, request *proto.BindPlatformRequest) (*emptypb.Empty, error) {
	ub := business.UserBusiness{Id: request.UserId}
	user, err := ub.GetDetail()
	if err != nil {
		return nil, err
	}
	b := business.AccountBusiness{Id: user.AccountId}
	if err := b.BlindPlatform(request.Code); err != nil {
		return nil, err
	}
	return nil, nil
}

// UnbindPlatform 取消绑定账号
func (s *AccountService) UnbindPlatform(ctx context.Context, request *proto.BindPlatformRequest) (*emptypb.Empty, error) {
	ub := business.UserBusiness{Id: request.UserId}
	user, err := ub.GetDetail()
	if err != nil {
		return nil, err
	}
	b := business.AccountBusiness{Id: user.AccountId}
	if err := b.UnBlindPlatform(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *AccountService) loginResponse(entity *model.Account) (*proto.AccountResponse, error) {
	// 获取用户信息
	ub := business.UserBusiness{AccountId: entity.ID}
	userEntity, _ := ub.GetDetail()
	userRes := proto.UserResponse{}
	if userEntity != nil {
		userRes.Id = userEntity.ID
		userRes.Code = userEntity.Code
		userRes.Nickname = userEntity.Nickname
		userRes.Avatar = userEntity.Avatar
		userRes.Gender = userEntity.Gender
		userRes.Level = nil
	}

	return &proto.AccountResponse{
		Token:    GenerateUserToken(&userRes),
		UserName: entity.UserName,
		Mobile:   entity.Mobile,
		Email:    entity.Email,
		User:     &userRes,
	}, nil
}
