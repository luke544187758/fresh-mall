package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"luke544187758/health"
	"luke544187758/order-srv/dao/mysql"
	"luke544187758/order-srv/logic"
	"luke544187758/order-srv/pkg/snowflake"
	"luke544187758/order-srv/proto"
	"luke544187758/order-srv/services"
	"luke544187758/order-srv/settings"
	"luke544187758/order-srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HealthImpl struct {
}

func (h *HealthImpl) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (h *HealthImpl) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}

func main() {
	if err := settings.Init(); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}

	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	fmt.Println("mysql init success...")

	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	if err := services.ServicesInit(); err != nil {
		fmt.Printf("init services failed, err:%v", err)
		return
	}
	fmt.Println("all services init success...")

	var srvPort int
	var err error
	if settings.Conf.Mode == "dev" {
		srvPort, err = utils.GetFreePort()
		if err != nil {
			fmt.Printf("utils.GetFreePort failed, err:%v\n", err)
		}
		if srvPort == 0 {
			srvPort = settings.Conf.Port
		}
	} else {
		srvPort = settings.Conf.Port
	}

	// 注册服务到consul
	ccfg := &health.ConsulConfig{
		Host: settings.Conf.ConsulConfig.Host,
		Port: settings.Conf.ConsulConfig.Port,
	}
	scfg := &health.ServiceConfig{
		Port:                           srvPort,
		Host:                           settings.Conf.Host,
		Name:                           settings.Conf.Name,
		Tags:                           settings.Conf.Tags,
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	health.Register(ccfg, scfg)

	grpcServer := grpc.NewServer()
	service := logic.NewOrderService()
	proto.RegisterOrderServer(grpcServer, service)
	grpc_health_v1.RegisterHealthServer(grpcServer, &HealthImpl{})
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.Conf.Host, srvPort))
	if err != nil {
		fmt.Printf("cannot start server, err:%v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("order service start....", "address:", listener.Addr().String())

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			fmt.Printf("cannot start server, err:%v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")
	if err = health.Deregister(ccfg, scfg); err != nil {
		zap.L().Error("order service deregister failed", zap.Error(err))
	}
	zap.L().Info("order service deregister success...")

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcServer.Stop()

	zap.L().Info("Server exiting")
}
