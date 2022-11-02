package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Account struct {
	IDModel
	UserName    string `gorm:"type:varchar(64);not null default '';comment:登陆账号;index"`
	Email       string `gorm:"type:varchar(64);not null default '';comment:邮箱;index"`
	Mobile      string `gorm:"type:varchar(11);not null default '';comment:手机号;index"`
	Password    string `gorm:"type:varchar(128);not null default '';comment:登陆密码"`
	LoginCount  int64  `gorm:"type:int(64);not null default 0;comment:登陆次数"`
	CreatedIp   string `gorm:"type:varchar(32);not null default '';comment:注册IP"`
	LastLoginIp string `gorm:"type:varchar(32);not null default '';comment:上次登陆IP"`
	LastLoginAt *time.Time
	DateModel
	DeletedModel
}

func (m *Account) Login(tx *gorm.DB) *gorm.DB {
	updates := map[string]interface{}{
		"login_count":   gorm.Expr("login_count + ?", 1),
		"last_login_ip": m.LastLoginIp,
		"last_login_at": time.Now().Format("2006-01-02 15:04.05"),
	}

	return tx.Model(Account{}).Where(Account{IDModel: IDModel{ID: m.ID}}).Updates(updates)
}

func (m *Account) ExistsMobile(tx *gorm.DB) bool {
	if m.Mobile == "" {
		return true
	}
	if res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(Account{Mobile: m.Mobile}).Select("mobile").First(&Account{}); res.RowsAffected == 0 {
		return false
	}
	return true
}

func (m *Account) ExistsEmail(tx *gorm.DB) bool {
	if m.Email == "" {
		return true
	}

	if res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(Account{Email: m.Email}).Select("email").First(&Account{}); res.RowsAffected == 0 {
		return false
	}
	return true
}

type AccountPlatform struct {
	IDModel
	AccountID     int64  `gorm:"type:int(64);not null;comment:账号ID;index"`
	PlatformID    string `gorm:"type:varchar(128);not null default '';comment:第三方平台id;index"`
	PlatformToken string `gorm:"type:varchar(128);not null default '';comment:第三方平台access_token"`
	Type          string `gorm:"type:varchar(64);not null default '';comment:类型"`
	DateModel
}

func (m *AccountPlatform) Exists(tx *gorm.DB) bool {
	if res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(AccountPlatform{PlatformID: m.PlatformID, Type: m.Type}).Select("id").First(&AccountPlatform{}); res.RowsAffected == 0 {
		return false
	}
	return true
}
