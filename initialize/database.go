package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
	"user/global"
)

func InitDatabase() {
	User := global.ServerConfig.DBConfig.User
	Password := global.ServerConfig.DBConfig.Password
	Host := global.ServerConfig.DBConfig.Host
	Port := global.ServerConfig.DBConfig.Port
	Database := global.ServerConfig.DBConfig.Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", User, Password, Host, Port, Database)
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名不带s
		},
		//Logger: newLogger(),
	})

	if err != nil {
		zap.S().Panicf("连接数据库异常: %s", err.Error())
	}
}

func newLogger() logger.Interface {
	log.Writer()
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,       // 禁用彩色打印
		},
	)
	return newLogger
}
