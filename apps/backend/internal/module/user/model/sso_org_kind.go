package model

import "time"

// SsoOrgKind 组织类型表模型
type SsoOrgKind struct {
	ID        string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	ParentID  string    `json:"parent_id" form:"parent_id" gorm:"type:varchar(64);not null;default:'';column:parent_id;comment:父ID"`
	AppCode   string    `json:"app_code" form:"app_code" gorm:"size:64;column:app_code;comment:应用编码"`
	KindName  string    `json:"kind_name" form:"kind_name" gorm:"size:100;default:'';column:kind_name;comment:分类名称"`
	KindCode  string    `json:"kind_code" form:"kind_code" gorm:"size:64;uniqueIndex;default:'';column:kind_code;comment:分类编码"`
	KindTag   string    `json:"kind_tag" form:"kind_tag" gorm:"size:64;default:'';column:kind_tag"`
	Status    int       `json:"status" form:"status" gorm:"not null;default:1;column:status;comment:状态"`
	Remark    string    `json:"remark" form:"remark" gorm:"size:512;default:'';column:remark;comment:备注"`
	Sort      int       `json:"sort" form:"sort" gorm:"default:0;column:sort;comment:序号"`
	IsDeleted bool      `json:"is_deleted" form:"is_deleted" gorm:"type:tinyint(1);default:0;column:is_deleted;comment:是否删除"`
	TenantID  string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'1';column:tenant_id;comment:租户ID"`
	CreateID  string    `json:"create_id" form:"create_id" gorm:"size:64;default:'';column:create_id;comment:创建人ID"`
	CreateBy  string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';column:create_by;comment:创建人"`
	CreateAt  time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID  string    `json:"update_id" form:"update_id" gorm:"size:64;default:'';column:update_id;comment:更新人ID"`
	UpdateBy  string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';column:update_by;comment:更新人"`
	UpdateAt  time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoOrgKind) TableName() string {
	return "sso_org_kind"
}
