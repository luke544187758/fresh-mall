package routes

import (
	"github.com/gin-gonic/gin"
	"luke544187758/user-web/api"
	"luke544187758/user-web/middlewares"
)

func InitUserRouter(r *gin.RouterGroup) {
	u := r.Group("/user")
	{
		u.GET("/list", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), api.GetUserList)
		u.POST("/login", api.PasswordLogin)
		u.POST("/register", api.Register)
	}
}
