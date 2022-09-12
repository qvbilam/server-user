package cache

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
	"user/global"
)

func setClient() {
	host := "127.0.0.1"
	port := 6379
	global.Redis = *redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
		DB:   0,
	})
}

func TestRedisServer_SetNX(t *testing.T) {
	setClient()
	s := RedisServer{}
	key := "test"
	res := s.SetNX(key, "1", 100)
	fmt.Println(res)
}

func TestRedisServer_GetUserCodes(t *testing.T) {
	count := 10
	setClient()
	s := RedisServer{}
	res, err := s.GetUserCodes(int64(count))
	if err != nil {
		fmt.Println("获取错误", err)
		return
	}
	fmt.Println(res)
}
