package api

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"
)

// MdConnHandler 数据连接API处理器
type MdConnHandler struct {
	connService service.MdConnService
}

// NewMdConnHandler 创建数据连接API处理器实例
func NewMdConnHandler(connService service.MdConnService) *MdConnHandler {
	return &MdConnHandler{connService: connService}
}

// CreateRequest 创建数据连接请求
type CreateConnRequest struct {
	ParentID     string `json:"parent_id"`
	ConnName     string `json:"conn_name" binding:"required"`
	ConnKind     string `json:"conn_kind" binding:"required"`
	ConnVersion  string `json:"conn_version"`
	ConnHost     string `json:"conn_host" binding:"required"`
	ConnPort     int    `json:"conn_port" binding:"required"`
	ConnUser     string `json:"conn_user" binding:"required"`
	ConnPassword string `json:"conn_password" binding:"required"`
	ConnDatabase string `json:"conn_database" binding:"required"`
	Remark       string `json:"remark"`
}

// UpdateConnRequest 更新数据连接请求
type UpdateConnRequest struct {
	ConnName     string `json:"conn_name"`
	ConnHost     string `json:"conn_host"`
	ConnPort     int    `json:"conn_port"`
	ConnUser     string `json:"conn_user"`
	ConnPassword string `json:"conn_password"`
	ConnDatabase string `json:"conn_database"`
	Remark       string `json:"remark"`
}

// CreateConn 创建数据连接
func (h *MdConnHandler) CreateConn(c context.Context, ctx *app.RequestContext) {
	var req CreateConnRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 从上下文获取租户ID
	tenantID, _ := ctx.Get("tenant_id")
	// 从上下文获取用户ID和用户名
	userID, _ := ctx.Get("user_id")
	username, _ := ctx.Get("username")

	conn := &model.MdConn{
		ParentID:     req.ParentID,
		ConnName:     req.ConnName,
		ConnKind:     req.ConnKind,
		ConnVersion:  req.ConnVersion,
		ConnHost:     req.ConnHost,
		ConnPort:     req.ConnPort,
		ConnUser:     req.ConnUser,
		ConnPassword: req.ConnPassword,
		ConnDatabase: req.ConnDatabase,
		TenantID:     strconv.FormatUint(uint64(tenantID.(uint)), 10),
		CreateID:     strconv.FormatInt(userID.(int64), 10),
		CreateBy:     username.(string),
		UpdateID:     strconv.FormatInt(userID.(int64), 10),
		UpdateBy:     username.(string),
	}

	// 密码加密处理（实际应在Service层处理，这里假设Service CreateConn前已处理或内部处理）
	// 目前简单处理，直接传递

	if err := h.connService.CreateConn(conn); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, conn)
}

// GetConnByID 获取数据连接
func (h *MdConnHandler) GetConnByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	
	conn, err := h.connService.GetConnByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "数据连接不存在")
		return
	}

	// 密码脱敏
	conn.ConnPassword = "******"

	utils.SuccessResponse(ctx, conn)
}

// UpdateConn 更新数据连接
func (h *MdConnHandler) UpdateConn(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req UpdateConnRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	conn, err := h.connService.GetConnByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "数据连接不存在")
		return
	}

	// 更新字段
	if req.ConnName != "" { conn.ConnName = req.ConnName }
	if req.ConnHost != "" { conn.ConnHost = req.ConnHost }
	if req.ConnPort != 0 { conn.ConnPort = req.ConnPort }
	if req.ConnUser != "" { conn.ConnUser = req.ConnUser }
	if req.ConnPassword != "" { conn.ConnPassword = req.ConnPassword }
	if req.ConnDatabase != "" { conn.ConnDatabase = req.ConnDatabase }

	// 从上下文获取用户ID和用户名
	userID, _ := ctx.Get("user_id")
	username, _ := ctx.Get("username")
	conn.UpdateID = strconv.FormatInt(userID.(int64), 10)
	conn.UpdateBy = username.(string)

	if err := h.connService.UpdateConn(conn); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, conn)
}

// DeleteConn 删除数据连接
func (h *MdConnHandler) DeleteConn(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	if err := h.connService.DeleteConn(id); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, nil)
}

// GetAllConns 获取所有数据连接
func (h *MdConnHandler) GetAllConns(c context.Context, ctx *app.RequestContext) {
	// 从上下文获取租户ID
	tenantID, _ := ctx.Get("tenant_id")
	
	conns, err := h.connService.GetAllConns(strconv.FormatUint(uint64(tenantID.(uint)), 10))
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	// 密码脱敏
	for i := range conns {
		conns[i].ConnPassword = "******"
	}

	utils.SuccessResponse(ctx, conns)
}

// GetConnsByParentID 根据父ID获取数据连接
func (h *MdConnHandler) GetConnsByParentID(c context.Context, ctx *app.RequestContext) {
	parentID := ctx.Param("parent_id")

	conns, err := h.connService.GetConnsByParentID(parentID)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	// 密码脱敏
	for i := range conns {
		conns[i].ConnPassword = "******"
	}

	utils.SuccessResponse(ctx, conns)
}

// TestConnection 测试数据连接
func (h *MdConnHandler) TestConnection(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	
	// 先获取连接信息
	conn, err := h.connService.GetConnByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "数据连接不存在")
		return
	}

	// 测试连接
	if err := h.connService.TestConnection(conn); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, "连接失败: "+err.Error())
		return
	}

	utils.SuccessResponse(ctx, "连接成功")
}

// GetTables 获取数据源表列表
func (h *MdConnHandler) GetTables(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	schema := ctx.Query("schema")

	conn, err := h.connService.GetConnByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "数据连接不存在")
		return
	}

	tables, err := h.connService.GetTables(conn, schema)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, tables)
}

// GetViews 获取数据源视图列表
func (h *MdConnHandler) GetViews(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	schema := ctx.Query("schema")

	conn, err := h.connService.GetConnByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "数据连接不存在")
		return
	}

	views, err := h.connService.GetViews(conn, schema)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, views)
}

// GetTableStructure 获取数据源表结构
func (h *MdConnHandler) GetTableStructure(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	table := ctx.Param("table")
	schema := ctx.Query("schema")

	conn, err := h.connService.GetConnByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "数据连接不存在")
		return
	}

	columns, err := h.connService.GetTableStructure(conn, schema, table)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, columns)
}

// PreviewTableData 预览数据源表数据
func (h *MdConnHandler) PreviewTableData(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	table := ctx.Param("table")
	schema := ctx.Query("schema")
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	conn, err := h.connService.GetConnByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "数据连接不存在")
		return
	}

	data, err := h.connService.PreviewTableData(conn, schema, table, limit)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, data)
}
