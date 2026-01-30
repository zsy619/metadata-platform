package api

import (
	"context"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// MdModelHandler 模型定义API处理器
type MdModelHandler struct {
	modelService service.MdModelService
}

// NewMdModelHandler 创建模型定义API处理器实例
func NewMdModelHandler(modelService service.MdModelService) *MdModelHandler {
	return &MdModelHandler{modelService: modelService}
}

// BuildFromTableRequest 从表构建模型请求
type BuildFromTableRequest struct {
	ConnID    string `json:"conn_id" binding:"required"`
	Schema    string `json:"schema"`
	Table     string `json:"table" binding:"required"`
	ModelName string `json:"model_name" binding:"required"`
	ModelCode string `json:"model_code" binding:"required"`
}

// BuildFromViewRequest 从视图构建模型请求
type BuildFromViewRequest struct {
	ConnID    string `json:"conn_id" binding:"required"`
	Schema    string `json:"schema"`
	View      string `json:"view" binding:"required"`
	ModelName string `json:"model_name" binding:"required"`
	ModelCode string `json:"model_code" binding:"required"`
}

// BuildFromSQLRequest 从 SQL 构建模型请求
type BuildFromSQLRequest struct {
	ConnID        string                 `json:"conn_id" binding:"required"`
	ModelName     string                 `json:"model_name" binding:"required"`
	ModelCode     string                 `json:"model_code" binding:"required"`
	SQLContent    string                 `json:"sql_content" binding:"required"`
	Parameters    []service.SQLParameter `json:"parameters"`
	FieldMappings []service.FieldMapping `json:"field_mappings"`
}

// UpdateSQLModelRequest 更新 SQL 模型请求
type UpdateSQLModelRequest struct {
	ModelID       string                 `json:"model_id" binding:"required"`
	ModelName     string                 `json:"model_name" binding:"required"`
	ModelCode     string                 `json:"model_code" binding:"required"`
	SQLContent    string                 `json:"sql_content" binding:"required"`
	Parameters    []service.SQLParameter `json:"parameters"`
	FieldMappings []service.FieldMapping `json:"field_mappings"`
	IsPublic      bool                   `json:"is_public"`
	Remark        string                 `json:"remark"`
}

// TestSQLRequest 测试 SQL 请求
type TestSQLRequest struct {
	ConnID     string                 `json:"conn_id" binding:"required"`
	SQLContent string                 `json:"sql_content" binding:"required"`
	Parameters []service.SQLParameter `json:"parameters"`
}

// CreateModelRequest 创建模型请求
type CreateModelRequest struct {
	ConnID       string `json:"conn_id" binding:"required"`
	ConnName     string `json:"conn_name"`
	ModelName    string `json:"model_name" binding:"required"`
	ModelCode    string `json:"model_code" binding:"required"`
	ModelVersion string `json:"model_version"`
	ModelLogo    string `json:"model_logo"`
	ModelKind    int    `json:"model_kind"`
	IsPublic     bool   `json:"is_public"`
	IsLocked     bool   `json:"is_locked"`
	ParentID     string `json:"parent_id"`
}

// UpdateModelRequest 更新模型请求
type UpdateModelRequest struct {
	ModelName    string `json:"model_name"`
	ModelVersion string `json:"model_version"`
	ModelLogo    string `json:"model_logo"`
	IsPublic     bool   `json:"is_public"`
	IsLocked     bool   `json:"is_locked"`
}

// CreateModelFieldRequest 创建模型字段请求
type CreateModelFieldRequest struct {
	ModelID     string `json:"model_id" binding:"required"`
	TableSchema string `json:"table_schema"`
	TableName   string `json:"table_name" binding:"required"`
	ColumnName  string `json:"column_name" binding:"required"`
	ColumnTitle string `json:"column_title"`
	Func        string `json:"func"`
	AggFunc     string `json:"agg_func"`
	ShowTitle   string `json:"show_title"`
	ShowWidth   int    `json:"show_width"`
}

// SaveVisualModelRequest 可视化构建保存模型请求
type SaveVisualModelRequest struct {
	ModelID      string               `json:"model_id"`
	ConnID       string               `json:"conn_id" binding:"required"`
	ModelName    string               `json:"model_name" binding:"required"`
	ModelCode    string               `json:"model_code" binding:"required"`
	ModelVersion string               `json:"model_version"`
	ModelKind    int                  `json:"model_kind"`
	IsPublic     bool                 `json:"is_public"`
	Remark       string               `json:"remark"`
	Parameters   string               `json:"parameters"`
	Tables       []model.MdModelTable     `json:"tables"`
	Fields       []model.MdModelField     `json:"fields"`
	Joins        []model.MdModelJoin      `json:"joins"`
	JoinFields   []model.MdModelJoinField `json:"join_fields"`
	Wheres       []model.MdModelWhere     `json:"wheres"`
	Orders       []model.MdModelOrder     `json:"orders"`
	Groups       []model.MdModelGroup     `json:"groups"`
	Havings      []model.MdModelHaving    `json:"havings"`
}

// UpdateModelFieldRequest 更新模型字段请求
type UpdateModelFieldRequest struct {
	ColumnTitle string `json:"column_title"`
	Func        string `json:"func"`
	AggFunc     string `json:"agg_func"`
	ShowTitle   string `json:"show_title"`
	ShowWidth   int    `json:"show_width"`
}

// BuildFromTable 从表构建模型
func (h *MdModelHandler) BuildFromTable(c context.Context, ctx *app.RequestContext) {
	var req BuildFromTableRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 从上下文获取租户ID和用户信息
	tenantID, _ := ctx.Get("tenant_id")
	userID, _ := ctx.Get("user_id")
	username, _ := ctx.Get("username")

	serviceReq := &service.BuildFromTableRequest{
		ConnID:    req.ConnID,
		Schema:    req.Schema,
		Table:     req.Table,
		ModelName: req.ModelName,
		ModelCode: req.ModelCode,
		TenantID:  strconv.FormatUint(uint64(tenantID.(uint)), 10),
		UserID:    userID.(string),
		Username:  username.(string),
	}

	if err := h.modelService.BuildFromTable(serviceReq); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "模型构建成功")
}

// BuildFromView 从视图构建模型
func (h *MdModelHandler) BuildFromView(c context.Context, ctx *app.RequestContext) {
	var req BuildFromViewRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 从上下文获取租户ID和用户信息
	tenantID, _ := ctx.Get("tenant_id")
	userID, _ := ctx.Get("user_id")
	username, _ := ctx.Get("username")

	serviceReq := &service.BuildFromViewRequest{
		ConnID:    req.ConnID,
		Schema:    req.Schema,
		View:      req.View,
		ModelName: req.ModelName,
		ModelCode: req.ModelCode,
		TenantID:  strconv.FormatUint(uint64(tenantID.(uint)), 10),
		UserID:    userID.(string),
		Username:  username.(string),
	}

	if err := h.modelService.BuildFromView(serviceReq); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "模型构建成功")
}

// BuildFromSQL 从 SQL 构建模型
func (h *MdModelHandler) BuildFromSQL(c context.Context, ctx *app.RequestContext) {
	var req BuildFromSQLRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	tenantID, _ := ctx.Get("tenant_id")
	userID, _ := ctx.Get("user_id")
	username, _ := ctx.Get("username")

	serviceReq := &service.BuildFromSQLRequest{
		ConnID:        req.ConnID,
		ModelName:     req.ModelName,
		ModelCode:     req.ModelCode,
		SQLContent:    req.SQLContent,
		Parameters:    req.Parameters,
		FieldMappings: req.FieldMappings,
		TenantID:      strconv.FormatUint(uint64(tenantID.(uint)), 10),
		UserID:        userID.(string),
		Username:      username.(string),
	}

	if err := h.modelService.BuildFromSQL(serviceReq); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "模型构建成功")
}

// UpdateSQLModel 更新 SQL 模型
func (h *MdModelHandler) UpdateSQLModel(c context.Context, ctx *app.RequestContext) {
	var req UpdateSQLModelRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	tenantID, _ := ctx.Get("tenant_id")
	userID, _ := ctx.Get("user_id")
	username, _ := ctx.Get("username")
	
	serviceReq := &service.UpdateSQLModelRequest{
		ModelID:       req.ModelID,
		ModelName:     req.ModelName,
		ModelCode:     req.ModelCode,
		SQLContent:    req.SQLContent,
		Parameters:    req.Parameters,
		FieldMappings: req.FieldMappings,
		IsPublic:      req.IsPublic,
		Remark:        req.Remark,
		TenantID:      strconv.FormatUint(uint64(tenantID.(uint)), 10),
		UserID:        userID.(string),
		Username:      username.(string),
	}

	if err := h.modelService.UpdateSQLModel(serviceReq); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "模型更新成功")
}

// TestSQL 测试/预览 SQL
func (h *MdModelHandler) TestSQL(c context.Context, ctx *app.RequestContext) {
	var req TestSQLRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	serviceReq := &service.TestSQLRequest{
		ConnID:     req.ConnID,
		SQLContent: req.SQLContent,
		Parameters: req.Parameters,
	}

	fields, err := h.modelService.TestSQL(serviceReq)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, map[string]interface{}{
		"message": "SQL 解析成功",
		"fields":  fields,
	})
}

// CreateModel 创建模型
func (h *MdModelHandler) CreateModel(c context.Context, ctx *app.RequestContext) {
	var req CreateModelRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	tenantID, _ := ctx.Get("tenant_id")
	userID, _ := ctx.Get("user_id")
	username, _ := ctx.Get("username")

	model := &model.MdModel{
		ParentID:     req.ParentID,
		ConnID:       req.ConnID,
		ConnName:     req.ConnName,
		ModelName:    req.ModelName,
		ModelCode:    req.ModelCode,
		ModelVersion: req.ModelVersion,
		ModelLogo:    req.ModelLogo,
		ModelKind:    req.ModelKind,
		IsPublic:     req.IsPublic,
		IsLocked:     req.IsLocked,
		TenantID:     strconv.FormatUint(uint64(tenantID.(uint)), 10),
		CreateID:     userID.(string),
		CreateBy:     username.(string),
		UpdateID:     userID.(string),
		UpdateBy:     username.(string),
	}

	if err := h.modelService.CreateModel(model); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, model)
}

// GetModelByID 获取模型
func (h *MdModelHandler) GetModelByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	model, err := h.modelService.GetModelByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "模型不存在")
		return
	}

	utils.SuccessResponse(ctx, model)
}

// GetSQLByModelID 获取模型 SQL 内容
func (h *MdModelHandler) GetSQLByModelID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	sql, err := h.modelService.GetSQLByModelID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "SQL内容不存在")
		return
	}

	utils.SuccessResponse(ctx, sql)
}

// UpdateModel 更新模型
func (h *MdModelHandler) UpdateModel(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req UpdateModelRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	model, err := h.modelService.GetModelByID(id)
	if err != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "模型不存在",
		})
		return
	}

	if req.ModelName != "" {
		model.ModelName = req.ModelName
	}
	if req.ModelVersion != "" {
		model.ModelVersion = req.ModelVersion
	}
	if req.ModelLogo != "" {
		model.ModelLogo = req.ModelLogo
	}
	model.IsPublic = req.IsPublic
	model.IsLocked = req.IsLocked

	userID, _ := ctx.Get("user_id")
	username, _ := ctx.Get("username")
	model.UpdateID = userID.(string)
	model.UpdateBy = username.(string)

	if err := h.modelService.UpdateModel(model); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, model)
}

// DeleteModel 删除模型
func (h *MdModelHandler) DeleteModel(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	if err := h.modelService.DeleteModel(id); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, nil)
}

// Generate16Code 生成 16 位模型编码
func (h *MdModelHandler) Generate16Code(c context.Context, ctx *app.RequestContext) {
	code := h.modelService.Generate16Code()
	code = strings.ToLower(code)
	utils.SuccessResponse(ctx, map[string]string{
		"code": code,
	})
}

// GenerateCode 生成 32 位模型编码
func (h *MdModelHandler) GenerateCode(c context.Context, ctx *app.RequestContext) {
	code := h.modelService.Generate32Code()
	code = strings.ToLower(code)
	utils.SuccessResponse(ctx, map[string]string{
		"code": code,
	})
}

// Generate64Code 生成 64 位模型编码
func (h *MdModelHandler) Generate64Code(c context.Context, ctx *app.RequestContext) {
	code := h.modelService.Generate64Code()
	code = strings.ToLower(code)
	utils.SuccessResponse(ctx, map[string]string{
		"code": code,
	})
}

// ListModels 获取模型列表（分页和搜索）
func (h *MdModelHandler) ListModels(c context.Context, ctx *app.RequestContext) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("page_size"))
	search := ctx.Query("search")
	modelKind, _ := strconv.Atoi(ctx.Query("model_kind"))

	tenantID, _ := ctx.Get("tenant_id")
	models, total, err := h.modelService.GetModels(strconv.FormatUint(uint64(tenantID.(uint)), 10), page, pageSize, search, modelKind)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	utils.SuccessWithPagination(ctx, models, total, page, pageSize)
}

// GetAllModels 获取所有模型
func (h *MdModelHandler) GetAllModels(c context.Context, ctx *app.RequestContext) {
	tenantID, _ := ctx.Get("tenant_id")
	models, err := h.modelService.GetAllModels(strconv.FormatUint(uint64(tenantID.(uint)), 10))
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, models)
}

// GetModelsByConnID 根据连接ID获取模型
func (h *MdModelHandler) GetModelsByConnID(c context.Context, ctx *app.RequestContext) {
	connID := ctx.Param("conn_id")
	models, err := h.modelService.GetModelsByConnID(connID)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, models)
}

// GetFieldsByModelID 获取模型字段列表
func (h *MdModelHandler) GetFieldsByModelID(c context.Context, ctx *app.RequestContext) {
	modelID := ctx.Param("id")
	fields, err := h.modelService.GetFieldsByModelID(modelID)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, fields)
}

// GetModelParams 获取模型参数列表
func (h *MdModelHandler) GetModelParams(c context.Context, ctx *app.RequestContext) {
	modelID := ctx.Param("id")
	params, err := h.modelService.GetModelParams(modelID)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, params)
}

// CreateModelField 添加模型字段
func (h *MdModelHandler) CreateModelField(c context.Context, ctx *app.RequestContext) {
	var req CreateModelFieldRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	tenantID, _ := ctx.Get("tenant_id")
	userID, _ := ctx.Get("user_id")
	username, _ := ctx.Get("username")

	field := &model.MdModelField{
		ModelID:      req.ModelID,
		TableSchema:  req.TableSchema,
		TableNameStr: req.TableName,
		ColumnName:   req.ColumnName,
		ColumnTitle:  req.ColumnTitle,
		Func:         req.Func,
		AggFunc:      req.AggFunc,
		ShowTitle:    req.ShowTitle,
		ShowWidth:    req.ShowWidth,
		TenantID:     strconv.FormatUint(uint64(tenantID.(uint)), 10),
		CreateID:     userID.(string),
		CreateBy:     username.(string),
		UpdateID:     userID.(string),
		UpdateBy:     username.(string),
	}

	if err := h.modelService.CreateField(field); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, field)
}

// UpdateModelField 更新模型字段
func (h *MdModelHandler) UpdateModelField(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("fieldId")
	var req UpdateModelFieldRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 为简单起见，这里直接创建一个包含要更新字段的对象
	// 实际生产中可能需要先查询再更新
	userID, _ := ctx.Get("user_id")
	username, _ := ctx.Get("username")

	field := &model.MdModelField{
		ID:          id,
		ColumnTitle: req.ColumnTitle,
		Func:        req.Func,
		AggFunc:     req.AggFunc,
		ShowTitle:   req.ShowTitle,
		ShowWidth:   req.ShowWidth,
		UpdateID:    userID.(string),
		UpdateBy:    username.(string),
	}

	if err := h.modelService.UpdateField(field); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, field)
}

// DeleteModelField 删除模型字段
func (h *MdModelHandler) DeleteModelField(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("fieldId")
	if err := h.modelService.DeleteField(id); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, nil)
}
// SaveVisualModel 全量保存可视化构建模型
func (h *MdModelHandler) SaveVisualModel(c context.Context, ctx *app.RequestContext) {
	var req SaveVisualModelRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	tenantID, _ := ctx.Get("tenant_id")
	userID, _ := ctx.Get("user_id")
	username, _ := ctx.Get("username")

	serviceReq := &service.SaveVisualModelRequest{
		ModelID:      req.ModelID,
		ConnID:       req.ConnID,
		ModelName:    req.ModelName,
		ModelCode:    req.ModelCode,
		ModelVersion: req.ModelVersion,
		ModelKind:    req.ModelKind,
		IsPublic:     req.IsPublic,
		Remark:       req.Remark,
		Parameters:   req.Parameters,
		Tables:       req.Tables,
		Fields:       req.Fields,
		Joins:        req.Joins,
		JoinFields:   req.JoinFields,
		Wheres:       req.Wheres,
		Orders:       req.Orders,
		Groups:       req.Groups,
		Havings:      req.Havings,
		TenantID:     strconv.FormatUint(uint64(tenantID.(uint)), 10),
		UserID:       userID.(string),
		Username:     username.(string),
	}

	res, err := h.modelService.SaveVisualModel(serviceReq)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, res)
}
