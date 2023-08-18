package config

type Config struct {
	Address         string   `json:"address"`
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
	//web页面根目录
	Webroot     string `json:"webroot"`
	MysqlTarget string `json:"mysqlTarget"`
}

var conf Config

func InitConf() {
	conf = Config{
		Address:         "0.0.0.0:8899",
		EtcdEndpoints:   []string{"127.0.0.1:2379"},
		EtcdDialTimeout: 5000,
		Webroot:         "./webroot",
		MysqlTarget:     "root:123456@tcp(127.0.0.1:3306)/crontab?charset=utf8mb4&parseTime=True&loc=Local",
	}
}

func GetConf() Config {

	return conf
}
