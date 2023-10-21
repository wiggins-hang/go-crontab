package config

// 程序配置
type WorkerConf struct {
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
	// 为了减少mysql网络往返, 打包成一批写入
	JobLogBatchSize int `json:"jobLogBatchSize"`
	// 在批次未达到阀值之前, 超时会自动提交batch
	JobLogCommitTimeout int    `json:"jobLogCommitTimeout"`
	MysqlTarget         string `json:"mysqlTarget"`
}

var conf WorkerConf

func InitConf() {
	conf = WorkerConf{
		EtcdEndpoints:       []string{"127.0.0.1:2379"},
		EtcdDialTimeout:     5000,
		JobLogBatchSize:     100,
		JobLogCommitTimeout: 100,
		MysqlTarget:         "test:test@tcp(127.0.0.1:3306)/cron?charset=utf8mb4&parseTime=True&loc=Local",
	}
}

func GetWorkerConf() WorkerConf {

	return conf
}
