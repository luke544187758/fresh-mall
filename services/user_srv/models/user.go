package models

import (
	"database/sql"
)

type User struct {
	ID           int64          `db:"id"`           // id
	Role         int32          `db:"role"`         // 角色
	Mobile       string         `db:"mobile"`       // 手机号
	Password     string         `db:"password"`     // 密码
	NickName     string         `db:"nick_name"`    // 昵称
	Avatar       sql.NullString `db:"avatar"`       // 头像
	Address      sql.NullString `db:"address"`      // 地址
	Introduction sql.NullString `db:"introduction"` // 简介
	Gender       sql.NullString `db:"gender"`       // 性别
	Birthday     sql.NullTime   `db:"birthday"`     // 生日
}
