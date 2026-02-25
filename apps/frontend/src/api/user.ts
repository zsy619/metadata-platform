/**
 * 用户管理模块API服务
 */
import type { App, LoginRequest, LoginResponse, Menu, Pos, Role, Tenant, Unit, User } from '@/types/user'
import request from '@/utils/request'

// ==================== 用户相关API ====================

/**
 * 获取用户列表
 * @returns 用户列表
 */
export const getUsers = async (): Promise<User[]> => {
  return request({
      url: '/api/sso/user',
    method: 'get'
  })
}

/**
 * 根据ID获取用户
 * @param id 用户ID
 * @returns 用户详情
 */
export const getUserById = async (id: string): Promise<User> => {
  return request({
    url: `/api/sso/user/${id}`,
    method: 'get'
  })
}

/**
 * 创建用户
 * @param data 用户数据
 * @returns 创建结果
 */
export const createUser = async (data: Partial<User>): Promise<User> => {
  return request({
      url: '/api/sso/user',
    method: 'post',
    data
  })
}

/**
 * 更新用户
 * @param id 用户ID
 * @param data 用户数据
 * @returns 更新结果
 */
export const updateUser = async (id: string, data: Partial<User>): Promise<User> => {
  return request({
    url: `/api/sso/user/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除用户
 * @param id 用户ID
 * @returns 删除结果
 */
export const deleteUser = async (id: string): Promise<void> => {
  return request({
    url: `/api/sso/user/${id}`,
    method: 'delete'
  })
}

// ==================== 用户关联API ====================

/**
 * 获取用户的角色ID列表
 * @param userId 用户ID
 * @returns 角色ID列表
 */
export const getUserRoles = async (userId: string): Promise<{ user_id: string; role_ids: string[] }> => {
  return request({
    url: `/api/sso/user/${userId}/roles`,
    method: 'get'
  })
}

/**
 * 更新用户的角色关联
 * @param userId 用户ID
 * @param roleIds 角色ID列表
 * @returns 更新结果
 */
export const updateUserRoles = async (userId: string, roleIds: string[]): Promise<void> => {
  return request({
    url: `/api/sso/user/${userId}/roles`,
    method: 'put',
    data: { role_ids: roleIds }
  })
}

/**
 * 获取用户的职位ID列表
 * @param userId 用户ID
 * @returns 职位ID列表
 */
export const getUserPos = async (userId: string): Promise<{ user_id: string; pos_ids: string[] }> => {
  return request({
    url: `/api/sso/user/${userId}/pos`,
    method: 'get'
  })
}

/**
 * 更新用户的职位关联
 * @param userId 用户ID
 * @param posIds 职位ID列表
 * @returns 更新结果
 */
export const updateUserPos = async (userId: string, posIds: string[]): Promise<void> => {
  return request({
    url: `/api/sso/user/${userId}/pos`,
    method: 'put',
    data: { pos_ids: posIds }
  })
}

/**
 * 获取用户的用户组ID列表
 * @param userId 用户ID
 * @returns 用户组ID列表
 */
export const getUserGroups = async (userId: string): Promise<{ user_id: string; group_ids: string[] }> => {
  return request({
    url: `/api/sso/user/${userId}/groups`,
    method: 'get'
  })
}

/**
 * 更新用户的用户组关联
 * @param userId 用户ID
 * @param groupIds 用户组ID列表
 * @returns 更新结果
 */
export const updateUserGroups = async (userId: string, groupIds: string[]): Promise<void> => {
  return request({
    url: `/api/sso/user/${userId}/groups`,
    method: 'put',
    data: { group_ids: groupIds }
  })
}

/**
 * 获取用户的角色组ID列表
 * @param userId 用户ID
 * @returns 角色组ID列表
 */
export const getUserRoleGroups = async (userId: string): Promise<{ user_id: string; role_group_ids: string[] }> => {
  return request({
    url: `/api/sso/user/${userId}/role-groups`,
    method: 'get'
  })
}

/**
 * 更新用户的角色组关联
 * @param userId 用户ID
 * @param roleGroupIds 角色组ID列表
 * @returns 更新结果
 */
export const updateUserRoleGroups = async (userId: string, roleGroupIds: string[]): Promise<void> => {
  return request({
    url: `/api/sso/user/${userId}/role-groups`,
    method: 'put',
    data: { role_group_ids: roleGroupIds }
  })
}

/**
 * 获取用户的组织ID列表
 * @param userId 用户ID
 * @returns 组织ID列表
 */
export const getUserOrgs = async (userId: string): Promise<{ user_id: string; org_ids: string[] }> => {
  return request({
    url: `/api/sso/user/${userId}/orgs`,
    method: 'get'
  })
}

/**
 * 更新用户的组织关联
 * @param userId 用户ID
 * @param orgIds 组织ID列表
 * @returns 更新结果
 */
export const updateUserOrgs = async (userId: string, orgIds: string[]): Promise<void> => {
  return request({
    url: `/api/sso/user/${userId}/orgs`,
    method: 'put',
    data: { org_ids: orgIds }
  })
}

/**
 * 用户登录
 * @param data 登录信息
 * @returns 登录结果
 */
export const login = async (data: LoginRequest): Promise<LoginResponse> => {
  return request({
    url: '/api/sso/user/login',
    method: 'post',
    data
  })
}

// ==================== 租户相关API ====================

/**
 * 获取租户列表
 * @returns 租户列表
 */
export const getTenants = async (): Promise<Tenant[]> => {
  return request({
      url: '/api/sso/tenant',
    method: 'get'
  })
}

/**
 * 根据ID获取租户
 * @param id 租户ID
 * @returns 租户详情
 */
export const getTenantById = async (id: string): Promise<Tenant> => {
  return request({
    url: `/api/sso/tenant/${id}`,
    method: 'get'
  })
}

/**
 * 创建租户
 * @param data 租户数据
 * @returns 创建结果
 */
export const createTenant = async (data: Partial<Tenant>): Promise<Tenant> => {
  return request({
      url: '/api/sso/tenant',
    method: 'post',
    data
  })
}

/**
 * 更新租户
 * @param id 租户ID
 * @param data 租户数据
 * @returns 更新结果
 */
export const updateTenant = async (id: string, data: Partial<Tenant>): Promise<Tenant> => {
  return request({
    url: `/api/sso/tenant/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除租户
 * @param id 租户ID
 * @returns 删除结果
 */
export const deleteTenant = async (id: string): Promise<void> => {
  return request({
    url: `/api/sso/tenant/${id}`,
    method: 'delete'
  })
}

// ==================== 应用相关API ====================

/**
 * 获取应用列表
 * @returns 应用列表
 */
export const getApps = async (): Promise<App[]> => {
  return request({
      url: '/api/sso/app',
    method: 'get'
  })
}

/**
 * 根据ID获取应用
 * @param id 应用ID
 * @returns 应用详情
 */
export const getAppById = async (id: string): Promise<App> => {
  return request({
    url: `/api/sso/app/${id}`,
    method: 'get'
  })
}

/**
 * 创建应用
 * @param data 应用数据
 * @returns 创建结果
 */
export const createApp = async (data: Partial<App>): Promise<App> => {
  return request({
      url: '/api/sso/app',
    method: 'post',
    data
  })
}

/**
 * 更新应用
 * @param id 应用ID
 * @param data 应用数据
 * @returns 更新结果
 */
export const updateApp = async (id: string, data: Partial<App>): Promise<App> => {
  return request({
    url: `/api/sso/app/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除应用
 * @param id 应用ID
 * @returns 删除结果
 */
export const deleteApp = async (id: string): Promise<void> => {
  return request({
    url: `/api/sso/app/${id}`,
    method: 'delete'
  })
}

// ==================== 菜单相关API ====================

/**
 * 获取菜单列表
 * @returns 菜单列表
 */
export const getMenus = async (): Promise<Menu[]> => {
  return request({
      url: '/api/sso/menu',
    method: 'get'
  })
}

/**
 * 根据ID获取菜单
 * @param id 菜单ID
 * @returns 菜单详情
 */
export const getMenuById = async (id: string): Promise<Menu> => {
  return request({
    url: `/api/sso/menu/${id}`,
    method: 'get'
  })
}

/**
 * 创建菜单
 * @param data 菜单数据
 * @returns 创建结果
 */
export const createMenu = async (data: Partial<Menu>): Promise<Menu> => {
  return request({
      url: '/api/sso/menu',
    method: 'post',
    data
  })
}

/**
 * 更新菜单
 * @param id 菜单ID
 * @param data 菜单数据
 * @returns 更新结果
 */
export const updateMenu = async (id: string, data: Partial<Menu>): Promise<Menu> => {
  return request({
    url: `/api/sso/menu/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除菜单
 * @param id 菜单ID
 * @returns 删除结果
 */
export const deleteMenu = async (id: string): Promise<void> => {
  return request({
    url: `/api/sso/menu/${id}`,
    method: 'delete'
  })
}

// ==================== 角色相关API ====================

/**
 * 获取角色列表
 * @returns 角色列表
 */
export const getRoles = async (): Promise<Role[]> => {
  return request({
      url: '/api/sso/role',
    method: 'get'
  })
}

/**
 * 根据ID获取角色
 * @param id 角色ID
 * @returns 角色详情
 */
export const getRoleById = async (id: string): Promise<Role> => {
  return request({
    url: `/api/sso/role/${id}`,
    method: 'get'
  })
}

/**
 * 创建角色
 * @param data 角色数据
 * @returns 创建结果
 */
export const createRole = async (data: Partial<Role>): Promise<Role> => {
  return request({
      url: '/api/sso/role',
    method: 'post',
    data
  })
}

/**
 * 更新角色
 * @param id 角色ID
 * @param data 角色数据
 * @returns 更新结果
 */
export const updateRole = async (id: string, data: Partial<Role>): Promise<Role> => {
  return request({
    url: `/api/sso/role/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除角色
 * @param id 角色ID
 * @returns 删除结果
 */
export const deleteRole = async (id: string): Promise<void> => {
  return request({
    url: `/api/sso/role/${id}`,
    method: 'delete'
  })
}

/**
 * 获取角色的菜单ID列表
 * @param roleId 角色ID
 * @returns 菜单ID列表
 */
export const getRoleMenus = async (roleId: string): Promise<{ role_id: string; menu_ids: string[] }> => {
  return request({
    url: `/api/sso/role/${roleId}/menus`,
    method: 'get'
  })
}

/**
 * 更新角色的菜单关联
 * @param roleId 角色ID
 * @param menuIds 菜单ID列表
 * @returns 更新结果
 */
export const updateRoleMenus = async (roleId: string, menuIds: string[]): Promise<void> => {
  return request({
    url: `/api/sso/role/${roleId}/menus`,
    method: 'put',
    data: { menu_ids: menuIds }
  })
}

// ==================== 组织相关API ====================

/**
 * 获取组织列表
 * @returns 组织列表
 */
export const getUnits = async (): Promise<Unit[]> => {
  return request({
      url: '/api/sso/org',
    method: 'get'
  })
}

/**
 * 根据ID获取组织
 * @param id 组织ID
 * @returns 组织详情
 */
export const getUnitById = async (id: string): Promise<Unit> => {
  return request({
    url: `/api/sso/org/${id}`,
    method: 'get'
  })
}

/**
 * 创建组织
 * @param data 组织数据
 * @returns 创建结果
 */
export const createUnit = async (data: Partial<Unit>): Promise<Unit> => {
  return request({
      url: '/api/sso/org',
    method: 'post',
    data
  })
}

/**
 * 更新组织
 * @param id 组织ID
 * @param data 组织数据
 * @returns 更新结果
 */
export const updateUnit = async (id: string, data: Partial<Unit>): Promise<Unit> => {
  return request({
    url: `/api/sso/org/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除组织
 * @param id 组织ID
 * @returns 删除结果
 */
export const deleteUnit = async (id: string): Promise<void> => {
  return request({
    url: `/api/sso/org/${id}`,
    method: 'delete'
  })
}

// ==================== 职位相关API ====================

/**
 * 获取职位列表
 * @returns 职位列表
 */
export const getPos = async (): Promise<Pos[]> => {
  return request({
      url: '/api/sso/pos',
    method: 'get'
  })
}

/**
 * 根据ID获取职位
 * @param id 职位ID
 * @returns 职位详情
 */
export const getPosById = async (id: string): Promise<Pos> => {
  return request({
    url: `/api/sso/pos/${id}`,
    method: 'get'
  })
}

/**
 * 创建职位
 * @param data 职位数据
 * @returns 创建结果
 */
export const createPos = async (data: Partial<Pos>): Promise<Pos> => {
  return request({
      url: '/api/sso/pos',
    method: 'post',
    data
  })
}

/**
 * 更新职位
 * @param id 职位ID
 * @param data 职位数据
 * @returns 更新结果
 */
export const updatePos = async (id: string, data: Partial<Pos>): Promise<Pos> => {
  return request({
    url: `/api/sso/pos/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除职位
 * @param id 职位ID
 * @returns 删除结果
 */
export const deletePos = async (id: string): Promise<void> => {
  return request({
    url: `/api/sso/pos/${id}`,
    method: 'delete'
  })
}

/**
 * 获取职位的角色ID列表
 * @param posId 职位ID
 * @returns 角色ID列表
 */
export const getPosRoles = async (posId: string): Promise<string[]> => {
  const res: any = await request({
    url: `/api/sso/pos/${posId}/roles`,
    method: 'get'
  })
  return res.role_ids || []
}

/**
 * 更新职位的角色关联
 * @param posId 职位ID
 * @param roleIds 角色ID列表
 * @returns 更新结果
 */
export const updatePosRoles = async (posId: string, roleIds: string[]): Promise<void> => {
  return request({
    url: `/api/sso/pos/${posId}/roles`,
    method: 'put',
    data: { role_ids: roleIds }
  })
}

/**
 * 获取角色组的角色ID列表
 * @param groupId 角色组ID
 * @returns 角色ID列表
 */
export const getRoleGroupRoles = async (groupId: string): Promise<string[]> => {
  const res: any = await request({
    url: `/api/sso/role-group/${groupId}/roles`,
    method: 'get'
  })
  return res.role_ids || []
}

/**
 * 更新角色组的角色关联
 * @param groupId 角色组ID
 * @param roleIds 角色ID列表
 * @returns 更新结果
 */
export const updateRoleGroupRoles = async (groupId: string, roleIds: string[]): Promise<void> => {
  return request({
    url: `/api/sso/role-group/${groupId}/roles`,
    method: 'put',
    data: { role_ids: roleIds }
  })
}

/**
 * 获取用户组的角色ID列表
 * @param groupId 用户组ID
 * @returns 角色ID列表
 */
export const getUserGroupRoles = async (groupId: string): Promise<string[]> => {
  const res: any = await request({
    url: `/api/sso/user-group/${groupId}/roles`,
    method: 'get'
  })
  return res.role_ids || []
}

/**
 * 更新用户组的角色关联
 * @param groupId 用户组ID
 * @param roleIds 角色ID列表
 * @returns 更新结果
 */
export const updateUserGroupRoles = async (groupId: string, roleIds: string[]): Promise<void> => {
  return request({
    url: `/api/sso/user-group/${groupId}/roles`,
    method: 'put',
    data: { role_ids: roleIds }
  })
}

// ==================== 组织类型相关API ====================

export const getOrgKinds = async (): Promise<any[]> => {
  return request({ url: '/api/sso/org-kind', method: 'get' })
}

export const getOrgKindById = async (id: string): Promise<any> => {
  return request({ url: `/api/sso/org-kind/${id}`, method: 'get' })
}

export const createOrgKind = async (data: any): Promise<any> => {
    return request({ url: '/api/sso/org-kind', method: 'post', data })
}

export const updateOrgKind = async (id: string, data: any): Promise<any> => {
  return request({ url: `/api/sso/org-kind/${id}`, method: 'put', data })
}

export const deleteOrgKind = async (id: string): Promise<void> => {
  return request({ url: `/api/sso/org-kind/${id}`, method: 'delete' })
}

// ==================== 角色组相关API ====================

export const getRoleGroups = async (): Promise<any[]> => {
    return request({ url: '/api/sso/role-group', method: 'get' })
}

export const createRoleGroup = async (data: any): Promise<any> => {
    return request({ url: '/api/sso/role-group', method: 'post', data })
}

export const updateRoleGroup = async (id: string, data: any): Promise<any> => {
    return request({ url: `/api/sso/role-group/${id}`, method: 'put', data })
}

export const deleteRoleGroup = async (id: string): Promise<void> => {
    return request({ url: `/api/sso/role-group/${id}`, method: 'delete' })
}

// ==================== 用户组相关API ====================

export const getAllUserGroups = async (): Promise<any[]> => {
    return request({ url: '/api/sso/user-group', method: 'get' })
}

export const createUserGroup = async (data: any): Promise<any> => {
    return request({ url: '/api/sso/user-group', method: 'post', data })
}

export const updateUserGroup = async (id: string, data: any): Promise<any> => {
    return request({ url: `/api/sso/user-group/${id}`, method: 'put', data })
}

export const deleteUserGroup = async (id: string): Promise<void> => {
    return request({ url: `/api/sso/user-group/${id}`, method: 'delete' })
}

// ==================== 仪表板相关API ====================

export const getDashboardStats = async (): Promise<{
    user_count: number
    role_count: number
    org_count: number
    menu_count: number
    pos_count: number
    user_group_count: number
    role_group_count: number
}> => {
    return request({ url: '/api/sso/dashboard/stats', method: 'get' })
}

export const getRecentLoginLogs = async (): Promise<any[]> => {
    return request({ url: '/api/sso/dashboard/login-logs', method: 'get' })
}

export const getRecentOperationLogs = async (): Promise<any[]> => {
    return request({ url: '/api/sso/dashboard/operation-logs', method: 'get' })
}

// 登录趋势项
export interface LoginTrendItem {
    date: string
    success: number
    fail: number
}

// 用户状态分布
export interface UserStatusDistribution {
    active: number
    inactive: number
    locked: number
    pending: number
}

// 操作统计
export interface OperationStats {
    create: number
    update: number
    delete: number
    query: number
    export: number
}

// 组织分布项
export interface OrgDistribution {
    name: string
    value: number
}

// 获取登录趋势（最近7天）
export const getLoginTrend = async (): Promise<LoginTrendItem[]> => {
    return request({ url: '/api/sso/dashboard/login-trend', method: 'get' })
}

// 获取用户状态分布
export const getUserStatusDistribution = async (): Promise<UserStatusDistribution> => {
    return request({ url: '/api/sso/dashboard/user-status', method: 'get' })
}

// 获取操作统计
export const getOperationStats = async (): Promise<OperationStats> => {
    return request({ url: '/api/sso/dashboard/operation-stats', method: 'get' })
}

// 获取组织分布
export const getOrgDistribution = async (): Promise<OrgDistribution[]> => {
    return request({ url: '/api/sso/dashboard/org-distribution', method: 'get' })
}
