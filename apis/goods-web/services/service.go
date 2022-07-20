package services

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"luke544187758/goods-web/global"
	"luke544187758/goods-web/proto"
	"luke544187758/goods-web/settings"
)

func ServiceInit() error {
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port, fmt.Sprintf("grpc.health.v1.%s", settings.Conf.GoodsService.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		return err
	}
	global.GoodsServiceClient = proto.NewGoodsClient(conn)
	return nil
}
