package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

type ssoRoleGroupRoleRepository struct {
	db *gorm.DB
}

func NewSsoRoleGroupRoleRepository(db *gorm.DB) SsoRoleGroupRoleRepository {
	return &ssoRoleGroupRoleRepository{db: db}
}

func (r *ssoRoleGroupRoleRepository) CreateRoleGroupRole(item *model.SsoRoleGroupRole) error {
	var existing model.SsoRoleGroupRole
	result := r.db.Where("group_id = ? AND role_id = ?", item.GroupID, item.RoleID).First(&existing)
	if result.Error == nil {
		return r.db.Model(&existing).Updates(map[string]interface{}{
			"is_deleted": false,
			"create_by":  item.CreateBy,
		}).Error
	}
	return r.db.Create(item).Error
}

func (r *ssoRoleGroupRoleRepository) GetRoleGroupRoleByID(id string) (*model.SsoRoleGroupRole, error) {
	var item model.SsoRoleGroupRole
	result := r.db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoRoleGroupRoleRepository) GetRoleGroupRolesByGroupID(groupID string) ([]model.SsoRoleGroupRole, error) {
	var items []model.SsoRoleGroupRole
	result := r.db.Where("group_id = ? AND is_deleted = false", groupID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoRoleGroupRoleRepository) GetRoleGroupRolesByRoleID(roleID string) ([]model.SsoRoleGroupRole, error) {
	var items []model.SsoRoleGroupRole
	result := r.db.Where("role_id = ? AND is_deleted = false", roleID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoRoleGroupRoleRepository) DeleteRoleGroupRole(id string) error {
	return r.db.Model(&model.SsoRoleGroupRole{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoRoleGroupRoleRepository) DeleteRoleGroupRolesByGroupID(groupID string) error {
	return r.db.Model(&model.SsoRoleGroupRole{}).Where("group_id = ?", groupID).Update("is_deleted", true).Error
}

func (r *ssoRoleGroupRoleRepository) DeleteRoleGroupRolesByRoleID(roleID string) error {
	return r.db.Model(&model.SsoRoleGroupRole{}).Where("role_id = ?", roleID).Update("is_deleted", true).Error
}
