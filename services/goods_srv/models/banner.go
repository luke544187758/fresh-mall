package models

import (
	"time"
)

type Banner struct {
	ID         int32     `db:"id"`
	Index      int32     `db:"index"`
	Image      string    `db:"image"`
	Url        string    `db:"url"`
	IsDeleted  bool      `db:"is_deleted"`
	AddTime    time.Time `db:"add_time"`
	UpdateTime time.Time `db:"update_time"`
}
