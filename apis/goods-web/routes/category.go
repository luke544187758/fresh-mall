package routes

import (
	"github.com/gin-gonic/gin"
	"luke544187758/goods-web/api/category"
	"luke544187758/goods-web/middlewares"
)

func InitCategoryRouter(r *gin.RouterGroup) {
	c := r.Group("/category")
	{
		c.GET("/list", category.List)
		c.POST("/create", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), category.New)
		c.GET("/:id", category.Detail)
		c.DELETE("/:id", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), category.Delete)
		c.PUT("/:id", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), category.Update)
	}
}
