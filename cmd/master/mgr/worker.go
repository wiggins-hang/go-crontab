package mgr

import (
	"context"
	"time"

	"go-crontab/cmd/master/config"
	"go-crontab/common"
	"go-crontab/log"
	"go-crontab/shutdown"

	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
)

type WorkerMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	GWorkerMgr *WorkerMgr
)

// 获取在线worker列表
func (workerMgr *WorkerMgr) ListWorkers() ([]string, error) {
	var (
		kv       *mvccpb.KeyValue
		workerIP string
	)

	// 初始化数组
	workerArr := make([]string, 0)

	// 获取目录下所有Kv
	getResp, err := workerMgr.kv.Get(context.TODO(), common.JOB_WORKER_DIR, clientv3.WithPrefix())
	if err != nil {
		return workerArr, err
	}

	// 解析每个节点的IP
	for _, kv = range getResp.Kvs {
		// kv.Key : /cron/workers/192.168.2.1
		workerIP = common.ExtractWorkerIP(string(kv.Key))
		workerArr = append(workerArr, workerIP)
	}
	return workerArr, nil
}

func InitWorkerMgr() error {

	// 初始化配置
	etcdConfig := clientv3.Config{
		Endpoints:   config.GetConf().EtcdEndpoints,                                     // 集群地址
		DialTimeout: time.Duration(config.GetConf().EtcdDialTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	client, err := clientv3.New(etcdConfig)
	if err != nil {
		return err
	}

	// 得到KV和Lease的API子集
	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)

	GWorkerMgr = &WorkerMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	// 注册 etcd 关闭事件
	shutdown.ConnectResourceListeners.RegisterStopListener(func() {
		log.Info("worker close etcd connect start")
		GWorkerMgr.client.Close()
		log.Info("worker close etcd connect stop")
	})
	return nil
}
