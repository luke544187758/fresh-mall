package routes

import (
	"github.com/gin-gonic/gin"
	"luke544187758/order-web/api/order"
	"luke544187758/order-web/middlewares"
)

func InitOrderRouter(r *gin.RouterGroup) {
	o := r.Group("/orders").Use(middlewares.JWTAuthMiddleware())
	{
		o.GET("/list", order.List)
		o.POST("/create", order.New)
		o.GET("/:id", order.Detail)
	}
}
