package api

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"

	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/utils"
)

// ListClients 获取客户端配置列表
func (h *SsoHandler) ListClients(c context.Context, ctx *app.RequestContext) {
	clients, err := h.services.Client.GetAllClients()
	if err != nil {
		utils.ErrorResponse(ctx, 500, "获取客户端配置列表失败")
		return
	}
	utils.SuccessResponse(ctx, clients)
}

// GetClient 根据ID获取客户端配置
func (h *SsoHandler) GetClient(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	client, err := h.services.Client.GetClientByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, 404, "客户端配置不存在")
		return
	}
	utils.SuccessResponse(ctx, client)
}

// CreateClient 创建客户端配置
func (h *SsoHandler) CreateClient(c context.Context, ctx *app.RequestContext) {
	var client model.SsoClient
	if !utils.BindAndValidate(ctx, &client) {
		return
	}

	client.ID = uuid.New().String()
	client.TenantID = "1"

	if err := h.services.Client.CreateClient(&client); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("创建客户端配置失败: %v", err))
		return
	}

	utils.SuccessResponse(ctx, client)
}

// UpdateClient 更新客户端配置
func (h *SsoHandler) UpdateClient(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	
	var client model.SsoClient
	if !utils.BindAndValidate(ctx, &client) {
		return
	}

	client.ID = id
	if err := h.services.Client.UpdateClient(&client); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("更新客户端配置失败: %v", err))
		return
	}

	updated, err := h.services.Client.GetClientByID(id)
	if err != nil {
		utils.SuccessResponse(ctx, map[string]string{"message": "更新成功"})
		return
	}

	utils.SuccessResponse(ctx, updated)
}

// DeleteClient 删除客户端配置
func (h *SsoHandler) DeleteClient(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	if err := h.services.Client.DeleteClient(id); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("删除客户端配置失败: %v", err))
		return
	}
	utils.SuccessResponse(ctx, map[string]string{"message": "删除成功"})
}
