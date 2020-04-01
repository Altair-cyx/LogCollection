package conf

// AppConf 总配置
type AppConf struct {
	KafkaConf `ini:"kafka"`
	EtcdConf  `ini:"etcd"`
}

// KafkaConf Kafka配置
type KafkaConf struct {
	Address     string `ini:"address"`
	ChanMaxSize int    `ini:"chan_max_size"`
}

// EtcdConf Etcd配置
type EtcdConf struct {
	Address string `ini:"address"`
	Key     string `ini:"collect_log_key"`
	Timeout int    `ini:"timeout"`
}

//------------- unused ------------

// TaillogConf taillog配置
type TaillogConf struct {
	FileName string `ini:"path"`
}
