package models

import (
	"database/sql"
)

type Message struct {
	Id          int64          `db:"id"`
	UserId      int64          `db:"user_id"`
	MessageType int32          `db:"message_type"`
	Subject     sql.NullString `db:"subject"`
	Message     sql.NullString `db:"message"`
	File        sql.NullString `db:"file"`
	AddTime     sql.NullTime   `db:"add_time"`
	UpdateTime  sql.NullTime   `db:"update_time"`
	IsDeleted   bool           `db:"is_deleted"`
}
