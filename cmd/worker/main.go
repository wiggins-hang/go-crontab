package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-crontab/cmd/worker/config"
	"go-crontab/cmd/worker/internal"
	"go-crontab/log"
	"go-crontab/model/crontab"
	"go-crontab/shutdown"
)

func main() {
	InitDepend()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	log.Info("start to cron worker success")
	<-quit
	// 释放资源
	log.Info("start to release source ")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// 调用注册好的关闭func
	shutdown.ConnectResourceListeners.NotifyStopListener(ctx)
	log.Info("shutdown release success")
}

func InitDepend() {
	// 初始化配置
	config.InitConf()
	// 初始化db
	crontab.InitDb(config.GetWorkerConf().MysqlTarget)
	// 服务注册
	if err := internal.InitRegister(); err != nil {
		log.Fatalf("InitRegister error ", err)
	}
	// 启动日志协程
	if err := internal.InitLogSink(); err != nil {
		log.Fatalf("InitLogSink error ", err)
	}

	// 启动执行器
	internal.InitExecutor()

	// 启动调度器
	if err := internal.InitScheduler(); err != nil {
		log.Fatalf("InitScheduler error ", err)
	}

	// 初始化任务管理器
	if err := internal.InitJobMgr(); err != nil {
		log.Fatalf("InitJobMgr error ", err)
	}

}
