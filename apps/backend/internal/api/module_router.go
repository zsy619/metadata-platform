package api

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"

	audit "metadata-platform/internal/module/audit"
	"metadata-platform/internal/module/audit/queue"
	metadata "metadata-platform/internal/module/metadata"
	"metadata-platform/internal/module/monitor"
	sso "metadata-platform/internal/module/sso"
	user "metadata-platform/internal/module/user"
)

// RegisterModuleRoutes 注册所有模块的路由
func RegisterModuleRoutes(h *server.Hertz, metadataDB, userDB, auditDB *gorm.DB, auditQueue *queue.AuditLogQueue) {
	fmt.Fprintln(os.Stderr, ">>> Starting route registration...")
	registerMetadataRoutes(h, metadataDB, auditDB, auditQueue)
	fmt.Fprintln(os.Stderr, ">>> Metadata routes registered")
	registerUserRoutes(h, userDB, auditDB, auditQueue)
	fmt.Fprintln(os.Stderr, ">>> User routes registered")
	fmt.Fprintln(os.Stderr, ">>> Registering SSO routes...")
	registerSSORoutes(h, userDB)
	fmt.Fprintln(os.Stderr, ">>> SSO routes registered successfully")
	registerAuditRoutes(h, auditDB, auditQueue)
	fmt.Fprintln(os.Stderr, ">>> Audit routes registered")
	monitor.RegisterRoutesWithDB(h, userDB, auditDB)
	fmt.Fprintln(os.Stderr, ">>> Monitor routes registered")

	h.GET("/debug/routes", func(c context.Context, ctx *app.RequestContext) {
		routes := h.Routes()
		// 创建简化的路由信息，避免序列化 HandlerFunc
		type SimpleRoute struct {
			Method string `json:"method"`
			Path   string `json:"path"`
		}
		var simpleRoutes []SimpleRoute
		for _, route := range routes {
			simpleRoutes = append(simpleRoutes, SimpleRoute{
				Method: route.Method,
				Path:   route.Path,
			})
		}
		ctx.JSON(200, simpleRoutes)
	})
}

func registerMetadataRoutes(h *server.Hertz, db *gorm.DB, auditDB *gorm.DB, auditQueue *queue.AuditLogQueue) {
	metadata.RegisterRoutes(h, db, auditDB, auditQueue)
}

func registerUserRoutes(h *server.Hertz, db *gorm.DB, auditDB *gorm.DB, auditQueue *queue.AuditLogQueue) {
	user.RegisterRoutes(h, db, auditDB, auditQueue)
}

func registerAuditRoutes(h *server.Hertz, db *gorm.DB, auditQueue *queue.AuditLogQueue) {
	audit.RegisterRoutes(h, db, auditQueue)
}

func registerSSORoutes(h *server.Hertz, db *gorm.DB) {
	sso.RegisterRoutes(h, db)
}
