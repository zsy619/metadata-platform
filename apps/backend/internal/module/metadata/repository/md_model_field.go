package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/metadata/model"
)

// MdModelFieldRepository 模型字段定义仓库接口
type MdModelFieldRepository interface {
	CreateField(field *model.MdModelField) error
	GetFieldByID(id string) (*model.MdModelField, error)
	GetFieldsByModelID(modelID string) ([]model.MdModelField, error)
	UpdateField(field *model.MdModelField) error
	DeleteField(id string) error
	DeleteFieldsByModelID(modelID string) error
	GetAllFields(tenantID string) ([]model.MdModelField, error)
}

// mdModelFieldRepository 模型字段定义仓库实现
type mdModelFieldRepository struct {
	db *gorm.DB
}

// NewMdModelFieldRepository 创建模型字段定义仓库实例
func NewMdModelFieldRepository(db *gorm.DB) MdModelFieldRepository {
	return &mdModelFieldRepository{db: db}
}

// CreateField 创建模型字段定义
func (r *mdModelFieldRepository) CreateField(field *model.MdModelField) error {
	return r.db.Create(field).Error
}

// GetFieldByID 根据ID获取模型字段定义
func (r *mdModelFieldRepository) GetFieldByID(id string) (*model.MdModelField, error) {
	var field model.MdModelField
	result := r.db.Where("id = ?", id).First(&field)
	if result.Error != nil {
		return nil, result.Error
	}
	return &field, nil
}

// GetFieldsByModelID 根据模型ID获取所有字段定义
func (r *mdModelFieldRepository) GetFieldsByModelID(modelID string) ([]model.MdModelField, error) {
	var fields []model.MdModelField
	result := r.db.Where("model_id = ?", modelID).Find(&fields)
	if result.Error != nil {
		return nil, result.Error
	}
	return fields, nil
}

// UpdateField 更新模型字段定义
func (r *mdModelFieldRepository) UpdateField(field *model.MdModelField) error {
	return r.db.Save(field).Error
}

// DeleteField 删除模型字段定义
func (r *mdModelFieldRepository) DeleteField(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.MdModelField{}).Error
}

// DeleteFieldsByModelID 根据模型ID删除所有字段定义
func (r *mdModelFieldRepository) DeleteFieldsByModelID(modelID string) error {
	return r.db.Where("model_id = ?", modelID).Delete(&model.MdModelField{}).Error
}

// GetAllFields 获取所有模型字段定义
func (r *mdModelFieldRepository) GetAllFields(tenantID string) ([]model.MdModelField, error) {
	var fields []model.MdModelField
	result := r.db.Where("tenant_id = ?", tenantID).Find(&fields)
	if result.Error != nil {
		return nil, result.Error
	}
	return fields, nil
}
