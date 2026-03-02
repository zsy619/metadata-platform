package document

import (
	"metadata-platform/internal/module/document/model"
	"metadata-platform/internal/utils"

	"gorm.io/gorm"
)

// MigrateDatabase 执行文档模块数据库迁移
func MigrateDatabase(db *gorm.DB) error {
	utils.SugarLogger.Info("Starting document module database migration...")
	
	// 自动迁移模型
	err := db.AutoMigrate(
		&model.Document{},
		&model.DocumentCategory{},
		&model.DocumentVersion{},
		&model.DocumentFavorite{},
		&model.DocumentReadProgress{},
		&model.DocumentFolder{},
	)
	
	if err != nil {
		utils.SugarLogger.Errorf("Failed to migrate document module: %v", err)
		return err
	}
	
	utils.SugarLogger.Info("Document module database migration completed successfully")
	return nil
}

// SeedDatabase 播种子数据
func SeedDatabase(db *gorm.DB) error {
	utils.SugarLogger.Info("Seeding document module data...")
	
	// 清理路径不正确的文件夹数据（临时修复）
	var folderCount int64
	db.Model(&model.DocumentFolder{}).Where("path = '/' OR path = ''").Count(&folderCount)
	if folderCount > 0 {
		utils.SugarLogger.Warnf("Found %d folders with incorrect paths, deleting...", folderCount)
		db.Where("path = '/' OR path = ''").Delete(&model.DocumentFolder{})
	}
	
	// 检查是否已有分类
	var count int64
	db.Model(&model.DocumentCategory{}).Count(&count)
	if count > 0 {
		utils.SugarLogger.Info("Document categories already exist, skipping seed")
		return nil
	}
	
	// 创建默认分类
	categories := []*model.DocumentCategory{
		{
			ID:          "cat_system_overview",
			Name:        "系统概述",
			Description: "系统整体介绍、核心特性、技术指标等",
			Icon:        "fa-info-circle",
			SortOrder:   1,
			IsEnabled:   true,
		},
		{
			ID:          "cat_technical_arch",
			Name:        "技术架构",
			Description: "系统架构设计、技术栈、组件说明等",
			Icon:        "fa-sitemap",
			SortOrder:   2,
			IsEnabled:   true,
		},
		{
			ID:          "cat_core_features",
			Name:        "核心功能",
			Description: "核心功能实现、协议支持、功能说明等",
			Icon:        "fa-cogs",
			SortOrder:   3,
			IsEnabled:   true,
		},
		{
			ID:          "cat_security",
			Name:        "安全增强",
			Description: "安全特性、证书验证、签名验证等",
			Icon:        "fa-shield-alt",
			SortOrder:   4,
			IsEnabled:   true,
		},
		{
			ID:          "cat_performance",
			Name:        "性能优化",
			Description: "性能优化、缓存策略、基准测试等",
			Icon:        "fa-tachometer-alt",
			SortOrder:   5,
			IsEnabled:   true,
		},
		{
			ID:          "cat_test_reports",
			Name:        "测试报告",
			Description: "单元测试、集成测试、性能测试等",
			Icon:        "fa-vial",
			SortOrder:   6,
			IsEnabled:   true,
		},
		{
			ID:          "cat_deployment",
			Name:        "部署指南",
			Description: "部署步骤、环境要求、配置说明等",
			Icon:        "fa-rocket",
			SortOrder:   7,
			IsEnabled:   true,
		},
		{
			ID:          "cat_best_practices",
			Name:        "最佳实践",
			Description: "使用建议、配置优化、常见问题等",
			Icon:        "fa-lightbulb",
			SortOrder:   8,
			IsEnabled:   true,
		},
	}
	
	for _, cat := range categories {
		if err := db.Create(cat).Error; err != nil {
			utils.SugarLogger.Errorf("Failed to create category %s: %v", cat.Name, err)
			continue
		}
		utils.SugarLogger.Infof("Created document category: %s", cat.Name)
	}
	
	utils.SugarLogger.Info("Document module seed completed successfully")
	return nil
}
