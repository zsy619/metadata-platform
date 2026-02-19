package user

import (
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/utils"

	"gorm.io/gorm"
)

// Migrate 自动迁移用户模块表结构
func Migrate(db *gorm.DB) error {
	utils.SugarLogger.Info("Starting user migration...")

	// 配置多对多关联表模型
	_ = db.SetupJoinTable(&model.SsoUser{}, "Roles", &model.SsoUserRole{})
	_ = db.SetupJoinTable(&model.SsoUser{}, "Positions", &model.SsoPos{})
	_ = db.SetupJoinTable(&model.SsoUser{}, "Organizations", &model.SsoOrgUser{})

	var err error
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
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='用户职位关联'").AutoMigrate(&model.SsoUserPos{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='用户角色关联'").AutoMigrate(&model.SsoUserRole{}); err != nil {
		return err
	}
	if err = db.Set("gorm:table_options", "ENGINE=InnoDB COMMENT='用户角色组关联'").AutoMigrate(&model.SsoUserRoleGroup{}); err != nil {
		return err
	}
	return nil
}
