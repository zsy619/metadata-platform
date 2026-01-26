package model

import "time"

// MdModelTable 模型-表模型
type MdModelTable struct {
	ID           string    `json:"id" gorm:"primary_key;type:varchar(64)"`
	TenantID     string    `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	ModelID      string    `json:"model_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	ConnID       string    `json:"conn_id" gorm:"type:varchar(64);not null;default:'0'"`
	TableSchema  string    `json:"table_schema" gorm:"size:64;default:''"`
	TableID      string    `json:"table_id" gorm:"type:varchar(64);not null;default:'0'"`
	TableNameStr string    `json:"table_name" gorm:"column:table_name;size:256;not null;default:''"`
	TableTitle   string    `json:"table_title" gorm:"size:256;default:''"`
	IsMain       bool      `json:"is_main" gorm:"not null;default:false"`
	IsDeleted    bool      `json:"is_deleted" gorm:"default:false"`
	CreateID     string    `json:"create_id" gorm:"size:64;default:'0'"`
	CreateBy     string    `json:"create_by" gorm:"size:64;default:''"`
	CreateAt     time.Time `json:"create_at"`
	UpdateID     string    `json:"update_id" gorm:"size:64;default:'0'"`
	UpdateBy     string    `json:"update_by" gorm:"size:64;default:''"`
	UpdateAt     time.Time `json:"update_at"`
}

// TableName 指定表名
func (MdModelTable) TableName() string {
	return "md_model_table"
}
