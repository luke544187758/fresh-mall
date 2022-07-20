package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"luke544187758/oss-web/logger"
	"luke544187758/oss-web/routes"
	"luke544187758/oss-web/settings"
	"luke544187758/oss-web/utils"
	"luke544187758/oss-web/utils/register/consul"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings config failed, err:%v\n", err)
		return
	}

	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Info("logger init success...")

	cli := consul.NewRegistry(settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port)
	err := cli.Register(settings.Conf.Host, settings.Conf.Port, settings.Conf.Name, settings.Conf.Tags)
	if err != nil {
		zap.L().Panic("oss-web register failed", zap.Error(err))
	}
	zap.L().Info("oss-web register success...")

	// 注册路由
	r := routes.Init(settings.Conf)

	srvPort := settings.Conf.Port
	// 生产环境下，获取动态端口
	if settings.Conf.Mode == "release" {
		port, err := utils.GetFreePort()
		if err != nil {
			zap.L().Error("get free port failed", zap.Error(err))
		} else {
			srvPort = port
		}
	}
	//  启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", srvPort),
		Handler: r,
	}
	go func() {
		// service connection
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen:", zap.Error(err))
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	if err := cli.DeRegister(settings.Conf.Name, settings.Conf.Host, settings.Conf.Port); err != nil {
		zap.L().Error("oss-web deregister failed", zap.Error(err))
	}
	zap.L().Info("oss-web deregister success...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
