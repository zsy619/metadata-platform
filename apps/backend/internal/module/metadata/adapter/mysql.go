package adapter

import (
	"database/sql"
	"fmt"
	"strings"

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

func (e *MySQLExtractor) GetSchemas() ([]string, error) {
	var currentDB string
	// MySQL 中 Schema 等同于 Database
	// 按照需求，只显示当前连接的数据库
	err := e.db.QueryRow("SELECT DATABASE()").Scan(&currentDB)
	if err != nil {
		return nil, err
	}
	return []string{currentDB}, nil
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
		SELECT 
			t.TABLE_NAME,
			t.TABLE_COMMENT
		FROM information_schema.TABLES t
		WHERE t.TABLE_SCHEMA = ? AND t.TABLE_TYPE = 'VIEW'
	`
	rows, err := e.db.Query(query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []ViewInfo
	for rows.Next() {
		var v ViewInfo
		var comment sql.NullString
		if err := rows.Scan(&v.Name, &comment); err != nil {
			return nil, err
		}
		v.Comment = comment.String
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

// GetQueryColumns 获取查询结果的列信息
func (e *MySQLExtractor) GetQueryColumns(query string, params []interface{}) ([]ColumnInfo, error) {
	// 针对 SELECT 语句，拼接 LIMIT 0 防止查询大量数据
	// 注意：这是一个简单处理，可能不适用于所有复杂 SQL (如已有 LIMIT)
	// 但 Go sql.DB 提供的 ColumnTypes 可以即使为空结果集也能返回列信息
	// 最好还是包裹一层: SELECT * FROM (UserQuery) AS tmp LIMIT 0
	wrappedQuery := fmt.Sprintf("SELECT * FROM (%s) AS tmp_metadata_extractor LIMIT 0", query)

	rows, err := e.db.Query(wrappedQuery, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	var columns []ColumnInfo
	for _, ct := range columnTypes {
		c := ColumnInfo{
			Name:       ct.Name(),
			Type:       ct.DatabaseTypeName(),
			IsNullable: true, // 难以准确获取，默认为 true
		}

		if length, ok := ct.Length(); ok {
			c.Length = int(length)
		}
		if nullable, ok := ct.Nullable(); ok {
			c.IsNullable = nullable
		}

		columns = append(columns, c)
	}
	return columns, nil
}

// GetProcedures 获取存储过程列表
func (e *MySQLExtractor) GetProcedures(schema string) ([]ProcedureInfo, error) {
	query := `
		SELECT 
			ROUTINE_NAME,
			ROUTINE_DEFINITION,
			ROUTINE_COMMENT,
			'DET SQL' AS LANGUAGE
		FROM information_schema.ROUTINES
		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'PROCEDURE'
	`
	rows, err := e.db.Query(query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var procedures []ProcedureInfo
	for rows.Next() {
		var p ProcedureInfo
		var definition, comment sql.NullString
		var language sql.NullString

		if err := rows.Scan(&p.Name, &definition, &comment, &language); err != nil {
			return nil, err
		}

		p.Type = "PROCEDURE"
		p.Schema = schema
		p.Definition = definition.String
		p.Comment = comment.String
		p.Language = language.String

		// 获取存储过程的参数信息
		params, err := e.getProcedureParameters(schema, p.Name)
		if err == nil {
			p.Parameters = params
		}

		procedures = append(procedures, p)
	}
	return procedures, nil
}

// GetFunctions 获取函数列表
func (e *MySQLExtractor) GetFunctions(schema string) ([]ProcedureInfo, error) {
	query := `
		SELECT 
			ROUTINE_NAME,
			ROUTINE_DEFINITION,
			ROUTINE_COMMENT,
			DTD_IDENTIFIER AS RETURN_TYPE,
			'SQL' AS LANGUAGE
		FROM information_schema.ROUTINES
		WHERE ROUTINE_SCHEMA = ? AND ROUTINE_TYPE = 'FUNCTION'
	`
	rows, err := e.db.Query(query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var functions []ProcedureInfo
	for rows.Next() {
		var f ProcedureInfo
		var definition, comment, returnType, language sql.NullString

		if err := rows.Scan(&f.Name, &definition, &comment, &returnType, &language); err != nil {
			return nil, err
		}

		f.Type = "FUNCTION"
		f.Schema = schema
		f.Definition = definition.String
		f.Comment = comment.String
		f.ReturnType = returnType.String
		f.Language = language.String

		// 获取函数的参数信息
		params, err := e.getProcedureParameters(schema, f.Name)
		if err == nil {
			f.Parameters = params
		}

		functions = append(functions, f)
	}
	return functions, nil
}

// getProcedureParameters 获取存储过程/函数的参数列表
func (e *MySQLExtractor) getProcedureParameters(schema, name string) (string, error) {
	query := `
		SELECT 
			PARAMETER_NAME,
			DATA_TYPE,
			CHARACTER_MAXIMUM_LENGTH,
			NUMERIC_PRECISION,
			NUMERIC_SCALE,
			PARAMETER_MODE
		FROM information_schema.PARAMETERS
		WHERE SPECIFIC_SCHEMA = ? AND SPECIFIC_NAME = ?
		ORDER BY ORDINAL_POSITION
	`
	rows, err := e.db.Query(query, schema, name)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var params []string
	for rows.Next() {
		var paramName, dataType, paramMode sql.NullString
		var charLength, numPrecision, numScale sql.NullInt64

		if err := rows.Scan(&paramName, &dataType, &charLength, &numPrecision, &numScale, &paramMode); err != nil {
			continue
		}

		paramStr := ""
		if paramMode.Valid {
			paramStr += paramMode.String + " "
		}
		if paramName.Valid {
			paramStr += paramName.String + " "
		}
		if dataType.Valid {
			paramStr += dataType.String
			if charLength.Valid && charLength.Int64 > 0 {
				paramStr += fmt.Sprintf("(%d)", charLength.Int64)
			} else if numPrecision.Valid && numPrecision.Int64 > 0 {
				if numScale.Valid && numScale.Int64 > 0 {
					paramStr += fmt.Sprintf("(%d,%d)", numPrecision.Int64, numScale.Int64)
				} else {
					paramStr += fmt.Sprintf("(%d)", numPrecision.Int64)
				}
			}
		}
		if paramStr != "" {
			params = append(params, paramStr)
		}
	}

	return strings.Join(params, ", "), nil
}

// Close 关闭连接
func (e *MySQLExtractor) Close() error {
	return e.db.Close()
}
