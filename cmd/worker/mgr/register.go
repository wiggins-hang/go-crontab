package mgr

import (
	"context"
	"net"
	"time"

	"go-crontab/cmd/worker/config"
	"go-crontab/common"
	"go-crontab/log"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 注册节点到etcd： /cron/workers/IP地址
type Register struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease

	localIP string // 本机IP
}

var (
	GRegister *Register
)

func InitRegister() (err error) {

	// 初始化配置
	etcdConf := clientv3.Config{
		Endpoints:   config.GetWorkerConf().EtcdEndpoints,                                     // 集群地址
		DialTimeout: time.Duration(config.GetWorkerConf().EtcdDialTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	client, err := clientv3.New(etcdConf)
	if err != nil {
		log.Error("register new etcd client error ", err)
		return err
	}

	// 本机IP
	localIp, err := getLocalIP()
	if err != nil {
		return err
	}

	// 得到KV和Lease的API子集
	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)

	GRegister = &Register{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIP: localIp,
	}

	// 服务注册
	go GRegister.keepOnline()

	return
}

// 获取本机网卡IP
func getLocalIP() (ipv4 string, err error) {

	// 获取所有网卡
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Error("net . InterfaceAddrs error ", err)
		return "", err
	}
	// 取第一个非lo的网卡IP
	for _, addr := range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet := addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() // 192.168.1.1
				return
			}
		}
	}

	err = common.ERR_NO_LOCAL_IP_FOUND
	return
}

// 注册到/cron/workers/IP, 并自动续租
func (register *Register) keepOnline() {
	var (
		keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
		keepAliveResp *clientv3.LeaseKeepAliveResponse
		cancelCtx     context.Context
		cancelFunc    context.CancelFunc
	)

	// 注册路径
	regKey := common.JobWorkersDir + register.localIP
	for {

		cancelFunc = nil
		// 创建租约
		leaseGrantResp, err := register.lease.Grant(context.TODO(), 10)
		if err != nil {
			Retry(cancelFunc)
			continue
		}

		// 自动续租
		keepAliveChan, err = register.lease.KeepAlive(context.TODO(), leaseGrantResp.ID)
		if err != nil {
			Retry(cancelFunc)
			continue
		}

		cancelCtx, cancelFunc = context.WithCancel(context.TODO())

		// 注册到etcd
		if _, err = register.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
			Retry(cancelFunc)
			continue
		}

		// 处理续租应答
		for {
			select {
			case keepAliveResp = <-keepAliveChan:
				if keepAliveResp == nil { // 续租失败
					Retry(cancelFunc)
					continue
				}
			}
		}

	}
}

func Retry(cancelFunc context.CancelFunc) {
	time.Sleep(1 * time.Second)
	if cancelFunc != nil {
		cancelFunc()
	}
}
