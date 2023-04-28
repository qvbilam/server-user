package model

type Level struct {
	IDModel
	Level int64  `gorm:"type:int;not null default 0;comment:等级;index"`
	Name  string `gorm:"type:varchar(20); not null default '';comment:名称"`
	Icon  string `gorm:"type:varchar(20); not null default '';comment:图标"`
	Exp   int64  `gorm:"type:int;not null default 0;comment:经验;index"`
	DateModel
}

type LevelExpLog struct {
	IDModel
	UserID       int64  `gorm:"type:int not null default 0;comment:用户id;index:idx_user"`
	Exp          int64  `gorm:"type:int;not null default 0;comment:经验"`
	ExpAfter     int64  `gorm:"type:int;not null default 0;comment:变更后经验"`
	Level        int64  `gorm:"type:int;not null default 0;comment:等级"`
	LevelAfter   int64  `gorm:"type:int;not null default 0;comment:变更后等级"`
	BusinessType string `gorm:"type:varchar(20); not null default '';comment:业务类型"`
	BusinessID   int64  `gorm:"type:int;not null default 0;comment:业务id"`
	Desc         string `gorm:"type:varchar(64); not null default '';comment:说明"`
	DateModel
}
