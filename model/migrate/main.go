package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"user/global"
	"user/initialize"
	"user/model"
)

func main() {
	initialize.InitConfig()

	user := global.ServerConfig.DBConfig.User
	password := global.ServerConfig.DBConfig.Password
	host := global.ServerConfig.DBConfig.Host
	port := global.ServerConfig.DBConfig.Port
	database := global.ServerConfig.DBConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //不带表名
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键关联
	})
	if err != nil {
		panic(any(err))
	}

	var u model.User
	db.First(&u)
	fmt.Println(host)
	fmt.Println(u)
	//_ = db.AutoMigrate(
	//	&model.User{},
	//	&model.Level{},
	//	&model.LevelExpLog{},
	//	&model.AccountPlatform{},
	//	&model.Account{},
	//	&model.AccountLog{},
	//)
}
