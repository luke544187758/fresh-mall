package message

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorUserNotAuth     = errors.New("没有权限登录")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
	ErrorInvalidID       = errors.New("无效的ID")
	ErrorUserNotLogin    = errors.New("用户未登录")
	ErrorNotAdmin        = errors.New("没有权限")
	ErrorInvalidCaptcha  = errors.New("验证码错误")
)
