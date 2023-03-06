package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 小寫
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tmpStr := h.Sum(nil)
	return hex.EncodeToString(tmpStr)
}

// 大寫
func Md5Decode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// 加密
func MakePasssword(plainPwd, salt string) string {
	return Md5Encode(plainPwd + salt)
}

// 解密
func ValidPassword(plainPwd, salt string, password string) bool {
	return Md5Encode(plainPwd+salt) == password
}
