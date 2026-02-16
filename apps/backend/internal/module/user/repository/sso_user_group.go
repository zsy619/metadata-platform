package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

type ssoUserGroupRepository struct {
	db *gorm.DB
}

func NewSsoUserGroupRepository(db *gorm.DB) SsoUserGroupRepository {
	return &ssoUserGroupRepository{db: db}
}

func (r *ssoUserGroupRepository) CreateUserGroup(item *model.SsoUserGroup) error {
	return r.db.Create(item).Error
}

func (r *ssoUserGroupRepository) GetUserGroupByID(id string) (*model.SsoUserGroup, error) {
	var item model.SsoUserGroup
	result := r.db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoUserGroupRepository) GetUserGroupByCode(code string) (*model.SsoUserGroup, error) {
	var item model.SsoUserGroup
	result := r.db.Where("group_code = ?", code).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoUserGroupRepository) UpdateUserGroup(item *model.SsoUserGroup) error {
	return r.db.Save(item).Error
}

func (r *ssoUserGroupRepository) UpdateUserGroupFields(id string, fields map[string]any) error {
	return r.db.Model(&model.SsoUserGroup{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ssoUserGroupRepository) DeleteUserGroup(id string) error {
	return r.db.Model(&model.SsoUserGroup{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoUserGroupRepository) HasChildren(parentID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.SsoUserGroup{}).Where("parent_id = ? AND is_deleted = false", parentID).Count(&count).Error
	return count > 0, err
}

func (r *ssoUserGroupRepository) GetAllUserGroups() ([]model.SsoUserGroup, error) {
	var items []model.SsoUserGroup
	result := r.db.Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoUserGroupRepository) GetMaxSort() (int, error) {
	var maxSort int
	err := r.db.Model(&model.SsoUserGroup{}).Select("COALESCE(MAX(sort), 0)").Scan(&maxSort).Error
	return maxSort, err
}
