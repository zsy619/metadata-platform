package utils

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
)

// BaseHandler 基础处理器结构体
type BaseHandler struct{}

// NewBaseHandler 创建基础处理器实例
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

type HeaderUser struct {
	UserID      string
	UserAccount string
	TenantID    string
	TraceID     string
}

func (h *BaseHandler) GetHeaderUser(c context.Context, ctx *app.RequestContext) (string, string, string) {
	userID := string(ctx.Request.Header.Get("X-User-ID"))
	userAccount := string(ctx.Request.Header.Get("X-User-Account"))
	tenantID := string(ctx.Request.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		tenantID = "1"
	}
	fmt.Println(userID, userAccount, tenantID)
	return userID, userAccount, tenantID
}

func (h *BaseHandler) GetHeaderUserStruct(c context.Context, ctx *app.RequestContext) *HeaderUser {
	userID, userAccount, tenantID := h.GetHeaderUser(c, ctx)
	traceID := string(ctx.Request.Header.Get("X-Trace-ID"))
	if traceID == "" {
		traceID = string(ctx.Request.Header.Get("trace_id"))
	}
	if traceID == "" {
		traceID = ctx.GetString("trace_id")
	}
	if traceID == "" {
		// 获取 trace_id（如果 handler 中有设置则使用 handler 中的值）
		if storedTraceID, exists := ctx.Get("trace_id"); exists {
			traceID = storedTraceID.(string)
		}
	}
	return &HeaderUser{
		UserID:      userID,
		UserAccount: userAccount,
		TenantID:    tenantID,
		TraceID:     traceID,
	}
}
