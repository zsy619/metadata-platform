package audit

import (
	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/utils"

	"gorm.io/gorm"
)

// Migrate 自动迁移审计模块表结构
func Migrate(db *gorm.DB) error {
	utils.SugarLogger.Info("Starting audit migration...")
	var err error
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='系统操作日志'").AutoMigrate(&model.SysOperationLog{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='系统数据变更日志'").AutoMigrate(&model.SysDataChangeLog{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='系统登录日志'").AutoMigrate(&model.SysLoginLog{}); err != nil {
		return err
	}
	return nil
}
