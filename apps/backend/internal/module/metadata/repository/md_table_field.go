package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/metadata/model"
)

// MdTableFieldRepository 数据连接表字段仓库接口
type MdTableFieldRepository interface {
	CreateField(field *model.MdTableField) error
	GetFieldByID(id string) (*model.MdTableField, error)
	GetFieldsByTableID(tableID string) ([]model.MdTableField, error)
	UpdateField(field *model.MdTableField) error
	DeleteField(id string) error
	DeleteFieldsByTableID(tableID string) error
	GetAllFields(tenantID string) ([]model.MdTableField, error)
}

// mdTableFieldRepository 数据连接表字段仓库实现
type mdTableFieldRepository struct {
	db *gorm.DB
}

// NewMdTableFieldRepository 创建数据连接表字段仓库实例
func NewMdTableFieldRepository(db *gorm.DB) MdTableFieldRepository {
	return &mdTableFieldRepository{db: db}
}

// CreateField 创建数据连接表字段
func (r *mdTableFieldRepository) CreateField(field *model.MdTableField) error {
	return r.db.Create(field).Error
}

// GetFieldByID 根据ID获取数据连接表字段
func (r *mdTableFieldRepository) GetFieldByID(id string) (*model.MdTableField, error) {
	var field model.MdTableField
	result := r.db.Where("id = ?", id).First(&field)
	if result.Error != nil {
		return nil, result.Error
	}
	return &field, nil
}

// GetFieldsByTableID 根据表ID获取所有字段
func (r *mdTableFieldRepository) GetFieldsByTableID(tableID string) ([]model.MdTableField, error) {
	var fields []model.MdTableField
	result := r.db.Where("table_id = ?", tableID).Order("sort asc").Find(&fields)
	if result.Error != nil {
		return nil, result.Error
	}
	return fields, nil
}

// UpdateField 更新数据连接表字段
func (r *mdTableFieldRepository) UpdateField(field *model.MdTableField) error {
	return r.db.Save(field).Error
}

// DeleteField 删除数据连接表字段
func (r *mdTableFieldRepository) DeleteField(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.MdTableField{}).Error
}

// DeleteFieldsByTableID 根据表ID删除所有字段
func (r *mdTableFieldRepository) DeleteFieldsByTableID(tableID string) error {
	return r.db.Where("table_id = ?", tableID).Delete(&model.MdTableField{}).Error
}

// GetAllFields 获取所有数据连接表字段
func (r *mdTableFieldRepository) GetAllFields(tenantID string) ([]model.MdTableField, error) {
	var fields []model.MdTableField
	result := r.db.Where("tenant_id = ?", tenantID).Find(&fields)
	if result.Error != nil {
		return nil, result.Error
	}
	return fields, nil
}
