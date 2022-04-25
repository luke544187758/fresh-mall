package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"luke544187758/health"
	"luke544187758/user-srv/dao/mysql"
	"luke544187758/user-srv/logger"
	"luke544187758/user-srv/logic"
	"luke544187758/user-srv/proto"
	"luke544187758/user-srv/settings"
	"luke544187758/user-srv/utils"
	"net"
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

	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	zap.L().Info("logger init success...")

	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	zap.L().Info("mysql init success...")

	var srvPort int
	var err error
	if settings.Conf.Mode == "dev" {
		srvPort, err = utils.GetFreePort()
		if err != nil {
			zap.L().Error("utils.GetFreePort failed", zap.Error(err))
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
		Tags:                           []string{"user-service"},
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	health.Register(ccfg, scfg)

	grpcServer := grpc.NewServer()
	service := logic.NewUserService()
	proto.RegisterUserServer(grpcServer, service)
	grpc_health_v1.RegisterHealthServer(grpcServer, &HealthImpl{})
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.Conf.Host, srvPort))
	if err != nil {
		zap.L().Error("cannot start server", zap.Error(err))
		return
	}
	defer listener.Close()

	zap.L().Info("user service start....", zap.String("address", listener.Addr().String()))
	if err := grpcServer.Serve(listener); err != nil {
		zap.L().Error("cannot start server", zap.Error(err))
		return
	}
}
