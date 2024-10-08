package business

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
type AccountDevice struct {
	Version string
	Client  string
	Device  string
}

type AccountBusiness struct {
	Ctx             context.Context
	Id              int64
	Username        string
	Mobile          string
	Email           string
	Password        string
	Ip              string
	LoginMethod     string
	Code            string
	AccountPlatform *AccountPlatform
	AccountDevice   *AccountDevice
}

func (b *AccountBusiness) Create() (*model.Account, error) {
	parentSpan := opentracing.SpanFromContext(b.Ctx)
	spanCheckMobile := opentracing.GlobalTracer().StartSpan("checkMobileExists", opentracing.ChildOf(parentSpan.Context()))
	spanSqlBegin := opentracing.GlobalTracer().StartSpan("startSqlBegin", opentracing.ChildOf(parentSpan.Context()))

	tx := global.DB.Begin()
	m := model.Account{Mobile: b.Mobile}
	// 验证手机号

	existsMobile := b.ExistsMobile(tx)
	spanCheckMobile.Finish() // 验证手机号基数

	if existsMobile {
		tx.Rollback()
		spanSqlBegin.Finish()
		return nil, status.Errorf(codes.AlreadyExists, "手机号已存在")
	}

	// 验证邮箱
	if b.Email != "" {
		spanCheckEmail := opentracing.GlobalTracer().StartSpan("checkEmailExists", opentracing.ChildOf(parentSpan.Context()))
		m.Email = b.Email
		existsEmail := b.ExistsEmail(tx)
		spanCheckEmail.Finish() // 验证邮箱结束
		if existsEmail {
			tx.Rollback()
			spanSqlBegin.Finish() // sql结束
			return nil, status.Errorf(codes.AlreadyExists, "邮箱已存在")
		}
	}

	// 密码
	if b.Password != "" {
		spanGeneratePassword := opentracing.GlobalTracer().StartSpan("generatePassword", opentracing.ChildOf(parentSpan.Context()))
		m.Password = utils.GeneratePassword(b.Password)
		spanGeneratePassword.Finish() // 生成密码结束
	}

	m.CreatedIp = b.Ip
	// 创建账号
	spanCreateAccount := opentracing.GlobalTracer().StartSpan("createAccount", opentracing.ChildOf(parentSpan.Context()))
	err := tx.Save(&m)
	spanCreateAccount.Finish() // 创建账号结束

	if err.RowsAffected == 0 {
		tx.Rollback()
		spanSqlBegin.Finish()
		return nil, status.Errorf(codes.Internal, "创建账号失败")
	}

	spanCreateAccountLog := opentracing.GlobalTracer().StartSpan("createAccountLog", opentracing.ChildOf(parentSpan.Context()))
	b.Id = m.ID
	b.accountLog(tx, enum.AccountTypeSignin)
	spanCreateAccountLog.Finish()

	// 创建默认用户信息
	spanCreateAccountUser := opentracing.GlobalTracer().StartSpan("createAccountUser", opentracing.ChildOf(parentSpan.Context()))
	spanCreateAccountUser.LogFields(log.Int64("accountId", m.ID))
	ub := UserBusiness{AccountId: m.ID}
	user, cErr := ub.Create(tx)
	spanCreateAccountUser.LogFields(
		log.Int64("user.id", user.ID),
		log.Int64("user.code", user.Code))
	spanCreateAccountUser.Finish()

	if cErr != nil {
		tx.Rollback()
		spanSqlBegin.Finish()
		return nil, cErr
	}

	tx.Commit()
	spanSqlBegin.Finish()
	return &m, nil
}

func (b *AccountBusiness) Update() (*model.Account, error) {
	entity := model.Account{}
	tx := global.DB.Begin()
	if res := tx.Where(&model.Account{IDModel: model.IDModel{ID: b.Id}}).First(&entity); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "账号不存在")
	}

	if (b.Username != "") && (b.Username != entity.Username) {
		if b.ExistsUserName(tx) {
			return nil, status.Errorf(codes.AlreadyExists, "账号已存在")
		}
		entity.Username = b.Username
	}
	if (b.Mobile != "") && (b.Mobile != entity.Mobile) {
		if b.ExistsMobile(tx) {
			return nil, status.Errorf(codes.AlreadyExists, "手机号已注册")
		}
		entity.Mobile = b.Mobile
	}
	if (b.Email != "") && (b.Email != entity.Email) {
		if b.ExistsEmail(tx) {
			return nil, status.Errorf(codes.AlreadyExists, "邮箱已注册")
		}
		entity.Email = b.Email
	}
	if b.Password != "" {
		entity.Password = utils.GeneratePassword(b.Password)
	}

	tx.Save(&entity)
	tx.Commit()

	return &entity, nil
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
		if b.Username == "" {
			return nil, status.Errorf(codes.InvalidArgument, "请输入账号")
		}
		condition.Username = b.Username
	default:
		return nil, status.Errorf(codes.InvalidArgument, "非法请求")
	}

	tx := global.DB.Begin()
	if res := tx.Where(condition).Select("id, password").First(&entity); res.RowsAffected == 0 {
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
		accountEntity.Username = b.AccountPlatform.Type + u.PlatformId
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
	// 新增账户日志
	b.accountLog(tx, enum.AccountTypeLogin)
	return nil
}

func (b *AccountBusiness) BlindPlatform(code string) error {
	if b.Id == 0 {
		return status.Errorf(codes.InvalidArgument, "参数错误")
	}

	if b.AccountPlatform.Type == "" || code == "" {
		return status.Errorf(codes.Internal, "第三方登陆参数错误")
	}

	tx := global.DB.Begin()

	// 验证是否绑定
	if b.ExistsPlatform(tx) {
		return status.Errorf(codes.AlreadyExists, "已绑定当前平台登录账号")
	}
	//b.existsPlatform(tx)
	//entity := model.AccountPlatform{}
	//if res := global.DB.Where(&model.AccountPlatform{
	//	AccountID: b.Id,
	//	Type:      b.AccountPlatform.Type,
	//}).First(&entity); res.RowsAffected == 1 {
	//	return status.Errorf(codes.AlreadyExists, "已绑定当前平台登录账号")
	//}

	// 登陆
	u, err := b.getPlatformUser(code)
	if err != nil {
		return err
	}
	if u == nil {
		return status.Errorf(codes.Internal, "第三方登陆失败")
	}

	// 绑定
	if res := tx.Save(&model.AccountPlatform{
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

	if res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(&model.Account{Email: b.Email}).
		Select("email").
		First(&model.Account{}); res.RowsAffected == 0 {
		return false
	}
	return true
}

func (b *AccountBusiness) ExistsPlatform(tx *gorm.DB) bool {
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

func (b *AccountBusiness) ExistsUserName(tx *gorm.DB) bool {
	if tx == nil {
		tx = global.DB
	}
	if b.Username == "" {
		return true
	}

	if res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where(&model.Account{Username: b.Username}).
		Select("username").
		First(&model.Account{}); res.RowsAffected == 0 {
		return false
	}
	return true
}

func (b *AccountBusiness) accountLog(tx *gorm.DB, accountType string) bool {
	client, version, device := "", "", ""
	if b.AccountDevice != nil {
		client = b.AccountDevice.Client
		version = b.AccountDevice.Version
		device = b.AccountDevice.Device
	}

	entity := model.AccountLog{
		AccountId: b.Id,
		Type:      accountType,
		Method:    b.LoginMethod,
		Client:    client,
		Version:   version,
		Device:    device,
		Ip:        b.Ip,
	}

	if res := tx.Save(&entity); res.RowsAffected == 0 {
		return false
	}

	return true
}
