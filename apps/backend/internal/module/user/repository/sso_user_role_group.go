package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

type ssoUserRoleGroupRepository struct {
	db *gorm.DB
}

func NewSsoUserRoleGroupRepository(db *gorm.DB) SsoUserRoleGroupRepository {
	return &ssoUserRoleGroupRepository{db: db}
}

func (r *ssoUserRoleGroupRepository) CreateUserRoleGroup(item *model.SsoUserRoleGroup) error {
	return r.db.Create(item).Error
}

func (r *ssoUserRoleGroupRepository) GetUserRoleGroupByID(id string) (*model.SsoUserRoleGroup, error) {
	var item model.SsoUserRoleGroup
	result := r.db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoUserRoleGroupRepository) GetUserRoleGroupsByUserID(userID string) ([]model.SsoUserRoleGroup, error) {
	var items []model.SsoUserRoleGroup
	result := r.db.Where("user_id = ? AND is_deleted = false", userID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoUserRoleGroupRepository) GetUserRoleGroupsByGroupID(groupID string) ([]model.SsoUserRoleGroup, error) {
	var items []model.SsoUserRoleGroup
	result := r.db.Where("group_id = ? AND is_deleted = false", groupID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoUserRoleGroupRepository) DeleteUserRoleGroup(id string) error {
	return r.db.Model(&model.SsoUserRoleGroup{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoUserRoleGroupRepository) DeleteUserRoleGroupsByUserID(userID string) error {
	return r.db.Model(&model.SsoUserRoleGroup{}).Where("user_id = ?", userID).Update("is_deleted", true).Error
}

func (r *ssoUserRoleGroupRepository) DeleteUserRoleGroupsByGroupID(groupID string) error {
	return r.db.Model(&model.SsoUserRoleGroup{}).Where("group_id = ?", groupID).Update("is_deleted", true).Error
}
