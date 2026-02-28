package repository

import (
	"metadata-platform/internal/module/sso/model"

	"gorm.io/gorm"
)

type ssoFieldMappingRepository struct {
	db *gorm.DB
}

func NewSsoFieldMappingRepository(db *gorm.DB) SsoFieldMappingRepository {
	return &ssoFieldMappingRepository{db: db}
}

func (r *ssoFieldMappingRepository) CreateMapping(mapping *model.SsoFieldMapping) error {
	return r.db.Create(mapping).Error
}

func (r *ssoFieldMappingRepository) GetMappingByID(id string) (*model.SsoFieldMapping, error) {
	var mapping model.SsoFieldMapping
	result := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&mapping)
	if result.Error != nil {
		return nil, result.Error
	}
	return &mapping, nil
}

func (r *ssoFieldMappingRepository) GetMappingsByProtocolConfigID(protocolConfigID string) ([]model.SsoFieldMapping, error) {
	var mappings []model.SsoFieldMapping
	result := r.db.Where("protocol_config_id = ? AND is_deleted = ?", protocolConfigID, false).
		Order("sort ASC, create_at DESC").Find(&mappings)
	if result.Error != nil {
		return nil, result.Error
	}
	return mappings, nil
}

func (r *ssoFieldMappingRepository) GetMappingsByClientID(clientID string) ([]model.SsoFieldMapping, error) {
	var mappings []model.SsoFieldMapping
	result := r.db.Where("client_id = ? AND is_deleted = ?", clientID, false).
		Order("sort ASC, create_at DESC").Find(&mappings)
	if result.Error != nil {
		return nil, result.Error
	}
	return mappings, nil
}

func (r *ssoFieldMappingRepository) GetMappingBySourceAndTarget(protocolConfigID, sourceField, targetField string) (*model.SsoFieldMapping, error) {
	var mapping model.SsoFieldMapping
	result := r.db.Where("protocol_config_id = ? AND source_field = ? AND target_field = ? AND is_deleted = ?",
		protocolConfigID, sourceField, targetField, false).First(&mapping)
	if result.Error != nil {
		return nil, result.Error
	}
	return &mapping, nil
}

func (r *ssoFieldMappingRepository) UpdateMapping(mapping *model.SsoFieldMapping) error {
	return r.db.Save(mapping).Error
}

func (r *ssoFieldMappingRepository) UpdateMappingFields(id string, fields map[string]any) error {
	return r.db.Model(&model.SsoFieldMapping{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ssoFieldMappingRepository) DeleteMapping(id string) error {
	return r.db.Model(&model.SsoFieldMapping{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoFieldMappingRepository) GetAllMappings() ([]model.SsoFieldMapping, error) {
	var mappings []model.SsoFieldMapping
	result := r.db.Where("is_deleted = ?", false).
		Order("sort ASC, create_at DESC").Find(&mappings)
	if result.Error != nil {
		return nil, result.Error
	}
	return mappings, nil
}
