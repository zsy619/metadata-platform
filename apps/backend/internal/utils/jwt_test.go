package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken(t *testing.T) {
	userID := "user123"
	username := "testuser"
	isAdmin := true

	token, err := GenerateToken(userID, username, isAdmin)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, isAdmin, claims.IsAdmin)
}

func TestGenerateAndParseRefreshToken(t *testing.T) {
	userID := "user123"

	token, err := GenerateRefreshToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ParseRefreshToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.Subject)
}

func TestValidateToken(t *testing.T) {
	userID := "user123"
	username := "testuser"
	isAdmin := false

	token, _ := GenerateToken(userID, username, isAdmin)
	assert.True(t, ValidateToken(token))

	assert.False(t, ValidateToken("invalid.token.string"))
}

func TestExpiredToken(t *testing.T) {
	// Temporarily set a very short expiry for testing if possible, 
	// but TokenExpireIn is a package variable.
	oldExpiry := TokenExpireIn
	TokenExpireIn = -1 * time.Second
	defer func() { TokenExpireIn = oldExpiry }()

	token, _ := GenerateToken("123", "test", false)
	_, err := ParseToken(token)
	assert.Error(t, err)
}

func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateToken("user123", "testuser", true)
	}
}

func BenchmarkParseToken(b *testing.B) {
	token, _ := GenerateToken("user123", "testuser", true)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseToken(token)
	}
}
