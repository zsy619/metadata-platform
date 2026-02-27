package adapter

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteExtractor SQLite元数据提取器
type SQLiteExtractor struct {
	db *sql.DB
}

// NewSQLiteExtractor 创建SQLite元数据提取器
func NewSQLiteExtractor(dsn string) (*SQLiteExtractor, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	return &SQLiteExtractor{db: db}, nil
}

// TestConnection 测试连接
func (e *SQLiteExtractor) TestConnection() error {
	return e.db.Ping()
}

func (e *SQLiteExtractor) GetSchemas() ([]string, error) {
	return []string{"main"}, nil
}

// GetTables 获取表列表
func (e *SQLiteExtractor) GetTables(schema string) ([]TableInfo, error) {
	// SQLite 不支持 schema，忽略该参数
	query := `
		SELECT name, '' as comment
		FROM sqlite_master 
		WHERE type = 'table' AND name NOT LIKE 'sqlite_%'
		ORDER BY name
	`
	
	rows, err := e.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var t TableInfo
		if err := rows.Scan(&t.Name, &t.Comment); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, nil
}

// GetViews 获取视图列表
func (e *SQLiteExtractor) GetViews(schema string) ([]ViewInfo, error) {
	query := `
		SELECT name, sql as definition
		FROM sqlite_master 
		WHERE type = 'view'
		ORDER BY name
	`
	
	rows, err := e.db.Query(query)
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
func (e *SQLiteExtractor) GetColumns(schema, table string) ([]ColumnInfo, error) {
	query := fmt.Sprintf("PRAGMA table_info(%s)", table)
	
	rows, err := e.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var c ColumnInfo
		var cid int
		var notNull int
		var pk int
		var defaultValue sql.NullString
		
		// PRAGMA table_info 返回: cid, name, type, notnull, dflt_value, pk
		if err := rows.Scan(&cid, &c.Name, &c.Type, &notNull, &defaultValue, &pk); err != nil {
			return nil, err
		}
		
		c.IsNullable = notNull == 0
		if defaultValue.Valid {
			c.DefaultValue = defaultValue.String
		}
		c.IsPrimaryKey = pk > 0
		
		// 检查是否为自增字段
		c.IsAutoIncrement = c.IsPrimaryKey && strings.ToUpper(c.Type) == "INTEGER"
		
		columns = append(columns, c)
	}
	return columns, nil
}

// GetIndexes 获取表索引信息
func (e *SQLiteExtractor) GetIndexes(schema, table string) ([]IndexInfo, error) {
	// 获取索引列表
	query := fmt.Sprintf("PRAGMA index_list(%s)", table)
	
	rows, err := e.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []IndexInfo
	for rows.Next() {
		var seq int
		var indexName string
		var unique int
		var origin string
		var partial int
		
		// PRAGMA index_list 返回: seq, name, unique, origin, partial
		if err := rows.Scan(&seq, &indexName, &unique, &origin, &partial); err != nil {
			return nil, err
		}
		
		// 获取索引的列信息
		colQuery := fmt.Sprintf("PRAGMA index_info(%s)", indexName)
		colRows, err := e.db.Query(colQuery)
		if err != nil {
			continue
		}
		
		var columns []string
		for colRows.Next() {
			var seqno, cid int
			var name sql.NullString
			
			if err := colRows.Scan(&seqno, &cid, &name); err != nil {
				continue
			}
			if name.Valid {
				columns = append(columns, name.String)
			}
		}
		colRows.Close()
		
		indexes = append(indexes, IndexInfo{
			Name:      indexName,
			Columns:   columns,
			IsUnique:  unique == 1,
			IsPrimary: origin == "pk",
			Type:      "BTREE",
		})
	}
	
	return indexes, nil
}

// PreviewData 预览数据
func (e *SQLiteExtractor) PreviewData(schema, table string, limit int) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s LIMIT %d", table, limit)
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
func (e *SQLiteExtractor) GetQueryColumns(query string, params []interface{}) ([]ColumnInfo, error) {
	return nil, fmt.Errorf("method GetQueryColumns not implemented for this adapter")
}

// GetProcedures 获取存储过程列表
func (e *SQLiteExtractor) GetProcedures(schema string) ([]ProcedureInfo, error) {
	return []ProcedureInfo{}, nil
}

// GetFunctions 获取函数列表
func (e *SQLiteExtractor) GetFunctions(schema string) ([]ProcedureInfo, error) {
	return []ProcedureInfo{}, nil
}

// Close 关闭连接
func (e *SQLiteExtractor) Close() error {
	return e.db.Close()
}
