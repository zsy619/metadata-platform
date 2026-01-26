package model

import "time"

// SsoOrganization 组织表
type SsoOrganization struct {
	ID              string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id"`
	TenantID        string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';column:tenant_id"`
	ParentID        string    `json:"parent_id" form:"parent_id" gorm:"type:varchar(64);default:'0';column:parent_id"`
	FromID          string    `json:"from_id" form:"from_id" gorm:"size:64;column:from_id"`
	ApplicationCode string    `json:"application_code" form:"application_code" gorm:"size:64;column:application_code"`
	UnitName        string    `json:"unit_name" form:"unit_name" gorm:"size:128;column:unit_name"`
	UnitShort       string    `json:"unit_short" form:"unit_short" gorm:"size:128;column:unit_short"`
	UnitEn          string    `json:"unit_en" form:"unit_en" gorm:"size:128;column:unit_en"`
	UnitEnShort     string    `json:"unit_en_short" form:"unit_en_short" gorm:"size:128;column:unit_en_short"`
	UnitCode        string    `json:"unit_code" form:"unit_code" gorm:"size:64;uniqueIndex;column:unit_code"`
	KindCode        string    `json:"kind_code" form:"kind_code" gorm:"size:64;column:kind_code"`
	Logo            string    `json:"logo" form:"logo" gorm:"size:256;column:logo"`
	Host            string    `json:"host" form:"host" gorm:"size:512;column:host"`
	Contact         string    `json:"contact" form:"contact" gorm:"size:128;column:contact"`
	Phone           string    `json:"phone" form:"phone" gorm:"size:128;column:phone"`
	Address         string    `json:"address" form:"address" gorm:"size:256;column:address"`
	Postcode        string    `json:"postcode" form:"postcode" gorm:"size:16;column:postcode"`
	State           int       `json:"state" form:"state" gorm:"default:1;column:state"`
	Remark          string    `json:"remark" form:"remark" gorm:"size:512;column:remark"`
	Sort            int       `json:"sort" form:"sort" gorm:"default:0;column:sort"`
	IsDeleted       bool      `json:"is_deleted" form:"is_deleted" gorm:"type:tinyint(1);default:0;column:is_deleted"`
	CreateID        string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';column:create_id"`
	CreateBy        string    `json:"create_by" form:"create_by" gorm:"size:64;column:create_by"`
	CreateAt        time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime"`
	UpdateID        string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';column:update_id"`
	UpdateBy        string    `json:"update_by" form:"update_by" gorm:"size:64;column:update_by"`
	UpdateAt        time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime"`
}

// TableName 指定表名
func (SsoOrganization) TableName() string {
	return "sso_organization"
}
