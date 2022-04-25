package services

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"luke544187758/user-web/global"
	"luke544187758/user-web/proto"
	"luke544187758/user-web/settings"
)

func ServiceInit() error {
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port, fmt.Sprintf("grpc.health.v1.%s", settings.Conf.UserService.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		return err
	}
	global.UserServiceClient = proto.NewUserClient(conn)
	return nil
}

func ServiceInit_BAK() error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port)
	cli, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	userSrvHost := ""
	userSrvPort := 0
	data, err := cli.Agent().ServicesWithFilter(fmt.Sprintf("Service==\"%s\"", settings.Conf.UserService.Name))
	if err != nil {
		return err
	}
	for _, val := range data {
		userSrvHost = val.Address
		userSrvPort = val.Port
		break
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	global.UserServiceClient = proto.NewUserClient(conn)
	return nil
}
