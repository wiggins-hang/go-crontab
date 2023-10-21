package routers

import (
	"go-crontab/cmd/master/api"
	"go-crontab/cmd/master/config"

	"github.com/gin-gonic/gin"
)

func SetJobRouter(e *gin.Engine) {
	jobGroup := e.Group("/job")
	jobGroup.POST("/job/save", api.JobSave)
	jobGroup.POST("/job/delete", api.JobDelete)
	jobGroup.GET("/job/list", api.JobList)
	jobGroup.POST("/job/kill", api.JobKill)
	jobGroup.GET("/job/log", api.JobLog)
	jobGroup.GET("/worker/list", api.WorkerList)
	// 静态文件目录
	// ./webroot/index.html
	e.StaticFile("/", config.GetConf().Webroot)
}
