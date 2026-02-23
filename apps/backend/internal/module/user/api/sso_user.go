package api

import (
	"context"
	"encoding/json"
	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/module/user/service"
	"metadata-platform/internal/utils"
	"time"

	auditModel "metadata-platform/internal/module/audit/model"

	userModel "metadata-platform/internal/module/user/model"

	"github.com/cloudwego/hertz/pkg/app"
)

type SsoUserHandler struct {
	*utils.BaseHandler
	userService service.SsoUserService
	audit       AuditService
}

func NewSsoUserHandler(userService service.SsoUserService, auditQueue *queue.AuditLogQueue) *SsoUserHandler {
	return &SsoUserHandler{
		BaseHandler: utils.NewBaseHandler(),
		userService: userService,
		audit:       &auditServiceImpl{queue: auditQueue},
	}
}

type SsoCreateUserRequest struct {
	Account  string     `json:"account" form:"account" binding:"required"`
	Password string     `json:"password" form:"password" binding:"required"`
	Name     string     `json:"name" form:"name"`
	Mobile   string     `json:"mobile" form:"mobile"`
	Email    string     `json:"email" form:"email"`
	Kind     int        `json:"kind" form:"kind"`
	Status   int        `json:"status" form:"status"`
	Remark   string     `json:"remark" form:"remark"`
	EndTime  *time.Time `json:"end_time" form:"end_time"`
}

type SsoUpdateUserRequest struct {
	Account  string     `json:"account" form:"account"`
	Password string     `json:"password" form:"password"`
	Name     string     `json:"name" form:"name"`
	Mobile   string     `json:"mobile" form:"mobile"`
	Email    string     `json:"email" form:"email"`
	Kind     int        `json:"kind" form:"kind"`
	Status   *int       `json:"status" form:"status"`
	Remark   string     `json:"remark" form:"remark"`
	EndTime  *time.Time `json:"end_time" form:"end_time"`
}

func (h *SsoUserHandler) CreateUser(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateUserRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	headerUser := h.GetHeaderUserStruct(c, ctx)

	user := &userModel.SsoUser{
		TenantID: headerUser.TenantID,
		Account:  req.Account,
		Password: req.Password,
		Name:     req.Name,
		Mobile:   req.Mobile,
		Email:    req.Email,
		Kind:     req.Kind,
		Status:   req.Status,
		EndTime:  req.EndTime,
		Remark:   req.Remark,
		CreateID: headerUser.UserID,
		CreateBy: headerUser.UserAccount,
		UpdateID: "",
		UpdateBy: "",
	}

	if err := h.userService.CreateUser(user); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	afterData, _ := json.Marshal(user)

	h.audit.RecordDataChange(c, &auditModel.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "sso_user",
		RecordID:   user.ID,
		Action:     "CREATE",
		BeforeData: "",
		AfterData:  string(afterData),
		CreateBy:   headerUser.UserAccount,
		Source:     "sso_user_service",
		Remark:     "新增用户",
	})

	user.Password = ""
	ctx.JSON(201, user)
}

func (h *SsoUserHandler) GetUserByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "用户不存在"})
		return
	}
	user.Password = ""
	ctx.JSON(200, user)
}

func (h *SsoUserHandler) UpdateUser(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req SsoUpdateUserRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "用户不存在"})
		return
	}

	// 获取修改前的数据快照
	beforeData, _ := json.Marshal(user)

	user.Account = req.Account
	// 密码：前端传了新密码则赋值（service 会加密），留空则清空（service 会保持原密码）
	user.Password = req.Password
	user.Name = req.Name
	// 手机号/邮箱允许清空，所以始终覆盖
	user.Mobile = req.Mobile
	user.Email = req.Email

	// status 使用指针，支持设为 0（禁用）
	if req.Status != nil {
		user.Status = *req.Status
	}
	user.Remark = req.Remark
	user.EndTime = req.EndTime

	if err := h.userService.UpdateUser(user); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取操作人信息和修改后的数据
	headerUser := h.GetHeaderUserStruct(c, ctx)
	// 由于 user 对象已包含最新属性信息（其中包含屏蔽密码后的 user）
	// 为准确记录修改后数据，重新从 db 获取最新数据
	afterUser, _ := h.userService.GetUserByID(id)
	afterData, _ := json.Marshal(afterUser)

	h.audit.RecordDataChange(c, &auditModel.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "sso_user",
		RecordID:   id,
		Action:     "UPDATE",
		BeforeData: string(beforeData),
		AfterData:  string(afterData),
		CreateBy:   headerUser.UserAccount,
		Source:     "sso_user_service",
		Remark:     "更新用户信息",
	})

	user.Password = ""
	ctx.JSON(200, user)
}

func (h *SsoUserHandler) DeleteUser(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 获取当前操作人信息
	headerUser := h.GetHeaderUserStruct(c, ctx)

	// 记录操作前查询用户信息
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "用户不存在，无法删除"})
		return
	}
	beforeData, _ := json.Marshal(user)

	if err := h.userService.DeleteUser(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 记录数据变更日志 (DELETE)
	h.audit.RecordDataChange(c, &auditModel.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "sso_user",
		RecordID:   id,
		Action:     "DELETE",
		BeforeData: string(beforeData),
		AfterData:  "",
		CreateBy:   headerUser.UserAccount, // 统一使用 UserAccount
		Source:     "sso_user_service",
		Remark:     "物理删除用户",
	})

	ctx.JSON(200, map[string]string{"message": "用户删除成功"})
}

func (h *SsoUserHandler) GetAllUsers(c context.Context, ctx *app.RequestContext) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取用户列表失败"})
		return
	}
	for i := range users {
		users[i].Password = ""
	}
	ctx.JSON(200, users)
}

func (h *SsoUserHandler) GetUserRoles(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	roleIDs, err := h.userService.GetUserRoles(userID)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	ctx.JSON(200, map[string]interface{}{"user_id": userID, "role_ids": roleIDs})
}

type UpdateUserRolesRequest struct {
	RoleIDs []string `json:"role_ids" binding:"required"`
}

func (h *SsoUserHandler) UpdateUserRoles(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	var req UpdateUserRolesRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	headerUser := h.GetHeaderUserStruct(c, ctx)

	beforeIDs, _ := h.userService.GetUserRoles(userID)
	beforeJSON, _ := json.Marshal(beforeIDs)

	if err := h.userService.UpdateUserRoles(userID, req.RoleIDs, headerUser.UserAccount); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	afterJSON, _ := json.Marshal(req.RoleIDs)
	h.audit.RecordDataChange(c, &auditModel.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "user_role",
		RecordID:   userID,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "user_service",
	})

	ctx.JSON(200, map[string]interface{}{"user_id": userID, "role_ids": req.RoleIDs, "message": "用户角色更新成功"})
}

func (h *SsoUserHandler) GetUserPos(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	posIDs, err := h.userService.GetUserPos(userID)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	ctx.JSON(200, map[string]interface{}{"user_id": userID, "pos_ids": posIDs})
}

type UpdateUserPosRequest struct {
	PosIDs []string `json:"pos_ids" binding:"required"`
}

func (h *SsoUserHandler) UpdateUserPos(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	var req UpdateUserPosRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	headerUser := h.GetHeaderUserStruct(c, ctx)

	beforeIDs, _ := h.userService.GetUserPos(userID)
	beforeJSON, _ := json.Marshal(beforeIDs)

	if err := h.userService.UpdateUserPos(userID, req.PosIDs, headerUser.UserAccount); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	afterJSON, _ := json.Marshal(req.PosIDs)
	h.audit.RecordDataChange(c, &auditModel.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "user_pos",
		RecordID:   userID,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "user_service",
	})

	ctx.JSON(200, map[string]interface{}{"user_id": userID, "pos_ids": req.PosIDs, "message": "用户职位更新成功"})
}

func (h *SsoUserHandler) GetUserGroups(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	groupIDs, err := h.userService.GetUserGroups(userID)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	ctx.JSON(200, map[string]interface{}{"user_id": userID, "group_ids": groupIDs})
}

type UpdateUserGroupsRequest struct {
	GroupIDs []string `json:"group_ids" binding:"required"`
}

func (h *SsoUserHandler) UpdateUserGroups(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	var req UpdateUserGroupsRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	headerUser := h.GetHeaderUserStruct(c, ctx)

	beforeIDs, _ := h.userService.GetUserGroups(userID)
	beforeJSON, _ := json.Marshal(beforeIDs)

	if err := h.userService.UpdateUserGroups(userID, req.GroupIDs, headerUser.UserAccount); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	afterJSON, _ := json.Marshal(req.GroupIDs)
	h.audit.RecordDataChange(c, &auditModel.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "user_group",
		RecordID:   userID,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "user_service",
	})

	ctx.JSON(200, map[string]interface{}{"user_id": userID, "group_ids": req.GroupIDs, "message": "用户用户组更新成功"})
}

func (h *SsoUserHandler) GetUserRoleGroups(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	roleGroupIDs, err := h.userService.GetUserRoleGroups(userID)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	ctx.JSON(200, map[string]interface{}{"user_id": userID, "role_group_ids": roleGroupIDs})
}

type UpdateUserRoleGroupsRequest struct {
	RoleGroupIDs []string `json:"role_group_ids" binding:"required"`
}

func (h *SsoUserHandler) UpdateUserRoleGroups(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	var req UpdateUserRoleGroupsRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	headerUser := h.GetHeaderUserStruct(c, ctx)

	beforeIDs, _ := h.userService.GetUserRoleGroups(userID)
	beforeJSON, _ := json.Marshal(beforeIDs)

	if err := h.userService.UpdateUserRoleGroups(userID, req.RoleGroupIDs, headerUser.UserAccount); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	afterJSON, _ := json.Marshal(req.RoleGroupIDs)
	h.audit.RecordDataChange(c, &auditModel.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "user_role_group",
		RecordID:   userID,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "user_service",
	})

	ctx.JSON(200, map[string]interface{}{"user_id": userID, "role_group_ids": req.RoleGroupIDs, "message": "用户角色组更新成功"})
}

func (h *SsoUserHandler) GetUserOrgs(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	orgIDs, err := h.userService.GetUserOrgs(userID)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	ctx.JSON(200, map[string]interface{}{"user_id": userID, "org_ids": orgIDs})
}

type UpdateUserOrgsRequest struct {
	OrgIDs []string `json:"org_ids" binding:"required"`
}

func (h *SsoUserHandler) UpdateUserOrgs(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	var req UpdateUserOrgsRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
	headerUser := h.GetHeaderUserStruct(c, ctx)

	beforeIDs, _ := h.userService.GetUserOrgs(userID)
	beforeJSON, _ := json.Marshal(beforeIDs)

	if err := h.userService.UpdateUserOrgs(userID, req.OrgIDs, headerUser.UserAccount); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	afterJSON, _ := json.Marshal(req.OrgIDs)
	h.audit.RecordDataChange(c, &auditModel.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "user_org",
		RecordID:   userID,
		Action:     "UPDATE",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "user_service",
	})

	ctx.JSON(200, map[string]interface{}{"user_id": userID, "org_ids": req.OrgIDs, "message": "用户组织更新成功"})
}
