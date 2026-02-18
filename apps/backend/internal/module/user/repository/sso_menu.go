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

// UpdateMenuFields 更新菜单指定字段
// 使用 map 方式只更新指定的字段，避免全量更新
func (r *ssoMenuRepository) UpdateMenuFields(id string, fields map[string]any) error {
	// 更新菜单字段
	if err := r.db.Model(&model.SsoMenu{}).Where("id = ?", id).Updates(fields).Error; err != nil {
		return err
	}
	// // 更新完毕，计算当前的层级及子节点的层级
	// var menu model.SsoMenu
	// if err := r.db.Where("id = ?", id).First(&menu).Error; err != nil {
	// 	return err
	// }
	// // 递归更新子节点的层级
	// if menu.Tier > 0 {
	// 	if err := r.updateChildLevels(id, menu.Tier+1); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// DeleteMenu 删除菜单（物理删除）
func (r *ssoMenuRepository) DeleteMenu(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.SsoMenu{}).Error
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

// GetMaxSort 获取最大排序值
func (r *ssoMenuRepository) GetMaxSort() (int, error) {
	var maxSort int
	result := r.db.Model(&model.SsoMenu{}).Select("COALESCE(MAX(sort), 0)").Scan(&maxSort)
	if result.Error != nil {
		return 0, result.Error
	}
	return maxSort, nil
}
