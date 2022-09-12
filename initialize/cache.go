package initialize

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"user/global"
)

func InitRedis() {
	global.Redis = *redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServerConfig.RedisConfig.Host, global.ServerConfig.RedisConfig.Port),
		Username: global.ServerConfig.RedisConfig.User,
		Password: global.ServerConfig.RedisConfig.Password,
		DB:       global.ServerConfig.RedisConfig.Database,
	})
}
