package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
)

type ssoOrgKindRepository struct {
	db *gorm.DB
}

func NewSsoOrgKindRepository(db *gorm.DB) SsoOrgKindRepository {
	return &ssoOrgKindRepository{db: db}
}

func (r *ssoOrgKindRepository) CreateOrgKind(orgKind *model.SsoOrgKind) error {
	return r.db.Create(orgKind).Error
}

func (r *ssoOrgKindRepository) GetOrgKindByID(id string) (*model.SsoOrgKind, error) {
	var item model.SsoOrgKind
	result := r.db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoOrgKindRepository) GetOrgKindByCode(code string) (*model.SsoOrgKind, error) {
	var item model.SsoOrgKind
	result := r.db.Where("kind_code = ?", code).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoOrgKindRepository) UpdateOrgKind(orgKind *model.SsoOrgKind) error {
	return r.db.Save(orgKind).Error
}

func (r *ssoOrgKindRepository) UpdateOrgKindFields(id string, fields map[string]any) error {
	return r.db.Model(&model.SsoOrgKind{}).Where("id = ?", id).Updates(fields).Error
}

func (r *ssoOrgKindRepository) DeleteOrgKind(id string) error {
	return r.db.Model(&model.SsoOrgKind{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *ssoOrgKindRepository) HasChildren(parentID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.SsoOrgKind{}).Where("parent_id = ? AND is_deleted = false", parentID).Count(&count).Error
	return count > 0, err
}

func (r *ssoOrgKindRepository) GetAllOrgKinds() ([]model.SsoOrgKind, error) {
	var items []model.SsoOrgKind
	result := r.db.Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoOrgKindRepository) GetMaxSort() (int, error) {
	var maxSort int
	err := r.db.Model(&model.SsoOrgKind{}).Select("COALESCE(MAX(sort), 0)").Scan(&maxSort).Error
	return maxSort, err
}
