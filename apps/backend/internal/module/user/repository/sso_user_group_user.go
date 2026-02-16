package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

type ssoUserGroupUserRepository struct {
	db *gorm.DB
}

func NewSsoUserGroupUserRepository(db *gorm.DB) SsoUserGroupUserRepository {
	return &ssoUserGroupUserRepository{db: db}
}

func (r *ssoUserGroupUserRepository) CreateUserGroupUser(item *model.SsoUserGroupUser) error {
	return r.db.Create(item).Error
}

func (r *ssoUserGroupUserRepository) GetUserGroupUserByID(id string) (*model.SsoUserGroupUser, error) {
	var item model.SsoUserGroupUser
	result := r.db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoUserGroupUserRepository) GetUserGroupUsersByGroupID(groupID string) ([]model.SsoUserGroupUser, error) {
	var items []model.SsoUserGroupUser
	result := r.db.Where("group_id = ? AND is_deleted = false", groupID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoUserGroupUserRepository) GetUserGroupUsersByUserID(userID string) ([]model.SsoUserGroupUser, error) {
	var items []model.SsoUserGroupUser
	result := r.db.Where("user_id = ? AND is_deleted = false", userID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoUserGroupUserRepository) DeleteUserGroupUser(id string) error {
	return r.db.Model(&model.SsoUserGroupUser{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoUserGroupUserRepository) DeleteUserGroupUsersByGroupID(groupID string) error {
	return r.db.Model(&model.SsoUserGroupUser{}).Where("group_id = ?", groupID).Update("is_deleted", true).Error
}

func (r *ssoUserGroupUserRepository) DeleteUserGroupUsersByUserID(userID string) error {
	return r.db.Model(&model.SsoUserGroupUser{}).Where("user_id = ?", userID).Update("is_deleted", true).Error
}
