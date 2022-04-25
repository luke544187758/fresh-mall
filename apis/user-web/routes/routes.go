package routes

import (
	"github.com/gin-gonic/gin"
	"luke544187758/user-web/logger"
	"luke544187758/user-web/middlewares"
	"luke544187758/user-web/settings"
	"net/http"
)

func Init(cfg *settings.AppConfig) *gin.Engine {
	gin.SetMode(cfg.Mode) // gin设置成发布模式
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.Cross())
	r.GET("/health", func(ctx *gin.Context) {

	})
	v1 := r.Group("/v1")
	InitUserRouter(v1)
	InitBaseRouter(v1)
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{"msg": "404"})
	})
	return r
}
