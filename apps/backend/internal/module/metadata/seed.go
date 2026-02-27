package metadata

import (
	"strings"
	"time"

	"gorm.io/gorm"

	"metadata-platform/configs"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/utils"
)

// SeedData 初始化元数据模块种子数据
func SeedData(db *gorm.DB, cfg *configs.Config) {
	utils.SugarLogger.Info("Seeding metadata database...")

	sf := utils.NewSnowflake(1, 1)

	now := time.Now()

	// 初始化三个数据库连接
	seedDatabaseConnections(db, cfg, sf, now)

	utils.SugarLogger.Info("Metadata database seeding completed")
}

// seedDatabaseConnections 初始化数据库连接种子数据
func seedDatabaseConnections(db *gorm.DB, cfg *configs.Config, sf *utils.Snowflake, now time.Time) {
	// 定义三个数据库连接配置
	connections := []struct {
		name string
		cfg  configs.DBConfig
	}{
		{"metadata", cfg.MetadataDB},
		{"user", cfg.UserDB},
		{"audit", cfg.AuditDB},
	}

	for _, conn := range connections {
		seedSingleConnection(db, conn.name, conn.cfg, sf, now)
	}
}

// seedSingleConnection 初始化单个数据库连接
func seedSingleConnection(db *gorm.DB, connName string, dbCfg configs.DBConfig, sf *utils.Snowflake, now time.Time) {
	// 根据数据库类型设置端口
	connPort := dbCfg.Port
	if connPort == 0 {
		if strings.ToLower(dbCfg.Type) == "postgres" {
			connPort = 5432
		} else {
			connPort = 3306
		}
	}

	// 创建连接记录
	connRecord := model.MdConn{
		ID:           sf.GenerateIDString(),
		TenantID:     utils.SystemTenantID,
		ParentID:     "",
		ConnName:     connName,
		ConnKind:     strings.ToUpper(dbCfg.Type),
		ConnHost:     dbCfg.Host,
		ConnPort:     connPort,
		ConnUser:     dbCfg.User,
		ConnPassword: dbCfg.Password,
		ConnDatabase: dbCfg.Name,
		IsDeleted:    false,
		CreateBy:     "system",
		CreateAt:     now,
		UpdateBy:     "system",
		UpdateAt:     now,
	}

	// 检查是否已存在
	var count int64
	db.Model(&model.MdConn{}).Where("conn_name = ? AND parent_id = ?", connName, "").Count(&count)
	if count == 0 {
		if err := db.Create(&connRecord).Error; err != nil {
			utils.SugarLogger.Errorf("Failed to seed %s connection: %v", connName, err)
		} else {
			utils.SugarLogger.Infof("Seeded %s connection: %s@%s:%d/%s (%s)",
				connName, dbCfg.User, dbCfg.Host, connPort, dbCfg.Name, dbCfg.Type)
		}
	} else {
		utils.SugarLogger.Debugf("%s connection already exists, skipping", connName)
	}
}
