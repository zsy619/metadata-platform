package service

import (
	"errors"
	"fmt"

	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/repository"
	"metadata-platform/internal/utils"
)

// MdModelService 模型定义服务接口
type MdModelService interface {
	CreateModel(model *model.MdModel) error
	GetModelByID(id string) (*model.MdModel, error)
	GetModelByCode(code string) (*model.MdModel, error)
	UpdateModel(model *model.MdModel) error
	DeleteModel(id string) error
	GetModelsByConnID(connID string) ([]model.MdModel, error)
	GetAllModels(tenantID string) ([]model.MdModel, error)
	BuildFromTable(req *BuildFromTableRequest) error
	BuildFromView(req *BuildFromViewRequest) error
	BuildFromSQL(req *BuildFromSQLRequest) error
	TestSQL(req *TestSQLRequest) ([]FieldMapping, error)
	// ModelField operations
	GetFieldsByModelID(modelID string) ([]model.MdModelField, error)
	CreateField(field *model.MdModelField) error
	UpdateField(field *model.MdModelField) error
	DeleteField(id string) error
}

type BuildFromViewRequest struct {
	ConnID      string
	Schema      string
	View        string
	ModelName   string
	ModelCode   string
	TenantID    string
	UserID      string
	Username    string
}

type BuildFromSQLRequest struct {
	ConnID       string
	ModelName    string
	ModelCode    string
	SQLContent   string
	Parameters   []SQLParameter
	FieldMappings []FieldMapping
	TenantID     string
	UserID       string
	Username     string
}

type TestSQLRequest struct {
	ConnID     string
	SQLContent string
	Parameters []SQLParameter
}

type SQLParameter struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Default  string `json:"default"`
}

type FieldMapping struct {
	ColumnName string `json:"column_name"`
	ShowTitle  string `json:"show_title"`
	ShowWidth  int    `json:"show_width"`
	Format     string `json:"format"`
}

type BuildFromTableRequest struct {
	ConnID      string
	Schema      string
	Table       string
	ModelName   string
	ModelCode   string
	TenantID    string
	UserID      string
	Username    string
}

// mdModelService 模型定义服务实现
type mdModelService struct {
	modelRepo   repository.MdModelRepository
	fieldRepo   repository.MdModelFieldRepository
	modelSqlRepo repository.MdModelSqlRepository
	connService MdConnService
	snowflake   *utils.Snowflake
}

// NewMdModelService 创建模型定义服务实例
func NewMdModelService(modelRepo repository.MdModelRepository, fieldRepo repository.MdModelFieldRepository, modelSqlRepo repository.MdModelSqlRepository, connService MdConnService) MdModelService {
	// 创建雪花算法生成器实例，使用默认数据中心ID和机器ID
	snowflake := utils.NewSnowflake(1, 1)
	return &mdModelService{
		modelRepo:    modelRepo,
		fieldRepo:    fieldRepo,
		modelSqlRepo: modelSqlRepo,
		connService:  connService,
		snowflake:    snowflake,
	}
}

// CreateModel 创建模型定义
func (s *mdModelService) CreateModel(model *model.MdModel) error {
	// 使用雪花算法生成唯一ID
	model.ID = s.snowflake.GenerateIDString()

	// 检查模型编码是否已存在
	existingModel, err := s.modelRepo.GetModelByCode(model.ModelCode)
	if err == nil && existingModel != nil {
		return errors.New("模型编码已存在")
	}

	// 创建模型
	return s.modelRepo.CreateModel(model)
}

// GetModelByID 根据ID获取模型定义
func (s *mdModelService) GetModelByID(id string) (*model.MdModel, error) {
	return s.modelRepo.GetModelByID(id)
}

// GetModelByCode 根据编码获取模型定义
func (s *mdModelService) GetModelByCode(code string) (*model.MdModel, error) {
	return s.modelRepo.GetModelByCode(code)
}

// UpdateModel 更新模型定义
func (s *mdModelService) UpdateModel(model *model.MdModel) error {
	// 检查模型是否存在
	existingModel, err := s.modelRepo.GetModelByID(model.ID)
	if err != nil {
		return errors.New("模型不存在")
	}

	// 如果模型编码发生变化，检查新编码是否已存在
	if existingModel.ModelCode != model.ModelCode {
		otherModel, err := s.modelRepo.GetModelByCode(model.ModelCode)
		if err == nil && otherModel != nil {
			return errors.New("模型编码已存在")
		}
	}

	// 更新模型定义
	return s.modelRepo.UpdateModel(model)
}

// DeleteModel 删除模型定义
func (s *mdModelService) DeleteModel(id string) error {
	// 检查模型是否存在
	_, err := s.modelRepo.GetModelByID(id)
	if err != nil {
		return errors.New("模型不存在")
	}

	// 删除模型定义
	return s.modelRepo.DeleteModel(id)
}

// GetModelsByConnID 根据连接ID获取模型定义列表
func (s *mdModelService) GetModelsByConnID(connID string) ([]model.MdModel, error) {
	return s.modelRepo.GetModelsByConnID(connID)
}

// GetAllModels 获取所有模型定义
// GetAllModels 获取所有模型定义
func (s *mdModelService) GetAllModels(tenantID string) ([]model.MdModel, error) {
	return s.modelRepo.GetAllModels(tenantID)
}

// GetFieldsByModelID 获取模型下的所有字段
func (s *mdModelService) GetFieldsByModelID(modelID string) ([]model.MdModelField, error) {
	return s.fieldRepo.GetFieldsByModelID(modelID)
}

// CreateField 创建模型字段
func (s *mdModelService) CreateField(field *model.MdModelField) error {
	return s.fieldRepo.CreateField(field)
}

// UpdateField 更新模型字段
func (s *mdModelService) UpdateField(field *model.MdModelField) error {
	return s.fieldRepo.UpdateField(field)
}

// DeleteField 删除模型字段
func (s *mdModelService) DeleteField(id string) error {
	return s.fieldRepo.DeleteField(id)
}

// BuildFromView 从视图构建模型
func (s *mdModelService) BuildFromView(req *BuildFromViewRequest) error {
	// 复用 BuildFromTable 的逻辑，但由于视图也是表结构，流程基本一致
	// 唯一的区别可能在于模型类型标记或者某些属性
	tableReq := &BuildFromTableRequest{
		ConnID:    req.ConnID,
		Schema:    req.Schema,
		Table:     req.View,
		ModelName: req.ModelName,
		ModelCode: req.ModelCode,
		TenantID:  req.TenantID,
		UserID:    req.UserID,
		Username:  req.Username,
	}
	return s.BuildFromTable(tableReq)
}

// BuildFromTable 从表构建模型
func (s *mdModelService) BuildFromTable(req *BuildFromTableRequest) error {
	// 1. 获取数据连接信息
	conn, err := s.connService.GetConnByID(req.ConnID)
	if err != nil {
		return err
	}

	// 2. 获取表结构
	columns, err := s.connService.GetTableStructure(conn, req.Schema, req.Table)
	if err != nil {
		return err
	}

	// 3. 创建模型头信息
	modelID := s.snowflake.GenerateIDString()
	mdModel := &model.MdModel{
		ID:           modelID,
		TenantID:     req.TenantID,
		ParentID:     "0", // 默认根目录
		ConnID:       req.ConnID,
		ConnName:     conn.ConnName,
		ModelName:    req.ModelName,
		ModelCode:    req.ModelCode,
		ModelVersion: "1.0.0",
		ModelKind:    2, // 2: 视图/表
		IsPublic:     false,
		IsLocked:     false,
		IsDeleted:    false,
		CreateID:     req.UserID,
		CreateBy:     req.Username,
		UpdateID:     req.UserID,
		UpdateBy:     req.Username,
	}

	// 检查模型编码是否已存在
	existingModel, err := s.modelRepo.GetModelByCode(req.ModelCode)
	if err == nil && existingModel != nil {
		return errors.New("模型编码已存在")
	}

	if err := s.modelRepo.CreateModel(mdModel); err != nil {
		return err
	}

	// 4. 创建模型字段信息
	for _, col := range columns {
		// 类型映射
		dataType := mapDataType(col.Type)
		
		field := &model.MdModelField{
			ID:              s.snowflake.GenerateIDString(),
			TenantID:        req.TenantID,
			ModelID:         modelID,
			TableSchema:     req.Schema,
			TableID:         "",
			TableNameStr:    req.Table,
			TableTitle:      req.Table,
			ColumnID:        "",
			ColumnName:      col.Name,
			ColumnTitle:     col.Comment,
			Func:            "", // 默认无函数
			AggFunc:         "", // 默认无聚合
			ColumnType:      col.Type,
			ColumnLength:    col.Length,
			IsNullable:      col.IsNullable,
			IsPrimaryKey:    col.IsPrimaryKey,
			IsAutoIncrement: col.IsAutoIncrement,
			FieldType:       dataType,
			MaxLength:       col.Length,
			ShowTitle:       col.Comment,
			ShowWidth:       150,
			IsDeleted:       false,
			CreateID:        req.UserID,
			CreateBy:        req.Username,
			UpdateID:        req.UserID,
			UpdateBy:        req.Username,
		}

		// 默认值处理
		if col.DefaultValue != nil {
			field.DefaultValue = fmt.Sprintf("%v", col.DefaultValue)
		}
		
		// 如果没有注释，使用列名作为显示标题
		if field.ColumnTitle == "" {
			field.ColumnTitle = col.Name
			field.ShowTitle = col.Name
		}

		if err := s.fieldRepo.CreateField(field); err != nil {
			// TODO: 考虑回滚模型创建
			return err
		}
	}

	return nil
}

// BuildFromSQL 从 SQL 构建模型
func (s *mdModelService) BuildFromSQL(req *BuildFromSQLRequest) error {
	// 1. 获取数据连接信息
	conn, err := s.connService.GetConnByID(req.ConnID)
	if err != nil {
		return err
	}

	// 2. 创建模型头信息
	modelID := s.snowflake.GenerateIDString()
	mdModel := &model.MdModel{
		ID:           modelID,
		TenantID:     req.TenantID,
		ParentID:     "0",
		ConnID:       req.ConnID,
		ConnName:     conn.ConnName,
		ModelName:    req.ModelName,
		ModelCode:    req.ModelCode,
		ModelVersion: "1.0.0",
		ModelKind:    1, // 1: SQL 语句
		IsPublic:     false,
		IsLocked:     false,
		IsDeleted:    false,
		CreateID:     req.UserID,
		CreateBy:     req.Username,
		UpdateID:     req.UserID,
		UpdateBy:     req.Username,
	}

	// 检查模型编码是否已存在
	existingModel, err := s.modelRepo.GetModelByCode(req.ModelCode)
	if err == nil && existingModel != nil {
		return errors.New("模型编码已存在")
	}

	if err := s.modelRepo.CreateModel(mdModel); err != nil {
		return err
	}

	// 3. 保存 SQL 内容
	modelSql := &model.MdModelSql{
		ID:        s.snowflake.GenerateIDString(),
		TenantID:  req.TenantID,
		ModelID:   modelID,
		Content:   req.SQLContent,
		Remark:    "",
		CreateID:  req.UserID,
		CreateBy:  req.Username,
		UpdateID:  req.UserID,
		UpdateBy:  req.Username,
	}
	if err := s.modelSqlRepo.Create(modelSql); err != nil {
		return err
	}

	// 4. 创建模型字段信息 (基于映射配置)
	for _, mapping := range req.FieldMappings {
		field := &model.MdModelField{
			ID:           s.snowflake.GenerateIDString(),
			TenantID:     req.TenantID,
			ModelID:      modelID,
			ColumnName:   mapping.ColumnName,
			ColumnTitle:  mapping.ShowTitle,
			ShowTitle:    mapping.ShowTitle,
			ShowWidth:    mapping.ShowWidth,
			FieldType:    "string", // 默认 string，后续可通过预览解析更精确
			IsDeleted:    false,
			CreateID:     req.UserID,
			CreateBy:     req.Username,
			UpdateID:     req.UserID,
			UpdateBy:     req.Username,
		}
		if err := s.fieldRepo.CreateField(field); err != nil {
			return err
		}
	}

	return nil
	return nil
}

// TestSQL 测试/预览 SQL，返回字段映射
func (s *mdModelService) TestSQL(req *TestSQLRequest) ([]FieldMapping, error) {
	// 1. 获取数据连接
	conn, err := s.connService.GetConnByID(req.ConnID)
	if err != nil {
		return nil, err
	}

	// 2. 执行 SQL 获取列信息 (仅获取结构，不获取数据)
	// 构建参数 map
	params := make(map[string]interface{})
	for _, p := range req.Parameters {
		params[p.Name] = p.Default
	}

	// 3. 执行查询获取列信息
	columns, err := s.connService.ExecuteSQLForColumns(conn, req.SQLContent, params)
	if err != nil {
		return nil, err
	}

	// 4. 转换为 FieldMapping
	mappings := make([]FieldMapping, 0, len(columns))
	for _, col := range columns {
		mappings = append(mappings, FieldMapping{
			ColumnName: col.Name,
			ShowTitle:  col.Name, // 默认显示标题为列名
			ShowWidth:  150,
			Format:     "",
		})
	}

	return mappings, nil
}

// mapDataType 映射数据库类型到标准类型
func mapDataType(dbType string) string {
	switch dbType {
	case "int", "tinyint", "smallint", "mediumint", "integer", "bit":
		return "integer"
	case "bigint":
		return "long"
	case "varchar", "char", "text", "mediumtext", "longtext", "tinytext":
		return "string"
	case "datetime", "timestamp":
		return "datetime"
	case "date":
		return "date"
	case "time":
		return "time"
	case "decimal", "double", "float", "real", "numeric":
		return "decimal"
	case "boolean", "bool":
		return "boolean"
	default:
		return "string"
	}
}
