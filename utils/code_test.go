package utils

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
	})
}

func TestGenerateUserCode(t *testing.T) {
	setClient()
	err := GenerateUserCode(1)
	if err != nil {
		fmt.Println("生成错误", err)
		return
	}
	return
}
