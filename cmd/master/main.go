package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-crontab/cmd/master/config"
	"go-crontab/cmd/master/mgr"
	"go-crontab/cmd/master/routers"
	"go-crontab/common"
	"go-crontab/log"
	"go-crontab/model"
	"go-crontab/shutdown"
)

func main() {

	InitDepend()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
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
	model.InitDb(config.GetConf().MysqlTarget)

	if err := mgr.InitWorkerMgr(); err != nil {
		log.Fatalf("InitWorkerMgr error %s", err.Error())
	}

	if err := mgr.InitJobMgr(); err != nil {
		log.Fatalf("InitJobMgr error %s", err.Error())
	}

	InitHttpServer()

}

func InitHttpServer() {
	// 注册http 服务
	routers.Include(routers.SetJobRouter)
	router := routers.Init()

	server := http.Server{
		Addr:    config.GetConf().Address,
		Handler: router,
	}

	common.SafelyGo(func() {
		log.Infof(" master start listening addr %s ", server.Addr)
		shutdown.ConnectResourceListeners.RegisterStopListener(func() {
			log.Info("start to close  http connect start")
			if err := server.Shutdown(context.Background()); err != nil {
				log.Error("close http server error ", err)
			}
			log.Info("close http connect stop")
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(" master service listening err: ", err)
		}
	})
}
