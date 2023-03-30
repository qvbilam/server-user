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

func (s *AccountService) Create(ctx context.Context, request *proto.UpdateAccountRequest) (*proto.AccountResponse, error) {
	b := business.AccountBusiness{
		Mobile:   request.Mobile,
		Email:    request.Email,
		Password: request.Password,
		Ip:       request.Ip,
	}
	if request.AccountPlatform != nil && request.AccountPlatform.PlatformID != "" {
		b.AccountPlatform = &business.AccountPlatform{
			PlatformID:    request.AccountPlatform.PlatformID,
			PlatformToken: request.AccountPlatform.PlatformToken,
			Type:          request.AccountPlatform.Type,
		}
	}
	// 创建账号
	entity, err := b.Create()
	if err != nil {
		return nil, err
	}
	// 创建默认用户信息
	ub := business.UserBusiness{AccountId: entity.ID}
	if _, err := ub.Create(); err != nil {
		return nil, err
	}

	return &proto.AccountResponse{Id: entity.ID}, nil
}

func (s *AccountService) Update(ctx context.Context, request *proto.UpdateAccountRequest) (*emptypb.Empty, error) {
	return nil, nil
}

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
