package model

import "time"

// SsoTenant 租户模型
type SsoTenant struct {
	ID         string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	TenantName string    `json:"tenant_name" form:"tenant_name" gorm:"size:128;default:'';column:tenant_name;comment:租户名称"`
	TenantCode string    `json:"tenant_code" form:"tenant_code" gorm:"size:64;default:'';column:tenant_code;comment:租户编码"`
	Linkman    string    `json:"linkman" form:"linkman" gorm:"size:64;default:'';column:linkman;comment:联系人"`
	Contact    string    `json:"contact" form:"contact" gorm:"size:64;default:'';column:contact;comment:联系方式"`
	Address    string    `json:"address" form:"address" gorm:"size:128;default:'';column:address;comment:地址"`
	State      int       `json:"state" form:"state" gorm:"not null;default:1;column:state;comment:状态"`
	Remark     string    `json:"remark" form:"remark" gorm:"size:256;default:'';column:remark;comment:备注"`
	IsDeleted  bool      `json:"is_deleted" form:"is_deleted" gorm:"type:tinyint(1);default:0;column:is_deleted;comment:是否删除"`
	TenantID   string    `json:"tenant_id" form:"tenant_id" gorm:"type:varchar(64);not null;default:'0';column:tenant_id;comment:租户ID"`
	CreateID   string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';column:create_id;comment:创建人ID"`
	CreateBy   string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';column:create_by;comment:创建人"`
	CreateAt   time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID   string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';column:update_id;comment:更新人ID"`
	UpdateBy   string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';column:update_by;comment:更新人"`
	UpdateAt   time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoTenant) TableName() string {
	return "sso_tenant"
}
