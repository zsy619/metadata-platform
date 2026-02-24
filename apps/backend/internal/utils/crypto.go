package utils

import (
	"encoding/hex"
	"fmt"

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
	pass := EncryptPasswordSM3(password, salt)
	fmt.Printf("ComparePasswordSM3===>%s %s %s %s \n", password, salt, pass, hashedPassword)
	return pass == hashedPassword
}
