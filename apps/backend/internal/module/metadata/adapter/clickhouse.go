package adapter

import (
	"database/sql"
	"fmt"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

// ClickHouseExtractor ClickHouse元数据提取器
type ClickHouseExtractor struct {
	db *sql.DB
}

// NewClickHouseExtractor 创建ClickHouse元数据提取器
func NewClickHouseExtractor(dsn string) (*ClickHouseExtractor, error) {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	return &ClickHouseExtractor{db: db}, nil
}

// TestConnection 测试连接
func (e *ClickHouseExtractor) TestConnection() error {
	return e.db.Ping()
}

func (e *ClickHouseExtractor) GetSchemas() ([]string, error) {
	var currentDB string
	// ClickHouse 中 Schema 等同于 Database
	// 按照需求，只显示当前连接的数据库
	err := e.db.QueryRow("SELECT database()").Scan(&currentDB)
	if err != nil {
		return nil, err
	}
	return []string{currentDB}, nil
}

// GetTables 获取表列表
func (e *ClickHouseExtractor) GetTables(schema string) ([]TableInfo, error) {
	query := `
		SELECT 
			name,
			comment,
			engine,
			create_table_query
		FROM system.tables
		WHERE database = ?
		ORDER BY name
	`
	
	rows, err := e.db.Query(query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var t TableInfo
		var createQuery string
		if err := rows.Scan(&t.Name, &t.Comment, &t.Engine, &createQuery); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, nil
}

// GetViews 获取视图列表
func (e *ClickHouseExtractor) GetViews(schema string) ([]ViewInfo, error) {
	query := `
		SELECT 
			name,
			as_select as definition
		FROM system.tables
		WHERE database = ? AND engine = 'View'
		ORDER BY name
	`
	
	rows, err := e.db.Query(query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []ViewInfo
	for rows.Next() {
		var v ViewInfo
		if err := rows.Scan(&v.Name, &v.Definition); err != nil {
			return nil, err
		}
		views = append(views, v)
	}
	return views, nil
}

// GetColumns 获取表字段信息
func (e *ClickHouseExtractor) GetColumns(schema, table string) ([]ColumnInfo, error) {
	query := `
		SELECT 
			name,
			type,
			0 as length,
			CASE WHEN type LIKE '%Nullable%' THEN 1 ELSE 0 END as is_nullable,
			default_expression,
			comment,
			is_in_primary_key
		FROM system.columns
		WHERE database = ? AND table = ?
		ORDER BY position
	`
	
	rows, err := e.db.Query(query, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var c ColumnInfo
		var defaultExpr string
		var isInPrimaryKey uint8
		
		if err := rows.Scan(
			&c.Name, &c.Type, &c.Length,
			&c.IsNullable, &defaultExpr, &c.Comment,
			&isInPrimaryKey,
		); err != nil {
			return nil, err
		}
		
		if defaultExpr != "" {
			c.DefaultValue = defaultExpr
		}
		c.IsPrimaryKey = isInPrimaryKey == 1
		c.IsAutoIncrement = false // ClickHouse 不支持自增
		
		columns = append(columns, c)
	}
	return columns, nil
}

// GetIndexes 获取表索引信息
func (e *ClickHouseExtractor) GetIndexes(schema, table string) ([]IndexInfo, error) {
	// ClickHouse 的索引信息较为特殊，这里简化实现
	// 主键信息已在 GetColumns 中返回
	return []IndexInfo{}, nil
}

// PreviewData 预览数据
func (e *ClickHouseExtractor) PreviewData(schema, table string, limit int) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM `%s`.`%s` LIMIT %d", schema, table, limit)
	rows, err := e.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	var result []map[string]interface{}
	
	for rows.Next() {
		rows.Scan(valuePtrs...)
		
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		result = append(result, entry)
	}
	
	return result, nil
}
// GetQueryColumns 获取查询结果的列信息
func (e *ClickHouseExtractor) GetQueryColumns(query string, params []interface{}) ([]ColumnInfo, error) {
	return nil, fmt.Errorf("method GetQueryColumns not implemented for this adapter")
}

// Close 关闭连接
func (e *ClickHouseExtractor) Close() error {
	return e.db.Close()
}
