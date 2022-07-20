package banners

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
	"luke544187758/goods-web/api"
	"luke544187758/goods-web/forms"
	"luke544187758/goods-web/global"
	"luke544187758/goods-web/message"
	"luke544187758/goods-web/proto"
	"strconv"
)

func List(ctx *gin.Context) {
	res, err := global.GoodsServiceClient.BannerList(context.Background(), &emptypb.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	data := make([]interface{}, 0)
	for _, v := range res.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = fmt.Sprintf("%d", v.Id)
		reMap["index"] = v.Index
		reMap["image"] = v.Image
		reMap["url"] = v.Url
		data = append(data, reMap)
	}
	message.ResponseSuccess(ctx, data)
}

func New(ctx *gin.Context) {
	form := new(forms.BannerForm)
	if err := ctx.ShouldBindJSON(form); err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}

	res, err := global.GoodsServiceClient.CreateBanner(context.Background(), &proto.BannerRequest{
		Index: form.Index,
		Image: form.Image,
		Url:   form.Url,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	data := make(map[string]interface{})
	data["id"] = res.Id
	data["index"] = res.Index
	data["url"] = res.Url
	data["image"] = res.Image

	message.ResponseSuccess(ctx, data)
}

func Update(ctx *gin.Context) {
	form := new(forms.BannerForm)
	if err := ctx.ShouldBindJSON(form); err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}

	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}
	_, err = global.GoodsServiceClient.UpdateBanner(context.Background(), &proto.BannerRequest{
		Id:    int32(idInt),
		Index: form.Index,
		Image: form.Image,
		Url:   form.Url,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	message.ResponseSuccess(ctx, nil)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}

	_, err = global.GoodsServiceClient.DeleteBanner(context.Background(), &proto.BannerRequest{Id: int32(idInt)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	message.ResponseSuccess(ctx, nil)
}
