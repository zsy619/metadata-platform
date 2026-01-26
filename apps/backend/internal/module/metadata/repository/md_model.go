package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/metadata/model"
)

// MdModelRepository 模型定义仓库接口
type MdModelRepository interface {
	CreateModel(model *model.MdModel) error
	GetModelByID(id string) (*model.MdModel, error)
	GetModelByCode(code string) (*model.MdModel, error)
	UpdateModel(model *model.MdModel) error
	DeleteModel(id string) error
	GetModelsByConnID(connID string) ([]model.MdModel, error)
	GetAllModels(tenantID string) ([]model.MdModel, error)
}

// mdModelRepository 模型定义仓库实现
type mdModelRepository struct {
	db *gorm.DB
}

// NewMdModelRepository 创建模型定义仓库实例
func NewMdModelRepository(db *gorm.DB) MdModelRepository {
	return &mdModelRepository{db: db}
}

// CreateModel 创建模型定义
func (r *mdModelRepository) CreateModel(model *model.MdModel) error {
	return r.db.Create(model).Error
}

// GetModelByID 根据ID获取模型定义
func (r *mdModelRepository) GetModelByID(id string) (*model.MdModel, error) {
	var model model.MdModel
	result := r.db.Where("id = ?", id).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return &model, nil
}

// GetModelByCode 根据编码获取模型定义
func (r *mdModelRepository) GetModelByCode(code string) (*model.MdModel, error) {
	var model model.MdModel
	result := r.db.Where("model_code = ?", code).First(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return &model, nil
}

// UpdateModel 更新模型定义
func (r *mdModelRepository) UpdateModel(model *model.MdModel) error {
	return r.db.Save(model).Error
}

// DeleteModel 删除模型定义
func (r *mdModelRepository) DeleteModel(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.MdModel{}).Error
}

// GetModelsByConnID 根据连接ID获取模型定义列表
func (r *mdModelRepository) GetModelsByConnID(connID string) ([]model.MdModel, error) {
	var models []model.MdModel
	result := r.db.Where("conn_id = ?", connID).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}
	return models, nil
}

// GetAllModels 获取所有模型定义
func (r *mdModelRepository) GetAllModels(tenantID string) ([]model.MdModel, error) {
	var models []model.MdModel
	result := r.db.Where("tenant_id = ?", tenantID).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}
	return models, nil
}
