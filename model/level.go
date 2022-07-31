package model

type Level struct {
	IDModel
	Level int64  `gorm:"type:int;not null default 0;comment:等级"`
	Name  string `gorm:"type:varchar(20); not null default '';comment:名称"`
	Icon  string `gorm:"type:varchar(20); not null default '';comment:图标"`
	DateModel
}
