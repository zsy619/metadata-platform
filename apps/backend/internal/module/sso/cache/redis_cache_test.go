package cache

import (
	"testing"
	"time"
)

// TestRedisConfig 测试 Redis 配置
func TestRedisConfig(t *testing.T) {
	config := &RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 100,
	}

	if config.Addr != "localhost:6379" {
		t.Error("Redis 地址配置错误")
	}
	if config.PoolSize != 100 {
		t.Error("连接池大小配置错误")
	}
}

// TestDistributedCRLCache_Skip 跳过实际 Redis 测试（需要 Redis 服务）
func TestDistributedCRLCache_Skip(t *testing.T) {
	t.Skip("跳过实际 Redis 测试，需要 Redis 服务")
}

// TestDistributedOCSPCache_Skip 跳过实际 Redis 测试
func TestDistributedOCSPCache_Skip(t *testing.T) {
	t.Skip("跳过实际 Redis 测试，需要 Redis 服务")
}

// TestDistributedAssertionCache_Skip 跳过实际 Redis 测试
func TestDistributedAssertionCache_Skip(t *testing.T) {
	t.Skip("跳过实际 Redis 测试，需要 Redis 服务")
}

// TestCacheStats 测试缓存统计
func TestCacheStats(t *testing.T) {
	stats := &CacheStats{
		Connected: true,
		Prefix:    "test:sso:",
		TTL:       1 * time.Hour,
		Info:      "test info",
	}

	if !stats.Connected {
		t.Error("连接状态错误")
	}
	if stats.Prefix != "test:sso:" {
		t.Error("前缀错误")
	}
	if stats.TTL != 1*time.Hour {
		t.Error("TTL 错误")
	}
}

// TestParseCRLItem 测试 CRL 项解析
func TestParseCRLItem(t *testing.T) {
	t.Run("Nil Item", func(t *testing.T) {
		_, err := ParseCRLItem(nil)
		if err == nil {
			t.Error("nil 项应该返回错误")
		}
	})

	t.Run("Empty CRL Bytes", func(t *testing.T) {
		item := &CRLItem{
			URL:       "http://example.com/crl.pem",
			CRLBytes:  []byte{},
			Timestamp: time.Now(),
		}
		_, err := ParseCRLItem(item)
		if err == nil {
			t.Error("空 CRL 字节应该返回错误")
		}
	})
}

// BenchmarkCacheStats 基准测试：缓存统计
func BenchmarkCacheStats(b *testing.B) {
	stats := &CacheStats{
		Connected: true,
		Prefix:    "test:sso:",
		TTL:       1 * time.Hour,
		Info:      "test info",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = stats.Connected
		_ = stats.Prefix
		_ = stats.TTL
	}
}

// BenchmarkCRLItem 基准测试：CRL 项创建
func BenchmarkCRLItem(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		item := &CRLItem{
			URL:       "http://example.com/crl.pem",
			CRLBytes:  []byte("test crl data"),
			Timestamp: time.Now(),
		}
		_ = item
	}
}

// BenchmarkOCSPItem 基准测试：OCSP 项创建
func BenchmarkOCSPItem(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		item := &OCSPItem{
			CertURL:    "http://example.com/ocsp",
			Response:   []byte("test ocsp response"),
			Status:     0,
			ThisUpdate: time.Now(),
			NextUpdate: time.Now().Add(24 * time.Hour),
			Timestamp:  time.Now(),
		}
		_ = item
	}
}
