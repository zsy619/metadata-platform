/**
 * SSO单点登录模块API服务
 */
import type { SsoClient, SsoFieldMapping, SsoKey, SsoProtocolConfig, SsoSession } from '@/types/sso'
import request from '@/utils/request'

// ==================== 协议配置API ====================

/**
 * 获取协议配置列表
 */
export const getProtocolConfigs = async (): Promise<SsoProtocolConfig[]> => {
  return request({
    url: '/api/sso/config',
    method: 'get'
  })
}

/**
 * 根据ID获取协议配置
 */
export const getProtocolConfigById = async (id: string): Promise<SsoProtocolConfig> => {
  return request({
    url: `/api/sso/config/${id}`,
    method: 'get'
  })
}

/**
 * 根据协议类型获取配置列表
 */
export const getProtocolConfigsByType = async (type: string): Promise<SsoProtocolConfig[]> => {
  return request({
    url: `/api/sso/config/type/${type}`,
    method: 'get'
  })
}

/**
 * 创建协议配置
 */
export const createProtocolConfig = async (data: Partial<SsoProtocolConfig>): Promise<SsoProtocolConfig> => {
  return request({
    url: '/api/sso/config',
    method: 'post',
    data
  })
}

/**
 * 更新协议配置
 */
export const updateProtocolConfig = async (id: string, data: Partial<SsoProtocolConfig>): Promise<SsoProtocolConfig> => {
  return request({
    url: `/api/sso/config/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除协议配置
 */
export const deleteProtocolConfig = async (id: string): Promise<void> => {
  return request({
    url: `/api/sso/config/${id}`,
    method: 'delete'
  })
}

// ==================== 客户端配置API ====================

/**
 * 获取客户端列表
 */
export const getClients = async (): Promise<SsoClient[]> => {
  return request({
    url: '/api/sso/client',
    method: 'get'
  })
}

/**
 * 根据ID获取客户端
 */
export const getClientById = async (id: string): Promise<SsoClient> => {
  return request({
    url: `/api/sso/client/${id}`,
    method: 'get'
  })
}

/**
 * 创建客户端
 */
export const createClient = async (data: Partial<SsoClient>): Promise<SsoClient> => {
  return request({
    url: '/api/sso/client',
    method: 'post',
    data
  })
}

/**
 * 更新客户端
 */
export const updateClient = async (id: string, data: Partial<SsoClient>): Promise<SsoClient> => {
  return request({
    url: `/api/sso/client/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除客户端
 */
export const deleteClient = async (id: string): Promise<void> => {
  return request({
    url: `/api/sso/client/${id}`,
    method: 'delete'
  })
}

// ==================== 密钥管理API ====================

/**
 * 获取密钥列表
 */
export const getKeys = async (): Promise<SsoKey[]> => {
  return request({
    url: '/api/sso/key',
    method: 'get'
  })
}

/**
 * 根据ID获取密钥
 */
export const getKeyById = async (id: string): Promise<SsoKey> => {
  return request({
    url: `/api/sso/key/${id}`,
    method: 'get'
  })
}

/**
 * 创建密钥
 */
export const createKey = async (data: Partial<SsoKey>): Promise<SsoKey> => {
  return request({
    url: '/api/sso/key',
    method: 'post',
    data
  })
}

/**
 * 生成密钥对
 */
export const generateKeyPair = async (keyType: string, algorithm?: string): Promise<SsoKey> => {
  return request({
    url: '/api/sso/key/generate',
    method: 'post',
    data: { key_type: keyType, algorithm }
  })
}

/**
 * 更新密钥
 */
export const updateKey = async (id: string, data: Partial<SsoKey>): Promise<SsoKey> => {
  return request({
    url: `/api/sso/key/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除密钥
 */
export const deleteKey = async (id: string): Promise<void> => {
  return request({
    url: `/api/sso/key/${id}`,
    method: 'delete'
  })
}

// ==================== 字段映射API ====================

/**
 * 获取字段映射列表
 */
export const getFieldMappings = async (): Promise<SsoFieldMapping[]> => {
  return request({
    url: '/api/sso/mapping',
    method: 'get'
  })
}

/**
 * 根据ID获取字段映射
 */
export const getFieldMappingById = async (id: string): Promise<SsoFieldMapping> => {
  return request({
    url: `/api/sso/mapping/${id}`,
    method: 'get'
  })
}

/**
 * 创建字段映射
 */
export const createFieldMapping = async (data: Partial<SsoFieldMapping>): Promise<SsoFieldMapping> => {
  return request({
    url: '/api/sso/mapping',
    method: 'post',
    data
  })
}

/**
 * 更新字段映射
 */
export const updateFieldMapping = async (id: string, data: Partial<SsoFieldMapping>): Promise<SsoFieldMapping> => {
  return request({
    url: `/api/sso/mapping/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除字段映射
 */
export const deleteFieldMapping = async (id: string): Promise<void> => {
  return request({
    url: `/api/sso/mapping/${id}`,
    method: 'delete'
  })
}

// ==================== 会话管理API ====================

/**
 * 获取会话列表
 */
export const getSessions = async (userId?: string): Promise<SsoSession[]> => {
  const params = userId ? { user_id: userId } : {}
  return request({
    url: '/api/sso/session',
    method: 'get',
    params
  })
}

/**
 * 根据ID获取会话
 */
export const getSessionById = async (id: string): Promise<SsoSession> => {
  return request({
    url: `/api/sso/session/${id}`,
    method: 'get'
  })
}

/**
 * 删除会话
 */
export const deleteSession = async (id: string): Promise<void> => {
  return request({
    url: `/api/sso/session/${id}`,
    method: 'delete'
  })
}

/**
 * 撤销会话
 */
export const revokeSession = async (id: string): Promise<void> => {
  return request({
    url: `/api/sso/session/${id}/revoke`,
    method: 'post'
  })
}
