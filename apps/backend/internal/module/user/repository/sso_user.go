package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

// ssoUserRepository 用户仓库实现
type ssoUserRepository struct {
	db *gorm.DB
}

// NewSsoUserRepository 创建用户仓库实例
func NewSsoUserRepository(db *gorm.DB) SsoUserRepository {
	return &ssoUserRepository{db: db}
}

// CreateUser 创建用户
func (r *ssoUserRepository) CreateUser(user *model.SsoUser) error {
	return r.db.Create(user).Error
}

// GetUserByID 根据ID获取用户
func (r *ssoUserRepository) GetUserByID(id string) (*model.SsoUser, error) {
	var user model.SsoUser
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByAccount 根据账号获取用户
func (r *ssoUserRepository) GetUserByAccount(account string) (*model.SsoUser, error) {
	var user model.SsoUser
	result := r.db.Where("account = ?", account).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser 更新用户
func (r *ssoUserRepository) UpdateUser(user *model.SsoUser) error {
	return r.db.Save(user).Error
}

// DeleteUser 删除用户
func (r *ssoUserRepository) DeleteUser(id string) error {
	return r.db.Model(&model.SsoUser{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// GetAllUsers 获取所有用户
func (r *ssoUserRepository) GetAllUsers() ([]model.SsoUser, error) {
	var users []model.SsoUser
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
