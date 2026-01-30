package service

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/repository"
	"metadata-platform/internal/utils"
	"strings"
	"time"
)

// MdModelService 模型定义服务接口
type MdModelService interface {
	CreateModel(model *model.MdModel) error
	GetModelByID(id string) (*model.MdModel, error)
	GetModelByCode(code string) (*model.MdModel, error)
	UpdateModel(model *model.MdModel) error
	DeleteModel(id string) error
	GetModels(tenantID string, page, pageSize int, search string, modelKind int) ([]model.MdModel, int64, error)
	GetModelsByConnID(connID string) ([]model.MdModel, error)
	GetAllModels(tenantID string) ([]model.MdModel, error)
	BuildFromTable(req *BuildFromTableRequest) error
	BuildFromView(req *BuildFromViewRequest) error
	BuildFromSQL(req *BuildFromSQLRequest) error
	UpdateSQLModel(req *UpdateSQLModelRequest) error
	GetSQLByModelID(modelID string) (*model.MdModelSql, error)
	TestSQL(req *TestSQLRequest) ([]FieldMapping, error)
	// ModelField operations
	GetFieldsByModelID(modelID string) ([]model.MdModelField, error)
	CreateField(field *model.MdModelField) error
	UpdateField(field *model.MdModelField) error
	DeleteField(id string) error
	Generate16Code() string
	Generate32Code() string
	Generate64Code() string
	GetModelParams(modelID string) ([]model.MdModelParam, error)
	SaveVisualModel(req *SaveVisualModelRequest) (*model.MdModel, error)
}

type SaveVisualModelRequest struct {
	ModelID      string
	ConnID       string
	ModelName    string
	ModelCode    string
	ModelVersion string
	ModelKind    int
	IsPublic     bool
	Remark       string
	Parameters   string
	Tables       []model.MdModelTable
	Fields       []model.MdModelField
	Joins        []model.MdModelJoin
	JoinFields   []model.MdModelJoinField
	Wheres       []model.MdModelWhere
	Orders       []model.MdModelOrder
	Groups       []model.MdModelGroup
	Havings      []model.MdModelHaving
	TenantID     string
	UserID       string
	Username     string
}

type BuildFromViewRequest struct {
	ConnID    string
	Schema    string
	View      string
	ModelName string
	ModelCode string
	TenantID  string
	UserID    string
	Username  string
}

type BuildFromSQLRequest struct {
	ConnID        string
	ModelName     string
	ModelCode     string
	SQLContent    string
	Parameters    []SQLParameter
	FieldMappings []FieldMapping
	TenantID      string
	UserID        string
	Username      string
}

type UpdateSQLModelRequest struct {
	ModelID       string                 `json:"model_id" binding:"required"`
	ModelName     string                 `json:"model_name" binding:"required"`
	ModelCode     string                 `json:"model_code" binding:"required"`
	SQLContent    string                 `json:"sql_content" binding:"required"`
	Parameters    []SQLParameter         `json:"parameters"`
	FieldMappings []FieldMapping         `json:"field_mappings"`
	IsPublic      bool                   `json:"is_public"`
	Remark        string                 `json:"remark"`
	TenantID      string
	UserID        string
	Username      string
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
	ConnID    string
	Schema    string
	Table     string
	ModelName string
	ModelCode string
	TenantID  string
	UserID    string
	Username  string
}

// mdModelService 模型定义服务实现
type mdModelService struct {
	modelRepo      repository.MdModelRepository
	fieldRepo      repository.MdModelFieldRepository
	modelSqlRepo   repository.MdModelSqlRepository
	modelParamRepo repository.MdModelParamRepository
	connService    MdConnService
	snowflake      *utils.Snowflake
}

// NewMdModelService 创建模型定义服务实例
func NewMdModelService(modelRepo repository.MdModelRepository, fieldRepo repository.MdModelFieldRepository, modelSqlRepo repository.MdModelSqlRepository, modelParamRepo repository.MdModelParamRepository, connService MdConnService) MdModelService {
	// 创建雪花算法生成器实例，使用默认数据中心ID和机器ID
	snowflake := utils.NewSnowflake(1, 1)
	return &mdModelService{
		modelRepo:      modelRepo,
		fieldRepo:      fieldRepo,
		modelSqlRepo:   modelSqlRepo,
		modelParamRepo: modelParamRepo,
		connService:    connService,
		snowflake:      snowflake,
	}
}

// CreateModel 创建模型定义
func (s *mdModelService) CreateModel(model *model.MdModel) error {
	// 使用雪花算法生成唯一ID
	model.ID = s.snowflake.GenerateIDString()

	// 如果编码为空，自动生成 32 位编码
	if model.ModelCode == "" {
		model.ModelCode = s.Generate32Code()
	}

	// 检查模型编码是否已存在
	existingModel, err := s.modelRepo.GetModelByCode(model.ModelCode)
	if err == nil && existingModel != nil {
		return errors.New("模型编码已存在")
	}

	// 创建模型
	return s.modelRepo.CreateModel(model)
}

// Generate16Code 生成 16 位唯一编码 (MD5(SnowflakeID) 的前 16 位)
func (s *mdModelService) Generate16Code() string {
	id := s.snowflake.GenerateIDString()
	hash := md5.Sum([]byte(id))
	return strings.ToUpper(hex.EncodeToString(hash[:8]))
}

// Generate32Code 生成 32 位唯一编码 (MD5(SnowflakeID))
func (s *mdModelService) Generate32Code() string {
	id := s.snowflake.GenerateIDString()
	hash := md5.Sum([]byte(id))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// Generate64Code 生成 64 位唯一编码 (SHA256(SnowflakeID))
func (s *mdModelService) Generate64Code() string {
	id := s.snowflake.GenerateIDString()
	hash := sha256.Sum256([]byte(id))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// GetModelParams 获取模型参数列表
func (s *mdModelService) GetModelParams(modelID string) ([]model.MdModelParam, error) {
	return s.modelParamRepo.GetByModelID(modelID)
}

// GetSQLByModelID 根据模型ID获取 SQL 内容
func (s *mdModelService) GetSQLByModelID(modelID string) (*model.MdModelSql, error) {
	return s.modelSqlRepo.GetByModelID(modelID)
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

// GetModels 获取模型定义列表（支持分页、搜索和类型过滤）
func (s *mdModelService) GetModels(tenantID string, page, pageSize int, search string, modelKind int) ([]model.MdModel, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.modelRepo.GetModels(tenantID, offset, pageSize, search, modelKind)
}

// GetModelsByConnID 根据连接ID获取模型定义列表
func (s *mdModelService) GetModelsByConnID(connID string) ([]model.MdModel, error) {
	return s.modelRepo.GetModelsByConnID(connID)
}

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
	conn, err := s.connService.GetConnByID(req.ConnID)
	if err != nil {
		return err
	}

	columns, err := s.connService.GetTableStructure(conn, req.Schema, req.Table)
	if err != nil {
		return err
	}

	modelID := s.snowflake.GenerateIDString()
	if req.ModelCode == "" {
		req.ModelCode = s.Generate32Code()
	}

	mdModel := &model.MdModel{
		ID:           modelID,
		TenantID:     req.TenantID,
		ParentID:     "0",
		ConnID:       req.ConnID,
		ConnName:     conn.ConnName,
		ModelName:    req.ModelName,
		ModelCode:    req.ModelCode,
		ModelVersion: "1.0.0",
		ModelKind:    2,
		IsPublic:     false,
		IsLocked:     false,
		IsDeleted:    false,
		CreateID:     req.UserID,
		CreateBy:     req.Username,
		UpdateID:     req.UserID,
		UpdateBy:     req.Username,
	}

	existingModel, err := s.modelRepo.GetModelByCode(req.ModelCode)
	if err == nil && existingModel != nil {
		return errors.New("模型编码已存在")
	}

	if err := s.modelRepo.CreateModel(mdModel); err != nil {
		return err
	}

	for _, col := range columns {
		dataType := mapDataType(col.Type)
		field := &model.MdModelField{
			ID:           s.snowflake.GenerateIDString(),
			TenantID:     req.TenantID,
			ModelID:      modelID,
			TableSchema:  req.Schema,
			TableNameStr: req.Table,
			TableTitle:   req.Table,
			ColumnName:   col.Name,
			ColumnTitle:  col.Comment,
			FieldType:    dataType,
			MaxLength:    col.Length,
			ShowTitle:    col.Comment,
			ShowWidth:    150,
			IsDeleted:    false,
			CreateID:     req.UserID,
			CreateBy:     req.Username,
			UpdateID:     req.UserID,
			UpdateBy:     req.Username,
		}

		if col.DefaultValue != nil {
			field.DefaultValue = fmt.Sprintf("%v", col.DefaultValue)
		}
		if field.ColumnTitle == "" {
			field.ColumnTitle = col.Name
			field.ShowTitle = col.Name
		}

		if err := s.fieldRepo.CreateField(field); err != nil {
			return err
		}
	}
	return nil
}

// BuildFromSQL 从 SQL 构建模型
func (s *mdModelService) BuildFromSQL(req *BuildFromSQLRequest) error {
	conn, err := s.connService.GetConnByID(req.ConnID)
	if err != nil {
		return err
	}

	modelID := s.snowflake.GenerateIDString()
	if req.ModelCode == "" {
		req.ModelCode = s.Generate32Code()
	}

	paramsJson, _ := json.Marshal(req.Parameters)

	mdModel := &model.MdModel{
		ID:           modelID,
		TenantID:     req.TenantID,
		ParentID:     "0",
		ConnID:       req.ConnID,
		ConnName:     conn.ConnName,
		ModelName:    req.ModelName,
		ModelCode:    req.ModelCode,
		ModelVersion: "1.0.0",
		ModelKind:    1,
		Parameters:   string(paramsJson),
		IsPublic:     false,
		IsLocked:     false,
		IsDeleted:    false,
		CreateID:     req.UserID,
		CreateBy:     req.Username,
		UpdateID:     req.UserID,
		UpdateBy:     req.Username,
	}

	existingModel, err := s.modelRepo.GetModelByCode(req.ModelCode)
	if err == nil && existingModel != nil {
		return errors.New("模型编码已存在")
	}

	if err := s.modelRepo.CreateModel(mdModel); err != nil {
		return err
	}

	modelSql := &model.MdModelSql{
		ID:       s.snowflake.GenerateIDString(),
		TenantID: req.TenantID,
		ModelID:  modelID,
		Content:  req.SQLContent,
		CreateID: req.UserID,
		CreateBy: req.Username,
		UpdateID: req.UserID,
		UpdateBy: req.Username,
	}
	if err := s.modelSqlRepo.Create(modelSql); err != nil {
		return err
	}

	for _, mapping := range req.FieldMappings {
		field := &model.MdModelField{
			ID:          s.snowflake.GenerateIDString(),
			TenantID:    req.TenantID,
			ModelID:     modelID,
			ColumnName:  mapping.ColumnName,
			ColumnTitle: mapping.ShowTitle,
			ShowTitle:   mapping.ShowTitle,
			ShowWidth:   mapping.ShowWidth,
			FieldType:   "string",
			IsDeleted:   false,
			CreateID:    req.UserID,
			CreateBy:    req.Username,
			UpdateID:    req.UserID,
			UpdateBy:    req.Username,
		}
		if err := s.fieldRepo.CreateField(field); err != nil {
			return err
		}
	}

	// 保存参数到独立的参数表
	for _, param := range req.Parameters {
		modelParam := &model.MdModelParam{
			ID:       s.snowflake.GenerateIDString(),
			TenantID: req.TenantID,
			ModelID:  modelID,
			Name:     param.Name,
			Type:     param.Type,
			Required: param.Required,
			Default:  param.Default,
			CreateID: req.UserID,
			CreateBy: req.Username,
			UpdateID: req.UserID,
			UpdateBy: req.Username,
		}
		if err := s.modelParamRepo.Create(modelParam); err != nil {
			return err
		}
	}

	return nil
}

// UpdateSQLModel 更新 SQL 模型
func (s *mdModelService) UpdateSQLModel(req *UpdateSQLModelRequest) error {
	// 1. 获取原模型
	mod, err := s.modelRepo.GetModelByID(req.ModelID)
	if err != nil {
		return errors.New("模型不存在")
	}

	// 2. 检查编码冲突
	if mod.ModelCode != req.ModelCode {
		existing, err := s.modelRepo.GetModelByCode(req.ModelCode)
		if err == nil && existing != nil {
			return errors.New("模型编码已存在")
		}
	}

	// 3. 更新基本信息
	mod.ModelName = req.ModelName
	mod.ModelCode = req.ModelCode
	mod.IsPublic = req.IsPublic
	mod.Remark = req.Remark
	paramsJson, _ := json.Marshal(req.Parameters)
	mod.Parameters = string(paramsJson)
	mod.UpdateID = req.UserID
	mod.UpdateBy = req.Username

	if err := s.modelRepo.UpdateModel(mod); err != nil {
		return err
	}

	// 4. 更新 SQL 内容
	// 先获取 sql 记录
	modelSql, err := s.modelSqlRepo.GetByModelID(req.ModelID)
	// 如果不存在，则创建
	if err != nil || modelSql == nil {
		// 创建
		modelSql = &model.MdModelSql{
			ID:       s.snowflake.GenerateIDString(),
			TenantID: req.TenantID,
			ModelID:  req.ModelID,
			Content:  req.SQLContent,
			CreateID: req.UserID,
			CreateBy: req.Username,
			UpdateID: req.UserID,
			UpdateBy: req.Username,
		}
		if err := s.modelSqlRepo.Create(modelSql); err != nil {
			return err
		}
	} else {
		modelSql.Content = req.SQLContent
		modelSql.UpdateID = req.UserID
		modelSql.UpdateBy = req.Username
		if err := s.modelSqlRepo.Update(modelSql); err != nil {
			return err
		}
	}
	
	// 5. 更新字段 (删除原有，重新创建)
	if err := s.fieldRepo.DeleteFieldsByModelID(req.ModelID); err != nil {
		return err
	}

	for _, mapping := range req.FieldMappings {
		field := &model.MdModelField{
			ID:          s.snowflake.GenerateIDString(),
			TenantID:    req.TenantID,
			ModelID:     req.ModelID,
			ColumnName:  mapping.ColumnName,
			ColumnTitle: mapping.ShowTitle,
			ShowTitle:   mapping.ShowTitle,
			ShowWidth:   mapping.ShowWidth,
			FieldType:   "string",
			IsDeleted:   false,
			CreateID:    req.UserID,
			CreateBy:    req.Username,
			UpdateID:    req.UserID,
			UpdateBy:    req.Username,
		}
		if err := s.fieldRepo.CreateField(field); err != nil {
			return err
		}
	}

	// 6. 更新参数 (删除原有，重新创建)
	if err := s.modelParamRepo.DeleteByModelID(req.ModelID); err != nil {
		return err
	}
	
	for _, param := range req.Parameters {
		modelParam := &model.MdModelParam{
			ID:       s.snowflake.GenerateIDString(),
			TenantID: req.TenantID,
			ModelID:  req.ModelID,
			Name:     param.Name,
			Type:     param.Type,
			Required: param.Required,
			Default:  param.Default,
			CreateID: req.UserID,
			CreateBy: req.Username,
			UpdateID: req.UserID,
			UpdateBy: req.Username,
		}
		if err := s.modelParamRepo.Create(modelParam); err != nil {
			return err
		}
	}

	return nil
}

// TestSQL 测试/预览 SQL
func (s *mdModelService) TestSQL(req *TestSQLRequest) ([]FieldMapping, error) {
	conn, err := s.connService.GetConnByID(req.ConnID)
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	for _, p := range req.Parameters {
		params[p.Name] = p.Default
	}
	columns, err := s.connService.ExecuteSQLForColumns(conn, req.SQLContent, params)
	if err != nil {
		return nil, err
	}
	mappings := make([]FieldMapping, 0, len(columns))
	for _, col := range columns {
		mappings = append(mappings, FieldMapping{
			ColumnName: col.Name,
			ShowTitle:  col.Name,
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

// SaveVisualModel 全量保存可视化模型
func (s *mdModelService) SaveVisualModel(req *SaveVisualModelRequest) (*model.MdModel, error) {
	// 1. 准备模型主表数据
	mdModel := &model.MdModel{
		ID:           req.ModelID,
		TenantID:     req.TenantID,
		ParentID:     "0",
		ConnID:       req.ConnID,
		ModelName:    req.ModelName,
		ModelCode:    req.ModelCode,
		ModelVersion: req.ModelVersion,
		ModelKind:    req.ModelKind,
		ModelLogo:    "", // logo 暂时留空
		IsPublic:     req.IsPublic,
		Remark:       req.Remark,
		Parameters:   req.Parameters,
		IsLocked:     false,
		IsDeleted:    false,
		UpdateID:     req.UserID,
		UpdateBy:     req.Username,
	}

	// 如果是新建，生成ID
	if mdModel.ID == "" {
		mdModel.ID = s.snowflake.GenerateIDString()
		mdModel.CreateID = req.UserID
		mdModel.CreateBy = req.Username
		mdModel.CreateAt = time.Now()
	} else {
		// 如果是更新，尝试获取原记录以保留 Create 信息
		oldModel, err := s.modelRepo.GetModelByID(req.ModelID)
		if err == nil && oldModel != nil {
			mdModel.CreateID = oldModel.CreateID
			mdModel.CreateBy = oldModel.CreateBy
			mdModel.CreateAt = oldModel.CreateAt
		}
	}
	mdModel.UpdateAt = time.Now()

	// 2. 准备关联表数据
	// 需要为子表项生成/补全 ID (如果是新加的)
	// 同时补全关联外键 ModelID
	for i := range req.Tables {
		if req.Tables[i].ID == "" || strings.HasPrefix(req.Tables[i].ID, "tmp_") { // 假设前端临时ID
			req.Tables[i].ID = s.snowflake.GenerateIDString()
		}
		req.Tables[i].ModelID = mdModel.ID
		req.Tables[i].TenantID = req.TenantID
		req.Tables[i].CreateID = req.UserID
		req.Tables[i].CreateBy = req.Username
	}

	for i := range req.Fields {
		if req.Fields[i].ID == "" || strings.HasPrefix(req.Fields[i].ID, "tmp_") {
			req.Fields[i].ID = s.snowflake.GenerateIDString()
		}
		req.Fields[i].ModelID = mdModel.ID
		req.Fields[i].TenantID = req.TenantID
		req.Fields[i].CreateID = req.UserID
		req.Fields[i].CreateBy = req.Username
	}

	for i := range req.Joins {
		if req.Joins[i].ID == "" || strings.HasPrefix(req.Joins[i].ID, "tmp_") {
			req.Joins[i].ID = s.snowflake.GenerateIDString()
		}
		req.Joins[i].ModelID = mdModel.ID
		req.Joins[i].TenantID = req.TenantID
		req.Joins[i].CreateID = req.UserID
		req.Joins[i].CreateBy = req.Username
	}

	for i := range req.JoinFields {
		if req.JoinFields[i].ID == "" || strings.HasPrefix(req.JoinFields[i].ID, "tmp_") {
			req.JoinFields[i].ID = s.snowflake.GenerateIDString()
		}
		req.JoinFields[i].TenantID = req.TenantID
		req.JoinFields[i].CreateID = req.UserID
		req.JoinFields[i].CreateBy = req.Username
	}

	for i := range req.Wheres {
		if req.Wheres[i].ID == "" || strings.HasPrefix(req.Wheres[i].ID, "tmp_") {
			req.Wheres[i].ID = s.snowflake.GenerateIDString()
		}
		req.Wheres[i].ModelID = mdModel.ID
		req.Wheres[i].TenantID = req.TenantID
		req.Wheres[i].CreateID = req.UserID
		req.Wheres[i].CreateBy = req.Username
	}

	for i := range req.Orders {
		if req.Orders[i].ID == "" || strings.HasPrefix(req.Orders[i].ID, "tmp_") {
			req.Orders[i].ID = s.snowflake.GenerateIDString()
		}
		req.Orders[i].ModelID = mdModel.ID
		req.Orders[i].TenantID = req.TenantID
		req.Orders[i].CreateID = req.UserID
		req.Orders[i].CreateBy = req.Username
	}

	for i := range req.Groups {
		if req.Groups[i].ID == "" || strings.HasPrefix(req.Groups[i].ID, "tmp_") {
			req.Groups[i].ID = s.snowflake.GenerateIDString()
		}
		req.Groups[i].ModelID = mdModel.ID
		req.Groups[i].TenantID = req.TenantID
		req.Groups[i].CreateID = req.UserID
		req.Groups[i].CreateBy = req.Username
	}

	for i := range req.Havings {
		if req.Havings[i].ID == "" || strings.HasPrefix(req.Havings[i].ID, "tmp_") {
			req.Havings[i].ID = s.snowflake.GenerateIDString()
		}
		req.Havings[i].ModelID = mdModel.ID
		req.Havings[i].TenantID = req.TenantID
		req.Havings[i].CreateID = req.UserID
		req.Havings[i].CreateBy = req.Username
	}

	// 3. 调用 Repo 执行全量保存
	if err := s.modelRepo.SaveVisualModel(mdModel, req.Tables, req.Fields, req.Joins, req.JoinFields, req.Wheres, req.Orders, req.Groups, req.Havings); err != nil {
		return nil, err
	}

	return mdModel, nil
}
