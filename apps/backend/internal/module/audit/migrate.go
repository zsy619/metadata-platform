package audit

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/utils"
)

// Migrate 自动迁移审计模块表结构
func Migrate(db *gorm.DB) error {
	utils.SugarLogger.Info("Starting audit migration...")

	helper := utils.NewMigrationHelper(db)
	var err error

	// 使用 MigrationHelper 自动迁移
	models := []any{
		&model.SysOperationLog{},
		&model.SysDataChangeLog{},
		&model.SysLoginLog{},
		&model.SysAccessLog{},
	}

	if err = helper.AutoMigrate(models...); err != nil {
		panic(err)
	}

	// 添加表注释
	comments := map[string]string{
		"sys_operation_log":   "系统操作日志",
		"sys_data_change_log": "系统数据变更日志",
		"sys_login_log":       "系统登录日志",
		"sys_access_log":      "系统访问日志",
	}
	helper.AddComments(comments)

	utils.SugarLogger.Info("Audit migration completed successfully")
	return nil
}
