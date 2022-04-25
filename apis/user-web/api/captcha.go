package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"luke544187758/user-web/message"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(ctx *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 200, 5, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := captcha.Generate()
	if err != nil {
		zap.L().Error("gen capcha failed", zap.Error(err))
		message.ResponseError(ctx, message.CodeServerBusy)
		return
	}
	message.ResponseSuccess(ctx, gin.H{
		"captcha_id": id,
		"pic_path":   b64s,
	})
}
