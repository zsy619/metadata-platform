package api

import (
	"context"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"

	"github.com/cloudwego/hertz/pkg/app"
)

// SsoTenantHandler 租户处理器结构体
type SsoTenantHandler struct {
	tenantService service.SsoTenantService
}

// NewSsoTenantHandler 创建租户处理器实例
func NewSsoTenantHandler(tenantService service.SsoTenantService) *SsoTenantHandler {
	return &SsoTenantHandler{tenantService: tenantService}
}

// SsoCreateTenantRequest 创建租户请求结构
type SsoCreateTenantRequest struct {
	TenantName string `json:"tenant_name" form:"tenant_name" binding:"required"`
	TenantCode string `json:"tenant_code" form:"tenant_code" binding:"required"`
	Linkman    string `json:"linkman" form:"linkman"`
	Contact    string `json:"contact" form:"contact"`
	Address    string `json:"address" form:"address"`
	State      int    `json:"state" form:"state"`
	Remark     string `json:"remark" form:"remark"`
}

// SsoUpdateTenantRequest 更新租户请求结构
type SsoUpdateTenantRequest struct {
	TenantName string `json:"tenant_name" form:"tenant_name"`
	TenantCode string `json:"tenant_code" form:"tenant_code"`
	Linkman    string `json:"linkman" form:"linkman"`
	Contact    string `json:"contact" form:"contact"`
	Address    string `json:"address" form:"address"`
	State      int    `json:"state" form:"state"`
	Remark     string `json:"remark" form:"remark"`
}

// CreateTenant 创建租户
func (h *SsoTenantHandler) CreateTenant(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateTenantRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 创建租户模型
	tenant := &model.SsoTenant{
		TenantName: req.TenantName,
		TenantCode: req.TenantCode,
		Linkman:    req.Linkman,
		Contact:    req.Contact,
		Address:    req.Address,
		State:      req.State,
		Remark:     req.Remark,
	}

	// 调用服务层创建租户
	if err := h.tenantService.CreateTenant(tenant); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

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

	// 调用服务层获取租户
	tenant, err := h.tenantService.GetTenantByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "租户不存在"})
		return
	}

	// 更新租户字段
	if req.TenantName != "" {
		tenant.TenantName = req.TenantName
	}
	if req.TenantCode != "" {
		tenant.TenantCode = req.TenantCode
	}
	if req.Linkman != "" {
		tenant.Linkman = req.Linkman
	}
	if req.Contact != "" {
		tenant.Contact = req.Contact
	}
	if req.Address != "" {
		tenant.Address = req.Address
	}
	if req.State != 0 {
		tenant.State = req.State
	}
	if req.Remark != "" {
		tenant.Remark = req.Remark
	}

	// 调用服务层更新租户
	if err := h.tenantService.UpdateTenant(tenant); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, tenant)
}

// DeleteTenant 删除租户
func (h *SsoTenantHandler) DeleteTenant(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层删除租户
	if err := h.tenantService.DeleteTenant(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

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
