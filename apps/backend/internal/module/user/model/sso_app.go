package model

import "time"

// SsoApp 应用服务模块模型
type SsoApp struct {
	ID        string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	ParentID  string    `json:"parent_id" form:"parent_id" gorm:"type:varchar(64);not null;default:'';column:parent_id;comment:父ID"`
	AppName   string    `json:"app_name" form:"app_name" gorm:"size:32;default:'';column:app_name;comment:应用名称"`
	AppCode   string    `json:"app_code" form:"app_code" gorm:"size:64;uniqueIndex;default:'';column:app_code;comment:应用编码"`
	Status    int       `json:"status" form:"status" gorm:"not null;default:1;column:status;comment:状态"`
	Host      string    `json:"host" form:"host" gorm:"size:512;default:'';column:host"`
	Logo      string    `json:"logo" form:"logo" gorm:"size:256;default:'';column:logo"`
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
func (SsoApp) TableName() string {
	return "sso_app"
}
