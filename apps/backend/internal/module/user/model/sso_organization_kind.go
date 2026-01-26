package model

import "time"

// SsoOrganizationKind 组织类型表模型
type SsoOrganizationKind struct {
	ID              string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id"`
	ParentID        string    `json:"parent_id" form:"parent_id" gorm:"type:varchar(64);not null;default:'';column:parent_id"`
	ApplicationCode string    `json:"application_code" form:"application_code" gorm:"size:64;column:application_code"`
	KindName        string    `json:"kind_name" form:"kind_name" gorm:"size:100;default:'';column:kind_name"`
	KindCode        string    `json:"kind_code" form:"kind_code" gorm:"size:64;uniqueIndex;default:'';column:kind_code"`
	KindTag         string    `json:"kind_tag" form:"kind_tag" gorm:"size:64;default:'';column:kind_tag"`
	State           int       `json:"state" form:"state" gorm:"not null;default:1;column:state"`
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
func (SsoOrganizationKind) TableName() string {
	return "sso_organization_kind"
}
