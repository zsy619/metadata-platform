package api

import (
	"context"
	"encoding/json"

	auditModel "metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/queue"
	userModel "metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"
	"metadata-platform/internal/utils"

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
	Account  string `json:"account" form:"account" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Name     string `json:"name" form:"name"`
	Mobile   string `json:"mobile" form:"mobile"`
	Email    string `json:"email" form:"email"`
	Kind     int    `json:"kind" form:"kind"`
	Status   int    `json:"status" form:"status"`
	Remark   string `json:"remark" form:"remark"`
}

type SsoUpdateUserRequest struct {
	Account  string `json:"account" form:"account"`
	Password string `json:"password" form:"password"`
	Name     string `json:"name" form:"name"`
	Mobile   string `json:"mobile" form:"mobile"`
	Email    string `json:"email" form:"email"`
	Kind     int    `json:"kind" form:"kind"`
	Status   int    `json:"status" form:"status"`
	Remark   string `json:"remark" form:"remark"`
}

func (h *SsoUserHandler) CreateUser(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateUserRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	user := &userModel.SsoUser{
		Account:  req.Account,
		Password: req.Password,
		Name:     req.Name,
		Mobile:   req.Mobile,
		Email:    req.Email,
		Kind:     req.Kind,
		Status:   req.Status,
		Remark:   req.Remark,
	}

	if err := h.userService.CreateUser(user); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

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

	if req.Account != "" {
		user.Account = req.Account
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Mobile != "" {
		user.Mobile = req.Mobile
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Kind != 0 {
		user.Kind = req.Kind
	}
	if req.Status != 0 {
		user.Status = req.Status
	}
	if req.Remark != "" {
		user.Remark = req.Remark
	}

	if err := h.userService.UpdateUser(user); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	user.Password = ""
	ctx.JSON(200, user)
}

func (h *SsoUserHandler) DeleteUser(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	if err := h.userService.DeleteUser(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}
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
