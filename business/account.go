package business

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
	"user/enum"
	"user/global"
	"user/model"
	"user/utils"
)

type AccountPlatform struct {
	PlatformID string
	Type       string
}

type AccountBusiness struct {
	Id              int64
	UserName        string
	Mobile          string
	Email           string
	Password        string
	Ip              string
	LoginMethod     string
	Code            string
	AccountPlatform *AccountPlatform
}

func (b *AccountBusiness) Create() (*model.Account, error) {
	tx := global.DB.Begin()
	m := model.Account{Mobile: b.Mobile}
	// 验证手机号
	if b.ExistsMobile(tx) {
		tx.Rollback()
		return nil, status.Errorf(codes.AlreadyExists, "手机号已存在")
	}

	// 验证邮箱
	if b.Email != "" {
		m.Email = b.Email
		if b.ExistsEmail(tx) {
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

	tx.Commit()
	return &m, nil
}

func (b *AccountBusiness) Update() (*model.Account, error) {
	return nil, nil
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
		tx.Rollback()
		return nil, status.Errorf(codes.InvalidArgument, "密码错误")
	}
	b.Id = entity.ID

	// 更新登陆后相关状态
	if err := b.login(tx); err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "更新失败")
	}

	tx.Commit()
	return &entity, nil
}

func (b *AccountBusiness) LoginMobileCode() (*model.Account, error) {
	entity := model.Account{}
	if b.Mobile == "" {
		return nil, status.Errorf(codes.InvalidArgument, "请输入手机号")
	}

	tx := global.DB.Begin()
	if res := tx.Where(&model.Account{Mobile: b.Mobile}).First(&entity); res.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.NotFound, "账号不存在")
	}
	b.Id = entity.ID
	b.LoginMethod = enum.LoginMethodSms

	// todo 验证短信

	// 更新登陆后相关状态
	if err := b.login(tx); err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "更新失败")
	}

	tx.Commit()
	return &entity, nil
}

func (b *AccountBusiness) getPlatformUser(code string) (*OAuthUserResponse, error) {
	var u *OAuthUserResponse
	switch b.AccountPlatform.Type {
	case enum.LoginMethodPlatformQQ:
		oauth := OAuthQQBusiness{
			AppId:     global.ServerConfig.OauthQQConfig.AppId,
			AppSecret: global.ServerConfig.OauthQQConfig.AppSecret,
			Uri:       global.ServerConfig.OauthQQConfig.Uri,
		}
		u = oauth.User(code)
	case enum.LoginMethodPlatformWGitHub:
		oauth := OAuthGitHubBusiness{
			AppId:      global.ServerConfig.OauthGithubConfig.AppId,
			AppSecrete: global.ServerConfig.OauthGithubConfig.AppSecret,
			Uri:        global.ServerConfig.OauthGithubConfig.Uri,
		}
		u = oauth.User(code)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "非法请求")
	}

	return u, nil
}

func (b *AccountBusiness) LoginPlatform(code string) (*model.Account, error) {
	if b.AccountPlatform.Type == "" || code == "" {
		return nil, status.Errorf(codes.Internal, "第三方登陆参数错误")
	}
	u, err := b.getPlatformUser(code)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, status.Errorf(codes.Internal, "第三方登陆失败")
	}

	tx := global.DB.Begin()
	accountEntity := model.Account{}
	entity := model.AccountPlatform{}
	if res := tx.Where(&model.AccountPlatform{
		PlatformID: u.PlatformId,
		Type:       b.AccountPlatform.Type,
	}).First(&entity); res.RowsAffected == 0 {
		// 创建 account
		accountEntity.UserName = b.AccountPlatform.Type + u.PlatformId
		if res := tx.Save(&accountEntity); res.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "创建账号失败")
		}

		// 创建第三方账号
		entity.AccountID = accountEntity.ID
		entity.PlatformID = u.PlatformId
		entity.Type = b.AccountPlatform.Type
		if res := tx.Save(&entity); res.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "创建账号失败")
		}

		// 创建用户
		ub := UserBusiness{
			AccountId: accountEntity.ID,
			Gender:    u.Gender,
			Nickname:  u.Nickname,
			Avatar:    u.Avatar,
		}
		if _, err := ub.Create(tx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	b.Id = entity.AccountID
	b.LoginMethod = b.AccountPlatform.Type

	// 更新登陆后相关状态
	if err := b.login(tx); err != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "更新失败")
	}

	accountEntity.ID = entity.AccountID
	tx.Commit()
	return &accountEntity, nil

}

func (b *AccountBusiness) login(tx *gorm.DB) error {
	if tx == nil {
		tx = global.DB
	}

	updates := map[string]interface{}{
		"login_count":   gorm.Expr("login_count + ?", 1),
		"last_login_ip": b.Ip,
		"last_login_at": time.Now().Format("2006-01-02 15:04.05"),
	}
	// 更新登陆信息
	if res := tx.Model(model.Account{}).Where(model.Account{IDModel: model.IDModel{ID: b.Id}}).Updates(updates); res.Error != nil || res.RowsAffected == 0 {
		return status.Errorf(codes.Internal, "更新失败")
	}
	// todo 新增账户日志

	return nil
}

func (b *AccountBusiness) BlindPlatform(code string) error {
	if b.Id == 0 {
		return status.Errorf(codes.InvalidArgument, "参数错误")
	}

	if b.AccountPlatform.Type == "" || code == "" {
		return status.Errorf(codes.Internal, "第三方登陆参数错误")
	}

	// 验证是否绑定
	entity := model.AccountPlatform{}
	if res := global.DB.Where(&model.AccountPlatform{
		AccountID: b.Id,
		Type:      b.AccountPlatform.Type,
	}).First(&entity); res.RowsAffected == 1 {
		return status.Errorf(codes.AlreadyExists, "已绑定当前平台登录账号")
	}

	// 登陆
	u, err := b.getPlatformUser(code)
	if err != nil {
		return err
	}
	if u == nil {
		return status.Errorf(codes.Internal, "第三方登陆失败")
	}

	// 绑定
	if res := global.DB.Save(&model.AccountPlatform{
		AccountID:  b.Id,
		PlatformID: u.PlatformId,
		Type:       u.Type,
	}); res.RowsAffected == 0 {
		return status.Errorf(codes.Internal, "绑定失败")
	}

	return nil
}

func (b *AccountBusiness) UnBlindPlatform() error {
	accountEntity := model.Account{}
	global.DB.Where(&model.Account{IDModel: model.IDModel{ID: b.Id}}).First(&accountEntity)
	if accountEntity.Mobile == "" || accountEntity.Email == "" {
		return status.Errorf(codes.InvalidArgument, "请先绑定手机号或邮箱")
	}

	entity := model.AccountPlatform{}
	if res := global.DB.Where(&model.AccountPlatform{
		AccountID: b.Id,
		Type:      b.AccountPlatform.Type,
	}).First(&entity); res.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "未绑定平台账号")
	}

	if entity.AccountID != 0 {
		entity.OriginAccountID = entity.AccountID
	}
	entity.AccountID = 0

	if res := global.DB.Save(&entity); res.RowsAffected == 0 {
		return status.Errorf(codes.Internal, "解绑失败")
	}

	return nil
}

func (b *AccountBusiness) ExistsMobile(tx *gorm.DB) bool {
	if tx == nil {
		tx = global.DB
	}
	if b.Mobile == "" {
		return true
	}
	if res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Account{Mobile: b.Mobile}).Select("mobile").First(&model.Account{}); res.RowsAffected == 0 {
		return false
	}
	return true
}

func (b *AccountBusiness) ExistsEmail(tx *gorm.DB) bool {
	if tx == nil {
		tx = global.DB
	}
	if b.Email == "" {
		return true
	}

	if res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Account{Email: b.Email}).Select("email").First(&model.Account{}); res.RowsAffected == 0 {
		return false
	}
	return true
}

func (b *AccountBusiness) existsPlatform(tx *gorm.DB) bool {
	if tx == nil {
		tx = global.DB
	}
	if b.AccountPlatform == nil {
		return true
	}

	if res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(AccountPlatform{
		PlatformID: b.AccountPlatform.PlatformID,
		Type:       b.AccountPlatform.Type,
	}).Select("id").First(&AccountPlatform{}); res.RowsAffected == 0 {
		return false
	}
	return true
}
