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
	GetModels(tenantID string, offset, limit int, search string, modelKind int) ([]model.MdModel, int64, error)
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

// DeleteModel 删除模型定义（级联删除关联子表）
func (r *mdModelRepository) DeleteModel(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 定义需要清理的关联模型列表
		// 注意：这里按顺序清理所有关联 md_model 的表
		relatedModels := []interface{}{
			&model.MdModelField{},
			&model.MdModelFieldEnhancement{},
			&model.MdModelGroup{},
			&model.MdModelHaving{},
			&model.MdModelJoin{},
			&model.MdModelLimit{},
			&model.MdModelOrder{},
			&model.MdModelRelation{},
			&model.MdModelSql{},
			&model.MdModelTable{},
			&model.MdModelWhere{},
			&model.MdQueryTemplate{},
		}

		// 循环删除各关联表中的记录
		for _, m := range relatedModels {
			if err := tx.Where("model_id = ?", id).Delete(m).Error; err != nil {
				return err
			}
		}

		// 最后删除模型主表记录
		if err := tx.Where("id = ?", id).Delete(&model.MdModel{}).Error; err != nil {
			return err
		}

		return nil
	})
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

// GetModels 获取所有模型定义（支持分页、搜索和类型过滤）
func (r *mdModelRepository) GetModels(tenantID string, offset, limit int, search string, modelKind int) ([]model.MdModel, int64, error) {
	var models []model.MdModel
	var total int64

	// 构建基础查询
	query := r.db.Model(&model.MdModel{}).Where("tenant_id = ?", tenantID).Where("is_deleted = ?", false)

	// 添加检索条件
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("(model_name LIKE ? OR model_code LIKE ?)", searchPattern, searchPattern)
	}
	if modelKind > 0 {
		query = query.Where("model_kind = ?", modelKind)
	}

	// 统计总数 (使用 Session 确保 query 状态不被 Count 消耗)
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset(offset).Limit(limit).Order("create_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return models, total, nil
}

// GetAllModels 获取所有模型定义
func (r *mdModelRepository) GetAllModels(tenantID string) ([]model.MdModel, error) {
	var models []model.MdModel
	err := r.db.Where("tenant_id = ? AND is_deleted = ?", tenantID, false).Order("create_at DESC").Find(&models).Error
	return models, err
}
