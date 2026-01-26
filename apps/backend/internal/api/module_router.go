package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"

	metadata "metadata-platform/internal/module/metadata"

	user "metadata-platform/internal/module/user"
)

// RegisterModuleRoutes 注册所有模块的路由
func RegisterModuleRoutes(h *server.Hertz, metadataDB, userDB *gorm.DB) {
	// 注册元数据模块，使用元数据数据库
	registerMetadataRoutes(h, metadataDB)

	// 注册用户管理模块，使用用户管理数据库
	registerUserRoutes(h, userDB)
}

// registerMetadataRoutes 注册元数据模块路由
func registerMetadataRoutes(h *server.Hertz, db *gorm.DB) {
	// 注册元数据模块路由
	metadata.RegisterRoutes(h, db)
}

// registerUserRoutes 注册用户管理模块路由
func registerUserRoutes(h *server.Hertz, db *gorm.DB) {
	// 调用用户模块的统一路由注册
	user.RegisterRoutes(h, db)
}