package forms

type ShopCartItem struct {
	GoodsId int64 `json:"goods_id,string" binding:"required"`
	Nums    int32 `json:"nums" binding:"required,min=1"`
}

type ShopCartItemUpdate struct {
	Nums    int32 `json:"nums" binding:"required,min=1"`
	Checked *bool `json:"checked"`
}
