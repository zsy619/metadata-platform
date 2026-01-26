package user

import (
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/utils"

	"gorm.io/gorm"
)

// Migrate 自动迁移用户模块表结构
func Migrate(db *gorm.DB) error {
	utils.SugarLogger.Info("Starting user migration...")
	return db.AutoMigrate(
		&model.SsoCasbinRule{},
		&model.SsoTenant{},
		&model.SsoApplication{},
		&model.SsoMenu{},
		&model.SsoOrganization{},
		&model.SsoOrganizationMenu{},
		&model.SsoOrganizationRole{},
		&model.SsoOrganizationKind{},
		&model.SsoOrganizationKindRole{},
		&model.SsoOrganizationUser{},
		&model.SsoPosition{},
		&model.SsoPositionRole{},
		&model.SsoRole{},
		&model.SsoRoleGroup{},
		&model.SsoRoleGroupRole{},
		&model.SsoRoleMenu{},
		&model.SsoUser{},
		&model.SsoUserGroup{},
		&model.SsoUserGroupUser{},
		&model.SsoUserPosition{},
		&model.SsoUserRole{},
		&model.SsoUserRoleGroup{},
	)
}
