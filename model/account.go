package model

import (
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

type AccountPlatform struct {
	IDModel
	AccountID     int64  `gorm:"type:int(64);not null;comment:账号ID;index"`
	PlatformID    string `gorm:"type:varchar(128);not null default '';comment:第三方平台id;index"`
	PlatformToken string `gorm:"type:varchar(128);not null default '';comment:第三方平台access_token"`
	Type          string `gorm:"type:varchar(64);not null default '';comment:类型"`
	DateModel
}

type AccountLog struct {
	IDModel
	AccountId int64  `gorm:"type:int(64);not null;comment:账号ID;index"`
	Type      string `gorm:"type:varchar(128);not null default '';comment:类型,login,logout,sign-in,sign-out;index"`
	Method    string `gorm:"type:varchar(128);not null default '';comment:方式"`
	Client    string `gorm:"type:varchar(128);not null default '';comment:客户端web,iOS,Android"`
	Version   string `gorm:"type:varchar(128);not null default '';comment:应用版本"`
	Device    string `gorm:"type:varchar(128);not null default '';comment:设备"`
	Ip        string `gorm:"type:varchar(128);not null default '';comment:设备"`
	DateModel
}
