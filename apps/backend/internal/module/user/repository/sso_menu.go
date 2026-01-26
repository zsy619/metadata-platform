package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoMenuRepository 菜单仓库实现
type ssoMenuRepository struct {
	db *gorm.DB
}

// NewSsoMenuRepository 创建菜单仓库实例
func NewSsoMenuRepository(db *gorm.DB) SsoMenuRepository {
	return &ssoMenuRepository{db: db}
}

// CreateMenu 创建菜单
func (r *ssoMenuRepository) CreateMenu(menu *model.SsoMenu) error {
	return r.db.Create(menu).Error
}

// GetMenuByID 根据ID获取菜单
func (r *ssoMenuRepository) GetMenuByID(id string) (*model.SsoMenu, error) {
	var menu model.SsoMenu
	result := r.db.Where("id = ?", id).First(&menu)
	if result.Error != nil {
		return nil, result.Error
	}
	return &menu, nil
}

// GetMenuByCode 根据编码获取菜单
func (r *ssoMenuRepository) GetMenuByCode(code string) (*model.SsoMenu, error) {
	var menu model.SsoMenu
	result := r.db.Where("menu_code = ?", code).First(&menu)
	if result.Error != nil {
		return nil, result.Error
	}
	return &menu, nil
}

// UpdateMenu 更新菜单
func (r *ssoMenuRepository) UpdateMenu(menu *model.SsoMenu) error {
	return r.db.Save(menu).Error
}

// DeleteMenu 删除菜单
func (r *ssoMenuRepository) DeleteMenu(id string) error {
	return r.db.Model(&model.SsoMenu{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// GetAllMenus 获取所有菜单
func (r *ssoMenuRepository) GetAllMenus() ([]model.SsoMenu, error) {
	var menus []model.SsoMenu
	result := r.db.Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}
	return menus, nil
}
