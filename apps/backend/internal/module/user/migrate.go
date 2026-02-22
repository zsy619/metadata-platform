package user

import (
	"fmt"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/utils"
	"strings"

	"gorm.io/gorm"
)

// Migrate 自动迁移用户模块表结构
func Migrate(db *gorm.DB) error {
	utils.SugarLogger.Info("Starting user migration...")

	var err error
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='用户职位关联'").AutoMigrate(&model.SsoUserPos{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='用户角色关联'").AutoMigrate(&model.SsoUserRole{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='Casbin规则'").AutoMigrate(&model.SsoCasbinRule{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='租户'").AutoMigrate(&model.SsoTenant{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='应用服务'").AutoMigrate(&model.SsoApp{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='菜单权限'").AutoMigrate(&model.SsoMenu{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='组织机构'").AutoMigrate(&model.SsoOrg{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='组织菜单关联'").AutoMigrate(&model.SsoOrgMenu{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='组织角色关联'").AutoMigrate(&model.SsoOrgRole{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='组织类型'").AutoMigrate(&model.SsoOrgKind{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='组织类型角色关联'").AutoMigrate(&model.SsoOrgKindRole{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='组织用户关联'").AutoMigrate(&model.SsoOrgUser{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='职位'").AutoMigrate(&model.SsoPos{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='职位角色关联'").AutoMigrate(&model.SsoPosRole{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='角色'").AutoMigrate(&model.SsoRole{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='角色组'").AutoMigrate(&model.SsoRoleGroup{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='角色组角色关联'").AutoMigrate(&model.SsoRoleGroupRole{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='角色菜单关联'").AutoMigrate(&model.SsoRoleMenu{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='用户'").AutoMigrate(&model.SsoUser{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='用户组'").AutoMigrate(&model.SsoUserGroup{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='用户组用户关联'").AutoMigrate(&model.SsoUserGroupUser{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='用户组角色关联'").AutoMigrate(&model.SsoUserGroupRole{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='用户角色组关联'").AutoMigrate(&model.SsoUserRoleGroup{}); err != nil {
		return err
	}

	// 配置多对多关联表模型（在 AutoMigrate 之后注册，避免 GORM 以 join table 路径处理导致 COMMENT 丢失）
	_ = db.SetupJoinTable(&model.SsoUser{}, "Roles", &model.SsoUserRole{})
	_ = db.SetupJoinTable(&model.SsoUser{}, "Positions", &model.SsoUserPos{})
	_ = db.SetupJoinTable(&model.SsoUser{}, "Organizations", &model.SsoOrgUser{})

	// 根据结构，创建所有的外键关联关系
	utils.SugarLogger.Info("Creating foreign key constraints...")
	type fkDef struct {
		table      string
		constraint string
		column     string
		refTable   string
		refColumn  string
	}
	fks := []fkDef{
		// sso_org_menu
		{"sso_org_menu", "fk_org_menu_org_id", "org_id", "sso_org", "id"},
		{"sso_org_menu", "fk_org_menu_menu_id", "menu_id", "sso_menu", "id"},
		// sso_org_role
		{"sso_org_role", "fk_org_role_org_id", "org_id", "sso_org", "id"},
		{"sso_org_role", "fk_org_role_role_id", "role_id", "sso_role", "id"},
		// sso_org_user
		{"sso_org_user", "fk_org_user_org_id", "org_id", "sso_org", "id"},
		{"sso_org_user", "fk_org_user_user_id", "user_id", "sso_user", "id"},
		// sso_org_kind_role（kind_code 引用 sso_org_kind 唯一索引字段）
		{"sso_org_kind_role", "fk_org_kind_role_kind_code", "kind_code", "sso_org_kind", "kind_code"},
		{"sso_org_kind_role", "fk_org_kind_role_role_id", "role_id", "sso_role", "id"},
		// sso_pos_role
		{"sso_pos_role", "fk_pos_role_pos_id", "pos_id", "sso_pos", "id"},
		{"sso_pos_role", "fk_pos_role_role_id", "role_id", "sso_role", "id"},
		// sso_role_menu
		{"sso_role_menu", "fk_role_menu_role_id", "role_id", "sso_role", "id"},
		{"sso_role_menu", "fk_role_menu_menu_id", "menu_id", "sso_menu", "id"},
		// sso_role_group_role
		{"sso_role_group_role", "fk_role_group_role_group_id", "group_id", "sso_role_group", "id"},
		{"sso_role_group_role", "fk_role_group_role_role_id", "role_id", "sso_role", "id"},
		// sso_user_pos
		{"sso_user_pos", "fk_user_pos_user_id", "user_id", "sso_user", "id"},
		{"sso_user_pos", "fk_user_pos_pos_id", "pos_id", "sso_pos", "id"},
		// sso_user_role
		{"sso_user_role", "fk_user_role_user_id", "user_id", "sso_user", "id"},
		{"sso_user_role", "fk_user_role_role_id", "role_id", "sso_role", "id"},
		// sso_user_role_group
		{"sso_user_role_group", "fk_user_role_group_user_id", "user_id", "sso_user", "id"},
		{"sso_user_role_group", "fk_user_role_group_group_id", "group_id", "sso_role_group", "id"},
		// sso_user_group_role
		{"sso_user_group_role", "fk_user_group_role_group_id", "group_id", "sso_user_group", "id"},
		{"sso_user_group_role", "fk_user_group_role_role_id", "role_id", "sso_role", "id"},
		// sso_user_group_user
		{"sso_user_group_user", "fk_user_group_user_group_id", "group_id", "sso_user_group", "id"},
		{"sso_user_group_user", "fk_user_group_user_user_id", "user_id", "sso_user", "id"},
	}
	for _, fk := range fks {
		if e := addForeignKey(db, fk.table, fk.constraint, fk.column, fk.refTable, fk.refColumn); e != nil {
			utils.SugarLogger.Warnf("添加外键 %s 失败: %v", fk.constraint, e)
		}
	}

	return nil
}

// addForeignKey 幂等地添加外键约束，若已存在则跳过。
// 支持数据库：MySQL/MariaDB、PostgreSQL、人大金仓(KingbaseES)、openGauss、Oracle、达梦(DM)。
func addForeignKey(db *gorm.DB, table, constraint, column, refTable, refColumn string) error {
	dialect := db.Dialector.Name()
	dialect = strings.ToLower(dialect)

	// ── 1. 检查外键是否已存在 ──────────────────────────────────────────────────
	var exists bool
	switch dialect {
	case "mysql", "mariadb":
		var count int64
		db.Raw(
			`SELECT COUNT(*) FROM information_schema.TABLE_CONSTRAINTS
			 WHERE CONSTRAINT_SCHEMA = DATABASE()
			   AND TABLE_NAME = ? AND CONSTRAINT_NAME = ? AND CONSTRAINT_TYPE = 'FOREIGN KEY'`,
			table, constraint,
		).Scan(&count)
		exists = count > 0

	case "oracle", "dm": // 达梦兼容 Oracle 数据字典
		var count int64
		db.Raw(
			`SELECT COUNT(*) FROM USER_CONSTRAINTS
			 WHERE CONSTRAINT_NAME = UPPER(:1) AND CONSTRAINT_TYPE = 'R'`,
			constraint,
		).Scan(&count)
		exists = count > 0

	default: // postgres、kingbase(人大金仓)、opengauss 等 PG 兼容方言
		var count int64
		db.Raw(
			`SELECT COUNT(*) FROM information_schema.table_constraints
			 WHERE table_name = $1 AND constraint_name = $2 AND constraint_type = 'FOREIGN KEY'`,
			table, constraint,
		).Scan(&count)
		exists = count > 0
	}

	if exists {
		return nil // 已存在，跳过
	}

	// ── 2. 构造 ALTER TABLE 语句 ────────────────────────────────────────────────
	var ddl string
	switch dialect {
	case "mysql", "mariadb":
		// MySQL/MariaDB：反引号标识符，支持 ON UPDATE CASCADE
		ddl = fmt.Sprintf(
			"ALTER TABLE `%s` ADD CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`) ON DELETE RESTRICT ON UPDATE CASCADE",
			table, constraint, column, refTable, refColumn,
		)
	case "oracle", "dm":
		// Oracle/达梦：双引号标识符，不支持 ON UPDATE CASCADE
		ddl = fmt.Sprintf(
			`ALTER TABLE "%s" ADD CONSTRAINT "%s" FOREIGN KEY ("%s") REFERENCES "%s"("%s") ON DELETE RESTRICT`,
			table, constraint, column, refTable, refColumn,
		)
	default:
		// PostgreSQL / 人大金仓 / openGauss：双引号标识符，支持 ON UPDATE CASCADE
		ddl = fmt.Sprintf(
			`ALTER TABLE "%s" ADD CONSTRAINT "%s" FOREIGN KEY ("%s") REFERENCES "%s"("%s") ON DELETE RESTRICT ON UPDATE CASCADE`,
			table, constraint, column, refTable, refColumn,
		)
	}

	return db.Exec(ddl).Error
}
