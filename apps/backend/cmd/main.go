package main

import (
	"fmt"
	"os"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"

	"metadata-platform/configs"
	"metadata-platform/internal/api"
	"metadata-platform/internal/middleware"
	audit "metadata-platform/internal/module/audit"
	auditQueuePkg "metadata-platform/internal/module/audit/queue"
	metadata "metadata-platform/internal/module/metadata"
	sso "metadata-platform/internal/module/sso"
	user "metadata-platform/internal/module/user"
	"metadata-platform/internal/utils"
)

// ... comments ...

func main() {
	fmt.Fprintln(os.Stderr, "DEBUG: Starting main function...")
	// 1. 加载配置
	fmt.Fprintln(os.Stderr, "DEBUG: Loading config...")
	cfg, err := configs.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// 2. 初始化日志
	fmt.Fprintln(os.Stderr, "DEBUG: Config loaded. Initializing logger...")
	utils.InitLogger(cfg.LogLevel, cfg.LogFilePath)
	defer utils.SyncLogger()

	// 3. 初始化数据库管理器
	fmt.Fprintln(os.Stderr, "DEBUG: Logger initialized. Creating DB manager...")
	dbManager, err := utils.NewDBManager(cfg)
	if err != nil {
		utils.SugarLogger.Fatalf("Failed to initialize database manager: %v", err)
	}

	// 4. 获取各个数据库连接
	fmt.Fprintln(os.Stderr, "DEBUG: DB Manager created. Getting DB connections...")
	metadataDB := dbManager.GetMetadataDB()
	userDB := dbManager.GetUserDB()
	auditDB := dbManager.GetAuditDB()

	// 5. 执行数据库迁移和种子数据
	fmt.Fprintln(os.Stderr, "DEBUG: DB connections retrieved. Starting migrations...")
	migrateMetadataDatabase(metadataDB)
	migrateUserDatabase(userDB)
	migrateSSODatabase(userDB)
	migrateAuditDatabase(auditDB)

	// 5.1 初始化 Casbin 权限管理
	fmt.Fprintln(os.Stderr, "DEBUG: Migrations completed. Initializing Casbin...")
	if err := middleware.InitCasbin(userDB, "./configs/rbac_model.conf"); err != nil {
		utils.SugarLogger.Warnf("Failed to initialize Casbin: %v, continuing without RBAC", err)
	}

	fmt.Fprintln(os.Stderr, "DEBUG: Seeding completed. Starting seeding...")
	seedMetadataDatabase(metadataDB, cfg)
	seedUserDatabase(userDB)
	// 6. 初始化Hertz引擎
	fmt.Fprintln(os.Stderr, "DEBUG: Seeding completed. Initializing Hertz server...")
	r := server.Default(
		server.WithHostPorts(fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)),
	)

	// 6.4 初始化审计日志队列
	fmt.Fprintln(os.Stderr, "DEBUG: Hertz initialized. Starting audit log queue...")
	auditLogQueue := auditQueuePkg.NewAuditLogQueue(auditDB, 1000, 5)
	auditLogQueue.Start()
	defer auditLogQueue.Stop()

	// 6.5 执行权限同步
	fmt.Fprintln(os.Stderr, "DEBUG: Audit queue started. Syncing Casbin policies...")
	syncCasbinPolicies(userDB, auditDB, auditLogQueue)

	// ...

	// 7. 初始化中间件加载器
	fmt.Fprintln(os.Stderr, "DEBUG: Casbin synced. Loading middlewares...")
	middlewareLoader := middleware.NewMiddlewareLoader()
	middlewareLoader.RegisterDefaultMiddlewares()
	middlewareLoader.LoadMiddlewareConfig(middleware.GetDefaultMiddlewareConfig())
	middlewareLoader.UseMiddlewareChain(r)

	// 8. 注册模块化路由
	fmt.Fprintln(os.Stderr, "DEBUG: Middlewares loaded. Registering routes...")
	api.RegisterModuleRoutes(r, metadataDB, userDB, auditDB, auditLogQueue)

	// 9. 启动服务器
	fmt.Fprintln(os.Stderr, "DEBUG: Routes registered. Starting server run loop...")
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

func migrateSSODatabase(db *gorm.DB) {
	if err := sso.Migrate(db); err != nil {
		utils.SugarLogger.Errorf("Failed to migrate SSO database: %v", err)
	} else {
		utils.SugarLogger.Info("SSO database migration completed successfully")
	}
}

func seedMetadataDatabase(db *gorm.DB, cfg *configs.Config) {
	metadata.SeedData(db, cfg)
}

func seedUserDatabase(db *gorm.DB) {
	user.SeedData(db)
}

// syncCasbinPolicies 系统启动时同步 Casbin 策略
func syncCasbinPolicies(db *gorm.DB, auditDB *gorm.DB, auditQueue *auditQueuePkg.AuditLogQueue) {
	// 初始化内部服务来执行同步
	repos := user.GetRepositories(db)
	services := user.GetServices(repos, db, auditDB, auditQueue)
	if err := services.CasbinSync.SyncAll(); err != nil {
		utils.SugarLogger.Errorf("Failed to sync Casbin policies: %v", err)
	}
}
