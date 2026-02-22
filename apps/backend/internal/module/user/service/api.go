package service

import (
	"context"

	"gorm.io/gorm"

	"metadata-platform/internal/module/audit/queue"
	auditService "metadata-platform/internal/module/audit/service"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/repository"
	"metadata-platform/internal/utils"
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
	// 用户关联相关
	GetUserRoles(userID string) ([]string, error)
	UpdateUserRoles(userID string, roleIDs []string, createBy string) error
	GetUserPos(userID string) ([]string, error)
	UpdateUserPos(userID string, posIDs []string, createBy string) error
	GetUserGroups(userID string) ([]string, error)
	UpdateUserGroups(userID string, groupIDs []string, createBy string) error
	GetUserRoleGroups(userID string) ([]string, error)
	UpdateUserRoleGroups(userID string, roleGroupIDs []string, createBy string) error
	GetUserOrgs(userID string) ([]string, error)
	UpdateUserOrgs(userID string, orgIDs []string, createBy string) error
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
	UpdateAppFields(id string, fields map[string]any) error
	DeleteApp(id string) error
	GetAllApps() ([]model.SsoApp, error)
}

// SsoMenuService 菜单服务接口
type SsoMenuService interface {
	CreateMenu(menu *model.SsoMenu) error
	GetMenuByID(id string) (*model.SsoMenu, error)
	GetMenuByCode(code string) (*model.SsoMenu, error)
	UpdateMenu(menu *model.SsoMenu) error
	UpdateMenuFields(id string, fields map[string]any) error
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
	HasChildren(parentID string) (bool, error)
	// 角色菜单相关
	GetRoleMenus(roleID string) ([]string, error)
	UpdateRoleMenus(roleID string, menuIDs []string, createBy string) error
}

// SsoOrgService 组织服务接口
type SsoOrgService interface {
	CreateOrg(org *model.SsoOrg) error
	GetOrgByID(id string) (*model.SsoOrg, error)
	GetOrgByCode(code string) (*model.SsoOrg, error)
	UpdateOrg(org *model.SsoOrg) error
	UpdateOrgFields(id string, fields map[string]any) error
	DeleteOrg(id string) error
	GetAllOrgs() ([]model.SsoOrg, error)
}

// SsoOrgKindService 组织类型服务接口
type SsoOrgKindService interface {
	CreateOrgKind(item *model.SsoOrgKind) error
	GetOrgKindByID(id string) (*model.SsoOrgKind, error)
	GetOrgKindByCode(code string) (*model.SsoOrgKind, error)
	UpdateOrgKind(item *model.SsoOrgKind) error
	UpdateOrgKindFields(id string, fields map[string]any) error
	DeleteOrgKind(id string) error
	GetAllOrgKinds() ([]model.SsoOrgKind, error)
}

// SsoPosService 职位服务接口
type SsoPosService interface {
	CreatePos(pos *model.SsoPos) error
	GetPosByID(id string) (*model.SsoPos, error)
	GetPosByCode(code string) (*model.SsoPos, error)
	UpdatePos(pos *model.SsoPos) error
	UpdatePosFields(id string, fields map[string]any) error
	DeletePos(id string) error
	GetAllPoss() ([]model.SsoPos, error)
	// 职位角色关联相关
	GetPosRoles(posID string) ([]string, error)
	UpdatePosRoles(posID string, roleIDs []string, createBy string) error
}

// SsoRoleGroupService 角色分组服务接口
type SsoRoleGroupService interface {
	CreateRoleGroup(item *model.SsoRoleGroup) error
	GetRoleGroupByID(id string) (*model.SsoRoleGroup, error)
	GetRoleGroupByCode(code string) (*model.SsoRoleGroup, error)
	UpdateRoleGroup(item *model.SsoRoleGroup) error
	UpdateRoleGroupFields(id string, fields map[string]any) error
	DeleteRoleGroup(id string) error
	GetAllRoleGroups() ([]model.SsoRoleGroup, error)
	// 角色组角色关联相关
	GetRoleGroupRoles(groupID string) ([]string, error)
	UpdateRoleGroupRoles(groupID string, roleIDs []string, createBy string) error
}

// SsoUserGroupService 用户组服务接口
type SsoUserGroupService interface {
	CreateUserGroup(item *model.SsoUserGroup) error
	GetUserGroupByID(id string) (*model.SsoUserGroup, error)
	GetUserGroupByCode(code string) (*model.SsoUserGroup, error)
	UpdateUserGroup(item *model.SsoUserGroup) error
	UpdateUserGroupFields(id string, fields map[string]any) error
	DeleteUserGroup(id string) error
	GetAllUserGroups() ([]model.SsoUserGroup, error)
	// 用户组角色关联相关
	GetUserGroupRoles(groupID string) ([]string, error)
	UpdateUserGroupRoles(groupID string, roleIDs []string, createBy string) error
}

// SsoAuthService 认证服务接口
type SsoAuthService interface {
	Login(account string, password string, tenantID uint, clientInfo utils.ClientInfo) (accessToken string, refreshToken string, user *model.SsoUser, err error)
	Logout(ctx context.Context, userID string, clientInfo utils.ClientInfo) error
	Refresh(refreshToken string) (newAccessToken string, err error)
	GetUserInfo(userID string) (*model.SsoUser, error)
	ChangePassword(userID string, oldPassword string, newPassword string) error
}

// SsoUserProfileService 用户档案服务接口
type SsoUserProfileService interface {
	GetByUserID(userID string) (*model.SsoUserProfile, error)
	Upsert(profile *model.SsoUserProfile) error
}

// SsoUserAddressService 用户地址服务接口
type SsoUserAddressService interface {
	GetByUserID(userID string) ([]model.SsoUserAddress, error)
	Create(addr *model.SsoUserAddress) error
	UpdateFields(id string, fields map[string]any) error
	SetDefault(userID, id string) error
	Delete(userID, id string) error
}

// SsoUserContactService 用户联系方式服务接口
type SsoUserContactService interface {
	GetByUserID(userID string) ([]model.SsoUserContact, error)
	Create(contact *model.SsoUserContact) error
	UpdateFields(id string, fields map[string]any) error
	Delete(userID, id string) error
}

// SsoUserSocialService 用户第三方账号服务接口
type SsoUserSocialService interface {
	GetByUserID(userID string) ([]model.SsoUserSocial, error)
	Bind(social *model.SsoUserSocial) error
	Unbind(userID, id string) error
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
	OrgKind    SsoOrgKindService
	RoleGroup  SsoRoleGroupService
	UserGroup  SsoUserGroupService
	CasbinSync SsoCasbinSyncService
	Audit      auditService.AuditService
	UserProfile SsoUserProfileService
	UserAddress SsoUserAddressService
	UserContact SsoUserContactService
	UserSocial  SsoUserSocialService
}

// NewServices 创建用户模块服务集合
func NewServices(repos *repository.Repositories, db *gorm.DB, auditDB *gorm.DB, auditQueue *queue.AuditLogQueue) *Services {
	auditSvc := auditService.NewAuditService(auditDB, auditQueue)
	return &Services{
		User:      NewSsoUserService(repos.User, repos.Org, repos.OrgUser, repos.UserRole, repos.UserPos, repos.UserGroupUser, repos.UserRoleGroup, repos.UserProfile, repos.UserAddress, repos.UserContact, repos.UserSocial),
		Tenant:    NewSsoTenantService(repos.Tenant),
		App:       NewSsoAppService(repos.App),
		Menu:      NewSsoMenuService(repos.Menu),
		Role:      NewSsoRoleService(repos.Role, repos.RoleMenu, repos.UserRole, repos.PosRole, repos.OrgRole, repos.RoleGroupRole, repos.UserGroupRole, repos.OrgKindRole),
		Org:       NewSsoOrgService(repos.Org, repos.OrgUser, repos.OrgRole, repos.OrgMenu),
		Pos:       NewSsoPosService(repos.Pos, repos.PosRole, repos.UserPos),
		Auth:      NewSsoAuthService(repos.User, repos.Role, repos.UserRole, auditSvc),
		OrgKind:   NewSsoOrgKindService(repos.OrgKind),
		RoleGroup: NewSsoRoleGroupService(repos.RoleGroup, repos.RoleGroupRole),
		UserGroup: NewSsoUserGroupService(repos.UserGroup, repos.UserGroupUser, repos.UserGroupRole),
		CasbinSync: NewSsoCasbinSyncService(
			repos.UserRole,
			repos.RoleMenu,
			repos.Role,
			repos.Menu,
		),
		Audit:       auditSvc,
		UserProfile: NewSsoUserProfileService(repos.UserProfile),
		UserAddress: NewSsoUserAddressService(repos.UserAddress),
		UserContact: NewSsoUserContactService(repos.UserContact),
		UserSocial:  NewSsoUserSocialService(repos.UserSocial),
	}
}
