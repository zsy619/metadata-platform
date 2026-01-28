package adapter

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
)

// SQLServerExtractor SQL Server元数据提取器
type SQLServerExtractor struct {
	db *sql.DB
}

// NewSQLServerExtractor 创建SQL Server元数据提取器
func NewSQLServerExtractor(dsn string) (*SQLServerExtractor, error) {
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}
	return &SQLServerExtractor{db: db}, nil
}

// TestConnection 测试连接
func (e *SQLServerExtractor) TestConnection() error {
	return e.db.Ping()
}

func (e *SQLServerExtractor) GetSchemas() ([]string, error) {
	rows, err := e.db.Query("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []string
	systemSchemas := map[string]bool{
		"INFORMATION_SCHEMA": true,
		"sys":               true,
		"guest":             true,
		"db_owner":          true,
		"db_accessadmin":    true,
		"db_securityadmin":  true,
		"db_ddladmin":       true,
		"db_backupoperator": true,
		"db_datareader":     true,
		"db_datawriter":     true,
		"db_denydatareader": true,
		"db_denydatawriter": true,
	}
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, err
		}
		if !systemSchemas[schema] {
			schemas = append(schemas, schema)
		}
	}
	return schemas, nil
}

// GetTables 获取表列表
func (e *SQLServerExtractor) GetTables(schema string) ([]TableInfo, error) {
	if schema == "" {
		schema = "dbo"
	}
	
	query := `
		SELECT 
			t.TABLE_NAME,
			ISNULL(ep.value, '') as TABLE_COMMENT,
			t.TABLE_TYPE
		FROM INFORMATION_SCHEMA.TABLES t
		LEFT JOIN sys.extended_properties ep 
			ON ep.major_id = OBJECT_ID(t.TABLE_SCHEMA + '.' + t.TABLE_NAME)
			AND ep.minor_id = 0
			AND ep.name = 'MS_Description'
		WHERE t.TABLE_SCHEMA = @p1 AND t.TABLE_TYPE = 'BASE TABLE'
		ORDER BY t.TABLE_NAME
	`
	
	rows, err := e.db.Query(query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var t TableInfo
		var tableType string
		if err := rows.Scan(&t.Name, &t.Comment, &tableType); err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, nil
}

// GetViews 获取视图列表
func (e *SQLServerExtractor) GetViews(schema string) ([]ViewInfo, error) {
	if schema == "" {
		schema = "dbo"
	}
	
	query := `
		SELECT 
			TABLE_NAME,
			VIEW_DEFINITION
		FROM INFORMATION_SCHEMA.VIEWS
		WHERE TABLE_SCHEMA = @p1
		ORDER BY TABLE_NAME
	`
	
	rows, err := e.db.Query(query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []ViewInfo
	for rows.Next() {
		var v ViewInfo
		var definition sql.NullString
		if err := rows.Scan(&v.Name, &definition); err != nil {
			return nil, err
		}
		v.Definition = definition.String
		views = append(views, v)
	}
	return views, nil
}

// GetColumns 获取表字段信息
func (e *SQLServerExtractor) GetColumns(schema, table string) ([]ColumnInfo, error) {
	if schema == "" {
		schema = "dbo"
	}
	
	query := `
		SELECT 
			c.COLUMN_NAME,
			c.DATA_TYPE,
			ISNULL(c.CHARACTER_MAXIMUM_LENGTH, 0) as LENGTH,
			CASE WHEN c.IS_NULLABLE = 'YES' THEN 1 ELSE 0 END as IS_NULLABLE,
			c.COLUMN_DEFAULT,
			ISNULL(ep.value, '') as COLUMN_COMMENT,
			CASE WHEN pk.COLUMN_NAME IS NOT NULL THEN 1 ELSE 0 END as IS_PRIMARY_KEY,
			CASE WHEN c.COLUMN_DEFAULT LIKE '%IDENTITY%' THEN 1 ELSE 0 END as IS_AUTO_INCREMENT
		FROM INFORMATION_SCHEMA.COLUMNS c
		LEFT JOIN sys.extended_properties ep 
			ON ep.major_id = OBJECT_ID(c.TABLE_SCHEMA + '.' + c.TABLE_NAME)
			AND ep.minor_id = c.ORDINAL_POSITION
			AND ep.name = 'MS_Description'
		LEFT JOIN (
			SELECT ku.TABLE_SCHEMA, ku.TABLE_NAME, ku.COLUMN_NAME
			FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE ku
				ON tc.CONSTRAINT_NAME = ku.CONSTRAINT_NAME
			WHERE tc.CONSTRAINT_TYPE = 'PRIMARY KEY'
		) pk ON c.TABLE_SCHEMA = pk.TABLE_SCHEMA 
			AND c.TABLE_NAME = pk.TABLE_NAME 
			AND c.COLUMN_NAME = pk.COLUMN_NAME
		WHERE c.TABLE_SCHEMA = @p1 AND c.TABLE_NAME = @p2
		ORDER BY c.ORDINAL_POSITION
	`
	
	rows, err := e.db.Query(query, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var c ColumnInfo
		var defaultValue sql.NullString
		var isPrimaryKey, isAutoIncrement int
		
		if err := rows.Scan(
			&c.Name, &c.Type, &c.Length,
			&c.IsNullable, &defaultValue, &c.Comment,
			&isPrimaryKey, &isAutoIncrement,
		); err != nil {
			return nil, err
		}
		
		if defaultValue.Valid {
			c.DefaultValue = defaultValue.String
		}
		c.IsPrimaryKey = isPrimaryKey == 1
		c.IsAutoIncrement = isAutoIncrement == 1
		
		columns = append(columns, c)
	}
	return columns, nil
}

// GetIndexes 获取表索引信息
func (e *SQLServerExtractor) GetIndexes(schema, table string) ([]IndexInfo, error) {
	if schema == "" {
		schema = "dbo"
	}
	
	query := `
		SELECT 
			i.name as INDEX_NAME,
			COL_NAME(ic.object_id, ic.column_id) as COLUMN_NAME,
			i.is_unique,
			i.is_primary_key,
			i.type_desc
		FROM sys.indexes i
		JOIN sys.index_columns ic ON i.object_id = ic.object_id AND i.index_id = ic.index_id
		WHERE i.object_id = OBJECT_ID(@p1 + '.' + @p2)
		ORDER BY i.name, ic.key_ordinal
	`
	
	rows, err := e.db.Query(query, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indexMap := make(map[string]*IndexInfo)
	for rows.Next() {
		var indexName, columnName, indexType string
		var isUnique, isPrimary bool
		
		if err := rows.Scan(&indexName, &columnName, &isUnique, &isPrimary, &indexType); err != nil {
			return nil, err
		}
		
		if idx, exists := indexMap[indexName]; exists {
			idx.Columns = append(idx.Columns, columnName)
		} else {
			indexMap[indexName] = &IndexInfo{
				Name:      indexName,
				Columns:   []string{columnName},
				IsUnique:  isUnique,
				IsPrimary: isPrimary,
				Type:      indexType,
			}
		}
	}
	
	var indexes []IndexInfo
	for _, idx := range indexMap {
		indexes = append(indexes, *idx)
	}
	return indexes, nil
}

// PreviewData 预览数据
func (e *SQLServerExtractor) PreviewData(schema, table string, limit int) ([]map[string]interface{}, error) {
	if schema == "" {
		schema = "dbo"
	}
	
	query := fmt.Sprintf("SELECT TOP %d * FROM [%s].[%s]", limit, schema, table)
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
func (e *SQLServerExtractor) GetQueryColumns(query string, params []interface{}) ([]ColumnInfo, error) {
	return nil, fmt.Errorf("method GetQueryColumns not implemented for this adapter")
}

// Close 关闭连接
func (e *SQLServerExtractor) Close() error {
	return e.db.Close()
}
