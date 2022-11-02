package api

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	proto "user/api/qvbilam/user/v1"
	"user/business"
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
	if request.AccountPlatform.PlatformID != "" {
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
	ub := business.UserBusiness{}
	_, _ = ub.Create()

	return &proto.AccountResponse{Id: entity.ID}, nil
}

func (s *AccountService) Update(ctx context.Context, request *proto.UpdateAccountRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *AccountService) CheckPassword(ctx context.Context, response *proto.CheckPasswordRequest) (*proto.CheckPasswordResponse, error) {
	b := business.AccountBusiness{
		Mobile:          "",
		Email:           "",
		Password:        "",
		Ip:              "",
		AccountPlatform: nil,
	}
	return nil, nil
}
