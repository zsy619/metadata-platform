package main

import (
	"fmt"
	"metadata-platform/configs"
	"metadata-platform/internal/api"
	"metadata-platform/internal/middleware"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"

	metadata "metadata-platform/internal/module/metadata"
	user "metadata-platform/internal/module/user"
)

// @title 元数据管理平台API
// @version 1.0
// @description 元数据管理平台后端API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http

func main() {
	// 1. 加载配置
	cfg, err := configs.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	// 2. 初始化日志系统
	utils.InitLogger(cfg.LogLevel, cfg.LogFilePath)
	defer utils.SyncLogger()
	utils.SugarLogger.Info("Application starting...")

	// 3. 初始化数据库连接管理器
	dbManager, err := utils.NewDBManager(cfg)
	if err != nil {
		utils.SugarLogger.Fatalf("Failed to initialize database manager: %v", err)
	}
	defer dbManager.Close()

	// 4. 初始化权限引擎 (需在 DB 之后)
	userDB := dbManager.GetUserDB()
	if err := middleware.InitCasbin(userDB, "configs/rbac_model.conf"); err != nil {
		utils.SugarLogger.Errorf("Failed to initialize Casbin: %v", err)
	}

	// 4. 自动迁移数据库表结构
	metadataDB := dbManager.GetMetadataDB()

	// 迁移元数据数据库
	migrateMetadataDatabase(metadataDB)
	// 迁移用户管理数据库
	migrateUserDatabase(userDB)

	// 5. 初始化数据库种子数据
	seedMetadataDatabase(metadataDB)
	seedUserDatabase(userDB)

	// 6. 初始化Hertz引擎
	r := server.Default(
		server.WithHostPorts(fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)),
	)

	// 6.5 执行权限同步
	syncCasbinPolicies(userDB)

	// 6. 初始化中间件加载器
	middlewareLoader := middleware.NewMiddlewareLoader()
	middlewareLoader.RegisterDefaultMiddlewares()
	middlewareLoader.LoadMiddlewareConfig(middleware.GetDefaultMiddlewareConfig())
	middlewareLoader.UseMiddlewareChain(r)

	// 7. 注册模块化路由
	api.RegisterModuleRoutes(r, metadataDB, userDB)

	// 8. 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	utils.SugarLogger.Infof("Server is running on %s", addr)
	utils.SugarLogger.Infof("Application started successfully in %s mode", cfg.AppMode)
	if err := r.Run(); err != nil {
		utils.SugarLogger.Fatalf("Failed to start server: %v", err)
	}
}

// migrateMetadataDatabase 自动迁移元数据数据库表结构
func migrateMetadataDatabase(db *gorm.DB) {
	if err := metadata.Migrate(db); err != nil {
		utils.SugarLogger.Errorf("Failed to migrate metadata database: %v", err)
	} else {
		utils.SugarLogger.Info("Metadata database migration completed successfully")
	}
}

// migrateUserDatabase 自动迁移用户管理数据库表结构
func migrateUserDatabase(db *gorm.DB) {
	if err := user.Migrate(db); err != nil {
		utils.SugarLogger.Errorf("Failed to migrate user database: %v", err)
	} else {
		utils.SugarLogger.Info("User database migration completed successfully")
	}
}

// seedMetadataDatabase 初始化元数据数据库种子数据
func seedMetadataDatabase(db *gorm.DB) {
	metadata.SeedData(db)
}

// seedUserDatabase 初始化用户管理数据库种子数据
func seedUserDatabase(db *gorm.DB) {
	user.SeedData(db)
}

// syncCasbinPolicies 系统启动时同步 Casbin 策略
func syncCasbinPolicies(db *gorm.DB) {
	// 初始化内部服务来执行同步
	repos := user.GetRepositories(db)
	services := user.GetServices(repos)
	if err := services.CasbinSync.SyncAll(); err != nil {
		utils.SugarLogger.Errorf("Failed to sync Casbin policies: %v", err)
	}
}
