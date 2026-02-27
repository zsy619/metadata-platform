package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/metadata/model"
)

// APIRepository API仓库接口
type APIRepository interface {
	CreateAPI(api *model.API) error
	GetAPIByID(id string) (*model.API, error)
	GetAPIByCode(code string) (*model.API, error)
	UpdateAPI(api *model.API) error
	DeleteAPI(id string) error
	GetAllAPIs() ([]model.API, error)
}

// Repositories 元数据模块仓库集合
type Repositories struct {
	API              APIRepository
	Conn             MdConnRepository
	Table            MdTableRepository
	TableField       MdTableFieldRepository
	Model            MdModelRepository
	ModelField       MdModelFieldRepository
	Procedure        MdModelProcedureRepository
	FieldEnhancement MdModelFieldEnhancementRepository
	QueryTemplate    MdQueryTemplateRepository
	QueryCondition   MdQueryConditionRepository
	ModelRelation    MdModelRelationRepository
	ModelSql         MdModelSqlRepository
	ModelParam       MdModelParamRepository
}

// NewRepositories 创建元数据模块仓库集合
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		API:              NewAPIRepository(db),
		Conn:             NewMdConnRepository(db),
		Table:            NewMdTableRepository(db),
		TableField:       NewMdTableFieldRepository(db),
		Model:            NewMdModelRepository(db),
		ModelField:       NewMdModelFieldRepository(db),
		Procedure:        NewMdModelProcedureRepository(db),
		FieldEnhancement: NewMdModelFieldEnhancementRepository(db),
		QueryTemplate:    NewMdQueryTemplateRepository(db),
		QueryCondition:   NewMdQueryConditionRepository(db),
		ModelRelation:    NewMdModelRelationRepository(db),
		ModelSql:         NewMdModelSqlRepository(db),
		ModelParam:       NewMdModelParamRepository(db),
	}
}

// apiRepository API仓库实现
type apiRepository struct {
	db *gorm.DB
}

// NewAPIRepository 创建API仓库实例
func NewAPIRepository(db *gorm.DB) APIRepository {
	return &apiRepository{db: db}
}

// CreateAPI 创建API
func (r *apiRepository) CreateAPI(api *model.API) error {
	return r.db.Create(api).Error
}

// GetAPIByID 根据ID获取API
func (r *apiRepository) GetAPIByID(id string) (*model.API, error) {
	var api model.API
	result := r.db.Where("id = ?", id).First(&api)
	if result.Error != nil {
		return nil, result.Error
	}
	return &api, nil
}

// GetAPIByCode 根据编码获取API
func (r *apiRepository) GetAPIByCode(code string) (*model.API, error) {
	var api model.API
	result := r.db.Where("code = ?", code).First(&api)
	if result.Error != nil {
		return nil, result.Error
	}
	return &api, nil
}

// UpdateAPI 更新API
func (r *apiRepository) UpdateAPI(api *model.API) error {
	return r.db.Save(api).Error
}

// DeleteAPI 删除API
func (r *apiRepository) DeleteAPI(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.API{}).Error
}

// GetAllAPIs 获取所有API
func (r *apiRepository) GetAllAPIs() ([]model.API, error) {
	var apis []model.API
	result := r.db.Find(&apis)
	if result.Error != nil {
		return nil, result.Error
	}
	return apis, nil
}
