package api

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"

	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/utils"
)

// ListFieldMappings 获取字段映射列表
func (h *SsoHandler) ListFieldMappings(c context.Context, ctx *app.RequestContext) {
	mappings, err := h.services.FieldMapping.GetAllMappings()
	if err != nil {
		utils.ErrorResponse(ctx, 500, "获取字段映射列表失败")
		return
	}
	utils.SuccessResponse(ctx, mappings)
}

// GetFieldMapping 根据ID获取字段映射
func (h *SsoHandler) GetFieldMapping(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	mapping, err := h.services.FieldMapping.GetMappingByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, 404, "字段映射不存在")
		return
	}
	utils.SuccessResponse(ctx, mapping)
}

// CreateFieldMapping 创建字段映射
func (h *SsoHandler) CreateFieldMapping(c context.Context, ctx *app.RequestContext) {
	var mapping model.SsoFieldMapping
	if !utils.BindAndValidate(ctx, &mapping) {
		return
	}

	mapping.ID = uuid.New().String()
	mapping.TenantID = "1"

	if err := h.services.FieldMapping.CreateMapping(&mapping); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("创建字段映射失败: %v", err))
		return
	}

	utils.SuccessResponse(ctx, mapping)
}

// UpdateFieldMapping 更新字段映射
func (h *SsoHandler) UpdateFieldMapping(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	
	var mapping model.SsoFieldMapping
	if !utils.BindAndValidate(ctx, &mapping) {
		return
	}

	mapping.ID = id
	if err := h.services.FieldMapping.UpdateMapping(&mapping); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("更新字段映射失败: %v", err))
		return
	}

	updated, err := h.services.FieldMapping.GetMappingByID(id)
	if err != nil {
		utils.SuccessResponse(ctx, map[string]string{"message": "更新成功"})
		return
	}

	utils.SuccessResponse(ctx, updated)
}

// DeleteFieldMapping 删除字段映射
func (h *SsoHandler) DeleteFieldMapping(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	if err := h.services.FieldMapping.DeleteMapping(id); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("删除字段映射失败: %v", err))
		return
	}
	utils.SuccessResponse(ctx, map[string]string{"message": "删除成功"})
}
