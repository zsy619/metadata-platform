package adapter

import "time"

// ColumnInfo 列信息
type ColumnInfo struct {
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	Length        int         `json:"length"`
	IsNullable    bool        `json:"is_nullable"`
	DefaultValue  interface{} `json:"default_value"`
	Comment       string      `json:"comment"`
	IsPrimaryKey  bool        `json:"is_primary_key"`
	IsAutoIncrement bool      `json:"is_auto_increment"`
	Sort          int         `json:"sort"`
}

// TableInfo 表信息
type TableInfo struct {
	Name        string    `json:"name"`
	Comment     string    `json:"comment"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
	Engine      string    `json:"engine"`
	Collation   string    `json:"collation"`
}

// ViewInfo 视图信息
type ViewInfo struct {
	Name        string `json:"name"`
	Definition  string `json:"definition"`
	Comment     string `json:"comment"`
}

// IndexInfo 索引信息
type IndexInfo struct {
	Name        string   `json:"name"`
	Columns     []string `json:"columns"`
	IsUnique    bool     `json:"is_unique"`
	IsPrimary   bool     `json:"is_primary"`
	Type        string   `json:"type"`
}

// ProcedureInfo 存储过程/函数信息
type ProcedureInfo struct {
	Name        string `json:"name"`        // 名称
	Type        string `json:"type"`        // 类型: PROCEDURE 或 FUNCTION
	Definition  string `json:"definition"`  // 定义/代码
	Comment     string `json:"comment"`     // 注释
	ReturnType  string `json:"return_type"` // 返回类型 (仅函数有)
	Parameters  string `json:"parameters"`  // 参数列表
	Schema      string `json:"schema"`      // 模式
	Language    string `json:"language"`    // 语言
}

// MetadataExtractor 元数据提取接口
type MetadataExtractor interface {
	// TestConnection 测试连接
	TestConnection() error

	// GetSchemas 获取特定的命名空间/模式列表
	GetSchemas() ([]string, error)

	
	// GetTables 获取表列表
	GetTables(schema string) ([]TableInfo, error)
	
	// GetViews 获取视图列表
	GetViews(schema string) ([]ViewInfo, error)
	
	// GetColumns 获取表字段信息
	GetColumns(schema, table string) ([]ColumnInfo, error)
	
	// GetIndexes 获取表索引信息
	GetIndexes(schema, table string) ([]IndexInfo, error)
	
	// PreviewData 预览数据
	PreviewData(schema, table string, limit int) ([]map[string]interface{}, error)
	
	// GetQueryColumns 获取查询结果的列信息（不返回数据）
	GetQueryColumns(query string, params []interface{}) ([]ColumnInfo, error)

	// GetProcedures 获取存储过程列表
	GetProcedures(schema string) ([]ProcedureInfo, error)

	// GetFunctions 获取函数列表
	GetFunctions(schema string) ([]ProcedureInfo, error)

	// Close 关闭连接
	Close() error
}
