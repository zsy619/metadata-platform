package model

import "time"

// SsoOrg 组织表
type SsoOrg struct {
	ID         string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	TenantID   string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';column:tenant_id;comment:租户ID"`
	ParentID   string    `json:"parent_id" form:"parent_id" gorm:"type:varchar(64);default:'0';column:parent_id;comment:父ID"`
	FromID     string    `json:"from_id" form:"from_id" gorm:"size:64;column:from_id"`
	AppCode    string    `json:"app_code" form:"app_code" gorm:"size:64;column:app_code;comment:应用编码"`
	OrgName    string    `json:"org_name" form:"org_name" gorm:"size:128;column:org_name;comment:组织名称"`
	OrgShort   string    `json:"org_short" form:"org_short" gorm:"size:128;column:org_short;comment:组织简称"`
	OrgEn      string    `json:"org_en" form:"org_en" gorm:"size:128;column:org_en"`
	OrgEnShort string    `json:"org_en_short" form:"org_en_short" gorm:"size:128;column:org_en_short"`
	OrgCode    string    `json:"org_code" form:"org_code" gorm:"size:64;uniqueIndex;column:org_code;comment:组织编码"`
	KindCode   string    `json:"kind_code" form:"kind_code" gorm:"size:64;column:kind_code;comment:类型编码"`
	Logo       string    `json:"logo" form:"logo" gorm:"size:256;column:logo"`
	Host       string    `json:"host" form:"host" gorm:"size:512;column:host"`
	Contact    string    `json:"contact" form:"contact" gorm:"size:128;column:contact"`
	Phone      string    `json:"phone" form:"phone" gorm:"size:128;column:phone"`
	Address    string    `json:"address" form:"address" gorm:"size:256;column:address"`
	Postcode   string    `json:"postcode" form:"postcode" gorm:"size:16;column:postcode"`
	Status     int       `json:"status" form:"status" gorm:"default:1;column:status;comment:状态"`
	Remark     string    `json:"remark" form:"remark" gorm:"size:512;column:remark;comment:备注"`
	Sort       int       `json:"sort" form:"sort" gorm:"default:0;column:sort;comment:排序"`
	IsDeleted  bool      `json:"is_deleted" form:"is_deleted" gorm:"type:tinyint(1);default:0;column:is_deleted;comment:是否删除"`
	CreateID   string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';column:create_id;comment:创建人ID"`
	CreateBy   string    `json:"create_by" form:"create_by" gorm:"size:64;column:create_by;comment:创建人"`
	CreateAt   time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID   string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';column:update_id;comment:更新人ID"`
	UpdateBy   string    `json:"update_by" form:"update_by" gorm:"size:64;column:update_by;comment:更新人"`
	UpdateAt   time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoOrg) TableName() string {
	return "sso_org"
}
