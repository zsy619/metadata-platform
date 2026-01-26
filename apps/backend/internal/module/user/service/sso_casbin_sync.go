package service

import (
	"metadata-platform/internal/middleware"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// SsoCasbinSyncService Casbin 同步服务接口
type SsoCasbinSyncService interface {
	SyncAll() error
}

type ssoCasbinSyncService struct {
	userRoleRepo repository.SsoUserRoleRepository
	roleMenuRepo repository.SsoRoleMenuRepository
	roleRepo     repository.SsoRoleRepository
	menuRepo     repository.SsoMenuRepository
}

// NewSsoCasbinSyncService 创建 Casbin 同步服务实例
func NewSsoCasbinSyncService(
	userRoleRepo repository.SsoUserRoleRepository,
	roleMenuRepo repository.SsoRoleMenuRepository,
	roleRepo     repository.SsoRoleRepository,
	menuRepo     repository.SsoMenuRepository,
) SsoCasbinSyncService {
	return &ssoCasbinSyncService{
		userRoleRepo: userRoleRepo,
		roleMenuRepo: roleMenuRepo,
		roleRepo:     roleRepo,
		menuRepo:     menuRepo,
	}
}

// SyncAll 执行全量同步：将业务表映射到 Casbin 规则
func (s *ssoCasbinSyncService) SyncAll() error {
	enforcer := middleware.GetEnforcer()
	if enforcer == nil {
		return nil
	}

	utils.SugarLogger.Info("Starting Casbin full sync...")

	// 1. 清空当前内存和数据库中的规则 (如果使用的是适配器)
	// 注意：gorm-adapter 的同步通常通过 enforcer.SavePolicy() 完成
	// 这里我们直接清理内存并重新加载
	enforcer.ClearPolicy()

	// 2. 同步 用户-角色 映射 (g)
	userRoles, err := s.userRoleRepo.GetAllUserRoles()
	if err != nil {
		return err
	}

	for _, ur := range userRoles {
		role, err := s.roleRepo.GetRoleByID(ur.RoleID)
		if err != nil || role == nil {
			continue
		}
		// g, userID, roleCode
		_, _ = enforcer.AddGroupingPolicy(ur.UserID, role.RoleCode)
	}

	// 3. 同步 角色-权限 映射 (p)
	roleMenus, err := s.roleMenuRepo.GetAllRoleMenus()
	if err != nil {
		return err
	}

	for _, rm := range roleMenus {
		role, err := s.roleRepo.GetRoleByID(rm.RoleID)
		if err != nil || role == nil {
			continue
		}
		menu, err := s.menuRepo.GetMenuByID(rm.MenuID)
		if err != nil || menu == nil || menu.URL == "" {
			continue
		}
		// p, roleCode, path, method
		// Method 为空时默认为 *
		method := menu.Method
		if method == "" {
			method = "*"
		}
		_, _ = enforcer.AddPolicy(role.RoleCode, menu.URL, method)
	}

	// 4. 持久化到数据库
	if err := enforcer.SavePolicy(); err != nil {
		utils.SugarLogger.Errorf("Failed to save Casbin policy: %v", err)
		return err
	}

	utils.SugarLogger.Info("Casbin full sync completed successfully")
	return nil
}
