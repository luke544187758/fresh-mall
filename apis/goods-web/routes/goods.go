package routes

import (
	"github.com/gin-gonic/gin"
	"luke544187758/goods-web/api/goods"
	"luke544187758/goods-web/middlewares"
)

func InitGoodsRouter(r *gin.RouterGroup) {
	g := r.Group("/goods")
	{
		g.GET("/list", goods.List)
		g.POST("/create", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), goods.New)
		g.GET("/:id", goods.Detail)
		g.DELETE("/:id", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), goods.Delete)
		g.GET("/:id/stocks", goods.Stocks)
		g.PUT("/:id", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), goods.Update)
		g.PATCH("/:id", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), goods.UpdateStatus)
	}
}
