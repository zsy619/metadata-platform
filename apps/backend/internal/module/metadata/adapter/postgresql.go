package adapter

import (
	"database/sql"
	"fmt"

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
		SELECT tablename, '' as comment
		FROM pg_catalog.pg_tables
		WHERE schemaname = $1
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
		SELECT column_name, udt_name as data_type, 
		       COALESCE(character_maximum_length, 0) as length,
		       is_nullable = 'YES' as is_nullable
		FROM information_schema.columns 
		WHERE table_schema = $1 AND table_name = $2
		ORDER BY ordinal_position
	`
	rows, err := e.db.Query(query, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		if err := rows.Scan(&col.Name, &col.Type, &col.Length, &col.IsNullable); err != nil {
			return nil, err
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

func (e *postgreSQLExtractor) Close() error {
	return e.db.Close()
}
