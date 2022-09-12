package business

import (
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	proto "user/api/qvbilam/user/v1"
	"user/global"
	"user/model"
)

type UserBusiness struct {
	Id       int64
	Code     int64
	Mobile   string
	Gender   string
	Nickname string
	Password string
	Ids      []int64
	Keyword  string
}

// ExistsMobile 验证手机号
func (b *UserBusiness) ExistsMobile() bool {
	entity := model.User{}
	res := global.DB.Where(&model.User{
		Mobile: b.Mobile,
	}).Select("mobile").First(&entity)
	if res.RowsAffected == 0 {
		return false
	}
	return true
}

// CreateUser 创建用户
func (b *UserBusiness) CreateUser() (*proto.UserResponse, error) {
	if b.ExistsMobile() {
		return nil, status.Errorf(codes.AlreadyExists, "手机号已注册")
	}

	Entity := model.User{
		Code:     generateUserCode(),
		Mobile:   b.Mobile,
		Password: b.Password,
		Nickname: b.Nickname,
		Gender:   b.Gender,
		Level:    0,
		Visible: model.Visible{
			IsVisible: true,
		},
	}
	if res := global.DB.Save(&Entity); res.RowsAffected == 0 {
		zap.S().Errorf("创建用户失败: %s", res.Error)
		return nil, status.Errorf(codes.Internal, "创建失败")
	}

	return &proto.UserResponse{Id: Entity.ID}, nil
}

// GetUserById 获取用户
func (b *UserBusiness) GetUserById() (*proto.UserResponse, error) {
	entity := model.User{}
	if res := global.DB.First(&entity, b.Id); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	return &proto.UserResponse{
		Id:       entity.ID,
		Code:     entity.Code,
		Nickname: entity.Nickname,
		Avatar:   entity.Avatar,
		Gender:   entity.Gender,
		Level:    nil,
	}, nil
}

func (b *UserBusiness) GetUserList() (*proto.UsersResponse, error) {
	// todo elasticsearch 搜索用户
	return b.GetUserByIds()
}

func (b *UserBusiness) GetUserByIds() (*proto.UsersResponse, error) {
	var entity []model.User

	res := global.DB.Find(&entity, b.Ids)

	response := proto.UsersResponse{}
	response.Total = res.RowsAffected
	for _, user := range entity {
		response.Users = append(response.Users, &proto.UserResponse{
			Id:       user.ID,
			Code:     user.Code,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Gender:   user.Gender,
			Level:    nil,
		})
	}

	return &response, nil
}

func generateUserCode() int64 {
	return 534511019
}
