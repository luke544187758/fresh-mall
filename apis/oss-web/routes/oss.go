package routes

import (
	"github.com/gin-gonic/gin"
	"luke544187758/oss-web/api"
	"luke544187758/oss-web/middlewares"
)

func InitOssRouter(r *gin.RouterGroup) {
	oss := r.Group("/oss")
	{
		oss.GET("/token", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), api.Token)
		oss.POST("/callback", api.HandlerRequest)
	}
}
