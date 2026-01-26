package api

import (
	"context"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// FieldEnhancementHandler 字段增强处理器
type FieldEnhancementHandler struct {
	enhService service.MdModelFieldEnhancementService
}

// NewFieldEnhancementHandler 创建字段增强处理器实例
func NewFieldEnhancementHandler(enhService service.MdModelFieldEnhancementService) *FieldEnhancementHandler {
	return &FieldEnhancementHandler{enhService: enhService}
}

// GetEnhancementsByModelID 获取模型的所有字段增强配置
func (h *FieldEnhancementHandler) GetEnhancementsByModelID(c context.Context, ctx *app.RequestContext) {
	modelID := ctx.Param("id")
	enhancements, err := h.enhService.GetEnhancementsByModelID(modelID)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, enhancements)
}

// UpdateEnhancements 更新单个字段增强配置
func (h *FieldEnhancementHandler) UpdateEnhancements(c context.Context, ctx *app.RequestContext) {
	var enh model.MdModelFieldEnhancement
	if err := ctx.BindJSON(&enh); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// 确保 URL 中的 modelID 与 body 中的一致 (虽然 enhService 可能不强校验，但最好一致)
	// 这里实际上 UpdateEnhancement 是根据 ID 更新，或者根据 field_id
	// 假设 RESTful path: /api/models/:id/fields/enhancements
	// 我们通常更新时 ID 必填。

	if err := h.enhService.UpdateEnhancement(&enh); err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, enh)
}

// BatchUpdateEnhancements 批量更新
func (h *FieldEnhancementHandler) BatchUpdateEnhancements(c context.Context, ctx *app.RequestContext) {
	var enhs []model.MdModelFieldEnhancement
	if err := ctx.BindJSON(&enhs); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, "Invalid JSON array")
		return
	}

	if err := h.enhService.BatchUpdateEnhancements(enhs); err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, "Batch update successful")
}
