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
	return &HeaderUser{
		UserID:      userID,
		UserAccount: userAccount,
		TenantID:    tenantID,
	}
}
