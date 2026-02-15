package api

import (
	"context"
	"encoding/json"
	"fmt"
	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/module/user/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"

	userModel "metadata-platform/internal/module/user/model"
)

// SsoTenantHandler 租户处理器结构体
type SsoTenantHandler struct {
	*utils.BaseHandler
	tenantService service.SsoTenantService
	audit         AuditService
}

// NewSsoTenantHandler 创建租户处理器实例
func NewSsoTenantHandler(tenantService service.SsoTenantService, auditQueue *queue.AuditLogQueue) *SsoTenantHandler {
	return &SsoTenantHandler{
		BaseHandler:   utils.NewBaseHandler(),
		tenantService: tenantService,
		audit:         &auditServiceImpl{queue: auditQueue},
	}
}

// SsoCreateTenantRequest 创建租户请求结构
type SsoCreateTenantRequest struct {
	TenantName string `json:"tenant_name" form:"tenant_name" binding:"required"`
	TenantCode string `json:"tenant_code" form:"tenant_code" binding:"required"`
	Status     int    `json:"status" form:"status"`
	Remark     string `json:"remark" form:"remark"`
}

// SsoUpdateTenantRequest 更新租户请求结构
type SsoUpdateTenantRequest struct {
	TenantName *string `json:"tenant_name" form:"tenant_name"`
	TenantCode *string `json:"tenant_code" form:"tenant_code"`
	Status     *int    `json:"status" form:"status"`
	Remark     *string `json:"remark" form:"remark"`
}

// CreateTenant 创建租户
func (h *SsoTenantHandler) CreateTenant(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateTenantRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 创建租户模型
	tenant := &userModel.SsoTenant{
		TenantName: req.TenantName,
		TenantCode: req.TenantCode,
		Status:     req.Status,
		Remark:     req.Remark,
		CreateID:   headerUser.UserID,
		CreateBy:   headerUser.UserAccount,
		TenantID:   headerUser.TenantID,
	}

	// 调用服务层创建租户
	if err := h.tenantService.CreateTenant(tenant); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	afterData, _ := json.Marshal(tenant)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:        utils.GetSnowflake().GenerateIDString(),
		TraceID:   ctx.GetString("trace_id"),
		ModelID:   "tenant",
		RecordID:  tenant.ID,
		Action:    "CREATE",
		AfterData: string(afterData),
		CreateBy:  headerUser.UserAccount,
		Source:    "tenant_service",
	})

	ctx.JSON(201, tenant)
}

// GetTenantByID 根据ID获取租户
func (h *SsoTenantHandler) GetTenantByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层获取租户
	tenant, err := h.tenantService.GetTenantByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "租户不存在"})
		return
	}

	ctx.JSON(200, tenant)
}

// UpdateTenant 更新租户
func (h *SsoTenantHandler) UpdateTenant(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req SsoUpdateTenantRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 获取变更前数据
	beforeData, err := h.tenantService.GetTenantByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "租户不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	// 更新租户字段
	if req.TenantName != nil {
		beforeData.TenantName = *req.TenantName
	}
	if req.TenantCode != nil {
		beforeData.TenantCode = *req.TenantCode
	}
	if req.Status != nil {
		beforeData.Status = *req.Status
	}
	if req.Remark != nil {
		beforeData.Remark = *req.Remark
	}

	// 设置更新人信息
	beforeData.UpdateID = headerUser.UserID
	beforeData.UpdateBy = headerUser.UserAccount

	// 调用服务层更新租户
	if err := h.tenantService.UpdateTenant(beforeData); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	afterJSON, _ := json.Marshal(beforeData)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    ctx.GetString("trace_id"),
		ModelID:    "tenant",
		RecordID:   id,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "tenant_service",
	})

	ctx.JSON(200, beforeData)
}

// DeleteTenant 删除租户
func (h *SsoTenantHandler) DeleteTenant(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 获取删除前数据
	beforeData, err := h.tenantService.GetTenantByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "租户不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	// 调用服务层删除租户
	if err := h.tenantService.DeleteTenant(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    ctx.GetString("trace_id"),
		ModelID:    "tenant",
		RecordID:   id,
		Action:     "DELETE",
		BeforeData: string(beforeJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "tenant_service",
	})

	ctx.JSON(200, map[string]string{"message": "租户删除成功"})
}

// GetAllTenants 获取所有租户
func (h *SsoTenantHandler) GetAllTenants(c context.Context, ctx *app.RequestContext) {
	// 调用服务层获取所有租户
	tenants, err := h.tenantService.GetAllTenants()
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取租户列表失败"})
		return
	}

	ctx.JSON(200, tenants)
}
