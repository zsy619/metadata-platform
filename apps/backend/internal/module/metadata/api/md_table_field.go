package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"
)

// MdTableFieldHandler 数据连接表字段处理器结构体
type MdTableFieldHandler struct {
	fieldService service.MdTableFieldService
}

// NewMdTableFieldHandler 创建数据连接表字段处理器实例
func NewMdTableFieldHandler(fieldService service.MdTableFieldService) *MdTableFieldHandler {
	return &MdTableFieldHandler{fieldService: fieldService}
}

// CreateFieldRequest 创建数据连接表字段请求结构
type CreateFieldRequest struct {
	TenantID        string `json:"tenant_id" binding:"required"`
	ConnID          string `json:"conn_id" binding:"required"`
	TableID         string `json:"table_id" binding:"required"`
	TableTitle      string `json:"table_title"`
	ColumnName      string `json:"column_name" binding:"required"`
	ColumnTitle     string `json:"column_title"`
	ColumnType      string `json:"column_type" binding:"required"`
	ColumnLength    int    `json:"column_length"`
	ColumnComment   string `json:"column_comment"`
	IsNullable      bool   `json:"is_nullable"`
	IsPrimaryKey    bool   `json:"is_primary_key"`
	IsAutoIncrement bool   `json:"is_auto_increment"`
	DefaultValue    string `json:"default_value"`
	ExtraInfo       string `json:"extra_info"`
}

// UpdateFieldRequest 更新数据连接表字段请求结构
type UpdateFieldRequest struct {
	ConnID          string `json:"conn_id"`
	TableID         string `json:"table_id"`
	TableTitle      string `json:"table_title"`
	ColumnName      string `json:"column_name"`
	ColumnTitle     string `json:"column_title"`
	ColumnType      string `json:"column_type"`
	ColumnLength    int    `json:"column_length"`
	ColumnComment   string `json:"column_comment"`
	IsNullable      bool   `json:"is_nullable"`
	IsPrimaryKey    bool   `json:"is_primary_key"`
	IsAutoIncrement bool   `json:"is_auto_increment"`
	DefaultValue    string `json:"default_value"`
	ExtraInfo       string `json:"extra_info"`
}

// CreateField 创建数据连接表字段
func (h *MdTableFieldHandler) CreateField(c context.Context, ctx *app.RequestContext) {
	var req CreateFieldRequest
	if err := ctx.BindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	// 创建数据连接表字段模型
	field := &model.MdTableField{
		TenantID:        req.TenantID,
		ConnID:          req.ConnID,
		TableID:         req.TableID,
		TableTitle:      req.TableTitle,
		ColumnName:      req.ColumnName,
		ColumnTitle:     req.ColumnTitle,
		ColumnType:      req.ColumnType,
		ColumnLength:    req.ColumnLength,
		ColumnComment:   req.ColumnComment,
		IsNullable:      req.IsNullable,
		IsPrimaryKey:    req.IsPrimaryKey,
		IsAutoIncrement: req.IsAutoIncrement,
		DefaultValue:    req.DefaultValue,
		ExtraInfo:       req.ExtraInfo,
	}

	// 调用服务层创建数据连接表字段
	if err := h.fieldService.CreateField(field); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, field)
}

// GetFieldByID 根据ID获取数据连接表字段
func (h *MdTableFieldHandler) GetFieldByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层获取数据连接表字段
	field, err := h.fieldService.GetFieldByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "字段不存在")
		return
	}

	utils.SuccessResponse(ctx, field)
}

// GetFieldsByTableID 根据表ID获取数据连接表字段列表
func (h *MdTableFieldHandler) GetFieldsByTableID(c context.Context, ctx *app.RequestContext) {
	tableID := ctx.Param("table_id")

	// 调用服务层根据表ID获取数据连接表字段列表
	fields, err := h.fieldService.GetFieldsByTableID(tableID)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, "获取字段列表失败")
		return
	}

	utils.SuccessResponse(ctx, fields)
}

// UpdateField 更新数据连接表字段
func (h *MdTableFieldHandler) UpdateField(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req UpdateFieldRequest
	if err := ctx.BindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	// 调用服务层获取数据连接表字段
	field, err := h.fieldService.GetFieldByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "字段不存在"})
		return
	}

	// 更新数据连接表字段
	if req.ConnID != "" {
		field.ConnID = req.ConnID
	}
	if req.TableID != "" {
		field.TableID = req.TableID
	}
	if req.TableTitle != "" {
		field.TableTitle = req.TableTitle
	}
	if req.ColumnName != "" {
		field.ColumnName = req.ColumnName
	}
	if req.ColumnTitle != "" {
		field.ColumnTitle = req.ColumnTitle
	}
	if req.ColumnType != "" {
		field.ColumnType = req.ColumnType
	}
	if req.ColumnLength != 0 {
		field.ColumnLength = req.ColumnLength
	}
	if req.ColumnComment != "" {
		field.ColumnComment = req.ColumnComment
	}
	field.IsNullable = req.IsNullable
	field.IsPrimaryKey = req.IsPrimaryKey
	field.IsAutoIncrement = req.IsAutoIncrement
	if req.DefaultValue != "" {
		field.DefaultValue = req.DefaultValue
	}
	if req.ExtraInfo != "" {
		field.ExtraInfo = req.ExtraInfo
	}

	// 调用服务层更新数据连接表字段
	if err := h.fieldService.UpdateField(field); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, field)
}

// DeleteField 删除数据连接表字段
func (h *MdTableFieldHandler) DeleteField(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层删除数据连接表字段
	if err := h.fieldService.DeleteField(id); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "字段删除成功")
}

// DeleteFieldsByTableID 根据表ID删除所有数据连接表字段
func (h *MdTableFieldHandler) DeleteFieldsByTableID(c context.Context, ctx *app.RequestContext) {
	tableID := ctx.Param("table_id")

	// 调用服务层根据表ID删除所有数据连接表字段
	if err := h.fieldService.DeleteFieldsByTableID(tableID); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "表字段批量删除成功")
}

// GetAllFields 获取所有数据连接表字段
func (h *MdTableFieldHandler) GetAllFields(c context.Context, ctx *app.RequestContext) {
	tenantID := ctx.Query("tenant_id")

	// 调用服务层获取所有数据连接表字段
	fields, err := h.fieldService.GetAllFields(tenantID)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, "获取字段列表失败")
		return
	}

	utils.SuccessResponse(ctx, fields)
}
