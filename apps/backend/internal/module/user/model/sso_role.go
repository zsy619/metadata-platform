package model

import "time"

// SsoRole 角色管理模型
type SsoRole struct {
	ID              string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id"`
	ParentID        string    `json:"parent_id" form:"parent_id" gorm:"type:varchar(64);default:'';column:parent_id"`
	ApplicationCode string    `json:"application_code" form:"application_code" gorm:"size:64;default:'';column:application_code"`
	OrganizationID  string    `json:"organization_id" form:"organization_id" gorm:"type:varchar(64);default:'0';column:organization_id"`
	KindCode        string    `json:"kind_code" form:"kind_code" gorm:"size:64;default:'';column:kind_code"`
	RoleName        string    `json:"role_name" form:"role_name" gorm:"size:128;default:'';column:role_name"`
	RoleCode        string    `json:"role_code" form:"role_code" gorm:"size:128;uniqueIndex;default:'';column:role_code"`
	State           int       `json:"state" form:"state" gorm:"not null;default:1;column:state"`
	DataScope       string    `json:"data_scope" form:"data_scope" gorm:"type:char(1);default:'1';column:data_scope"`
	Remark          string    `json:"remark" form:"remark" gorm:"size:512;default:'';column:remark"`
	Sort            int       `json:"sort" form:"sort" gorm:"default:0;column:sort"`
	IsDeleted       bool      `json:"is_deleted" form:"is_deleted" gorm:"type:tinyint(1);default:0;column:is_deleted"`
	TenantID        string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';column:tenant_id"`
	CreateID        string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';column:create_id"`
	CreateBy        string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';column:create_by"`
	CreateAt        time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime"`
	UpdateID        string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';column:update_id"`
	UpdateBy        string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';column:update_by"`
	UpdateAt        time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime"`
}

// TableName 指定表名
func (SsoRole) TableName() string {
	return "sso_role"
}
