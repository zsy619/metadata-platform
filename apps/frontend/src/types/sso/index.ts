/**
 * SSO单点登录模块类型定义
 */

/**
 * 协议类型枚举
 */
export type ProtocolType = 'oidc' | 'saml' | 'ldap' | 'cas'

/**
 * 客户端类型枚举
 */
export type ClientType = 'web' | 'spa' | 'mobile' | 'backend'

/**
 * 密钥类型枚举
 */
export type KeyType = 'rsa' | 'ec' | 'octet' | 'cert'

/**
 * 密钥用途枚举
 */
export type KeyUsage = 'signing' | 'encryption' | 'both'

/**
 * 会话状态枚举
 */
export type SessionStatus = 'active' | 'inactive' | 'expired' | 'revoked'

/**
 * SSO协议配置
 */
export interface SsoProtocolConfig {
    id: string
    config_name: string
    protocol_type: ProtocolType
    is_enabled: boolean
    config_data: string
    remark?: string
    sort: number
    is_deleted: boolean
    tenant_id: string
    create_id: string
    create_by: string
    create_at: string
    update_id: string
    update_by: string
    update_at: string
}

/**
 * SSO客户端配置
 */
export interface SsoClient {
    id: string
    client_id: string
    client_name: string
    client_type: ClientType
    protocol_config_id: string
    client_secret?: string
    redirect_uris: string
    post_logout_redirect_uris?: string
    scopes?: string
    grant_types?: string
    response_types?: string
    app_logo?: string
    app_description?: string
    homepage_url?: string
    is_public: boolean
    require_consent: boolean
    is_enabled: boolean
    status: number
    remark?: string
    sort: number
    is_deleted: boolean
    tenant_id: string
    create_id: string
    create_by: string
    create_at: string
    update_id: string
    update_by: string
    update_at: string
}

/**
 * SSO密钥
 */
export interface SsoKey {
    id: string
    key_id: string
    key_name: string
    key_type: KeyType
    key_usage: KeyUsage
    algorithm: string
    public_key?: string
    private_key?: string
    secret_key?: string
    certificate?: string
    is_primary: boolean
    is_enabled: boolean
    valid_from: string
    valid_to: string
    protocol_config_id?: string
    remark?: string
    is_deleted: boolean
    tenant_id: string
    create_id: string
    create_by: string
    create_at: string
    update_id: string
    update_by: string
    update_at: string
}

/**
 * SSO字段映射
 */
export interface SsoFieldMapping {
    id: string
    mapping_name: string
    protocol_config_id: string
    client_id?: string
    source_field: string
    target_field: string
    field_type: string
    is_required: boolean
    default_value?: string
    transform_script?: string
    is_enabled: boolean
    sort: number
    remark?: string
    is_deleted: boolean
    tenant_id: string
    create_id: string
    create_by: string
    create_at: string
    update_id: string
    update_by: string
    update_at: string
}

/**
 * SSO会话
 */
export interface SsoSession {
    id: string
    session_id: string
    user_id: string
    client_id?: string
    protocol_config_id?: string
    protocol_type?: ProtocolType
    access_token?: string
    refresh_token?: string
    id_token?: string
    status: SessionStatus
    auth_time: string
    expires_at: string
    last_activity_at: string
    ip_address?: string
    user_agent?: string
    device_info?: string
    location_info?: string
    scopes?: string
    is_deleted: boolean
    tenant_id: string
    create_at: string
    update_at: string
}
