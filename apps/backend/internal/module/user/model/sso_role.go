package model

import "time"

// SsoRole 角色管理模型
type SsoRole struct {
	ID              string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	ParentID        string    `json:"parent_id" form:"parent_id" gorm:"type:varchar(64);default:'';column:parent_id;comment:父ID"`
	ApplicationCode string    `json:"application_code" form:"application_code" gorm:"size:64;default:'';column:application_code;comment:应用编码"`
	OrganizationID  string    `json:"organization_id" form:"organization_id" gorm:"type:varchar(64);default:'0';column:organization_id;comment:组织ID"`
	KindCode        string    `json:"kind_code" form:"kind_code" gorm:"size:64;default:'';column:kind_code;comment:分类编码"`
	RoleName        string    `json:"role_name" form:"role_name" gorm:"size:128;default:'';column:role_name;comment:角色名称"`
	RoleCode        string    `json:"role_code" form:"role_code" gorm:"size:128;uniqueIndex;default:'';column:role_code;comment:角色编码"`
	State           int       `json:"state" form:"state" gorm:"not null;default:1;column:state;comment:状态"`
	DataScope       string    `json:"data_scope" form:"data_scope" gorm:"type:char(1);default:'1';column:data_scope;comment:数据范围"`
	Remark          string    `json:"remark" form:"remark" gorm:"size:512;default:'';column:remark;comment:备注"`
	Sort            int       `json:"sort" form:"sort" gorm:"default:0;column:sort;comment:排序"`
	IsDeleted       bool      `json:"is_deleted" form:"is_deleted" gorm:"type:tinyint(1);default:0;column:is_deleted;comment:是否删除"`
	TenantID        string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';column:tenant_id;comment:租户ID"`
	CreateID        string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';column:create_id;comment:创建人ID"`
	CreateBy        string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';column:create_by;comment:创建人"`
	CreateAt        time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID        string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';column:update_id;comment:更新人ID"`
	UpdateBy        string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';column:update_by;comment:更新人"`
	UpdateAt        time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoRole) TableName() string {
	return "sso_role"
}
