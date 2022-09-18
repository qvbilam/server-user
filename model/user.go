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
	Code         int64  `gorm:"type:int;not null default 0;comment:用户code;index"`
	Mobile       string `gorm:"type:varchar(11);not null default '';comment:手机号;index"`
	Nickname     string `gorm:"type:varchar(20); not null default '';comment:用户昵称"`
	Password     string `gorm:"type:varchar(128); not null default '';comment:登陆密码"`
	Avatar       string `gorm:"type:varchar(255); not null default '';comment:头像"`
	Gender       string `gorm:"column:gender;default:male;type:varchar(6);comment:female.女,male.男"`
	Level        int64  `gorm:"type:int;not null default 0;comment:等级"`
	FansCount    int64  `gorm:"type:int;not null default 0;comment:粉丝数量"`
	FollowCount  int64  `gorm:"type:int;not null default 0;comment:关注数量"`
	GetLikeCount int64  `gorm:"type:int;not null default 0;comment:获得点赞数量"`
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
		ID:           user.ID,
		Code:         user.Code,
		Level:        user.Level,
		FansCount:    user.FansCount,
		FollowCount:  user.FollowCount,
		GetLikeCount: user.GetLikeCount,
		Mobile:       user.Mobile,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Gender:       user.Gender,
		IsVisible:    user.IsVisible,
	}
}
