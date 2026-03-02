package api

import (
	"context"
	"metadata-platform/internal/module/document/model"
	"metadata-platform/internal/module/document/repository"
	"metadata-platform/internal/module/document/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"
)

// DocumentFolderHandler 文档目录处理器
type DocumentFolderHandler struct {
	svc service.DocumentFolderService
}

// NewDocumentFolderHandler 创建文档目录处理器
func NewDocumentFolderHandler(db *gorm.DB) *DocumentFolderHandler {
	repo := repository.NewDocumentFolderRepository(db)
	svc := service.NewDocumentFolderService(repo)
	return &DocumentFolderHandler{
		svc: svc,
	}
}

// CreateFolder 创建文件夹
// POST /api/documents/folders
func (h *DocumentFolderHandler) CreateFolder(c context.Context, ctx *app.RequestContext) {
	var folder model.DocumentFolder
	if err := ctx.BindAndValidate(&folder); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}
	
	created, err := h.svc.CreateFolder(ctx, &folder)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "创建文件夹失败",
			"error":   err.Error(),
		})
		return
	}
	
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"success": true,
		"message": "文件夹创建成功",
		"data":    created,
	})
}

// UpdateFolder 更新文件夹
// PUT /api/documents/folders/:id
func (h *DocumentFolderHandler) UpdateFolder(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	
	var folder model.DocumentFolder
	if err := ctx.BindAndValidate(&folder); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}
	
	updated, err := h.svc.UpdateFolder(ctx, id, &folder)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "更新文件夹失败",
			"error":   err.Error(),
		})
		return
	}
	
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"success": true,
		"message": "文件夹更新成功",
		"data":    updated,
	})
}

// DeleteFolder 删除文件夹
// DELETE /api/documents/folders/:id
func (h *DocumentFolderHandler) DeleteFolder(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	
	if err := h.svc.DeleteFolder(id); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "删除文件夹失败",
			"error":   err.Error(),
		})
		return
	}
	
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"success": true,
		"message": "文件夹删除成功",
	})
}

// GetFolderByID 根据 ID 获取文件夹
// GET /api/documents/folders/:id
func (h *DocumentFolderHandler) GetFolderByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	
	folder, err := h.svc.GetFolderByID(id)
	if err != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "文件夹不存在",
			"error":   err.Error(),
		})
		return
	}
	
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"success": true,
		"data":    folder,
	})
}

// GetFolderByPath 根据路径获取文件夹
// GET /api/documents/folders/by-path
func (h *DocumentFolderHandler) GetFolderByPath(c context.Context, ctx *app.RequestContext) {
	path := ctx.Query("path")
	if path == "" {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "缺少 path 参数",
		})
		return
	}
	
	folder, err := h.svc.GetFolderByPath(path)
	if err != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "文件夹不存在",
			"error":   err.Error(),
		})
		return
	}
	
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"success": true,
		"data":    folder,
	})
}

// GetFolderList 获取文件夹列表
// GET /api/documents/folders
func (h *DocumentFolderHandler) GetFolderList(c context.Context, ctx *app.RequestContext) {
	parentID := ctx.Query("parentID")
	enabledOnly := ctx.Query("enabledOnly") == "true"
	
	folders, err := h.svc.GetFolderList(ctx, parentID, enabledOnly)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "获取文件夹列表失败",
			"error":   err.Error(),
		})
		return
	}
	
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"success": true,
		"data":    folders,
	})
}

// GetFolderTree 获取文件夹树
// GET /api/documents/folders/tree
func (h *DocumentFolderHandler) GetFolderTree(c context.Context, ctx *app.RequestContext) {
	tree, err := h.svc.GetFolderTree()
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "获取文件夹树失败",
			"error":   err.Error(),
		})
		return
	}
	
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"success": true,
		"data":    tree,
	})
}

// MoveFolder 移动文件夹
// POST /api/documents/folders/:id/move
func (h *DocumentFolderHandler) MoveFolder(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	
	var req struct {
		NewParentID string `json:"newParentId"`
	}
	
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}
	
	if err := h.svc.MoveFolder(id, req.NewParentID); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "移动文件夹失败",
			"error":   err.Error(),
		})
		return
	}
	
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"success": true,
		"message": "文件夹移动成功",
	})
}

// CopyFolder 复制文件夹
// POST /api/documents/folders/:id/copy
func (h *DocumentFolderHandler) CopyFolder(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	
	var req struct {
		NewParentID string `json:"newParentId"`
	}
	
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}
	
	copied, err := h.svc.CopyFolder(id, req.NewParentID)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "复制文件夹失败",
			"error":   err.Error(),
		})
		return
	}
	
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"success": true,
		"message": "文件夹复制成功",
		"data":    copied,
	})
}
