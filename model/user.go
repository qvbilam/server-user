package model

import (
	"context"
	"gorm.io/gorm"
	"strconv"
	"time"
	"user/global"
)

type User struct {
	IDModel
	AccountId   int64  `gorm:"type:int;not null;comment:用户code;index"`
	Code        int64  `gorm:"type:int;not null default 0;comment:用户code;index"`
	Nickname    string `gorm:"type:varchar(20); not null default '';comment:用户昵称"`
	Avatar      string `gorm:"type:varchar(255); not null default '';comment:头像"`
	Introduce   string `gorm:"type:varchar(2048); not null default '';comment:简介"`
	Gender      string `gorm:"column:gender;default:male;type:varchar(6);comment:female.女,male.男"`
	VipLevel    int64  `gorm:"type:int;not null default 0;comment:会员等级"`
	Level       int64  `gorm:"type:int;not null default 0;comment:等级"`
	LevelExp    int64  `gorm:"type:int;not null default 0;comment:等级经验"`
	FansCount   int64  `gorm:"type:int;not null default 0;comment:粉丝数量"`
	FollowCount int64  `gorm:"type:int;not null default 0;comment:关注数量"`
	//LevelModel *Level
	Visible
	DateModel
	DeletedAt *time.Time `gorm:"type:datetime"`
}

func (user *User) AfterCreate(tx *gorm.DB) error {
	esModel := userModelToEsIndex(user)
	// 写入es
	_, err := global.ES.
		Index().
		Index(UserES{}.GetIndexName()).
		BodyJson(esModel).
		Id(strconv.Itoa(int(user.ID))).
		Do(context.Background())
	return err
}

func (user *User) AfterUpdate(tx *gorm.DB) error {
	esModel := userModelToEsIndex(user)
	// 更新es. 指定 id 防止重复
	_, err := global.ES.
		Update().
		Index(UserES{}.GetIndexName()).
		Doc(esModel).
		Id(strconv.Itoa(int(user.ID))).
		Do(context.Background())

	return err
}

func (user *User) AfterDelete(tx *gorm.DB) error {
	// 删除 es 数据
	_, err := global.ES.
		Delete().
		Index(UserES{}.GetIndexName()).
		Id(strconv.Itoa(int(user.ID))).
		Do(context.Background())

	return err
}

func userModelToEsIndex(user *User) *UserES {
	return &UserES{
		ID:          user.ID,
		Code:        user.Code,
		Level:       user.Level,
		FansCount:   user.FansCount,
		FollowCount: user.FollowCount,
		Nickname:    user.Nickname,
		Introduce:   user.Introduce,
		Gender:      user.Gender,
		IsVisible:   user.IsVisible,
	}
}
