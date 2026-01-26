package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoRoleRepository 角色仓库实现
type ssoRoleRepository struct {
	db *gorm.DB
}

// NewSsoRoleRepository 创建角色仓库实例
func NewSsoRoleRepository(db *gorm.DB) SsoRoleRepository {
	return &ssoRoleRepository{db: db}
}

// CreateRole 创建角色
func (r *ssoRoleRepository) CreateRole(role *model.SsoRole) error {
	return r.db.Create(role).Error
}

// GetRoleByID 根据ID获取角色
func (r *ssoRoleRepository) GetRoleByID(id string) (*model.SsoRole, error) {
	var role model.SsoRole
	result := r.db.Where("id = ?", id).First(&role)
	if result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

// GetRoleByCode 根据编码获取角色
func (r *ssoRoleRepository) GetRoleByCode(code string) (*model.SsoRole, error) {
	var role model.SsoRole
	result := r.db.Where("role_code = ?", code).First(&role)
	if result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

// UpdateRole 更新角色
func (r *ssoRoleRepository) UpdateRole(role *model.SsoRole) error {
	return r.db.Save(role).Error
}

// DeleteRole 删除角色
func (r *ssoRoleRepository) DeleteRole(id string) error {
	return r.db.Model(&model.SsoRole{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// GetAllRoles 获取所有角色
func (r *ssoRoleRepository) GetAllRoles() ([]model.SsoRole, error) {
	var roles []model.SsoRole
	result := r.db.Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}
