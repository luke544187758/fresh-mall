package models

import "database/sql"

type Address struct {
	Id           int64        `db:"id"`
	UserId       int64        `db:"user_id"`
	Province     string       `db:"province"`
	City         string       `db:"city"`
	District     string       `db:"district"`
	Address      string       `db:"address"`
	SignerName   string       `db:"signer_name"`
	SignerMobile string       `db:"signer_mobile"`
	AddTime      sql.NullTime `db:"add_time"`
	UpdateTime   sql.NullTime `db:"update_time"`
	IsDeleted    bool         `db:"is_deleted"`
}
