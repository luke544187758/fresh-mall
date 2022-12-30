package order

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	_ "github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"luke544187758/order-web/api"
	"luke544187758/order-web/forms"
	"luke544187758/order-web/global"
	"luke544187758/order-web/message"
	"luke544187758/order-web/pkg/jwt"
	"luke544187758/order-web/proto"
	"luke544187758/order-web/settings"
	"net/http"
	"strconv"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("user_id")
	claims, _ := ctx.Get("claims")
	request := new(proto.OrderFilterRequest)

	//如果是管理员用户，则返回所有的订单
	model := claims.(*jwt.MyClaims)
	if model.AuthorityId == 1 {
		request.UserId = int64(userId.(uint))
	}

	pageSize := ctx.DefaultQuery("page_size", "0")
	pageSizeInt, _ := strconv.ParseInt(pageSize, 64, 10)
	request.PerNum = int32(pageSizeInt)

	page := ctx.DefaultQuery("page", "0")
	pageInt, _ := strconv.ParseInt(page, 64, 10)
	request.Page = int32(pageInt)

	resp, err := global.OrderServiceClient.OrderList(context.Background(), request)
	if err != nil {
		zap.L().Error("Get order list failed", zap.Error(err))
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	data := make([]interface{}, 0)
	for _, item := range resp.Data {
		tmp := map[string]interface{}{}
		tmp["id"] = item.Id
		tmp["user_id"] = item.UserId
		tmp["order_sn"] = item.OrderSn
		tmp["pay_type"] = item.PayType
		tmp["status"] = item.Status
		tmp["remark"] = item.Remark
		tmp["trade_no"] = item.TradeNo
		tmp["order_mount"] = item.OrderMount
		tmp["address"] = item.Address
		tmp["name"] = item.Name
		tmp["mobile"] = item.Mobile
		tmp["pay_time"] = item.PayTime
		tmp["add_time"] = item.AddTime
		data = append(data, tmp)
	}
	resMap := gin.H{
		"total": resp.Total,
		"data":  data,
	}
	message.ResponseSuccess(ctx, resMap)
}

func New(ctx *gin.Context) {
	form := new(forms.CreateOrderForm)
	if err := ctx.ShouldBindJSON(form); err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}
	userId, _ := ctx.Get("user_id")
	rsp, err := global.OrderServiceClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  int64(userId.(uint)),
		Address: form.Address,
		Mobile:  form.Mobile,
		Name:    form.Name,
		Remark:  form.Remark,
	})
	if err != nil {
		zap.L().Error("create a new order failed", zap.Error(err))
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 生成支付宝的支付url
	client, err := alipay.New(settings.Conf.AliPayConfig.AppId, settings.Conf.AliPayConfig.PrivateKey, false)
	if err != nil {
		zap.L().Error("An error occurred when creating an Alipay instance", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if err = client.LoadAliPayPublicKey(settings.Conf.AliPayConfig.AliPublicKey); err != nil {
		zap.L().Error("An error occurred when loading public key of Alipay", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	p := alipay.TradePagePay{}
	p.NotifyURL = settings.Conf.AliPayConfig.NotifyUrl
	p.ReturnURL = settings.Conf.AliPayConfig.ReturnUrl
	p.Subject = "生鲜商城 - " + rsp.OrderSn
	p.OutTradeNo = rsp.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(rsp.OrderMount), 'f', 2, 64)
	p.ProductCode = settings.Conf.AliPayConfig.ProductCode

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.L().Error("An error occurred in generating the payment URL", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	message.ResponseSuccess(ctx, gin.H{
		"id":         rsp.Id,
		"alipay_url": url.String(),
	})
}

func Detail(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		message.ResponseError(ctx, message.CodeInvalidParam)
		return
	}

	//如果是管理员用户，则返回所有的订单
	userId, _ := ctx.Get("user_id")
	claims, _ := ctx.Get("claims")
	request := &proto.OrderRequest{
		Id: id,
	}
	model := claims.(*jwt.MyClaims)
	if model.AuthorityId == 1 {
		request.UserId = int64(userId.(uint))
	}

	resp, err := global.OrderServiceClient.OrderDetail(context.Background(), request)
	if err != nil {
		zap.L().Error("get order detail failed", zap.Error(err))
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"id":          resp.OrderInfo.Id,
		"pay_type":    resp.OrderInfo.PayType,
		"user_id":     resp.OrderInfo.UserId,
		"order_sn":    resp.OrderInfo.OrderSn,
		"status":      resp.OrderInfo.Status,
		"remark":      resp.OrderInfo.Remark,
		"trade_no":    resp.OrderInfo.TradeNo,
		"order_mount": resp.OrderInfo.OrderMount,
		"address":     resp.OrderInfo.Address,
		"name":        resp.OrderInfo.Name,
		"mobile":      resp.OrderInfo.Mobile,
		"pay_time":    resp.OrderInfo.PayTime,
		"add_time":    resp.OrderInfo.AddTime,
	}

	goods := make([]interface{}, 0)
	for _, item := range resp.Data {
		tmp := gin.H{
			"id":    item.GoodsId,
			"name":  item.GoodsName,
			"image": item.GoodsImage,
			"price": item.GoodsPrice,
			"nums":  item.Nums,
		}
		goods = append(goods, tmp)
	}
	reMap["goods"] = goods

	// 生成支付宝的支付url
	client, err := alipay.New(settings.Conf.AliPayConfig.AppId, settings.Conf.AliPayConfig.PrivateKey, false)
	if err != nil {
		zap.L().Error("An error occurred when creating an Alipay instance", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if err = client.LoadAliPayPublicKey(settings.Conf.AliPayConfig.AliPublicKey); err != nil {
		zap.L().Error("An error occurred when loading public key of Alipay", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	p := alipay.TradePagePay{}
	p.NotifyURL = settings.Conf.AliPayConfig.NotifyUrl
	p.ReturnURL = settings.Conf.AliPayConfig.ReturnUrl
	p.Subject = "生鲜商城 - " + resp.OrderInfo.OrderSn
	p.OutTradeNo = resp.OrderInfo.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(resp.OrderInfo.OrderMount), 'f', 2, 64)
	p.ProductCode = settings.Conf.AliPayConfig.ProductCode

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.L().Error("An error occurred in generating the payment URL", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	reMap["alipay_url"] = url.String()

	message.ResponseSuccess(ctx, reMap)
}
