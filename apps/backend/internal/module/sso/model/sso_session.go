package model

import "time"

// SessionStatus 会话状态
type SessionStatus string

const (
	SessionStatusActive    SessionStatus = "active"     // 活跃
	SessionStatusInactive  SessionStatus = "inactive"   // 不活跃
	SessionStatusExpired   SessionStatus = "expired"    // 已过期
	SessionStatusRevoked   SessionStatus = "revoked"    // 已撤销
	SessionStatusPending   SessionStatus = "pending"    // 待处理
	SessionStatusLoggedOut SessionStatus = "logged_out" // 已登出
)

// SsoSession SSO会话模型
type SsoSession struct {
	ID               string        `json:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	SessionID        string        `json:"session_id" gorm:"size:128;uniqueIndex;not null;column:session_id;comment:会话ID"`
	UserID           string        `json:"user_id" gorm:"type:varchar(64);index;not null;column:user_id;comment:用户ID"`
	ClientID         string        `json:"client_id" gorm:"size:128;index;column:client_id;comment:客户端ID"`
	ProtocolConfigID string        `json:"protocol_config_id" gorm:"type:varchar(64);index;column:protocol_config_id;comment:协议配置ID"`
	ProtocolType     ProtocolType  `json:"protocol_type" gorm:"size:32;column:protocol_type;comment:协议类型"`
	AccessToken      string        `json:"access_token" gorm:"type:text;column:access_token;comment:访问令牌(加密存储)"`
	RefreshToken     string        `json:"refresh_token" gorm:"type:text;column:refresh_token;comment:刷新令牌(加密存储)"`
	IDToken          string        `json:"id_token" gorm:"type:text;column:id_token;comment:ID令牌"`
	Status           SessionStatus `json:"status" gorm:"size:32;default:'active';column:status;comment:会话状态"`
	AuthTime         time.Time     `json:"auth_time" gorm:"column:auth_time;comment:认证时间"`
	ExpiresAt        time.Time     `json:"expires_at" gorm:"index;column:expires_at;comment:过期时间"`
	LastActivityAt   time.Time     `json:"last_activity_at" gorm:"column:last_activity_at;comment:最后活动时间"`
	IPAddress        string        `json:"ip_address" gorm:"size:64;column:ip_address;comment:IP地址"`
	UserAgent        string        `json:"user_agent" gorm:"size:512;column:user_agent;comment:用户代理"`
	DeviceInfo       string        `json:"device_info" gorm:"type:text;column:device_info;comment:设备信息(JSON)"`
	LocationInfo     string        `json:"location_info" gorm:"type:text;column:location_info;comment:位置信息 (JSON)"`
	Scopes           string        `json:"scopes" gorm:"type:text;column:scopes;comment:授权范围 (JSON 数组)"`
	ExtraData        string        `json:"extra_data" gorm:"type:text;column:extra_data;comment:额外数据 (JSON)"`
	IsDeleted        bool          `json:"is_deleted" gorm:"default:false;column:is_deleted;comment:是否删除"`
	TenantID         string        `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'1';column:tenant_id;comment:租户ID"`
	CreateAt         time.Time     `json:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateAt         time.Time     `json:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoSession) TableName() string {
	return "sso_session"
}

// IsExpired 检查会话是否过期
func (s *SsoSession) IsExpired() bool {
	if s.Status == SessionStatusExpired || s.Status == SessionStatusRevoked {
		return true
	}
	if s.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(s.ExpiresAt)
}

// IsActive 检查会话是否活跃
func (s *SsoSession) IsActive() bool {
	if s.Status != SessionStatusActive {
		return false
	}
	return !s.IsExpired()
}
