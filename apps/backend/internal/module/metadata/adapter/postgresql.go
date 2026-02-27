package adapter

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type postgreSQLExtractor struct {
	db *sql.DB
}

// NewPostgreSQLExtractor 创建 PostgreSQL 提取器
func NewPostgreSQLExtractor(dsn string) (MetadataExtractor, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &postgreSQLExtractor{db: db}, nil
}

func (e *postgreSQLExtractor) TestConnection() error {
	return e.db.Ping()
}

func (e *postgreSQLExtractor) GetSchemas() ([]string, error) {
	rows, err := e.db.Query("SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('information_schema', 'pg_catalog', 'pg_toast')")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, err
		}
		schemas = append(schemas, schema)
	}
	return schemas, nil
}

func (e *postgreSQLExtractor) GetTables(schema string) ([]TableInfo, error) {
	query := `
		SELECT 
			c.relname AS tablename,
			COALESCE(obj_description(c.oid, 'pg_class'), '') AS comment
		FROM pg_catalog.pg_class c
		JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = $1
		AND c.relkind = 'r'
		ORDER BY c.relname
	`
	rows, err := e.db.Query(query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var table TableInfo
		if err := rows.Scan(&table.Name, &table.Comment); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func (e *postgreSQLExtractor) GetViews(schema string) ([]ViewInfo, error) {
	query := `
		SELECT viewname, definition
		FROM pg_catalog.pg_views
		WHERE schemaname = $1
	`
	rows, err := e.db.Query(query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []ViewInfo
	for rows.Next() {
		var view ViewInfo
		if err := rows.Scan(&view.Name, &view.Definition); err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, nil
}

func (e *postgreSQLExtractor) GetColumns(schema, table string) ([]ColumnInfo, error) {
	query := `
		SELECT 
			a.attname AS column_name,
			t.typname AS data_type,
			COALESCE(a.atttypmod - 4, 0) AS length,
			NOT a.attnotnull AS is_nullable,
			COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') AS column_default,
			COALESCE(col_description(a.attrelid, a.attnum), '') AS column_comment
		FROM pg_catalog.pg_attribute a
		JOIN pg_catalog.pg_class c ON c.oid = a.attrelid
		JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		JOIN pg_catalog.pg_type t ON t.oid = a.atttypid
		LEFT JOIN pg_catalog.pg_attrdef ad ON ad.adrelid = c.oid AND ad.adnum = a.attnum
		WHERE n.nspname = $1
		AND c.relname = $2
		AND a.attnum > 0
		AND NOT a.attisdropped
		ORDER BY a.attnum
	`
	rows, err := e.db.Query(query, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		var defaultValue sql.NullString
		if err := rows.Scan(&col.Name, &col.Type, &col.Length, &col.IsNullable, &defaultValue, &col.Comment); err != nil {
			return nil, err
		}
		if defaultValue.Valid {
			col.DefaultValue = defaultValue.String
		}
		columns = append(columns, col)
	}
	return columns, nil
}

func (e *postgreSQLExtractor) GetIndexes(schema, table string) ([]IndexInfo, error) {
	// 简化实现
	return []IndexInfo{}, nil
}

func (e *postgreSQLExtractor) PreviewData(schema, table string, limit int) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s LIMIT %d", schema, table, limit)
	rows, err := e.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	var result []map[string]interface{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columns[i]
			if b, ok := val.([]byte); ok {
				m[colName] = string(b)
			} else {
				m[colName] = val
			}
		}
		result = append(result, m)
	}
	return result, nil
}

// GetQueryColumns 获取查询结果的列信息
func (e *postgreSQLExtractor) GetQueryColumns(query string, params []interface{}) ([]ColumnInfo, error) {
	return nil, fmt.Errorf("method GetQueryColumns not implemented for this adapter")
}

// GetProcedures 获取存储过程列表
func (e *postgreSQLExtractor) GetProcedures(schema string) ([]ProcedureInfo, error) {
	query := `
		SELECT 
			p.proname AS routine_name,
			pg_get_functiondef(p.oid) AS definition,
			COALESCE(obj_description(p.oid, 'pg_proc'), '') AS comment,
			l.lanname AS language
		FROM pg_catalog.pg_proc p
		JOIN pg_catalog.pg_namespace n ON n.oid = p.pronamespace
		JOIN pg_catalog.pg_language l ON l.oid = p.prolang
		WHERE n.nspname = $1
		AND p.prokind = 'p'
		ORDER BY p.proname
	`
	rows, err := e.db.Query(query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var procedures []ProcedureInfo
	for rows.Next() {
		var p ProcedureInfo
		var definition, comment, language sql.NullString

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
func (e *postgreSQLExtractor) GetFunctions(schema string) ([]ProcedureInfo, error) {
	query := `
		SELECT 
			p.proname AS routine_name,
			pg_get_functiondef(p.oid) AS definition,
			COALESCE(obj_description(p.oid, 'pg_proc'), '') AS comment,
			l.lanname AS language,
			pg_get_function_result(p.oid) AS return_type
		FROM pg_catalog.pg_proc p
		JOIN pg_catalog.pg_namespace n ON n.oid = p.pronamespace
		JOIN pg_catalog.pg_language l ON l.oid = p.prolang
		WHERE n.nspname = $1
		AND p.prokind = 'f'
		ORDER BY p.proname
	`
	rows, err := e.db.Query(query, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var functions []ProcedureInfo
	for rows.Next() {
		var f ProcedureInfo
		var definition, comment, language, returnType sql.NullString

		if err := rows.Scan(&f.Name, &definition, &comment, &language, &returnType); err != nil {
			return nil, err
		}

		f.Type = "FUNCTION"
		f.Schema = schema
		f.Definition = definition.String
		f.Comment = comment.String
		f.Language = language.String
		f.ReturnType = returnType.String

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
func (e *postgreSQLExtractor) getProcedureParameters(schema, name string) (string, error) {
	query := `
		SELECT 
			p.proname,
			unnest(p.proargnames) AS param_name,
			unnest(p.proargtypes::regtype[]) AS param_type,
			unnest(p.proargmodes) AS param_mode
		FROM pg_catalog.pg_proc p
		JOIN pg_catalog.pg_namespace n ON n.oid = p.pronamespace
		WHERE n.nspname = $1 AND p.proname = $2
	`
	rows, err := e.db.Query(query, schema, name)
	if err != nil {
		// 如果参数查询失败,尝试简化的查询方式
		return e.getSimpleParameters(schema, name)
	}
	defer rows.Close()

	var params []string
	for rows.Next() {
		var procName, paramName, paramType, paramMode sql.NullString

		if err := rows.Scan(&procName, &paramName, &paramType, &paramMode); err != nil {
			continue
		}

		paramStr := ""
		if paramMode.Valid {
			modeMap := map[string]string{
				"i": "IN",
				"o": "OUT",
				"b": "INOUT",
				"v": "VARIADIC",
				"t": "TABLE",
			}
			if mode, ok := modeMap[paramMode.String]; ok {
				paramStr += mode + " "
			}
		}
		if paramName.Valid {
			paramStr += paramName.String + " "
		}
		if paramType.Valid {
			paramStr += paramType.String
		}
		if paramStr != "" {
			params = append(params, paramStr)
		}
	}

	if len(params) > 0 {
		return strings.Join(params, ", "), nil
	}

	// 如果没有获取到参数,尝试简化方式
	return e.getSimpleParameters(schema, name)
}

// getSimpleParameters 简化的参数获取方式
func (e *postgreSQLExtractor) getSimpleParameters(schema, name string) (string, error) {
	query := `
		SELECT 
			pg_get_function_arguments(p.oid) AS arguments
		FROM pg_catalog.pg_proc p
		JOIN pg_catalog.pg_namespace n ON n.oid = p.pronamespace
		WHERE n.nspname = $1 AND p.proname = $2
	`
	var arguments sql.NullString
	err := e.db.QueryRow(query, schema, name).Scan(&arguments)
	if err != nil {
		return "", err
	}
	return arguments.String, nil
}

// Close 关闭连接
func (e *postgreSQLExtractor) Close() error {
	return e.db.Close()
}
