package global

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"user/config"
)

var (
	DB           *gorm.DB
	Redis        redis.Client
	ServerConfig config.ServerConfig
)
