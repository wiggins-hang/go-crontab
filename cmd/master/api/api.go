package api

import (
	"context"
	"net/http"
	"strconv"

	"go-crontab/cmd/master/mgr"
	"go-crontab/common"
	"go-crontab/log"
	"go-crontab/tools/jsoner"

	"github.com/gin-gonic/gin"
)

// WorkerList 获取健康worker节点列表
func WorkerList(ctx *gin.Context) {
	var (
		err error
	)
	workerArr := make([]string, 0)
	if workerArr, err = mgr.GWorkerMgr.ListWorkers(); err != nil {
		log.ErrorContextf(context.Background(), "mgr get list workers error ", err)
		ctx.JSON(http.StatusOK, common.Response{
			Errno: -1, Msg: err.Error(), Data: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, common.Response{
		Errno: 0,
		Msg:   "",
		Data:  workerArr,
	})
	return
}

func JobList(ctx *gin.Context) {

	// 获取任务列表
	jobList, err := mgr.GJobMgr.ListJobs(ctx)
	if err != nil {
		log.ErrorContext(ctx, "list jobs error ", err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, common.Response{
		Errno: 0,
		Msg:   "",
		Data:  jobList,
	})
}

// JobSave 保存任务接口
// POST job={"name": "job1", "command": "echo hello", "cronExpr": "* * * * *"}
func JobSave(ctx *gin.Context) {
	req := common.Job{}
	job := ctx.PostForm("job")
	if err := jsoner.Unmarshal(job, &req); err != nil {
		log.ErrorContext(ctx, "unmarshal job error ", err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	if req.CronExpr == "" || req.Name == "" || req.Command == "" {
		log.ErrorContext(ctx, "param is error ")
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	// 保存到etcd
	oldJob, err := mgr.GJobMgr.SaveJob(ctx, &req)
	if err != nil {
		log.ErrorContextf(ctx, "job save to etcd error ", err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, common.Response{
		Errno: 0,
		Msg:   "",
		Data:  oldJob,
	})

}

func JobDelete(ctx *gin.Context) {
	name := ctx.PostForm("name")

	if name == "" {
		log.ErrorContext(ctx, "name is empty ")
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	oldJob, err := mgr.GJobMgr.DeleteJob(ctx, name)
	if err != nil {
		log.ErrorContextf(ctx, "delete job error name is %s error : %v", name, err)
		ctx.JSON(http.StatusInternalServerError, oldJob)
		return
	}

	ctx.JSON(http.StatusOK, common.Response{
		Errno: 0,
		Msg:   "",
		Data:  oldJob,
	})
}

func JobKill(ctx *gin.Context) {
	name := ctx.PostForm("name")

	if name == "" {
		log.ErrorContext(ctx, "name is empty ")
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	if _, err := mgr.GJobMgr.DeleteJob(ctx, name); err != nil {
		log.ErrorContextf(ctx, "delete job error ", err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.JSON(http.StatusOK, common.Response{})
}

func JobLog(ctx *gin.Context) {
	// 获取请求参数 /job/log?name=job10&skip=0&limit=10
	var err error
	name := ctx.Query("name")
	skip := 0
	limit := 20
	if skip, err = strconv.Atoi(ctx.DefaultQuery("skip", "0")); err != nil {
		skip = 0
	}
	if limit, err = strconv.Atoi(ctx.Query("limit")); err != nil {
		limit = 20
	}

	logArr, err := mgr.GLogMgr.ListLog(ctx, name, skip, limit)
	if err != nil {
		log.ErrorContext(ctx, "get log list mysql error ", err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusOK, common.Response{
		Errno: 0,
		Msg:   "",
		Data:  logArr,
	})
}
