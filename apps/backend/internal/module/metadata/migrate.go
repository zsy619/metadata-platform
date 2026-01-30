package metadata

import (
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/utils"

	"gorm.io/gorm"
)

// Migrate 自动迁移元数据模块表结构
func Migrate(db *gorm.DB) error {
	utils.SugarLogger.Info("Starting metadata migration...")
	var err error
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='API接口'").AutoMigrate(&model.API{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='数据连接'").AutoMigrate(&model.MdConn{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='物理表'").AutoMigrate(&model.MdTable{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='物理表字段'").AutoMigrate(&model.MdTableField{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='模型定义'").AutoMigrate(&model.MdModel{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='模型字段'").AutoMigrate(&model.MdModelField{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='模型字段分组'").AutoMigrate(&model.MdModelGroup{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='模型聚合过滤'").AutoMigrate(&model.MdModelHaving{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='模型关联JOIN'").AutoMigrate(&model.MdModelJoin{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='模型关联JOIN字段'").AutoMigrate(&model.MdModelJoinField{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='模型分页限制'").AutoMigrate(&model.MdModelLimit{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='模型排序'").AutoMigrate(&model.MdModelOrder{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='模型自定义SQL'").AutoMigrate(&model.MdModelSql{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='模型关联物理表'").AutoMigrate(&model.MdModelTable{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='模型筛选条件'").AutoMigrate(&model.MdModelWhere{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='模型字段增强'").AutoMigrate(&model.MdModelFieldEnhancement{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='模型间关联'").AutoMigrate(&model.MdModelRelation{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='查询模板'").AutoMigrate(&model.MdQueryTemplate{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='查询条件'").AutoMigrate(&model.MdQueryCondition{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='模型参数'").AutoMigrate(&model.MdModelParam{}); err != nil {
		return err
	}
	return nil
}
