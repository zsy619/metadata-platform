package service

import (
	"errors"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// ssoRoleService 角色服务实现
type ssoRoleService struct {
	roleRepo          repository.SsoRoleRepository
	roleMenuRepo      repository.SsoRoleMenuRepository
	userRoleRepo      repository.SsoUserRoleRepository
	posRoleRepo       repository.SsoPosRoleRepository
	orgRoleRepo       repository.SsoOrgRoleRepository
	roleGroupRoleRepo repository.SsoRoleGroupRoleRepository
	userGroupRoleRepo repository.SsoUserGroupRoleRepository
	orgKindRoleRepo   repository.SsoOrgKindRoleRepository
}

// NewSsoRoleService 创建角色服务实例
func NewSsoRoleService(roleRepo repository.SsoRoleRepository, roleMenuRepo repository.SsoRoleMenuRepository, userRoleRepo repository.SsoUserRoleRepository, posRoleRepo repository.SsoPosRoleRepository, orgRoleRepo repository.SsoOrgRoleRepository, roleGroupRoleRepo repository.SsoRoleGroupRoleRepository, userGroupRoleRepo repository.SsoUserGroupRoleRepository, orgKindRoleRepo repository.SsoOrgKindRoleRepository) SsoRoleService {
	return &ssoRoleService{
		roleRepo:          roleRepo,
		roleMenuRepo:      roleMenuRepo,
		userRoleRepo:      userRoleRepo,
		posRoleRepo:       posRoleRepo,
		orgRoleRepo:       orgRoleRepo,
		roleGroupRoleRepo: roleGroupRoleRepo,
		userGroupRoleRepo: userGroupRoleRepo,
		orgKindRoleRepo:   orgKindRoleRepo,
	}
}

// CreateRole 创建角色
func (s *ssoRoleService) CreateRole(role *model.SsoRole) error {
	// 检查角色编码是否已存在
	existingRole, err := s.roleRepo.GetRoleByCode(role.RoleCode)
	if err == nil && existingRole != nil {
		return errors.New("角色编码已存在")
	}

	// 检查父角色是否存在（如果有）
	if role.ParentID != "" {
		_, err := s.roleRepo.GetRoleByID(role.ParentID)
		if err != nil {
			return errors.New("父角色不存在")
		}
	}

	// 使用全局雪花算法生成ID
	role.ID = utils.GetSnowflake().GenerateIDString()

	// 自动获取 Sort 最大值并加1
	if role.Sort == 0 {
		maxSort, err := s.roleRepo.GetMaxSort()
		if err == nil {
			role.Sort = maxSort + 1
		}
	}

	// 创建角色
	return s.roleRepo.CreateRole(role)
}

// GetRoleByID 根据ID获取角色
func (s *ssoRoleService) GetRoleByID(id string) (*model.SsoRole, error) {
	return s.roleRepo.GetRoleByID(id)
}

// GetRoleByCode 根据编码获取角色
func (s *ssoRoleService) GetRoleByCode(code string) (*model.SsoRole, error) {
	return s.roleRepo.GetRoleByCode(code)
}

// UpdateRole 更新角色
func (s *ssoRoleService) UpdateRole(role *model.SsoRole) error {
	// 检查角色是否存在
	existingRole, err := s.roleRepo.GetRoleByID(role.ID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 如果角色编码发生变化，检查新编码是否已存在
	if existingRole.RoleCode != role.RoleCode {
		anotherRole, err := s.roleRepo.GetRoleByCode(role.RoleCode)
		if err == nil && anotherRole != nil {
			return errors.New("角色编码已存在")
		}
	}

	// 更新角色
	return s.roleRepo.UpdateRole(role)
}

// DeleteRole 删除角色
func (s *ssoRoleService) DeleteRole(id string) error {
	// 检查角色是否存在
	role, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 检查是否为系统内置角色
	if role.IsSystem {
		return errors.New("系统内置角色不允许删除")
	}

	// 删除角色关联的菜单
	if err := s.roleMenuRepo.DeleteRoleMenusByRoleID(id); err != nil {
		utils.SugarLogger.Errorw("删除角色菜单关联失败", "roleID", id, "error", err)
	}

	// 删除角色关联的用户
	if err := s.userRoleRepo.DeleteUserRolesByRoleID(id); err != nil {
		utils.SugarLogger.Errorw("删除用户角色关联失败", "roleID", id, "error", err)
	}

	// 删除角色关联的职位
	if err := s.posRoleRepo.DeletePosRolesByRoleID(id); err != nil {
		utils.SugarLogger.Errorw("删除职位角色关联失败", "roleID", id, "error", err)
	}

	// 删除角色关联的组织
	if err := s.orgRoleRepo.DeleteOrgRolesByRoleID(id); err != nil {
		utils.SugarLogger.Errorw("删除组织角色关联失败", "roleID", id, "error", err)
	}

	// 删除角色关联的角色组
	if err := s.roleGroupRoleRepo.DeleteRoleGroupRolesByRoleID(id); err != nil {
		utils.SugarLogger.Errorw("删除角色组角色关联失败", "roleID", id, "error", err)
	}

	// 删除角色关联的用户组
	if err := s.userGroupRoleRepo.DeleteUserGroupRolesByRoleID(id); err != nil {
		utils.SugarLogger.Errorw("删除用户组角色关联失败", "roleID", id, "error", err)
	}

	// 删除角色关联的组织类型
	if err := s.orgKindRoleRepo.DeleteOrgKindRoleByRoleID(id); err != nil {
		utils.SugarLogger.Errorw("删除组织类型角色关联失败", "roleID", id, "error", err)
	}

	// 删除角色
	return s.roleRepo.DeleteRole(id)
}

// GetAllRoles 获取所有角色
func (s *ssoRoleService) GetAllRoles() ([]model.SsoRole, error) {
	return s.roleRepo.GetAllRoles()
}

// HasChildren 检查角色是否有子角色
func (s *ssoRoleService) HasChildren(parentID string) (bool, error) {
	return s.roleRepo.HasChildren(parentID)
}

// GetRoleMenus 获取角色的菜单ID列表
func (s *ssoRoleService) GetRoleMenus(roleID string) ([]string, error) {
	// 检查角色是否存在
	_, err := s.roleRepo.GetRoleByID(roleID)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	// 获取角色菜单关联
	roleMenus, err := s.roleMenuRepo.GetRoleMenusByRoleID(roleID)
	if err != nil {
		return nil, err
	}

	// 提取菜单ID列表
	menuIDs := make([]string, 0, len(roleMenus))
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	return menuIDs, nil
}

// UpdateRoleMenus 更新角色的菜单关联
func (s *ssoRoleService) UpdateRoleMenus(roleID string, menuIDs []string, createBy string) error {
	// 检查角色是否存在
	_, err := s.roleRepo.GetRoleByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 删除原有的角色菜单关联
	if err := s.roleMenuRepo.DeleteRoleMenusByRoleID(roleID); err != nil {
		return err
	}

	// 创建新的角色菜单关联
	for _, menuID := range menuIDs {
		roleMenu := &model.SsoRoleMenu{
			ID:       utils.GetSnowflake().GenerateIDString(),
			RoleID:   roleID,
			MenuID:   menuID,
			CreateBy: createBy,
		}
		if err := s.roleMenuRepo.CreateRoleMenu(roleMenu); err != nil {
			utils.SugarLogger.Errorw("创建角色菜单关联失败", "roleID", roleID, "menuID", menuID, "error", err)
		}
	}

	return nil
}
