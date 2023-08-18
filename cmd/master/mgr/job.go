package mgr

import (
	"context"
	"time"

	"go-crontab/cmd/master/config"
	"go-crontab/common"
	"go-crontab/log"
	"go-crontab/shutdown"
	"go-crontab/tools/jsoner"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 任务管理器
type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	GJobMgr *JobMgr
)

// 初始化管理器
func InitJobMgr() (err error) {

	// 初始化配置
	etcdConfig := clientv3.Config{
		Endpoints:   config.GetConf().EtcdEndpoints,                                     // 集群地址
		DialTimeout: time.Duration(config.GetConf().EtcdDialTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	client, err := clientv3.New(etcdConfig)
	if err != nil {
		return
	}

	// 得到KV和Lease的API子集

	// 赋值单例
	GJobMgr = &JobMgr{
		client: client,
		kv:     clientv3.NewKV(client),
		lease:  clientv3.NewLease(client),
	}
	// 注册 etcd 关闭事件
	shutdown.ConnectResourceListeners.RegisterStopListener(func() {
		log.Info("job close etcd connect start")
		GJobMgr.client.Close()
		log.Info("job close etcd connect stop")
	})
	return nil
}

func (jobMgr *JobMgr) ListJobs(ctx context.Context) ([]*common.Job, error) {
	// 任务保存的目录
	dirKey := common.JOB_SAVE_DIR

	// 获取目录下所有任务信息
	getResp, err := jobMgr.kv.Get(ctx, dirKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	jobList := make([]*common.Job, 0)

	for _, kvPair := range getResp.Kvs {
		job := &common.Job{}
		if err = jsoner.UnmarshalByte(kvPair.Value, job); err != nil {
			err = nil
			continue
		}
		jobList = append(jobList, job)
	}
	return jobList, nil
}

func (jobMgr *JobMgr) SaveJob(ctx context.Context, job *common.Job) (*common.Job, error) {
	// etcd的保存key
	jobKey := common.JOB_SAVE_DIR + job.Name

	jobVal, err := jsoner.MarshalToString(job)
	if err != nil {
		log.ErrorContext(ctx, "json marshal to string error ", err)
		return nil, err
	}
	// 保存到etcd
	putRsp, err := jobMgr.kv.Put(ctx, jobKey, jobVal, clientv3.WithPrevKV())
	if err != nil {
		log.ErrorContext(ctx, "put error ", err)
		return nil, err
	}
	oldJob := &common.Job{}
	// 如果是更新, 那么返回旧值
	if putRsp.PrevKv != nil {
		// 对旧值做一个反序列化
		if err = jsoner.UnmarshalByte(putRsp.PrevKv.Value, &oldJob); err != nil {
			log.ErrorContext(ctx, "unmarshal error ", err)
			return nil, err
		}
	}
	return oldJob, nil
}

func (jobMgr *JobMgr) DeleteJob(ctx context.Context, name string) (*common.Job, error) {
	// etcd中保存任务的key
	jobKey := common.JOB_SAVE_DIR + name

	oldJob := &common.Job{}
	// 从etcd中删除它
	delResp, err := jobMgr.kv.Delete(ctx, jobKey, clientv3.WithPrevKV())
	if err != nil {
		log.ErrorContextf(ctx, "del job etcd error job key is %s error : %v", jobKey, err)
		return oldJob, err
	}
	// 返回被删除的任务信息
	if len(delResp.PrevKvs) != 0 {
		// 解析一下旧值, 返回它
		if err = jsoner.UnmarshalByte(delResp.PrevKvs[0].Value, &oldJob); err != nil {
			log.ErrorContext(ctx, "json umarshal error ", err)
			return oldJob, err
		}
	}
	return oldJob, nil
}
