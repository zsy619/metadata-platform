package model

import "time"

// MdTableField 数据连接表字段模型
type MdTableField struct {
	ID              string     `json:"id" gorm:"primary_key;type:varchar(64)"`
	TenantID        string     `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	ConnID          string     `json:"conn_id" gorm:"type:varchar(64);not null;default:'0'"`
	TableID         string     `json:"table_id" gorm:"type:varchar(64);not null;default:'0'"`
	TableTitle      string     `json:"table_title" gorm:"size:256;default:''"`
	ColumnName      string     `json:"column_name" gorm:"size:256;default:''"`
	ColumnTitle     string     `json:"column_title" gorm:"size:256;default:''"`
	ColumnType      string     `json:"column_type" gorm:"size:64;default:''"`
	ColumnLength    int        `json:"column_length" gorm:"default:0"`
	ColumnComment   string     `json:"column_comment" gorm:"size:256;default:''"`
	IsNullable      bool       `json:"is_nullable" gorm:"not null;default:false"`
	IsPrimaryKey    bool       `json:"is_primary_key" gorm:"not null;default:false"`
	IsAutoIncrement bool       `json:"is_auto_increment" gorm:"not null;default:false"`
	DefaultValue    string     `json:"default_value" gorm:"size:256;default:''"`
	ExtraInfo       string     `json:"extra_info" gorm:"size:1024;default:''"`
	IsDeleted       bool       `json:"is_deleted" gorm:"default:false"`
	CreateID        string     `json:"create_id" gorm:"size:64;default:'0'"`
	CreateBy        string     `json:"create_by" gorm:"size:64;default:''"`
	CreateAt        time.Time  `json:"create_at" gorm:"autoCreateTime"`
	UpdateID        string     `json:"update_id" gorm:"size:64;default:'0'"`
	UpdateBy        string     `json:"update_by" gorm:"size:64;default:''"`
	UpdateAt        time.Time  `json:"update_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (MdTableField) TableName() string {
	return "md_table_field"
}
