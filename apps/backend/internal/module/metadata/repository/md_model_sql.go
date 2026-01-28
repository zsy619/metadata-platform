package repository

import (
	"metadata-platform/internal/module/metadata/model"

	"gorm.io/gorm"
)

// MdModelSqlRepository 模型-sql模型仓库接口
type MdModelSqlRepository interface {
	Create(sql *model.MdModelSql) error
	GetByModelID(modelID string) (*model.MdModelSql, error)
	Update(sql *model.MdModelSql) error
	DeleteByModelID(modelID string) error
}

type mdModelSqlRepository struct {
	db *gorm.DB
}

// NewMdModelSqlRepository 创建模型-sql模型仓库实例
func NewMdModelSqlRepository(db *gorm.DB) MdModelSqlRepository {
	return &mdModelSqlRepository{db: db}
}

func (r *mdModelSqlRepository) Create(sql *model.MdModelSql) error {
	return r.db.Create(sql).Error
}

func (r *mdModelSqlRepository) GetByModelID(modelID string) (*model.MdModelSql, error) {
	var res model.MdModelSql
	err := r.db.Where("model_id = ? AND is_deleted = ?", modelID, false).First(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *mdModelSqlRepository) Update(sql *model.MdModelSql) error {
	return r.db.Save(sql).Error
}

func (r *mdModelSqlRepository) DeleteByModelID(modelID string) error {
	return r.db.Model(&model.MdModelSql{}).Where("model_id = ?", modelID).Update("is_deleted", true).Error
}
