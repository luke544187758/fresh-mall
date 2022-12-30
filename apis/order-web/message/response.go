package message

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseData struct {
	Code    ResCode     `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseError(ctx *gin.Context, code ResCode) {
	rd := &ResponseData{
		Code:    code,
		Message: code.Msg(),
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, rd)
}

func ResponseErrorWithMsg(ctx *gin.Context, code ResCode, msg interface{}) {
	rd := &ResponseData{
		Code:    code,
		Message: msg,
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, rd)
}

func ResponseSuccess(ctx *gin.Context, data interface{}) {
	rd := &ResponseData{
		Code:    CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data:    data,
	}
	ctx.JSON(http.StatusOK, rd)
}
