package routes

import (
	"github.com/gin-gonic/gin"
	"luke544187758/order-web/logger"
	"luke544187758/order-web/middlewares"
	"luke544187758/order-web/settings"
	"net/http"
)

func Init(cfg *settings.AppConfig) *gin.Engine {
	gin.SetMode(cfg.Mode) // gin设置成发布模式
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.Cross())
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "success",
		})
	})
	v1 := r.Group("/v1")
	InitOrderRouter(v1)
	InitShopCartRouter(v1)
	InitPayRouter(v1)
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{"msg": "404"})
	})
	return r
}
