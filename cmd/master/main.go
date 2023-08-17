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
	"go-crontab/log"
	"go-crontab/model"
)

func main() {
	InitDepend()
}

func InitDepend() {
	// 初始化配置
	config.InitConf()
	// 初始化db
	model.InitDb(config.GetConf().MysqlTarget)

	if err := mgr.InitWorkerMgr(); err != nil {
		log.Fatalf("InitWorkerMgr error %s", err.Error())
	}

	// 注册http 服务
	routers.Include(routers.SetJobRouter)
	router := routers.Init()

	server := http.Server{
		Addr:    config.GetConf().Address,
		Handler: router,
	}

	go func() {
		log.Infof(" master start listening addr %s ", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(" master service listening err: ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGHUP)
	<-quit
	time.Sleep(10 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("master server shutdown:", err)
	}
}
