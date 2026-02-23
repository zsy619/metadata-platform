package model

import "time"

// SsoMenu 菜单权限模型
type SsoMenu struct {
	ID        string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	ParentID  string    `json:"parent_id" form:"parent_id" gorm:"type:varchar(64);default:'';column:parent_id;comment:父ID"`
	AppCode   string    `json:"app_code" form:"app_code" gorm:"size:64;default:'';column:app_code;comment:应用编码"`
	MenuName  string    `json:"menu_name" form:"menu_name" gorm:"size:128;default:'';column:menu_name;comment:菜单名称"`
	MenuCode  string    `json:"menu_code" form:"menu_code" gorm:"size:128;uniqueIndex;default:'';column:menu_code;comment:菜单编码"`
	Status    int       `json:"status" form:"status" gorm:"not null;default:1;column:status;comment:状态"`
	DataRange string    `json:"data_range" form:"data_range" gorm:"type:char(1);default:'1';column:data_range;comment:数据权限范围"`
	DataScope string    `json:"data_scope" form:"data_scope" gorm:"type:varchar(2048);default:'';column:data_scope;comment:自定义数据范围(组织ID组合)"`
	IsVisible bool      `json:"is_visible" form:"is_visible" gorm:"not null;default:1;column:is_visible;comment:是否可见"`
	MenuType  string    `json:"menu_type" form:"menu_type" gorm:"type:char(1);default:'';column:menu_type;comment:菜单类型（M目录 C菜单 F按钮 Z资源）"`
	Icon      string    `json:"icon" form:"icon" gorm:"size:128;default:'';column:icon;comment:图标"`
	URL       string    `json:"url" form:"url" gorm:"size:512;default:'#';column:url;comment:请求地址"`
	Method    string    `json:"method" form:"method" gorm:"size:16;default:'';column:method;comment:请求方法"`
	Target    string    `json:"target" form:"target" gorm:"size:64;default:'';column:target;comment:打开方式"`
	Remark    string    `json:"remark" form:"remark" gorm:"size:512;default:'';column:remark;comment:备注"`
	Sort      int       `json:"sort" form:"sort" gorm:"default:0;column:sort;comment:序号"`
	Tier      int       `json:"tier" form:"tier" gorm:"default:0;column:tier;comment:层级"`
	IsDeleted bool      `json:"is_deleted" form:"is_deleted" gorm:"type:tinyint(1);default:0;column:is_deleted;comment:是否删除"`
	IsSystem  bool      `json:"is_system" form:"is_system" gorm:"type:tinyint(1);default:0;comment:是否系统内置"`
	TenantID  string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'1';column:tenant_id;comment:租户ID"`
	CreateID  string    `json:"create_id" form:"create_id" gorm:"size:64;default:'';column:create_id;comment:创建人ID"`
	CreateBy  string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';column:create_by;comment:创建人"`
	CreateAt  time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID  string    `json:"update_id" form:"update_id" gorm:"size:64;default:'';column:update_id;comment:更新人ID"`
	UpdateBy  string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';column:update_by;comment:更新人"`
	UpdateAt  time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoMenu) TableName() string {
	return "sso_menu"
}
