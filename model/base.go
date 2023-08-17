package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb(target string) {

	connect, err := gorm.Open(mysql.Open(target), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = connect
}

func GetDb() *gorm.DB {
	return db
}
