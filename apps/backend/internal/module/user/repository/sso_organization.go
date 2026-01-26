package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoOrganizationRepository 组织仓库实现
type ssoOrganizationRepository struct {
	db *gorm.DB
}

// NewSsoOrganizationRepository 创建组织仓库实例
func NewSsoOrganizationRepository(db *gorm.DB) SsoOrganizationRepository {
	return &ssoOrganizationRepository{db: db}
}

// CreateOrganization 创建组织
func (r *ssoOrganizationRepository) CreateOrganization(unit *model.SsoOrganization) error {
	return r.db.Create(unit).Error
}

// GetOrganizationByID 根据ID获取组织
func (r *ssoOrganizationRepository) GetOrganizationByID(id string) (*model.SsoOrganization, error) {
	var unit model.SsoOrganization
	result := r.db.Where("id = ?", id).First(&unit)
	if result.Error != nil {
		return nil, result.Error
	}
	return &unit, nil
}

// GetOrganizationByCode 根据编码获取组织
func (r *ssoOrganizationRepository) GetOrganizationByCode(code string) (*model.SsoOrganization, error) {
	var unit model.SsoOrganization
	result := r.db.Where("unit_code = ?", code).First(&unit)
	if result.Error != nil {
		return nil, result.Error
	}
	return &unit, nil
}

// UpdateOrganization 更新组织
func (r *ssoOrganizationRepository) UpdateOrganization(unit *model.SsoOrganization) error {
	return r.db.Save(unit).Error
}

// DeleteOrganization 删除组织
func (r *ssoOrganizationRepository) DeleteOrganization(id string) error {
	return r.db.Model(&model.SsoOrganization{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// GetAllOrganizations 获取所有组织
func (r *ssoOrganizationRepository) GetAllOrganizations() ([]model.SsoOrganization, error) {
	var units []model.SsoOrganization
	result := r.db.Find(&units)
	if result.Error != nil {
		return nil, result.Error
	}
	return units, nil
}
