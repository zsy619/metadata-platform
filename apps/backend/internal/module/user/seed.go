package user

import (
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/utils"
	"time"

	"gorm.io/gorm"
)

// SeedData 初始化用户模块种子数据
func SeedData(db *gorm.DB) {
	utils.SugarLogger.Info("Seeding user database...")

	// sf := utils.NewSnowflake(1, 1)
	now := time.Now()

	// 1. 初始化超管租户
	defaultTenant := model.SsoTenant{
		ID:         utils.SystemTenantID,
		TenantID:   utils.SystemTenantID,
		TenantName: "system",
		TenantCode: "system",
		Status:     1,
		IsSystem:   true,
		CreateBy:   "system",
		CreateAt:   now,
		UpdateBy:   "system",
		UpdateAt:   now,
	}

	var count int64
	db.Model(&model.SsoTenant{}).Where("id = ?", "1").Count(&count)
	if count == 0 {
		if err := db.Create(&defaultTenant).Error; err != nil {
			utils.SugarLogger.Errorf("Failed to seed default tenant: %v", err)
		} else {
			utils.SugarLogger.Info("Seeded default tenant")
		}
	}

	// 2. 初始化管理员账号
	salt := utils.GenerateSalt()
	adminUser := model.SsoUser{
		ID:       utils.SuperAdminID,
		TenantID: utils.SystemTenantID,
		Account:  "admin",
		Password: utils.EncryptPasswordSM3("Admin@2026", salt),
		Salt:     salt,
		Name:     "系统管理员",
		Status:   1,
		Kind:     1, // 1: 系统管理员
		IsSystem: true,
		CreateBy: "system",
		CreateAt: now,
		UpdateBy: "system",
		UpdateAt: now,
	}

	db.Model(&model.SsoUser{}).Where("account = ?", "admin").Count(&count)
	if count == 0 {
		if err := db.Create(&adminUser).Error; err != nil {
			utils.SugarLogger.Errorf("Failed to seed admin user: %v", err)
		} else {
			utils.SugarLogger.Info("Seeded admin user")
		}
	}

	// 3. 初始化默认应用
	apps := []model.SsoApp{
		{
			ID:       "1",
			TenantID: utils.SystemTenantID,
			AppName:  "元数据管理",
			AppCode:  "metadata",
			Status:   1,
			IsSystem: true,
			Remark:   "元数据管理",
			CreateBy: "system",
			CreateAt: now,
			UpdateBy: "system",
			UpdateAt: now,
		},
		{
			ID:       "2",
			TenantID: utils.SystemTenantID,
			AppName:  "用户管理",
			AppCode:  "sso",
			Status:   1,
			IsSystem: true,
			Remark:   "账号认证与权限管理系统",
			CreateBy: "system",
			CreateAt: now,
			UpdateBy: "system",
			UpdateAt: now,
		},
	}

	for _, app := range apps {
		db.Model(&model.SsoApp{}).Where("app_code = ?", app.AppCode).Count(&count)
		if count == 0 {
			if err := db.Create(&app).Error; err != nil {
				utils.SugarLogger.Errorf("Failed to seed app %s: %v", app.AppName, err)
			} else {
				utils.SugarLogger.Infof("Seeded app: %s", app.AppName)
			}
		}
	}

	// 4. 初始化默认角色
	roles := []model.SsoRole{
		{
			ID:        utils.SuperRoleID,
			TenantID:  utils.SystemTenantID,
			RoleName:  "超级管理员",
			RoleCode:  "super_role",
			Status:    1,
			DataRange: utils.DataRangeAll,
			Remark:    "系统最高权限管理员",
			IsSystem:  true,
			CreateBy:  "system",
			CreateAt:  now,
			UpdateBy:  "system",
			UpdateAt:  now,
		},
	}

	for _, role := range roles {
		db.Model(&model.SsoRole{}).Where("role_code = ?", role.RoleCode).Count(&count)
		if count == 0 {
			if err := db.Create(&role).Error; err != nil {
				utils.SugarLogger.Errorf("Failed to seed role %s: %v", role.RoleName, err)
			} else {
				utils.SugarLogger.Infof("Seeded role: %s", role.RoleName)
			}
		}
	}

	// 5. 初始化默认用户角色
	userRoles := []model.SsoUserRole{
		{
			ID:       "1",
			TenantID: utils.SystemTenantID,
			UserID:   utils.SuperAdminID,
			RoleID:   utils.SuperRoleID,
			IsSystem: true,
			CreateBy: "system",
			CreateAt: now,
			UpdateBy: "system",
			UpdateAt: now,
		},
	}

	for _, userRole := range userRoles {
		db.Model(&model.SsoUserRole{}).Where("user_id = ? AND role_id = ?", userRole.UserID, userRole.RoleID).Count(&count)
		if count == 0 {
			if err := db.Create(&userRole).Error; err != nil {
				utils.SugarLogger.Errorf("Failed to seed user role %s: %v", userRole.UserID, err)
			} else {
				utils.SugarLogger.Infof("Seeded user role: %s", userRole.UserID)
			}
		}
	}

	utils.SugarLogger.Info("User database seeding completed")
}
