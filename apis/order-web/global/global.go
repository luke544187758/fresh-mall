package global

import (
	ut "github.com/go-playground/universal-translator"
	"luke544187758/order-web/proto"
)

var (
	Trans ut.Translator

	OrderServiceClient     proto.OrderClient
	GoodsServiceClient     proto.GoodsClient
	InventoryServiceClient proto.InventoryClient
)
