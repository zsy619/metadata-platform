package model

import "time"

// ClientType 客户端类型
type ClientType string

const (
	ClientTypeWeb     ClientType = "web"     // Web应用
	ClientTypeSPA     ClientType = "spa"     // 单页应用
	ClientTypeMobile  ClientType = "mobile"  // 移动应用
	ClientTypeBackend ClientType = "backend" // 后端服务
)

// SsoClient SSO客户端配置模型
type SsoClient struct {
	ID             string     `json:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	ClientID       string     `json:"client_id" gorm:"size:128;uniqueIndex;not null;column:client_id;comment:客户端ID"`
	ClientName     string     `json:"client_name" gorm:"size:128;not null;column:client_name;comment:客户端名称"`
	ClientType     ClientType `json:"client_type" gorm:"size:32;not null;column:client_type;comment:客户端类型"`
	ProtocolConfigID string   `json:"protocol_config_id" gorm:"type:varchar(64);index;column:protocol_config_id;comment:关联的协议配置ID"`
	ClientSecret   string     `json:"client_secret" gorm:"size:256;column:client_secret;comment:客户端密钥(加密存储)"`
	RedirectURIs   string     `json:"redirect_uris" gorm:"type:text;column:redirect_uris;comment:重定向URI列表(JSON数组)"`
	PostLogoutRedirectURIs string `json:"post_logout_redirect_uris" gorm:"type:text;column:post_logout_redirect_uris;comment:登出重定向URI列表(JSON数组)"`
	Scopes         string     `json:"scopes" gorm:"type:text;column:scopes;comment:允许的范围(JSON数组)"`
	GrantTypes     string     `json:"grant_types" gorm:"type:text;column:grant_types;comment:授权类型(JSON数组)"`
	ResponseTypes  string     `json:"response_types" gorm:"type:text;column:response_types;comment:响应类型(JSON数组)"`
	AppLogo        string     `json:"app_logo" gorm:"size:256;default:'';column:app_logo;comment:应用Logo"`
	AppDescription string     `json:"app_description" gorm:"size:512;default:'';column:app_description;comment:应用描述"`
	HomepageURL    string     `json:"homepage_url" gorm:"size:512;default:'';column:homepage_url;comment:应用主页"`
	IsPublic       bool       `json:"is_public" gorm:"default:false;column:is_public;comment:是否公开客户端"`
	RequireConsent bool       `json:"require_consent" gorm:"default:true;column:require_consent;comment:是否需要用户同意"`
	IsEnabled      bool       `json:"is_enabled" gorm:"default:true;column:is_enabled;comment:是否启用"`
	Status         int        `json:"status" gorm:"not null;default:1;column:status;comment:状态"`
	Remark         string     `json:"remark" gorm:"size:512;default:'';column:remark;comment:备注"`
	Sort           int        `json:"sort" gorm:"default:0;column:sort;comment:排序"`
	IsDeleted      bool       `json:"is_deleted" gorm:"default:false;column:is_deleted;comment:是否删除"`
	TenantID       string     `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'1';column:tenant_id;comment:租户ID"`
	CreateID       string     `json:"create_id" gorm:"size:64;default:'';column:create_id;comment:创建人ID"`
	CreateBy       string     `json:"create_by" gorm:"size:64;default:'';column:create_by;comment:创建人"`
	CreateAt       time.Time  `json:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID       string     `json:"update_id" gorm:"size:64;default:'';column:update_id;comment:更新人ID"`
	UpdateBy       string     `json:"update_by" gorm:"size:64;default:'';column:update_by;comment:更新人"`
	UpdateAt       time.Time  `json:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoClient) TableName() string {
	return "sso_client"
}
