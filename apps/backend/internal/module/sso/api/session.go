package api

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"

	"metadata-platform/internal/utils"
)

// ListSessions 获取会话列表
func (h *SsoHandler) ListSessions(c context.Context, ctx *app.RequestContext) {
	sessions, err := h.services.Session.GetAllSessions()
	if err != nil {
		utils.ErrorResponse(ctx, 500, "获取会话列表失败")
		return
	}
	utils.SuccessResponse(ctx, sessions)
}

// GetSession 根据ID获取会话
func (h *SsoHandler) GetSession(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	session, err := h.services.Session.GetSessionByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, 404, "会话不存在")
		return
	}
	utils.SuccessResponse(ctx, session)
}

// DeleteSession 删除会话
func (h *SsoHandler) DeleteSession(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	if err := h.services.Session.DeleteSession(id); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("删除会话失败: %v", err))
		return
	}
	utils.SuccessResponse(ctx, map[string]string{"message": "删除成功"})
}

// RevokeSession 撤销会话
func (h *SsoHandler) RevokeSession(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	if err := h.services.Session.RevokeSession(id); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("撤销会话失败: %v", err))
		return
	}
	utils.SuccessResponse(ctx, map[string]string{"message": "撤销成功"})
}
