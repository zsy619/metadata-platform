package model

import "time"

// SsoRoleGroupRole 角色组关联角色表模型
type SsoRoleGroupRole struct {
	ID        string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id"`
	GroupID   string    `json:"group_id" form:"group_id" gorm:"type:varchar(64);default:'';uniqueIndex:idx_group_role;column:group_id"`
	RoleID    string    `json:"role_id" form:"role_id" gorm:"type:varchar(64);default:'';uniqueIndex:idx_group_role;column:role_id"`
	Remark    string    `json:"remark" gorm:"size:1024;default:'';column:remark;comment:备注"`
	IsDeleted bool      `json:"is_deleted" form:"is_deleted" gorm:"default:0;column:is_deleted"`
	IsSystem  bool      `json:"is_system" form:"is_system" gorm:"default:0;comment:是否系统内置"`
	TenantID  string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'';column:tenant_id"`
	CreateID  string    `json:"create_id" form:"create_id" gorm:"size:64;default:'';column:create_id"`
	CreateBy  string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';column:create_by"`
	CreateAt  time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime"`
	UpdateID  string    `json:"update_id" form:"update_id" gorm:"size:64;default:'';column:update_id"`
	UpdateBy  string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';column:update_by"`
	UpdateAt  time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime"`
}

// TableName 指定表名
func (SsoRoleGroupRole) TableName() string {
	return "sso_role_group_role"
}
