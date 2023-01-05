package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"luke544187758/order-web/global"
	"luke544187758/order-web/message"
	"luke544187758/order-web/validators"
	"net/http"
)

func HandleGrpcErrorToHttp(err error, ctx *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				message.ResponseError(ctx, message.CodeNotFound)
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Message(),
				})
			case codes.InvalidArgument:
				message.ResponseError(ctx, message.CodeInvalidParam)
			case codes.Unavailable:
				message.ResponseError(ctx, message.CodeServerBusy)
			case codes.AlreadyExists:
				message.ResponseErrorWithMsg(ctx, message.CodeServerBusy, gin.H{
					"msg": e.Message(),
				})

			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其它错误",
					"err": err.Error(),
				})
			}
			return
		}
	}
}

func HandleValidatorError(ctx *gin.Context, err error) {
	fmt.Println(err)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		message.ResponseErrorWithMsg(ctx, message.CodeInvalidParam, gin.H{"msg": err.Error()})
		return
	}
	message.ResponseErrorWithMsg(ctx, message.CodeInvalidParam, gin.H{"msg": validators.RemoveTopStruct(errs.Translate(global.Trans))})
}
