package models

import "time"

type User struct {
	ID           int64     `db:"id"`           // id
	Role         int32     `db:"role"`         // 角色
	Mobile       string    `db:"mobile"`       // 手机号
	Password     string    `db:"password"`     // 密码
	NickName     string    `db:"nick_name"`    // 昵称
	Avatar       string    `db:"avatar"`       // 头像
	Address      string    `db:"address"`      // 地址
	Introduction string    `db:"introduction"` // 简介
	Gender       string    `db:"gender"`       // 性别
	Birthday     time.Time `db:"birthday"`     // 生日
}
