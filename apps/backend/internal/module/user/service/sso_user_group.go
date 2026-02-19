package service

import (
	"errors"

	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

type ssoUserGroupService struct {
	userGroupRepo     repository.SsoUserGroupRepository
	userGroupUserRepo repository.SsoUserGroupUserRepository
	userGroupRoleRepo repository.SsoUserGroupRoleRepository
}

func NewSsoUserGroupService(userGroupRepo repository.SsoUserGroupRepository, userGroupUserRepo repository.SsoUserGroupUserRepository, userGroupRoleRepo repository.SsoUserGroupRoleRepository) *ssoUserGroupService {
	return &ssoUserGroupService{
		userGroupRepo:     userGroupRepo,
		userGroupUserRepo: userGroupUserRepo,
		userGroupRoleRepo: userGroupRoleRepo,
	}
}

func (s *ssoUserGroupService) CreateUserGroup(item *model.SsoUserGroup) error {
	existingItem, err := s.userGroupRepo.GetUserGroupByCode(item.GroupCode)
	if err == nil && existingItem != nil {
		return errors.New("用户组编码已存在")
	}

	if item.ParentID != "" {
		_, err := s.userGroupRepo.GetUserGroupByID(item.ParentID)
		if err != nil {
			return errors.New("父用户组不存在")
		}
	}

	item.ID = utils.GetSnowflake().GenerateIDString()

	if item.Sort == 0 {
		maxSort, err := s.userGroupRepo.GetMaxSort()
		if err == nil {
			item.Sort = maxSort + 1
		}
	}

	return s.userGroupRepo.CreateUserGroup(item)
}

func (s *ssoUserGroupService) GetUserGroupByID(id string) (*model.SsoUserGroup, error) {
	return s.userGroupRepo.GetUserGroupByID(id)
}

func (s *ssoUserGroupService) GetUserGroupByCode(code string) (*model.SsoUserGroup, error) {
	return s.userGroupRepo.GetUserGroupByCode(code)
}

func (s *ssoUserGroupService) UpdateUserGroup(item *model.SsoUserGroup) error {
	return s.userGroupRepo.UpdateUserGroup(item)
}

func (s *ssoUserGroupService) UpdateUserGroupFields(id string, fields map[string]any) error {
	return s.userGroupRepo.UpdateUserGroupFields(id, fields)
}

func (s *ssoUserGroupService) DeleteUserGroup(id string) error {
	hasChildren, err := s.userGroupRepo.HasChildren(id)
	if err != nil {
		return err
	}
	if hasChildren {
		return errors.New("该用户组下存在子用户组，无法删除")
	}

	if err := s.userGroupUserRepo.DeleteUserGroupUsersByGroupID(id); err != nil {
		utils.SugarLogger.Errorw("删除用户组用户关联失败", "groupID", id, "error", err)
	}

	if err := s.userGroupRoleRepo.DeleteUserGroupRolesByGroupID(id); err != nil {
		utils.SugarLogger.Errorw("删除用户组角色关联失败", "groupID", id, "error", err)
	}

	return s.userGroupRepo.DeleteUserGroup(id)
}

func (s *ssoUserGroupService) GetAllUserGroups() ([]model.SsoUserGroup, error) {
	return s.userGroupRepo.GetAllUserGroups()
}

// GetUserGroupRoles 获取用户组的角色ID列表
func (s *ssoUserGroupService) GetUserGroupRoles(groupID string) ([]string, error) {
	// 检查用户组是否存在
	_, err := s.userGroupRepo.GetUserGroupByID(groupID)
	if err != nil {
		return nil, errors.New("用户组不存在")
	}

	// 获取用户组角色关联
	groupRoles, err := s.userGroupRoleRepo.GetUserGroupRolesByGroupID(groupID)
	if err != nil {
		return nil, err
	}

	// 提取角色ID列表
	roleIDs := make([]string, 0, len(groupRoles))
	for _, gr := range groupRoles {
		roleIDs = append(roleIDs, gr.RoleID)
	}

	return roleIDs, nil
}

// UpdateUserGroupRoles 更新用户组的角色关联
func (s *ssoUserGroupService) UpdateUserGroupRoles(groupID string, roleIDs []string, createBy string) error {
	// 检查用户组是否存在
	_, err := s.userGroupRepo.GetUserGroupByID(groupID)
	if err != nil {
		return errors.New("用户组不存在")
	}

	// 删除原有的用户组角色关联
	if err := s.userGroupRoleRepo.DeleteUserGroupRolesByGroupID(groupID); err != nil {
		return err
	}

	// 创建新的用户组角色关联
	for _, roleID := range roleIDs {
		groupRole := &model.SsoUserGroupRole{
			ID:       utils.GetSnowflake().GenerateIDString(),
			GroupID:  groupID,
			RoleID:   roleID,
			CreateBy: createBy,
		}
		if err := s.userGroupRoleRepo.CreateUserGroupRole(groupRole); err != nil {
			utils.SugarLogger.Errorw("创建用户组角色关联失败", "groupID", groupID, "roleID", roleID, "error", err)
		}
	}

	return nil
}
