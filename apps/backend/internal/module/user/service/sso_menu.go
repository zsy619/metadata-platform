package service

import (
	"errors"

	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// ssoMenuService 菜单服务实现
type ssoMenuService struct {
	menuRepo repository.SsoMenuRepository
}

// NewSsoMenuService 创建菜单服务实例
func NewSsoMenuService(menuRepo repository.SsoMenuRepository) SsoMenuService {
	return &ssoMenuService{menuRepo: menuRepo}
}

// CreateMenu 创建菜单
func (s *ssoMenuService) CreateMenu(menu *model.SsoMenu) error {
	// 检查菜单编码是否已存在
	existingMenu, err := s.menuRepo.GetMenuByCode(menu.MenuCode)
	if err == nil && existingMenu != nil {
		return errors.New("菜单编码已存在")
	}

	// 检查父菜单是否存在（如果有）
	if menu.ParentID != "" {
		_, err := s.menuRepo.GetMenuByID(menu.ParentID)
		if err != nil {
			return errors.New("父菜单不存在")
		}
	}

	// 创建ID
	menu.ID = utils.GetSnowflake().GenerateIDString()

	// 自动获取 Sort 最大值并加1
	if menu.Sort == 0 {
		maxSort, err := s.menuRepo.GetMaxSort()
		if err == nil {
			menu.Sort = maxSort + 1
		}
	}

	// 创建菜单
	return s.menuRepo.CreateMenu(menu)
}

// GetMenuByID 根据ID获取菜单
func (s *ssoMenuService) GetMenuByID(id string) (*model.SsoMenu, error) {
	return s.menuRepo.GetMenuByID(id)
}

// GetMenuByCode 根据编码获取菜单
func (s *ssoMenuService) GetMenuByCode(code string) (*model.SsoMenu, error) {
	return s.menuRepo.GetMenuByCode(code)
}

// UpdateMenu 更新菜单
func (s *ssoMenuService) UpdateMenu(menu *model.SsoMenu) error {
	// 检查菜单是否存在
	existingMenu, err := s.menuRepo.GetMenuByID(menu.ID)
	if err != nil {
		return errors.New("菜单不存在")
	}

	// 如果菜单编码发生变化，检查新编码是否已存在
	if existingMenu.MenuCode != menu.MenuCode {
		anotherMenu, err := s.menuRepo.GetMenuByCode(menu.MenuCode)
		if err == nil && anotherMenu != nil {
			return errors.New("菜单编码已存在")
		}
	}

	// 更新菜单
	return s.menuRepo.UpdateMenu(menu)
}

// UpdateMenuFields 更新菜单指定字段
// 使用 map 方式只更新指定的字段，避免全量更新
// 会检查菜单是否存在，以及如果更新了菜单编码，会检查新编码是否已存在
func (s *ssoMenuService) UpdateMenuFields(id string, fields map[string]any) error {
	// 检查菜单是否存在
	_, err := s.menuRepo.GetMenuByID(id)
	if err != nil {
		return errors.New("菜单不存在")
	}

	// 如果更新了菜单编码，检查新编码是否已存在
	if menuCode, ok := fields["menu_code"]; ok && menuCode != "" {
		anotherMenu, err := s.menuRepo.GetMenuByCode(menuCode.(string))
		if err == nil && anotherMenu != nil && anotherMenu.ID != id {
			return errors.New("菜单编码已存在")
		}
	}

	return s.menuRepo.UpdateMenuFields(id, fields)
}

// DeleteMenu 删除菜单
func (s *ssoMenuService) DeleteMenu(id string) error {
	// 检查菜单是否存在
	menu, err := s.menuRepo.GetMenuByID(id)
	if err != nil {
		return errors.New("菜单不存在")
	}

	// 检查是否为系统内置菜单
	if menu.IsSystem {
		return errors.New("系统内置菜单不允许删除")
	}

	// 删除菜单
	return s.menuRepo.DeleteMenu(id)
}

// GetAllMenus 获取所有菜单
func (s *ssoMenuService) GetAllMenus() ([]model.SsoMenu, error) {
	return s.menuRepo.GetAllMenus()
}
