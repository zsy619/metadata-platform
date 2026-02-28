package model

import "time"

// ProtocolType 协议类型枚举
type ProtocolType string

const (
	ProtocolTypeOIDC  ProtocolType = "oidc"  // OIDC/OAuth 2.1
	ProtocolTypeSAML  ProtocolType = "saml"  // SAML 2.0
	ProtocolTypeLDAP  ProtocolType = "ldap"  // LDAP/LDAPS
	ProtocolTypeCAS   ProtocolType = "cas"   // CAS
)

// SsoProtocolConfig SSO协议配置模型
type SsoProtocolConfig struct {
	ID          string       `json:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	ConfigName  string       `json:"config_name" gorm:"size:128;not null;column:config_name;comment:配置名称"`
	ProtocolType ProtocolType `json:"protocol_type" gorm:"size:32;not null;index;column:protocol_type;comment:协议类型"`
	IsEnabled   bool         `json:"is_enabled" gorm:"default:true;column:is_enabled;comment:是否启用"`
	ConfigData  string       `json:"config_data" gorm:"type:text;column:config_data;comment:协议配置JSON数据"`
	Remark      string       `json:"remark" gorm:"size:512;default:'';column:remark;comment:备注"`
	Sort        int          `json:"sort" gorm:"default:0;column:sort;comment:排序"`
	IsDeleted   bool         `json:"is_deleted" gorm:"default:false;column:is_deleted;comment:是否删除"`
	TenantID    string       `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'1';column:tenant_id;comment:租户ID"`
	CreateID    string       `json:"create_id" gorm:"size:64;default:'';column:create_id;comment:创建人ID"`
	CreateBy    string       `json:"create_by" gorm:"size:64;default:'';column:create_by;comment:创建人"`
	CreateAt    time.Time    `json:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID    string       `json:"update_id" gorm:"size:64;default:'';column:update_id;comment:更新人ID"`
	UpdateBy    string       `json:"update_by" gorm:"size:64;default:'';column:update_by;comment:更新人"`
	UpdateAt    time.Time    `json:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoProtocolConfig) TableName() string {
	return "sso_protocol_config"
}

// OIDCConfig OIDC协议配置结构
type OIDCConfig struct {
	Issuer           string   `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint    string   `json:"token_endpoint"`
	UserinfoEndpoint string   `json:"userinfo_endpoint"`
	JwksURI          string   `json:"jwks_uri"`
	ScopesSupported  []string `json:"scopes_supported"`
	ResponseTypesSupported []string `json:"response_types_supported"`
	GrantTypesSupported []string `json:"grant_types_supported"`
	RequirePKCE      bool     `json:"require_pkce"`
	AccessTokenTTL   int      `json:"access_token_ttl"`   // 秒
	RefreshTokenTTL  int      `json:"refresh_token_ttl"`  // 秒
	IDTokenTTL       int      `json:"id_token_ttl"`       // 秒
	SigningAlgorithm string   `json:"signing_algorithm"`  // RS256, HS256等
}

// SAMLConfig SAML协议配置结构
type SAMLConfig struct {
	EntityID                  string `json:"entity_id"`
	SSOServiceURL             string `json:"sso_service_url"`
	SLOServiceURL             string `json:"slo_service_url"`
	AssertionConsumerServiceURL string `json:"assertion_consumer_service_url"`
	NameIDFormat              string `json:"name_id_format"`
	SignAssertions            bool   `json:"sign_assertions"`
	SignRequests              bool   `json:"sign_requests"`
	EncryptAssertions         bool   `json:"encrypt_assertions"`
	SignatureAlgorithm        string `json:"signature_algorithm"`
	AssertionTTL              int    `json:"assertion_ttl"` // 秒
	MetadataXML               string `json:"metadata_xml"`    // 元数据XML
}

// LDAPConfig LDAP协议配置结构
type LDAPConfig struct {
	ServerURL       string `json:"server_url"`       // ldap://或ldaps://
	BaseDN          string `json:"base_dn"`
	BindDN          string `json:"bind_dn"`
	BindPassword    string `json:"bind_password"`    // 加密存储
	UserSearchBase  string `json:"user_search_base"`
	UserSearchFilter string `json:"user_search_filter"`
	GroupSearchBase string `json:"group_search_base"`
	GroupSearchFilter string `json:"group_search_filter"`
	UserAttrMapping map[string]string `json:"user_attr_mapping"` // 属性映射
	UseTLS          bool   `json:"use_tls"`
	SkipTLSVerify   bool   `json:"skip_tls_verify"`
	ConnectionTimeout int   `json:"connection_timeout"` // 秒
	PoolSize        int    `json:"pool_size"`
}

// CASConfig CAS协议配置结构
type CASConfig struct {
	ServerURL      string `json:"server_url"`
	ServiceURL     string `json:"service_url"`
	ProtocolVersion string `json:"protocol_version"` // 1.0, 2.0, 3.0, 4.0
	EnableProxy    bool   `json:"enable_proxy"`
	ProxyCallbackURL string `json:"proxy_callback_url"`
	STTTL          int    `json:"st_ttl"`  // Service Ticket TTL 秒
	PTTTL          int    `json:"pt_ttl"`  // Proxy Ticket TTL 秒
}
