package middlewares

import (
	"github.com/gin-gonic/gin"
	"luke544187758/goods-web/message"
	"luke544187758/goods-web/pkg/jwt"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mc, _ := ctx.Get("claims")
		myClaims := mc.(*jwt.MyClaims)
		if myClaims.AuthorityId != 2 {
			message.ResponseError(ctx, message.CodeNotAdmin)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
