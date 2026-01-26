package model

import "time"

// MdModelHaving 模型-having条件模型
type MdModelHaving struct {
	ID                 string    `json:"id" gorm:"primary_key;type:varchar(64)"`
	TenantID           string    `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	ModelID            string    `json:"model_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	Operator1          string    `json:"operator1" gorm:"size:64;not null;default:''"`
	Brackets1          string    `json:"brackets1" gorm:"size:64;not null;default:''"`
	TableSchema        string    `json:"table_schema" gorm:"size:64;default:''"`
	TableID            string    `json:"table_id" gorm:"type:varchar(64);not null;default:'0'"`
	TableNameStr       string    `json:"table_name" gorm:"column:table_name;size:256;not null;default:''"`
	TableTitle         string    `json:"table_title" gorm:"size:256;default:''"`
	ColumnID           string    `json:"column_id" gorm:"type:varchar(64);not null;default:'0'"`
	ColumnName         string    `json:"column_name" gorm:"size:256;default:''"`
	ColumnTitle        string    `json:"column_title" gorm:"size:256;default:''"`
	Func               string    `json:"func" gorm:"size:256;not null;default:''"`
	Operator2          string    `json:"operator2" gorm:"size:64;not null;default:''"`
	HavingTableSchema  string    `json:"having_table_schema" gorm:"size:64;default:''"`
	HavingTableID      string    `json:"having_table_id" gorm:"type:varchar(64);not null;default:'0'"`
	HavingTableNameStr string    `json:"having_table_name" gorm:"column:having_table_name;size:256;not null;default:''"`
	HavingTableTitle   string    `json:"having_table_title" gorm:"size:256;default:''"`
	HavingColumnID     string    `json:"having_column_id" gorm:"type:varchar(64);not null;default:'0'"`
	HavingColumnName   string    `json:"having_column_name" gorm:"size:256;not null;default:''"`
	HavingColumnTitle  string    `json:"having_column_title" gorm:"size:256;not null;default:''"`
	HavingFunc         string    `json:"having_func" gorm:"size:64;not null;default:''"`
	Value1             string    `json:"value1" gorm:"size:128;not null;default:''"`
	Value2             string    `json:"value2" gorm:"size:128;not null;default:''"`
	ParamKey           string    `json:"param_key" gorm:"size:128;not null;default:''"`
	Brackets2          string    `json:"brackets2" gorm:"size:64;not null;default:''"`
	IsDeleted          bool      `json:"is_deleted" gorm:"default:false"`
	CreateID           string    `json:"create_id" gorm:"size:64;default:'0'"`
	CreateBy           string    `json:"create_by" gorm:"size:64;default:''"`
	CreateAt           time.Time `json:"create_at"`
	UpdateID           string    `json:"update_id" gorm:"size:64;default:'0'"`
	UpdateBy           string    `json:"update_by" gorm:"size:64;default:''"`
	UpdateAt           time.Time `json:"update_at"`
}

// TableName 指定表名
func (MdModelHaving) TableName() string {
	return "md_model_having"
}
