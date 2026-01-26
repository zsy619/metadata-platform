package service

import (
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
)

// ssoApplicationService 应用服务实现
type ssoApplicationService struct {
	appRepo repository.SsoApplicationRepository
}

// NewSsoApplicationService 创建应用服务实例
func NewSsoApplicationService(appRepo repository.SsoApplicationRepository) SsoApplicationService {
	return &ssoApplicationService{appRepo: appRepo}
}

// CreateApplication 创建应用
func (s *ssoApplicationService) CreateApplication(app *model.SsoApplication) error {
	return s.appRepo.CreateApplication(app)
}

// GetApplicationByID 根据ID获取应用
func (s *ssoApplicationService) GetApplicationByID(id string) (*model.SsoApplication, error) {
	return s.appRepo.GetApplicationByID(id)
}

// GetApplicationByCode 根据编码获取应用
func (s *ssoApplicationService) GetApplicationByCode(code string) (*model.SsoApplication, error) {
	return s.appRepo.GetApplicationByCode(code)
}

// UpdateApplication 更新应用
func (s *ssoApplicationService) UpdateApplication(app *model.SsoApplication) error {
	return s.appRepo.UpdateApplication(app)
}

// DeleteApplication 删除应用
func (s *ssoApplicationService) DeleteApplication(id string) error {
	return s.appRepo.DeleteApplication(id)
}

// GetAllApplications 获取所有应用
func (s *ssoApplicationService) GetAllApplications() ([]model.SsoApplication, error) {
	return s.appRepo.GetAllApplications()
}
