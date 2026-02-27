package api

import (
	"context"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// MdModelProcedureHandler 存储过程/函数API处理器
type MdModelProcedureHandler struct {
	*utils.BaseHandler
	procService service.MdModelProcedureService
}

// NewMdModelProcedureHandler 创建存储过程/函数API处理器实例
func NewMdModelProcedureHandler(procService service.MdModelProcedureService) *MdModelProcedureHandler {
	return &MdModelProcedureHandler{
		BaseHandler: utils.NewBaseHandler(),
		procService: procService,
	}
}

// CreateProcedureRequest 创建存储过程/函数请求
type CreateProcedureRequest struct {
	ConnID      string `json:"conn_id" binding:"required"`
	ConnName    string `json:"conn_name"`
	ProcSchema  string `json:"proc_schema"`
	ProcName    string `json:"proc_name" binding:"required"`
	ProcTitle   string `json:"proc_title"`
	ProcType    string `json:"proc_type" binding:"required"`
	ProcComment string `json:"proc_comment"`
	Definition  string `json:"definition"`
	ReturnType  string `json:"return_type"`
	Language    string `json:"language"`
}

// UpdateProcedureRequest 更新存储过程/函数请求
type UpdateProcedureRequest struct {
	ConnID      string `json:"conn_id"`
	ConnName    string `json:"conn_name"`
	ProcSchema  string `json:"proc_schema"`
	ProcName    string `json:"proc_name"`
	ProcTitle   string `json:"proc_title"`
	ProcType    string `json:"proc_type"`
	ProcComment string `json:"proc_comment"`
	Definition  string `json:"definition"`
	ReturnType  string `json:"return_type"`
	Language    string `json:"language"`
}

// ImportProceduresRequest 导入存储过程/函数请求
type ImportProceduresRequest struct {
	ConnID     string                   `json:"conn_id" binding:"required"`
	ProcSchema string                   `json:"proc_schema"`
	Procedures []map[string]interface{} `json:"procedures" binding:"required"`
}

// CreateProcedure 创建存储过程/函数
func (h *MdModelProcedureHandler) CreateProcedure(c context.Context, ctx *app.RequestContext) {
	var req CreateProcedureRequest
	if err := ctx.BindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	tenantID, _ := ctx.Get("tenant_id")
	userID, _ := ctx.Get("user_id")
	username, _ := ctx.Get("username")

	proc := &model.MdModelProcedure{
		ConnID:      req.ConnID,
		ConnName:    req.ConnName,
		ProcSchema:  req.ProcSchema,
		ProcName:    req.ProcName,
		ProcTitle:   req.ProcTitle,
		ProcType:    req.ProcType,
		ProcComment: req.ProcComment,
		Definition:  req.Definition,
		ReturnType:  req.ReturnType,
		Language:    req.Language,
		TenantID:    strconv.FormatUint(uint64(tenantID.(uint)), 10),
		CreateID:    userID.(string),
		CreateBy:    username.(string),
		UpdateID:    userID.(string),
		UpdateBy:    username.(string),
	}

	if err := h.procService.CreateProcedure(proc); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, proc)
}

// GetProcedureByID 根据ID获取存储过程/函数
func (h *MdModelProcedureHandler) GetProcedureByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	proc, err := h.procService.GetProcedureByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "存储过程/函数不存在")
		return
	}
	utils.SuccessResponse(ctx, proc)
}

// UpdateProcedure 更新存储过程/函数
func (h *MdModelProcedureHandler) UpdateProcedure(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req UpdateProcedureRequest
	if err := ctx.BindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	proc, err := h.procService.GetProcedureByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusNotFound, "存储过程/函数不存在")
		return
	}

	if req.ConnID != "" {
		proc.ConnID = req.ConnID
	}
	if req.ConnName != "" {
		proc.ConnName = req.ConnName
	}
	if req.ProcSchema != "" {
		proc.ProcSchema = req.ProcSchema
	}
	if req.ProcName != "" {
		proc.ProcName = req.ProcName
	}
	if req.ProcTitle != "" {
		proc.ProcTitle = req.ProcTitle
	}
	if req.ProcType != "" {
		proc.ProcType = req.ProcType
	}
	if req.ProcComment != "" {
		proc.ProcComment = req.ProcComment
	}
	if req.Definition != "" {
		proc.Definition = req.Definition
	}
	if req.ReturnType != "" {
		proc.ReturnType = req.ReturnType
	}
	if req.Language != "" {
		proc.Language = req.Language
	}

	userID, _ := ctx.Get("user_id")
	username, _ := ctx.Get("username")
	proc.UpdateID = userID.(string)
	proc.UpdateBy = username.(string)

	if err := h.procService.UpdateProcedure(proc); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(ctx, proc)
}

// DeleteProcedure 删除存储过程/函数
func (h *MdModelProcedureHandler) DeleteProcedure(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	if err := h.procService.DeleteProcedure(id); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(ctx, "删除成功")
}

// GetProceduresByConnID 根据连接ID获取存储过程/函数列表
func (h *MdModelProcedureHandler) GetProceduresByConnID(c context.Context, ctx *app.RequestContext) {
	connID := ctx.Param("conn_id")
	procedures, err := h.procService.GetProceduresByConnID(connID)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, "获取存储过程/函数列表失败")
		return
	}
	utils.SuccessResponse(ctx, procedures)
}

// GetAllProcedures 获取所有存储过程/函数
func (h *MdModelProcedureHandler) GetAllProcedures(c context.Context, ctx *app.RequestContext) {
	tenantID, _ := ctx.Get("tenant_id")
	procedures, err := h.procService.GetAllProcedures(strconv.FormatUint(uint64(tenantID.(uint)), 10))
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, "获取存储过程/函数列表失败")
		return
	}
	utils.SuccessResponse(ctx, procedures)
}

// GetParamsByProcID 根据存储过程ID获取参数列表
func (h *MdModelProcedureHandler) GetParamsByProcID(c context.Context, ctx *app.RequestContext) {
	procID := ctx.Param("id")
	params, err := h.procService.GetParamsByProcID(procID)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, "获取参数列表失败")
		return
	}
	utils.SuccessResponse(ctx, params)
}
