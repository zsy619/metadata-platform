package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"

	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/queue"
	userModel "metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"
	"metadata-platform/internal/utils"
)

// SsoRoleHandler 角色处理器结构体
type SsoRoleHandler struct {
	*utils.BaseHandler
	roleService service.SsoRoleService
	audit       AuditService
}

// NewSsoRoleHandler 创建角色处理器实例
func NewSsoRoleHandler(roleService service.SsoRoleService, auditQueue *queue.AuditLogQueue) *SsoRoleHandler {
	return &SsoRoleHandler{
		BaseHandler: utils.NewBaseHandler(),
		roleService: roleService,
		audit:       &auditServiceImpl{queue: auditQueue},
	}
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
	ParentID  *string `json:"parent_id" form:"parent_id"`
	AppCode   *string `json:"app_code" form:"application_code"`
	OrgID     *string `json:"org_id" form:"org_id"`
	KindCode  *string `json:"kind_code" form:"kind_code"`
	RoleName  *string `json:"role_name" form:"role_name"`
	RoleCode  *string `json:"role_code" form:"role_code"`
	Status    *int    `json:"status" form:"status"`
	DataRange *string `json:"data_range" form:"data_range"`
	DataScope *string `json:"data_scope" form:"data_scope"`
	Remark    *string `json:"remark" form:"remark"`
	Sort      *int    `json:"sort" form:"sort"`
}

// CreateRole 创建角色
func (h *SsoRoleHandler) CreateRole(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateRoleRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 创建角色模型
	role := &userModel.SsoRole{
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
		CreateID:  headerUser.UserID,
		CreateBy:  headerUser.UserAccount,
		TenantID:  headerUser.TenantID,
	}

	// 调用服务层创建角色
	if err := h.roleService.CreateRole(role); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	afterData, _ := json.Marshal(role)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:        utils.GetSnowflake().GenerateIDString(),
		TraceID:   headerUser.TraceID,
		ModelID:   "role",
		RecordID:  role.ID,
		Action:    "CREATE",
		AfterData: string(afterData),
		CreateBy:  headerUser.UserAccount,
		Source:    "role_service",
	})

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

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 获取变更前数据
	beforeData, err := h.roleService.GetRoleByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "角色不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	// 更新角色字段
	if req.ParentID != nil {
		beforeData.ParentID = *req.ParentID
	}
	if req.AppCode != nil {
		beforeData.AppCode = *req.AppCode
	}
	if req.OrgID != nil {
		beforeData.OrgID = *req.OrgID
	}
	if req.KindCode != nil {
		beforeData.KindCode = *req.KindCode
	}
	if req.RoleName != nil {
		beforeData.RoleName = *req.RoleName
	}
	if req.RoleCode != nil {
		beforeData.RoleCode = *req.RoleCode
	}
	if req.Status != nil {
		beforeData.Status = *req.Status
	}
	if req.DataRange != nil {
		beforeData.DataRange = *req.DataRange
	}
	if req.DataScope != nil {
		beforeData.DataScope = *req.DataScope
	}
	if req.Remark != nil {
		beforeData.Remark = *req.Remark
	}
	if req.Sort != nil {
		beforeData.Sort = *req.Sort
	}

	// 设置更新人信息
	beforeData.UpdateID = headerUser.UserID
	beforeData.UpdateBy = headerUser.UserAccount

	// 调用服务层更新角色
	if err := h.roleService.UpdateRole(beforeData); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	afterJSON, _ := json.Marshal(beforeData)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "role",
		RecordID:   id,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "role_service",
	})

	ctx.JSON(200, beforeData)
}

// DeleteRole 删除角色
func (h *SsoRoleHandler) DeleteRole(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 获取删除前数据
	beforeData, err := h.roleService.GetRoleByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "角色不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	// 调用服务层删除角色
	if err := h.roleService.DeleteRole(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "role",
		RecordID:   id,
		Action:     "DELETE",
		BeforeData: string(beforeJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "role_service",
	})

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

// GetRoleMenus 获取角色的菜单ID列表
func (h *SsoRoleHandler) GetRoleMenus(c context.Context, ctx *app.RequestContext) {
	roleID := ctx.Param("id")

	// 调用服务层获取角色菜单
	menuIDs, err := h.roleService.GetRoleMenus(roleID)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"role_id":  roleID,
		"menu_ids": menuIDs,
	})
}

// UpdateRoleMenusRequest 更新角色菜单请求结构
type UpdateRoleMenusRequest struct {
	MenuIDs []string `json:"menu_ids" binding:"required"`
}

// UpdateRoleMenus 更新角色的菜单关联
func (h *SsoRoleHandler) UpdateRoleMenus(c context.Context, ctx *app.RequestContext) {
	roleID := ctx.Param("id")
	var req UpdateRoleMenusRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取当前用户信息
	headerUser := h.GetHeaderUserStruct(c, ctx)

	// 获取更新前的菜单ID列表
	beforeMenuIDs, _ := h.roleService.GetRoleMenus(roleID)
	beforeJSON, _ := json.Marshal(beforeMenuIDs)

	// 调用服务层更新角色菜单
	if err := h.roleService.UpdateRoleMenus(roleID, req.MenuIDs, headerUser.UserAccount); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	afterJSON, _ := json.Marshal(req.MenuIDs)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "role_menu",
		RecordID:   roleID,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "role_service",
	})

	ctx.JSON(200, map[string]interface{}{
		"role_id":  roleID,
		"menu_ids": req.MenuIDs,
		"message":  "角色菜单更新成功",
	})
}
