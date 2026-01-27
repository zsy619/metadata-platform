package api

import (
	"metadata-platform/internal/module/audit/queue"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"

	audit "metadata-platform/internal/module/audit"

	metadata "metadata-platform/internal/module/metadata"

	user "metadata-platform/internal/module/user"
)

// RegisterModuleRoutes 注册所有模块的路由
func RegisterModuleRoutes(h *server.Hertz, metadataDB, userDB, auditDB *gorm.DB, auditQueue *queue.AuditLogQueue) {
	// 注册元数据模块，使用元数据数据库
	registerMetadataRoutes(h, metadataDB, auditDB, auditQueue)

	// 注册用户管理模块，使用用户管理数据库
	registerUserRoutes(h, userDB, auditDB, auditQueue)

	// Register Audit Module
	registerAuditRoutes(h, auditDB, auditQueue)
}

// registerMetadataRoutes 注册元数据模块路由
func registerMetadataRoutes(h *server.Hertz, db *gorm.DB, auditDB *gorm.DB, auditQueue *queue.AuditLogQueue) {
	// 注册元数据模块路由
	metadata.RegisterRoutes(h, db, auditDB, auditQueue)
}

// registerUserRoutes 注册用户管理模块路由
func registerUserRoutes(h *server.Hertz, db *gorm.DB, auditDB *gorm.DB, auditQueue *queue.AuditLogQueue) {
	// 调用用户模块的统一路由注册
	user.RegisterRoutes(h, db, auditDB, auditQueue)
}

// registerAuditRoutes 注册审计模块路由
func registerAuditRoutes(h *server.Hertz, db *gorm.DB, auditQueue *queue.AuditLogQueue) {
	audit.RegisterRoutes(h, db, auditQueue)
}
