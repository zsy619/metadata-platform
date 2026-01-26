package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoRoleMenuRepository 角色菜单仓库实现
type ssoRoleMenuRepository struct {
	db *gorm.DB
}

// NewSsoRoleMenuRepository 创建角色菜单仓库实例
func NewSsoRoleMenuRepository(db *gorm.DB) SsoRoleMenuRepository {
	return &ssoRoleMenuRepository{db: db}
}

// CreateRoleMenu 创建关联
func (r *ssoRoleMenuRepository) CreateRoleMenu(roleMenu *model.SsoRoleMenu) error {
	return r.db.Create(roleMenu).Error
}

// GetRoleMenuByID 根据ID获取关联
func (r *ssoRoleMenuRepository) GetRoleMenuByID(id string) (*model.SsoRoleMenu, error) {
	var roleMenu model.SsoRoleMenu
	result := r.db.Where("id = ?", id).First(&roleMenu)
	if result.Error != nil {
		return nil, result.Error
	}
	return &roleMenu, nil
}

// GetRoleMenusByRoleID 根据角色ID获取所有关联
func (r *ssoRoleMenuRepository) GetRoleMenusByRoleID(roleID string) ([]model.SsoRoleMenu, error) {
	var roleMenus []model.SsoRoleMenu
	result := r.db.Where("role_id = ?", roleID).Find(&roleMenus)
	if result.Error != nil {
		return nil, result.Error
	}
	return roleMenus, nil
}

// GetRoleMenusByMenuID 根据菜单ID获取所有关联
func (r *ssoRoleMenuRepository) GetRoleMenusByMenuID(menuID string) ([]model.SsoRoleMenu, error) {
	var roleMenus []model.SsoRoleMenu
	result := r.db.Where("menu_id = ?", menuID).Find(&roleMenus)
	if result.Error != nil {
		return nil, result.Error
	}
	return roleMenus, nil
}

// GetAllRoleMenus 获取所有角色菜单关联
func (r *ssoRoleMenuRepository) GetAllRoleMenus() ([]model.SsoRoleMenu, error) {
	var roleMenus []model.SsoRoleMenu
	result := r.db.Find(&roleMenus)
	if result.Error != nil {
		return nil, result.Error
	}
	return roleMenus, nil
}

// DeleteRoleMenu 删除关联
func (r *ssoRoleMenuRepository) DeleteRoleMenu(id string) error {
	return r.db.Delete(&model.SsoRoleMenu{}, "id = ?", id).Error
}

// DeleteRoleMenusByRoleID 根据角色ID删除所有关联
func (r *ssoRoleMenuRepository) DeleteRoleMenusByRoleID(roleID string) error {
	return r.db.Delete(&model.SsoRoleMenu{}, "role_id = ?", roleID).Error
}

// DeleteRoleMenusByMenuID 根据菜单ID删除所有关联
func (r *ssoRoleMenuRepository) DeleteRoleMenusByMenuID(menuID string) error {
	return r.db.Delete(&model.SsoRoleMenu{}, "menu_id = ?", menuID).Error
}
