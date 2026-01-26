package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSalt 生成随机盐值
func GenerateSalt() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "static_salt_fallback_long_enough_to_be_secure_even_if_not_random"
	}
	return hex.EncodeToString(b)
}

// HashPassword 对密码进行哈希加密（使用 SM3 和随机生成的盐）
// 注意：为了向后兼容和接口简洁，这里假设外部会处理盐的持久化
// 但为了符合 SM3 规范，我们应该显式处理盐
func HashPassword(password string) (string, string, error) {
	salt := GenerateSalt()
	hash := EncryptPasswordSM3(password, salt)
	return hash, salt, nil
}

// CheckPasswordHash 验证密码是否与哈希值匹配
func CheckPasswordHash(password, hash, salt string) bool {
	return ComparePasswordSM3(hash, password, salt)
}
