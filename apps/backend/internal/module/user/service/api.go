package service

import (
	"context"
	"metadata-platform/internal/module/audit/queue"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"

	"gorm.io/gorm"

	auditService "metadata-platform/internal/module/audit/service"
)

// SsoUserService 用户服务接口
type SsoUserService interface {
	CreateUser(user *model.SsoUser) error
	GetUserByID(id string) (*model.SsoUser, error)
	GetUserByAccount(account string) (*model.SsoUser, error)
	UpdateUser(user *model.SsoUser) error
	DeleteUser(id string) error
	GetAllUsers() ([]model.SsoUser, error)
	Login(account, password string) (string, error)
}

// SsoTenantService 租户服务接口
type SsoTenantService interface {
	CreateTenant(tenant *model.SsoTenant) error
	GetTenantByID(id string) (*model.SsoTenant, error)
	GetTenantByCode(code string) (*model.SsoTenant, error)
	UpdateTenant(tenant *model.SsoTenant) error
	DeleteTenant(id string) error
	GetAllTenants() ([]model.SsoTenant, error)
}

// SsoAppService 应用服务接口
type SsoAppService interface {
	CreateApp(app *model.SsoApp) error
	GetAppByID(id string) (*model.SsoApp, error)
	GetAppByCode(code string) (*model.SsoApp, error)
	UpdateApp(app *model.SsoApp) error
	DeleteApp(id string) error
	GetAllApps() ([]model.SsoApp, error)
}

// SsoMenuService 菜单服务接口
type SsoMenuService interface {
	CreateMenu(menu *model.SsoMenu) error
	GetMenuByID(id string) (*model.SsoMenu, error)
	GetMenuByCode(code string) (*model.SsoMenu, error)
	UpdateMenu(menu *model.SsoMenu) error
	DeleteMenu(id string) error
	GetAllMenus() ([]model.SsoMenu, error)
}

// SsoRoleService 角色服务接口
type SsoRoleService interface {
	CreateRole(role *model.SsoRole) error
	GetRoleByID(id string) (*model.SsoRole, error)
	GetRoleByCode(code string) (*model.SsoRole, error)
	UpdateRole(role *model.SsoRole) error
	DeleteRole(id string) error
	GetAllRoles() ([]model.SsoRole, error)
}

// SsoOrgService 组织服务接口
type SsoOrgService interface {
	CreateOrg(org *model.SsoOrg) error
	GetOrgByID(id string) (*model.SsoOrg, error)
	GetOrgByCode(code string) (*model.SsoOrg, error)
	UpdateOrg(org *model.SsoOrg) error
	DeleteOrg(id string) error
	GetAllOrgs() ([]model.SsoOrg, error)
}

// SsoPosService 职位服务接口
type SsoPosService interface {
	CreatePos(pos *model.SsoPos) error
	GetPosByID(id string) (*model.SsoPos, error)
	GetPosByCode(code string) (*model.SsoPos, error)
	UpdatePos(pos *model.SsoPos) error
	DeletePos(id string) error
	GetAllPoss() ([]model.SsoPos, error)
}

// SsoAuthService 认证服务接口
type SsoAuthService interface {
	Login(account string, password string, tenantID uint, clientInfo utils.ClientInfo) (accessToken string, refreshToken string, err error)
	Logout(ctx context.Context, userID string, clientInfo utils.ClientInfo) error
	Refresh(refreshToken string) (newAccessToken string, err error)
	GetUserInfo(userID string) (*model.SsoUser, error)
	ChangePassword(userID string, oldPassword string, newPassword string) error
}

// Services 用户模块服务集合
type Services struct {
	User       SsoUserService
	Tenant     SsoTenantService
	App        SsoAppService
	Menu       SsoMenuService
	Role       SsoRoleService
	Org        SsoOrgService
	Pos        SsoPosService
	Auth       SsoAuthService
	CasbinSync SsoCasbinSyncService
	Audit      auditService.AuditService
}

// NewServices 创建用户模块服务集合
func NewServices(repos *repository.Repositories, auditDB *gorm.DB, auditQueue *queue.AuditLogQueue) *Services {
	auditSvc := auditService.NewAuditService(auditDB, auditQueue)
	return &Services{
		User:   NewSsoUserService(repos.User, repos.Org),
		Tenant: NewSsoTenantService(repos.Tenant),
		App:    NewSsoAppService(repos.App),
		Menu:   NewSsoMenuService(repos.Menu),
		Role:   NewSsoRoleService(repos.Role),
		Org:    NewSsoOrgService(repos.Org),
		Pos:    NewSsoPosService(repos.Pos),
		Auth:   NewSsoAuthService(repos.User, auditSvc),
		CasbinSync: NewSsoCasbinSyncService(
			repos.UserRole,
			repos.RoleMenu,
			repos.Role,
			repos.Menu,
		),
		Audit: auditSvc,
	}
}
