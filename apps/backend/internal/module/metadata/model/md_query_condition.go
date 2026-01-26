package model

import "time"

// MdQueryCondition 查询模板条件模型
type MdQueryCondition struct {
	ID                string    `json:"id" gorm:"primary_key;type:varchar(64)"`
	TenantID          string    `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	TemplateID        string    `json:"template_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	Operator1         string    `json:"operator1" gorm:"size:64;not null;default:''"` // AND/OR
	Brackets1         string    `json:"brackets1" gorm:"size:64;not null;default:''"` // (
	TableSchema       string    `json:"table_schema" gorm:"size:64;default:''"`
	TableNameStr      string    `json:"table_name" gorm:"column:table_name;size:256;not null;default:''"`
	ColumnName        string    `json:"column_name" gorm:"size:256;default:''"`
	Func              string    `json:"func" gorm:"size:256;not null;default:''"`
	Operator2         string    `json:"operator2" gorm:"size:64;not null;default:''"` // =, >, <, LIKE, etc.
	Value1            string    `json:"value1" gorm:"size:128;not null;default:''"`
	Value2            string    `json:"value2" gorm:"size:128;not null;default:''"`
	Brackets2         string    `json:"brackets2" gorm:"size:64;not null;default:''"` // )
	Sort              int       `json:"sort" gorm:"default:0"`
	IsDeleted         bool      `json:"is_deleted" gorm:"default:false"`
	CreateID          string    `json:"create_id" gorm:"size:64;default:'0'"`
	CreateBy          string    `json:"create_by" gorm:"size:64;default:''"`
	CreateAt          time.Time `json:"create_at"`
	UpdateID          string    `json:"update_id" gorm:"size:64;default:'0'"`
	UpdateBy          string    `json:"update_by" gorm:"size:64;default:''"`
	UpdateAt          time.Time `json:"update_at"`
}

// TableName 指定表名
func (MdQueryCondition) TableName() string {
	return "md_query_condition"
}
