package repository

import (
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/utils"

	"gorm.io/gorm"
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
	SaveVisualModel(md *model.MdModel, tables []model.MdModelTable, fields []model.MdModelField, joins []model.MdModelJoin, wheres []model.MdModelWhere, orders []model.MdModelOrder, groups []model.MdModelGroup, havings []model.MdModelHaving) error
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
			&model.MdModelParam{},
		}

		// 循环删除各关联表中的记录
		for _, m := range relatedModels {
			if err := tx.Where("model_id = ?", id).Delete(m).Error; err != nil {
				utils.SugarLogger.Fatalf("删除关联表记录失败: %v", err)
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
// SaveVisualModel 全量事务保存可视化模型相关配置
func (r *mdModelRepository) SaveVisualModel(md *model.MdModel, tables []model.MdModelTable, fields []model.MdModelField, joins []model.MdModelJoin, wheres []model.MdModelWhere, orders []model.MdModelOrder, groups []model.MdModelGroup, havings []model.MdModelHaving) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 保存/更新模型主表
		if err := tx.Save(md).Error; err != nil {
			return err
		}

		// 2. 清理旧的关联记录
		relatedModels := []interface{}{
			&model.MdModelTable{},
			&model.MdModelField{},
			&model.MdModelJoin{},
			&model.MdModelWhere{},
			&model.MdModelOrder{},
			&model.MdModelGroup{},
			&model.MdModelHaving{},
			// 对可视化模型构建器，我们暂时只管理这几张表。其它如 enhancement/limit 可根据需要后续加入。
		}
		for _, m := range relatedModels {
			if err := tx.Where("model_id = ?", md.ID).Delete(m).Error; err != nil {
				return err
			}
		}

		// 3. 批量插入新记录 (如果存在)
		if len(tables) > 0 {
			if err := tx.Create(&tables).Error; err != nil {
				return err
			}
		}
		if len(fields) > 0 {
			if err := tx.Create(&fields).Error; err != nil {
				return err
			}
		}
		if len(joins) > 0 {
			if err := tx.Create(&joins).Error; err != nil {
				return err
			}
		}
		if len(wheres) > 0 {
			if err := tx.Create(&wheres).Error; err != nil {
				return err
			}
		}
		if len(orders) > 0 {
			if err := tx.Create(&orders).Error; err != nil {
				return err
			}
		}
		if len(groups) > 0 {
			if err := tx.Create(&groups).Error; err != nil {
				return err
			}
		}
		if len(havings) > 0 {
			if err := tx.Create(&havings).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
