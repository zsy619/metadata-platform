package api

import (
	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/module/user/service"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"

	globalMiddleware "metadata-platform/internal/middleware"

	auditSvc "metadata-platform/internal/module/audit/service"
)

// SsoHandler 用户模块处理器集合
type SsoHandler struct {
	TenantHandler      *SsoTenantHandler
	AppHandler         *SsoAppHandler
	MenuHandler        *SsoMenuHandler
	RoleHandler        *SsoRoleHandler
	OrgHandler         *SsoOrgHandler
	UserHandler        *SsoUserHandler
	PosHandler         *SsoPosHandler
	AuthHandler        *SsoAuthHandler
	OrgKindHandler     *SsoOrgKindHandler
	RoleGroupHandler   *SsoRoleGroupHandler
	UserGroupHandler   *SsoUserGroupHandler
	UserProfileHandler *SsoUserProfileHandler
	UserAddressHandler *SsoUserAddressHandler
	UserContactHandler *SsoUserContactHandler
	UserSocialHandler  *SsoUserSocialHandler
	DashboardHandler   *DashboardHandler
	AuditService       auditSvc.AuditService
}

// NewSsoHandler 创建用户模块处理器集合
func NewSsoHandler(db *gorm.DB, auditDB *gorm.DB, services *service.Services, auditQueue *queue.AuditLogQueue, repos *repository.Repositories) *SsoHandler {
	return &SsoHandler{
		UserHandler:        NewSsoUserHandler(services.User, auditQueue),
		TenantHandler:      NewSsoTenantHandler(services.Tenant, auditQueue),
		AppHandler:         NewSsoAppHandler(services.App, auditQueue),
		MenuHandler:        NewSsoMenuHandler(services.Menu, auditQueue),
		RoleHandler:        NewSsoRoleHandler(services.Role, auditQueue),
		OrgHandler:         NewSsoOrgHandler(services.Org, auditQueue),
		PosHandler:         NewSsoPosHandler(services.Pos, auditQueue),
		AuthHandler:        NewSsoAuthHandler(services.Auth),
		OrgKindHandler:     NewSsoOrgKindHandler(services.OrgKind, auditQueue),
		RoleGroupHandler:   NewSsoRoleGroupHandler(services.RoleGroup, auditQueue),
		UserGroupHandler:   NewSsoUserGroupHandler(services.UserGroup, auditQueue),
		UserProfileHandler: NewSsoUserProfileHandler(services.UserProfile, auditQueue),
		UserAddressHandler: NewSsoUserAddressHandler(services.UserAddress, auditQueue),
		UserContactHandler: NewSsoUserContactHandler(services.UserContact, auditQueue),
		UserSocialHandler:  NewSsoUserSocialHandler(services.UserSocial, auditQueue),
		DashboardHandler:   NewDashboardHandler(db, auditDB, repos, services.Audit),
		AuditService:       services.Audit,
	}
}

// RegisterRoutes 注册路由
func (h *SsoHandler) RegisterRoutes(router *server.Hertz) {
	// 认证相关路由
	authRouter := router.Group("/api/auth")
	{
		authRouter.POST("/login", h.AuthHandler.Login)
		authRouter.POST("/refresh", h.AuthHandler.Refresh)
		authRouter.POST("/logout", globalMiddleware.AuthMiddleware(), h.AuthHandler.Logout)
		authRouter.POST("/password", globalMiddleware.AuthMiddleware(), h.AuthHandler.ChangePassword)
		authRouter.GET("/profile", globalMiddleware.AuthMiddleware(), h.AuthHandler.GetProfile)
		authRouter.GET("/captcha", h.AuthHandler.GetCaptcha)
	}

	group := router.Group("/api/sso")
	group.Use(globalMiddleware.TenantMiddleware())
	group.Use(globalMiddleware.AuthMiddleware())
	group.Use(globalMiddleware.AuditMiddleware(h.AuditService, "sso"))
	// 用户相关路由
	userRouter := group.Group("/user")
	{
		userRouter.POST("", h.UserHandler.CreateUser)
		userRouter.GET("/:id", h.UserHandler.GetUserByID)
		userRouter.PUT("/:id", h.UserHandler.UpdateUser)
		userRouter.DELETE("/:id", h.UserHandler.DeleteUser)
		userRouter.GET("", h.UserHandler.GetAllUsers)
		// 用户关联路由
		userRouter.GET("/:id/roles", h.UserHandler.GetUserRoles)
		userRouter.PUT("/:id/roles", h.UserHandler.UpdateUserRoles)
		userRouter.GET("/:id/pos", h.UserHandler.GetUserPos)
		userRouter.PUT("/:id/pos", h.UserHandler.UpdateUserPos)
		userRouter.GET("/:id/groups", h.UserHandler.GetUserGroups)
		userRouter.PUT("/:id/groups", h.UserHandler.UpdateUserGroups)
		userRouter.GET("/:id/role-groups", h.UserHandler.GetUserRoleGroups)
		userRouter.PUT("/:id/role-groups", h.UserHandler.UpdateUserRoleGroups)
		userRouter.GET("/:id/orgs", h.UserHandler.GetUserOrgs)
		userRouter.PUT("/:id/orgs", h.UserHandler.UpdateUserOrgs)
		// 用户扩展信息路由
		userRouter.GET("/:id/profile", h.UserProfileHandler.GetUserProfile)
		userRouter.PUT("/:id/profile", h.UserProfileHandler.UpsertUserProfile)
		userRouter.GET("/:id/addresses", h.UserAddressHandler.GetAddresses)
		userRouter.POST("/:id/addresses", h.UserAddressHandler.CreateAddress)
		userRouter.PUT("/:id/addresses/:aid", h.UserAddressHandler.UpdateAddress)
		userRouter.DELETE("/:id/addresses/:aid", h.UserAddressHandler.DeleteAddress)
		userRouter.PUT("/:id/addresses/:aid/default", h.UserAddressHandler.SetDefaultAddress)
		userRouter.GET("/:id/contacts", h.UserContactHandler.GetContacts)
		userRouter.POST("/:id/contacts", h.UserContactHandler.CreateContact)
		userRouter.PUT("/:id/contacts/:cid", h.UserContactHandler.UpdateContact)
		userRouter.DELETE("/:id/contacts/:cid", h.UserContactHandler.DeleteContact)
		userRouter.GET("/:id/socials", h.UserSocialHandler.GetSocials)
		userRouter.POST("/:id/socials", h.UserSocialHandler.BindSocial)
		userRouter.DELETE("/:id/socials/:sid", h.UserSocialHandler.UnbindSocial)
	}

	// 租户相关路由
	tenantRouter := group.Group("/tenant")
	{
		tenantRouter.POST("", h.TenantHandler.CreateTenant)
		tenantRouter.GET("/:id", h.TenantHandler.GetTenantByID)
		tenantRouter.PUT("/:id", h.TenantHandler.UpdateTenant)
		tenantRouter.DELETE("/:id", h.TenantHandler.DeleteTenant)
		tenantRouter.GET("", h.TenantHandler.GetAllTenants)
	}

	// 应用相关路由
	appRouter := group.Group("/app")
	{
		appRouter.POST("", h.AppHandler.CreateApp)
		appRouter.GET("/:id", h.AppHandler.GetAppByID)
		appRouter.PUT("/:id", h.AppHandler.UpdateApp)
		appRouter.DELETE("/:id", h.AppHandler.DeleteApp)
		appRouter.GET("", h.AppHandler.GetAllApps)
	}

	// 菜单相关路由
	menuRouter := group.Group("/menu")
	{
		menuRouter.POST("", h.MenuHandler.CreateMenu)
		menuRouter.GET("/:id", h.MenuHandler.GetMenuByID)
		menuRouter.PUT("/:id", h.MenuHandler.UpdateMenu)
		menuRouter.DELETE("/:id", h.MenuHandler.DeleteMenu)
		menuRouter.GET("", h.MenuHandler.GetAllMenus)
	}

	// 角色相关路由
	roleRouter := group.Group("/role")
	{
		roleRouter.POST("", h.RoleHandler.CreateRole)
		roleRouter.GET("/:id", h.RoleHandler.GetRoleByID)
		roleRouter.PUT("/:id", h.RoleHandler.UpdateRole)
		roleRouter.DELETE("/:id", h.RoleHandler.DeleteRole)
		roleRouter.GET("", h.RoleHandler.GetAllRoles)
		// 角色菜单相关路由
		roleRouter.GET("/:id/menus", h.RoleHandler.GetRoleMenus)
		roleRouter.PUT("/:id/menus", h.RoleHandler.UpdateRoleMenus)
	}

	// 组织相关路由
	orgRouter := group.Group("/org")
	{
		orgRouter.POST("", h.OrgHandler.CreateOrg)
		orgRouter.GET("/:id", h.OrgHandler.GetOrgByID)
		orgRouter.PUT("/:id", h.OrgHandler.UpdateOrg)
		orgRouter.DELETE("/:id", h.OrgHandler.DeleteOrg)
		orgRouter.GET("", h.OrgHandler.GetAllOrgs)
	}

	// 组织类型相关路由
	orgKindRouter := group.Group("/org-kind")
	{
		orgKindRouter.POST("", h.OrgKindHandler.CreateOrgKind)
		orgKindRouter.GET("/:id", h.OrgKindHandler.GetOrgKindByID)
		orgKindRouter.PUT("/:id", h.OrgKindHandler.UpdateOrgKind)
		orgKindRouter.DELETE("/:id", h.OrgKindHandler.DeleteOrgKind)
		orgKindRouter.GET("", h.OrgKindHandler.GetAllOrgKinds)
	}

	// 职位相关路由
	posRouter := group.Group("/pos")
	{
		posRouter.POST("", h.PosHandler.CreatePos)
		posRouter.GET("/:id", h.PosHandler.GetPosByID)
		posRouter.PUT("/:id", h.PosHandler.UpdatePos)
		posRouter.DELETE("/:id", h.PosHandler.DeletePos)
		posRouter.GET("", h.PosHandler.GetAllPoss)
		// 职位角色关联路由
		posRouter.GET("/:id/roles", h.PosHandler.GetPosRoles)
		posRouter.PUT("/:id/roles", h.PosHandler.UpdatePosRoles)
	}

	// 角色分组相关路由
	roleGroupRouter := group.Group("/role-group")
	{
		roleGroupRouter.POST("", h.RoleGroupHandler.CreateRoleGroup)
		roleGroupRouter.GET("/:id", h.RoleGroupHandler.GetRoleGroupByID)
		roleGroupRouter.PUT("/:id", h.RoleGroupHandler.UpdateRoleGroup)
		roleGroupRouter.DELETE("/:id", h.RoleGroupHandler.DeleteRoleGroup)
		roleGroupRouter.GET("", h.RoleGroupHandler.GetAllRoleGroups)
		// 角色组角色关联路由
		roleGroupRouter.GET("/:id/roles", h.RoleGroupHandler.GetRoleGroupRoles)
		roleGroupRouter.PUT("/:id/roles", h.RoleGroupHandler.UpdateRoleGroupRoles)
	}

	// 用户组相关路由
	userGroupRouter := group.Group("/user-group")
	{
		userGroupRouter.POST("", h.UserGroupHandler.CreateUserGroup)
		userGroupRouter.GET("/:id", h.UserGroupHandler.GetUserGroupByID)
		userGroupRouter.PUT("/:id", h.UserGroupHandler.UpdateUserGroup)
		userGroupRouter.DELETE("/:id", h.UserGroupHandler.DeleteUserGroup)
		userGroupRouter.GET("", h.UserGroupHandler.GetAllUserGroups)
		// 用户组角色关联路由
		userGroupRouter.GET("/:id/roles", h.UserGroupHandler.GetUserGroupRoles)
		userGroupRouter.PUT("/:id/roles", h.UserGroupHandler.UpdateUserGroupRoles)
	}

	// 仪表板相关路由
	dashboardRouter := group.Group("/dashboard")
	{
		dashboardRouter.GET("/stats", h.DashboardHandler.GetStats)
		dashboardRouter.GET("/login-logs", h.DashboardHandler.GetRecentLoginLogs)
		dashboardRouter.GET("/operation-logs", h.DashboardHandler.GetRecentOperationLogs)
		dashboardRouter.GET("/login-trend", h.DashboardHandler.GetLoginTrend)
		dashboardRouter.GET("/user-status", h.DashboardHandler.GetUserStatusDistribution)
		dashboardRouter.GET("/operation-stats", h.DashboardHandler.GetOperationStats)
		dashboardRouter.GET("/org-distribution", h.DashboardHandler.GetOrgDistribution)
	}
}
