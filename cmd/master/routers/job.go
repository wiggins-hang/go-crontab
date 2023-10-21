package routers

import (
	"go-crontab/cmd/master/api"
	"go-crontab/cmd/master/config"

	"github.com/gin-gonic/gin"
)

func SetJobRouter(e *gin.Engine) {
	jobGroup := e.Group("/job")
	jobGroup.POST("/save", api.JobSave)
	jobGroup.POST("/delete", api.JobDelete)
	jobGroup.GET("/list", api.JobList)
	jobGroup.POST("/kill", api.JobKill)
	jobGroup.GET("/log", api.JobLog)
	e.GET("/worker/list", api.WorkerList)
	// 静态文件目录
	// ./webroot/index.html
	e.StaticFile("/", config.GetConf().Webroot)
}
