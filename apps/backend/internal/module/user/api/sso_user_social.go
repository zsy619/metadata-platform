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

// SsoUserSocialHandler 用户第三方账号处理器
type SsoUserSocialHandler struct {
	*utils.BaseHandler
	socialSvc service.SsoUserSocialService
	audit     AuditService
}

// NewSsoUserSocialHandler 创建用户第三方账号处理器实例
func NewSsoUserSocialHandler(socialSvc service.SsoUserSocialService, auditQueue *queue.AuditLogQueue) *SsoUserSocialHandler {
	return &SsoUserSocialHandler{
		BaseHandler: utils.NewBaseHandler(),
		socialSvc:   socialSvc,
		audit:       &auditServiceImpl{queue: auditQueue},
	}
}

// GetSocials 获取第三方账号绑定列表
func (h *SsoUserSocialHandler) GetSocials(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	list, err := h.socialSvc.GetByUserID(userID)
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取第三方账号失败"})
		return
	}
	ctx.JSON(200, list)
}

// BindSocialRequest 绑定第三方账号请求
type BindSocialRequest struct {
	Provider    string `json:"provider" form:"provider" binding:"required"`
	OpenID      string `json:"open_id" form:"open_id" binding:"required"`
	UnionID     string `json:"union_id" form:"union_id"`
	Nickname    string `json:"nickname" form:"nickname"`
	Avatar      string `json:"avatar" form:"avatar"`
	ProfileJSON string `json:"profile_json" form:"profile_json"`
}

// BindSocial 绑定第三方账号
func (h *SsoUserSocialHandler) BindSocial(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	headerUser := h.GetHeaderUserStruct(c, ctx)

	var req BindSocialRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	social := &userModel.SsoUserSocial{
		UserID:      userID,
		Provider:    req.Provider,
		OpenID:      req.OpenID,
		UnionID:     req.UnionID,
		Nickname:    req.Nickname,
		Avatar:      req.Avatar,
		ProfileJSON: req.ProfileJSON,
		TenantID:    headerUser.TenantID,
		CreateID:    headerUser.UserID,
		CreateBy:    headerUser.UserAccount,
		UpdateID:    headerUser.UserID,
		UpdateBy:    headerUser.UserAccount,
	}

	if err := h.socialSvc.Bind(social); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	afterJSON, _ := json.Marshal(social)
	h.audit.RecordDataChange(c, &auditModel.SysDataChangeLog{
		ID:        utils.GetSnowflake().GenerateIDString(),
		TraceID:   headerUser.TraceID,
		ModelID:   "user_social",
		RecordID:  social.ID,
		Action:    "BIND",
		AfterData: string(afterJSON),
		CreateBy:  headerUser.UserAccount,
		Source:    "user_social_service",
	})

	ctx.JSON(201, social)
}

// UnbindSocial 解绑第三方账号
func (h *SsoUserSocialHandler) UnbindSocial(c context.Context, ctx *app.RequestContext) {
	userID := ctx.Param("id")
	socialID := ctx.Param("sid")
	headerUser := h.GetHeaderUserStruct(c, ctx)

	if err := h.socialSvc.Unbind(userID, socialID); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	h.audit.RecordDataChange(c, &auditModel.SysDataChangeLog{
		ID:       utils.GetSnowflake().GenerateIDString(),
		TraceID:  headerUser.TraceID,
		ModelID:  "user_social",
		RecordID: socialID,
		Action:   "UNBIND",
		CreateBy: headerUser.UserAccount,
		Source:   "user_social_service",
	})

	ctx.JSON(200, map[string]string{"message": "解绑成功"})
}
