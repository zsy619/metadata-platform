package utils

import (
	"encoding/hex"

	"github.com/tjfoc/gmsm/sm3"
)

// EncryptPasswordSM3 使用国密 SM3 算法加盐加密密码
func EncryptPasswordSM3(password string, salt string) string {
	h := sm3.New()
	// 混合密码和盐值
	h.Write([]byte(password + salt))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

// ComparePasswordSM3 校验密码是否正确
func ComparePasswordSM3(hashedPassword, password, salt string) bool {
	return EncryptPasswordSM3(password, salt) == hashedPassword
}
