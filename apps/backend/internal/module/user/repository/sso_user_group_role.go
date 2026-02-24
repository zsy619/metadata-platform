package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

type ssoUserGroupRoleRepository struct {
	db *gorm.DB
}

func NewSsoUserGroupRoleRepository(db *gorm.DB) SsoUserGroupRoleRepository {
	return &ssoUserGroupRoleRepository{db: db}
}

func (r *ssoUserGroupRoleRepository) CreateUserGroupRole(item *model.SsoUserGroupRole) error {
	return r.db.Create(item).Error
}

func (r *ssoUserGroupRoleRepository) GetUserGroupRoleByID(id string) (*model.SsoUserGroupRole, error) {
	var item model.SsoUserGroupRole
	result := r.db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoUserGroupRoleRepository) GetUserGroupRolesByGroupID(groupID string) ([]model.SsoUserGroupRole, error) {
	var items []model.SsoUserGroupRole
	result := r.db.Where("group_id = ? AND is_deleted = false", groupID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoUserGroupRoleRepository) GetUserGroupRolesByRoleID(roleID string) ([]model.SsoUserGroupRole, error) {
	var items []model.SsoUserGroupRole
	result := r.db.Where("role_id = ? AND is_deleted = false", roleID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoUserGroupRoleRepository) DeleteUserGroupRole(id string) error {
	return r.db.Model(&model.SsoUserGroupRole{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoUserGroupRoleRepository) DeleteUserGroupRolesByGroupID(groupID string) error {
	return r.db.Model(&model.SsoUserGroupRole{}).Where("group_id = ?", groupID).Update("is_deleted", true).Error
}

func (r *ssoUserGroupRoleRepository) DeleteUserGroupRolesByRoleID(roleID string) error {
	return r.db.Model(&model.SsoUserGroupRole{}).Where("role_id = ?", roleID).Update("is_deleted", true).Error
}

func (r *ssoUserGroupRoleRepository) GetAllUserGroupRoles() ([]model.SsoUserGroupRole, error) {
	var items []model.SsoUserGroupRole
	result := r.db.Where("is_deleted = false").Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}
