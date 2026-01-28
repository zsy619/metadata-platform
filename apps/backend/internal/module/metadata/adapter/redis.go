package adapter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisExtractor Redis元数据提取器
type RedisExtractor struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisExtractor 创建Redis元数据提取器
func NewRedisExtractor(dsn string) (*RedisExtractor, error) {
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	return &RedisExtractor{
		client: client,
		ctx:    ctx,
	}, nil
}

// TestConnection 测试连接
func (e *RedisExtractor) TestConnection() error {
	ctx, cancel := context.WithTimeout(e.ctx, 5*time.Second)
	defer cancel()
	return e.client.Ping(ctx).Err()
}

func (e *RedisExtractor) GetSchemas() ([]string, error) {
	// Redis 默认返回单数据库标识
	return []string{"0"}, nil
}

// GetTables 获取键空间信息 (Redis 没有表的概念,返回数据库信息)
func (e *RedisExtractor) GetTables(schema string) ([]TableInfo, error) {
	// Redis 没有表的概念,可以返回键的统计信息
	info, err := e.client.Info(e.ctx, "keyspace").Result()
	if err != nil {
		return nil, err
	}
	_ = info

	// 简化实现:返回当前数据库的键数量信息
	dbSize, err := e.client.DBSize(e.ctx).Result()
	if err != nil {
		return nil, err
	}

	return []TableInfo{
		{
			Name:    "redis_keys",
			Comment: "Redis键值对存储 (总键数: " + string(rune(dbSize)) + ")",
		},
	}, nil
}

// GetViews Redis 不支持视图
func (e *RedisExtractor) GetViews(schema string) ([]ViewInfo, error) {
	return []ViewInfo{}, nil
}

// GetColumns Redis 不支持固定列结构
func (e *RedisExtractor) GetColumns(schema, table string) ([]ColumnInfo, error) {
	// Redis 是键值存储,没有固定的列结构
	// 可以返回一些通用的字段说明
	return []ColumnInfo{
		{
			Name:    "key",
			Type:    "string",
			Comment: "Redis 键名",
		},
		{
			Name:    "value",
			Type:    "string/list/set/zset/hash",
			Comment: "Redis 值 (支持多种数据类型)",
		},
		{
			Name:    "ttl",
			Type:    "int",
			Comment: "过期时间 (秒)",
		},
	}, nil
}

// GetIndexes Redis 不支持索引
func (e *RedisExtractor) GetIndexes(schema, table string) ([]IndexInfo, error) {
	return []IndexInfo{}, nil
}

// PreviewData 预览键值数据
func (e *RedisExtractor) PreviewData(schema, table string, limit int) ([]map[string]interface{}, error) {
	// 使用 SCAN 命令获取键
	var cursor uint64
	var keys []string
	var err error

	for len(keys) < limit {
		var batch []string
		batch, cursor, err = e.client.Scan(e.ctx, cursor, "*", int64(limit)).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, batch...)
		if cursor == 0 {
			break
		}
	}

	// 限制返回数量
	if len(keys) > limit {
		keys = keys[:limit]
	}

	var result []map[string]interface{}
	for _, key := range keys {
		// 获取键的类型
		keyType, err := e.client.Type(e.ctx, key).Result()
		if err != nil {
			continue
		}

		// 获取 TTL
		ttl, _ := e.client.TTL(e.ctx, key).Result()

		// 获取值 (简化处理,只获取字符串类型)
		var value interface{}
		switch keyType {
		case "string":
			value, _ = e.client.Get(e.ctx, key).Result()
		case "list":
			value, _ = e.client.LRange(e.ctx, key, 0, 10).Result()
		case "set":
			value, _ = e.client.SMembers(e.ctx, key).Result()
		case "zset":
			value, _ = e.client.ZRange(e.ctx, key, 0, 10).Result()
		case "hash":
			value, _ = e.client.HGetAll(e.ctx, key).Result()
		default:
			value = "unsupported type"
		}

		result = append(result, map[string]interface{}{
			"key":   key,
			"type":  keyType,
			"value": value,
			"ttl":   ttl.Seconds(),
		})
	}

	if len(result) == 0 {
		return nil, errors.New("Redis 数据库为空或无法访问")
	}

	return result, nil
}
// GetQueryColumns 获取查询结果的列信息
func (e *RedisExtractor) GetQueryColumns(query string, params []interface{}) ([]ColumnInfo, error) {
	return nil, fmt.Errorf("method GetQueryColumns not implemented for this adapter")
}

// Close 关闭连接
func (e *RedisExtractor) Close() error {
	return e.client.Close()
}
