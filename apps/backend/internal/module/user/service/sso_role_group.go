package service

import (
	"errors"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

type ssoRoleGroupService struct {
	roleGroupRepo     repository.SsoRoleGroupRepository
	roleGroupRoleRepo repository.SsoRoleGroupRoleRepository
}

func NewSsoRoleGroupService(roleGroupRepo repository.SsoRoleGroupRepository, roleGroupRoleRepo repository.SsoRoleGroupRoleRepository) *ssoRoleGroupService {
	return &ssoRoleGroupService{
		roleGroupRepo:     roleGroupRepo,
		roleGroupRoleRepo: roleGroupRoleRepo,
	}
}

func (s *ssoRoleGroupService) CreateRoleGroup(item *model.SsoRoleGroup) error {
	existingItem, err := s.roleGroupRepo.GetRoleGroupByCode(item.GroupCode)
	if err == nil && existingItem != nil {
		return errors.New("角色组编码已存在")
	}

	if item.ParentID != "" {
		_, err := s.roleGroupRepo.GetRoleGroupByID(item.ParentID)
		if err != nil {
			return errors.New("父角色组不存在")
		}
	}

	item.ID = utils.GetSnowflake().GenerateIDString()

	if item.Sort == 0 {
		maxSort, err := s.roleGroupRepo.GetMaxSort()
		if err == nil {
			item.Sort = maxSort + 1
		}
	}

	return s.roleGroupRepo.CreateRoleGroup(item)
}

func (s *ssoRoleGroupService) GetRoleGroupByID(id string) (*model.SsoRoleGroup, error) {
	return s.roleGroupRepo.GetRoleGroupByID(id)
}

func (s *ssoRoleGroupService) GetRoleGroupByCode(code string) (*model.SsoRoleGroup, error) {
	return s.roleGroupRepo.GetRoleGroupByCode(code)
}

func (s *ssoRoleGroupService) UpdateRoleGroup(item *model.SsoRoleGroup) error {
	return s.roleGroupRepo.UpdateRoleGroup(item)
}

func (s *ssoRoleGroupService) UpdateRoleGroupFields(id string, fields map[string]any) error {
	return s.roleGroupRepo.UpdateRoleGroupFields(id, fields)
}

func (s *ssoRoleGroupService) DeleteRoleGroup(id string) error {
	hasChildren, err := s.roleGroupRepo.HasChildren(id)
	if err != nil {
		return err
	}
	if hasChildren {
		return errors.New("该角色组下存在子角色组，无法删除")
	}

	if err := s.roleGroupRoleRepo.DeleteRoleGroupRolesByGroupID(id); err != nil {
		utils.SugarLogger.Errorw("删除角色组角色关联失败", "groupID", id, "error", err)
	}

	return s.roleGroupRepo.DeleteRoleGroup(id)
}

func (s *ssoRoleGroupService) GetAllRoleGroups() ([]model.SsoRoleGroup, error) {
	return s.roleGroupRepo.GetAllRoleGroups()
}
