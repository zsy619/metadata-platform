package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoUserContactRepository 用户联系方式仓库实现
type ssoUserContactRepository struct {
	db *gorm.DB
}

// NewSsoUserContactRepository 创建用户联系方式仓库实例
func NewSsoUserContactRepository(db *gorm.DB) SsoUserContactRepository {
	return &ssoUserContactRepository{db: db}
}

// Create 创建联系方式
func (r *ssoUserContactRepository) Create(contact *model.SsoUserContact) error {
	return r.db.Create(contact).Error
}

// GetByID 根据ID获取联系方式
func (r *ssoUserContactRepository) GetByID(id string) (*model.SsoUserContact, error) {
	var contact model.SsoUserContact
	if err := r.db.Where("id = ?", id).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// GetByUserID 根据用户ID获取联系方式列表
func (r *ssoUserContactRepository) GetByUserID(userID string) ([]model.SsoUserContact, error) {
	var list []model.SsoUserContact
	if err := r.db.Where("user_id = ?", userID).Order("type ASC, create_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateFields 更新联系方式指定字段
func (r *ssoUserContactRepository) UpdateFields(id string, fields map[string]any) error {
	return r.db.Model(&model.SsoUserContact{}).Where("id = ?", id).Updates(fields).Error
}

// Delete 删除联系方式
func (r *ssoUserContactRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.SsoUserContact{}).Error
}
