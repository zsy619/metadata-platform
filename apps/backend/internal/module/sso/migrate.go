package sso

import (
	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/utils"

	"gorm.io/gorm"
)

// Migrate 自动迁移SSO模块表结构
func Migrate(db *gorm.DB) error {
	utils.SugarLogger.Info("Starting SSO migration...")

	helper := utils.NewMigrationHelper(db)
	var err error

	// 使用 MigrationHelper 自动迁移
	models := []interface{}{
		&model.SsoProtocolConfig{},
		&model.SsoClient{},
		&model.SsoKey{},
		&model.SsoFieldMapping{},
		&model.SsoSession{},
	}

	if err = helper.AutoMigrate(models...); err != nil {
		utils.SugarLogger.Errorf("SSO AutoMigrate failed: %v", err)
		return err
	}

	// 添加表注释
	comments := map[string]string{
		"sso_protocol_config": "SSO协议配置",
		"sso_client":          "SSO客户端配置",
		"sso_key":             "SSO密钥管理",
		"sso_field_mapping":   "SSO字段映射配置",
		"sso_session":         "SSO会话",
	}
	helper.AddComments(comments)

	// 创建外键关联关系
	utils.SugarLogger.Info("Creating SSO foreign key constraints...")
	fks := []utils.ForeignKeyDef{
		{Table: "sso_client", Constraint: "fk_sso_client_protocol_config_id", Column: "protocol_config_id", RefTable: "sso_protocol_config", RefColumn: "id"},
		{Table: "sso_key", Constraint: "fk_sso_key_protocol_config_id", Column: "protocol_config_id", RefTable: "sso_protocol_config", RefColumn: "id"},
		{Table: "sso_field_mapping", Constraint: "fk_sso_field_mapping_protocol_config_id", Column: "protocol_config_id", RefTable: "sso_protocol_config", RefColumn: "id"},
		{Table: "sso_field_mapping", Constraint: "fk_sso_field_mapping_client_id", Column: "client_id", RefTable: "sso_client", RefColumn: "client_id"},
		{Table: "sso_session", Constraint: "fk_sso_session_user_id", Column: "user_id", RefTable: "sso_user", RefColumn: "id"},
		{Table: "sso_session", Constraint: "fk_sso_session_client_id", Column: "client_id", RefTable: "sso_client", RefColumn: "client_id"},
		{Table: "sso_session", Constraint: "fk_sso_session_protocol_config_id", Column: "protocol_config_id", RefTable: "sso_protocol_config", RefColumn: "id"},
	}
	helper.AddForeignKeys(fks)

	utils.SugarLogger.Info("SSO migration completed successfully")
	return nil
}
