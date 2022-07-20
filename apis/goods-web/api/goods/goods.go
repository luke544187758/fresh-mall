package goods

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"luke544187758/goods-web/api"
	"luke544187758/goods-web/forms"
	"luke544187758/goods-web/global"
	"luke544187758/goods-web/message"
	"luke544187758/goods-web/proto"
	"strconv"
)

func List(ctx *gin.Context) {
	// 商品列表

	request := &proto.GoodsFilterRequest{}

	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.ParseInt(priceMin, 64, 10)
	request.PriceMin = int32(priceMinInt)

	priceMax := ctx.DefaultQuery("pmax", "0")
	prictMaxInt, _ := strconv.ParseInt(priceMax, 64, 10)
	request.PriceMax = int32(prictMaxInt)

	isHot := ctx.DefaultQuery("ishot", "0")
	if isHot == "1" {
		request.IsHot = true
	}
	isNew := ctx.DefaultQuery("isnew", "0")
	if isNew == "1" {
		request.IsNew = true
	}

	isTab := ctx.DefaultQuery("istab", "0")
	if isTab == "1" {
		request.IsTab = true
	}

	categoryId := ctx.DefaultQuery("cid", "0")
	categoryIdInt, _ := strconv.ParseInt(categoryId, 64, 10)
	request.TopCategory = categoryIdInt

	pageSize := ctx.DefaultQuery("page_size", "0")
	pageSizeInt, _ := strconv.ParseInt(pageSize, 64, 10)
	request.PerSize = int32(pageSizeInt)

	page := ctx.DefaultQuery("page", "0")
	pageInt, _ := strconv.ParseInt(page, 64, 10)
	request.Page = int32(pageInt)

	keywords := ctx.DefaultQuery("key", "")
	request.KeyWords = keywords

	brand := ctx.DefaultQuery("brand", "0")
	brandInt, _ := strconv.ParseInt(brand, 64, 10)
	request.Brand = brandInt

	// 请求商品的service服务
	r, err := global.GoodsServiceClient.GoodsList(context.Background(), request)
	if err != nil {
		zap.L().Error("Get goods list failed", zap.Error(err))
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := map[string]interface{}{
		"total": r.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, v := range r.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          v.Id,
			"name":        v.Name,
			"goods_brief": v.GoodsBrief,
			"ship_free":   v.ShipFree,
			"images":      v.Images,
			"desc_images": v.DescImages,
			"front_image": v.GoodsFrontImage,
			"shop_price":  v.ShopPrice,
			"category": map[string]interface{}{
				"id":   v.Category.Id,
				"name": v.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   v.Brand.Id,
				"name": v.Brand.Name,
				"logo": v.Brand.Logo,
			},
			"is_hot":  v.IsHot,
			"is_new":  v.IsNew,
			"on_sale": v.OnSale,
		})
	}
	reMap["items"] = goodsList
	message.ResponseSuccess(ctx, reMap)
}

func New(ctx *gin.Context) {
	form := new(forms.GoodsForm)
	if err := ctx.ShouldBindJSON(form); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	rsp, err := global.GoodsServiceClient.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            form.Name,
		GoodsSn:         form.GoodsSn,
		MarketPrice:     form.MarketPrice,
		ShopPrice:       form.ShopPrice,
		GoodsBrief:      form.GoodsBrief,
		ShipFree:        *form.ShipFree,
		Images:          form.Images,
		DescImages:      form.DescImages,
		GoodsFrontImage: form.GoodsFrontImage,
		CategoryId:      form.CategoryId,
		BrandId:         form.BrandId,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	//TODO 商品的库存处理
	message.ResponseSuccess(ctx, rsp)
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}

	rsp, err := global.GoodsServiceClient.GetGoodsDetail(context.Background(), &proto.GoodsInfoRequest{Id: idInt})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
	}

	data := map[string]interface{}{
		"id":                fmt.Sprintf("%d", rsp.Id),
		"name":              rsp.Name,
		"goods_brief":       rsp.GoodsBrief,
		"ship_free":         rsp.ShipFree,
		"images":            rsp.Images,
		"desc_images":       rsp.DescImages,
		"goods_front_image": rsp.GoodsFrontImage,
		"shop_price":        rsp.ShopPrice,
		"category": map[string]interface{}{
			"id":   fmt.Sprintf("%d", rsp.Category.Id),
			"name": rsp.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   fmt.Sprintf("%d", rsp.Brand.Id),
			"name": rsp.Brand.Name,
			"logo": rsp.Brand.Logo,
		},
		"is_hot":  rsp.IsHot,
		"is_new":  rsp.IsNew,
		"on_sale": rsp.OnSale,
	}
	message.ResponseSuccess(ctx, data)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}

	_, err = global.GoodsServiceClient.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{Id: idInt})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	message.ResponseSuccess(ctx, nil)
}

func Stocks(ctx *gin.Context) {
	id := ctx.Param("id")
	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}

	//TODO 商品的库存

	return
}

func Update(ctx *gin.Context) {
	form := new(forms.GoodsForm)
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
	_, err = global.GoodsServiceClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              idInt,
		Name:            form.Name,
		GoodsSn:         form.GoodsSn,
		MarketPrice:     form.MarketPrice,
		ShopPrice:       form.ShopPrice,
		GoodsBrief:      form.GoodsBrief,
		ShipFree:        *form.ShipFree,
		Images:          form.Images,
		DescImages:      form.DescImages,
		GoodsFrontImage: form.GoodsFrontImage,
		CategoryId:      form.CategoryId,
		BrandId:         form.BrandId,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	message.ResponseSuccess(ctx, nil)
}

func UpdateStatus(ctx *gin.Context) {
	form := new(forms.GoodsStatus)
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
	_, err = global.GoodsServiceClient.UpdateGoodsStatus(context.Background(), &proto.GoodsStatusInfo{
		Id:     idInt,
		IsNew:  *form.IsNew,
		IsHot:  *form.IsHot,
		OnSale: *form.OnSale,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	message.ResponseSuccess(ctx, nil)
}
