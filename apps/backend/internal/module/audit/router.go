package audit

import (
	"metadata-platform/internal/middleware"
	"metadata-platform/internal/module/audit/api"
	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/module/audit/service"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"
)

// RegisterRoutes 注册审计模块路由
func RegisterRoutes(r *server.Hertz, db *gorm.DB, auditQueue *queue.AuditLogQueue) {
	// 初始化服务
	svc := service.NewAuditService(db, auditQueue)
	handler := api.NewAuditHandler(svc)

	// 路由组
	g := r.Group("/api/audit")
	g.Use(middleware.AuthMiddleware()) // Protected

	// Login Logs
	g.GET("/login", handler.GetLoginLogs)
	g.GET("/login/export", handler.ExportLoginLogs)

	// Operation Logs
	g.GET("/operation", handler.GetOperationLogs)
	g.GET("/operation/export", handler.ExportOperationLogs)

	// Data Change Logs
	g.GET("/data", handler.GetDataChangeLogs)
	g.GET("/data/export", handler.ExportDataChangeLogs)
}
