package service

import (
	"metadata-platform/internal/middleware"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// SsoCasbinSyncService Casbin 同步服务接口
type SsoCasbinSyncService interface {
	SyncAll() error
	SyncUser(userID string) error

	SyncRole(roleID string) error
	ClearRole(roleCode string) error

	SyncMenu(menuID string) error
	ClearMenu(url, method string) error

	SyncPos(posID string) error
	SyncOrg(orgID string) error
	SyncUserGroup(groupID string) error
	SyncRoleGroup(groupID string) error
}

type ssoCasbinSyncService struct {
	userRoleRepo repository.SsoUserRoleRepository
	roleMenuRepo repository.SsoRoleMenuRepository
	roleRepo     repository.SsoRoleRepository
	menuRepo     repository.SsoMenuRepository

	userGroupRoleRepo repository.SsoUserGroupRoleRepository
	roleGroupRoleRepo repository.SsoRoleGroupRoleRepository
	userGroupUserRepo repository.SsoUserGroupUserRepository
	userRoleGroupRepo repository.SsoUserRoleGroupRepository
	userPosRepo       repository.SsoUserPosRepository
	posRoleRepo       repository.SsoPosRoleRepository
	orgUserRepo       repository.SsoOrgUserRepository
	orgRoleRepo       repository.SsoOrgRoleRepository
	orgMenuRepo       repository.SsoOrgMenuRepository
}

// NewSsoCasbinSyncService 创建 Casbin 同步服务实例
func NewSsoCasbinSyncService(
	userRoleRepo repository.SsoUserRoleRepository,
	roleMenuRepo repository.SsoRoleMenuRepository,
	roleRepo repository.SsoRoleRepository,
	menuRepo repository.SsoMenuRepository,

	userGroupRoleRepo repository.SsoUserGroupRoleRepository,
	roleGroupRoleRepo repository.SsoRoleGroupRoleRepository,
	userGroupUserRepo repository.SsoUserGroupUserRepository,
	userRoleGroupRepo repository.SsoUserRoleGroupRepository,
	userPosRepo repository.SsoUserPosRepository,
	posRoleRepo repository.SsoPosRoleRepository,
	orgUserRepo repository.SsoOrgUserRepository,
	orgRoleRepo repository.SsoOrgRoleRepository,
	orgMenuRepo repository.SsoOrgMenuRepository,
) SsoCasbinSyncService {
	return &ssoCasbinSyncService{
		userRoleRepo:      userRoleRepo,
		roleMenuRepo:      roleMenuRepo,
		roleRepo:          roleRepo,
		menuRepo:          menuRepo,
		userGroupRoleRepo: userGroupRoleRepo,
		roleGroupRoleRepo: roleGroupRoleRepo,
		userGroupUserRepo: userGroupUserRepo,
		userRoleGroupRepo: userRoleGroupRepo,
		userPosRepo:       userPosRepo,
		posRoleRepo:       posRoleRepo,
		orgUserRepo:       orgUserRepo,
		orgRoleRepo:       orgRoleRepo,
		orgMenuRepo:       orgMenuRepo,
	}
}

// SyncAll 执行全量同步：将业务表映射到 Casbin 规则
func (s *ssoCasbinSyncService) SyncAll() error {
	enforcer := middleware.GetEnforcer()
	if enforcer == nil {
		return nil
	}

	utils.SugarLogger.Info("Starting Casbin full sync...")

	// 1. 清空当前内存和数据库中的规则
	enforcer.ClearPolicy()

	// ---- 预加载 Role 和 Menu 以提高同步性能 ----
	roles, err := s.roleRepo.GetAllRoles()
	if err != nil {
		return err
	}
	roleMap := make(map[string]*model.SsoRole)
	for i := range roles {
		roleMap[roles[i].ID] = &roles[i]
	}

	menus, err := s.menuRepo.GetAllMenus()
	if err != nil {
		return err
	}
	menuMap := make(map[string]*model.SsoMenu)
	for i := range menus {
		menuMap[menus[i].ID] = &menus[i]
	}
	// ------------------------------------------

	// 2. 同步基础的 用户-角色 映射 (g)
	if userRoles, err := s.userRoleRepo.GetAllUserRoles(); err == nil {
		for _, ur := range userRoles {
			if role := roleMap[ur.RoleID]; role != nil {
				_, _ = enforcer.AddGroupingPolicy(ur.UserID, role.RoleCode)
			}
		}
	}

	// 3. 同步基础的 角色-权限 映射 (p)
	if roleMenus, err := s.roleMenuRepo.GetAllRoleMenus(); err == nil {
		for _, rm := range roleMenus {
			role := roleMap[rm.RoleID]
			menu := menuMap[rm.MenuID]
			if role != nil && menu != nil && menu.URL != "" {
				method := menu.Method
				if method == "" {
					method = "*"
				}
				_, _ = enforcer.AddPolicy(role.RoleCode, menu.URL, method)
			}
		}
	}

	// 4. 同步多层级层级依赖 (g & p)
	// UserGroup 关联
	if userGroupUsers, err := s.userGroupUserRepo.GetAllUserGroupUsers(); err == nil {
		for _, ugu := range userGroupUsers {
			_, _ = enforcer.AddGroupingPolicy(ugu.UserID, "UG_"+ugu.GroupID)
		}
	}
	if userGroupRoles, err := s.userGroupRoleRepo.GetAllUserGroupRoles(); err == nil {
		for _, ugr := range userGroupRoles {
			if role := roleMap[ugr.RoleID]; role != nil {
				_, _ = enforcer.AddGroupingPolicy("UG_"+ugr.GroupID, role.RoleCode)
			}
		}
	}

	// RoleGroup 关联
	if userRoleGroups, err := s.userRoleGroupRepo.GetAllUserRoleGroups(); err == nil {
		for _, urg := range userRoleGroups {
			_, _ = enforcer.AddGroupingPolicy(urg.UserID, "RG_"+urg.GroupID)
		}
	}
	if roleGroupRoles, err := s.roleGroupRoleRepo.GetAllRoleGroupRoles(); err == nil {
		for _, rgr := range roleGroupRoles {
			if role := roleMap[rgr.RoleID]; role != nil {
				_, _ = enforcer.AddGroupingPolicy("RG_"+rgr.GroupID, role.RoleCode)
			}
		}
	}

	// Pos 关联
	if userPoss, err := s.userPosRepo.GetAllUserPoss(); err == nil {
		for _, up := range userPoss {
			_, _ = enforcer.AddGroupingPolicy(up.UserID, "POS_"+up.PosID)
		}
	}
	if posRoles, err := s.posRoleRepo.GetAllPosRoles(); err == nil {
		for _, pr := range posRoles {
			if role := roleMap[pr.RoleID]; role != nil {
				_, _ = enforcer.AddGroupingPolicy("POS_"+pr.PosID, role.RoleCode)
			}
		}
	}

	// Org 关联
	if orgUsers, err := s.orgUserRepo.GetAllOrgUsers(); err == nil {
		for _, ou := range orgUsers {
			_, _ = enforcer.AddGroupingPolicy(ou.UserID, "ORG_"+ou.OrgID)
		}
	}
	if orgRoles, err := s.orgRoleRepo.GetAllOrgRoles(); err == nil {
		for _, or := range orgRoles {
			if role := roleMap[or.RoleID]; role != nil {
				_, _ = enforcer.AddGroupingPolicy("ORG_"+or.OrgID, role.RoleCode)
			}
		}
	}
	if orgMenus, err := s.orgMenuRepo.GetAllOrgMenus(); err == nil {
		for _, om := range orgMenus {
			menu := menuMap[om.MenuID]
			if menu != nil && menu.URL != "" {
				method := menu.Method
				if method == "" {
					method = "*"
				}
				_, _ = enforcer.AddPolicy("ORG_"+om.OrgID, menu.URL, method)
			}
		}
	}

	// 5. 持久化到数据库
	if err := enforcer.SavePolicy(); err != nil {
		utils.SugarLogger.Errorf("Failed to save Casbin policy: %v", err)
		return err
	}

	utils.SugarLogger.Info("Casbin full sync completed successfully")
	return nil
}

// SyncUser 执行单用户的 Casbin 规则增量同步
func (s *ssoCasbinSyncService) SyncUser(userID string) error {
	enforcer := middleware.GetEnforcer()
	if enforcer == nil {
		return nil
	}

	// 1. 清理该用户在内存中的所有角色关联 (g, userID, ...)
	if _, err := enforcer.RemoveFilteredGroupingPolicy(0, userID); err != nil {
		utils.SugarLogger.Errorf("Failed to clear Casbin policy for user %s: %v", userID, err)
		return err
	}

	// ---- 预加载必要的基础信息 ----
	roles, err := s.roleRepo.GetAllRoles()
	if err != nil {
		return err
	}
	roleMap := make(map[string]*model.SsoRole)
	for i := range roles {
		roleMap[roles[i].ID] = &roles[i]
	}

	// 2. 重新加载该用户的各项关联
	// 直系角色
	if userRoles, err := s.userRoleRepo.GetUserRolesByUserID(userID); err == nil {
		for _, ur := range userRoles {
			if role := roleMap[ur.RoleID]; role != nil {
				_, _ = enforcer.AddGroupingPolicy(userID, role.RoleCode)
			}
		}
	}

	// 用户组
	if userGroupUsers, err := s.userGroupUserRepo.GetUserGroupUsersByUserID(userID); err == nil {
		for _, ugu := range userGroupUsers {
			_, _ = enforcer.AddGroupingPolicy(userID, "UG_"+ugu.GroupID)
		}
	}

	// 角色组
	if userRoleGroups, err := s.userRoleGroupRepo.GetUserRoleGroupsByUserID(userID); err == nil {
		for _, urg := range userRoleGroups {
			_, _ = enforcer.AddGroupingPolicy(userID, "RG_"+urg.GroupID)
		}
	}

	// 职位
	if userPoss, err := s.userPosRepo.GetUserPosByUserID(userID); err == nil {
		for _, up := range userPoss {
			_, _ = enforcer.AddGroupingPolicy(userID, "POS_"+up.PosID)
		}
	}

	// 组织
	if orgUsers, err := s.orgUserRepo.GetOrgUsersByUserID(userID); err == nil {
		for _, ou := range orgUsers {
			_, _ = enforcer.AddGroupingPolicy(userID, "ORG_"+ou.OrgID)
		}
	}

	// 3. 持久化到数据库
	if err := enforcer.SavePolicy(); err != nil {
		utils.SugarLogger.Errorf("Failed to save Casbin policy for user %s: %v", userID, err)
		return err
	}

	utils.SugarLogger.Infof("Casbin sync for user %s completed successfully", userID)
	return nil
}

// SyncRole 执行单角色(涉及菜单权限p)的 Casbin 增量同步
func (s *ssoCasbinSyncService) SyncRole(roleID string) error {
	enforcer := middleware.GetEnforcer()
	if enforcer == nil {
		return nil
	}

	role, err := s.roleRepo.GetRoleByID(roleID)
	if err != nil || role == nil {
		return err
	}

	// 清理该角色绑定的所有旧菜单 p 规则
	if _, err := enforcer.RemoveFilteredPolicy(0, role.RoleCode); err != nil {
		return err
	}

	menus, err := s.menuRepo.GetAllMenus()
	if err != nil {
		return err
	}
	menuMap := make(map[string]*model.SsoMenu)
	for i := range menus {
		menuMap[menus[i].ID] = &menus[i]
	}

	// 重新写入包含的所有权限
	if roleMenus, err := s.roleMenuRepo.GetRoleMenusByRoleID(roleID); err == nil {
		for _, rm := range roleMenus {
			if menu := menuMap[rm.MenuID]; menu != nil && menu.URL != "" {
				method := menu.Method
				if method == "" {
					method = "*"
				}
				_, _ = enforcer.AddPolicy(role.RoleCode, menu.URL, method)
			}
		}
	}
	_ = enforcer.SavePolicy()
	return nil
}

// ClearRole 完全清理被删除角色的有关关联
func (s *ssoCasbinSyncService) ClearRole(roleCode string) error {
	enforcer := middleware.GetEnforcer()
	if enforcer == nil {
		return nil
	}

	// 清除角色的菜单绑定，以及绑定在这个角色的用户或组等分组
	_, _ = enforcer.RemoveFilteredPolicy(0, roleCode)
	_, _ = enforcer.RemoveFilteredGroupingPolicy(1, roleCode)

	_ = enforcer.SavePolicy()
	return nil
}

// SyncMenu 更新单个菜单产生的散列路径变更
func (s *ssoCasbinSyncService) SyncMenu(menuID string) error {
	// （提示：若菜单 URL 修改，通常需要在外部 Service 先执行一次 ClearMenu 清除旧地址路由的规则，再执行这里 SyncMenu 来生成新的！）
	enforcer := middleware.GetEnforcer()
	if enforcer == nil {
		return nil
	}
	menu, err := s.menuRepo.GetMenuByID(menuID)
	if err != nil || menu == nil || menu.URL == "" {
		return err
	}
	method := menu.Method
	if method == "" {
		method = "*"
	}
	roles, err := s.roleRepo.GetAllRoles()
	if err != nil {
		return err
	}
	roleMap := make(map[string]*model.SsoRole)
	for i := range roles {
		roleMap[roles[i].ID] = &roles[i]
	}

	// 清除之前因为这个菜单的变动已经有的旧规则是外部做的，这里只管无脑加入
	// 找到被哪些角色直接绑定
	if roleMenus, err := s.roleMenuRepo.GetAllRoleMenus(); err == nil {
		for _, rm := range roleMenus {
			if rm.MenuID == menuID {
				if role := roleMap[rm.RoleID]; role != nil {
					_, _ = enforcer.AddPolicy(role.RoleCode, menu.URL, method)
				}
			}
		}
	}

	// 找到被哪些组织直连
	if orgMenus, err := s.orgMenuRepo.GetAllOrgMenus(); err == nil {
		for _, om := range orgMenus {
			if om.MenuID == menuID {
				_, _ = enforcer.AddPolicy("ORG_"+om.OrgID, menu.URL, method)
			}
		}
	}

	_ = enforcer.SavePolicy()
	return nil
}

// ClearMenu 清空某个特定路由菜单在全局建立的 P 规则
func (s *ssoCasbinSyncService) ClearMenu(url, method string) error {
	enforcer := middleware.GetEnforcer()
	if enforcer == nil {
		return nil
	}
	if method == "" {
		method = "*"
	}
	_, _ = enforcer.RemoveFilteredPolicy(1, url, method)
	_ = enforcer.SavePolicy()
	return nil
}

// SyncPos 同步单一职位的最新角色配置
func (s *ssoCasbinSyncService) SyncPos(posID string) error {
	enforcer := middleware.GetEnforcer()
	if enforcer == nil {
		return nil
	}

	_, _ = enforcer.RemoveFilteredGroupingPolicy(0, "POS_"+posID)

	roles, err := s.roleRepo.GetAllRoles()
	if err != nil {
		return err
	}
	roleMap := make(map[string]*model.SsoRole)
	for i := range roles {
		roleMap[roles[i].ID] = &roles[i]
	}

	if posRoles, err := s.posRoleRepo.GetPosRolesByPosID(posID); err == nil {
		for _, pr := range posRoles {
			if role := roleMap[pr.RoleID]; role != nil {
				_, _ = enforcer.AddGroupingPolicy("POS_"+posID, role.RoleCode)
			}
		}
	}
	_ = enforcer.SavePolicy()
	return nil
}

// SyncOrg 同步单一组织的最新角色与菜单配置
func (s *ssoCasbinSyncService) SyncOrg(orgID string) error {
	enforcer := middleware.GetEnforcer()
	if enforcer == nil {
		return nil
	}

	_, _ = enforcer.RemoveFilteredGroupingPolicy(0, "ORG_"+orgID)
	_, _ = enforcer.RemoveFilteredPolicy(0, "ORG_"+orgID)

	roles, err := s.roleRepo.GetAllRoles()
	if err != nil {
		return err
	}
	roleMap := make(map[string]*model.SsoRole)
	for i := range roles {
		roleMap[roles[i].ID] = &roles[i]
	}

	menus, err := s.menuRepo.GetAllMenus()
	if err != nil {
		return err
	}
	menuMap := make(map[string]*model.SsoMenu)
	for i := range menus {
		menuMap[menus[i].ID] = &menus[i]
	}

	if orgRoles, err := s.orgRoleRepo.GetOrgRolesByOrgID(orgID); err == nil {
		for _, or := range orgRoles {
			if role := roleMap[or.RoleID]; role != nil {
				_, _ = enforcer.AddGroupingPolicy("ORG_"+orgID, role.RoleCode)
			}
		}
	}

	if orgMenus, err := s.orgMenuRepo.GetOrgMenusByOrgID(orgID); err == nil {
		for _, om := range orgMenus {
			if menu := menuMap[om.MenuID]; menu != nil && menu.URL != "" {
				method := menu.Method
				if method == "" {
					method = "*"
				}
				_, _ = enforcer.AddPolicy("ORG_"+orgID, menu.URL, method)
			}
		}
	}

	_ = enforcer.SavePolicy()
	return nil
}

// SyncUserGroup 同步单一用户组最新角色配置
func (s *ssoCasbinSyncService) SyncUserGroup(groupID string) error {
	enforcer := middleware.GetEnforcer()
	if enforcer == nil {
		return nil
	}

	_, _ = enforcer.RemoveFilteredGroupingPolicy(0, "UG_"+groupID)

	roles, err := s.roleRepo.GetAllRoles()
	if err != nil {
		return err
	}
	roleMap := make(map[string]*model.SsoRole)
	for i := range roles {
		roleMap[roles[i].ID] = &roles[i]
	}

	if groupRoles, err := s.userGroupRoleRepo.GetUserGroupRolesByGroupID(groupID); err == nil {
		for _, gr := range groupRoles {
			if role := roleMap[gr.RoleID]; role != nil {
				_, _ = enforcer.AddGroupingPolicy("UG_"+groupID, role.RoleCode)
			}
		}
	}
	_ = enforcer.SavePolicy()
	return nil
}

// SyncRoleGroup 同步单一角色组最新角色配置
func (s *ssoCasbinSyncService) SyncRoleGroup(groupID string) error {
	enforcer := middleware.GetEnforcer()
	if enforcer == nil {
		return nil
	}

	_, _ = enforcer.RemoveFilteredGroupingPolicy(0, "RG_"+groupID)

	roles, err := s.roleRepo.GetAllRoles()
	if err != nil {
		return err
	}
	roleMap := make(map[string]*model.SsoRole)
	for i := range roles {
		roleMap[roles[i].ID] = &roles[i]
	}

	if groupRoles, err := s.roleGroupRoleRepo.GetRoleGroupRolesByGroupID(groupID); err == nil {
		for _, gr := range groupRoles {
			if role := roleMap[gr.RoleID]; role != nil {
				_, _ = enforcer.AddGroupingPolicy("RG_"+groupID, role.RoleCode)
			}
		}
	}
	_ = enforcer.SavePolicy()
	return nil
}
