package user

import (
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/utils"

	"gorm.io/gorm"
)

// Migrate 自动迁移用户模块表结构
func Migrate(db *gorm.DB) error {
	utils.SugarLogger.Info("Starting user migration...")

	helper := utils.NewMigrationHelper(db)
	var err error

	// 使用 MigrationHelper 自动迁移
	models := []interface{}{
		&model.SsoUserProfile{},
		&model.SsoUserAddress{},
		&model.SsoUserContact{},
		&model.SsoUserSocial{},
		&model.SsoUserPos{},
		&model.SsoUserRole{},
		&model.SsoCasbinRule{},
		&model.SsoTenant{},
		&model.SsoApp{},
		&model.SsoMenu{},
		&model.SsoOrg{},
		&model.SsoOrgMenu{},
		&model.SsoOrgRole{},
		&model.SsoOrgKind{},
		&model.SsoOrgKindRole{},
		&model.SsoOrgUser{},
		&model.SsoPos{},
		&model.SsoPosRole{},
		&model.SsoRole{},
		&model.SsoRoleGroup{},
		&model.SsoRoleGroupRole{},
		&model.SsoRoleMenu{},
		&model.SsoUser{},
		&model.SsoUserGroup{},
		&model.SsoUserGroupUser{},
		&model.SsoUserGroupRole{},
		&model.SsoUserRoleGroup{},
	}

	if err = helper.AutoMigrate(models...); err != nil {
		utils.SugarLogger.Errorf("AutoMigrate failed: %v", err)
		return err
	}

	// 配置多对多关联表模型
	_ = db.SetupJoinTable(&model.SsoUser{}, "Roles", &model.SsoUserRole{})
	_ = db.SetupJoinTable(&model.SsoUser{}, "Positions", &model.SsoUserPos{})
	_ = db.SetupJoinTable(&model.SsoUser{}, "Organizations", &model.SsoOrgUser{})

	// 添加表注释
	comments := map[string]string{
		"sso_user_profile":      "用户档案",
		"sso_user_address":      "用户地址簿",
		"sso_user_contact":      "用户联系方式",
		"sso_user_social":       "用户第三方账号",
		"sso_user_pos":          "用户职位关联",
		"sso_user_role":         "用户角色关联",
		"sso_casbin_rule":       "Casbin规则",
		"sso_tenant":            "租户",
		"sso_app":               "应用服务",
		"sso_menu":              "菜单权限",
		"sso_org":               "组织机构",
		"sso_org_menu":          "组织菜单关联",
		"sso_org_role":          "组织角色关联",
		"sso_org_kind":          "组织类型",
		"sso_org_kind_role":     "组织类型角色关联",
		"sso_org_user":          "组织用户关联",
		"sso_pos":               "职位",
		"sso_pos_role":          "职位角色关联",
		"sso_role":              "角色",
		"sso_role_group":        "角色组",
		"sso_role_group_role":   "角色组角色关联",
		"sso_role_menu":         "角色菜单关联",
		"sso_user":              "用户",
		"sso_user_group":        "用户组",
		"sso_user_group_user":   "用户组用户关联",
		"sso_user_group_role":   "用户组角色关联",
		"sso_user_role_group":   "用户角色组关联",
	}
	helper.AddComments(comments)

	// 创建外键关联关系
	utils.SugarLogger.Info("Creating foreign key constraints...")
	fks := []utils.ForeignKeyDef{
		{Table: "sso_org_menu", Constraint: "fk_org_menu_org_id", Column: "org_id", RefTable: "sso_org", RefColumn: "id"},
		{Table: "sso_org_menu", Constraint: "fk_org_menu_menu_id", Column: "menu_id", RefTable: "sso_menu", RefColumn: "id"},
		{Table: "sso_org_role", Constraint: "fk_org_role_org_id", Column: "org_id", RefTable: "sso_org", RefColumn: "id"},
		{Table: "sso_org_role", Constraint: "fk_org_role_role_id", Column: "role_id", RefTable: "sso_role", RefColumn: "id"},
		{Table: "sso_org_user", Constraint: "fk_org_user_org_id", Column: "org_id", RefTable: "sso_org", RefColumn: "id"},
		{Table: "sso_org_user", Constraint: "fk_org_user_user_id", Column: "user_id", RefTable: "sso_user", RefColumn: "id"},
		{Table: "sso_org_kind_role", Constraint: "fk_org_kind_role_kind_code", Column: "kind_code", RefTable: "sso_org_kind", RefColumn: "kind_code"},
		{Table: "sso_org_kind_role", Constraint: "fk_org_kind_role_role_id", Column: "role_id", RefTable: "sso_role", RefColumn: "id"},
		{Table: "sso_pos_role", Constraint: "fk_pos_role_pos_id", Column: "pos_id", RefTable: "sso_pos", RefColumn: "id"},
		{Table: "sso_pos_role", Constraint: "fk_pos_role_role_id", Column: "role_id", RefTable: "sso_role", RefColumn: "id"},
		{Table: "sso_role_menu", Constraint: "fk_role_menu_role_id", Column: "role_id", RefTable: "sso_role", RefColumn: "id"},
		{Table: "sso_role_menu", Constraint: "fk_role_menu_menu_id", Column: "menu_id", RefTable: "sso_menu", RefColumn: "id"},
		{Table: "sso_role_group_role", Constraint: "fk_role_group_role_group_id", Column: "group_id", RefTable: "sso_role_group", RefColumn: "id"},
		{Table: "sso_role_group_role", Constraint: "fk_role_group_role_role_id", Column: "role_id", RefTable: "sso_role", RefColumn: "id"},
		{Table: "sso_user_pos", Constraint: "fk_user_pos_user_id", Column: "user_id", RefTable: "sso_user", RefColumn: "id"},
		{Table: "sso_user_pos", Constraint: "fk_user_pos_pos_id", Column: "pos_id", RefTable: "sso_pos", RefColumn: "id"},
		{Table: "sso_user_role", Constraint: "fk_user_role_user_id", Column: "user_id", RefTable: "sso_user", RefColumn: "id"},
		{Table: "sso_user_role", Constraint: "fk_user_role_role_id", Column: "role_id", RefTable: "sso_role", RefColumn: "id"},
		{Table: "sso_user_role_group", Constraint: "fk_user_role_group_user_id", Column: "user_id", RefTable: "sso_user", RefColumn: "id"},
		{Table: "sso_user_role_group", Constraint: "fk_user_role_group_group_id", Column: "group_id", RefTable: "sso_role_group", RefColumn: "id"},
		{Table: "sso_user_group_role", Constraint: "fk_user_group_role_group_id", Column: "group_id", RefTable: "sso_user_group", RefColumn: "id"},
		{Table: "sso_user_group_role", Constraint: "fk_user_group_role_role_id", Column: "role_id", RefTable: "sso_role", RefColumn: "id"},
		{Table: "sso_user_group_user", Constraint: "fk_user_group_user_group_id", Column: "group_id", RefTable: "sso_user_group", RefColumn: "id"},
		{Table: "sso_user_group_user", Constraint: "fk_user_group_user_user_id", Column: "user_id", RefTable: "sso_user", RefColumn: "id"},
	}
	helper.AddForeignKeys(fks)

	utils.SugarLogger.Info("User migration completed successfully")
	return nil
}
