package repository

import (
	"metadata-platform/internal/module/user/model"

	"gorm.io/gorm"
)

// SsoUserRepository 用户仓库接口
type SsoUserRepository interface {
	CreateUser(user *model.SsoUser) error
	GetUserByID(id string) (*model.SsoUser, error)
	GetUserByAccount(account string) (*model.SsoUser, error)
	UpdateUser(user *model.SsoUser) error
	DeleteUser(id string) error
	GetAllUsers() ([]model.SsoUser, error)
}

// SsoTenantRepository 租户仓库接口
type SsoTenantRepository interface {
	CreateTenant(tenant *model.SsoTenant) error
	GetTenantByID(id string) (*model.SsoTenant, error)
	GetTenantByCode(code string) (*model.SsoTenant, error)
	UpdateTenant(tenant *model.SsoTenant) error
	DeleteTenant(id string) error
	GetAllTenants() ([]model.SsoTenant, error)
}

// SsoApplicationRepository 应用仓库接口
type SsoApplicationRepository interface {
	CreateApplication(app *model.SsoApplication) error
	GetApplicationByID(id string) (*model.SsoApplication, error)
	GetApplicationByCode(code string) (*model.SsoApplication, error)
	UpdateApplication(app *model.SsoApplication) error
	DeleteApplication(id string) error
	GetAllApplications() ([]model.SsoApplication, error)
}

// SsoMenuRepository 菜单仓库接口
type SsoMenuRepository interface {
	CreateMenu(menu *model.SsoMenu) error
	GetMenuByID(id string) (*model.SsoMenu, error)
	GetMenuByCode(code string) (*model.SsoMenu, error)
	UpdateMenu(menu *model.SsoMenu) error
	DeleteMenu(id string) error
	GetAllMenus() ([]model.SsoMenu, error)
}

// SsoRoleRepository 角色仓库接口
type SsoRoleRepository interface {
	CreateRole(role *model.SsoRole) error
	GetRoleByID(id string) (*model.SsoRole, error)
	GetRoleByCode(code string) (*model.SsoRole, error)
	UpdateRole(role *model.SsoRole) error
	DeleteRole(id string) error
	GetAllRoles() ([]model.SsoRole, error)
}

// SsoOrganizationRepository 组织仓库接口
type SsoOrganizationRepository interface {
	CreateOrganization(unit *model.SsoOrganization) error
	GetOrganizationByID(id string) (*model.SsoOrganization, error)
	GetOrganizationByCode(code string) (*model.SsoOrganization, error)
	UpdateOrganization(unit *model.SsoOrganization) error
	DeleteOrganization(id string) error
	GetAllOrganizations() ([]model.SsoOrganization, error)
}

// SsoPositionRepository 职位仓库接口
type SsoPositionRepository interface {
	CreatePosition(pos *model.SsoPosition) error
	GetPositionByID(id string) (*model.SsoPosition, error)
	GetPositionByCode(code string) (*model.SsoPosition, error)
	UpdatePosition(pos *model.SsoPosition) error
	DeletePosition(id string) error
	GetAllPositions() ([]model.SsoPosition, error)
}

// SsoUserRoleRepository 用户角色仓库接口
type SsoUserRoleRepository interface {
	CreateUserRole(userRole *model.SsoUserRole) error
	GetUserRoleByID(id string) (*model.SsoUserRole, error)
	GetUserRolesByUserID(userID string) ([]model.SsoUserRole, error)
	GetUserRolesByRoleID(roleID string) ([]model.SsoUserRole, error)
	GetAllUserRoles() ([]model.SsoUserRole, error)
	DeleteUserRole(id string) error
	DeleteUserRolesByUserID(userID string) error
	DeleteUserRolesByRoleID(roleID string) error
}

// SsoUserPositionRepository 用户职位仓库接口
type SsoUserPositionRepository interface {
	CreateUserPosition(userPos *model.SsoUserPosition) error
	GetUserPositionByID(id string) (*model.SsoUserPosition, error)
	GetUserPositionsByUserID(userID string) ([]model.SsoUserPosition, error)
	GetUserPositionsByPosID(posID string) ([]model.SsoUserPosition, error)
	DeleteUserPosition(id string) error
	DeleteUserPositionsByUserID(userID string) error
	DeleteUserPositionsByPosID(posID string) error
}

// SsoRoleMenuRepository 角色菜单仓库接口
type SsoRoleMenuRepository interface {
	CreateRoleMenu(roleMenu *model.SsoRoleMenu) error
	GetRoleMenuByID(id string) (*model.SsoRoleMenu, error)
	GetRoleMenusByRoleID(roleID string) ([]model.SsoRoleMenu, error)
	GetRoleMenusByMenuID(menuID string) ([]model.SsoRoleMenu, error)
	GetAllRoleMenus() ([]model.SsoRoleMenu, error)
	DeleteRoleMenu(id string) error
	DeleteRoleMenusByRoleID(roleID string) error
	DeleteRoleMenusByMenuID(menuID string) error
}

// SsoPositionRoleRepository 职位角色仓库接口
type SsoPositionRoleRepository interface {
	CreatePositionRole(posRole *model.SsoPositionRole) error
	GetPositionRoleByID(id string) (*model.SsoPositionRole, error)
	GetPositionRolesByPosID(posID string) ([]model.SsoPositionRole, error)
	GetPositionRolesByRoleID(roleID string) ([]model.SsoPositionRole, error)
	DeletePositionRole(id string) error
	DeletePositionRolesByPosID(posID string) error
	DeletePositionRolesByRoleID(roleID string) error
}

// Repositories 用户模块仓库集合
type Repositories struct {
	User         SsoUserRepository
	Tenant       SsoTenantRepository
	Application  SsoApplicationRepository
	Menu         SsoMenuRepository
	Role         SsoRoleRepository
	Organization SsoOrganizationRepository
	Position     SsoPositionRepository
	UserRole     SsoUserRoleRepository
	UserPosition SsoUserPositionRepository
	RoleMenu     SsoRoleMenuRepository
	PositionRole SsoPositionRoleRepository
}

// NewRepositories 创建用户模块仓库集合
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:         NewSsoUserRepository(db),
		Tenant:       NewSsoTenantRepository(db),
		Application:  NewSsoApplicationRepository(db),
		Menu:         NewSsoMenuRepository(db),
		Role:         NewSsoRoleRepository(db),
		Organization: NewSsoOrganizationRepository(db),
		Position:     NewSsoPositionRepository(db),
		UserRole:     NewSsoUserRoleRepository(db),
		UserPosition: NewSsoUserPositionRepository(db),
		RoleMenu:     NewSsoRoleMenuRepository(db),
		PositionRole: NewSsoPositionRoleRepository(db),
	}
}
