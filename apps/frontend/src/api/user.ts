/**
 * 用户管理模块API服务
 */
import type { LoginRequest, LoginResponse, Tenant, User } from '@/types/user'
import request from '@/utils/request'

// ==================== 用户相关API ====================

/**
 * 获取用户列表
 * @returns 用户列表
 */
export const getUsers = async (): Promise<User[]> => {
  return request({
    url: '/api/user',
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
    url: `/api/user/${id}`,
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
    url: '/api/user',
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
    url: `/api/user/${id}`,
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
    url: `/api/user/${id}`,
    method: 'delete'
  })
}

/**
 * 用户登录
 * @param data 登录信息
 * @returns 登录结果
 */
export const login = async (data: LoginRequest): Promise<LoginResponse> => {
  return request({
    url: '/api/user/login',
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
    url: '/api/tenant',
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
    url: `/api/tenant/${id}`,
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
    url: '/api/tenant',
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
    url: `/api/tenant/${id}`,
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
    url: `/api/tenant/${id}`,
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
    url: '/api/app',
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
    url: `/api/app/${id}`,
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
    url: '/api/app',
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
    url: `/api/app/${id}`,
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
    url: `/api/app/${id}`,
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
    url: '/api/menu',
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
    url: `/api/menu/${id}`,
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
    url: '/api/menu',
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
    url: `/api/menu/${id}`,
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
    url: `/api/menu/${id}`,
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
    url: '/api/role',
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
    url: `/api/role/${id}`,
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
    url: '/api/role',
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
    url: `/api/role/${id}`,
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
    url: `/api/role/${id}`,
    method: 'delete'
  })
}

// ==================== 组织相关API ====================

/**
 * 获取组织列表
 * @returns 组织列表
 */
export const getUnits = async (): Promise<Unit[]> => {
  return request({
    url: '/api/unit',
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
    url: `/api/unit/${id}`,
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
    url: '/api/unit',
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
    url: `/api/unit/${id}`,
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
    url: `/api/unit/${id}`,
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
    url: '/api/pos',
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
    url: `/api/pos/${id}`,
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
    url: '/api/pos',
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
    url: `/api/pos/${id}`,
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
    url: `/api/pos/${id}`,
    method: 'delete'
  })
}
