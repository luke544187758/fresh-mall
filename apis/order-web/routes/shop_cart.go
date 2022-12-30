package routes

import (
	"github.com/gin-gonic/gin"
	shop_cart "luke544187758/order-web/api/shop-cart"
	"luke544187758/order-web/middlewares"
)

func InitShopCartRouter(r *gin.RouterGroup) {
	cart := r.Group("/shopcarts").Use(middlewares.JWTAuthMiddleware())
	{
		cart.GET("/list", shop_cart.List)
		cart.POST("/create", shop_cart.New)
		cart.PATCH("/:id", shop_cart.Update)
		cart.DELETE("/:id", shop_cart.Delete)
	}
}
