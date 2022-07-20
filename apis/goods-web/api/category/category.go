package category

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
	res, err := global.GoodsServiceClient.GetAllCategoryList(context.Background(), &emptypb.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	data := make([]interface{}, 0)
	for _, v := range res.Data {
		item := map[string]interface{}{
			"id":     fmt.Sprintf("%d", v.Id),
			"name":   v.Name,
			"level":  v.Level,
			"is_tab": v.IsTab,
		}
		if v.Level != 1 {
			item["parent_category"] = fmt.Sprintf("%d", v.ParentCategory)
		}
		data = append(data, item)
	}

	message.ResponseSuccess(ctx, gin.H{
		"total": len(data),
		"items": data,
	})
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}
	reMap := make(map[string]interface{})
	subCategories := make([]interface{}, 0)
	res, err := global.GoodsServiceClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{Id: idInt})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	for _, v := range res.SubCategorys {
		subCategories = append(subCategories, map[string]interface{}{
			"id":              fmt.Sprintf("%d", v.Id),
			"name":            v.Name,
			"level":           v.Level,
			"parent_category": fmt.Sprintf("%d", v.ParentCategory),
			"is_tab":          v.IsTab,
		})
	}
	reMap["id"] = fmt.Sprintf("%d", res.Info.Id)
	reMap["name"] = res.Info.Name
	reMap["level"] = res.Info.Level
	reMap["parent_category"] = fmt.Sprintf("%d", res.Info.ParentCategory)
	reMap["is_tab"] = res.Info.IsTab
	if len(subCategories) > 0 {
		reMap["sub_categories"] = subCategories
	}

	message.ResponseSuccess(ctx, reMap)
}

func New(ctx *gin.Context) {
	form := new(forms.CategoryForm)
	if err := ctx.ShouldBindJSON(&form); err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}

	res, err := global.GoodsServiceClient.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name:           form.Name,
		ParentCategory: form.ParentCategory,
		Level:          form.Level,
		IsTab:          *form.IsTab,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	data := make(map[string]interface{})
	data["id"] = fmt.Sprintf("%d", res.Id)
	data["name"] = res.Name
	data["parent_category"] = fmt.Sprintf("%d", res.ParentCategory)
	data["level"] = res.Level
	data["is_tab"] = res.IsTab

	message.ResponseSuccess(ctx, data)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}
	_, err = global.GoodsServiceClient.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{Id: idInt})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	message.ResponseSuccess(ctx, nil)
}

func Update(ctx *gin.Context) {
	form := new(forms.UpdateCategoryForm)
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
	_, err = global.GoodsServiceClient.UpdateCategory(context.Background(), &proto.CategoryInfoRequest{
		Id:    idInt,
		Name:  form.Name,
		IsTab: *form.IsTab,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	message.ResponseSuccess(ctx, nil)
}
