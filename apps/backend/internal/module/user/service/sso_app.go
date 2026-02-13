package service

import (
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
)

// ssoAppService 应用服务实现
type ssoAppService struct {
	appRepo repository.SsoAppRepository
}

// NewSsoAppService 创建应用服务实例
func NewSsoAppService(appRepo repository.SsoAppRepository) SsoAppService {
	return &ssoAppService{appRepo: appRepo}
}

// CreateApp 创建应用
func (s *ssoAppService) CreateApp(app *model.SsoApp) error {
	return s.appRepo.CreateApp(app)
}

// GetAppByID 根据ID获取应用
func (s *ssoAppService) GetAppByID(id string) (*model.SsoApp, error) {
	return s.appRepo.GetAppByID(id)
}

// GetAppByCode 根据编码获取应用
func (s *ssoAppService) GetAppByCode(code string) (*model.SsoApp, error) {
	return s.appRepo.GetAppByCode(code)
}

// UpdateApp 更新应用
func (s *ssoAppService) UpdateApp(app *model.SsoApp) error {
	return s.appRepo.UpdateApp(app)
}

// DeleteApp 删除应用
func (s *ssoAppService) DeleteApp(id string) error {
	return s.appRepo.DeleteApp(id)
}

// GetAllApps 获取所有应用
func (s *ssoAppService) GetAllApps() ([]model.SsoApp, error) {
	return s.appRepo.GetAllApps()
}
