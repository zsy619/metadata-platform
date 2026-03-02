package handlers

import (
	"crypto/x509/pkix"
	"testing"
	"time"

	"golang.org/x/crypto/ocsp"
)

// TestCRLCache 测试 CRL 缓存
func TestCRLCache(t *testing.T) {
	cache := NewCRLCache(1 * time.Hour)

	t.Run("Set and Get", func(t *testing.T) {
		url := "http://example.com/crl.pem"
		crl := &pkix.CertificateList{}

		cache.Set(url, crl)
		got, ok := cache.Get(url)
		if !ok {
			t.Error("应该获取到缓存的 CRL")
		}
		if got != crl {
			t.Error("获取的 CRL 与设置的不一致")
		}
	})

	t.Run("Get Non-existent", func(t *testing.T) {
		_, ok := cache.Get("http://non-existent.com/crl.pem")
		if ok {
			t.Error("不应该获取到不存在的 CRL")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		url := "http://example.com/crl2.pem"
		crl := &pkix.CertificateList{}

		cache.Set(url, crl)
		cache.Delete(url)
		_, ok := cache.Get(url)
		if ok {
			t.Error("删除后不应该获取到 CRL")
		}
	})
}

// TestOCSPCache 测试 OCSP 缓存
func TestOCSPCache(t *testing.T) {
	cache := NewOCSPCache(1 * time.Hour)

	t.Run("Set and Get", func(t *testing.T) {
		url := "http://example.com/ocsp"
		response := &ocsp.Response{
			Status: 0,
			ThisUpdate: time.Now(),
			NextUpdate: time.Now().Add(24 * time.Hour),
		}

		cache.Set(url, response)
		got, ok := cache.Get(url)
		if !ok {
			t.Error("应该获取到缓存的 OCSP 响应")
		}
		if got.Status != response.Status {
			t.Error("获取的 OCSP 响应状态不一致")
		}
	})

	t.Run("Get Non-existent", func(t *testing.T) {
		_, ok := cache.Get("http://non-existent.com/ocsp")
		if ok {
			t.Error("不应该获取到不存在的 OCSP 响应")
		}
	})
}

// TestSAMLAssertionCache 测试 SAML 断言缓存
func TestSAMLAssertionCache(t *testing.T) {
	cache := NewSAMLAssertionCache(3600) // 1 小时 = 3600 秒

	t.Run("IsDuplicate", func(t *testing.T) {
		assertionID := "test-assertion-1"

		// 第一次检查应该不是重复
		if cache.IsDuplicate(assertionID) {
			t.Error("新的断言 ID 不应该是重复的")
		}

		// 添加断言
		cache.Add(assertionID)

		// 第二次检查应该是重复
		if !cache.IsDuplicate(assertionID) {
			t.Error("已存在的断言 ID 应该是重复的")
		}
	})

	t.Run("Multiple Assertions", func(t *testing.T) {
		ids := []string{"id1", "id2", "id3"}

		for _, id := range ids {
			if cache.IsDuplicate(id) {
				t.Errorf("断言 ID %s 不应该是重复的", id)
			}
			cache.Add(id)
		}

		// 验证所有 ID 都是重复的
		for _, id := range ids {
			if !cache.IsDuplicate(id) {
				t.Errorf("断言 ID %s 应该是重复的", id)
			}
		}
	})
}

// BenchmarkCRLCache 基准测试：CRL 缓存性能
func BenchmarkCRLCache(b *testing.B) {
	cache := NewCRLCache(1 * time.Hour)
	url := "http://example.com/crl.pem"
	crl := &pkix.CertificateList{}

	// 预填充缓存
	cache.Set(url, crl)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(url)
	}
}

// BenchmarkOCSPCache 基准测试：OCSP 缓存性能
func BenchmarkOCSPCache(b *testing.B) {
	cache := NewOCSPCache(1 * time.Hour)
	url := "http://example.com/ocsp"
	response := &ocsp.Response{Status: 0}

	// 预填充缓存
	cache.Set(url, response)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(url)
	}
}

// BenchmarkSAMLAssertionCache 基准测试：断言缓存性能
func BenchmarkSAMLAssertionCache(b *testing.B) {
	cache := NewSAMLAssertionCache(3600)
	assertionID := "test-assertion"

	// 预填充缓存
	cache.Add(assertionID)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.IsDuplicate(assertionID)
	}
}
