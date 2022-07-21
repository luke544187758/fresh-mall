package models

import (
	"database/sql/driver"
	"time"
)

type ShoppingCart struct {
	ID         int64     `db:"id"`
	User       int64     `db:"user"`
	Goods      int64     `db:"goods"`
	Nums       int32     `db:"nums"`
	Checked    bool      `db:"checked"`
	IsDeleted  bool      `db:"is_deleted"`
	AddTime    time.Time `db:"add_time"`
	UpdateTime time.Time `db:"update_time"`
}

type OrderInfo struct {
	ID           int64     `db:"id"`
	User         int64     `db:"user"`
	OrderMount   float32   `db:"order_mount"`
	IsDeleted    bool      `db:"is_deleted"`
	OrderSn      string    `db:"order_sn"`
	PayType      string    `db:"pay_type"`
	Status       string    `db:"status"`
	TradeNo      string    `db:"trade_no"`
	Address      string    `db:"address"`
	SignerName   string    `db:"signer_name"`
	SignerMobile string    `db:"signer_mobile"`
	Remark       string    `db:"remark"`
	PayTime      time.Time `db:"pay_time"`
	AddTime      time.Time `db:"add_time"`
	UpdateTime   time.Time `db:"update_time"`
}

type OrderGoods struct {
	ID         int64     `db:"id"`
	Order      int64     `db:"order"`
	Goods      int64     `db:"goods"`
	Nums       int32     `db:"nums"`
	GoodsPrice float32   `db:"goods_price"`
	IsDeleted  bool      `db:"is_deleted"`
	GoodsName  string    `db:"goods_name"`
	GoodsImage string    `db:"goods_image"`
	AddTime    time.Time `db:"add_time"`
	UpdateTime time.Time `db:"update_time"`
}

func (this *OrderGoods) Value() (driver.Value, error) {
	return []interface{}{
		this.ID, this.Order, this.Goods, this.Nums, this.GoodsPrice, this.IsDeleted, this.GoodsName, this.GoodsImage, this.AddTime, this.UpdateTime,
	}, nil
}
