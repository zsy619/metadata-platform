package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

type ssoOrgRoleRepository struct {
	db *gorm.DB
}

func NewSsoOrgRoleRepository(db *gorm.DB) SsoOrgRoleRepository {
	return &ssoOrgRoleRepository{db: db}
}

func (r *ssoOrgRoleRepository) CreateOrgRole(item *model.SsoOrgRole) error {
	return r.db.Create(item).Error
}

func (r *ssoOrgRoleRepository) GetOrgRoleByID(id string) (*model.SsoOrgRole, error) {
	var item model.SsoOrgRole
	result := r.db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoOrgRoleRepository) GetOrgRolesByOrgID(orgID string) ([]model.SsoOrgRole, error) {
	var items []model.SsoOrgRole
	result := r.db.Where("org_id = ? AND is_deleted = false", orgID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoOrgRoleRepository) GetOrgRolesByRoleID(roleID string) ([]model.SsoOrgRole, error) {
	var items []model.SsoOrgRole
	result := r.db.Where("role_id = ? AND is_deleted = false", roleID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoOrgRoleRepository) DeleteOrgRole(id string) error {
	return r.db.Model(&model.SsoOrgRole{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoOrgRoleRepository) DeleteOrgRoleByOrgID(orgID string) error {
	return r.db.Model(&model.SsoOrgRole{}).Where("org_id = ?", orgID).Update("is_deleted", true).Error
}

func (r *ssoOrgRoleRepository) DeleteOrgRolesByRoleID(roleID string) error {
	return r.db.Model(&model.SsoOrgRole{}).Where("role_id = ?", roleID).Update("is_deleted", true).Error
}

func (r *ssoOrgRoleRepository) GetAllOrgRoles() ([]model.SsoOrgRole, error) {
	var items []model.SsoOrgRole
	result := r.db.Where("is_deleted = false").Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}
