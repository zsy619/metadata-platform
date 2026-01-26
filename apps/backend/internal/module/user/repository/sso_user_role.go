package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoUserRoleRepository 用户角色仓库实现
type ssoUserRoleRepository struct {
	db *gorm.DB
}

// NewSsoUserRoleRepository 创建用户角色仓库实例
func NewSsoUserRoleRepository(db *gorm.DB) SsoUserRoleRepository {
	return &ssoUserRoleRepository{db: db}
}

// CreateUserRole 创建用户角色关联
func (r *ssoUserRoleRepository) CreateUserRole(userRole *model.SsoUserRole) error {
	return r.db.Create(userRole).Error
}

// GetUserRoleByID 根据ID获取关联
func (r *ssoUserRoleRepository) GetUserRoleByID(id string) (*model.SsoUserRole, error) {
	var userRole model.SsoUserRole
	result := r.db.Where("id = ?", id).First(&userRole)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userRole, nil
}

// GetUserRolesByUserID 根据用户ID获取所有角色关联
func (r *ssoUserRoleRepository) GetUserRolesByUserID(userID string) ([]model.SsoUserRole, error) {
	var userRoles []model.SsoUserRole
	result := r.db.Where("user_id = ?", userID).Find(&userRoles)
	if result.Error != nil {
		return nil, result.Error
	}
	return userRoles, nil
}

// GetUserRolesByRoleID 根据角色ID获取所有用户关联
func (r *ssoUserRoleRepository) GetUserRolesByRoleID(roleID string) ([]model.SsoUserRole, error) {
	var userRoles []model.SsoUserRole
	result := r.db.Where("role_id = ?", roleID).Find(&userRoles)
	if result.Error != nil {
		return nil, result.Error
	}
	return userRoles, nil
}

// GetAllUserRoles 获取所有用户角色关联
func (r *ssoUserRoleRepository) GetAllUserRoles() ([]model.SsoUserRole, error) {
	var userRoles []model.SsoUserRole
	result := r.db.Find(&userRoles)
	if result.Error != nil {
		return nil, result.Error
	}
	return userRoles, nil
}

// DeleteUserRole 删除关联
func (r *ssoUserRoleRepository) DeleteUserRole(id string) error {
	return r.db.Delete(&model.SsoUserRole{}, "id = ?", id).Error
}

// DeleteUserRolesByUserID 根据用户ID删除其所有角色关联
func (r *ssoUserRoleRepository) DeleteUserRolesByUserID(userID string) error {
	return r.db.Delete(&model.SsoUserRole{}, "user_id = ?", userID).Error
}

// DeleteUserRolesByRoleID 根据角色ID删除其所有用户关联
func (r *ssoUserRoleRepository) DeleteUserRolesByRoleID(roleID string) error {
	return r.db.Delete(&model.SsoUserRole{}, "role_id = ?", roleID).Error
}
