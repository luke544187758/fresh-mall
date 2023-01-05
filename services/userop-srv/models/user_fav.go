package models

import "database/sql"

type UserFav struct {
	UserId  int64        `db:"user_id"`
	GoodsId int64        `db:"goods_id"`
	AddTime sql.NullTime `db:"add_time"`
}
