package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

const secret = "没有蛀牙"

//encryptPassword 对密码进行哈希编码
func EncryptPassword(oPassword string) string {
	h := sha1.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
