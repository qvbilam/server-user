package business

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user/global"
	"user/model"
)

type LevelBusiness struct {
	UserID       int64
	Exp          int64
	BusinessType string
	BusinessID   int64
	Desc         string
}

func (b *LevelBusiness) LevelExp() (bool, *model.Level, error) {
	tx := global.DB.Begin()
	// 获取经验
	var levels []model.Level
	if res := tx.Model(&model.Level{}).Order("level asc").Find(&levels); res.RowsAffected == 0 {
		tx.Rollback()
		return false, nil, nil
	}

	// 获取用户经验等级
	var user model.User
	if res := tx.Model(&model.User{}).Where(&model.User{IDModel: model.IDModel{ID: b.UserID}}).First(&user); res.RowsAffected == 0 {
		tx.Rollback()
		return false, nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	// 查找等级
	userLevel := user.Level
	currentExp := user.LevelExp + b.Exp
	var currentLevel model.Level

	for _, level := range levels {
		if currentExp < level.Exp {
			break
		}
		userLevel = level.Level
		currentLevel = level
	}

	// 增加经验日志
	log := model.LevelExpLog{
		IDModel:      model.IDModel{},
		UserID:       b.UserID,
		Exp:          b.Exp,
		ExpAfter:     user.LevelExp,
		Level:        user.Level,
		LevelAfter:   userLevel,
		BusinessType: b.BusinessType,
		BusinessID:   b.BusinessID,
		Desc:         b.Desc,
	}
	if res := tx.Save(&log); res.RowsAffected == 0 {
		tx.Rollback()
		return false, nil, status.Errorf(codes.Internal, "记录经验日志失败")
	}

	// 更新用户等级
	isUpgrade := userLevel == user.Level
	user.Level = userLevel
	user.LevelExp = currentExp
	if res := tx.Save(&user); res.RowsAffected == 0 {
		tx.Rollback()
		return false, nil, status.Errorf(codes.Internal, "更新等级失败")
	}

	tx.Commit()
	return isUpgrade, &currentLevel, nil
}
