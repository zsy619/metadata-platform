package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT配置
var (
	SecretKey            = []byte("your-secret-key") // 实际使用时应从配置文件读取
	TokenExpireIn        = 24 * time.Hour            // access token过期时间
	RefreshTokenExpireIn = 7 * 24 * time.Hour        // refresh token过期时间
	RefreshThreshold     = 30 * time.Minute          // 自动续约阈值：剩余时间小于此值时续约
)

// JWTClaims 自定义JWT声明结构
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID string, username string, isAdmin bool) (string, error) {
	// 创建声明
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整的编码后的字符串token
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpireIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*JWTClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// 验证token并获取声明
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateToken 验证token有效性
func ValidateToken(tokenString string) bool {
	_, err := ParseToken(tokenString)
	return err == nil
}

// ParseRefreshToken 解析刷新令牌
func ParseRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid refresh token")
}

// ShouldRefresh 检查令牌是否由于接近过期而需要刷新
func ShouldRefresh(claims *JWTClaims) bool {
	if claims == nil {
		return false
	}
	expiresAt := claims.ExpiresAt.Time
	return time.Until(expiresAt) < RefreshThreshold
}
