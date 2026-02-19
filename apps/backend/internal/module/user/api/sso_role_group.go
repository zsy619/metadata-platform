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

// SsoRoleGroupHandler 角色分组处理器结构体
type SsoRoleGroupHandler struct {
	*utils.BaseHandler
	roleGroupService service.SsoRoleGroupService
	audit            AuditService
}

// NewSsoRoleGroupHandler 创建角色分组处理器实例
func NewSsoRoleGroupHandler(roleGroupService service.SsoRoleGroupService, auditQueue *queue.AuditLogQueue) *SsoRoleGroupHandler {
	return &SsoRoleGroupHandler{
		BaseHandler:      utils.NewBaseHandler(),
		roleGroupService: roleGroupService,
		audit:            &auditServiceImpl{queue: auditQueue},
	}
}

// SsoCreateRoleGroupRequest 创建角色分组请求结构
type SsoCreateRoleGroupRequest struct {
	ParentID  string `json:"parent_id" form:"parent_id"`
	AppCode   string `json:"app_code" form:"app_code"`
	OrgID     string `json:"org_id" form:"org_id"`
	KindCode  string `json:"kind_code" form:"kind_code" binding:"required"`
	GroupName string `json:"group_name" form:"group_name" binding:"required"`
	GroupCode string `json:"group_code" form:"group_code" binding:"required"`
	Status    int    `json:"status" form:"status"`
	Remark    string `json:"remark" form:"remark"`
	Sort      int    `json:"sort" form:"sort"`
}

// SsoUpdateRoleGroupRequest 更新角色分组请求结构
type SsoUpdateRoleGroupRequest struct {
	ParentID  *string `json:"parent_id" form:"parent_id"`
	AppCode   *string `json:"app_code" form:"app_code"`
	OrgID     *string `json:"org_id" form:"org_id"`
	KindCode  *string `json:"kind_code" form:"kind_code"`
	GroupName *string `json:"group_name" form:"group_name"`
	GroupCode *string `json:"group_code" form:"group_code"`
	Status    *int    `json:"status" form:"status"`
	Remark    *string `json:"remark" form:"remark"`
	Sort      *int    `json:"sort" form:"sort"`
}

// CreateRoleGroup 创建角色分组
func (h *SsoRoleGroupHandler) CreateRoleGroup(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateRoleGroupRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 创建角色分组模型
	roleGroup := &userModel.SsoRoleGroup{
		ParentID:  req.ParentID,
		AppCode:   req.AppCode,
		OrgID:     req.OrgID,
		KindCode:  req.KindCode,
		GroupName: req.GroupName,
		GroupCode: req.GroupCode,
		Status:    req.Status,
		Remark:    req.Remark,
		Sort:      req.Sort,
		CreateID:  headerUser.UserID,
		CreateBy:  headerUser.UserAccount,
		TenantID:  headerUser.TenantID,
	}

	// 调用服务层创建角色分组
	if err := h.roleGroupService.CreateRoleGroup(roleGroup); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	afterData, _ := json.Marshal(roleGroup)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:        utils.GetSnowflake().GenerateIDString(),
		TraceID:   headerUser.TraceID,
		ModelID:   "role_group",
		RecordID:  roleGroup.ID,
		Action:    "CREATE",
		AfterData: string(afterData),
		CreateBy:  headerUser.UserAccount,
		Source:    "role_group_service",
	})

	ctx.JSON(201, roleGroup)
}

// GetRoleGroupByID 根据ID获取角色分组
func (h *SsoRoleGroupHandler) GetRoleGroupByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层获取角色分组
	roleGroup, err := h.roleGroupService.GetRoleGroupByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "角色分组不存在"})
		return
	}

	ctx.JSON(200, roleGroup)
}

// UpdateRoleGroup 更新角色分组
func (h *SsoRoleGroupHandler) UpdateRoleGroup(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req SsoUpdateRoleGroupRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 获取变更前数据
	beforeData, err := h.roleGroupService.GetRoleGroupByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "角色分组不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	// 更新角色分组字段
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
	if req.GroupName != nil {
		beforeData.GroupName = *req.GroupName
	}
	if req.GroupCode != nil {
		beforeData.GroupCode = *req.GroupCode
	}
	if req.Status != nil {
		beforeData.Status = *req.Status
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

	// 调用服务层更新角色分组
	if err := h.roleGroupService.UpdateRoleGroup(beforeData); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	afterJSON, _ := json.Marshal(beforeData)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "role_group",
		RecordID:   id,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "role_group_service",
	})

	ctx.JSON(200, beforeData)
}

// DeleteRoleGroup 删除角色分组
func (h *SsoRoleGroupHandler) DeleteRoleGroup(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 获取当前用户ID和账户
	headerUser := h.GetHeaderUserStruct(c, ctx)
	fmt.Println(headerUser.UserID, headerUser.UserAccount, headerUser.TenantID)

	// 获取删除前数据
	beforeData, err := h.roleGroupService.GetRoleGroupByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "角色分组不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	// 调用服务层删除角色分组
	if err := h.roleGroupService.DeleteRoleGroup(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "role_group",
		RecordID:   id,
		Action:     "DELETE",
		BeforeData: string(beforeJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "role_group_service",
	})

	ctx.JSON(200, map[string]string{"message": "角色分组删除成功"})
}

// GetAllRoleGroups 获取所有角色分组
func (h *SsoRoleGroupHandler) GetAllRoleGroups(c context.Context, ctx *app.RequestContext) {
	// 调用服务层获取所有角色分组
	roleGroups, err := h.roleGroupService.GetAllRoleGroups()
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取角色分组列表失败"})
		return
	}

	ctx.JSON(200, roleGroups)
}

// GetRoleGroupRoles 获取角色组的角色列表
func (h *SsoRoleGroupHandler) GetRoleGroupRoles(c context.Context, ctx *app.RequestContext) {
	groupID := ctx.Param("id")

	// 调用服务层获取角色组的角色ID列表
	roleIDs, err := h.roleGroupService.GetRoleGroupRoles(groupID)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, map[string]any{"role_ids": roleIDs})
}

// UpdateRoleGroupRolesRequest 更新角色组角色请求结构
type UpdateRoleGroupRolesRequest struct {
	RoleIDs []string `json:"role_ids"`
}

// UpdateRoleGroupRoles 更新角色组的角色关联
func (h *SsoRoleGroupHandler) UpdateRoleGroupRoles(c context.Context, ctx *app.RequestContext) {
	groupID := ctx.Param("id")
	var req UpdateRoleGroupRolesRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取当前用户信息
	headerUser := h.GetHeaderUserStruct(c, ctx)

	// 获取更新前的角色列表（用于审计日志）
	beforeRoleIDs, _ := h.roleGroupService.GetRoleGroupRoles(groupID)
	beforeData, _ := json.Marshal(map[string]any{"role_ids": beforeRoleIDs})

	// 调用服务层更新角色组角色关联
	if err := h.roleGroupService.UpdateRoleGroupRoles(groupID, req.RoleIDs, headerUser.UserAccount); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取更新后的角色列表（用于审计日志）
	afterRoleIDs, _ := h.roleGroupService.GetRoleGroupRoles(groupID)
	afterData, _ := json.Marshal(map[string]any{"role_ids": afterRoleIDs})

	// 记录数据变更日志
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "role_group_role",
		RecordID:   groupID,
		Action:     "UPDATE_ROLES",
		BeforeData: string(beforeData),
		AfterData:  string(afterData),
		CreateBy:   headerUser.UserAccount,
		Source:     "role_group_service",
	})

	ctx.JSON(200, map[string]string{"message": "更新成功"})
}
