package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	// 生成随机字节
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// 如果出错，返回空字符串
		return ""
	}
	return hex.EncodeToString(bytes)[:length]
}