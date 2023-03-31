package cache

import "strconv"

type RedisKey struct {
}

const prefix = "qvbilam:"

// GetGeneratorUserCodeLockKey 用户 code 生成器
func (k RedisKey) GetGeneratorUserCodeLockKey(count int64) string {
	return prefix + "generator:user:codes:lock:" + strconv.Itoa(int(count))
}

func (k RedisKey) GetGeneratorUserCodeMaxDigit() string {
	return prefix + "generator:user:codes:ordinary:digit"
}

func (k RedisKey) GetGeneratorUserSpecialCodeMaxDigit() string {
	return prefix + "generator:user:codes:special:digit"
}

// GetUserCodesKey 获取用户code
func (k RedisKey) GetUserCodesKey() string {
	return prefix + "user:codes:ordinary"
}

// GetUserSpecialCodesKey 获取用户特殊code
func (k RedisKey) GetUserSpecialCodesKey() string {
	return prefix + "user:codes:special"
}

// GetUserTokenKey 获取用户token
func (k RedisKey) GetUserTokenKey(userId int64) string {
	return prefix + "user:token:" + strconv.Itoa(int(userId))
}

// GetOAuthQQTokenKey 获取 oauth qq token
func (k RedisKey) GetOAuthQQTokenKey(appId string) string {
	return prefix + "oauth:qq:token:" + appId
}
