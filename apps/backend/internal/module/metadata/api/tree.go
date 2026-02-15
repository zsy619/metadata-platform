package api

import (
	"context"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// TreeHandler 树形结构处理器
type TreeHandler struct {
	*utils.BaseHandler
	treeService service.TreeService
}

// NewTreeHandler 创建树形结构处理器实例
func NewTreeHandler(treeService service.TreeService) *TreeHandler {
	return &TreeHandler{
		BaseHandler: utils.NewBaseHandler(),
		treeService: treeService,
	}
}

// GetTree 获取完整树结构
func (h *TreeHandler) GetTree(c context.Context, ctx *app.RequestContext) {
	modelID := ctx.Param("model_id")
	tree, err := h.treeService.GetTree(modelID)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, tree)
}

// GetChildren 获取子节点
func (h *TreeHandler) GetChildren(c context.Context, ctx *app.RequestContext) {
	modelID := ctx.Param("model_id")
	id := ctx.Param("id")

	children, err := h.treeService.GetChildren(modelID, id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, children)
}

// GetPath 获取节点路径
func (h *TreeHandler) GetPath(c context.Context, ctx *app.RequestContext) {
	modelID := ctx.Param("model_id")
	id := ctx.Param("id")

	path, err := h.treeService.GetPath(modelID, id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, path)
}

// AddNode 添加节点
func (h *TreeHandler) AddNode(c context.Context, ctx *app.RequestContext) {
	modelID := ctx.Param("model_id")
	var data map[string]any
	if err := ctx.BindJSON(&data); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, "Invalid JSON payload")
		return
	}

	newNode, err := h.treeService.AddNode(modelID, data)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, newNode)
}

// MoveNode 移动节点
func (h *TreeHandler) MoveNode(c context.Context, ctx *app.RequestContext) {
	modelID := ctx.Param("model_id")
	id := ctx.Param("id")

	var req struct {
		TargetParentID string `json:"target_parent_id"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, consts.StatusBadRequest, "Invalid JSON payload")
		return
	}

	err := h.treeService.MoveNode(modelID, id, req.TargetParentID)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, "Moved successfully")
}

// DeleteNode 删除节点
func (h *TreeHandler) DeleteNode(c context.Context, ctx *app.RequestContext) {
	modelID := ctx.Param("model_id")
	id := ctx.Param("id")

	err := h.treeService.DeleteNode(modelID, id)
	if err != nil {
		utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, "Deleted successfully")
}
