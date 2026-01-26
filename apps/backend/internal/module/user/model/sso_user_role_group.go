package model

import "time"

// SsoUserRoleGroup 用户角色组关联表模型
type SsoUserRoleGroup struct {
	ID        string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id"`
	UserID    string    `json:"user_id" form:"user_id" gorm:"type:varchar(64);default:'0';uniqueIndex:idx_user_role_group;column:user_id"`
	GroupID   string    `json:"group_id" form:"group_id" gorm:"type:varchar(64);default:'0';uniqueIndex:idx_user_role_group;column:group_id"`
	IsDeleted bool      `json:"is_deleted" form:"is_deleted" gorm:"type:tinyint(1);default:0;column:is_deleted"`
	TenantID  string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';column:tenant_id"`
	CreateID  string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';column:create_id"`
	CreateBy  string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';column:create_by"`
	CreateAt  time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime"`
	UpdateID  string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';column:update_id"`
	UpdateBy  string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';column:update_by"`
	UpdateAt  time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime"`
}

// TableName 指定表名
func (SsoUserRoleGroup) TableName() string {
	return "sso_user_role_group"
}
