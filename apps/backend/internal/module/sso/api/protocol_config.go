package api

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"

	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/utils"
)

// ListProtocolConfigs 获取协议配置列表
func (h *SsoHandler) ListProtocolConfigs(c context.Context, ctx *app.RequestContext) {
	configs, err := h.services.ProtocolConfig.GetAllConfigs()
	if err != nil {
		utils.ErrorResponse(ctx, 500, "获取协议配置列表失败")
		return
	}
	utils.SuccessResponse(ctx, configs)
}

// GetProtocolConfig 根据ID获取协议配置
func (h *SsoHandler) GetProtocolConfig(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	config, err := h.services.ProtocolConfig.GetConfigByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, 404, "协议配置不存在")
		return
	}
	utils.SuccessResponse(ctx, config)
}

// CreateProtocolConfig 创建协议配置
func (h *SsoHandler) CreateProtocolConfig(c context.Context, ctx *app.RequestContext) {
	var config model.SsoProtocolConfig
	if !utils.BindAndValidate(ctx, &config) {
		return
	}

	config.ID = uuid.New().String()
	config.TenantID = "1"

	if err := h.services.ProtocolConfig.CreateConfig(&config); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("创建协议配置失败: %v", err))
		return
	}

	utils.SuccessResponse(ctx, config)
}

// UpdateProtocolConfig 更新协议配置
func (h *SsoHandler) UpdateProtocolConfig(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	
	var config model.SsoProtocolConfig
	if !utils.BindAndValidate(ctx, &config) {
		return
	}

	config.ID = id
	if err := h.services.ProtocolConfig.UpdateConfig(&config); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("更新协议配置失败: %v", err))
		return
	}

	updated, err := h.services.ProtocolConfig.GetConfigByID(id)
	if err != nil {
		utils.SuccessResponse(ctx, map[string]string{"message": "更新成功"})
		return
	}

	utils.SuccessResponse(ctx, updated)
}

// DeleteProtocolConfig 删除协议配置
func (h *SsoHandler) DeleteProtocolConfig(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	if err := h.services.ProtocolConfig.DeleteConfig(id); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("删除协议配置失败: %v", err))
		return
	}
	utils.SuccessResponse(ctx, map[string]string{"message": "删除成功"})
}

// GetConfigsByType 根据协议类型获取配置列表
func (h *SsoHandler) GetConfigsByType(c context.Context, ctx *app.RequestContext) {
	protocolType := ctx.Param("type")
	configs, err := h.services.ProtocolConfig.GetConfigByProtocolType(model.ProtocolType(protocolType))
	if err != nil {
		utils.ErrorResponse(ctx, 500, "获取协议配置列表失败")
		return
	}
	utils.SuccessResponse(ctx, configs)
}
