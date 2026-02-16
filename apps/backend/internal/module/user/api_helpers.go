package user

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/module/user/service"
)

// GetRepositories 暴露仓库集创建方法
func GetRepositories(db *gorm.DB) *repository.Repositories {
	return repository.NewRepositories(db)
}

// GetServices 暴露服务集创建方法
func GetServices(repos *repository.Repositories, db *gorm.DB, auditDB *gorm.DB, auditQueue *queue.AuditLogQueue) *service.Services {
	return service.NewServices(repos, db, auditDB, auditQueue)
}
