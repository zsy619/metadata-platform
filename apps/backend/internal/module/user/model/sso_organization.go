package model

import "time"

// SsoOrganization 组织表
type SsoOrganization struct {
	ID              string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	TenantID        string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';column:tenant_id;comment:租户ID"`
	ParentID        string    `json:"parent_id" form:"parent_id" gorm:"type:varchar(64);default:'0';column:parent_id;comment:父ID"`
	FromID          string    `json:"from_id" form:"from_id" gorm:"size:64;column:from_id"`
	ApplicationCode string    `json:"application_code" form:"application_code" gorm:"size:64;column:application_code;comment:应用编码"`
	UnitName        string    `json:"unit_name" form:"unit_name" gorm:"size:128;column:unit_name;comment:单位名称"`
	UnitShort       string    `json:"unit_short" form:"unit_short" gorm:"size:128;column:unit_short;comment:单位简称"`
	UnitEn          string    `json:"unit_en" form:"unit_en" gorm:"size:128;column:unit_en"`
	UnitEnShort     string    `json:"unit_en_short" form:"unit_en_short" gorm:"size:128;column:unit_en_short"`
	UnitCode        string    `json:"unit_code" form:"unit_code" gorm:"size:64;uniqueIndex;column:unit_code;comment:单位编码"`
	KindCode        string    `json:"kind_code" form:"kind_code" gorm:"size:64;column:kind_code;comment:类型编码"`
	Logo            string    `json:"logo" form:"logo" gorm:"size:256;column:logo"`
	Host            string    `json:"host" form:"host" gorm:"size:512;column:host"`
	Contact         string    `json:"contact" form:"contact" gorm:"size:128;column:contact"`
	Phone           string    `json:"phone" form:"phone" gorm:"size:128;column:phone"`
	Address         string    `json:"address" form:"address" gorm:"size:256;column:address"`
	Postcode        string    `json:"postcode" form:"postcode" gorm:"size:16;column:postcode"`
	State           int       `json:"state" form:"state" gorm:"default:1;column:state;comment:状态"`
	Remark          string    `json:"remark" form:"remark" gorm:"size:512;column:remark;comment:备注"`
	Sort            int       `json:"sort" form:"sort" gorm:"default:0;column:sort;comment:排序"`
	IsDeleted       bool      `json:"is_deleted" form:"is_deleted" gorm:"type:tinyint(1);default:0;column:is_deleted;comment:是否删除"`
	CreateID        string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';column:create_id;comment:创建人ID"`
	CreateBy        string    `json:"create_by" form:"create_by" gorm:"size:64;column:create_by;comment:创建人"`
	CreateAt        time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID        string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';column:update_id;comment:更新人ID"`
	UpdateBy        string    `json:"update_by" form:"update_by" gorm:"size:64;column:update_by;comment:更新人"`
	UpdateAt        time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoOrganization) TableName() string {
	return "sso_organization"
}
