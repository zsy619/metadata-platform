package user

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"

	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/module/user/api"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/module/user/service"
)

// RegisterRoutes 注册用户管理模块路由
func RegisterRoutes(r *server.Hertz, db *gorm.DB, auditDB *gorm.DB, auditQueue *queue.AuditLogQueue) {
	fmt.Println(">>> Initializing User Routes...")

	// 初始化仓库集合
	repos := repository.NewRepositories(db)

	// 初始化服务集合
	services := service.NewServices(repos, db, auditDB, auditQueue)

	// 初始化处理器集合
	handlers := api.NewSsoHandler(services, auditQueue)

	// 注册路由
	handlers.RegisterRoutes(r)
}
