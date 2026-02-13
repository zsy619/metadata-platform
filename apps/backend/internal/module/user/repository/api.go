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
	GetUserWithDetails(id string) (*model.SsoUser, error)
	UpdateUser(user *model.SsoUser) error
	UpdateLoginInfo(id string, ip string) error
	IncrementLoginError(id string) error
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

// SsoAppRepository 应用仓库接口
type SsoAppRepository interface {
	CreateApp(app *model.SsoApp) error
	GetAppByID(id string) (*model.SsoApp, error)
	GetAppByCode(code string) (*model.SsoApp, error)
	UpdateApp(app *model.SsoApp) error
	DeleteApp(id string) error
	GetAllApps() ([]model.SsoApp, error)
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

// SsoOrgRepository 组织仓库接口
type SsoOrgRepository interface {
	CreateOrg(org *model.SsoOrg) error
	GetOrgByID(id string) (*model.SsoOrg, error)
	GetOrgByCode(code string) (*model.SsoOrg, error)
	UpdateOrg(org *model.SsoOrg) error
	DeleteOrg(id string) error
	GetAllOrgs() ([]model.SsoOrg, error)
}

// SsoPosRepository 职位仓库接口
type SsoPosRepository interface {
	CreatePos(pos *model.SsoPos) error
	GetPosByID(id string) (*model.SsoPos, error)
	GetPosByCode(code string) (*model.SsoPos, error)
	UpdatePos(pos *model.SsoPos) error
	DeletePos(id string) error
	GetAllPoss() ([]model.SsoPos, error)
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
type SsoUserPosRepository interface {
	CreateUserPos(userPos *model.SsoUserPos) error
	GetUserPosByID(id string) (*model.SsoUserPos, error)
	GetUserPosByUserID(userID string) ([]model.SsoUserPos, error)
	GetUserPosByPosID(posID string) ([]model.SsoUserPos, error)
	DeleteUserPos(id string) error
	DeleteUserPosByUserID(userID string) error
	DeleteUserPosByPosID(posID string) error
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
type SsoPosRoleRepository interface {
	CreatePosRole(posRole *model.SsoPosRole) error
	GetPosRoleByID(id string) (*model.SsoPosRole, error)
	GetPosRolesByPosID(posID string) ([]model.SsoPosRole, error)
	GetPosRolesByRoleID(roleID string) ([]model.SsoPosRole, error)
	DeletePosRole(id string) error
	DeletePosRolesByPosID(posID string) error
	DeletePosRolesByRoleID(roleID string) error
}

// Repositories 用户模块仓库集合
type Repositories struct {
	User     SsoUserRepository
	Tenant   SsoTenantRepository
	App      SsoAppRepository
	Menu     SsoMenuRepository
	Role     SsoRoleRepository
	Org      SsoOrgRepository
	Pos      SsoPosRepository
	UserRole SsoUserRoleRepository
	UserPos  SsoUserPosRepository
	RoleMenu SsoRoleMenuRepository
	PosRole  SsoPosRoleRepository
}

// NewRepositories 创建用户模块仓库集合
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:     NewSsoUserRepository(db),
		Tenant:   NewSsoTenantRepository(db),
		App:      NewSsoAppRepository(db),
		Menu:     NewSsoMenuRepository(db),
		Role:     NewSsoRoleRepository(db),
		Org:      NewSsoOrgRepository(db),
		Pos:      NewSsoPosRepository(db),
		UserRole: NewSsoUserRoleRepository(db),
		UserPos:  NewSsoUserPosRepository(db),
		RoleMenu: NewSsoRoleMenuRepository(db),
		PosRole:  NewSsoPosRoleRepository(db),
	}
}
