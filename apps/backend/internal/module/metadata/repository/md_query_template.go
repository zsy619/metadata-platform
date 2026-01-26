package repository

import (
	"metadata-platform/internal/module/metadata/model"

	"gorm.io/gorm"
)

// MdQueryTemplateRepository 查询模板仓储接口
type MdQueryTemplateRepository interface {
	CreateTemplate(template *model.MdQueryTemplate) error
	GetTemplateByID(id string) (*model.MdQueryTemplate, error)
	GetTemplatesByModelID(modelID string) ([]model.MdQueryTemplate, error)
	UpdateTemplate(template *model.MdQueryTemplate) error
	DeleteTemplate(id string) error
	SetDefault(modelID, templateID string) error
}

// MdQueryConditionRepository 查询条件仓储接口
type MdQueryConditionRepository interface {
	BatchCreateConditions(conditions []model.MdQueryCondition) error
	DeleteConditionsByTemplateID(templateID string) error
	GetConditionsByTemplateID(templateID string) ([]model.MdQueryCondition, error)
}

type mdQueryTemplateRepository struct {
	db *gorm.DB
}

type mdQueryConditionRepository struct {
	db *gorm.DB
}

func NewMdQueryTemplateRepository(db *gorm.DB) MdQueryTemplateRepository {
	return &mdQueryTemplateRepository{db: db}
}

func NewMdQueryConditionRepository(db *gorm.DB) MdQueryConditionRepository {
	return &mdQueryConditionRepository{db: db}
}

// Implement MdQueryTemplateRepository methods
func (r *mdQueryTemplateRepository) CreateTemplate(template *model.MdQueryTemplate) error {
	return r.db.Create(template).Error
}

func (r *mdQueryTemplateRepository) GetTemplateByID(id string) (*model.MdQueryTemplate, error) {
	var template model.MdQueryTemplate
	err := r.db.Preload("Conditions").Where("id = ? AND is_deleted = ?", id, false).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

func (r *mdQueryTemplateRepository) GetTemplatesByModelID(modelID string) ([]model.MdQueryTemplate, error) {
	var templates []model.MdQueryTemplate
	err := r.db.Where("model_id = ? AND is_deleted = ?", modelID, false).Find(&templates).Error
	return templates, err
}

func (r *mdQueryTemplateRepository) UpdateTemplate(template *model.MdQueryTemplate) error {
	return r.db.Save(template).Error
}

func (r *mdQueryTemplateRepository) DeleteTemplate(id string) error {
	return r.db.Model(&model.MdQueryTemplate{}).Where("id = ?", id).Update("is_deleted", true).Error
}

func (r *mdQueryTemplateRepository) SetDefault(modelID, templateID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 取消旧的默认
		if err := tx.Model(&model.MdQueryTemplate{}).Where("model_id = ?", modelID).Update("is_default", false).Error; err != nil {
			return err
		}
		// 设置新的默认
		return tx.Model(&model.MdQueryTemplate{}).Where("id = ?", templateID).Update("is_default", true).Error
	})
}

// Implement MdQueryConditionRepository methods
func (r *mdQueryConditionRepository) BatchCreateConditions(conditions []model.MdQueryCondition) error {
	if len(conditions) == 0 {
		return nil
	}
	return r.db.Create(&conditions).Error
}

func (r *mdQueryConditionRepository) DeleteConditionsByTemplateID(templateID string) error {
	return r.db.Where("template_id = ?", templateID).Delete(&model.MdQueryCondition{}).Error
}

func (r *mdQueryConditionRepository) GetConditionsByTemplateID(templateID string) ([]model.MdQueryCondition, error) {
	var conditions []model.MdQueryCondition
	err := r.db.Where("template_id = ? AND is_deleted = ?", templateID, false).Order("sort asc").Find(&conditions).Error
	return conditions, err
}
