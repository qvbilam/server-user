package initialize

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"strconv"
	"user/config"
	"user/global"
)

func InitConfig() {
	initEnvConfig()
	initViperConfig()
}

func initEnvConfig() {

	serverPort, _ := strconv.Atoi(os.Getenv("PORT"))
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	esPort, _ := strconv.Atoi(os.Getenv("ES_PORT"))
	redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	redisDb, _ := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	jwtExpire, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE"))
	jaegerOutput := os.Getenv("JAEGER_OUTPUT")
	jaegerIsLog := os.Getenv("JAEGER_IS_LOG")
	jaegerOutputInt, _ := strconv.Atoi(jaegerOutput)
	jaegerOutputIsLog, _ := strconv.ParseBool(jaegerIsLog)

	if global.ServerConfig == nil {
		global.ServerConfig = &config.ServerConfig{}
	}

	global.ServerConfig.Name = os.Getenv("SERVER_NAME")
	global.ServerConfig.Port = serverPort

	global.ServerConfig.JWTConfig.SigningKey = os.Getenv("JWT_KEY")
	global.ServerConfig.JWTConfig.Issuer = os.Getenv("JWT_ISSUER")
	global.ServerConfig.JWTConfig.Expire = int64(jwtExpire)

	global.ServerConfig.DBConfig.Host = os.Getenv("DB_HOST")
	global.ServerConfig.DBConfig.Port = dbPort
	global.ServerConfig.DBConfig.User = os.Getenv("DB_USER")
	global.ServerConfig.DBConfig.Password = os.Getenv("DB_PASSWORD")
	global.ServerConfig.DBConfig.Database = os.Getenv("DB_DATABASE")

	global.ServerConfig.ESConfig.Host = os.Getenv("ES_HOST")
	global.ServerConfig.ESConfig.Port = esPort

	global.ServerConfig.RedisConfig.Host = os.Getenv("REDIS_HOST")
	global.ServerConfig.RedisConfig.Port = redisPort
	global.ServerConfig.RedisConfig.Password = os.Getenv("REDIS_PASSWORD")
	global.ServerConfig.RedisConfig.Database = redisDb

	global.ServerConfig.JaegerConfig.Server = os.Getenv("JAEGER_SERVER")
	global.ServerConfig.JaegerConfig.Host = os.Getenv("JAEGER_HOST")
	global.ServerConfig.JaegerConfig.Port = os.Getenv("JAEGER_PORT")
	global.ServerConfig.JaegerConfig.Output = int64(jaegerOutputInt)
	global.ServerConfig.JaegerConfig.IsLog = jaegerOutputIsLog
}

func initViperConfig() {
	file := "config.yaml"
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return
	}

	v := viper.New()
	v.SetConfigFile(file)
	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panicf("获取配置异常: %s", err)
	}
	// 映射配置文件
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		zap.S().Panicf("加载配置异常: %s", err)
	}
	// 动态监听配置
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&global.ServerConfig)
	})
}
