package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"luke544187758/order-web/settings"
	"time"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	ID          int64  `json:"id"`
	AuthorityId int32  `json:"authority_id"`
	NickName    string `json:"nick_name"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(uid int64, role int32, nickname string) (string, error) {
	// 创建一个我们自己的声明的数据
	c := MyClaims{
		uid,
		role,
		nickname, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(settings.Conf.JWTConfig.JWTExpire) * time.Hour).Unix(), // 过期时间
			Issuer: "user_app", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString([]byte(settings.Conf.Secret))
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(settings.Conf.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
