package metadata

import (
	globalMiddleware "metadata-platform/internal/middleware"
	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/module/metadata/api"
	"metadata-platform/internal/module/metadata/api/middleware"
	"metadata-platform/internal/module/metadata/repository"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"
)

// RegisterRoutes 注册元数据模块路由
func RegisterRoutes(r *server.Hertz, db *gorm.DB, auditDB *gorm.DB, auditQueue *queue.AuditLogQueue) {
	// 初始化仓库
	repos := repository.NewRepositories(db)

	// 初始化服务
	services := service.NewServices(db, repos, auditDB, auditQueue)

	// 初始化处理器
	apiHandler := api.NewAPIHandler(services.API)
	connHandler := api.NewMdConnHandler(services.Conn)
	tableHandler := api.NewMdTableHandler(services.Table)
	fieldHandler := api.NewMdTableFieldHandler(services.TableField)
	modelHandler := api.NewMdModelHandler(services.Model)
	templateHandler := api.NewQueryTemplateHandler(services.QueryTemplate)
	enhancementHandler := api.NewFieldEnhancementHandler(services.FieldEnhancement)
	treeHandler := api.NewTreeHandler(services.Tree)
	masterDetailHandler := api.NewMasterDetailHandler(services.MasterDetail)
	dataIOHandler := api.NewDataIOHandler(services.DataIO)

	// 元数据模块路由组
	metadataGroup := r.Group("/api/metadata")
	metadataGroup.Use(globalMiddleware.TenantMiddleware())
	metadataGroup.Use(globalMiddleware.AuthMiddleware())
	metadataGroup.Use(middleware.AuditMiddleware(services.Audit))

	// API路由
	apiGroup := metadataGroup.Group("/apis")
	{
		apiGroup.POST("", apiHandler.CreateAPI)
		apiGroup.GET("/:id", apiHandler.GetAPIByID)
		apiGroup.PUT("/:id", apiHandler.UpdateAPI)
		apiGroup.DELETE("/:id", apiHandler.DeleteAPI)
		apiGroup.GET("", apiHandler.GetAllAPIs)
		apiGroup.POST("/:id/enable", apiHandler.EnableAPI)
		apiGroup.POST("/:id/test", apiHandler.TestAPI)
	}

	// 数据连接路由
	connGroup := metadataGroup.Group("/conns")
	{
		connGroup.POST("", connHandler.CreateConn)
		connGroup.GET("/:id", connHandler.GetConnByID)
		connGroup.PUT("/:id", connHandler.UpdateConn)
		connGroup.DELETE("/:id", connHandler.DeleteConn)
		connGroup.GET("", connHandler.GetAllConns)
		connGroup.GET("/parent/:parent_id", connHandler.GetConnsByParentID)
		connGroup.POST("/test-raw", connHandler.TestRawConnection)
		connGroup.POST("/:id/test", connHandler.TestConnection)
		connGroup.GET("/:id/schemas", connHandler.GetSchemas)
		connGroup.GET("/:id/tables", connHandler.GetTables)
		connGroup.GET("/:id/views", connHandler.GetViews)
		connGroup.GET("/:id/tables/:table/structure", connHandler.GetTableStructure)
		connGroup.GET("/:id/tables/:table/preview", connHandler.PreviewTableData)
	}

	// 表路由
	tableGroup := metadataGroup.Group("/tables")
	{
		tableGroup.POST("", tableHandler.CreateTable)
		tableGroup.GET("/:id", tableHandler.GetTableByID)
		tableGroup.PUT("/:id", tableHandler.UpdateTable)
		tableGroup.DELETE("/:id", tableHandler.DeleteTable)
		tableGroup.GET("", tableHandler.GetAllTables)
		tableGroup.GET("/conn/:conn_id", tableHandler.GetTablesByConnID)
	}

	// 字段路由
	fieldGroup := metadataGroup.Group("/fields")
	{
		fieldGroup.POST("", fieldHandler.CreateField)
		fieldGroup.GET("/:id", fieldHandler.GetFieldByID)
		fieldGroup.PUT("/:id", fieldHandler.UpdateField)
		fieldGroup.DELETE("/:id", fieldHandler.DeleteField)
		fieldGroup.GET("", fieldHandler.GetAllFields)
		fieldGroup.GET("/table/:table_id", fieldHandler.GetFieldsByTableID)
		fieldGroup.DELETE("/table/:table_id", fieldHandler.DeleteFieldsByTableID)
	}

	// 模型路由
	modelGroup := metadataGroup.Group("/models")
	{
		modelGroup.POST("/build-from-table", modelHandler.BuildFromTable)
		modelGroup.POST("/build-from-view", modelHandler.BuildFromView)
		modelGroup.POST("/build-from-sql", modelHandler.BuildFromSQL)
		modelGroup.POST("/test-sql", modelHandler.TestSQL)
		modelGroup.GET("/generate-code", modelHandler.GenerateCode)
		modelGroup.POST("", modelHandler.CreateModel)
		modelGroup.GET("/:id", modelHandler.GetModelByID)
		modelGroup.PUT("/:id", modelHandler.UpdateModel)
		modelGroup.DELETE("/:id", modelHandler.DeleteModel)

		// 模型字段路由
		modelGroup.GET("/:id/fields", modelHandler.GetFieldsByModelID)
		modelGroup.POST("/:id/fields", modelHandler.CreateModelField)
		modelGroup.PUT("/:id/fields/:fieldId", modelHandler.UpdateModelField)
		modelGroup.DELETE("/:id/fields/:fieldId", modelHandler.DeleteModelField)

		modelGroup.GET("", modelHandler.ListModels)
		modelGroup.GET("/conn/:conn_id", modelHandler.GetModelsByConnID)

		// 查询模板路由
		templateGroup := modelGroup.Group("/:id/query-templates")
		{
			templateGroup.GET("", templateHandler.GetTemplatesByModelID)
			templateGroup.POST("", templateHandler.CreateTemplate)
			templateGroup.GET("/:templateId", templateHandler.GetTemplateByID)
			templateGroup.PUT("/:templateId", templateHandler.UpdateTemplate)
			templateGroup.DELETE("/:templateId", templateHandler.DeleteTemplate)
			templateGroup.POST("/:templateId/set-default", templateHandler.SetDefault)
			templateGroup.POST("/:templateId/duplicate", templateHandler.DuplicateTemplate)
			templateGroup.GET("/:templateId/preview", templateHandler.PreviewTemplate)
		}

		// 字段增强配置路由
		modelGroup.GET("/:id/fields/enhancements", enhancementHandler.GetEnhancementsByModelID)
		modelGroup.PUT("/:id/fields/enhancements", enhancementHandler.UpdateEnhancements)
		modelGroup.POST("/:id/fields/batch-enhancements", enhancementHandler.BatchUpdateEnhancements)
	}

	// 树形结构 API
	treeGroup := r.Group("/api/tree")
	treeGroup.Use(middleware.AuditMiddleware(services.Audit))
	{
		treeGroup.GET("/:model_id", treeHandler.GetTree)
		treeGroup.POST("/:model_id/node", treeHandler.AddNode)
		treeGroup.PUT("/:model_id/node/:id/move", treeHandler.MoveNode) // 注意：TargetParentID 在 body 中
		treeGroup.DELETE("/:model_id/node/:id", treeHandler.DeleteNode)
		treeGroup.GET("/:model_id/node/:id/children", treeHandler.GetChildren)
		treeGroup.GET("/:model_id/node/:id/path", treeHandler.GetPath)
	}

	// 主子表管理路由
	mdGroup := r.Group("/api/master-detail")
	mdGroup.Use(middleware.AuditMiddleware(services.Audit))
	{
		mdGroup.POST("/:master/:detail", masterDetailHandler.CreateMasterDetail)
	}

	// 数据导入导出路由
	ioGroup := r.Group("/api/data")
	ioGroup.Use(middleware.AuditMiddleware(services.Audit))
	{
		ioGroup.GET("/:model_id/export", dataIOHandler.ExportData)
		ioGroup.GET("/:model_id/import-template", dataIOHandler.ImportTemplate)
		ioGroup.POST("/:model_id/import", dataIOHandler.ImportData)
	}

	// 注册动态路由
	dynamicRouter := api.NewDynamicRouter(r, services)
	if err := dynamicRouter.LoadAndRegisterAll(); err != nil {
		utils.SugarLogger.Errorf("Failed to register dynamic routes: %v", err)
	}
}
