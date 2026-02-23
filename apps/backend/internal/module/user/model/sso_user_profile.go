package model

import "time"

// SsoUserProfile 用户档案模型（与 sso_user 一对一）
type SsoUserProfile struct {
	ID       string     `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	UserID   string     `json:"user_id" form:"user_id" gorm:"type:varchar(64);uniqueIndex;not null;column:user_id;comment:用户ID"`
	Nickname string     `json:"nickname" form:"nickname" gorm:"size:64;default:'';column:nickname;comment:昵称"`
	Avatar   string     `json:"avatar" form:"avatar" gorm:"size:512;default:'';column:avatar;comment:头像URL"`
	Gender   string     `json:"gender" form:"gender" gorm:"type:char(1);default:'';column:gender;comment:性别(M男/F女/U未知)"`
	Birthday *time.Time `json:"birthday" form:"birthday" gorm:"type:date;column:birthday;comment:生日"`
	Bio      string     `json:"bio" form:"bio" gorm:"type:text;column:bio;comment:个人简介"`
	Location string     `json:"location" form:"location" gorm:"size:256;default:'';column:location;comment:所在地"`
	TenantID string     `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'1';column:tenant_id;comment:租户ID"`
	CreateID string     `json:"create_id" form:"create_id" gorm:"size:64;default:'';column:create_id;comment:创建人ID"`
	CreateBy string     `json:"create_by" form:"create_by" gorm:"size:64;default:'';column:create_by;comment:创建人"`
	CreateAt time.Time  `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID string     `json:"update_id" form:"update_id" gorm:"size:64;default:'';column:update_id;comment:更新人ID"`
	UpdateBy string     `json:"update_by" form:"update_by" gorm:"size:64;default:'';column:update_by;comment:更新人"`
	UpdateAt time.Time  `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoUserProfile) TableName() string {
	return "sso_user_profile"
}
