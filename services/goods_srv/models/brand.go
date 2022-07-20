package models

import (
	"database/sql"
	"time"
)

type Brand struct {
	ID         int64          `db:"id"`
	Name       string         `db:"name"`
	Logo       sql.NullString `db:"logo"`
	IsDeleted  bool           `db:"is_deleted"`
	AddTime    time.Time      `db:"add_time"`
	UpdateTime time.Time      `db:"update_time"`
}
