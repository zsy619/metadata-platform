package model

import "time"

// MdModelLimit 模型-获取条数模型
type MdModelLimit struct {
	ID        string    `json:"id" gorm:"primary_key;type:varchar(64)"`
	TenantID  string    `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	ModelID   string    `json:"model_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	Page      int       `json:"page" gorm:"not null;default:0"`
	Limit     int       `json:"limit" gorm:"not null;default:0"`
	IsDeleted bool      `json:"is_deleted" gorm:"default:false"`
	CreateID  string    `json:"create_id" gorm:"size:64;default:'0'"`
	CreateBy  string    `json:"create_by" gorm:"size:64;default:''"`
	CreateAt  time.Time `json:"create_at"`
	UpdateID  string    `json:"update_id" gorm:"size:64;default:'0'"`
	UpdateBy  string    `json:"update_by" gorm:"size:64;default:''"`
	UpdateAt  time.Time `json:"update_at"`
}

// TableName 指定表名
func (MdModelLimit) TableName() string {
	return "md_model_limit"
}
