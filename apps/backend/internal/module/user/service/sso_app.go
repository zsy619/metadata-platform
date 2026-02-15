package service

import (
	"errors"
	"fmt"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
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
	// 检查应用编码是否已存在
	existingApp, err := s.appRepo.GetAppByCode(app.AppCode)
	if err == nil && existingApp != nil {
		return errors.New("应用编码已存在")
	}

	// 使用全局雪花算法生成ID
	app.ID = utils.GetSnowflake().GenerateIDString()

	// 自动获取 Sort 最大值并加1
	if app.Sort == 0 {
		maxSort, err := s.appRepo.GetMaxSort()
		if err == nil {
			app.Sort = maxSort + 1
		}
	}

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
	// 检查应用是否存在
	existingApp, err := s.appRepo.GetAppByID(app.ID)
	if err != nil {
		return errors.New("应用不存在")
	}

	// 如果应用编码发生变化，检查新编码是否已存在
	if existingApp.AppCode != app.AppCode {
		anotherApp, err := s.appRepo.GetAppByCode(app.AppCode)
		if err == nil && anotherApp != nil {
			return errors.New("应用编码已存在")
		}
	}

	return s.appRepo.UpdateApp(app)
}

// UpdateAppFields 更新应用指定字段
func (s *ssoAppService) UpdateAppFields(id string, fields map[string]any) error {
	// 检查应用是否存在
	_, err := s.appRepo.GetAppByID(id)
	if err != nil {
		return errors.New("应用不存在")
	}

	// 如果更新了应用编码，检查新编码是否已存在
	if appCode, ok := fields["app_code"]; ok && appCode != "" {
		anotherApp, err := s.appRepo.GetAppByCode(appCode.(string))
		if err == nil && anotherApp != nil && anotherApp.ID != id {
			return errors.New("应用编码已存在")
		}
	}

	return s.appRepo.UpdateAppFields(id, fields)
}

// DeleteApp 删除应用
func (s *ssoAppService) DeleteApp(id string) error {
	// 检查应用是否存在
	app, err := s.appRepo.GetAppByID(id)
	if err != nil {
		return errors.New("应用不存在")
	}

	// 检查是否为系统内置应用
	if app.IsSystem {
		return errors.New("系统内置应用不允许删除")
	}

	// 检查是否有子应用
	hasChildren, err := s.appRepo.HasChildren(id)
	if err != nil {
		fmt.Println(fmt.Errorf("检查子应用失败: %w", err))
	}
	if hasChildren {
		return errors.New("该应用下存在子应用，无法删除")
	}

	return s.appRepo.DeleteApp(id)
}

// GetAllApps 获取所有应用
func (s *ssoAppService) GetAllApps() ([]model.SsoApp, error) {
	return s.appRepo.GetAllApps()
}
