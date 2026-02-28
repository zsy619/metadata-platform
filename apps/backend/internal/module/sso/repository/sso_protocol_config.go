package repository

import (
	"metadata-platform/internal/module/sso/model"

	"gorm.io/gorm"
)

// ssoProtocolConfigRepository 协议配置仓库实现
type ssoProtocolConfigRepository struct {
	db *gorm.DB
}

// NewSsoProtocolConfigRepository 创建协议配置仓库实例
func NewSsoProtocolConfigRepository(db *gorm.DB) SsoProtocolConfigRepository {
	return &ssoProtocolConfigRepository{db: db}
}

// CreateConfig 创建协议配置
func (r *ssoProtocolConfigRepository) CreateConfig(config *model.SsoProtocolConfig) error {
	return r.db.Create(config).Error
}

// GetConfigByID 根据ID获取协议配置
func (r *ssoProtocolConfigRepository) GetConfigByID(id string) (*model.SsoProtocolConfig, error) {
	var config model.SsoProtocolConfig
	result := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&config)
	if result.Error != nil {
		return nil, result.Error
	}
	return &config, nil
}

// GetConfigByProtocolType 根据协议类型获取配置列表
func (r *ssoProtocolConfigRepository) GetConfigByProtocolType(protocolType model.ProtocolType) ([]model.SsoProtocolConfig, error) {
	var configs []model.SsoProtocolConfig
	result := r.db.Where("protocol_type = ? AND is_deleted = ?", protocolType, false).
		Order("sort ASC, create_at DESC").Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}
	return configs, nil
}

// GetEnabledConfigs 获取所有启用的配置
func (r *ssoProtocolConfigRepository) GetEnabledConfigs() ([]model.SsoProtocolConfig, error) {
	var configs []model.SsoProtocolConfig
	result := r.db.Where("is_enabled = ? AND is_deleted = ?", true, false).
		Order("sort ASC, create_at DESC").Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}
	return configs, nil
}

// UpdateConfig 更新协议配置
func (r *ssoProtocolConfigRepository) UpdateConfig(config *model.SsoProtocolConfig) error {
	return r.db.Save(config).Error
}

// UpdateConfigFields 更新协议配置指定字段
func (r *ssoProtocolConfigRepository) UpdateConfigFields(id string, fields map[string]any) error {
	return r.db.Model(&model.SsoProtocolConfig{}).Where("id = ?", id).Updates(fields).Error
}

// DeleteConfig 删除协议配置（逻辑删除）
func (r *ssoProtocolConfigRepository) DeleteConfig(id string) error {
	return r.db.Model(&model.SsoProtocolConfig{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// GetAllConfigs 获取所有协议配置
func (r *ssoProtocolConfigRepository) GetAllConfigs() ([]model.SsoProtocolConfig, error) {
	var configs []model.SsoProtocolConfig
	result := r.db.Where("is_deleted = ?", false).
		Order("sort ASC, create_at DESC").Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}
	return configs, nil
}
