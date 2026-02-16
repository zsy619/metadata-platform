package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

type ssoCasbinRuleRepository struct {
	db *gorm.DB
}

func NewSsoCasbinRuleRepository(db *gorm.DB) SsoCasbinRuleRepository {
	return &ssoCasbinRuleRepository{db: db}
}

func (r *ssoCasbinRuleRepository) CreateCasbinRule(item *model.SsoCasbinRule) error {
	return r.db.Create(item).Error
}

func (r *ssoCasbinRuleRepository) GetCasbinRuleByID(id string) (*model.SsoCasbinRule, error) {
	var item model.SsoCasbinRule
	result := r.db.Where("id = ?", id).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoCasbinRuleRepository) GetCasbinRulesByPType(ptype string) ([]model.SsoCasbinRule, error) {
	var items []model.SsoCasbinRule
	result := r.db.Where("ptype = ?", ptype).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *ssoCasbinRuleRepository) GetCasbinRule(pType, v0, v1 string) (*model.SsoCasbinRule, error) {
	var item model.SsoCasbinRule
	result := r.db.Where("ptype = ? AND v0 = ? AND v1 = ?", pType, v0, v1).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *ssoCasbinRuleRepository) DeleteCasbinRule(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.SsoCasbinRule{}).Error
}

func (r *ssoCasbinRuleRepository) DeleteCasbinRulesByPType(ptype string) error {
	return r.db.Where("ptype = ?", ptype).Delete(&model.SsoCasbinRule{}).Error
}

func (r *ssoCasbinRuleRepository) DeleteCasbinRules(pType, v0, v1 string) error {
	return r.db.Where("ptype = ? AND v0 = ? AND v1 = ?", pType, v0, v1).Delete(&model.SsoCasbinRule{}).Error
}

func (r *ssoCasbinRuleRepository) GetAllCasbinRules() ([]model.SsoCasbinRule, error) {
	var items []model.SsoCasbinRule
	result := r.db.Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}
