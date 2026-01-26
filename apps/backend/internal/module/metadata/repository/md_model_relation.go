package repository

import (
	"metadata-platform/internal/module/metadata/model"

	"gorm.io/gorm"
)

// MdModelRelationRepository 模型关联仓库接口
type MdModelRelationRepository interface {
	CreateRelation(rel *model.MdModelRelation) error
	GetRelationByMasterID(masterID string) ([]model.MdModelRelation, error)
	GetRelationByDetailID(detailID string) ([]model.MdModelRelation, error)
	GetRelation(masterID, detailID string) (*model.MdModelRelation, error)
	DeleteRelation(id string) error
}

type mdModelRelationRepository struct {
	db *gorm.DB
}

// NewMdModelRelationRepository 创建模型关联仓库实例
func NewMdModelRelationRepository(db *gorm.DB) MdModelRelationRepository {
	return &mdModelRelationRepository{db: db}
}

func (r *mdModelRelationRepository) CreateRelation(rel *model.MdModelRelation) error {
	return r.db.Create(rel).Error
}

func (r *mdModelRelationRepository) GetRelationByMasterID(masterID string) ([]model.MdModelRelation, error) {
	var rels []model.MdModelRelation
	err := r.db.Where("master_model_id = ?", masterID).Find(&rels).Error
	return rels, err
}

func (r *mdModelRelationRepository) GetRelationByDetailID(detailID string) ([]model.MdModelRelation, error) {
	var rels []model.MdModelRelation
	err := r.db.Where("detail_model_id = ?", detailID).Find(&rels).Error
	return rels, err
}

func (r *mdModelRelationRepository) GetRelation(masterID, detailID string) (*model.MdModelRelation, error) {
	var rel model.MdModelRelation
	err := r.db.Where("master_model_id = ? AND detail_model_id = ?", masterID, detailID).First(&rel).Error
	if err != nil {
		return nil, err
	}
	return &rel, nil
}

func (r *mdModelRelationRepository) DeleteRelation(id string) error {
	return r.db.Delete(&model.MdModelRelation{}, "id = ?", id).Error
}
