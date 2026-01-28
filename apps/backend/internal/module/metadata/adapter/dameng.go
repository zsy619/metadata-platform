package adapter

import (
	"database/sql"
	"fmt"

	_ "gitee.com/chunanyong/dm"
)

// DamengExtractor 达梦数据库元数据提取器
type DamengExtractor struct {
	db *sql.DB
}

// NewDamengExtractor 创建达梦数据库元数据提取器
func NewDamengExtractor(dsn string) (*DamengExtractor, error) {
	db, err := sql.Open("dm", dsn)
	if err != nil {
		return nil, err
	}
	return &DamengExtractor{db: db}, nil
}

// TestConnection 测试连接
func (e *DamengExtractor) TestConnection() error {
	return e.db.Ping()
}

func (e *DamengExtractor) GetSchemas() ([]string, error) {
	rows, err := e.db.Query("SELECT DISTINCT OWNER FROM ALL_TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []string
	systemSchemas := map[string]bool{
		"SYS": true, "SYSDBA": true, "SYSSSO": true, "SYSAUDITOR": true,
		"CTISYS": true, "SYSJOB": true, "DB_AUDIT": true, "DAMENG": true,
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
func (e *DamengExtractor) GetTables(schema string) ([]TableInfo, error) {
	// 达梦数据库的系统表结构类似 Oracle
	query := `
		SELECT 
			TABLE_NAME,
			NVL(COMMENTS, '') as TABLE_COMMENT
		FROM USER_TABLES t
		LEFT JOIN USER_TAB_COMMENTS c 
			ON t.TABLE_NAME = c.TABLE_NAME
		WHERE t.TABLESPACE_NAME IS NOT NULL
		ORDER BY TABLE_NAME
	`
	
	// 如果指定了 schema，使用 ALL_TABLES
	if schema != "" {
		query = `
			SELECT 
				TABLE_NAME,
				NVL(COMMENTS, '') as TABLE_COMMENT
			FROM ALL_TABLES t
			LEFT JOIN ALL_TAB_COMMENTS c 
				ON t.TABLE_NAME = c.TABLE_NAME AND t.OWNER = c.OWNER
			WHERE t.OWNER = :1
			ORDER BY TABLE_NAME
		`
		rows, err := e.db.Query(query, schema)
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
func (e *DamengExtractor) GetViews(schema string) ([]ViewInfo, error) {
	query := `
		SELECT 
			VIEW_NAME,
			TEXT as VIEW_DEFINITION
		FROM USER_VIEWS
		ORDER BY VIEW_NAME
	`
	
	if schema != "" {
		query = `
			SELECT 
				VIEW_NAME,
				TEXT as VIEW_DEFINITION
			FROM ALL_VIEWS
			WHERE OWNER = :1
			ORDER BY VIEW_NAME
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
func (e *DamengExtractor) GetColumns(schema, table string) ([]ColumnInfo, error) {
	query := `
		SELECT 
			c.COLUMN_NAME,
			c.DATA_TYPE,
			NVL(c.DATA_LENGTH, 0) as DATA_LENGTH,
			CASE WHEN c.NULLABLE = 'Y' THEN 1 ELSE 0 END as IS_NULLABLE,
			c.DATA_DEFAULT,
			NVL(cc.COMMENTS, '') as COLUMN_COMMENT,
			CASE WHEN pk.COLUMN_NAME IS NOT NULL THEN 1 ELSE 0 END as IS_PRIMARY_KEY,
			CASE WHEN c.DATA_DEFAULT LIKE '%IDENTITY%' THEN 1 ELSE 0 END as IS_AUTO_INCREMENT
		FROM USER_TAB_COLUMNS c
		LEFT JOIN USER_COL_COMMENTS cc 
			ON c.TABLE_NAME = cc.TABLE_NAME AND c.COLUMN_NAME = cc.COLUMN_NAME
		LEFT JOIN (
			SELECT acc.COLUMN_NAME, acc.TABLE_NAME
			FROM USER_CONSTRAINTS ac
			JOIN USER_CONS_COLUMNS acc 
				ON ac.CONSTRAINT_NAME = acc.CONSTRAINT_NAME
			WHERE ac.CONSTRAINT_TYPE = 'P'
		) pk ON c.TABLE_NAME = pk.TABLE_NAME AND c.COLUMN_NAME = pk.COLUMN_NAME
		WHERE c.TABLE_NAME = :1
		ORDER BY c.COLUMN_ID
	`
	
	if schema != "" {
		query = `
			SELECT 
				c.COLUMN_NAME,
				c.DATA_TYPE,
				NVL(c.DATA_LENGTH, 0) as DATA_LENGTH,
				CASE WHEN c.NULLABLE = 'Y' THEN 1 ELSE 0 END as IS_NULLABLE,
				c.DATA_DEFAULT,
				NVL(cc.COMMENTS, '') as COLUMN_COMMENT,
				CASE WHEN pk.COLUMN_NAME IS NOT NULL THEN 1 ELSE 0 END as IS_PRIMARY_KEY,
				CASE WHEN c.DATA_DEFAULT LIKE '%IDENTITY%' THEN 1 ELSE 0 END as IS_AUTO_INCREMENT
			FROM ALL_TAB_COLUMNS c
			LEFT JOIN ALL_COL_COMMENTS cc 
				ON c.TABLE_NAME = cc.TABLE_NAME 
				AND c.COLUMN_NAME = cc.COLUMN_NAME 
				AND c.OWNER = cc.OWNER
			LEFT JOIN (
				SELECT acc.COLUMN_NAME, acc.TABLE_NAME, acc.OWNER
				FROM ALL_CONSTRAINTS ac
				JOIN ALL_CONS_COLUMNS acc 
					ON ac.CONSTRAINT_NAME = acc.CONSTRAINT_NAME 
					AND ac.OWNER = acc.OWNER
				WHERE ac.CONSTRAINT_TYPE = 'P'
			) pk ON c.TABLE_NAME = pk.TABLE_NAME 
				AND c.COLUMN_NAME = pk.COLUMN_NAME 
				AND c.OWNER = pk.OWNER
			WHERE c.OWNER = :1 AND c.TABLE_NAME = :2
			ORDER BY c.COLUMN_ID
		`
		rows, err := e.db.Query(query, schema, table)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		return e.scanColumns(rows)
	}
	
	rows, err := e.db.Query(query, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return e.scanColumns(rows)
}

func (e *DamengExtractor) scanColumns(rows *sql.Rows) ([]ColumnInfo, error) {
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
func (e *DamengExtractor) GetIndexes(schema, table string) ([]IndexInfo, error) {
	query := `
		SELECT 
			i.INDEX_NAME,
			ic.COLUMN_NAME,
			CASE WHEN i.UNIQUENESS = 'UNIQUE' THEN 1 ELSE 0 END as IS_UNIQUE,
			CASE WHEN c.CONSTRAINT_TYPE = 'P' THEN 1 ELSE 0 END as IS_PRIMARY,
			i.INDEX_TYPE
		FROM USER_INDEXES i
		JOIN USER_IND_COLUMNS ic 
			ON i.INDEX_NAME = ic.INDEX_NAME
		LEFT JOIN USER_CONSTRAINTS c 
			ON i.INDEX_NAME = c.INDEX_NAME
		WHERE i.TABLE_NAME = :1
		ORDER BY i.INDEX_NAME, ic.COLUMN_POSITION
	`
	
	if schema != "" {
		query = `
			SELECT 
				i.INDEX_NAME,
				ic.COLUMN_NAME,
				CASE WHEN i.UNIQUENESS = 'UNIQUE' THEN 1 ELSE 0 END as IS_UNIQUE,
				CASE WHEN c.CONSTRAINT_TYPE = 'P' THEN 1 ELSE 0 END as IS_PRIMARY,
				i.INDEX_TYPE
			FROM ALL_INDEXES i
			JOIN ALL_IND_COLUMNS ic 
				ON i.INDEX_NAME = ic.INDEX_NAME AND i.OWNER = ic.INDEX_OWNER
			LEFT JOIN ALL_CONSTRAINTS c 
				ON i.INDEX_NAME = c.INDEX_NAME AND i.OWNER = c.OWNER
			WHERE i.TABLE_OWNER = :1 AND i.TABLE_NAME = :2
			ORDER BY i.INDEX_NAME, ic.COLUMN_POSITION
		`
		rows, err := e.db.Query(query, schema, table)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		return e.scanIndexes(rows)
	}
	
	rows, err := e.db.Query(query, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return e.scanIndexes(rows)
}

func (e *DamengExtractor) scanIndexes(rows *sql.Rows) ([]IndexInfo, error) {
	indexMap := make(map[string]*IndexInfo)
	for rows.Next() {
		var indexName, columnName, indexType string
		var isUnique, isPrimary int
		
		if err := rows.Scan(&indexName, &columnName, &isUnique, &isPrimary, &indexType); err != nil {
			return nil, err
		}
		
		if idx, exists := indexMap[indexName]; exists {
			idx.Columns = append(idx.Columns, columnName)
		} else {
			indexMap[indexName] = &IndexInfo{
				Name:      indexName,
				Columns:   []string{columnName},
				IsUnique:  isUnique == 1,
				IsPrimary: isPrimary == 1,
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
func (e *DamengExtractor) PreviewData(schema, table string, limit int) ([]map[string]interface{}, error) {
	var query string
	if schema != "" {
		query = fmt.Sprintf("SELECT * FROM %s.%s WHERE ROWNUM <= %d", schema, table, limit)
	} else {
		query = fmt.Sprintf("SELECT * FROM %s WHERE ROWNUM <= %d", table, limit)
	}
	
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
func (e *DamengExtractor) Close() error {
	return e.db.Close()
}
