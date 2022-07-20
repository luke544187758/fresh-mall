package forms

type GoodsForm struct {
	CategoryId      int64    `json:"category_id,string" binding:"required"`
	BrandId         int64    `json:"brand_id,string" binding:"required"`
	MarketPrice     float32  `json:"market_price" binding:"required,min=0"`
	ShopPrice       float32  `json:"shop_price" binding:"required,min=0"`
	GoodsSn         string   `json:"goods_sn" binding:"required,min=2,lt=20"`
	Name            string   `json:"name" binding:"required,min=2,max=100"`
	GoodsFrontImage string   `json:"goods_front_image" binding:"required,url"`
	GoodsBrief      string   `json:"goods_brief" binding:"required,min=3"`
	Images          []string `json:"images" binding:"required,min=1"`
	DescImages      []string `json:"desc_images" binding:"required,min=1"`
	ShipFree        *bool    `json:"ship_free" binding:"required"`
}

type GoodsStatus struct {
	IsNew  *bool `json:"is_new" binding:"required"`
	IsHot  *bool `json:"is_hot" binding:"required"`
	OnSale *bool `json:"on_sale" binding:"required"`
}
