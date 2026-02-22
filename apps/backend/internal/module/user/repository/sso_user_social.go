package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoUserSocialRepository 用户第三方账号仓库实现
type ssoUserSocialRepository struct {
	db *gorm.DB
}

// NewSsoUserSocialRepository 创建用户第三方账号仓库实例
func NewSsoUserSocialRepository(db *gorm.DB) SsoUserSocialRepository {
	return &ssoUserSocialRepository{db: db}
}

// Create 创建第三方账号绑定
func (r *ssoUserSocialRepository) Create(social *model.SsoUserSocial) error {
	return r.db.Create(social).Error
}

// GetByID 根据ID获取
func (r *ssoUserSocialRepository) GetByID(id string) (*model.SsoUserSocial, error) {
	var social model.SsoUserSocial
	if err := r.db.Where("id = ?", id).First(&social).Error; err != nil {
		return nil, err
	}
	return &social, nil
}

// GetByUserID 根据用户ID获取绑定列表
func (r *ssoUserSocialRepository) GetByUserID(userID string) ([]model.SsoUserSocial, error) {
	var list []model.SsoUserSocial
	if err := r.db.Where("user_id = ?", userID).Order("provider ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetByProviderAndOpenID 根据平台和 open_id 查找
func (r *ssoUserSocialRepository) GetByProviderAndOpenID(provider, openID string) (*model.SsoUserSocial, error) {
	var social model.SsoUserSocial
	if err := r.db.Where("provider = ? AND open_id = ?", provider, openID).First(&social).Error; err != nil {
		return nil, err
	}
	return &social, nil
}

// Delete 删除绑定记录
func (r *ssoUserSocialRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.SsoUserSocial{}).Error
}
