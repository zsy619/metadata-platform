package api

import (
	"metadata-platform/internal/module/user/service"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// SsoHandler 用户模块处理器集合
type SsoHandler struct {
	UserHandler         *SsoUserHandler
	TenantHandler       *SsoTenantHandler
	ApplicationHandler  *SsoApplicationHandler
	MenuHandler         *SsoMenuHandler
	RoleHandler         *SsoRoleHandler
	OrganizationHandler *SsoOrganizationHandler
	PositionHandler     *SsoPositionHandler
	AuthHandler         *SsoAuthHandler
}

// NewSsoHandler 创建用户模块处理器集合
func NewSsoHandler(services *service.Services) *SsoHandler {
	return &SsoHandler{
		UserHandler:         NewSsoUserHandler(services.User),
		TenantHandler:       NewSsoTenantHandler(services.Tenant),
		ApplicationHandler:  NewSsoApplicationHandler(services.Application),
		MenuHandler:         NewSsoMenuHandler(services.Menu),
		RoleHandler:         NewSsoRoleHandler(services.Role),
		OrganizationHandler: NewSsoOrganizationHandler(services.Organization),
		PositionHandler:     NewSsoPositionHandler(services.Position),
		AuthHandler:         NewSsoAuthHandler(services.Auth),
	}
}

// RegisterRoutes 注册路由
func (h *SsoHandler) RegisterRoutes(router *server.Hertz) {
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
		authRouter.POST("/logout", h.AuthHandler.Logout)
		authRouter.GET("/profile", h.AuthHandler.GetProfile)
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
		applicationRouter.POST("", h.ApplicationHandler.CreateApplication)
		applicationRouter.GET("/:id", h.ApplicationHandler.GetApplicationByID)
		applicationRouter.PUT("/:id", h.ApplicationHandler.UpdateApplication)
		applicationRouter.DELETE("/:id", h.ApplicationHandler.DeleteApplication)
		applicationRouter.GET("", h.ApplicationHandler.GetAllApplications)
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
		organizationRouter.POST("", h.OrganizationHandler.CreateOrganization)
		organizationRouter.GET("/:id", h.OrganizationHandler.GetOrganizationByID)
		organizationRouter.PUT("/:id", h.OrganizationHandler.UpdateOrganization)
		organizationRouter.DELETE("/:id", h.OrganizationHandler.DeleteOrganization)
		organizationRouter.GET("", h.OrganizationHandler.GetAllOrganizations)
	}

	// 职位相关路由
	positionRouter := router.Group("/api/pos")
	{
		positionRouter.POST("", h.PositionHandler.CreatePosition)
		positionRouter.GET("/:id", h.PositionHandler.GetPositionByID)
		positionRouter.PUT("/:id", h.PositionHandler.UpdatePosition)
		positionRouter.DELETE("/:id", h.PositionHandler.DeletePosition)
		positionRouter.GET("", h.PositionHandler.GetAllPositions)
	}
}
