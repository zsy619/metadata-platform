package api

import (
	"context"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// MdTableHandler 数据连接表处理器结构体
type MdTableHandler struct {
	tableService service.MdTableService
}

// NewMdTableHandler 创建数据连接表处理器实例
func NewMdTableHandler(tableService service.MdTableService) *MdTableHandler {
	return &MdTableHandler{tableService: tableService}
}

// CreateTableRequest 创建数据连接表请求结构
type CreateTableRequest struct {
	TenantID     string `json:"tenant_id" binding:"required"`
	ConnID       string `json:"conn_id" binding:"required"`
	ConnName     string `json:"conn_name" binding:"required"`
	TableSchema  string `json:"table_schema"`
	TableName    string `json:"table_name" binding:"required"`
	TableTitle   string `json:"table_title"`
	TableType    string `json:"table_type"`
	TableComment string `json:"table_comment"`
}

// UpdateTableRequest 更新数据连接表请求结构
type UpdateTableRequest struct {
	ConnID       string `json:"conn_id"`
	ConnName     string `json:"conn_name"`
	TableSchema  string `json:"table_schema"`
	TableName    string `json:"table_name"`
	TableTitle   string `json:"table_title"`
	TableType    string `json:"table_type"`
	TableComment string `json:"table_comment"`
}

// CreateTable 创建数据连接表
func (h *MdTableHandler) CreateTable(c context.Context, ctx *app.RequestContext) {
	var req CreateTableRequest
	if err := ctx.BindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	// 创建数据连接表模型
	table := &model.MdTable{
		TenantID:     req.TenantID,
		ConnID:       req.ConnID,
		ConnName:     req.ConnName,
		TableSchema:  req.TableSchema,
		TableNameStr: req.TableName,
		TableTitle:   req.TableTitle,
		TableType:    req.TableType,
		TableComment: req.TableComment,
	}

	// 调用服务层创建数据连接表
	if err := h.tableService.CreateTable(table); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, table)
}

// GetTableByID 根据ID获取数据连接表
func (h *MdTableHandler) GetTableByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层获取数据连接表
	table, err := h.tableService.GetTableByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "表不存在")
		return
	}

	utils.SuccessResponse(ctx, table)
}

// UpdateTable 更新数据连接表
func (h *MdTableHandler) UpdateTable(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req UpdateTableRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 调用服务层获取数据连接表
	table, err := h.tableService.GetTableByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "表不存在"})
		return
	}

	// 更新数据连接表字段
	if req.ConnID != "" {
		table.ConnID = req.ConnID
	}
	if req.ConnName != "" {
		table.ConnName = req.ConnName
	}
	if req.TableSchema != "" {
		table.TableSchema = req.TableSchema
	}
	if req.TableName != "" {
		table.TableNameStr = req.TableName
	}
	if req.TableTitle != "" {
		table.TableTitle = req.TableTitle
	}
	if req.TableType != "" {
		table.TableType = req.TableType
	}
	if req.TableComment != "" {
		table.TableComment = req.TableComment
	}

	// 调用服务层更新数据连接表
	if err := h.tableService.UpdateTable(table); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, table)
}

// DeleteTable 删除数据连接表
func (h *MdTableHandler) DeleteTable(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层删除数据连接表
	if err := h.tableService.DeleteTable(id); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, "表删除成功")
}

// GetTablesByConnID 根据连接ID获取数据连接表
func (h *MdTableHandler) GetTablesByConnID(c context.Context, ctx *app.RequestContext) {
	connID := ctx.Param("conn_id")

	// 调用服务层根据连接ID获取数据连接表
	tables, err := h.tableService.GetTablesByConnID(connID)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, "获取表列表失败")
		return
	}

	utils.SuccessResponse(ctx, tables)
}

// GetAllTables 获取所有数据连接表
func (h *MdTableHandler) GetAllTables(c context.Context, ctx *app.RequestContext) {
	tenantID := ctx.Query("tenant_id")

	// 调用服务层获取所有数据连接表
	tables, err := h.tableService.GetAllTables(tenantID)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, "获取表列表失败")
		return
	}

	utils.SuccessResponse(ctx, tables)
}
