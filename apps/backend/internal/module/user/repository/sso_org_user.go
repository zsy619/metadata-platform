package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

type ssoOrgUserRepository struct {
	db *gorm.DB
}

func NewSsoOrgUserRepository(db *gorm.DB) SsoOrgUserRepository {
	return &ssoOrgUserRepository{db: db}
}

func (r *ssoOrgUserRepository) CreateOrgUser(item *model.SsoOrgUser) error {
	return r.db.Create(item).Error
}

func (r *ssoOrgUserRepository) GetOrgUserByID(id string) (*model.SsoOrgUser, error) {
	var item model.SsoOrgUser
	result := r.db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoOrgUserRepository) GetOrgUsersByOrgID(orgID string) ([]model.SsoOrgUser, error) {
	var items []model.SsoOrgUser
	result := r.db.Where("org_id = ? AND is_deleted = false", orgID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoOrgUserRepository) GetOrgUsersByUserID(userID string) ([]model.SsoOrgUser, error) {
	var items []model.SsoOrgUser
	result := r.db.Where("user_id = ? AND is_deleted = false", userID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoOrgUserRepository) DeleteOrgUser(id string) error {
	return r.db.Model(&model.SsoOrgUser{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoOrgUserRepository) DeleteOrgUserByOrgID(orgID string) error {
	return r.db.Model(&model.SsoOrgUser{}).Where("org_id = ?", orgID).Update("is_deleted", true).Error
}

func (r *ssoOrgUserRepository) DeleteOrgUsersByUserID(userID string) error {
	return r.db.Model(&model.SsoOrgUser{}).Where("user_id = ?", userID).Update("is_deleted", true).Error
}

func (r *ssoOrgUserRepository) GetAllOrgUsers() ([]model.SsoOrgUser, error) {
	var items []model.SsoOrgUser
	result := r.db.Where("is_deleted = false").Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}
