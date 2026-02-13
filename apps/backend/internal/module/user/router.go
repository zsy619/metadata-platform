package user

import (
	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/module/user/api"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/module/user/service"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"
)

// RegisterRoutes 注册用户管理模块路由
func RegisterRoutes(r *server.Hertz, db *gorm.DB, auditDB *gorm.DB, auditQueue *queue.AuditLogQueue) {
	// 初始化仓库集合
	repos := repository.NewRepositories(db)

	// 初始化服务集合
	services := service.NewServices(repos, auditDB, auditQueue)

	// 初始化组织类型服务
	orgKindService := service.NewOrgKindService(db)
	orgKindHandler := service.NewOrgKindHandler(orgKindService)

	// 初始化处理器集合
	handlers := api.NewSsoHandler(services)

	// 注册路由
	handlers.RegisterRoutes(r, orgKindHandler)
}
