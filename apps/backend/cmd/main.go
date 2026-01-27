package main

import (
	"fmt"
	"metadata-platform/configs"
	"metadata-platform/internal/api"
	"metadata-platform/internal/middleware"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"

	audit "metadata-platform/internal/module/audit"
	auditQueuePkg "metadata-platform/internal/module/audit/queue"
	metadata "metadata-platform/internal/module/metadata"
	user "metadata-platform/internal/module/user"
)

// ... comments ...

func main() {
	// 1. 加载配置
	cfg, err := configs.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// 2. 初始化日志
	utils.InitLogger(cfg.LogLevel, cfg.LogFilePath)
	defer utils.SyncLogger()

	// 3. 初始化数据库管理器
	dbManager, err := utils.NewDBManager(cfg)
	if err != nil {
		utils.SugarLogger.Fatalf("Failed to initialize database manager: %v", err)
	}

	// 4. 获取各个数据库连接
	metadataDB := dbManager.GetMetadataDB()
	userDB := dbManager.GetUserDB()
	auditDB := dbManager.GetAuditDB()

	// 5. 执行数据库迁移和种子数据
	migrateMetadataDatabase(metadataDB)
	migrateUserDatabase(userDB)
	migrateAuditDatabase(auditDB)

	seedMetadataDatabase(metadataDB)
	seedUserDatabase(userDB)
	// 6. 初始化Hertz引擎
	r := server.Default(
		server.WithHostPorts(fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)),
	)

	// 6.4 初始化审计日志队列
	auditLogQueue := auditQueuePkg.NewAuditLogQueue(auditDB, 1000, 5)
	auditLogQueue.Start()
	defer auditLogQueue.Stop()

	// 6.5 执行权限同步
	syncCasbinPolicies(userDB, auditDB, auditLogQueue)

	// ...

	// 6. 初始化中间件加载器
	middlewareLoader := middleware.NewMiddlewareLoader()
	middlewareLoader.RegisterDefaultMiddlewares()
	middlewareLoader.LoadMiddlewareConfig(middleware.GetDefaultMiddlewareConfig())
	middlewareLoader.UseMiddlewareChain(r)

	// 7. 注册模块化路由
	api.RegisterModuleRoutes(r, metadataDB, userDB, auditDB, auditLogQueue)

	// 8. 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	utils.SugarLogger.Infof("Server is running on %s", addr)
	utils.SugarLogger.Infof("Application started successfully in %s mode", cfg.AppMode)
	if err := r.Run(); err != nil {
		utils.SugarLogger.Fatalf("Failed to start server: %v", err)
	}
}

// ... existing migration functions ...
// migrateMetadataDatabase, migrateUserDatabase, migrateAuditDatabase, seedMetadataDatabase, seedUserDatabase are unchanged.

func migrateMetadataDatabase(db *gorm.DB) {
	if err := metadata.Migrate(db); err != nil {
		utils.SugarLogger.Errorf("Failed to migrate metadata database: %v", err)
	} else {
		utils.SugarLogger.Info("Metadata database migration completed successfully")
	}
}

func migrateUserDatabase(db *gorm.DB) {
	if err := user.Migrate(db); err != nil {
		utils.SugarLogger.Errorf("Failed to migrate user database: %v", err)
	} else {
		utils.SugarLogger.Info("User database migration completed successfully")
	}
}

func migrateAuditDatabase(db *gorm.DB) {
	if err := audit.Migrate(db); err != nil {
		utils.SugarLogger.Errorf("Failed to migrate audit database: %v", err)
	} else {
		utils.SugarLogger.Info("Audit database migration completed successfully")
	}
}

func seedMetadataDatabase(db *gorm.DB) {
	metadata.SeedData(db)
}

func seedUserDatabase(db *gorm.DB) {
	user.SeedData(db)
}

// syncCasbinPolicies 系统启动时同步 Casbin 策略
func syncCasbinPolicies(db *gorm.DB, auditDB *gorm.DB, auditQueue *auditQueuePkg.AuditLogQueue) {
	// 初始化内部服务来执行同步
	repos := user.GetRepositories(db)
	services := user.GetServices(repos, auditDB, auditQueue)
	if err := services.CasbinSync.SyncAll(); err != nil {
		utils.SugarLogger.Errorf("Failed to sync Casbin policies: %v", err)
	}
}
