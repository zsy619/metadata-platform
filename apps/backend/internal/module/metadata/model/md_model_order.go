package model

import "time"

// MdModelOrder 模型-排序字段模型
type MdModelOrder struct {
	ID           string    `json:"id" gorm:"primary_key;type:varchar(64)"`
	TenantID     string    `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	ModelID      string    `json:"model_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	TableSchema  string    `json:"table_schema" gorm:"size:64;default:''"`
	TableID      string    `json:"table_id" gorm:"type:varchar(64);not null;default:'0'"`
	TableNameStr string    `json:"table_name" gorm:"column:table_name;size:256;not null;default:''"`
	TableTitle   string    `json:"table_title" gorm:"size:256;default:''"`
	ColumnID     string    `json:"column_id" gorm:"type:varchar(64);not null;default:'0'"`
	ColumnName   string    `json:"column_name" gorm:"size:256;default:''"`
	ColumnTitle  string    `json:"column_title" gorm:"size:256;default:''"`
	Func         string    `json:"func" gorm:"size:128;not null;default:''"`
	OrderType    string    `json:"order_type" gorm:"size:128;not null;default:'Asc'"`
	IsDeleted    bool      `json:"is_deleted" gorm:"default:false"`
	CreateID     string    `json:"create_id" gorm:"size:64;default:'0'"`
	CreateBy     string    `json:"create_by" gorm:"size:64;default:''"`
	CreateAt     time.Time `json:"create_at"`
	UpdateID     string    `json:"update_id" gorm:"size:64;default:'0'"`
	UpdateBy     string    `json:"update_by" gorm:"size:64;default:''"`
	UpdateAt     time.Time `json:"update_at"`
}

// TableName 指定表名
func (MdModelOrder) TableName() string {
	return "md_model_order"
}
