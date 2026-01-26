package service

import (
	"errors"

	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
)

// ssoTenantService 租户服务实现
type ssoTenantService struct {
	tenantRepo repository.SsoTenantRepository
}

// NewSsoTenantService 创建租户服务实例
func NewSsoTenantService(tenantRepo repository.SsoTenantRepository) SsoTenantService {
	return &ssoTenantService{tenantRepo: tenantRepo}
}

// CreateTenant 创建租户
func (s *ssoTenantService) CreateTenant(tenant *model.SsoTenant) error {
	// 检查租户编码是否已存在
	existingTenant, err := s.tenantRepo.GetTenantByCode(tenant.TenantCode)
	if err == nil && existingTenant != nil {
		return errors.New("租户编码已存在")
	}

	// 创建租户
	return s.tenantRepo.CreateTenant(tenant)
}

// GetTenantByID 根据ID获取租户
func (s *ssoTenantService) GetTenantByID(id string) (*model.SsoTenant, error) {
	return s.tenantRepo.GetTenantByID(id)
}

// GetTenantByCode 根据编码获取租户
func (s *ssoTenantService) GetTenantByCode(code string) (*model.SsoTenant, error) {
	return s.tenantRepo.GetTenantByCode(code)
}

// UpdateTenant 更新租户
func (s *ssoTenantService) UpdateTenant(tenant *model.SsoTenant) error {
	// 检查租户是否存在
	existingTenant, err := s.tenantRepo.GetTenantByID(tenant.ID)
	if err != nil {
		return errors.New("租户不存在")
	}

	// 如果租户编码发生变化，检查新编码是否已存在
	if existingTenant.TenantCode != tenant.TenantCode {
		otherTenant, err := s.tenantRepo.GetTenantByCode(tenant.TenantCode)
		if err == nil && otherTenant != nil {
			return errors.New("租户编码已存在")
		}
	}

	// 更新租户
	return s.tenantRepo.UpdateTenant(tenant)
}

// DeleteTenant 删除租户
func (s *ssoTenantService) DeleteTenant(id string) error {
	// 检查租户是否存在
	_, err := s.tenantRepo.GetTenantByID(id)
	if err != nil {
		return errors.New("租户不存在")
	}

	// 删除租户
	return s.tenantRepo.DeleteTenant(id)
}

// GetAllTenants 获取所有租户
func (s *ssoTenantService) GetAllTenants() ([]model.SsoTenant, error) {
	return s.tenantRepo.GetAllTenants()
}
