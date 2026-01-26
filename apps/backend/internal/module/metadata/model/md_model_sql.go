package model

import "time"

// MdModelSql 模型-sql模型
type MdModelSql struct {
	ID        string    `json:"id" gorm:"primary_key;type:varchar(64)"`
	TenantID  string    `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	ModelID   string    `json:"model_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	Content   string    `json:"content" gorm:"type:text"`
	IsDeleted bool      `json:"is_deleted" gorm:"default:false"`
	CreateID  string    `json:"create_id" gorm:"size:64;default:'0'"`
	CreateBy  string    `json:"create_by" gorm:"size:64;default:''"`
	CreateAt  time.Time `json:"create_at"`
	UpdateID  string    `json:"update_id" gorm:"size:64;default:'0'"`
	UpdateBy  string    `json:"update_by" gorm:"size:64;default:''"`
	UpdateAt  time.Time `json:"update_at"`
}

// TableName 指定表名
func (MdModelSql) TableName() string {
	return "md_model_sql"
}
