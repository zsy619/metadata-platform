package model

import "time"

// SsoOrgMenu 菜单和部门关联表模型
type SsoOrgMenu struct {
	ID        string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id"`
	MenuID    string    `json:"menu_id" form:"menu_id" gorm:"type:varchar(64);default:'0';uniqueIndex:idx_org_menu;column:menu_id;comment:菜单ID"`
	OrgID     string    `json:"org_id" form:"org_id" gorm:"type:varchar(64);default:'0';uniqueIndex:idx_org_menu;column:org_id;comment:组织ID"`
	Remark    string    `json:"remark" form:"remark" gorm:"size:1024;default:'';column:remark;comment:备注"`
	IsDeleted bool      `json:"is_deleted" form:"is_deleted" gorm:"type:tinyint(1);default:0;column:is_deleted;comment:是否删除"`
	IsSystem  bool      `json:"is_system" form:"is_system" gorm:"type:tinyint(1);default:0;comment:是否系统内置"`
	TenantID  string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'1';column:tenant_id;comment:租户ID"`
	CreateID  string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';column:create_id;comment:创建人ID"`
	CreateBy  string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';column:create_by;comment:创建人"`
	CreateAt  time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID  string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';column:update_id;comment:更新人ID"`
	UpdateBy  string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';column:update_by;comment:更新人"`
	UpdateAt  time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoOrgMenu) TableName() string {
	return "sso_org_menu"
}
