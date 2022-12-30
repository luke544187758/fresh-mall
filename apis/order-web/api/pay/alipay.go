package pay

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"luke544187758/order-web/global"
	"luke544187758/order-web/proto"
	"luke544187758/order-web/settings"
	"net/http"
)

func Notify(ctx *gin.Context) {
	// 支付宝回调通知
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

	noti, _ := client.GetTradeNotification(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	_, err = global.OrderServiceClient.UpdateOrderStatus(context.Background(), &proto.OrderStatus{
		OrderSn: noti.OutTradeNo,
		Status:  string(noti.TradeStatus),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}

	ctx.String(http.StatusOK, "success")
}
