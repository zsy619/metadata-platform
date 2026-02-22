package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoUserAddressRepository 用户地址仓库实现
type ssoUserAddressRepository struct {
	db *gorm.DB
}

// NewSsoUserAddressRepository 创建用户地址仓库实例
func NewSsoUserAddressRepository(db *gorm.DB) SsoUserAddressRepository {
	return &ssoUserAddressRepository{db: db}
}

// Create 创建地址
func (r *ssoUserAddressRepository) Create(addr *model.SsoUserAddress) error {
	return r.db.Create(addr).Error
}

// GetByID 根据ID获取地址
func (r *ssoUserAddressRepository) GetByID(id string) (*model.SsoUserAddress, error) {
	var addr model.SsoUserAddress
	if err := r.db.Where("id = ?", id).First(&addr).Error; err != nil {
		return nil, err
	}
	return &addr, nil
}

// GetByUserID 根据用户ID获取地址列表
func (r *ssoUserAddressRepository) GetByUserID(userID string) ([]model.SsoUserAddress, error) {
	var list []model.SsoUserAddress
	if err := r.db.Where("user_id = ?", userID).Order("is_default DESC, create_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateFields 更新地址指定字段
func (r *ssoUserAddressRepository) UpdateFields(id string, fields map[string]any) error {
	return r.db.Model(&model.SsoUserAddress{}).Where("id = ?", id).Updates(fields).Error
}

// ClearDefault 清除用户所有默认标记
func (r *ssoUserAddressRepository) ClearDefault(userID string) error {
	return r.db.Model(&model.SsoUserAddress{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
}

// Delete 删除地址
func (r *ssoUserAddressRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.SsoUserAddress{}).Error
}
