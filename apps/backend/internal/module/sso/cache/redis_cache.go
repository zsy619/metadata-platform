package cache

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string `json:"addr"`     // Redis 地址
	Password string `json:"password"` // 密码
	DB       int    `json:"db"`       // 数据库
	PoolSize int    `json:"pool_size"`// 连接池大小
}

// DistributedCache 分布式缓存基类
type DistributedCache struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

// NewDistributedCache 创建分布式缓存
func NewDistributedCache(config *RedisConfig, prefix string, ttl time.Duration) (*DistributedCache, error) {
	if config == nil {
		return nil, errors.New("Redis 配置不能为空")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接 Redis 失败：%w", err)
	}

	return &DistributedCache{
		client: client,
		prefix: prefix,
		ttl:    ttl,
	}, nil
}

// Close 关闭连接
func (c *DistributedCache) Close() error {
	return c.client.Close()
}

// DistributedCRLCache 分布式 CRL 缓存
type DistributedCRLCache struct {
	*DistributedCache
}

// NewDistributedCRLCache 创建分布式 CRL 缓存
func NewDistributedCRLCache(config *RedisConfig, ttl time.Duration) (*DistributedCRLCache, error) {
	cache, err := NewDistributedCache(config, "sso:crl:", ttl)
	if err != nil {
		return nil, err
	}

	return &DistributedCRLCache{
		DistributedCache: cache,
	}, nil
}

// CRLItem CRL 缓存项
type CRLItem struct {
	URL       string    `json:"url"`
	CRLBytes  []byte    `json:"crl_bytes"`
	Timestamp time.Time `json:"timestamp"`
}

// Get 从 Redis 获取 CRL
func (c *DistributedCRLCache) Get(ctx context.Context, url string) (*CRLItem, error) {
	key := c.prefix + url
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, err
	}

	var item CRLItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("解析 CRL 缓存失败：%w", err)
	}

	// 检查是否过期
	if time.Since(item.Timestamp) > c.ttl {
		c.Delete(ctx, url)
		return nil, nil
	}

	return &item, nil
}

// Set 设置 CRL 到 Redis
func (c *DistributedCRLCache) Set(ctx context.Context, url string, crlBytes []byte) error {
	key := c.prefix + url
	item := CRLItem{
		URL:       url,
		CRLBytes:  crlBytes,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("序列化 CRL 失败：%w", err)
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}

// Delete 从 Redis 删除 CRL
func (c *DistributedCRLCache) Delete(ctx context.Context, url string) error {
	key := c.prefix + url
	return c.client.Del(ctx, key).Err()
}

// DistributedOCSPCache 分布式 OCSP 缓存
type DistributedOCSPCache struct {
	*DistributedCache
}

// NewDistributedOCSPCache 创建分布式 OCSP 缓存
func NewDistributedOCSPCache(config *RedisConfig, ttl time.Duration) (*DistributedOCSPCache, error) {
	cache, err := NewDistributedCache(config, "sso:ocsp:", ttl)
	if err != nil {
		return nil, err
	}

	return &DistributedOCSPCache{
		DistributedCache: cache,
	}, nil
}

// OCSPItem OCSP 缓存项
type OCSPItem struct {
	CertURL    string    `json:"cert_url"`
	Response   []byte    `json:"response"`
	Status     int       `json:"status"`
	ThisUpdate time.Time `json:"this_update"`
	NextUpdate time.Time `json:"next_update"`
	Timestamp  time.Time `json:"timestamp"`
}

// Get 从 Redis 获取 OCSP 响应
func (c *DistributedOCSPCache) Get(ctx context.Context, certURL string) (*OCSPItem, error) {
	key := c.prefix + certURL
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, err
	}

	var item OCSPItem
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("解析 OCSP 缓存失败：%w", err)
	}

	// 检查是否过期
	if time.Since(item.Timestamp) > c.ttl {
		c.Delete(ctx, certURL)
		return nil, nil
	}

	return &item, nil
}

// Set 设置 OCSP 响应到 Redis
func (c *DistributedOCSPCache) Set(ctx context.Context, certURL string, response []byte, status int, thisUpdate, nextUpdate time.Time) error {
	key := c.prefix + certURL
	item := OCSPItem{
		CertURL:    certURL,
		Response:   response,
		Status:     status,
		ThisUpdate: thisUpdate,
		NextUpdate: nextUpdate,
		Timestamp:  time.Now(),
	}

	data, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("序列化 OCSP 响应失败：%w", err)
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}

// Delete 从 Redis 删除 OCSP 响应
func (c *DistributedOCSPCache) Delete(ctx context.Context, certURL string) error {
	key := c.prefix + certURL
	return c.client.Del(ctx, key).Err()
}

// DistributedAssertionCache 分布式断言缓存
type DistributedAssertionCache struct {
	*DistributedCache
}

// NewDistributedAssertionCache 创建分布式断言缓存
func NewDistributedAssertionCache(config *RedisConfig, ttl time.Duration) (*DistributedAssertionCache, error) {
	cache, err := NewDistributedCache(config, "sso:assertion:", ttl)
	if err != nil {
		return nil, err
	}

	return &DistributedAssertionCache{
		DistributedCache: cache,
	}, nil
}

// IsDuplicate 检查断言 ID 是否重复
func (c *DistributedAssertionCache) IsDuplicate(ctx context.Context, assertionID string) (bool, error) {
	key := c.prefix + assertionID
	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// Add 添加断言 ID 到 Redis
func (c *DistributedAssertionCache) Add(ctx context.Context, assertionID string) error {
	key := c.prefix + assertionID
	return c.client.Set(ctx, key, "1", c.ttl).Err()
}

// Delete 从 Redis 删除断言 ID
func (c *DistributedAssertionCache) Delete(ctx context.Context, assertionID string) error {
	key := c.prefix + assertionID
	return c.client.Del(ctx, key).Err()
}

// LocalCRLCacheInterface 本地 CRL 缓存接口
type LocalCRLCacheInterface interface {
	Get(url string) (interface{}, bool)
	Set(url string, crl interface{})
	Delete(url string)
}

// HybridCRLCache 混合 CRL 缓存（Redis + 本地）
type HybridCRLCache struct {
	redisCache    *DistributedCRLCache
	localCache    LocalCRLCacheInterface
	useRedisFirst bool
}

// NewHybridCRLCache 创建混合 CRL 缓存
func NewHybridCRLCache(redisConfig *RedisConfig, localCache LocalCRLCacheInterface, ttl time.Duration, useRedisFirst bool) (*HybridCRLCache, error) {
	var redisCache *DistributedCRLCache
	var err error

	if redisConfig != nil {
		redisCache, err = NewDistributedCRLCache(redisConfig, ttl)
		if err != nil {
			// Redis 连接失败，降级为本地缓存
			redisCache = nil
		}
	}

	return &HybridCRLCache{
		redisCache:    redisCache,
		localCache:    localCache,
		useRedisFirst: useRedisFirst,
	}, nil
}

// Get 获取 CRL（混合缓存）
func (c *HybridCRLCache) Get(ctx context.Context, url string) (interface{}, bool) {
	// 优先从 Redis 获取
	if c.useRedisFirst && c.redisCache != nil {
		item, err := c.redisCache.Get(ctx, url)
		if err == nil && item != nil {
			// 同步到本地缓存
			c.localCache.Set(url, item)
			return item.CRLBytes, true
		}
	}

	// 从本地缓存获取
	if item, ok := c.localCache.Get(url); ok {
		return item, true
	}

	// Redis 降级
	if !c.useRedisFirst && c.redisCache != nil {
		item, err := c.redisCache.Get(ctx, url)
		if err == nil && item != nil {
			// 同步到本地缓存
			c.localCache.Set(url, item)
			return item.CRLBytes, true
		}
	}

	return nil, false
}

// Set 设置 CRL（混合缓存）
func (c *HybridCRLCache) Set(ctx context.Context, url string, crlBytes []byte) {
	// 设置本地缓存
	c.localCache.Set(url, crlBytes)

	// 设置 Redis 缓存
	if c.redisCache != nil {
		c.redisCache.Set(ctx, url, crlBytes)
	}
}

// Delete 删除 CRL（混合缓存）
func (c *HybridCRLCache) Delete(ctx context.Context, url string) {
	// 删除本地缓存
	c.localCache.Delete(url)

	// 删除 Redis 缓存
	if c.redisCache != nil {
		c.redisCache.Delete(ctx, url)
	}
}

// GetStats 获取缓存统计信息
func (c *DistributedCache) GetStats(ctx context.Context) (*CacheStats, error) {
	info, err := c.client.Info(ctx, "stats").Result()
	if err != nil {
		return nil, err
	}

	// 解析 Redis 统计信息（简化版）
	return &CacheStats{
		Connected: true,
		Prefix:    c.prefix,
		TTL:       c.ttl,
		Info:      info,
	}, nil
}

// CacheStats 缓存统计信息
type CacheStats struct {
	Connected bool   `json:"connected"`
	Prefix    string `json:"prefix"`
	TTL       time.Duration `json:"ttl"`
	Info      string `json:"info"`
}

// HealthCheck 健康检查
func (c *DistributedCache) HealthCheck(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// ParseCRLItem 解析 CRL 缓存项为 x509.CRL
func ParseCRLItem(item *CRLItem) (*x509.RevocationList, error) {
	if item == nil || len(item.CRLBytes) == 0 {
		return nil, errors.New("CRL 数据为空")
	}

	return x509.ParseRevocationList(item.CRLBytes)
}
