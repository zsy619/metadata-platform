/**
 * SSO 单点登录系统类型定义
 */

/**
 * 协议配置
 */
export interface SsoProtocolConfig {
  id: string
  name: string
  protocolType: string
  enabled: boolean
  config: Record<string, any>
  createdAt: string
  updatedAt: string
}

/**
 * 客户端配置
 */
export interface SsoClient {
  id: string
  clientId: string
  clientSecret: string
  clientName: string
  redirectUris: string[]
  protocolType: string
  enabled: boolean
  config: Record<string, any>
  createdAt: string
  updatedAt: string
}

/**
 * 密钥
 */
export interface SsoKey {
  id: string
  name: string
  keyType: string
  algorithm: string
  publicKey?: string
  privateKey?: string
  certificate?: string
  validFrom: string
  validTo: string
  enabled: boolean
  createdAt: string
  updatedAt: string
}

/**
 * 字段映射
 */
export interface SsoFieldMapping {
  id: string
  protocolType: string
  externalField: string
  internalField: string
  transformation?: string
  required: boolean
  createdAt: string
  updatedAt: string
}

/**
 * 会话
 */
export interface SsoSession {
  id: string
  sessionId: string
  userId: string
  protocolType: string
  clientId: string
  loginTime: string
  lastActivity: string
  expiresAt: string
  isActive: boolean
  metadata?: Record<string, any>
}
