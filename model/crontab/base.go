package crontab

import (
	"go-crontab/log"
	"go-crontab/shutdown"

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

	// 注册db关闭事件
	shutdown.ConnectResourceListeners.RegisterStopListener(func() {
		log.Info("start to close db")
		a, _ := db.DB()
		a.Close()
		log.Info("close db success")
	})
}

func GetDb() *gorm.DB {
	return db
}
