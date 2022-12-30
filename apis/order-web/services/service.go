package services

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"luke544187758/order-web/global"
	"luke544187758/order-web/proto"
	"luke544187758/order-web/settings"
)

func ServicesInit() {
	orderConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port, fmt.Sprintf("grpc.health.v1.%s", settings.Conf.OrderService.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.L().Fatal("[ServicesInit] connect order service failed...")
	}
	global.OrderServiceClient = proto.NewOrderClient(orderConn)

	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port, fmt.Sprintf("grpc.health.v1.%s", settings.Conf.GoodsService.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.L().Fatal("[ServicesInit] connect order service failed...")
	}
	global.GoodsServiceClient = proto.NewGoodsClient(goodsConn)

	invConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port, fmt.Sprintf("grpc.health.v1.%s", settings.Conf.InventoryService.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.L().Fatal("[ServicesInit] connect inventory service failed...")
	}
	global.InventoryServiceClient = proto.NewInventoryClient(invConn)
}
