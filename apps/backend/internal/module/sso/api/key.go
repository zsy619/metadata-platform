package api

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"

	"metadata-platform/internal/module/sso/model"
	"metadata-platform/internal/utils"
)

// ListKeys 获取密钥列表
func (h *SsoHandler) ListKeys(c context.Context, ctx *app.RequestContext) {
	keys, err := h.services.Key.GetAllKeys()
	if err != nil {
		utils.ErrorResponse(ctx, 500, "获取密钥列表失败")
		return
	}
	utils.SuccessResponse(ctx, keys)
}

// GetKey 根据ID获取密钥
func (h *SsoHandler) GetKey(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	key, err := h.services.Key.GetKeyByID(id)
	if err != nil {
		utils.ErrorResponse(ctx, 404, "密钥不存在")
		return
	}
	utils.SuccessResponse(ctx, key)
}

// CreateKey 创建密钥
func (h *SsoHandler) CreateKey(c context.Context, ctx *app.RequestContext) {
	var key model.SsoKey
	if !utils.BindAndValidate(ctx, &key) {
		return
	}

	key.ID = uuid.New().String()
	key.TenantID = "1"

	if err := h.services.Key.CreateKey(&key); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("创建密钥失败: %v", err))
		return
	}

	utils.SuccessResponse(ctx, key)
}

// GenerateKey 生成密钥对
func (h *SsoHandler) GenerateKey(c context.Context, ctx *app.RequestContext) {
	var req struct {
		KeyType   model.KeyType `json:"key_type" form:"key_type" binding:"required"`
		Algorithm string        `json:"algorithm" form:"algorithm"`
	}
	if !utils.BindAndValidate(ctx, &req) {
		return
	}

	key, err := h.services.Key.GenerateKeyPair(req.KeyType, req.Algorithm)
	if err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("生成密钥失败: %v", err))
		return
	}

	key.TenantID = "1"
	if err := h.services.Key.CreateKey(key); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("保存密钥失败: %v", err))
		return
	}

	utils.SuccessResponse(ctx, key)
}

// UpdateKey 更新密钥
func (h *SsoHandler) UpdateKey(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	
	var key model.SsoKey
	if !utils.BindAndValidate(ctx, &key) {
		return
	}

	key.ID = id
	if err := h.services.Key.UpdateKey(&key); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("更新密钥失败: %v", err))
		return
	}

	updated, err := h.services.Key.GetKeyByID(id)
	if err != nil {
		utils.SuccessResponse(ctx, map[string]string{"message": "更新成功"})
		return
	}

	utils.SuccessResponse(ctx, updated)
}

// DeleteKey 删除密钥
func (h *SsoHandler) DeleteKey(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	if err := h.services.Key.DeleteKey(id); err != nil {
		utils.ErrorResponse(ctx, 400, fmt.Sprintf("删除密钥失败: %v", err))
		return
	}
	utils.SuccessResponse(ctx, map[string]string{"message": "删除成功"})
}
