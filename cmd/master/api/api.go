package api

import (
	"context"
	"net/http"

	"go-crontab/cmd/master/mgr"
	"go-crontab/common"
	"go-crontab/log"

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
