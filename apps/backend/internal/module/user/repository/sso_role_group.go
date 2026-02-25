package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

type ssoRoleGroupRepository struct {
	db *gorm.DB
}

func NewSsoRoleGroupRepository(db *gorm.DB) SsoRoleGroupRepository {
	return &ssoRoleGroupRepository{db: db}
}

func (r *ssoRoleGroupRepository) CreateRoleGroup(item *model.SsoRoleGroup) error {
	return r.db.Create(item).Error
}

func (r *ssoRoleGroupRepository) GetRoleGroupByID(id string) (*model.SsoRoleGroup, error) {
	var item model.SsoRoleGroup
	result := r.db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoRoleGroupRepository) GetRoleGroupByCode(code string) (*model.SsoRoleGroup, error) {
	var item model.SsoRoleGroup
	result := r.db.Where("group_code = ?", code).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoRoleGroupRepository) UpdateRoleGroup(item *model.SsoRoleGroup) error {
	return r.db.Save(item).Error
}

func (r *ssoRoleGroupRepository) UpdateRoleGroupFields(id string, fields map[string]any) error {
	return r.db.Model(&model.SsoRoleGroup{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ssoRoleGroupRepository) DeleteRoleGroup(id string) error {
	return r.db.Model(&model.SsoRoleGroup{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoRoleGroupRepository) HasChildren(parentID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.SsoRoleGroup{}).Where("parent_id = ? AND is_deleted = false", parentID).Count(&count).Error
	return count > 0, err
}

func (r *ssoRoleGroupRepository) GetAllRoleGroups() ([]model.SsoRoleGroup, error) {
	var items []model.SsoRoleGroup
	result := r.db.Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoRoleGroupRepository) GetMaxSort() (int, error) {
	var maxSort int
	err := r.db.Model(&model.SsoRoleGroup{}).Select("COALESCE(MAX(sort), 0)").Scan(&maxSort).Error
	return maxSort, err
}

func (r *ssoRoleGroupRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.SsoRoleGroup{}).Where("is_deleted = ?", false).Count(&count).Error
	return count, err
}
