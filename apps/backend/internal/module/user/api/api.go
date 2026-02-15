package api

import (
	"metadata-platform/internal/middleware"
	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/module/user/service"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// SsoHandler 用户模块处理器集合
type SsoHandler struct {
	TenantHandler  *SsoTenantHandler
	AppHandler     *SsoAppHandler
	MenuHandler    *SsoMenuHandler
	RoleHandler    *SsoRoleHandler
	OrgHandler     *SsoOrgHandler
	UserHandler    *SsoUserHandler
	PosHandler     *SsoPosHandler
	AuthHandler    *SsoAuthHandler
	OrgKindHandler *service.SsoOrgKindHandler
}

// NewSsoHandler 创建用户模块处理器集合
func NewSsoHandler(services *service.Services, auditQueue *queue.AuditLogQueue) *SsoHandler {
	return &SsoHandler{
		UserHandler:   NewSsoUserHandler(services.User),
		TenantHandler: NewSsoTenantHandler(services.Tenant, auditQueue),
		AppHandler:    NewSsoAppHandler(services.App, auditQueue),
		MenuHandler:   NewSsoMenuHandler(services.Menu, auditQueue),
		RoleHandler:   NewSsoRoleHandler(services.Role, auditQueue),
		OrgHandler:    NewSsoOrgHandler(services.Org, auditQueue),
		PosHandler:    NewSsoPosHandler(services.Pos, auditQueue),
		AuthHandler:   NewSsoAuthHandler(services.Auth),
	}
}

// RegisterRoutes 注册路由
func (h *SsoHandler) RegisterRoutes(router *server.Hertz, orgKindHandler *service.SsoOrgKindHandler) {
	// 用户相关路由
	userRouter := router.Group("/api/user")
	{
		userRouter.POST("", h.UserHandler.CreateUser)
		userRouter.GET("/:id", h.UserHandler.GetUserByID)
		userRouter.PUT("/:id", h.UserHandler.UpdateUser)
		userRouter.DELETE("/:id", h.UserHandler.DeleteUser)
		userRouter.GET("", h.UserHandler.GetAllUsers)
	}

	// 认证相关路由
	authRouter := router.Group("/api/auth")
	{
		authRouter.POST("/login", h.AuthHandler.Login)
		authRouter.POST("/refresh", h.AuthHandler.Refresh)
		authRouter.POST("/logout", middleware.AuthMiddleware(), h.AuthHandler.Logout)
		authRouter.POST("/password", middleware.AuthMiddleware(), h.AuthHandler.ChangePassword)
		authRouter.GET("/profile", middleware.AuthMiddleware(), h.AuthHandler.GetProfile)
		authRouter.GET("/captcha", h.AuthHandler.GetCaptcha)
	}

	// 租户相关路由
	tenantRouter := router.Group("/api/tenant")
	{
		tenantRouter.POST("", h.TenantHandler.CreateTenant)
		tenantRouter.GET("/:id", h.TenantHandler.GetTenantByID)
		tenantRouter.PUT("/:id", h.TenantHandler.UpdateTenant)
		tenantRouter.DELETE("/:id", h.TenantHandler.DeleteTenant)
		tenantRouter.GET("", h.TenantHandler.GetAllTenants)
	}

	// 应用相关路由
	appRouter := router.Group("/api/app")
	{
		appRouter.POST("", h.AppHandler.CreateApp)
		appRouter.GET("/:id", h.AppHandler.GetAppByID)
		appRouter.PUT("/:id", h.AppHandler.UpdateApp)
		appRouter.DELETE("/:id", h.AppHandler.DeleteApp)
		appRouter.GET("", h.AppHandler.GetAllApps)
	}

	// 菜单相关路由
	menuRouter := router.Group("/api/menu")
	{
		menuRouter.POST("", h.MenuHandler.CreateMenu)
		menuRouter.GET("/:id", h.MenuHandler.GetMenuByID)
		menuRouter.PUT("/:id", h.MenuHandler.UpdateMenu)
		menuRouter.DELETE("/:id", h.MenuHandler.DeleteMenu)
		menuRouter.GET("", h.MenuHandler.GetAllMenus)
	}

	// 角色相关路由
	roleRouter := router.Group("/api/role")
	{
		roleRouter.POST("", h.RoleHandler.CreateRole)
		roleRouter.GET("/:id", h.RoleHandler.GetRoleByID)
		roleRouter.PUT("/:id", h.RoleHandler.UpdateRole)
		roleRouter.DELETE("/:id", h.RoleHandler.DeleteRole)
		roleRouter.GET("", h.RoleHandler.GetAllRoles)
	}

	// 组织相关路由
	orgRouter := router.Group("/api/org")
	{
		orgRouter.POST("", h.OrgHandler.CreateOrg)
		orgRouter.GET("/:id", h.OrgHandler.GetOrgByID)
		orgRouter.PUT("/:id", h.OrgHandler.UpdateOrg)
		orgRouter.DELETE("/:id", h.OrgHandler.DeleteOrg)
		orgRouter.GET("", h.OrgHandler.GetAllOrgs)
	}

	// 组织类型相关路由
	orgKindRouter := router.Group("/api/org-kind")
	{
		orgKindRouter.POST("", orgKindHandler.CreateOrgKind)
		orgKindRouter.GET("/:id", orgKindHandler.GetOrgKindByID)
		orgKindRouter.PUT("/:id", orgKindHandler.UpdateOrgKind)
		orgKindRouter.DELETE("/:id", orgKindHandler.DeleteOrgKind)
		orgKindRouter.GET("", orgKindHandler.GetAllOrgKinds)
	}

	// 职位相关路由
	posRouter := router.Group("/api/pos")
	{
		posRouter.POST("", h.PosHandler.CreatePos)
		posRouter.GET("/:id", h.PosHandler.GetPosByID)
		posRouter.PUT("/:id", h.PosHandler.UpdatePos)
		posRouter.DELETE("/:id", h.PosHandler.DeletePos)
		posRouter.GET("", h.PosHandler.GetAllPoss)
	}
}
