package routes

import (
	"github.com/gin-gonic/gin"
	"luke544187758/user-web/api"
)

func InitBaseRouter(r *gin.RouterGroup) {
	base := r.Group("/base")
	{
		base.GET("captcha", api.GetCaptcha)
	}
}
