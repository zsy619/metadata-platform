package model

import "time"

// SsoUserContact 用户额外联系方式模型（与 sso_user 一对多）
type SsoUserContact struct {
	ID         string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	UserID     string    `json:"user_id" form:"user_id" gorm:"type:varchar(64);not null;index;column:user_id;comment:用户ID"`
	Type       string    `json:"type" form:"type" gorm:"size:16;not null;column:type;comment:联系方式类型(email/phone/wechat/qq)"`
	Value      string    `json:"value" form:"value" gorm:"size:256;not null;column:value;comment:联系方式值"`
	IsVerified bool      `json:"is_verified" form:"is_verified" gorm:"type:tinyint(1);default:0;column:is_verified;comment:是否已验证"`
	Remark     string    `json:"remark" form:"remark" gorm:"size:256;default:'';column:remark;comment:备注"`
	TenantID   string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'1';column:tenant_id;comment:租户ID"`
	CreateID   string    `json:"create_id" form:"create_id" gorm:"size:64;default:'';column:create_id;comment:创建人ID"`
	CreateBy   string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';column:create_by;comment:创建人"`
	CreateAt   time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID   string    `json:"update_id" form:"update_id" gorm:"size:64;default:'';column:update_id;comment:更新人ID"`
	UpdateBy   string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';column:update_by;comment:更新人"`
	UpdateAt   time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoUserContact) TableName() string {
	return "sso_user_contact"
}
