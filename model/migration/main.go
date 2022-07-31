package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"user/model"
)

func main() {
	user := "root"
	password := ""
	host := "127.0.0.1"
	port := 3306
	database := "qvbilam_user"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //不带表名
		},
	})
	if err != nil {
		panic(any(err))
	}

	_ = db.AutoMigrate(
		&model.User{},
		&model.Level{},
	)
}
