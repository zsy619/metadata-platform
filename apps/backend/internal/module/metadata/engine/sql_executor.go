package engine

import (
	"database/sql"
	"fmt"
	"metadata-platform/internal/module/metadata/repository"
	"metadata-platform/internal/utils"
	"sync"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SQLExecutor SQL 执行器
type SQLExecutor struct {
	connRepo repository.MdConnRepository
	db       *gorm.DB
	conns    sync.Map // 缓存数据库连接: map[string]*gorm.DB (key: connID)
}

// NewSQLExecutor 创建一个新的 SQLExecutor 实例
func NewSQLExecutor(db *gorm.DB, connRepo repository.MdConnRepository) *SQLExecutor {
	return &SQLExecutor{
		db:       db,
		connRepo: connRepo,
	}
}

// Execute 执行 SQL 查询并返回结果
func (e *SQLExecutor) Execute(connID string, sqlStr string, args ...any) ([]map[string]any, error) {
	db, err := e.GetConnection(connID)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	rows, err := db.Raw(sqlStr, args...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results, err := e.parseRows(rows)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start)
	utils.SugarLogger.Infof("SQL Executed [%v]: %s | Args: %v", duration, sqlStr, args)
	
	// 慢查询告警
	if duration > time.Second {
		utils.SugarLogger.Warnf("Slow query detected: %v", duration)
	}

	return results, nil
}

// ExecuteCount 执行 COUNT 查询并返回总数
func (e *SQLExecutor) ExecuteCount(connID string, sqlStr string, args ...any) (int64, error) {
	db, err := e.GetConnection(connID)
	if err != nil {
		return 0, err
	}

	// 简单的包装成 COUNT(*)
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS total", sqlStr)
	var count int64
	if err := db.Raw(countSQL, args...).Scan(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetConnection 获取或创建目标数据库连接
func (e *SQLExecutor) GetConnection(connID string) (*gorm.DB, error) {
	if val, ok := e.conns.Load(connID); ok {
		db := val.(*gorm.DB)
		// 健康检查
		sqlDB, err := db.DB()
		if err == nil && sqlDB.Ping() == nil {
			return db, nil
		}
	}

	// 从元数据库加载连接信息
	connInfo, err := e.connRepo.GetMDConnByID(connID)
	if err != nil {
		return nil, err
	}
	if connInfo == nil {
		return nil, fmt.Errorf("connection %s not found", connID)
	}

	// 构建 DSN (目前支持 MySQL and SQLite)
	var dialector gorm.Dialector
	switch connInfo.ConnKind {
	case "MySQL":
		dsn := connInfo.ConnConn
		if dsn == "" {
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				connInfo.ConnUser,
				connInfo.ConnPassword,
				connInfo.ConnHost,
				connInfo.ConnPort,
				connInfo.ConnDatabase,
			)
		}
		dialector = mysql.Open(dsn)
	case "SQLite":
		dialector = sqlite.Open(connInfo.ConnConn)
	default:
		return nil, fmt.Errorf("unsupported database kind: %s", connInfo.ConnKind)
	}

	// 打开新连接
	newDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to target database: %w", err)
	}

	e.conns.Store(connID, newDB)
	return newDB, nil
}

// SetCustomConnection 允许手动注入连接（主要用于单元测试）
func (e *SQLExecutor) SetCustomConnection(connID string, db *gorm.DB) {
	e.conns.Store(connID, db)
}

// parseRows 将 sql.Rows 解析为 map 切片
func (e *SQLExecutor) parseRows(rows *sql.Rows) ([]map[string]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := make([]map[string]any, 0)
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		entry := make(map[string]any)
		for i, col := range columns {
			val := values[i]

			// 类型转换处理 (主要是处理 []byte)
			b, ok := val.([]byte)
			if ok {
				entry[col] = string(b)
			} else {
				entry[col] = val
			}
		}
		results = append(results, entry)
	}
	return results, nil
}
