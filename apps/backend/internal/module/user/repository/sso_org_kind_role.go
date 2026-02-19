package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

type ssoOrgKindRoleRepository struct {
	db *gorm.DB
}

func NewSsoOrgKindRoleRepository(db *gorm.DB) SsoOrgKindRoleRepository {
	return &ssoOrgKindRoleRepository{db: db}
}

func (r *ssoOrgKindRoleRepository) CreateOrgKindRole(item *model.SsoOrgKindRole) error {
	return r.db.Create(item).Error
}

func (r *ssoOrgKindRoleRepository) GetOrgKindRoleByID(id string) (*model.SsoOrgKindRole, error) {
	var item model.SsoOrgKindRole
	result := r.db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoOrgKindRoleRepository) GetOrgKindRoleByKindCode(kindCode string) ([]model.SsoOrgKindRole, error) {
	var items []model.SsoOrgKindRole
	result := r.db.Where("kind_code = ? AND is_deleted = false", kindCode).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoOrgKindRoleRepository) GetOrgKindRoleByRoleID(roleID string) ([]model.SsoOrgKindRole, error) {
	var items []model.SsoOrgKindRole
	result := r.db.Where("role_id = ? AND is_deleted = false", roleID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoOrgKindRoleRepository) DeleteOrgKindRole(id string) error {
	return r.db.Model(&model.SsoOrgKindRole{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoOrgKindRoleRepository) DeleteOrgKindRoleByKindCode(kindCode string) error {
	return r.db.Model(&model.SsoOrgKindRole{}).Where("kind_code = ?", kindCode).Update("is_deleted", true).Error
}

func (r *ssoOrgKindRoleRepository) DeleteOrgKindRoleByRoleID(roleID string) error {
	return r.db.Model(&model.SsoOrgKindRole{}).Where("role_id = ?", roleID).Update("is_deleted", true).Error
}

func (r *ssoOrgKindRoleRepository) GetAllOrgKindRoles() ([]model.SsoOrgKindRole, error) {
	var items []model.SsoOrgKindRole
	result := r.db.Where("is_deleted = false").Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}
