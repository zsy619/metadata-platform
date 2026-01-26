package metadata

import (
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/utils"

	"gorm.io/gorm"
)

// Migrate 自动迁移元数据模块表结构
func Migrate(db *gorm.DB) error {
	utils.SugarLogger.Info("Starting metadata migration...")
	return db.AutoMigrate(
		&model.API{},
		&model.MdConn{},
		&model.MdTable{},
		&model.MdTableField{},
		&model.MdModel{},
		&model.MdModelField{},
		&model.MdModelGroup{},
		&model.MdModelHaving{},
		&model.MdModelJoin{},
		&model.MdModelLimit{},
		&model.MdModelOrder{},
		&model.MdModelSql{},
		&model.MdModelTable{},
		&model.MdModelWhere{},
		&model.MdModelFieldEnhancement{},
	)
}
