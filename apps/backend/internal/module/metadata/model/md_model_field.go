package model

import "time"

// MdModelField 模型字段模型
type MdModelField struct {
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
	Func         string    `json:"func" gorm:"size:256;not null;default:''"`
	AggFunc      string    `json:"agg_func" gorm:"size:64;not null;default:''"`
	ColumnType      string    `json:"column_type" gorm:"size:64;default:''"`
	ColumnLength    int       `json:"column_length" gorm:"default:0"`
	IsNullable      bool      `json:"is_nullable" gorm:"not null;default:true"`
	IsPrimaryKey    bool      `json:"is_primary_key" gorm:"not null;default:false"`
	IsAutoIncrement bool      `json:"is_auto_increment" gorm:"not null;default:false"`
	DefaultValue    string    `json:"default_value" gorm:"size:256;default:''"`
	FieldType       string    `json:"field_type" gorm:"size:64;default:''"`      // 逻辑类型 (string, number, date, etc.)
	Max             float64   `json:"max" gorm:"type:decimal(20,4);default:0"`   // 最大值 (数值)
	Min             float64   `json:"min" gorm:"type:decimal(20,4);default:0"`   // 最小值 (数值)
	MaxLength       int       `json:"max_length" gorm:"default:0"`               // 最大长度 (字符串)
	ValidationRule  string    `json:"validation_rule" gorm:"size:512;default:''"` // 正则校验
	ValidationExpr  string    `json:"validation_expr" gorm:"size:512;default:''"` // 表达式校验
	ShowTitle       string    `json:"show_title" gorm:"size:128;not null;default:''"`
	ShowWidth    int       `json:"show_width" gorm:"not null;default:100"`
	IsDeleted    bool      `json:"is_deleted" gorm:"default:false"`
	CreateID     string    `json:"create_id" gorm:"size:64;default:'0'"`
	CreateBy     string    `json:"create_by" gorm:"size:64;default:''"`
	CreateAt     time.Time `json:"create_at"`
	UpdateID     string    `json:"update_id" gorm:"size:64;default:'0'"`
	UpdateBy     string    `json:"update_by" gorm:"size:64;default:''"`
	UpdateAt     time.Time `json:"update_at"`
}

// TableName 指定表名
func (MdModelField) TableName() string {
	return "md_model_field"
}
