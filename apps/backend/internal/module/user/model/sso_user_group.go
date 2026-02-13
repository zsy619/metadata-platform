package model

import "time"

// SsoUserGroup 用户组表模型
type SsoUserGroup struct {
	ID        string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	ParentID  string    `json:"parent_id" form:"parent_id" gorm:"type:varchar(64);default:'';column:parent_id;comment:父ID"`
	AppCode   string    `json:"app_code" form:"app_code" gorm:"size:64;default:'';column:app_code;comment:应用编码"`
	OrgID     string    `json:"org_id" form:"org_id" gorm:"type:varchar(64);default:'0';column:org_id;comment:组织ID"`
	KindCode  string    `json:"kind_code" form:"kind_code" gorm:"size:64;default:'';column:kind_code;comment:分类编码"`
	GroupName string    `json:"group_name" form:"group_name" gorm:"size:128;default:'';column:group_name;comment:组名称"`
	GroupCode string    `json:"group_code" form:"group_code" gorm:"size:128;uniqueIndex;default:'';column:group_code;comment:组编码"`
	Status    int       `json:"status" form:"status" gorm:"not null;default:1;column:status;comment:状态"`
	Remark    string    `json:"remark" form:"remark" gorm:"size:512;default:'';column:remark;comment:备注"`
	Sort      int       `json:"sort" form:"sort" gorm:"default:0;column:sort;comment:排序"`
	IsDeleted bool      `json:"is_deleted" form:"is_deleted" gorm:"type:tinyint(1);default:0;column:is_deleted;comment:是否删除"`
	TenantID  string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';column:tenant_id;comment:租户ID"`
	CreateID  string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';column:create_id;comment:创建人ID"`
	CreateBy  string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';column:create_by;comment:创建人"`
	CreateAt  time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID  string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';column:update_id;comment:更新人ID"`
	UpdateBy  string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';column:update_by;comment:更新人"`
	UpdateAt  time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoUserGroup) TableName() string {
	return "sso_user_group"
}
