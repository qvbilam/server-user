package global

import (
	"gorm.io/gorm"
	"user/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
)
