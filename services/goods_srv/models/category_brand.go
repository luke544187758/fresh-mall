package models

import "time"

type CategoryBrand struct {
	ID         int64     `db:"id"`
	CategoryId int64     `db:"category_id"`
	BrandId    int64     `db:"brand_id"`
	IsDeleted  bool      `db:"is_deleted"`
	AddTime    time.Time `db:"add_time"`
	UpdateTime time.Time `db:"update_time"`
}
