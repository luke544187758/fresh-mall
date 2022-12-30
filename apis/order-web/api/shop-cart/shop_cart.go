package shop_cart

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"luke544187758/order-web/api"
	"luke544187758/order-web/forms"
	"luke544187758/order-web/global"
	"luke544187758/order-web/message"
	"luke544187758/order-web/proto"
	"strconv"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")
	rsp, err := global.OrderServiceClient.CartItemList(context.Background(), &proto.UserInfo{Id: userId.(int64)})
	if err != nil {
		zap.L().Error("[List] query cart list failed", zap.Error(err))
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ids := make([]int64, 0)
	for _, item := range rsp.Data {
		ids = append(ids, item.GoodsId)
	}
	if len(ids) == 0 {
		message.ResponseSuccess(ctx, gin.H{
			"total": 0,
		})
		return
	}

	// 请求商品服务获取商品信息
	goodsRes, err := global.GoodsServiceClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: ids})
	if err != nil {
		zap.L().Error("[List] batch query goods information failed", zap.Error(err))
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	data := make(map[string]interface{})
	items := make([]interface{}, 0)
	data["total"] = rsp.Total

	for _, item := range rsp.Data {
		for _, v := range goodsRes.Data {
			if v.Id == item.GoodsId {
				items = append(items, map[string]interface{}{
					"id":          fmt.Sprintf("%d", item.Id),
					"goods_id":    fmt.Sprintf("%d", item.GoodsId),
					"goods_name":  v.Name,
					"goods_image": v.GoodsFrontImage,
					"goods_price": v.ShopPrice,
					"nums":        item.Nums,
					"checked":     item.Checked,
				})
			}
		}
	}
	data["items"] = items
	message.ResponseSuccess(ctx, data)
}

func New(ctx *gin.Context) {
	form := new(forms.ShopCartItem)
	if err := ctx.ShouldBindJSON(form); err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}
	// 检查商品是否存在
	_, err := global.GoodsServiceClient.GetGoodsDetail(context.Background(), &proto.GoodsInfoRequest{Id: form.GoodsId})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 价差商品库存
	invRes, err := global.InventoryServiceClient.InventoryDetail(context.Background(), &proto.GoodsInventoryInfo{GoodsId: form.GoodsId})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	if invRes.Num < form.Nums {
		message.ResponseError(ctx, message.CodeInventoryShortage)
		return
	}

	userId, _ := ctx.Get("user_id")
	rsp, err := global.OrderServiceClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      0,
		UserId:  userId.(int64),
		GoodsId: form.GoodsId,
		Nums:    form.Nums,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	message.ResponseSuccess(ctx, gin.H{"id": fmt.Sprintf("%d", rsp.Id)})
}

func Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}

	userId, _ := ctx.Get("user_id")
	_, err = global.OrderServiceClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  userId.(int64),
		GoodsId: id,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	message.ResponseSuccess(ctx, nil)
}

func Update(ctx *gin.Context) {
	form := new(forms.ShopCartItemUpdate)
	if err := ctx.ShouldBindJSON(form); err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}

	userId, _ := ctx.Get("user_id")
	request := proto.CartItemRequest{
		UserId:  userId.(int64),
		GoodsId: id,
		Nums:    form.Nums,
		Checked: false,
	}
	if form.Checked != nil {
		request.Checked = *form.Checked
	}

	_, err = global.OrderServiceClient.UpdateCartItem(context.Background(), &request)
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	message.ResponseSuccess(ctx, nil)
}
