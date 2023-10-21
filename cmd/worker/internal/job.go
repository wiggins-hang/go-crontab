package internal

import (
	"context"
	"time"

	"go-crontab/cmd/worker/config"
	"go-crontab/common"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// 任务管理器
type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	// 单例
	GJobMgr *JobMgr
)

// 初始化管理器
func InitJobMgr() (err error) {
	var (
		etcdConfig clientv3.Config
		client     *clientv3.Client
		kv         clientv3.KV
		lease      clientv3.Lease
		watcher    clientv3.Watcher
	)

	// 初始化配置
	etcdConfig = clientv3.Config{
		Endpoints:   config.GetWorkerConf().EtcdEndpoints,                                     // 集群地址
		DialTimeout: time.Duration(config.GetWorkerConf().EtcdDialTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	if client, err = clientv3.New(etcdConfig); err != nil {
		return
	}

	// 得到KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	// 赋值单例
	GJobMgr = &JobMgr{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	// 启动任务监听
	GJobMgr.watchJobs()

	// 启动监听killer
	GJobMgr.watchKiller()

	return
}

// 监听任务变化
func (jobMgr *JobMgr) watchJobs() (err error) {
	var (
		getResp            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *common.JobEvent
	)

	// 1, get一下/cron/jobs/目录下的所有任务，并且获知当前集群的revision
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JobSaveDir, clientv3.WithPrefix()); err != nil {
		return
	}

	// 当前有哪些任务
	for _, kvpair = range getResp.Kvs {
		// 反序列化json得到Job
		if job, err = common.UnpackJob(kvpair.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JobEventType, job)
			// 同步给scheduler(调度协程)
			GScheduler.PushJobEvent(jobEvent)
		}
	}

	// 2, 从该revision向后监听变化事件
	go func() { // 监听协程
		// 从GET时刻的后续版本开始监听变化
		watchStartRevision = getResp.Header.Revision + 1
		// 监听/cron/jobs/目录的后续变化
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JobSaveDir, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		// 处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: // 任务保存事件
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					// 构建一个更新Event
					jobEvent = common.BuildJobEvent(common.JobEventType, job)
				case mvccpb.DELETE: // 任务被删除了
					// Delete /cron/jobs/job10
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))

					job = &common.Job{Name: jobName}

					// 构建一个删除Event
					jobEvent = common.BuildJobEvent(common.JobEventDelete, job)
				}
				// 变化推给scheduler
				GScheduler.PushJobEvent(jobEvent)
			}
		}
	}()
	return
}

// 监听强杀任务通知
func (jobMgr *JobMgr) watchKiller() {
	var (
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
		jobEvent   *common.JobEvent
		jobName    string
		job        *common.Job
	)
	// 监听/cron/killer目录
	go func() { // 监听协程
		// 监听/cron/killer/目录的变化
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JobKillerDir, clientv3.WithPrefix())
		// 处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: // 杀死任务事件
					jobName = common.ExtractKillerName(string(watchEvent.Kv.Key))
					job = &common.Job{Name: jobName}
					jobEvent = common.BuildJobEvent(common.JobEventKill, job)
					// 事件推给scheduler
					GScheduler.PushJobEvent(jobEvent)
				case mvccpb.DELETE: // killer标记过期, 被自动删除
				}
			}
		}
	}()
}
