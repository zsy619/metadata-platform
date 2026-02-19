package api

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"

	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/module/user/service"
	"metadata-platform/internal/utils"
	userModel "metadata-platform/internal/module/user/model"
)

type SsoUserGroupHandler struct {
	*utils.BaseHandler
	userGroupService service.SsoUserGroupService
	audit            AuditService
}

func NewSsoUserGroupHandler(userGroupService service.SsoUserGroupService, auditQueue *queue.AuditLogQueue) *SsoUserGroupHandler {
	return &SsoUserGroupHandler{
		BaseHandler:      utils.NewBaseHandler(),
		userGroupService: userGroupService,
		audit:            &auditServiceImpl{queue: auditQueue},
	}
}

type SsoCreateUserGroupRequest struct {
	ParentID  string `json:"parent_id" form:"parent_id"`
	AppCode   string `json:"app_code" form:"app_code"`
	OrgID     string `json:"org_id" form:"org_id"`
	KindCode  string `json:"kind_code" form:"kind_code"`
	GroupName string `json:"group_name" form:"group_name" binding:"required"`
	GroupCode string `json:"group_code" form:"group_code" binding:"required"`
	Status    int    `json:"status" form:"status"`
	Remark    string `json:"remark" form:"remark"`
	Sort      int    `json:"sort" form:"sort"`
}

type SsoUpdateUserGroupRequest struct {
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

func (h *SsoUserGroupHandler) CreateUserGroup(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateUserGroupRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	headerUser := h.GetHeaderUserStruct(c, ctx)

	item := &userModel.SsoUserGroup{
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

	item.ID = utils.GetSnowflake().GenerateIDString()

	if err := h.userGroupService.CreateUserGroup(item); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	afterData, _ := json.Marshal(item)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:        utils.GetSnowflake().GenerateIDString(),
		TraceID:   headerUser.TraceID,
		ModelID:   "user_group",
		RecordID:  item.ID,
		Action:    "CREATE",
		AfterData: string(afterData),
		CreateBy:  headerUser.UserAccount,
		Source:    "user_group_service",
	})

	ctx.JSON(201, item)
}

func (h *SsoUserGroupHandler) GetUserGroupByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	item, err := h.userGroupService.GetUserGroupByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "用户组不存在"})
		return
	}

	ctx.JSON(200, item)
}

func (h *SsoUserGroupHandler) UpdateUserGroup(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req SsoUpdateUserGroupRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	headerUser := h.GetHeaderUserStruct(c, ctx)

	beforeData, err := h.userGroupService.GetUserGroupByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "用户组不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

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

	beforeData.UpdateID = headerUser.UserID
	beforeData.UpdateBy = headerUser.UserAccount

	if err := h.userGroupService.UpdateUserGroup(beforeData); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	afterJSON, _ := json.Marshal(beforeData)
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "user_group",
		RecordID:   id,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "user_group_service",
	})

	ctx.JSON(200, beforeData)
}

func (h *SsoUserGroupHandler) DeleteUserGroup(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	headerUser := h.GetHeaderUserStruct(c, ctx)

	beforeData, err := h.userGroupService.GetUserGroupByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "用户组不存在"})
		return
	}
	beforeJSON, _ := json.Marshal(beforeData)

	if err := h.userGroupService.DeleteUserGroup(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "user_group",
		RecordID:   id,
		Action:     "DELETE",
		BeforeData: string(beforeJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "user_group_service",
	})

	ctx.JSON(200, map[string]string{"message": "用户组删除成功"})
}

func (h *SsoUserGroupHandler) GetAllUserGroups(c context.Context, ctx *app.RequestContext) {
	items, err := h.userGroupService.GetAllUserGroups()
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取用户组列表失败"})
		return
	}

	ctx.JSON(200, items)
}

// GetUserGroupRoles 获取用户组的角色列表
func (h *SsoUserGroupHandler) GetUserGroupRoles(c context.Context, ctx *app.RequestContext) {
	groupID := ctx.Param("id")

	// 调用服务层获取用户组的角色ID列表
	roleIDs, err := h.userGroupService.GetUserGroupRoles(groupID)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, map[string]any{"role_ids": roleIDs})
}

// UpdateUserGroupRolesRequest 更新用户组角色请求结构
type UpdateUserGroupRolesRequest struct {
	RoleIDs []string `json:"role_ids"`
}

// UpdateUserGroupRoles 更新用户组的角色关联
func (h *SsoUserGroupHandler) UpdateUserGroupRoles(c context.Context, ctx *app.RequestContext) {
	groupID := ctx.Param("id")
	var req UpdateUserGroupRolesRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取当前用户信息
	headerUser := h.GetHeaderUserStruct(c, ctx)

	// 获取更新前的角色列表（用于审计日志）
	beforeRoleIDs, _ := h.userGroupService.GetUserGroupRoles(groupID)
	beforeData, _ := json.Marshal(map[string]any{"role_ids": beforeRoleIDs})

	// 调用服务层更新用户组角色关联
	if err := h.userGroupService.UpdateUserGroupRoles(groupID, req.RoleIDs, headerUser.UserAccount); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取更新后的角色列表（用于审计日志）
	afterRoleIDs, _ := h.userGroupService.GetUserGroupRoles(groupID)
	afterData, _ := json.Marshal(map[string]any{"role_ids": afterRoleIDs})

	// 记录数据变更日志
	h.audit.RecordDataChange(c, &model.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "user_group_role",
		RecordID:   groupID,
		Action:     "UPDATE_ROLES",
		BeforeData: string(beforeData),
		AfterData:  string(afterData),
		CreateBy:   headerUser.UserAccount,
		Source:     "user_group_service",
	})

	ctx.JSON(200, map[string]string{"message": "更新成功"})
}
