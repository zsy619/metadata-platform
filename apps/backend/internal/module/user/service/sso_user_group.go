package service

import (
	"errors"

	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

type ssoUserGroupService struct {
	userGroupRepo       repository.SsoUserGroupRepository
	userGroupUserRepo   repository.SsoUserGroupUserRepository
}

func NewSsoUserGroupService(userGroupRepo repository.SsoUserGroupRepository, userGroupUserRepo repository.SsoUserGroupUserRepository) *ssoUserGroupService {
	return &ssoUserGroupService{
		userGroupRepo:     userGroupRepo,
		userGroupUserRepo: userGroupUserRepo,
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

	return s.userGroupRepo.DeleteUserGroup(id)
}

func (s *ssoUserGroupService) GetAllUserGroups() ([]model.SsoUserGroup, error) {
	return s.userGroupRepo.GetAllUserGroups()
}
