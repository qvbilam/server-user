package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
	"user/config"
)

var (
	DB           *gorm.DB
	ES           *elastic.Client
	Redis        redis.Client
	ServerConfig *config.ServerConfig
)
