package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"luke544187758/user-web/forms"
	"luke544187758/user-web/global"
	"luke544187758/user-web/global/response"
	"luke544187758/user-web/message"
	"luke544187758/user-web/pkg/jwt"
	"luke544187758/user-web/proto"
	"luke544187758/user-web/settings"
	"luke544187758/user-web/validators"
	"net/http"
	"strconv"
	"time"
)

func HandleGrpcErrorToHttp(err error, ctx *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				message.ResponseError(ctx, message.CodeUserNotExist)
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
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

func GetUserList(ctx *gin.Context) {

	pageStr := ctx.DefaultQuery("page", "0")
	page, _ := strconv.ParseInt(pageStr, 10, 64)
	sizeStr := ctx.DefaultQuery("size", "10")
	size, _ := strconv.ParseInt(sizeStr, 10, 64)

	rsp, err := global.UserServiceClient.GetUserList(context.Background(), &proto.PageInfoRequest{
		Page:     uint32(page),
		PageSize: uint32(size),
	})
	if err != nil {
		zap.L().Error("GetUserList failed", zap.Error(err))
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, val := range rsp.Data {
		user := &response.UserResponse{
			ID:       val.Id,
			Role:     val.Role,
			NickName: val.NickName,
			Birthday: val.Birthday,
			Gender:   val.Gender,
			Address:  val.Gender,
			Mobile:   val.Mobile,
			Password: val.Password,
		}
		result = append(result, user)
	}
	message.ResponseSuccess(ctx, result)
}

func PasswordLogin(ctx *gin.Context) {
	form := new(forms.PasswordLoginForm)
	if err := ctx.ShouldBindJSON(form); err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	// 验证码
	//if !store.Verify(form.CaptchaID, form.Captcha, false) {
	//	message.ResponseError(ctx, message.CodeInvalidCaptcha)
	//	return
	//}

	rsp, err := global.UserServiceClient.GetUserByMobile(context.Background(), &proto.UserMobileRequest{Mobile: form.Mobile})
	if err != nil {
		fmt.Println(err)
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				message.ResponseError(ctx, message.CodeUserNotExist)
			default:
				message.ResponseError(ctx, message.CodeInvalidParam)
			}
			return
		}
		message.ResponseError(ctx, message.CodeServerBusy)
		return
	}
	checkRes, err := global.UserServiceClient.CheckPassword(context.Background(), &proto.PasswordCheckRequest{
		Password:   form.Password,
		EncryptPwd: rsp.Password,
	})
	if err != nil {
		zap.L().Error("CheckPassword failed", zap.Error(err))
		return
	}
	if !checkRes.Success {
		message.ResponseError(ctx, message.CodeInvalidPassword)
		return
	}
	// 生成token
	token, err := jwt.GenToken(rsp.Id, rsp.Role, rsp.NickName)
	if err != nil {
		zap.L().Error("jwt.GenToken failed", zap.Error(err))
		message.ResponseError(ctx, message.CodeServerBusy)
		return
	}

	message.ResponseSuccess(ctx, gin.H{"id": rsp.Id, "nick_name": rsp.NickName, "token": token, "expired_at": time.Now().Add(
		time.Duration(settings.Conf.JWTConfig.JWTExpire) * time.Hour).Unix()})
}

func Register(ctx *gin.Context) {
	form := new(forms.RegisterForm)
	if err := ctx.ShouldBindJSON(form); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	res, err := global.UserServiceClient.CreateUser(context.Background(), &proto.UserInfoRequest{
		NickName: form.Mobile,
		Password: form.Password,
		Mobile:   form.Mobile,
	})
	if err != nil {
		zap.L().Error("cli.CreateUser failed", zap.Error(err))
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	token, err := jwt.GenToken(res.Id, 2, form.Mobile)
	if err != nil {
		zap.L().Error("jwt.GenToken failed", zap.Error(err))
		message.ResponseError(ctx, message.CodeInvalidAuth)
		return
	}

	message.ResponseSuccess(ctx, gin.H{
		"id":        res.Id,
		"nick_name": form.Mobile,
		"token":     token,
	})
}
