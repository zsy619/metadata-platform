package service

import (
	"errors"

	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

type ssoOrgKindService struct {
	orgKindRepo repository.SsoOrgKindRepository
}

func NewSsoOrgKindService(orgKindRepo repository.SsoOrgKindRepository) SsoOrgKindService {
	return &ssoOrgKindService{orgKindRepo: orgKindRepo}
}

func (s *ssoOrgKindService) CreateOrgKind(item *model.SsoOrgKind) error {
	existingItem, err := s.orgKindRepo.GetOrgKindByCode(item.KindCode)
	if err == nil && existingItem != nil {
		return errors.New("组织类型编码已存在")
	}

	if item.ParentID != "" {
		_, err := s.orgKindRepo.GetOrgKindByID(item.ParentID)
		if err != nil {
			return errors.New("父组织类型不存在")
		}
	}

	item.ID = utils.GetSnowflake().GenerateIDString()

	if item.Sort == 0 {
		maxSort, err := s.orgKindRepo.GetMaxSort()
		if err == nil {
			item.Sort = maxSort + 1
		}
	}

	return s.orgKindRepo.CreateOrgKind(item)
}

func (s *ssoOrgKindService) GetOrgKindByID(id string) (*model.SsoOrgKind, error) {
	return s.orgKindRepo.GetOrgKindByID(id)
}

func (s *ssoOrgKindService) GetOrgKindByCode(code string) (*model.SsoOrgKind, error) {
	return s.orgKindRepo.GetOrgKindByCode(code)
}

func (s *ssoOrgKindService) UpdateOrgKind(item *model.SsoOrgKind) error {
	return s.orgKindRepo.UpdateOrgKind(item)
}

func (s *ssoOrgKindService) UpdateOrgKindFields(id string, fields map[string]any) error {
	return s.orgKindRepo.UpdateOrgKindFields(id, fields)
}

func (s *ssoOrgKindService) DeleteOrgKind(id string) error {
	hasChildren, err := s.orgKindRepo.HasChildren(id)
	if err != nil {
		return err
	}
	if hasChildren {
		return errors.New("该组织类型下存在子组织类型，无法删除")
	}
	return s.orgKindRepo.DeleteOrgKind(id)
}

func (s *ssoOrgKindService) GetAllOrgKinds() ([]model.SsoOrgKind, error) {
	return s.orgKindRepo.GetAllOrgKinds()
}
