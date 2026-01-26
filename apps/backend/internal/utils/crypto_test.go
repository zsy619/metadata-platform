package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSM3EncryptionAndComparison(t *testing.T) {
	password := "user_password_123"
	salt := GenerateSalt()
	
	encrypted := EncryptPasswordSM3(password, salt)
	assert.Len(t, encrypted, 64) // SM3 hash is 256 bits = 64 hex chars

	// Correct comparison
	assert.True(t, ComparePasswordSM3(encrypted, password, salt))

	// Wrong password
	assert.False(t, ComparePasswordSM3(encrypted, "wrong_password", salt))

	// Wrong salt
	assert.False(t, ComparePasswordSM3(encrypted, password, "different_salt"))
}

func TestGenerateSalt_Length(t *testing.T) {
	salt := GenerateSalt()
	// Current requirement is 64 hex characters (32 random bytes)
	assert.Len(t, salt, 64)

	// Uniqueness (basic check)
	salt2 := GenerateSalt()
	assert.NotEqual(t, salt, salt2)
}

func TestPasswordHashingWrapper(t *testing.T) {
	password := "secure_pass"
	hash, salt, err := HashPassword(password)
	assert.NoError(t, err)
	assert.Len(t, hash, 64)
	assert.Len(t, salt, 64)

	assert.True(t, CheckPasswordHash(password, hash, salt))
	assert.False(t, CheckPasswordHash("wrong", hash, salt))
}

// BenchmarkGenerateSalt 测试盐值生成性能
func BenchmarkGenerateSalt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateSalt()
	}
}

// BenchmarkEncryptPasswordSM3 测试 SM3 加密性能
func BenchmarkEncryptPasswordSM3(b *testing.B) {
	password := "test_password_for_benchmarking_purposes"
	salt := "this_is_a_sixty_four_character_long_salt_string_exactly_64_chars!!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EncryptPasswordSM3(password, salt)
	}
}

// BenchmarkComparePasswordSM3 测试密码比对性能
func BenchmarkComparePasswordSM3(b *testing.B) {
	password := "test_password"
	salt := GenerateSalt()
	hashed := EncryptPasswordSM3(password, salt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ComparePasswordSM3(hashed, password, salt)
	}
}
