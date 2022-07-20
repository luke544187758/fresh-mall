package routes

import (
	"github.com/gin-gonic/gin"
	"luke544187758/goods-web/api/banners"
	"luke544187758/goods-web/middlewares"
)

func InitBannerRouter(g *gin.RouterGroup) {
	b := g.Group("/banners")
	{
		b.GET("/list", banners.List)
		b.DELETE("/:id", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), banners.Delete)
		b.POST("/create", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), banners.New)
		b.PUT("/:id", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), banners.Update)
	}
}
