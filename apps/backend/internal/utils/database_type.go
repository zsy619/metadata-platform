package utils

import (
	"strings"
)

// DatabaseType 数据库类型常量
const (
	DBTypeMySQL        = "MySQL"
	DBTypePostgreSQL   = "PostgreSQL"
	DBTypeSQLServer    = "SQL Server"
	DBTypeOracle       = "Oracle"
	DBTypeSQLite       = "SQLite"
	DBTypeClickHouse   = "ClickHouse"
	DBTypeDM           = "DM"
	DBTypeMongoDB      = "MongoDB"
	DBTypeRedis        = "Redis"
	DBTypeTiDB         = "TiDB"
	DBTypeOceanBase    = "OceanBase"
	DBTypeDoris        = "Doris"
	DBTypeStarRocks    = "StarRocks"
	DBTypeOpenGauss    = "OpenGauss"
	DBTypeKingbase     = "Kingbase"
)

// IsMySQL 是否为 MySQL 兼容数据库
func IsMySQL(dbType string) bool {
	normalized := NormalizeDBType(dbType)
	return normalized == DBTypeMySQL ||
		normalized == DBTypeTiDB ||
		normalized == DBTypeOceanBase ||
		normalized == DBTypeDoris ||
		normalized == DBTypeStarRocks
}

// IsPostgreSQL 是否为 PostgreSQL 兼容数据库
func IsPostgreSQL(dbType string) bool {
	normalized := NormalizeDBType(dbType)
	return normalized == DBTypePostgreSQL ||
		normalized == DBTypeOpenGauss ||
		normalized == DBTypeKingbase
}

// NormalizeDBType 标准化数据库类型（不区分大小写）
func NormalizeDBType(dbType string) string {
	switch strings.ToLower(strings.TrimSpace(dbType)) {
	case "mysql", "mariadb":
		return DBTypeMySQL
	case "postgres", "postgresql", "pg":
		return DBTypePostgreSQL
	case "sqlserver", "sql server", "mssql":
		return DBTypeSQLServer
	case "oracle":
		return DBTypeOracle
	case "sqlite", "sqlite3":
		return DBTypeSQLite
	case "clickhouse":
		return DBTypeClickHouse
	case "dm", "dameng":
		return DBTypeDM
	case "mongodb", "mongo":
		return DBTypeMongoDB
	case "redis":
		return DBTypeRedis
	case "tidb":
		return DBTypeTiDB
	case "oceanbase":
		return DBTypeOceanBase
	case "doris":
		return DBTypeDoris
	case "starrocks":
		return DBTypeStarRocks
	case "opengauss":
		return DBTypeOpenGauss
	case "kingbase":
		return DBTypeKingbase
	default:
		return dbType
	}
}

// GetDBTypeAliases 获取数据库类型别名列表
func GetDBTypeAliases(dbType string) []string {
	switch NormalizeDBType(dbType) {
	case DBTypeMySQL:
		return []string{DBTypeMySQL, "MariaDB", DBTypeTiDB, DBTypeOceanBase, DBTypeDoris, DBTypeStarRocks}
	case DBTypePostgreSQL:
		return []string{DBTypePostgreSQL, DBTypeOpenGauss, DBTypeKingbase}
	case DBTypeSQLServer:
		return []string{DBTypeSQLServer, "MSSQL"}
	case DBTypeSQLite:
		return []string{DBTypeSQLite, "SQLite3"}
	case DBTypeDM:
		return []string{DBTypeDM, "Dameng"}
	case DBTypeMongoDB:
		return []string{DBTypeMongoDB, "Mongo"}
	default:
		return []string{dbType}
	}
}
