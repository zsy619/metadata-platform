package adapter

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBExtractor MongoDB元数据提取器
type MongoDBExtractor struct {
	client *mongo.Client
	dbName string
	ctx    context.Context
}

// NewMongoDBExtractor 创建MongoDB元数据提取器
func NewMongoDBExtractor(dsn string) (*MongoDBExtractor, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	// 从 DSN 中提取数据库名
	clientOptions := options.Client().ApplyURI(dsn)
	dbName := "test" // 默认数据库
	if clientOptions.Auth != nil && clientOptions.Auth.AuthSource != "" {
		dbName = clientOptions.Auth.AuthSource
	}

	return &MongoDBExtractor{
		client: client,
		dbName: dbName,
		ctx:    ctx,
	}, nil
}

// TestConnection 测试连接
func (e *MongoDBExtractor) TestConnection() error {
	ctx, cancel := context.WithTimeout(e.ctx, 5*time.Second)
	defer cancel()
	return e.client.Ping(ctx, nil)
}

func (e *MongoDBExtractor) GetSchemas() ([]string, error) {
	// 按照需求，只显示当前连接的数据库
	return []string{e.dbName}, nil
}

// GetTables 获取集合列表 (MongoDB 中的集合对应关系型数据库的表)
func (e *MongoDBExtractor) GetTables(schema string) ([]TableInfo, error) {
	if schema == "" {
		schema = e.dbName
	}

	db := e.client.Database(schema)
	collections, err := db.ListCollectionNames(e.ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var tables []TableInfo
	for _, collName := range collections {
		tables = append(tables, TableInfo{
			Name:    collName,
			Comment: "MongoDB Collection",
		})
	}
	return tables, nil
}

// GetViews 获取视图列表
func (e *MongoDBExtractor) GetViews(schema string) ([]ViewInfo, error) {
	if schema == "" {
		schema = e.dbName
	}

	db := e.client.Database(schema)
	collections, err := db.ListCollections(e.ctx, bson.M{"type": "view"})
	if err != nil {
		return nil, err
	}
	defer collections.Close(e.ctx)

	var views []ViewInfo
	for collections.Next(e.ctx) {
		var result bson.M
		if err := collections.Decode(&result); err != nil {
			continue
		}
		if name, ok := result["name"].(string); ok {
			views = append(views, ViewInfo{
				Name: name,
			})
		}
	}
	return views, nil
}

// GetColumns 获取字段信息 (通过采样文档推断)
func (e *MongoDBExtractor) GetColumns(schema, table string) ([]ColumnInfo, error) {
	if schema == "" {
		schema = e.dbName
	}

	db := e.client.Database(schema)
	collection := db.Collection(table)

	// 采样前100个文档来推断字段
	cursor, err := collection.Find(e.ctx, bson.M{}, options.Find().SetLimit(100))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(e.ctx)

	// 收集所有字段
	fieldMap := make(map[string]string)
	for cursor.Next(e.ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}

		for key, value := range doc {
			if _, exists := fieldMap[key]; !exists {
				fieldMap[key] = inferBSONType(value)
			}
		}
	}

	var columns []ColumnInfo
	for name, dataType := range fieldMap {
		columns = append(columns, ColumnInfo{
			Name:       name,
			Type:       dataType,
			IsNullable: true, // MongoDB 字段都是可选的
			IsPrimaryKey: name == "_id",
		})
	}

	return columns, nil
}

// inferBSONType 推断 BSON 类型
func inferBSONType(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	case int, int32, int64:
		return "int"
	case float32, float64:
		return "double"
	case bool:
		return "bool"
	case time.Time:
		return "date"
	case bson.M, map[string]interface{}:
		return "object"
	case bson.A, []interface{}:
		return "array"
	default:
		return "mixed"
	}
}

// GetIndexes 获取索引信息
func (e *MongoDBExtractor) GetIndexes(schema, table string) ([]IndexInfo, error) {
	if schema == "" {
		schema = e.dbName
	}

	db := e.client.Database(schema)
	collection := db.Collection(table)

	cursor, err := collection.Indexes().List(e.ctx)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(e.ctx)

	var indexes []IndexInfo
	for cursor.Next(e.ctx) {
		var idx bson.M
		if err := cursor.Decode(&idx); err != nil {
			continue
		}

		indexName := idx["name"].(string)
		key := idx["key"].(bson.M)

		var columns []string
		for col := range key {
			columns = append(columns, col)
		}

		isUnique := false
		if unique, ok := idx["unique"].(bool); ok {
			isUnique = unique
		}

		indexes = append(indexes, IndexInfo{
			Name:      indexName,
			Columns:   columns,
			IsUnique:  isUnique,
			IsPrimary: indexName == "_id_",
			Type:      "BTREE",
		})
	}

	return indexes, nil
}

// PreviewData 预览数据
func (e *MongoDBExtractor) PreviewData(schema, table string, limit int) ([]map[string]interface{}, error) {
	if schema == "" {
		schema = e.dbName
	}

	db := e.client.Database(schema)
	collection := db.Collection(table)

	cursor, err := collection.Find(e.ctx, bson.M{}, options.Find().SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(e.ctx)

	var result []map[string]interface{}
	for cursor.Next(e.ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}

		// 转换为 map[string]interface{}
		entry := make(map[string]interface{})
		for k, v := range doc {
			entry[k] = fmt.Sprintf("%v", v)
		}
		result = append(result, entry)
	}

	return result, nil
}

// Close 关闭连接
func (e *MongoDBExtractor) Close() error {
	return e.client.Disconnect(e.ctx)
}
