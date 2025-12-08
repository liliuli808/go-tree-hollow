package utils

import (
	"math/rand"
	"time"
)

// GenerateCode 生成指定长度的数字验证码
func GenerateCode(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))

	const digits = "0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return string(code)
}
