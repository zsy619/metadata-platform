package service

import (
	"errors"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
)

// ssoOrgService 组织服务实现
type ssoOrgService struct {
	orgRepo repository.SsoOrgRepository
}

// NewSsoOrgService 创建组织服务实例
func NewSsoOrgService(orgRepo repository.SsoOrgRepository) SsoOrgService {
	return &ssoOrgService{orgRepo: orgRepo}
}

// CreateOrg 创建组织
func (s *ssoOrgService) CreateOrg(unit *model.SsoOrg) error {
	// 检查组织编码是否已存在
	existingUnit, err := s.orgRepo.GetOrgByCode(unit.OrgCode)
	if err == nil && existingUnit != nil {
		return errors.New("组织编码已存在")
	}

	// 检查父组织是否存在（如果有）
	if unit.ParentID != "" {
		_, err := s.orgRepo.GetOrgByID(unit.ParentID)
		if err != nil {
			return errors.New("父组织不存在")
		}
	}

	// 使用全局雪花算法生成ID
	unit.ID = utils.GetSnowflake().GenerateIDString()

	// 自动获取 Sort 最大值并加1
	if unit.Sort == 0 {
		maxSort, err := s.orgRepo.GetMaxSort()
		if err == nil {
			unit.Sort = maxSort + 1
		}
	}

	// 创建组织
	return s.orgRepo.CreateOrg(unit)
}

// GetOrgByID 根据ID获取组织
func (s *ssoOrgService) GetOrgByID(id string) (*model.SsoOrg, error) {
	return s.orgRepo.GetOrgByID(id)
}

// GetOrgByCode 根据编码获取组织
func (s *ssoOrgService) GetOrgByCode(code string) (*model.SsoOrg, error) {
	return s.orgRepo.GetOrgByCode(code)
}

// UpdateOrg 更新组织
func (s *ssoOrgService) UpdateOrg(org *model.SsoOrg) error {
	return s.orgRepo.UpdateOrg(org)
}

// DeleteOrg 删除组织
func (s *ssoOrgService) DeleteOrg(id string) error {
	// 检查组织是否存在
	org, err := s.orgRepo.GetOrgByID(id)
	if err != nil {
		return errors.New("组织不存在")
	}

	// 检查是否为系统内置组织
	if org.IsSystem {
		return errors.New("系统内置组织不允许删除")
	}

	return s.orgRepo.DeleteOrg(id)
}

// GetAllOrgs 获取所有组织
func (s *ssoOrgService) GetAllOrgs() ([]model.SsoOrg, error) {
	return s.orgRepo.GetAllOrgs()
}
