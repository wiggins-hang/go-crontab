package internal

import (
	"time"

	"go-crontab/cmd/worker/config"
	"go-crontab/common"
	"go-crontab/log"
	"go-crontab/model/crontab"
	"go-crontab/model/db_model"
)

// LogSink 存储日志
type LogSink struct {
	logChan        chan *db_model.JobLog
	autoCommitChan chan *common.LogBatch
}

var (
	GLogSink *LogSink
)

func InitLogSink() (err error) {

	GLogSink = &LogSink{
		logChan:        make(chan *db_model.JobLog, 1000),
		autoCommitChan: make(chan *common.LogBatch, 1000),
	}

	// 启动一个处理协程日志上报
	go GLogSink.writeLoop()
	return
}

// 日志存储协程
func (logSink *LogSink) writeLoop() {
	var (
		logBatch     *common.LogBatch // 当前的批次
		commitTimer  *time.Timer
		timeoutBatch *common.LogBatch // 超时批次
	)

	for {
		select {
		case jobLog := <-logSink.logChan:
			if logBatch == nil {
				logBatch = &common.LogBatch{}
				// 让这个批次超时自动提交(给1秒的时间）
				commitTimer = time.AfterFunc(
					time.Duration(config.GetWorkerConf().JobLogCommitTimeout)*time.Millisecond,
					func(batch *common.LogBatch) func() {
						return func() {
							logSink.autoCommitChan <- batch
						}
					}(logBatch),
				)
			}

			// 把新日志追加到批次中
			logBatch.Logs = append(logBatch.Logs, jobLog)

			// 如果批次满了, 就立即发送
			if len(logBatch.Logs) >= config.GetWorkerConf().JobLogBatchSize {
				// 发送日志
				logSink.saveLogs(logBatch)
				// 清空logBatch
				logBatch = nil
				// 取消定时器
				commitTimer.Stop()
			}
		case timeoutBatch = <-logSink.autoCommitChan: // 过期的批次
			// 判断过期批次是否仍旧是当前的批次
			if timeoutBatch != logBatch {
				continue // 跳过已经被提交的批次
			}
			// 把批次写入
			logSink.saveLogs(timeoutBatch)
			// 清空logBatch
			logBatch = nil
		}
	}
}

// 批量写入日志
func (logSink *LogSink) saveLogs(batch *common.LogBatch) {

	if err := crontab.CreateBatchJob(crontab.GetDb(), batch.Logs); err != nil {
		log.Error("save logs mysql error ", err)
	}
}

// 发送日志
func (logSink *LogSink) Append(jobLog *db_model.JobLog) {
	select {
	case logSink.logChan <- jobLog:
	default:
		// 队列满了就丢弃
	}
}
