package api

import (
	"context"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"

	"github.com/cloudwego/hertz/pkg/app"
)

// SsoRoleHandler 角色处理器结构体
type SsoRoleHandler struct {
	roleService service.SsoRoleService
}

// NewSsoRoleHandler 创建角色处理器实例
func NewSsoRoleHandler(roleService service.SsoRoleService) *SsoRoleHandler {
	return &SsoRoleHandler{roleService: roleService}
}

// SsoCreateRoleRequest 创建角色请求结构
type SsoCreateRoleRequest struct {
	ParentID  string `json:"parent_id" form:"parent_id"`
	AppCode   string `json:"app_code" form:"application_code"`
	OrgID     string `json:"org_id" form:"org_id"`
	KindCode  string `json:"kind_code" form:"kind_code" binding:"required"`
	RoleName  string `json:"role_name" form:"role_name" binding:"required"`
	RoleCode  string `json:"role_code" form:"role_code" binding:"required"`
	Status    int    `json:"status" form:"status"`
	DataRange string `json:"data_range" form:"data_range"`
	DataScope string `json:"data_scope" form:"data_scope"`
	Remark    string `json:"remark" form:"remark"`
	Sort      int    `json:"sort" form:"sort"`
}

// SsoUpdateRoleRequest 更新角色请求结构
type SsoUpdateRoleRequest struct {
	ParentID  string `json:"parent_id" form:"parent_id"`
	AppCode   string `json:"app_code" form:"application_code"`
	OrgID     string `json:"org_id" form:"org_id"`
	KindCode  string `json:"kind_code" form:"kind_code"`
	RoleName  string `json:"role_name" form:"role_name"`
	RoleCode  string `json:"role_code" form:"role_code"`
	Status    int    `json:"status" form:"status"`
	DataRange string `json:"data_range" form:"data_range"`
	DataScope string `json:"data_scope" form:"data_scope"`
	Remark    string `json:"remark" form:"remark"`
	Sort      int    `json:"sort" form:"sort"`
}

// CreateRole 创建角色
func (h *SsoRoleHandler) CreateRole(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateRoleRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 创建角色模型
	role := &model.SsoRole{
		ParentID:  req.ParentID,
		AppCode:   req.AppCode,
		OrgID:     req.OrgID,
		KindCode:  req.KindCode,
		RoleName:  req.RoleName,
		RoleCode:  req.RoleCode,
		Status:    req.Status,
		DataRange: req.DataRange,
		DataScope: req.DataScope,
		Remark:    req.Remark,
		Sort:      req.Sort,
	}

	// 调用服务层创建角色
	if err := h.roleService.CreateRole(role); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(201, role)
}

// GetRoleByID 根据ID获取角色
func (h *SsoRoleHandler) GetRoleByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层获取角色
	role, err := h.roleService.GetRoleByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "角色不存在"})
		return
	}

	ctx.JSON(200, role)
}

// UpdateRole 更新角色
func (h *SsoRoleHandler) UpdateRole(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req SsoUpdateRoleRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 调用服务层获取角色
	role, err := h.roleService.GetRoleByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "角色不存在"})
		return
	}

	// 更新角色字段
	if req.ParentID != "" {
		role.ParentID = req.ParentID
	}
	if req.AppCode != "" {
		role.AppCode = req.AppCode
	}
	if req.OrgID != "" {
		role.OrgID = req.OrgID
	}
	if req.KindCode != "" {
		role.KindCode = req.KindCode
	}
	if req.RoleName != "" {
		role.RoleName = req.RoleName
	}
	if req.RoleCode != "" {
		role.RoleCode = req.RoleCode
	}
	if req.Status != 0 {
		role.Status = req.Status
	}
	if req.DataRange != "" {
		role.DataRange = req.DataRange
	}
	if req.DataScope != "" {
		role.DataScope = req.DataScope
	}
	if req.Remark != "" {
		role.Remark = req.Remark
	}
	if req.Sort != 0 {
		role.Sort = req.Sort
	}

	// 调用服务层更新角色
	if err := h.roleService.UpdateRole(role); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, role)
}

// DeleteRole 删除角色
func (h *SsoRoleHandler) DeleteRole(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层删除角色
	if err := h.roleService.DeleteRole(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, map[string]string{"message": "角色删除成功"})
}

// GetAllRoles 获取所有角色
func (h *SsoRoleHandler) GetAllRoles(c context.Context, ctx *app.RequestContext) {
	// 调用服务层获取所有角色
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取角色列表失败"})
		return
	}

	ctx.JSON(200, roles)
}
