package models

import "time"

type Inventory struct {
	ID         int64     `db:"id"`
	Goods      int64     `db:"goods"`
	Stocks     int32     `db:"stocks"`
	Version    int32     `db:"version"`
	IsDeleted  bool      `db:"is_deleted"`
	AddTime    time.Time `db:"add_time"`
	UpdateTime time.Time `db:"update_time"`
}

type Deduct struct {
	Goods  int64 `db:"goods"`
	Stocks int32 `db:"stocks"`
}
