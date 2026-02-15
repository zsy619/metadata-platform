package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

// ssoOrgRepository 组织仓库实现
type ssoOrgRepository struct {
	db *gorm.DB
}

// NewSsoOrgRepository 创建组织仓库实例
func NewSsoOrgRepository(db *gorm.DB) SsoOrgRepository {
	return &ssoOrgRepository{db: db}
}

// CreateOrg 创建组织
func (r *ssoOrgRepository) CreateOrg(unit *model.SsoOrg) error {
	return r.db.Create(unit).Error
}

// GetOrgByID 根据ID获取组织
func (r *ssoOrgRepository) GetOrgByID(id string) (*model.SsoOrg, error) {
	var unit model.SsoOrg
	result := r.db.Where("id = ?", id).First(&unit)
	if result.Error != nil {
		return nil, result.Error
	}
	return &unit, nil
}

// GetOrgByCode 根据编码获取组织
func (r *ssoOrgRepository) GetOrgByCode(code string) (*model.SsoOrg, error) {
	var unit model.SsoOrg
	result := r.db.Where("unit_code = ?", code).First(&unit)
	if result.Error != nil {
		return nil, result.Error
	}
	return &unit, nil
}

// UpdateOrg 更新组织
func (r *ssoOrgRepository) UpdateOrg(unit *model.SsoOrg) error {
	return r.db.Save(unit).Error
}

// UpdateOrgFields 更新组织指定字段
func (r *ssoOrgRepository) UpdateOrgFields(id string, fields map[string]any) error {
	return r.db.Model(&model.SsoOrg{}).Where("id = ?", id).Updates(fields).Error
}

// DeleteOrg 删除组织
func (r *ssoOrgRepository) DeleteOrg(id string) error {
	return r.db.Model(&model.SsoOrg{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// GetAllOrgs 获取所有组织
func (r *ssoOrgRepository) GetAllOrgs() ([]model.SsoOrg, error) {
	var units []model.SsoOrg
	result := r.db.Find(&units)
	if result.Error != nil {
		return nil, result.Error
	}
	return units, nil
}
