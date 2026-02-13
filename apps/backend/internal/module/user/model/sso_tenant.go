package model

import "time"

// SsoTenant 租户模型
type SsoTenant struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(64)"`
	ParentID   string    `json:"parent_id" form:"parent_id" gorm:"type:varchar(64);default:'';column:parent_id;comment:父ID"`
	TenantName string    `json:"tenant_name" gorm:"type:varchar(128);not null;comment:租户名称"`
	TenantCode string    `json:"tenant_code" gorm:"type:varchar(64);not null;uniqueIndex;comment:租户编码"`
	Status     int       `json:"status" gorm:"type:tinyint;default:1;comment:状态: 1=有效, 0=禁用"`
	Remark     string    `json:"remark" gorm:"type:varchar(512);comment:备注"`
	IsDeleted  bool      `json:"is_deleted" gorm:"type:tinyint(1);default:0;comment:是否删除"`
	TenantID   string    `json:"tenant_id" gorm:"type:varchar(64);not null;default:'0';comment:租户ID"`
	CreateID   string    `json:"create_id" gorm:"type:varchar(64);default:'0';comment:创建人ID"`
	CreateBy   string    `json:"create_by" gorm:"type:varchar(64);comment:创建人"`
	CreateAt   time.Time `json:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID   string    `json:"update_id" gorm:"type:varchar(64);default:'0';comment:更新人ID"`
	UpdateBy   string    `json:"update_by" gorm:"type:varchar(64);comment:修改人"`
	UpdateAt   time.Time `json:"update_at" gorm:"autoUpdateTime;comment:修改时间"`
}

// TableName 指定表名
func (SsoTenant) TableName() string {
	return "sso_tenant"
}
