package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOrderSn(uid int64) string {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(90) + 10
	return fmt.Sprintf("%s%d%d", time.Now().Format("20060102150405"), uid, num)
}
