package cache

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
	"user/global"
)

type RedisServer struct {
}

const DefaultUserCodeDigit = 4

func (s *RedisServer) Exists(key string) int64 {
	result, _ := global.Redis.Exists(context.Background(), key).Result()
	return result
}

func (s *RedisServer) SetNX(key string, value string, expire int) bool {
	result, _ := global.Redis.SetNX(context.Background(), key, value, time.Duration(expire)*time.Second).Result()
	return result
}

func (s *RedisServer) Get(key string) string {
	result, _ := global.Redis.Get(context.Background(), key).Result()
	return result
}

func (s *RedisServer) Delete(key string, value string) int64 {
	if value != "" {
		result, _ := global.Redis.Get(context.Background(), key).Result()
		if result != value {
			return 0
		}
	}
	result, _ := global.Redis.Del(context.Background(), key).Result()
	return result
}

func (s *RedisServer) Clear(keys ...string) int64 {
	result, _ := global.Redis.Del(context.Background(), keys...).Result()
	return result
}

func (s *RedisServer) GetUserCodeDigit() int64 {
	key := RedisKey{}.GetGeneratorUserCodeMaxDigit()
	value := s.Get(key)
	i, _ := strconv.Atoi(value)
	if i == 0 {
		return DefaultUserCodeDigit
	}
	return int64(i)
}

func (s *RedisServer) GetUserSpecialCodeDigit() int64 {
	key := RedisKey{}.GetGeneratorUserSpecialCodeMaxDigit()
	value := s.Get(key)
	i, _ := strconv.Atoi(value)
	if i == 0 {
		return DefaultUserCodeDigit
	}
	return int64(i)
}

func (s *RedisServer) SetUserCodeDigit(digit int64) (string, error) {
	key := RedisKey{}.GetGeneratorUserCodeMaxDigit()
	maxDigit := s.GetUserCodeDigit()
	if maxDigit >= digit {
		return "", status.Errorf(codes.InvalidArgument, "无法设置指定位数, 当前最大位数: %d", maxDigit)
	}
	val := strconv.Itoa(int(digit))
	result, _ := global.Redis.Set(context.Background(), key, val, 0).Result()
	return result, nil
}

func (s *RedisServer) SetUserSpecialCodeDigit(digit int64) (string, error) {
	key := RedisKey{}.GetGeneratorUserSpecialCodeMaxDigit()
	maxDigit := s.GetUserSpecialCodeDigit()
	if maxDigit >= digit {
		return "", status.Errorf(codes.InvalidArgument, "无法设置指定位数, 当前最大位数: %d", maxDigit)
	}
	val := strconv.Itoa(int(digit))
	result, _ := global.Redis.Set(context.Background(), key, val, 0).Result()
	return result, nil
}

func (s *RedisServer) GenerateUserCodes(digit int64, data []interface{}) (int64, error) {
	// 验证位数
	if _, err := s.SetUserCodeDigit(digit); err != nil {
		return 0, err
	}

	// 设置锁
	lockKey := RedisKey{}.GetGeneratorUserCodeLockKey(digit)
	if res := s.SetNX(lockKey, strconv.Itoa(int(digit)), 0); res == false {
		return 0, status.Errorf(codes.AlreadyExists, "当前位数已存在")
	}

	// 生成集合
	key := RedisKey{}.GetUserCodesKey()
	result, _ := global.Redis.SAdd(context.Background(), key, data).Result()
	if result <= 0 {
		_ = s.Delete(lockKey, strconv.Itoa(int(digit))) // 添加失败, 删除锁
		return 0, status.Errorf(codes.InvalidArgument, "添加元素失败")
	}
	return result, nil
}

func (s *RedisServer) GenerateUserSpecialCodes(digit int64, data []interface{}) (int64, error) {
	// 验证位数
	if _, err := s.SetUserSpecialCodeDigit(digit); err != nil {
		return 0, err
	}

	// 设置锁
	lockKey := RedisKey{}.GetGeneratorUserCodeLockKey(digit)
	if res := s.SetNX(lockKey, strconv.Itoa(int(digit)), 0); res == false {
		return 0, status.Errorf(codes.AlreadyExists, "当前位数已存在")
	}

	// 生成集合
	key := RedisKey{}.GetUserSpecialCodesKey()
	result, _ := global.Redis.SAdd(context.Background(), key, data).Result()
	if result <= 0 {
		_ = s.Delete(lockKey, strconv.Itoa(int(digit))) // 添加失败, 删除锁
		return 0, status.Errorf(codes.InvalidArgument, "添加元素失败")
	}
	return result, nil
}

func (s *RedisServer) RandomUserCodes(count int64) ([]string, error) {
	key := RedisKey{}.GetUserCodesKey()
	return global.Redis.SPopN(context.Background(), key, count).Result()
}

func (s *RedisServer) RandomUserSpecialCodes(count int64) ([]string, error) {
	key := RedisKey{}.GetUserSpecialCodesKey()
	return global.Redis.SPopN(context.Background(), key, count).Result()
}
