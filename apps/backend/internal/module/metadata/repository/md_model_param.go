package repository

import (
	"metadata-platform/internal/module/metadata/model"

	"gorm.io/gorm"
)

// MdModelParamRepository 模型参数仓库接口
type MdModelParamRepository interface {
	Create(param *model.MdModelParam) error
	BatchCreate(params []model.MdModelParam) error
	GetByModelID(modelID string) ([]model.MdModelParam, error)
	DeleteByModelID(modelID string) error
}

type mdModelParamRepository struct {
	db *gorm.DB
}

// NewMdModelParamRepository 创建模型参数仓库实例
func NewMdModelParamRepository(db *gorm.DB) MdModelParamRepository {
	return &mdModelParamRepository{db: db}
}

func (r *mdModelParamRepository) Create(param *model.MdModelParam) error {
	return r.db.Create(param).Error
}

func (r *mdModelParamRepository) BatchCreate(params []model.MdModelParam) error {
	if len(params) == 0 {
		return nil
	}
	return r.db.Create(&params).Error
}

func (r *mdModelParamRepository) GetByModelID(modelID string) ([]model.MdModelParam, error) {
	var res []model.MdModelParam
	err := r.db.Where("model_id = ? AND is_deleted = ?", modelID, false).Find(&res).Error
	return res, err
}

func (r *mdModelParamRepository) DeleteByModelID(modelID string) error {
	return r.db.Model(&model.MdModelParam{}).Where("model_id = ?", modelID).Update("is_deleted", true).Error
}
