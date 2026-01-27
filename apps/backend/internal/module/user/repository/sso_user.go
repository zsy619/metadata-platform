package repository

import (
	"time"

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

// GetUserWithDetails 获取用户及其关联详情
func (r *ssoUserRepository) GetUserWithDetails(id string) (*model.SsoUser, error) {
	var user model.SsoUser
	result := r.db.Preload("Roles").Preload("Positions").Preload("Organizations").Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateLoginInfo 更新登录信息
func (r *ssoUserRepository) UpdateLoginInfo(id string, ip string) error {
	now := time.Now()
	// 使用 map 更新，确保零值也能更新
	updates := map[string]interface{}{
		"last_login_time":   now,
		"last_ip":           ip,
		"login_error_count": 0,
		"first_login":       1, // 标记为非首次登录（假设 0 是首次？逻辑可能反了，通常 1 是已登录过，或者 update logic 这里 simply set to 1 indicating "logged in at least once"）
	}
	// 修正：如果 first_login 定义是 0 为未登录，1 为已登录，则这里 update to 1。如果是 count，则应 increment。
	// 根据字段名 FirstLogin (int)，根据 SQL create table default 0。假设 1 表示已经登录过了。
	return r.db.Model(&model.SsoUser{}).Where("id = ?", id).Updates(updates).Error
}

// IncrementLoginError 增加登录错误次数
func (r *ssoUserRepository) IncrementLoginError(id string) error {
	return r.db.Model(&model.SsoUser{}).Where("id = ?", id).Update("login_error_count", gorm.Expr("login_error_count + ?", 1)).Error
}
