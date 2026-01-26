package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/metadata/model"
)

// MdModelFieldEnhancementRepository 字段增强规范仓库接口
type MdModelFieldEnhancementRepository interface {
	CreateEnhancement(enh *model.MdModelFieldEnhancement) error
	GetEnhancementByFieldID(fieldID string) (*model.MdModelFieldEnhancement, error)
	GetEnhancementsByModelID(modelID string) ([]model.MdModelFieldEnhancement, error)
	UpdateEnhancement(enh *model.MdModelFieldEnhancement) error
	DeleteEnhancement(id string) error
	BatchUpdateEnhancements(enhancements []model.MdModelFieldEnhancement) error
}

type mdModelFieldEnhancementRepository struct {
	db *gorm.DB
}

// NewMdModelFieldEnhancementRepository 创建字段增强规范仓库实例
func NewMdModelFieldEnhancementRepository(db *gorm.DB) MdModelFieldEnhancementRepository {
	return &mdModelFieldEnhancementRepository{db: db}
}

func (r *mdModelFieldEnhancementRepository) CreateEnhancement(enh *model.MdModelFieldEnhancement) error {
	return r.db.Create(enh).Error
}

func (r *mdModelFieldEnhancementRepository) GetEnhancementByFieldID(fieldID string) (*model.MdModelFieldEnhancement, error) {
	var enh model.MdModelFieldEnhancement
	result := r.db.Where("field_id = ?", fieldID).First(&enh)
	if result.Error != nil {
		return nil, result.Error
	}
	return &enh, nil
}

func (r *mdModelFieldEnhancementRepository) GetEnhancementsByModelID(modelID string) ([]model.MdModelFieldEnhancement, error) {
	var enhs []model.MdModelFieldEnhancement
	result := r.db.Where("model_id = ?", modelID).Order("display_order asc").Find(&enhs)
	if result.Error != nil {
		return nil, result.Error
	}
	return enhs, nil
}

func (r *mdModelFieldEnhancementRepository) UpdateEnhancement(enh *model.MdModelFieldEnhancement) error {
	return r.db.Save(enh).Error
}

func (r *mdModelFieldEnhancementRepository) DeleteEnhancement(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.MdModelFieldEnhancement{}).Error
}

func (r *mdModelFieldEnhancementRepository) BatchUpdateEnhancements(enhs []model.MdModelFieldEnhancement) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, enh := range enhs {
			if err := tx.Save(&enh).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
