package routes

import (
	"github.com/gin-gonic/gin"
	"luke544187758/order-web/api/pay"
)

func InitPayRouter(r *gin.RouterGroup) {
	group := r.Group("/pay")
	{
		group.POST("/alipay/notify", pay.Notify)
	}
}
