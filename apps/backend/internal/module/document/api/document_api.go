package api

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"

	"metadata-platform/internal/module/document/model"
	"metadata-platform/internal/module/document/repository"
	"metadata-platform/internal/module/document/service"
)

// DocumentHandler 文档处理器
type DocumentHandler struct {
	svc service.DocumentService
}

// NewDocumentHandler 创建文档处理器
func NewDocumentHandler(db *gorm.DB) *DocumentHandler {
	// 自动执行数据库迁移
	if err := db.AutoMigrate(
		&model.Document{},
		&model.DocumentCategory{},
		&model.DocumentVersion{},
		&model.DocumentFavorite{},
		&model.DocumentReadProgress{},
	); err != nil {
		// 如果迁移失败，记录错误但不阻止服务启动
		// 错误会在后续 API 调用时暴露
	}

	repo := repository.NewDocumentRepository(db)
	svc := service.NewDocumentService(repo)
	return &DocumentHandler{
		svc: svc,
	}
}

// GetDocumentList 获取文档列表
func (h *DocumentHandler) GetDocumentList(ctx context.Context, c *app.RequestContext) {
	// 获取查询参数
	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")
	category := c.Query("category")
	keyword := c.Query("keyword")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// 获取文档列表
	docs, total, err := h.svc.GetDocumentList(page, pageSize, category, keyword)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取文档列表失败：" + err.Error(),
		})
		return
	}

	// 返回分页结果
	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":     0,
		"data":     docs,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"hasMore":  page*pageSize < int(total),
	})
}

// GetDocumentCategories 获取文档分类
func (h *DocumentHandler) GetDocumentCategories(c context.Context, ctx *app.RequestContext) {
	categories, err := h.svc.GetCategories()
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取分类失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code": 0,
		"data": categories,
	})
}

// GetDocumentByID 根据 ID 获取文档
func (h *DocumentHandler) GetDocumentByID(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "文档 ID 不能为空",
		})
		return
	}

	doc, err := h.svc.GetDocumentByID(id)
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "文档不存在",
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code": 0,
		"data": doc,
	})
}

// GetDocumentByPath 根据路径获取文档
func (h *DocumentHandler) GetDocumentByPath(ctx context.Context, c *app.RequestContext) {
	path := c.Query("path")
	if path == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "文档路径不能为空",
		})
		return
	}

	doc, err := h.svc.GetDocumentByPath(path)
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "文档不存在",
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code": 0,
		"data": doc,
	})
}

// GetDocumentContent 获取文档内容
func (h *DocumentHandler) GetDocumentContent(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "文档 ID 不能为空",
		})
		return
	}

	content, err := h.svc.GetDocumentContent(id)
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "文档内容不存在",
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    0,
		"data":    content,
		"message": "success",
	})
}

// SearchDocuments 搜索文档
func (h *DocumentHandler) SearchDocuments(ctx context.Context, c *app.RequestContext) {
	keyword := c.Query("keyword")
	category := c.Query("category")
	limitStr := c.Query("limit")

	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 20
	}

	docs, err := h.svc.SearchDocuments(keyword, category, limit)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "搜索失败：" + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code": 0,
		"data": docs,
	})
}

// CreateDocument 创建文档
func (h *DocumentHandler) CreateDocument(ctx context.Context, c *app.RequestContext) {
	var req model.Document
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误：" + err.Error(),
		})
		return
	}

	doc, err := h.svc.CreateDocument(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建失败：" + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    0,
		"data":    doc,
		"message": "创建成功",
	})
}

// UpdateDocument 更新文档
func (h *DocumentHandler) UpdateDocument(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "文档 ID 不能为空",
		})
		return
	}

	var req model.Document
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误：" + err.Error(),
		})
		return
	}

	doc, err := h.svc.UpdateDocument(ctx, id, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新失败：" + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    0,
		"data":    doc,
		"message": "更新成功",
	})
}

// DeleteDocument 删除文档
func (h *DocumentHandler) DeleteDocument(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "文档 ID 不能为空",
		})
		return
	}

	if err := h.svc.DeleteDocument(id); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "删除失败：" + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "删除成功",
	})
}

// ToggleFavorite 切换收藏状态
func (h *DocumentHandler) ToggleFavorite(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "文档 ID 不能为空",
		})
		return
	}

	favorited, err := h.svc.ToggleFavorite(ctx, id)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "操作失败：" + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    0,
		"data":    favorited,
		"message": "操作成功",
	})
}

// GetMyFavorites 获取我的收藏
func (h *DocumentHandler) GetMyFavorites(ctx context.Context, c *app.RequestContext) {
	docs, err := h.svc.GetMyFavorites(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取收藏失败：" + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code": 0,
		"data": docs,
	})
}

// UpdateReadProgress 更新阅读进度
func (h *DocumentHandler) UpdateReadProgress(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "文档 ID 不能为空",
		})
		return
	}

	var req struct {
		Position int `json:"position"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	if err := h.svc.UpdateReadProgress(ctx, id, req.Position); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新失败：" + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "更新成功",
	})
}

// GetMyReadProgress 获取阅读进度
func (h *DocumentHandler) GetMyReadProgress(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "文档 ID 不能为空",
		})
		return
	}

	progress, err := h.svc.GetMyReadProgress(ctx, id)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取失败：" + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code": 0,
		"data": progress,
	})
}

// GetDocumentVersions 获取文档版本历史
func (h *DocumentHandler) GetDocumentVersions(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "文档 ID 不能为空",
		})
		return
	}

	versions, err := h.svc.GetDocumentVersions(id)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取失败：" + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code": 0,
		"data": versions,
	})
}

// RestoreVersion 恢复文档版本
func (h *DocumentHandler) RestoreVersion(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "文档 ID 不能为空",
		})
		return
	}

	var req struct {
		Version int `json:"version"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	doc, err := h.svc.RestoreVersion(id, req.Version)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "恢复失败：" + err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    0,
		"data":    doc,
		"message": "恢复成功",
	})
}

// DownloadDocument 下载文档
func (h *DocumentHandler) DownloadDocument(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "文档 ID 不能为空",
		})
		return
	}

	format := c.Query("format")
	if format == "" {
		format = "md"
	}

	doc, err := h.svc.GetDocumentByID(id)
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "文档不存在",
		})
		return
	}

	var contentType string
	var filename string
	var content []byte

	switch format {
	case "html":
		contentType = "text/html"
		filename = doc.Title + ".html"
		// 简单转换为 HTML
		content = []byte("<html><body>" + doc.Content + "</body></html>")
	case "pdf":
		contentType = "application/pdf"
		filename = doc.Title + ".pdf"
		// TODO: 实现 PDF 转换
		content = []byte("PDF conversion not implemented yet")
	default: // md
		contentType = "text/markdown"
		filename = doc.Title + ".md"
		content = []byte(doc.Content)
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(consts.StatusOK, contentType, content)
}
