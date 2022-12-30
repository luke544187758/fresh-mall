package message

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserNotAuth
	CodeServerBusy
	CodeNotFound
	CodeAlreadyExists

	CodeInvalidAuth
	CodeNeedLogin
	CodeNotAdmin

	CodeGoodsNotExists
	CodeGoodsAlreadyExists

	CodeInventoryShortage
)

type ResCode int64

var codeMsgMap = map[ResCode]string{
	CodeSuccess:            "success",
	CodeInvalidParam:       "请求参数错误",
	CodeUserNotAuth:        "用户没有权限登录",
	CodeServerBusy:         "服务器繁忙",
	CodeNotFound:           "没有找到记录",
	CodeAlreadyExists:      "记录已经存在",
	CodeInvalidAuth:        "请求中的token无效",
	CodeNeedLogin:          "需要登录",
	CodeNotAdmin:           "非管理员，没有权限",
	CodeGoodsNotExists:     "商品不存在",
	CodeGoodsAlreadyExists: "商品已经存在",
	CodeInventoryShortage:  "库存不足",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
