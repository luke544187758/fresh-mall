package models

import (
	"database/sql"
	"time"
)

type Category struct {
	ID               int64         `db:"id"`
	ParentCategoryId sql.NullInt64 `db:"parent_category_id"`
	Level            int32         `db:"level"`
	Name             string        `db:"name"`
	IsDeleted        bool          `db:"is_deleted"`
	IsTab            bool          `db:"is_tab"`
	AddTime          time.Time     `db:"add_time"`
	UpdateTime       time.Time     `db:"update_time"`
}

func (c *Category) CategoryModelToDict() map[string]interface{} {
	return map[string]interface{}{
		"id":          c.ID,
		"name":        c.Name,
		"level":       c.Level,
		"parent":      c.ParentCategoryId,
		"is_tab":      c.IsTab,
		"add_time":    c.AddTime,
		"update_time": c.UpdateTime,
	}
}
