package user

import (
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/utils"

	"gorm.io/gorm"
)

// SeedData 初始化用户模块种子数据
func SeedData(db *gorm.DB) {
	utils.SugarLogger.Info("Seeding user database...")

	// sf := utils.NewSnowflake(1, 1)

	// 1. 初始化超管租户
	defaultTenant := model.SsoTenant{
		ID:         "1",
		TenantID:   "0",
		TenantName: "系统管理组",
		TenantCode: "system",
		State:      1,
		CreateBy:   "system",
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
			Remark:          "核心元数据管理系统",
			CreateBy:        "system",
		},
		{
			ID:              "2",
			ApplicationName: "用户管理",
			ApplicationCode: "sso",
			State:           1,
			Remark:          "统一账号认证与权限管理系统",
			CreateBy:        "system",
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

	utils.SugarLogger.Info("User database seeding completed")
}
