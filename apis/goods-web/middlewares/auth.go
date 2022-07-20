package middlewares

import (
	"github.com/gin-gonic/gin"
	"luke544187758/goods-web/message"
	"luke544187758/goods-web/pkg/jwt"
	"strings"
)

func JWTAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// 客户端携带Token的三种方式 1、放在请求头 2、放在请求体 3、放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// Authorization：Bearer xxxx.xxxx.xxxx

		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			message.ResponseError(ctx, message.CodeNeedLogin)
			ctx.Abort()
			return
		}
		// 检查auth格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			message.ResponseError(ctx, message.CodeInvalidAuth)
			ctx.Abort()
			return
		}
		//解析auth中的token
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			message.ResponseError(ctx, message.CodeInvalidAuth)
			ctx.Abort()
			return
		}
		// 解析auth中的token，则通过验证，将信息保存在上下文中
		ctx.Set("claims", mc)
		ctx.Set("user_id", mc.ID)
		ctx.Next()
	}
}
