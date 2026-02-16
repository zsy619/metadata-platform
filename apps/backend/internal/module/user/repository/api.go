package repository

import (
	"gorm.io/gorm"

	"metadata-platform/internal/module/user/model"
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
	UpdateAppFields(id string, fields map[string]any) error
	DeleteApp(id string) error
	GetAllApps() ([]model.SsoApp, error)
	GetMaxSort() (int, error)
	HasChildren(parentID string) (bool, error)
}

// SsoMenuRepository 菜单仓库接口
type SsoMenuRepository interface {
	CreateMenu(menu *model.SsoMenu) error
	GetMenuByID(id string) (*model.SsoMenu, error)
	GetMenuByCode(code string) (*model.SsoMenu, error)
	UpdateMenu(menu *model.SsoMenu) error
	DeleteMenu(id string) error
	GetAllMenus() ([]model.SsoMenu, error)
	GetMaxSort() (int, error)
}

// SsoRoleRepository 角色仓库接口
type SsoRoleRepository interface {
	CreateRole(role *model.SsoRole) error
	GetRoleByID(id string) (*model.SsoRole, error)
	GetRoleByCode(code string) (*model.SsoRole, error)
	UpdateRole(role *model.SsoRole) error
	DeleteRole(id string) error
	GetAllRoles() ([]model.SsoRole, error)
	GetMaxSort() (int, error)
	GetRolesByUserID(userID string) ([]model.SsoRole, error)
}

// SsoOrgRepository 组织仓库接口
type SsoOrgRepository interface {
	CreateOrg(org *model.SsoOrg) error
	GetOrgByID(id string) (*model.SsoOrg, error)
	GetOrgByCode(code string) (*model.SsoOrg, error)
	UpdateOrg(org *model.SsoOrg) error
	UpdateOrgFields(id string, fields map[string]any) error
	DeleteOrg(id string) error
	GetAllOrgs() ([]model.SsoOrg, error)
	GetMaxSort() (int, error)
	HasChildren(parentID string) (bool, error)
}

type SsoOrgKindRepository interface {
	CreateOrgKind(orgKind *model.SsoOrgKind) error
	GetOrgKindByID(id string) (*model.SsoOrgKind, error)
	GetOrgKindByCode(code string) (*model.SsoOrgKind, error)
	UpdateOrgKind(orgKind *model.SsoOrgKind) error
	UpdateOrgKindFields(id string, fields map[string]any) error
	DeleteOrgKind(id string) error
	GetAllOrgKinds() ([]model.SsoOrgKind, error)
	GetMaxSort() (int, error)
	HasChildren(parentID string) (bool, error)
}

type SsoPosRepository interface {
	CreatePos(pos *model.SsoPos) error
	GetPosByID(id string) (*model.SsoPos, error)
	GetPosByCode(code string) (*model.SsoPos, error)
	UpdatePos(pos *model.SsoPos) error
	DeletePos(id string) error
	GetAllPoss() ([]model.SsoPos, error)
	GetMaxSort() (int, error)
}

type SsoRoleGroupRepository interface {
	CreateRoleGroup(item *model.SsoRoleGroup) error
	GetRoleGroupByID(id string) (*model.SsoRoleGroup, error)
	GetRoleGroupByCode(code string) (*model.SsoRoleGroup, error)
	UpdateRoleGroup(item *model.SsoRoleGroup) error
	UpdateRoleGroupFields(id string, fields map[string]any) error
	DeleteRoleGroup(id string) error
	HasChildren(parentID string) (bool, error)
	GetAllRoleGroups() ([]model.SsoRoleGroup, error)
	GetMaxSort() (int, error)
}

type SsoUserGroupRepository interface {
	CreateUserGroup(item *model.SsoUserGroup) error
	GetUserGroupByID(id string) (*model.SsoUserGroup, error)
	GetUserGroupByCode(code string) (*model.SsoUserGroup, error)
	UpdateUserGroup(item *model.SsoUserGroup) error
	UpdateUserGroupFields(id string, fields map[string]any) error
	DeleteUserGroup(id string) error
	HasChildren(parentID string) (bool, error)
	GetAllUserGroups() ([]model.SsoUserGroup, error)
	GetMaxSort() (int, error)
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

type SsoRoleGroupRoleRepository interface {
	CreateRoleGroupRole(item *model.SsoRoleGroupRole) error
	GetRoleGroupRoleByID(id string) (*model.SsoRoleGroupRole, error)
	GetRoleGroupRolesByGroupID(groupID string) ([]model.SsoRoleGroupRole, error)
	GetRoleGroupRolesByRoleID(roleID string) ([]model.SsoRoleGroupRole, error)
	DeleteRoleGroupRole(id string) error
	DeleteRoleGroupRolesByGroupID(groupID string) error
	DeleteRoleGroupRolesByRoleID(roleID string) error
}

type SsoUserGroupUserRepository interface {
	CreateUserGroupUser(item *model.SsoUserGroupUser) error
	GetUserGroupUserByID(id string) (*model.SsoUserGroupUser, error)
	GetUserGroupUsersByGroupID(groupID string) ([]model.SsoUserGroupUser, error)
	GetUserGroupUsersByUserID(userID string) ([]model.SsoUserGroupUser, error)
	DeleteUserGroupUser(id string) error
	DeleteUserGroupUsersByGroupID(groupID string) error
	DeleteUserGroupUsersByUserID(userID string) error
}

type SsoUserRoleGroupRepository interface {
	CreateUserRoleGroup(item *model.SsoUserRoleGroup) error
	GetUserRoleGroupByID(id string) (*model.SsoUserRoleGroup, error)
	GetUserRoleGroupsByUserID(userID string) ([]model.SsoUserRoleGroup, error)
	GetUserRoleGroupsByGroupID(groupID string) ([]model.SsoUserRoleGroup, error)
	DeleteUserRoleGroup(id string) error
	DeleteUserRoleGroupsByUserID(userID string) error
	DeleteUserRoleGroupsByGroupID(groupID string) error
}

type SsoCasbinRuleRepository interface {
	CreateCasbinRule(item *model.SsoCasbinRule) error
	GetCasbinRuleByID(id string) (*model.SsoCasbinRule, error)
	GetCasbinRulesByPType(ptype string) ([]model.SsoCasbinRule, error)
	GetCasbinRule(pType, v0, v1 string) (*model.SsoCasbinRule, error)
	DeleteCasbinRule(id string) error
	DeleteCasbinRulesByPType(ptype string) error
	DeleteCasbinRules(pType, v0, v1 string) error
	GetAllCasbinRules() ([]model.SsoCasbinRule, error)
}

type SsoOrgKindRoleRepository interface {
	CreateOrgKindRole(item *model.SsoOrgKindRole) error
	GetOrgKindRoleByID(id string) (*model.SsoOrgKindRole, error)
	GetOrgKindRoleByKindCode(kindCode string) ([]model.SsoOrgKindRole, error)
	GetOrgKindRoleByRoleID(roleID string) ([]model.SsoOrgKindRole, error)
	DeleteOrgKindRole(id string) error
	DeleteOrgKindRoleByKindCode(kindCode string) error
	GetAllOrgKindRoles() ([]model.SsoOrgKindRole, error)
}

type SsoOrgMenuRepository interface {
	CreateOrgMenu(item *model.SsoOrgMenu) error
	GetOrgMenuByID(id string) (*model.SsoOrgMenu, error)
	GetOrgMenusByOrgID(orgID string) ([]model.SsoOrgMenu, error)
	GetOrgMenusByMenuID(menuID string) ([]model.SsoOrgMenu, error)
	DeleteOrgMenu(id string) error
	DeleteOrgMenuByOrgID(orgID string) error
	GetAllOrgMenus() ([]model.SsoOrgMenu, error)
}

type SsoOrgRoleRepository interface {
	CreateOrgRole(item *model.SsoOrgRole) error
	GetOrgRoleByID(id string) (*model.SsoOrgRole, error)
	GetOrgRolesByOrgID(orgID string) ([]model.SsoOrgRole, error)
	GetOrgRolesByRoleID(roleID string) ([]model.SsoOrgRole, error)
	DeleteOrgRole(id string) error
	DeleteOrgRoleByOrgID(orgID string) error
	DeleteOrgRolesByRoleID(roleID string) error
	GetAllOrgRoles() ([]model.SsoOrgRole, error)
}

type SsoOrgUserRepository interface {
	CreateOrgUser(item *model.SsoOrgUser) error
	GetOrgUserByID(id string) (*model.SsoOrgUser, error)
	GetOrgUsersByOrgID(orgID string) ([]model.SsoOrgUser, error)
	GetOrgUsersByUserID(userID string) ([]model.SsoOrgUser, error)
	DeleteOrgUser(id string) error
	DeleteOrgUserByOrgID(orgID string) error
	GetAllOrgUsers() ([]model.SsoOrgUser, error)
}

// Repositories 用户模块仓库集合
type Repositories struct {
	User          SsoUserRepository
	Tenant        SsoTenantRepository
	App           SsoAppRepository
	Menu          SsoMenuRepository
	Role          SsoRoleRepository
	Org           SsoOrgRepository
	OrgKind       SsoOrgKindRepository
	OrgKindRole   SsoOrgKindRoleRepository
	OrgMenu       SsoOrgMenuRepository
	OrgRole       SsoOrgRoleRepository
	OrgUser       SsoOrgUserRepository
	RoleGroup     SsoRoleGroupRepository
	UserGroup     SsoUserGroupRepository
	RoleGroupRole SsoRoleGroupRoleRepository
	UserGroupUser SsoUserGroupUserRepository
	UserRoleGroup SsoUserRoleGroupRepository
	Pos           SsoPosRepository
	UserRole      SsoUserRoleRepository
	UserPos       SsoUserPosRepository
	RoleMenu      SsoRoleMenuRepository
	PosRole       SsoPosRoleRepository
	CasbinRule    SsoCasbinRuleRepository
}

// NewRepositories 创建用户模块仓库集合
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:          NewSsoUserRepository(db),
		Tenant:        NewSsoTenantRepository(db),
		App:           NewSsoAppRepository(db),
		Menu:          NewSsoMenuRepository(db),
		Role:          NewSsoRoleRepository(db),
		Org:           NewSsoOrgRepository(db),
		OrgKind:       NewSsoOrgKindRepository(db),
		OrgKindRole:   NewSsoOrgKindRoleRepository(db),
		OrgMenu:       NewSsoOrgMenuRepository(db),
		OrgRole:       NewSsoOrgRoleRepository(db),
		OrgUser:       NewSsoOrgUserRepository(db),
		RoleGroup:     NewSsoRoleGroupRepository(db),
		UserGroup:     NewSsoUserGroupRepository(db),
		RoleGroupRole: NewSsoRoleGroupRoleRepository(db),
		UserGroupUser: NewSsoUserGroupUserRepository(db),
		UserRoleGroup: NewSsoUserRoleGroupRepository(db),
		Pos:           NewSsoPosRepository(db),
		UserRole:      NewSsoUserRoleRepository(db),
		UserPos:       NewSsoUserPosRepository(db),
		RoleMenu:      NewSsoRoleMenuRepository(db),
		PosRole:       NewSsoPosRoleRepository(db),
		CasbinRule:    NewSsoCasbinRuleRepository(db),
	}
}
