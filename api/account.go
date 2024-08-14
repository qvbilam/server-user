package api

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		Ctx:      ctx,
		Mobile:   request.Mobile,
		Email:    request.Email,
		Password: request.Password,
		Ip:       request.Ip,
		AccountDevice: &business.AccountDevice{
			Version: request.Device.Version,
			Client:  request.Device.Client,
			Device:  request.Device.Device,
		},
	}
	// 获取span
	parentSpan := opentracing.SpanFromContext(ctx)
	BusSpan := opentracing.GlobalTracer().StartSpan("startAccountBusiness", opentracing.ChildOf(parentSpan.Context()))

	// 创建账号
	entity, err := b.Create()

	BusSpan.Finish()
	if err != nil {
		return nil, err
	}

	return &proto.AccountResponse{Id: entity.ID}, nil
}

// Update 更新账号
func (s *AccountService) Update(ctx context.Context, request *proto.UpdateAccountRequest) (*emptypb.Empty, error) {
	if request.Id == 0 {
		if request.UserId == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "参数错误")
		}
		ub := business.UserBusiness{Id: request.UserId}
		user, err := ub.GetDetail()
		if err != nil {
			return nil, err
		}
		request.Id = user.AccountId
	}

	b := business.AccountBusiness{
		Id:       request.Id,
		Username: request.Username,
		Mobile:   request.Mobile,
		Email:    request.Email,
		Password: request.Password,
	}

	if _, err := b.Update(); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// LoginPassword 密码登陆
func (s *AccountService) LoginPassword(ctx context.Context, request *proto.LoginPasswordRequest) (*proto.AccountResponse, error) {
	b := business.AccountBusiness{
		Username:    request.Username,
		Mobile:      request.Mobile,
		Email:       request.Email,
		Password:    request.Password,
		Ip:          request.Ip,
		LoginMethod: request.Method,
		AccountDevice: &business.AccountDevice{
			Version: request.Device.Version,
			Client:  request.Device.Client,
			Device:  request.Device.Device,
		},
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
		AccountDevice: &business.AccountDevice{
			Version: request.Device.Version,
			Client:  request.Device.Client,
			Device:  request.Device.Device,
		},
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
		AccountDevice: &business.AccountDevice{
			Version: "",
			Client:  "",
			Device:  "",
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
		Username: entity.Username,
		Mobile:   entity.Mobile,
		Email:    entity.Email,
		User:     &userRes,
	}, nil
}
