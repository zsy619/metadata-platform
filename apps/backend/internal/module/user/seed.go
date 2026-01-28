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
		ID:         "1",
		TenantID:   "0",
		TenantName: "系统管理组",
		TenantCode: "system",
		State:      1,
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
		ID:       "1",
		TenantID: "1",
		Account:  "admin",
		Password: utils.EncryptPasswordSM3("admin@123", salt),
		Salt:     salt,
		Name:     "系统管理员",
		State:    1,
		Kind:     1, // 1: 系统管理员
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
	apps := []model.SsoApplication{
		{
			ID:              "1",
			ApplicationName: "元数据管理",
			ApplicationCode: "metadata",
			State:           1,
			Remark:          "元数据管理",
			CreateBy:        "system",
			CreateAt:        now,
			UpdateBy:        "system",
			UpdateAt:        now,
		},
		{
			ID:              "2",
			ApplicationName: "用户管理",
			ApplicationCode: "sso",
			State:           1,
			Remark:          "账号认证与权限管理系统",
			CreateBy:        "system",
			CreateAt:        now,
			UpdateBy:        "system",
			UpdateAt:        now,
		},
	}

	for _, app := range apps {
		db.Model(&model.SsoApplication{}).Where("application_code = ?", app.ApplicationCode).Count(&count)
		if count == 0 {
			if err := db.Create(&app).Error; err != nil {
				utils.SugarLogger.Errorf("Failed to seed application %s: %v", app.ApplicationName, err)
			} else {
				utils.SugarLogger.Infof("Seeded application: %s", app.ApplicationName)
			}
		}
	}

	// 4. 初始化默认角色
	roles := []model.SsoRole{
		{
			ID:        "1",
			TenantID:  "1",
			RoleName:  "超级管理员",
			RoleCode:  "super_role",
			State:     1,
			DataScope: "1", // 1: 全部数据权限
			Remark:    "系统最高权限管理员",
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
			TenantID: "1",
			UserID:   "1",
			RoleID:   "1",
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
