package routers

import (
	"go-crontab/cmd/master/api"
	"go-crontab/cmd/master/config"

	"github.com/gin-gonic/gin"
)

func SetJobRouter(e *gin.Engine) {
	//mux.HandleFunc("/job/save", handleJobSave)
	//mux.HandleFunc("/job/delete", handleJobDelete)
	//mux.HandleFunc("/job/list", handleJobList)
	//mux.HandleFunc("/job/kill", handleJobKill)
	//mux.HandleFunc("/job/log", handleJobLog)
	e.GET("/worker/list", api.WorkerList)
	e.StaticFile("/", config.GetConf().Webroot)
	// 静态文件目录
	//staticDir = http.Dir(config.GetConf().Webroot)
	//staticHandler = http.FileServer(staticDir)
	//mux.Handle("/", http.StripPrefix("/", staticHandler)) //   ./webroot/index.html
}
