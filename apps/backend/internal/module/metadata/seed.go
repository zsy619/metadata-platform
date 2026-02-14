package metadata

import (
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/utils"
	"time"

	"gorm.io/gorm"
)

// SeedData 初始化元数据模块种子数据
func SeedData(db *gorm.DB) {
	utils.SugarLogger.Info("Seeding metadata database...")

	sf := utils.NewSnowflake(1, 1)

	now := time.Now()

	// 1. 初始化默认数据连接组 (ParentID = "0")
	defaultConnGroup := model.MdConn{
		ID:           sf.GenerateIDString(),
		TenantID:     "1",
		ParentID:     "0",
		ConnName:     "system",
		ConnKind:     "MySQL",
		ConnHost:     "localhost",
		ConnPort:     3306,
		ConnUser:     "root",
		ConnPassword: "123456",
		ConnDatabase: "metadata_platform",
		IsDeleted:    false,
		CreateBy:     "system",
		CreateAt:     now,
		UpdateBy:     "system",
		UpdateAt:     now,
	}

	// 检查是否已存在
	var count int64
	db.Model(&model.MdConn{}).Where("conn_name = ? AND parent_id = ?", "system", "0").Count(&count)
	if count == 0 {
		if err := db.Create(&defaultConnGroup).Error; err != nil {
			utils.SugarLogger.Errorf("Failed to seed default connection group: %v", err)
		} else {
			utils.SugarLogger.Info("Seeded default connection group")
		}
	}

	// // 2. 初始化默认模型组 (ParentID = "0")
	// defaultModelGroup := model.MdModel{
	// 	ID:        sf.GenerateIDString(),
	// 	TenantID:  "1",
	// 	ParentID:  "0",
	// 	ModelName: "system",
	// 	ModelCode: "system",
	// 	ModelKind: 0, // 0 as folder/group
	// 	IsDeleted: false,
	// 	CreateBy:  "system",
	// 	CreateAt:  now,
	// 	UpdateBy:  "system",
	// 	UpdateAt:  now,
	// }

	// db.Model(&model.MdModel{}).Where("model_code = ?", "system").Count(&count)
	// if count == 0 {
	// 	if err := db.Create(&defaultModelGroup).Error; err != nil {
	// 		utils.SugarLogger.Errorf("Failed to seed default model group: %v", err)
	// 	} else {
	// 		utils.SugarLogger.Info("Seeded default model group")
	// 	}
	// }

	utils.SugarLogger.Info("Metadata database seeding completed")
}
