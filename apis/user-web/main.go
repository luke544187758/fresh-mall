package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"luke544187758/user-web/global"
	"luke544187758/user-web/logger"
	"luke544187758/user-web/routes"
	"luke544187758/user-web/services"
	"luke544187758/user-web/settings"
	"luke544187758/user-web/utils"
	"luke544187758/user-web/utils/register/consul"
	"luke544187758/user-web/validators"
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

	if err := validators.InitTrans("zh"); err != nil {
		zap.L().Error("init trans failed", zap.Error(err))
	} else {
		zap.L().Info("trans init success...")
	}

	if err := services.ServiceInit(); err != nil {
		zap.L().Error("init user service failed", zap.Error(err))
		return
	}
	zap.L().Info("user service init success...")

	cli := consul.NewRegistry(settings.Conf.ConsulConfig.Host, settings.Conf.ConsulConfig.Port)
	err := cli.Register(settings.Conf.Host, settings.Conf.Port, settings.Conf.Name, settings.Conf.Tags)
	if err != nil {
		zap.L().Panic("user-web register failed", zap.Error(err))
	}
	zap.L().Info("user-web register success...")

	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", validators.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 无效的手机号码!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	// 注册路由
	r := routes.Init(settings.Conf)

	srvPort := settings.Conf.Port
	// 生产环境下，获取动态端口号
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
	// 等待中断信号以优雅的关闭服务器 (设置 5s 的超时的时间)
	quit := make(chan os.Signal, 1) // 创建一个接受信号的通道
	// kill 默认会发送 syscall.SIGTERM  信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的 Ctrl+C 就是触发系统 SIGINT 信号
	// kill -9 发送 syscall.SIGKILL信号，但是不能被捕获，所以不需要添加它
	// signal.Notify 把收到的 syscall.SIGINT 或 syscall.SIGTERM 信号转发给 quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	if err := cli.DeRegister(settings.Conf.Name, settings.Conf.Host, settings.Conf.Port); err != nil {
		zap.L().Error("user-web deregister failed", zap.Error(err))
	}
	zap.L().Info("user-web deregister success...")

	// 创建一个5秒超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
