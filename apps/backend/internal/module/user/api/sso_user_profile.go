package api

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"

	auditModel "metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/module/audit/queue"
	userModel "metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"
	"metadata-platform/internal/utils"
)

// SsoUserProfileHandler 用户档案处理器
type SsoUserProfileHandler struct {
	*utils.BaseHandler
	profileSvc service.SsoUserProfileService
	audit      AuditService
}

// NewSsoUserProfileHandler 创建用户档案处理器实例
func NewSsoUserProfileHandler(profileSvc service.SsoUserProfileService, auditQueue *queue.AuditLogQueue) *SsoUserProfileHandler {
	return &SsoUserProfileHandler{
		BaseHandler: utils.NewBaseHandler(),
		profileSvc:  profileSvc,
		audit:       &auditServiceImpl{queue: auditQueue},
	}
}

// GetUserProfile 获取用户档案
func (h *SsoUserProfileHandler) GetUserProfile(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	profile, err := h.profileSvc.GetByUserID(userID)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取用户档案失败"})
		return
	}
	ctx.JSON(200, profile)
}

// UpsertUserProfileRequest 创建/更新档案请求
type UpsertUserProfileRequest struct {
	Nickname string  `json:"nickname" form:"nickname"`
	Avatar   string  `json:"avatar" form:"avatar"`
	Gender   string  `json:"gender" form:"gender"`
	Birthday *string `json:"birthday" form:"birthday"`
	Bio      string  `json:"bio" form:"bio"`
	Location string  `json:"location" form:"location"`
}

// UpsertUserProfile 创建或更新用户档案
func (h *SsoUserProfileHandler) UpsertUserProfile(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	headerUser := h.GetHeaderUserStruct(c, ctx)

	var req UpsertUserProfileRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 获取变更前数据（用于审计）
	before, _ := h.profileSvc.GetByUserID(userID)
	beforeJSON, _ := json.Marshal(before)

	profile := &userModel.SsoUserProfile{
		UserID:   userID,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Gender:   req.Gender,
		Bio:      req.Bio,
		Location: req.Location,
		TenantID: headerUser.TenantID,
		CreateID: headerUser.UserID,
		CreateBy: headerUser.UserAccount,
		UpdateID: headerUser.UserID,
		UpdateBy: headerUser.UserAccount,
	}

	if err := h.profileSvc.Upsert(profile); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	afterJSON, _ := json.Marshal(profile)
	h.audit.RecordDataChange(c, &auditModel.SysDataChangeLog{
		ID:         utils.GetSnowflake().GenerateIDString(),
		TraceID:    headerUser.TraceID,
		ModelID:    "user_profile",
		RecordID:   userID,
		Action:     "UPSERT",
		BeforeData: string(beforeJSON),
		AfterData:  string(afterJSON),
		CreateBy:   headerUser.UserAccount,
		Source:     "user_profile_service",
	})

	ctx.JSON(200, profile)
}
