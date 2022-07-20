package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"luke544187758/oss-web/logger"
	"luke544187758/oss-web/middlewares"
	"luke544187758/oss-web/settings"
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

	r.LoadHTMLFiles(fmt.Sprintf("oss-web/templates/index.html"))
	// 配置静态文件夹路径 第一个参数是api，第二个是文件夹路径
	r.StaticFS("/static", http.Dir(fmt.Sprintf("oss-web/static")))
	// GET：请求方式；/hello：请求的路径
	// 当客户端以GET方法请求/hello路径时，会执行后面的匿名函数
	r.GET("", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "posts/index",
		})
	})

	v1 := r.Group("/v1")
	InitOssRouter(v1)

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{"msg": "404"})
	})
	return r
}
