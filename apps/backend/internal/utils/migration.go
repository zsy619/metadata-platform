package utils

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// DialectType 数据库方言类型
type DialectType string

const (
	DialectMySQL     DialectType = "mysql"
	DialectMariaDB   DialectType = "mariadb"
	DialectPostgres  DialectType = "postgres"
	DialectSQLite    DialectType = "sqlite"
	DialectOracle    DialectType = "oracle"
	DialectDM        DialectType = "dm"
	DialectKingbase  DialectType = "kingbase"
	DialectOpenGauss DialectType = "opengauss"
	DialectSQLServer DialectType = "sqlserver"
)

// MigrationHelper 迁移辅助工具
type MigrationHelper struct {
	db      *gorm.DB
	dialect DialectType
}

// NewMigrationHelper 创建迁移辅助工具
func NewMigrationHelper(db *gorm.DB) *MigrationHelper {
	dialectName := strings.ToLower(db.Dialector.Name())
	var dialect DialectType

	switch dialectName {
	case "mysql":
		dialect = DialectMySQL
	case "mariadb":
		dialect = DialectMariaDB
	case "postgres", "postgresql":
		dialect = DialectPostgres
	case "sqlite", "sqlite3":
		dialect = DialectSQLite
	case "oracle":
		dialect = DialectOracle
	case "dm", "dameng":
		dialect = DialectDM
	case "kingbase", "kingbasees":
		dialect = DialectKingbase
	case "opengauss":
		dialect = DialectOpenGauss
	case "sqlserver", "mssql":
		dialect = DialectSQLServer
	default:
		dialect = DialectType(dialectName)
	}

	return &MigrationHelper{
		db:      db,
		dialect: dialect,
	}
}

// GetDialect 获取数据库方言类型
func (h *MigrationHelper) GetDialect() DialectType {
	return h.dialect
}

// IsMySQL 是否为 MySQL/MariaDB
func (h *MigrationHelper) IsMySQL() bool {
	return h.dialect == DialectMySQL || h.dialect == DialectMariaDB
}

// IsPostgres 是否为 PostgreSQL 兼容数据库
func (h *MigrationHelper) IsPostgres() bool {
	return h.dialect == DialectPostgres || h.dialect == DialectKingbase || h.dialect == DialectOpenGauss
}

// IsOracle 是否为 Oracle 兼容数据库
func (h *MigrationHelper) IsOracle() bool {
	return h.dialect == DialectOracle || h.dialect == DialectDM
}

// IsSQLite 是否为 SQLite
func (h *MigrationHelper) IsSQLite() bool {
	return h.dialect == DialectSQLite
}

// GetTableOptions 获取表选项
func (h *MigrationHelper) GetTableOptions() string {
	switch {
	case h.IsMySQL():
		return "ENGINE=InnoDB"
	case h.IsPostgres(), h.IsOracle(), h.IsSQLite():
		return ""
	case h.dialect == DialectSQLServer:
		return ""
	default:
		return ""
	}
}

// AutoMigrate 自动迁移表结构
func (h *MigrationHelper) AutoMigrate(models ...any) error {
	tableOptions := h.GetTableOptions()
	for _, model := range models {
		if err := h.db.Set("gorm:table_options", tableOptions).AutoMigrate(model); err != nil {
			return err
		}
	}
	return nil
}

// AddTableComment 添加表注释
func (h *MigrationHelper) AddTableComment(table, comment string) error {
	var sql string
	escapedComment := strings.ReplaceAll(comment, "'", "''")

	switch {
	case h.IsMySQL():
		sql = fmt.Sprintf("ALTER TABLE `%s` COMMENT = '%s'", table, escapedComment)
	case h.IsPostgres():
		sql = fmt.Sprintf(`COMMENT ON TABLE "%s" IS '%s'`, table, escapedComment)
	case h.IsOracle():
		sql = fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", strings.ToUpper(table), escapedComment)
	case h.dialect == DialectSQLServer:
		sql = fmt.Sprintf("EXEC sp_addextendedproperty 'MS_Description', N'%s', 'SCHEMA', 'dbo', 'TABLE', '%s'", escapedComment, table)
	default:
		return nil
	}

	return h.db.Exec(sql).Error
}

// AddColumnComment 添加列注释
func (h *MigrationHelper) AddColumnComment(table, column, comment string) error {
	var sql string
	escapedComment := strings.ReplaceAll(comment, "'", "''")

	switch {
	case h.IsMySQL():
		return nil
	case h.IsPostgres():
		sql = fmt.Sprintf(`COMMENT ON COLUMN "%s"."%s" IS '%s'`, table, column, escapedComment)
	case h.IsOracle():
		sql = fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s'", strings.ToUpper(table), strings.ToUpper(column), escapedComment)
	case h.dialect == DialectSQLServer:
		sql = fmt.Sprintf("EXEC sp_addextendedproperty 'MS_Description', N'%s', 'SCHEMA', 'dbo', 'TABLE', '%s', 'COLUMN', '%s'", escapedComment, table, column)
	default:
		return nil
	}

	return h.db.Exec(sql).Error
}

// AddComments 批量添加表注释
func (h *MigrationHelper) AddComments(comments map[string]string) error {
	for table, comment := range comments {
		if err := h.AddTableComment(table, comment); err != nil {
			SugarLogger.Warnf("添加表注释失败 %s: %v", table, err)
		}
	}
	return nil
}

// ForeignKeyDef 外键定义
type ForeignKeyDef struct {
	Table      string
	Constraint string
	Column     string
	RefTable   string
	RefColumn  string
	OnDelete   string
	OnUpdate   string
}

// AddForeignKey 添加外键约束
func (h *MigrationHelper) AddForeignKey(fk ForeignKeyDef) error {
	if h.ForeignKeyExists(fk.Table, fk.Constraint) {
		return nil
	}

	if fk.OnDelete == "" {
		fk.OnDelete = "RESTRICT"
	}
	if fk.OnUpdate == "" {
		fk.OnUpdate = "CASCADE"
	}

	var sql string
	switch {
	case h.IsMySQL():
		sql = fmt.Sprintf(
			"ALTER TABLE `%s` ADD CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`) ON DELETE %s ON UPDATE %s",
			fk.Table, fk.Constraint, fk.Column, fk.RefTable, fk.RefColumn, fk.OnDelete, fk.OnUpdate,
		)
	case h.IsPostgres():
		sql = fmt.Sprintf(
			`ALTER TABLE "%s" ADD CONSTRAINT "%s" FOREIGN KEY ("%s") REFERENCES "%s"("%s") ON DELETE %s ON UPDATE %s`,
			fk.Table, fk.Constraint, fk.Column, fk.RefTable, fk.RefColumn, fk.OnDelete, fk.OnUpdate,
		)
	case h.IsOracle():
		sql = fmt.Sprintf(
			`ALTER TABLE "%s" ADD CONSTRAINT "%s" FOREIGN KEY ("%s") REFERENCES "%s"("%s") ON DELETE %s`,
			strings.ToUpper(fk.Table), strings.ToUpper(fk.Constraint), strings.ToUpper(fk.Column),
			strings.ToUpper(fk.RefTable), strings.ToUpper(fk.RefColumn), fk.OnDelete,
		)
	case h.dialect == DialectSQLServer:
		sql = fmt.Sprintf(
			"ALTER TABLE [%s] ADD CONSTRAINT [%s] FOREIGN KEY ([%s]) REFERENCES [%s]([%s]) ON DELETE %s ON UPDATE %s",
			fk.Table, fk.Constraint, fk.Column, fk.RefTable, fk.RefColumn, fk.OnDelete, fk.OnUpdate,
		)
	default:
		return nil
	}

	return h.db.Exec(sql).Error
}

// ForeignKeyExists 检查外键是否存在
func (h *MigrationHelper) ForeignKeyExists(table, constraint string) bool {
	var count int64

	switch {
	case h.IsMySQL():
		h.db.Raw(
			`SELECT COUNT(*) FROM information_schema.TABLE_CONSTRAINTS
			 WHERE CONSTRAINT_SCHEMA = DATABASE()
			   AND TABLE_NAME = ? AND CONSTRAINT_NAME = ? AND CONSTRAINT_TYPE = 'FOREIGN KEY'`,
			table, constraint,
		).Scan(&count)
	case h.IsPostgres():
		h.db.Raw(
			`SELECT COUNT(*) FROM information_schema.table_constraints
			 WHERE table_name = $1 AND constraint_name = $2 AND constraint_type = 'FOREIGN KEY'`,
			table, constraint,
		).Scan(&count)
	case h.IsOracle():
		h.db.Raw(
			`SELECT COUNT(*) FROM USER_CONSTRAINTS
			 WHERE CONSTRAINT_NAME = UPPER(:1) AND CONSTRAINT_TYPE = 'R'`,
			constraint,
		).Scan(&count)
	case h.dialect == DialectSQLServer:
		h.db.Raw(
			`SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
			 WHERE TABLE_NAME = @p1 AND CONSTRAINT_NAME = @p2 AND CONSTRAINT_TYPE = 'FOREIGN KEY'`,
			table, constraint,
		).Scan(&count)
	default:
		return false
	}

	return count > 0
}

// AddForeignKeys 批量添加外键
func (h *MigrationHelper) AddForeignKeys(fks []ForeignKeyDef) error {
	for _, fk := range fks {
		if err := h.AddForeignKey(fk); err != nil {
			SugarLogger.Warnf("添加外键 %s 失败: %v", fk.Constraint, err)
		}
	}
	return nil
}

// DropIndex 删除索引
func (h *MigrationHelper) DropIndex(table, indexName string) error {
	var sql string

	switch {
	case h.IsMySQL():
		sql = fmt.Sprintf("ALTER TABLE `%s` DROP INDEX `%s`", table, indexName)
	case h.IsPostgres():
		sql = fmt.Sprintf("DROP INDEX IF EXISTS %s", indexName)
	case h.IsOracle():
		sql = fmt.Sprintf("DROP INDEX %s", strings.ToUpper(indexName))
	case h.dialect == DialectSQLServer:
		sql = fmt.Sprintf("DROP INDEX [%s] ON [%s]", indexName, table)
	default:
		sql = fmt.Sprintf("DROP INDEX IF EXISTS %s", indexName)
	}

	return h.db.Exec(sql).Error
}

// CreateIndex 创建索引
func (h *MigrationHelper) CreateIndex(table, indexName string, columns []string, unique bool) error {
	uniqueStr := ""
	if unique {
		uniqueStr = "UNIQUE "
	}

	columnsStr := ""
	switch {
	case h.IsMySQL():
		for i, col := range columns {
			if i > 0 {
				columnsStr += ", "
			}
			columnsStr += fmt.Sprintf("`%s`", col)
		}
		sql := fmt.Sprintf("CREATE %sINDEX `%s` ON `%s` (%s)", uniqueStr, indexName, table, columnsStr)
		return h.db.Exec(sql).Error
	case h.IsPostgres():
		for i, col := range columns {
			if i > 0 {
				columnsStr += ", "
			}
			columnsStr += fmt.Sprintf(`"%s"`, col)
		}
		sql := fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS %s ON %s (%s)", uniqueStr, indexName, table, columnsStr)
		return h.db.Exec(sql).Error
	case h.IsOracle():
		for i, col := range columns {
			if i > 0 {
				columnsStr += ", "
			}
			columnsStr += fmt.Sprintf(`"%s"`, strings.ToUpper(col))
		}
		sql := fmt.Sprintf("CREATE %sINDEX %s ON %s (%s)", uniqueStr, strings.ToUpper(indexName), strings.ToUpper(table), columnsStr)
		return h.db.Exec(sql).Error
	default:
		return nil
	}
}

// QuoteIdentifier 引用标识符
func (h *MigrationHelper) QuoteIdentifier(name string) string {
	switch {
	case h.IsMySQL():
		return fmt.Sprintf("`%s`", name)
	case h.IsPostgres(), h.IsOracle():
		return fmt.Sprintf(`"%s"`, name)
	case h.dialect == DialectSQLServer:
		return fmt.Sprintf("[%s]", name)
	default:
		return name
	}
}
