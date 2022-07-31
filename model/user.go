package model

type User struct {
	IDModel
	Code     int64  `gorm:"type:int;not null default 0;comment:用户code;index"`
	Mobile   string `gorm:"type:varchar(11);not null default '';comment:手机号;index"`
	Nickname string `gorm:"type:varchar(20); not null default '';comment:用户昵称"`
	Password string `gorm:"type:varchar(50); not null default '';comment:登陆密码"`
	Avatar   string `gorm:"type:varchar(255); not null default '';comment:头像"`
	Gender   string `gorm:"column:gender;default:male;type:varchar(6);comment:female.女,male.男"`
	Level    int64  `gorm:"type:int;not null default 0;comment:等级"`
	//LevelModel *Level
	Visible
	DateModel
	DeletedModel
}
