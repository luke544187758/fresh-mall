package message

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeUserNotAuth
	CodeInvalidPassword
	CodeServerBusy

	CodeInvalidAuth
	CodeNeedLogin
	CodeNotAdmin
	CodeInvalidCaptcha
)

type ResCode int64

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeUserNotAuth:     "用户没有权限登录",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务器繁忙",
	CodeInvalidAuth:     "请求中的token无效",
	CodeNeedLogin:       "需要登录",
	CodeNotAdmin:        "非管理员，没有权限",
	CodeInvalidCaptcha:  "验证码错误",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
