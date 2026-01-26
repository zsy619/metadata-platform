package service

import (
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
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

// SsoApplicationService 应用服务接口
type SsoApplicationService interface {
	CreateApplication(app *model.SsoApplication) error
	GetApplicationByID(id string) (*model.SsoApplication, error)
	GetApplicationByCode(code string) (*model.SsoApplication, error)
	UpdateApplication(app *model.SsoApplication) error
	DeleteApplication(id string) error
	GetAllApplications() ([]model.SsoApplication, error)
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

// SsoOrganizationService 组织服务接口
type SsoOrganizationService interface {
	CreateOrganization(unit *model.SsoOrganization) error
	GetOrganizationByID(id string) (*model.SsoOrganization, error)
	GetOrganizationByCode(code string) (*model.SsoOrganization, error)
	UpdateOrganization(unit *model.SsoOrganization) error
	DeleteOrganization(id string) error
	GetAllOrganizations() ([]model.SsoOrganization, error)
}

// SsoPositionService 职位服务接口
type SsoPositionService interface {
	CreatePosition(pos *model.SsoPosition) error
	GetPositionByID(id string) (*model.SsoPosition, error)
	GetPositionByCode(code string) (*model.SsoPosition, error)
	UpdatePosition(pos *model.SsoPosition) error
	DeletePosition(id string) error
	GetAllPositions() ([]model.SsoPosition, error)
}

// SsoAuthService 认证服务接口
type SsoAuthService interface {
	Login(account string, password string, tenantID uint) (accessToken string, refreshToken string, err error)
	Refresh(refreshToken string) (newAccessToken string, err error)
	GetUserInfo(userID string) (*model.SsoUser, error)
}

// Services 用户模块服务集合
type Services struct {
	User         SsoUserService
	Tenant       SsoTenantService
	Application  SsoApplicationService
	Menu         SsoMenuService
	Role         SsoRoleService
	Organization SsoOrganizationService
	Position     SsoPositionService
	Auth         SsoAuthService
	CasbinSync   SsoCasbinSyncService
}

// NewServices 创建用户模块服务集合
func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		User:         NewSsoUserService(repos.User),
		Tenant:       NewSsoTenantService(repos.Tenant),
		Application:  NewSsoApplicationService(repos.Application),
		Menu:         NewSsoMenuService(repos.Menu),
		Role:         NewSsoRoleService(repos.Role),
		Organization: NewSsoOrganizationService(repos.Organization),
		Position:     NewSsoPositionService(repos.Position),
		Auth:         NewSsoAuthService(repos.User),
		CasbinSync: NewSsoCasbinSyncService(
			repos.UserRole,
			repos.RoleMenu,
			repos.Role,
			repos.Menu,
		),
	}
}
