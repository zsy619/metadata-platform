package service

import (
	"metadata-platform/internal/module/audit/queue"
	auditService "metadata-platform/internal/module/audit/service"
	"metadata-platform/internal/module/metadata/engine"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/repository"

	"gorm.io/gorm"
)

// APIService API服务接口
type APIService interface {
	CreateAPI(api *model.API) error
	GetAPIByID(id string) (*model.API, error)
	GetAPIByCode(code string) (*model.API, error)
	UpdateAPI(api *model.API) error
	DeleteAPI(id string) error
	GetAllAPIs() ([]model.API, error)
}

// Services 元数据模块服务集合
type Services struct {
	API              APIService
	Conn             MdConnService
	Table            MdTableService
	TableField       MdTableFieldService
	Model            MdModelService
	FieldEnhancement MdModelFieldEnhancementService
	CRUD             CRUDService
	APIGenerator     APIGenerator
	Validator        DataValidator
	QueryTemplate    QueryTemplateService
	Tree             TreeService
	MasterDetail     MasterDetailService
	DataIO           DataIOService
	Audit            auditService.AuditService
}

// NewServices 创建元数据模块服务集合
func NewServices(db *gorm.DB, repos *repository.Repositories, auditDB *gorm.DB, auditQueue *queue.AuditLogQueue) *Services {
	connService := NewMdConnService(repos.Conn)

	validator := NewDataValidator()
	queryTemplateService := NewQueryTemplateService(repos.QueryTemplate, repos.QueryCondition)

	// 初始化 SQL 引擎
	sqlBuilder := engine.NewSQLBuilder(db, repos.Model)
	sqlExecutor := engine.NewSQLExecutor(db, repos.Conn)
	auditSvc := auditService.NewAuditService(auditDB, auditQueue)
	crudSvc := NewCRUDService(sqlBuilder, sqlExecutor, validator, queryTemplateService, auditSvc)
	treeSvc := NewTreeService(repos.Model, crudSvc, sqlExecutor)
	masterDetailSvc := NewMasterDetailService(crudSvc, repos.ModelRelation, repos.Model, sqlExecutor)
	dataIOSvc := NewDataIOService(crudSvc, repos.Model, repos.ModelField, validator)

	return &Services{
		API:              NewAPIService(repos.API),
		Conn:             connService,
		Table:            NewMdTableService(repos.Table),
		TableField:       NewMdTableFieldService(repos.TableField),
		Model:            NewMdModelService(repos.Model, repos.ModelField, connService),
		FieldEnhancement: NewMdModelFieldEnhancementService(repos.FieldEnhancement),
		CRUD:             crudSvc,
		APIGenerator:     NewAPIGenerator(repos.Model, repos.API),
		Validator:        validator,
		QueryTemplate:    queryTemplateService,
		Tree:             treeSvc,
		MasterDetail:     masterDetailSvc,
		DataIO:           dataIOSvc,
		Audit:            auditSvc,
	}
}

type apiService struct {
	repo repository.APIRepository
}

// NewAPIService 创建API服务实例
func NewAPIService(repo repository.APIRepository) APIService {
	return &apiService{repo: repo}
}

func (s *apiService) CreateAPI(api *model.API) error {
	return s.repo.CreateAPI(api)
}

func (s *apiService) GetAPIByID(id string) (*model.API, error) {
	return s.repo.GetAPIByID(id)
}

func (s *apiService) GetAPIByCode(code string) (*model.API, error) {
	return s.repo.GetAPIByCode(code)
}

func (s *apiService) UpdateAPI(api *model.API) error {
	return s.repo.UpdateAPI(api)
}

func (s *apiService) DeleteAPI(id string) error {
	return s.repo.DeleteAPI(id)
}

func (s *apiService) GetAllAPIs() ([]model.API, error) {
	return s.repo.GetAllAPIs()
}
