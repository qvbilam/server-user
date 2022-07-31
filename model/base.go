package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type IDModel struct {
	ID int64 `gorm:"primarykey"`
}

type DateModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

//type UserModel struct {
//	UserID int64 `gorm:"type:int not null default 0;comment:用户id;index:idx_user"`
//}

type Visible struct {
	IsVisible bool `gorm:"not null default 0;comment:是否显示"`
}

type DeletedModel struct {
	DeletedAt gorm.DeletedAt
}

type GormList []string

func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

// Paginate 自定义优雅的分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
