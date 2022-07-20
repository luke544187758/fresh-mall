package models

import (
	"time"
)

type Goods struct {
	ClickNum        int32     `db:"click_num"`
	SoldNum         int32     `db:"sold_num"`
	FavNum          int32     `db:"fav_num"`
	ID              int64     `db:"id"`
	CategoryId      int64     `db:"category_id"`
	BrandId         int64     `db:"brand_id"`
	MarketPrice     float32   `db:"market_price"`
	ShopPrice       float32   `db:"shop_price"`
	GoodsSn         string    `db:"goods_sn"`
	Name            string    `db:"name"`
	GoodsFrontImage string    `db:"goods_front_image"`
	GoodsBrief      string    `db:"goods_brief"`
	Images          []byte    `db:"images"`
	DescImages      []byte    `db:"desc_images"`
	ShipFree        bool      `db:"ship_free"`
	OnSale          bool      `db:"on_sale"`
	IsDeleted       bool      `db:"is_deleted"`
	IsNew           bool      `db:"is_new"`
	IsHot           bool      `db:"is_hot"`
	AddTime         time.Time `db:"add_time"`
	UpdateTime      time.Time `db:"update_time"`
}

type GoodsStatus struct {
	ID     int64 `db:"id"`
	IsNew  bool  `db:"is_new"`
	IsHot  bool  `db:"is_hot"`
	OnSale bool  `db:"on_sale"`
}
