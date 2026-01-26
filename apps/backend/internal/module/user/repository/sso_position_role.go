package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoPositionRoleRepository 职位角色仓库实现
type ssoPositionRoleRepository struct {
	db *gorm.DB
}

// NewSsoPositionRoleRepository 创建职位角色仓库实例
func NewSsoPositionRoleRepository(db *gorm.DB) SsoPositionRoleRepository {
	return &ssoPositionRoleRepository{db: db}
}

// CreatePositionRole 创建职位角色关联
func (r *ssoPositionRoleRepository) CreatePositionRole(posRole *model.SsoPositionRole) error {
	return r.db.Create(posRole).Error
}

// GetPositionRoleByID 根据ID获取关联
func (r *ssoPositionRoleRepository) GetPositionRoleByID(id string) (*model.SsoPositionRole, error) {
	var posRole model.SsoPositionRole
	result := r.db.Where("id = ?", id).First(&posRole)
	if result.Error != nil {
		return nil, result.Error
	}
	return &posRole, nil
}

// GetPositionRolesByPosID 根据职位ID获取所有关联
func (r *ssoPositionRoleRepository) GetPositionRolesByPosID(posID string) ([]model.SsoPositionRole, error) {
	var posRoles []model.SsoPositionRole
	result := r.db.Where("pos_id = ?", posID).Find(&posRoles)
	if result.Error != nil {
		return nil, result.Error
	}
	return posRoles, nil
}

// GetPositionRolesByRoleID 根据角色ID获取所有关联
func (r *ssoPositionRoleRepository) GetPositionRolesByRoleID(roleID string) ([]model.SsoPositionRole, error) {
	var posRoles []model.SsoPositionRole
	result := r.db.Where("role_id = ?", roleID).Find(&posRoles)
	if result.Error != nil {
		return nil, result.Error
	}
	return posRoles, nil
}

// DeletePositionRole 删除关联
func (r *ssoPositionRoleRepository) DeletePositionRole(id string) error {
	return r.db.Delete(&model.SsoPositionRole{}, "id = ?", id).Error
}

// DeletePositionRolesByPosID 根据职位ID删除所有关联
func (r *ssoPositionRoleRepository) DeletePositionRolesByPosID(posID string) error {
	return r.db.Delete(&model.SsoPositionRole{}, "pos_id = ?", posID).Error
}

// DeletePositionRolesByRoleID 根据角色ID删除所有关联
func (r *ssoPositionRoleRepository) DeletePositionRolesByRoleID(roleID string) error {
	return r.db.Delete(&model.SsoPositionRole{}, "role_id = ?", roleID).Error
}
