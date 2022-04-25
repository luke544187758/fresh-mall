package health

type ConsulConfig struct {
	Host string // 注册中心ip
	Port int    // 注册中心端口号
}

type ServiceConfig struct {
	Port                           int      // 端口号
	Host                           string   // 服务的ip
	Name                           string   // 服务名称
	Tags                           []string // 服务的标签
	Timeout                        string   // 服务超时时间
	Interval                       string   // 服务的检查频率
	DeregisterCriticalServiceAfter string   // 服务注销时间
}
