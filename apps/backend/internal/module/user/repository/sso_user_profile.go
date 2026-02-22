package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoUserProfileRepository 用户档案仓库实现
type ssoUserProfileRepository struct {
	db *gorm.DB
}

// NewSsoUserProfileRepository 创建用户档案仓库实例
func NewSsoUserProfileRepository(db *gorm.DB) SsoUserProfileRepository {
	return &ssoUserProfileRepository{db: db}
}

// GetByUserID 根据用户ID获取档案
func (r *ssoUserProfileRepository) GetByUserID(userID string) (*model.SsoUserProfile, error) {
	var profile model.SsoUserProfile
	result := r.db.Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		return nil, result.Error
	}
	return &profile, nil
}

// Upsert 插入或更新档案
func (r *ssoUserProfileRepository) Upsert(profile *model.SsoUserProfile) error {
	var existing model.SsoUserProfile
	err := r.db.Where("user_id = ?", profile.UserID).First(&existing).Error
	if err != nil {
		// 不存在，插入
		return r.db.Create(profile).Error
	}
	// 已存在，更新
	profile.ID = existing.ID
	return r.db.Model(&existing).Updates(profile).Error
}

// Delete 删除档案
func (r *ssoUserProfileRepository) Delete(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.SsoUserProfile{}).Error
}
