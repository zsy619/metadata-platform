package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

// ssoPositionRoleRepository 职位角色仓库实现
type ssoPosRoleRepository struct {
	db *gorm.DB
}

// NewSsoPosRoleRepository 创建职位角色仓库实例
func NewSsoPosRoleRepository(db *gorm.DB) SsoPosRoleRepository {
	return &ssoPosRoleRepository{db: db}
}

// CreatePosRole 创建职位角色关联
func (r *ssoPosRoleRepository) CreatePosRole(posRole *model.SsoPosRole) error {
	return r.db.Create(posRole).Error
}

// GetPosRoleByID 根据ID获取关联
func (r *ssoPosRoleRepository) GetPosRoleByID(id string) (*model.SsoPosRole, error) {
	var posRole model.SsoPosRole
	result := r.db.Where("id = ?", id).First(&posRole)
	if result.Error != nil {
		return nil, result.Error
	}
	return &posRole, nil
}

// GetPosRolesByPosID 根据职位ID获取所有关联
func (r *ssoPosRoleRepository) GetPosRolesByPosID(posID string) ([]model.SsoPosRole, error) {
	var posRoles []model.SsoPosRole
	result := r.db.Where("pos_id = ?", posID).Find(&posRoles)
	if result.Error != nil {
		return nil, result.Error
	}
	return posRoles, nil
}

// GetPosRolesByRoleID 根据角色ID获取所有关联
func (r *ssoPosRoleRepository) GetPosRolesByRoleID(roleID string) ([]model.SsoPosRole, error) {
	var posRoles []model.SsoPosRole
	result := r.db.Where("role_id = ?", roleID).Find(&posRoles)
	if result.Error != nil {
		return nil, result.Error
	}
	return posRoles, nil
}

// DeletePosRole 删除关联
func (r *ssoPosRoleRepository) DeletePosRole(id string) error {
	return r.db.Delete(&model.SsoPosRole{}, "id = ?", id).Error
}

// DeletePosRolesByPosID 根据职位ID删除所有关联
func (r *ssoPosRoleRepository) DeletePosRolesByPosID(posID string) error {
	return r.db.Delete(&model.SsoPosRole{}, "pos_id = ?", posID).Error
}

// DeletePosRolesByRoleID 根据角色ID删除所有关联
func (r *ssoPosRoleRepository) DeletePosRolesByRoleID(roleID string) error {
	return r.db.Delete(&model.SsoPosRole{}, "role_id = ?", roleID).Error
}
