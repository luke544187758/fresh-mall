package brands

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"luke544187758/goods-web/api"
	"luke544187758/goods-web/forms"
	"luke544187758/goods-web/global"
	"luke544187758/goods-web/message"
	"luke544187758/goods-web/proto"
	"strconv"
)

func BrandList(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "0")
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	pageSize := ctx.DefaultQuery("page_size", "0")
	pageSizeInt, _ := strconv.ParseInt(pageSize, 10, 64)

	rsp, err := global.GoodsServiceClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		Page:    int32(pageInt),
		PerSize: int32(pageSizeInt),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	data := make(map[string]interface{})
	items := make([]interface{}, 0)
	data["total"] = rsp.Total

	if len(rsp.Data) == 0 {
		message.ResponseSuccess(ctx, data)
		return
	}

	for _, v := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = fmt.Sprintf("%d", v.Id)
		reMap["name"] = v.Name
		reMap["logo"] = v.Logo

		items = append(items, reMap)
	}
	data["items"] = items
	message.ResponseSuccess(ctx, data)
}

func NewBrand(ctx *gin.Context) {
	form := new(forms.BrandForm)
	if err := ctx.ShouldBindJSON(form); err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}
	rsp, err := global.GoodsServiceClient.CreateBrand(context.Background(), &proto.BrandRequest{
		Name: form.Name,
		Logo: form.Logo,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	message.ResponseSuccess(ctx, gin.H{
		"id":   fmt.Sprintf("%d", rsp.Id),
		"name": rsp.Name,
		"logo": rsp.Logo,
	})
}

func DeleteBrand(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}
	_, err = global.GoodsServiceClient.DeleteBrand(context.Background(), &proto.BrandRequest{Id: idInt})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	message.ResponseSuccess(ctx, nil)
}

func UpdateBrand(ctx *gin.Context) {
	form := new(forms.BrandForm)
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
	_, err = global.GoodsServiceClient.UpdateBrand(context.Background(), &proto.BrandRequest{
		Id:   idInt,
		Name: form.Name,
		Logo: form.Logo,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	message.ResponseSuccess(ctx, nil)
}

func GetCategoryBrandList(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}

	rsp, err := global.GoodsServiceClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{Id: idInt})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	data := make(map[string]interface{})
	data["total"] = rsp.Total
	items := make([]interface{}, 0)
	for _, v := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = fmt.Sprintf("%d", v.Id)
		reMap["name"] = v.Name
		reMap["logo"] = v.Logo
		items = append(items, reMap)
	}
	data["items"] = items
	message.ResponseSuccess(ctx, data)
}

func CategoryBrandList(ctx *gin.Context) {
	rsp, err := global.GoodsServiceClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{Page: 0, PerSize: 0})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	data := make(map[string]interface{})
	data["total"] = rsp.Total
	items := make([]interface{}, 0)
	for _, v := range rsp.Data {
		items = append(items, map[string]interface{}{
			"id": fmt.Sprintf("%d", v.Id),
			"category": map[string]interface{}{
				"id":   fmt.Sprintf("%d", v.Category.Id),
				"name": v.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   fmt.Sprintf("%d", v.Brand.Id),
				"name": v.Brand.Name,
				"logo": v.Brand.Logo,
			},
		})
	}
	data["items"] = items
	message.ResponseSuccess(ctx, data)
}

func NewCategoryBrand(ctx *gin.Context) {
	form := new(forms.CategoryBrandForm)
	if err := ctx.ShouldBindJSON(form); err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}

	rsp, err := global.GoodsServiceClient.CreateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		CategoryId: form.CategoryId,
		BrandId:    form.BrandId,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	message.ResponseSuccess(ctx, gin.H{
		"id": fmt.Sprintf("%d", rsp.Id),
	})
}

func UpdateCategoryBrand(ctx *gin.Context) {
	form := new(forms.CategoryBrandForm)
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
	_, err = global.GoodsServiceClient.UpdateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id:         idInt,
		CategoryId: form.CategoryId,
		BrandId:    form.BrandId,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	message.ResponseSuccess(ctx, nil)
}

func DeleteCategoryBrand(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}
	_, err = global.GoodsServiceClient.DeleteCategoryBrand(context.Background(), &proto.CategoryBrandRequest{Id: idInt})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	message.ResponseSuccess(ctx, nil)
}
