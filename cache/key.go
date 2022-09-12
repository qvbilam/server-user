package cache

import "strconv"

type RedisKey struct {
	Prefix string
}

// GetGeneratorUserCodeLockKey 用户 code 生成器
func (k RedisKey) GetGeneratorUserCodeLockKey(count int64) string {
	return k.Prefix + "generator:user:codes:lock:" + strconv.Itoa(int(count))
}

func (k RedisKey) GetGeneratorUserCodeMaxDigit() string {
	return k.Prefix + "generator:user:codes:ordinary:digit"
}

func (k RedisKey) GetGeneratorUserSpecialCodeMaxDigit() string {
	return k.Prefix + "generator:user:codes:special:digit"
}

// GetUserCodesKey 获取用户code
func (k RedisKey) GetUserCodesKey() string {
	return k.Prefix + "user:codes:ordinary"
}

// GetUserSpecialCodesKey 获取用户特殊code
func (k RedisKey) GetUserSpecialCodesKey() string {
	return k.Prefix + "user:codes:special"
}

// GetUserTokenKey 获取用户token
func (k RedisKey) GetUserTokenKey(userId int64) string {
	return k.Prefix + "user:token:" + strconv.Itoa(int(userId))
}
