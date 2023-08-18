package routers

import (
	"go-crontab/cmd/master/api"
	"go-crontab/cmd/master/config"

	"github.com/gin-gonic/gin"
)

func SetJobRouter(e *gin.Engine) {
	e.POST("/job/save", api.JobSave)
	e.POST("/job/delete", api.JobDelete)
	e.GET("/job/list", api.JobList)
	e.POST("/job/kill", api.JobKill)
	e.GET("/job/log", api.JobLog)
	e.GET("/worker/list", api.WorkerList)
	// 静态文件目录
	// ./webroot/index.html
	e.StaticFile("/", config.GetConf().Webroot)
}
