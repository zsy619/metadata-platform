package api

import (
	"metadata-platform/internal/middleware"
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
	OrgKindHandler *service.OrgKindHandler
}

// NewSsoHandler 创建用户模块处理器集合
func NewSsoHandler(services *service.Services) *SsoHandler {
	return &SsoHandler{
		UserHandler:   NewSsoUserHandler(services.User),
		TenantHandler: NewSsoTenantHandler(services.Tenant),
		AppHandler:    NewSsoAppHandler(services.App),
		MenuHandler:   NewSsoMenuHandler(services.Menu),
		RoleHandler:   NewSsoRoleHandler(services.Role),
		OrgHandler:    NewSsoOrgHandler(services.Org),
		PosHandler:    NewSsoPosHandler(services.Pos),
		AuthHandler:   NewSsoAuthHandler(services.Auth),
	}
}

// RegisterRoutes 注册路由
func (h *SsoHandler) RegisterRoutes(router *server.Hertz, orgKindHandler *service.OrgKindHandler) {
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
	applicationRouter := router.Group("/api/app")
	{
		applicationRouter.POST("", h.AppHandler.CreateApp)
		applicationRouter.GET("/:id", h.AppHandler.GetAppByID)
		applicationRouter.PUT("/:id", h.AppHandler.UpdateApp)
		applicationRouter.DELETE("/:id", h.AppHandler.DeleteApp)
		applicationRouter.GET("", h.AppHandler.GetAllApps)
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
	organizationRouter := router.Group("/api/unit")
	{
		organizationRouter.POST("", h.OrgHandler.CreateOrg)
		organizationRouter.GET("/:id", h.OrgHandler.GetOrgByID)
		organizationRouter.PUT("/:id", h.OrgHandler.UpdateOrg)
		organizationRouter.DELETE("/:id", h.OrgHandler.DeleteOrg)
		organizationRouter.GET("", h.OrgHandler.GetAllOrgs)
	}

	// 职位相关路由
	positionRouter := router.Group("/api/pos")
	{
		positionRouter.POST("", h.PosHandler.CreatePos)
		positionRouter.GET("/:id", h.PosHandler.GetPosByID)
		positionRouter.PUT("/:id", h.PosHandler.UpdatePos)
		positionRouter.DELETE("/:id", h.PosHandler.DeletePos)
		positionRouter.GET("", h.PosHandler.GetAllPoss)
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
}
