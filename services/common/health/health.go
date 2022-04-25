package health

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
)

func Register(ccfg *ConsulConfig, scfg *ServiceConfig) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", ccfg.Host, ccfg.Port)

	cli, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = fmt.Sprintf("grpc.health.v1.%s", scfg.Name)
	registration.ID = fmt.Sprintf("%s-%s-%d", scfg.Name, scfg.Host, scfg.Port)
	registration.Port = scfg.Port
	registration.Tags = scfg.Tags
	registration.Address = scfg.Host
	// 生成对应的检查对象
	registration.Check = &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d/%s", scfg.Host, scfg.Port, scfg.Name),
		GRPCUseTLS:                     false,
		Timeout:                        scfg.Timeout,
		Interval:                       scfg.Interval,
		DeregisterCriticalServiceAfter: scfg.DeregisterCriticalServiceAfter,
	}
	if err := cli.Agent().ServiceRegister(registration); err != nil {
		log.Fatal(err)
	}
}

func AllServices(ccfg *ConsulConfig) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", ccfg.Host, ccfg.Port)

	cli, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	data, err := cli.Agent().Services()
	if err != nil {
		log.Fatal(err)
	}
	for k, _ := range data {
		fmt.Println(k)
	}
}

func FilterService(ccfg *ConsulConfig, name string) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", ccfg.Host, ccfg.Port)

	cli, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	data, err := cli.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, name))
	if err != nil {
		log.Fatal(err)
	}
	for k, _ := range data {
		fmt.Println(k)
	}
}
