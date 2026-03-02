package document

import (
	"metadata-platform/internal/module/document/api"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"
)

// RegisterRoutes 注册文档模块路由
func RegisterRoutes(h *server.Hertz, db *gorm.DB) {
	utils.SugarLogger.Info("Initializing document module...")

	// 初始化 Handler（api 包内部会创建 Repository 和 Service）
	handler := api.NewDocumentHandler(db)
	folderHandler := api.NewDocumentFolderHandler(db)

	// 文档模块路由组
	docGroup := h.Group("/api/documents")
	{
		// 获取文档列表
		docGroup.GET("", handler.GetDocumentList)

		// 获取所有分类
		docGroup.GET("/categories", handler.GetDocumentCategories)

		// 搜索文档
		docGroup.GET("/search", handler.SearchDocuments)

		// 根据 ID 获取文档
		docGroup.GET("/:id", handler.GetDocumentByID)

		// 根据路径获取文档
		docGroup.GET("/by-path", handler.GetDocumentByPath)

		// 获取文档内容
		docGroup.GET("/:id/content", handler.GetDocumentContent)

		// 创建文档
		docGroup.POST("", handler.CreateDocument)

		// 更新文档
		docGroup.PUT("/:id", handler.UpdateDocument)

		// 删除文档
		docGroup.DELETE("/:id", handler.DeleteDocument)

		// 下载文档
		docGroup.GET("/:id/download", handler.DownloadDocument)

		// 获取版本历史
		docGroup.GET("/:id/versions", handler.GetDocumentVersions)

		// 恢复版本
		docGroup.POST("/:id/versions/restore", handler.RestoreVersion)

		// 切换收藏状态
		docGroup.POST("/:id/favorite", handler.ToggleFavorite)

		// 获取我的收藏
		docGroup.GET("/favorites/my", handler.GetMyFavorites)

		// 更新阅读进度
		docGroup.POST("/:id/progress", handler.UpdateReadProgress)

		// 获取阅读进度
		docGroup.GET("/:id/progress", handler.GetMyReadProgress)
	}

	// 文件夹管理路由组
	folderGroup := h.Group("/api/documents/folders")
	{
		// 获取文件夹树
		folderGroup.GET("/tree", folderHandler.GetFolderTree)

		// 获取文件夹列表
		folderGroup.GET("", folderHandler.GetFolderList)

		// 根据路径获取文件夹
		folderGroup.GET("/by-path", folderHandler.GetFolderByPath)

		// 根据 ID 获取文件夹
		folderGroup.GET("/:id", folderHandler.GetFolderByID)

		// 创建文件夹
		folderGroup.POST("", folderHandler.CreateFolder)

		// 更新文件夹
		folderGroup.PUT("/:id", folderHandler.UpdateFolder)

		// 删除文件夹
		folderGroup.DELETE("/:id", folderHandler.DeleteFolder)

		// 移动文件夹
		folderGroup.POST("/:id/move", folderHandler.MoveFolder)

		// 复制文件夹
		folderGroup.POST("/:id/copy", folderHandler.CopyFolder)
	}

	utils.SugarLogger.Info("Document routes registered successfully")
}
