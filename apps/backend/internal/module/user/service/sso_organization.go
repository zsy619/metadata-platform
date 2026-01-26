package service

import (
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
)

// ssoOrganizationService 组织服务实现
type ssoOrganizationService struct {
	orgRepo repository.SsoOrganizationRepository
}

// NewSsoOrganizationService 创建组织服务实例
func NewSsoOrganizationService(orgRepo repository.SsoOrganizationRepository) SsoOrganizationService {
	return &ssoOrganizationService{orgRepo: orgRepo}
}

// CreateOrganization 创建组织
func (s *ssoOrganizationService) CreateOrganization(unit *model.SsoOrganization) error {
	return s.orgRepo.CreateOrganization(unit)
}

// GetOrganizationByID 根据ID获取组织
func (s *ssoOrganizationService) GetOrganizationByID(id string) (*model.SsoOrganization, error) {
	return s.orgRepo.GetOrganizationByID(id)
}

// GetOrganizationByCode 根据编码获取组织
func (s *ssoOrganizationService) GetOrganizationByCode(code string) (*model.SsoOrganization, error) {
	return s.orgRepo.GetOrganizationByCode(code)
}

// UpdateOrganization 更新组织
func (s *ssoOrganizationService) UpdateOrganization(unit *model.SsoOrganization) error {
	return s.orgRepo.UpdateOrganization(unit)
}

// DeleteOrganization 删除组织
func (s *ssoOrganizationService) DeleteOrganization(id string) error {
	return s.orgRepo.DeleteOrganization(id)
}

// GetAllOrganizations 获取所有组织
func (s *ssoOrganizationService) GetAllOrganizations() ([]model.SsoOrganization, error) {
	return s.orgRepo.GetAllOrganizations()
}
