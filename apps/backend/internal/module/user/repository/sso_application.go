package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoApplicationRepository 应用仓库实现
type ssoApplicationRepository struct {
	db *gorm.DB
}

// NewSsoApplicationRepository 创建应用仓库实例
func NewSsoApplicationRepository(db *gorm.DB) SsoApplicationRepository {
	return &ssoApplicationRepository{db: db}
}

// CreateApplication 创建应用
func (r *ssoApplicationRepository) CreateApplication(app *model.SsoApplication) error {
	return r.db.Create(app).Error
}

// GetApplicationByID 根据ID获取应用
func (r *ssoApplicationRepository) GetApplicationByID(id string) (*model.SsoApplication, error) {
	var app model.SsoApplication
	result := r.db.Where("id = ?", id).First(&app)
	if result.Error != nil {
		return nil, result.Error
	}
	return &app, nil
}

// GetApplicationByCode 根据编码获取应用
func (r *ssoApplicationRepository) GetApplicationByCode(code string) (*model.SsoApplication, error) {
	var app model.SsoApplication
	result := r.db.Where("application_code = ?", code).First(&app)
	if result.Error != nil {
		return nil, result.Error
	}
	return &app, nil
}

// UpdateApplication 更新应用
func (r *ssoApplicationRepository) UpdateApplication(app *model.SsoApplication) error {
	return r.db.Save(app).Error
}

// DeleteApplication 删除应用
func (r *ssoApplicationRepository) DeleteApplication(id string) error {
	return r.db.Model(&model.SsoApplication{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// GetAllApplications 获取所有应用
func (r *ssoApplicationRepository) GetAllApplications() ([]model.SsoApplication, error) {
	var apps []model.SsoApplication
	result := r.db.Find(&apps)
	if result.Error != nil {
		return nil, result.Error
	}
	return apps, nil
}
