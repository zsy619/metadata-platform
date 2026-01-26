package adapter

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLExtractor MySQL元数据提取器
type MySQLExtractor struct {
	db *sql.DB
}

// NewMySQLExtractor 创建MySQL元数据提取器
func NewMySQLExtractor(dsn string) (*MySQLExtractor, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &MySQLExtractor{db: db}, nil
}

// TestConnection 测试连接
func (e *MySQLExtractor) TestConnection() error {
	return e.db.Ping()
}

// GetTables 获取表列表
func (e *MySQLExtractor) GetTables(schema string) ([]TableInfo, error) {
	query := `
		SELECT 
			TABLE_NAME, TABLE_COMMENT, CREATE_TIME, UPDATE_TIME, ENGINE, TABLE_COLLATION
		FROM information_schema.TABLES 
		WHERE TABLE_SCHEMA = ? AND TABLE_TYPE = 'BASE TABLE'
	`
	rows, err := e.db.Query(query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var t TableInfo
		var createTime, updateTime []uint8 // MySQL driver might return []byte for time
		var comment, engine, collation sql.NullString
		
		if err := rows.Scan(&t.Name, &comment, &createTime, &updateTime, &engine, &collation); err != nil {
			return nil, err
		}
		t.Comment = comment.String
		t.Engine = engine.String
		t.Collation = collation.String
		
		// Parse time if needed, simplified here
		
		tables = append(tables, t)
	}
	return tables, nil
}

// GetViews 获取视图列表
func (e *MySQLExtractor) GetViews(schema string) ([]ViewInfo, error) {
	query := `
		SELECT TABLE_NAME
		FROM information_schema.VIEWS
		WHERE TABLE_SCHEMA = ?
	`
	rows, err := e.db.Query(query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []ViewInfo
	for rows.Next() {
		var v ViewInfo
		if err := rows.Scan(&v.Name); err != nil {
			return nil, err
		}
		views = append(views, v)
	}
	return views, nil
}

// GetColumns 获取表字段信息
func (e *MySQLExtractor) GetColumns(schema, table string) ([]ColumnInfo, error) {
	query := `
		SELECT 
			COLUMN_NAME, DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, 
			IS_NULLABLE, COLUMN_DEFAULT, COLUMN_COMMENT, COLUMN_KEY, EXTRA
		FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION
	`
	rows, err := e.db.Query(query, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var c ColumnInfo
		var length sql.NullInt64
		var isNullable, key, extra string
		var defaultValue sql.NullString
		
		if err := rows.Scan(
			&c.Name, &c.Type, &length, 
			&isNullable, &defaultValue, &c.Comment, &key, &extra,
		); err != nil {
			return nil, err
		}
		
		if length.Valid {
			c.Length = int(length.Int64)
		}
		c.IsNullable = isNullable == "YES"
		if defaultValue.Valid {
			c.DefaultValue = defaultValue.String
		}
		c.IsPrimaryKey = key == "PRI"
		c.IsAutoIncrement = extra == "auto_increment"
		
		columns = append(columns, c)
	}
	return columns, nil
}

// GetIndexes 获取表索引信息
func (e *MySQLExtractor) GetIndexes(schema, table string) ([]IndexInfo, error) {
	// 简化实现，实际可能需要查询 information_schema.STATISTICS
	return []IndexInfo{}, nil
}

// PreviewData 预览数据
func (e *MySQLExtractor) PreviewData(schema, table string, limit int) ([]map[string]interface{}, error) {
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

// Close 关闭连接
func (e *MySQLExtractor) Close() error {
	return e.db.Close()
}
