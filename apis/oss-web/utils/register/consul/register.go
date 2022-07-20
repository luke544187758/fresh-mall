package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type RegistryClient interface {
	Register(addr string, port int, name string, tags []string) error
	DeRegister(name, addr string, port int) error
}

type Registry struct {
	Host string
	Port int
}

func NewRegistry(host string, port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

func (r *Registry) Register(addr string, port int, name string, tags []string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	cli, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = fmt.Sprintf("http.health.v1.%s", name)
	registration.ID = fmt.Sprintf("%s-%s-%d", name, addr, port)
	registration.Port = port
	registration.Tags = tags

	// 生成对应的检查对象
	registration.Check = &api.AgentServiceCheck{
		Interval:                       "5s",
		Timeout:                        "5s",
		HTTP:                           fmt.Sprintf("http://%s:%d/health", addr, port),
		DeregisterCriticalServiceAfter: "10s",
	}

	if err := cli.Agent().ServiceRegister(registration); err != nil {
		panic(err)
	}
	return nil
}

func (r *Registry) DeRegister(name, addr string, port int) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	cli, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	return cli.Agent().ServiceDeregister(fmt.Sprintf("%s-%s-%d", name, addr, port))
}
