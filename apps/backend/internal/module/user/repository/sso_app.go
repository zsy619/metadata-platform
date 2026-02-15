package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

// ssoAppRepository 应用仓库实现
type ssoAppRepository struct {
	db *gorm.DB
}

// NewSsoAppRepository 创建应用仓库实例
func NewSsoAppRepository(db *gorm.DB) SsoAppRepository {
	return &ssoAppRepository{db: db}
}

// CreateApp 创建应用
func (r *ssoAppRepository) CreateApp(app *model.SsoApp) error {
	return r.db.Create(app).Error
}

// GetAppByID 根据ID获取应用
func (r *ssoAppRepository) GetAppByID(id string) (*model.SsoApp, error) {
	var app model.SsoApp
	result := r.db.Where("id = ?", id).First(&app)
	if result.Error != nil {
		return nil, result.Error
	}
	return &app, nil
}

// GetAppByCode 根据编码获取应用
func (r *ssoAppRepository) GetAppByCode(code string) (*model.SsoApp, error) {
	var app model.SsoApp
	result := r.db.Where("app_code = ?", code).First(&app)
	if result.Error != nil {
		return nil, result.Error
	}
	return &app, nil
}

// UpdateApp 更新应用
func (r *ssoAppRepository) UpdateApp(app *model.SsoApp) error {
	return r.db.Save(app).Error
}

// UpdateAppFields 更新应用指定字段
func (r *ssoAppRepository) UpdateAppFields(id string, fields map[string]any) error {
	return r.db.Model(&model.SsoApp{}).Where("id = ?", id).Updates(fields).Error
}

// DeleteApp 删除应用（物理删除）
func (r *ssoAppRepository) DeleteApp(id string) error {
	return r.db.Unscoped().Delete(&model.SsoApp{}, "id = ?", id).Error
}

// GetAllApps 获取所有应用
func (r *ssoAppRepository) GetAllApps() ([]model.SsoApp, error) {
	var apps []model.SsoApp
	result := r.db.Find(&apps)
	if result.Error != nil {
		return nil, result.Error
	}
	return apps, nil
}

// GetMaxSort 获取最大排序值
func (r *ssoAppRepository) GetMaxSort() (int, error) {
	var maxSort int
	result := r.db.Model(&model.SsoApp{}).Select("COALESCE(MAX(sort), 0)").Scan(&maxSort)
	if result.Error != nil {
		return 0, result.Error
	}
	return maxSort, nil
}

// HasChildren 检查是否有子应用
func (r *ssoAppRepository) HasChildren(parentID string) (bool, error) {
	var count int64
	result := r.db.Model(&model.SsoApp{}).Where("parent_id = ?", parentID).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
