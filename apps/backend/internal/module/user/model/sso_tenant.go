package model

import "time"

// SsoTenant 租户模型
type SsoTenant struct {
	ID         string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id"`
	TenantName string    `json:"tenant_name" form:"tenant_name" gorm:"size:128;default:'';column:tenant_name"`
	TenantCode string    `json:"tenant_code" form:"tenant_code" gorm:"size:64;default:'';column:tenant_code"`
	Linkman    string    `json:"linkman" form:"linkman" gorm:"size:64;default:'';column:linkman"`
	Contact    string    `json:"contact" form:"contact" gorm:"size:64;default:'';column:contact"`
	Address    string    `json:"address" form:"address" gorm:"size:128;default:'';column:address"`
	State      int       `json:"state" form:"state" gorm:"not null;default:1;column:state"`
	Remark     string    `json:"remark" form:"remark" gorm:"size:256;default:'';column:remark"`
	IsDeleted  bool      `json:"is_deleted" form:"is_deleted" gorm:"type:tinyint(1);default:0;column:is_deleted"`
	TenantID   string    `json:"tenant_id" form:"tenant_id" gorm:"type:varchar(64);not null;default:'0';column:tenant_id"`
	CreateID   string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';column:create_id"`
	CreateBy   string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';column:create_by"`
	CreateAt   time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime"`
	UpdateID   string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';column:update_id"`
	UpdateBy   string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';column:update_by"`
	UpdateAt   time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime"`
}

// TableName 指定表名
func (SsoTenant) TableName() string {
	return "sso_tenant"
}
