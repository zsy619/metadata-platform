package utils

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"metadata-platform/configs"
)

// DBManager 数据库管理器
type DBManager struct {
	MetadataDB  *gorm.DB        // 元数据数据库连接
	UserDB      *gorm.DB        // 用户管理数据库连接
	Config      *configs.Config // 数据库配置
	HealthCheck chan bool       // 健康检查通道
	shutdown    chan bool       // 关闭通道
}

// NewDBManager 创建数据库管理器
func NewDBManager(cfg *configs.Config) (*DBManager, error) {
	// 初始化元数据数据库连接
	metadataDB, err := initDB(cfg.AppMode, cfg.MetadataDB)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metadata database: %w", err)
	}

	// 初始化用户管理数据库连接
	userDB, err := initDB(cfg.AppMode, cfg.UserDB)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user database: %w", err)
	}

	dbm := &DBManager{
		MetadataDB:  metadataDB,
		UserDB:      userDB,
		Config:      cfg,
		HealthCheck: make(chan bool, 1),
		shutdown:    make(chan bool, 1),
	}

	// 启动健康检查
	go dbm.startHealthCheck()

	return dbm, nil
}

// initDB 初始化数据库连接
func initDB(appMode string, dbCfg configs.DBConfig) (*gorm.DB, error) {
	// 构建DSN连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name,
	)

	// 配置日志级别
	logLevel := logger.Silent
	if appMode == "debug" {
		logLevel = logger.Info
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取底层SQL连接以设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	SugarLogger.Infof("Database connection established: %s:%d/%s", dbCfg.Host, dbCfg.Port, dbCfg.Name)

	return db, nil
}

// startHealthCheck 启动数据库健康检查
func (dbm *DBManager) startHealthCheck() {
	ticker := time.NewTicker(30 * time.Second) // 每30秒检查一次
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dbm.checkHealth()
		case <-dbm.shutdown:
			SugarLogger.Info("Database health check stopped")
			return
		}
	}
}

// checkHealth 检查数据库连接健康状态
func (dbm *DBManager) checkHealth() {
	// 检查元数据数据库
	metadataHealthy := dbm.checkSingleDBHealth(dbm.MetadataDB, "metadata")
	// 检查用户管理数据库
	userHealthy := dbm.checkSingleDBHealth(dbm.UserDB, "user")
	
	// 只有两个数据库都健康时才发送健康信号
	dbm.HealthCheck <- metadataHealthy && userHealthy
}

// checkSingleDBHealth 检查单个数据库连接健康状态
func (dbm *DBManager) checkSingleDBHealth(db *gorm.DB, dbType string) bool {
	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取底层SQL连接
	sqlDB, err := db.DB()
	if err != nil {
		SugarLogger.Errorf("Failed to get sql.DB for %s database: %v", dbType, err)
		return false
	}

	// 执行简单的查询来检查连接
	err = sqlDB.PingContext(ctx)
	if err != nil {
		SugarLogger.Errorf("%s database connection health check failed: %v", dbType, err)
		// 尝试重新连接
		if err := dbm.reconnectSingleDB(dbType); err != nil {
			SugarLogger.Errorf("Failed to reconnect to %s database: %v", dbType, err)
			return false
		} else {
			SugarLogger.Info("Successfully reconnected to %s database", dbType)
			return true
		}
	}
	SugarLogger.Debugf("%s database connection is healthy", dbType)
	return true
}

// reconnectSingleDB 重新连接单个数据库
func (dbm *DBManager) reconnectSingleDB(dbType string) error {
	var err error
	var newDB *gorm.DB
	
	// 根据数据库类型选择配置
	if dbType == "metadata" {
		// 关闭现有连接
		sqlDB, err := dbm.MetadataDB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				SugarLogger.Warnf("Failed to close old metadata database connection: %v", err)
			}
		}
		
		// 尝试重新连接，最多尝试3次
		for i := 0; i < 3; i++ {
			SugarLogger.Infof("Attempting to reconnect to metadata database (attempt %d/3)", i+1)
			newDB, err = initDB(dbm.Config.AppMode, dbm.Config.MetadataDB)
			if err == nil {
				break
			}
			// 等待一段时间后重试
			time.Sleep(time.Duration(i+1) * time.Second)
		}
		
		if err == nil {
			dbm.MetadataDB = newDB
		}
	} else if dbType == "user" {
		// 关闭现有连接
		sqlDB, err := dbm.UserDB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				SugarLogger.Warnf("Failed to close old user database connection: %v", err)
			}
		}
		
		// 尝试重新连接，最多尝试3次
		for i := 0; i < 3; i++ {
			SugarLogger.Infof("Attempting to reconnect to user database (attempt %d/3)", i+1)
			newDB, err = initDB(dbm.Config.AppMode, dbm.Config.UserDB)
			if err == nil {
				break
			}
			// 等待一段时间后重试
			time.Sleep(time.Duration(i+1) * time.Second)
		}
		
		if err == nil {
			dbm.UserDB = newDB
		}
	} else {
		return fmt.Errorf("unknown database type: %s", dbType)
	}
	
	return err
}

// Close 关闭数据库连接和健康检查
func (dbm *DBManager) Close() error {
	// 发送关闭信号
	dbm.shutdown <- true

	// 关闭健康检查通道
	close(dbm.HealthCheck)

	// 关闭元数据数据库连接
	sqlDB, err := dbm.MetadataDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get metadata sql.DB: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close metadata database connection: %w", err)
	}

	// 关闭用户管理数据库连接
	sqlDB, err = dbm.UserDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get user sql.DB: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close user database connection: %w", err)
	}

	SugarLogger.Info("All database connections closed")
	return nil
}

// GetMetadataDB 获取元数据数据库连接
func (dbm *DBManager) GetMetadataDB() *gorm.DB {
	return dbm.MetadataDB
}

// GetUserDB 获取用户管理数据库连接
func (dbm *DBManager) GetUserDB() *gorm.DB {
	return dbm.UserDB
}

// IsHealthy 检查数据库是否健康
func (dbm *DBManager) IsHealthy() bool {
	// 检查元数据数据库
	metadataHealthy := dbm.isSingleDBHealthy(dbm.MetadataDB)
	// 检查用户管理数据库
	userHealthy := dbm.isSingleDBHealthy(dbm.UserDB)
	
	return metadataHealthy && userHealthy
}

// isSingleDBHealthy 检查单个数据库是否健康
func (dbm *DBManager) isSingleDBHealthy(db *gorm.DB) bool {
	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 执行简单的查询来检查连接
	sqlDB, err := db.DB()
	if err != nil {
		return false
	}
	err = sqlDB.PingContext(ctx)
	return err == nil
}

// GetConnectionStats 获取数据库连接统计信息
func (dbm *DBManager) GetConnectionStats() map[string]any {
	// 获取元数据数据库统计信息
	metadataStats := dbm.getSingleDBStats(dbm.MetadataDB, "metadata")
	// 获取用户管理数据库统计信息
	userStats := dbm.getSingleDBStats(dbm.UserDB, "user")
	
	// 合并统计信息
	return map[string]any{
		"MetadataDB": metadataStats,
		"UserDB":     userStats,
		"Timestamp":  time.Now(),
	}
}

// getSingleDBStats 获取单个数据库连接统计信息
func (dbm *DBManager) getSingleDBStats(db *gorm.DB, dbType string) map[string]any {
	sqlDB, err := db.DB()
	if err != nil {
		return map[string]any{
			"Error":     fmt.Sprintf("Failed to get sql.DB for %s database", dbType),
			"DBType":    dbType,
		}
	}
	stats := sqlDB.Stats()
	return map[string]any{
		"DBType":              dbType,
		"MaxOpenConnections":  stats.MaxOpenConnections,
		"OpenConnections":     stats.OpenConnections,
		"InUse":               stats.InUse,
		"Idle":                stats.Idle,
		"WaitCount":           stats.WaitCount,
		"WaitDuration":        stats.WaitDuration,
		"MaxIdleClosed":       stats.MaxIdleClosed,
		"MaxLifetimeClosed":   stats.MaxLifetimeClosed,
	}
}
