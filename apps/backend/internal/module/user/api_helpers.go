package user

import (
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/module/user/service"

	"gorm.io/gorm"
)

// GetRepositories 暴露仓库集创建方法
func GetRepositories(db *gorm.DB) *repository.Repositories {
	return repository.NewRepositories(db)
}

// GetServices 暴露服务集创建方法
func GetServices(repos *repository.Repositories) *service.Services {
	return service.NewServices(repos)
}
