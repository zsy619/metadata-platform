package model

import "time"

// SsoUserSocial 用户第三方账号绑定模型（与 sso_user 一对多）
type SsoUserSocial struct {
	ID          string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	UserID      string    `json:"user_id" form:"user_id" gorm:"type:varchar(64);not null;index;column:user_id;comment:用户ID"`
	Provider    string    `json:"provider" form:"provider" gorm:"size:32;not null;column:provider;comment:平台(wechat/github/qq/dingtalk)"`
	OpenID      string    `json:"open_id" form:"open_id" gorm:"size:256;not null;uniqueIndex:idx_provider_openid;column:open_id;comment:平台用户唯一标识"`
	UnionID     string    `json:"union_id" form:"union_id" gorm:"size:256;default:'';column:union_id;comment:平台全局唯一标识"`
	Nickname    string    `json:"nickname" form:"nickname" gorm:"size:128;default:'';column:nickname;comment:平台昵称"`
	Avatar      string    `json:"avatar" form:"avatar" gorm:"size:512;default:'';column:avatar;comment:平台头像"`
	ProfileJSON string    `json:"profile_json" form:"profile_json" gorm:"type:text;column:profile_json;comment:平台原始数据(JSON)"`
	TenantID    string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'1';column:tenant_id;comment:租户ID"`
	CreateID    string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';column:create_id;comment:创建人ID"`
	CreateBy    string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';column:create_by;comment:创建人"`
	CreateAt    time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID    string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';column:update_id;comment:更新人ID"`
	UpdateBy    string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';column:update_by;comment:更新人"`
	UpdateAt    time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoUserSocial) TableName() string {
	return "sso_user_social"
}
