package business

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"
	"user/enum"
	"user/global"
	"user/model"
	"user/utils"
)

type AccountPlatform struct {
	PlatformID    string
	PlatformToken string
	Type          string
}

type AccountBusiness struct {
	Id              int64
	UserName        string
	Mobile          string
	Email           string
	Password        string
	Ip              string
	LoginMethod     string
	AccountPlatform *AccountPlatform
}

func (b *AccountBusiness) Create() (*model.Account, error) {
	tx := global.DB.Begin()
	m := model.Account{Mobile: b.Mobile}
	// 验证手机号
	if m.ExistsMobile(tx) {
		tx.Rollback()
		return nil, status.Errorf(codes.AlreadyExists, "手机号已存在")
	}

	// 验证邮箱
	if b.Email != "" {
		m.Email = b.Email
		if m.ExistsEmail(tx) {
			tx.Rollback()
			return nil, status.Errorf(codes.AlreadyExists, "邮箱已存在")
		}
	}

	// 密码
	if b.Password != "" {
		m.Password = utils.GeneratePassword(b.Password)
	}

	m.CreatedIp = b.Ip
	// 创建账号
	if err := tx.Save(&m); err.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建账号失败")
	}

	// 创建第三方账号
	if b.AccountPlatform != nil {
		apm := model.AccountPlatform{PlatformID: b.AccountPlatform.PlatformID, Type: b.AccountPlatform.Type}
		if apm.Exists(tx) { // 创建
			apm.AccountID = m.ID
			tx.Save(&apm)
		} else { // 绑定
			tx.Where(apm).Updates(model.AccountPlatform{AccountID: m.ID})
		}
	}

	tx.Commit()
	return &m, nil
}

func (b *AccountBusiness) LoginPassword() (*model.Account, error) {
	entity, condition := model.Account{}, model.Account{}
	switch b.LoginMethod {
	case enum.LoginMethodMobile:
		if b.Mobile == "" {
			return nil, status.Errorf(codes.InvalidArgument, "请输入账号")
		}
		condition.Mobile = b.Mobile
	case enum.LoginMethodEmail:
		if b.Email == "" {
			return nil, status.Errorf(codes.InvalidArgument, "请输入账号")
		}
		condition.Email = b.Email
	case enum.LoginMethodUserName:
		if b.UserName == "" {
			return nil, status.Errorf(codes.InvalidArgument, "请输入账号")
		}
		condition.UserName = b.UserName
	default:
		return nil, status.Errorf(codes.InvalidArgument, "非法请求")
	}

	tx := global.DB.Begin()
	if res := tx.Where(condition).Preload("id, password").First(&entity); res.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "账号不存在")
	}

	// 验证密码
	if s := utils.CheckPassword(b.Password, entity.Password); s == false {
		return nil, status.Errorf(codes.InvalidArgument, "密码错误")
	}

	// 登陆
	updates := map[string]interface{}{
		"login_count":   gorm.Expr("login_count + ?", 1),
		"last_login_ip": b.Ip,
		"last_login_at": time.Now().Format("2006-01-02 15:04.05"),
	}
	if res := tx.Model(model.Account{}).Where(model.Account{IDModel: model.IDModel{ID: entity.ID}}).Updates(updates); res.Error != nil || res.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "更新失败")
	}

	tx.Commit()
	return &entity, nil
}
