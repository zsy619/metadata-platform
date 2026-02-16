package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

type ssoOrgMenuRepository struct {
	db *gorm.DB
}

func NewSsoOrgMenuRepository(db *gorm.DB) SsoOrgMenuRepository {
	return &ssoOrgMenuRepository{db: db}
}

func (r *ssoOrgMenuRepository) CreateOrgMenu(item *model.SsoOrgMenu) error {
	return r.db.Create(item).Error
}

func (r *ssoOrgMenuRepository) GetOrgMenuByID(id string) (*model.SsoOrgMenu, error) {
	var item model.SsoOrgMenu
	result := r.db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoOrgMenuRepository) GetOrgMenusByOrgID(orgID string) ([]model.SsoOrgMenu, error) {
	var items []model.SsoOrgMenu
	result := r.db.Where("org_id = ? AND is_deleted = false", orgID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoOrgMenuRepository) GetOrgMenusByMenuID(menuID string) ([]model.SsoOrgMenu, error) {
	var items []model.SsoOrgMenu
	result := r.db.Where("menu_id = ? AND is_deleted = false", menuID).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoOrgMenuRepository) DeleteOrgMenu(id string) error {
	return r.db.Model(&model.SsoOrgMenu{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoOrgMenuRepository) DeleteOrgMenuByOrgID(orgID string) error {
	return r.db.Model(&model.SsoOrgMenu{}).Where("org_id = ?", orgID).Update("is_deleted", true).Error
}

func (r *ssoOrgMenuRepository) GetAllOrgMenus() ([]model.SsoOrgMenu, error) {
	var items []model.SsoOrgMenu
	result := r.db.Where("is_deleted = false").Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}
