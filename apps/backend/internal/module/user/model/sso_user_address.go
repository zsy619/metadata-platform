package model

import "time"

// SsoUserAddress 用户地址簿模型（与 sso_user 一对多）
type SsoUserAddress struct {
	ID           string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);column:id;comment:主键ID"`
	UserID       string    `json:"user_id" form:"user_id" gorm:"type:varchar(64);not null;index;column:user_id;comment:用户ID"`
	Label        string    `json:"label" form:"label" gorm:"size:32;default:'home';column:label;comment:标签(home/office/other)"`
	ReceiverName string    `json:"receiver_name" form:"receiver_name" gorm:"size:64;default:'';column:receiver_name;comment:收件人姓名"`
	Phone        string    `json:"phone" form:"phone" gorm:"size:32;default:'';column:phone;comment:联系电话"`
	Province     string    `json:"province" form:"province" gorm:"size:64;default:'';column:province;comment:省份"`
	City         string    `json:"city" form:"city" gorm:"size:64;default:'';column:city;comment:城市"`
	District     string    `json:"district" form:"district" gorm:"size:64;default:'';column:district;comment:区县"`
	Detail       string    `json:"detail" form:"detail" gorm:"size:512;default:'';column:detail;comment:详细地址"`
	IsDefault    bool      `json:"is_default" form:"is_default" gorm:";default:0;column:is_default;comment:是否默认地址"`
	Remark       string    `json:"remark" form:"remark" gorm:"size:512;default:'';column:remark;comment:备注"`
	TenantID     string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'1';column:tenant_id;comment:租户ID"`
	CreateID     string    `json:"create_id" form:"create_id" gorm:"size:64;default:'';column:create_id;comment:创建人ID"`
	CreateBy     string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';column:create_by;comment:创建人"`
	CreateAt     time.Time `json:"create_at" form:"create_at" gorm:"column:create_at;autoCreateTime;comment:创建时间"`
	UpdateID     string    `json:"update_id" form:"update_id" gorm:"size:64;default:'';column:update_id;comment:更新人ID"`
	UpdateBy     string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';column:update_by;comment:更新人"`
	UpdateAt     time.Time `json:"update_at" form:"update_at" gorm:"column:update_at;autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SsoUserAddress) TableName() string {
	return "sso_user_address"
}
