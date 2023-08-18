package mgr

import (
	"context"

	"go-crontab/model/crontab"
	"go-crontab/model/db_model"
)

// mongodb日志管理
type LogMgr struct {
}

var (
	GLogMgr *LogMgr
)

func init() {
	GLogMgr = &LogMgr{}
}

func (logMgr *LogMgr) ListLog(ctx context.Context, name string, skip int, limit int) ([]*db_model.JobLog, error) {
	return crontab.JobLogPageList(crontab.GetDb(), crontab.LogPageFilter{
		Skip:  skip,
		Limit: limit,
		Name:  name,
	})
}
