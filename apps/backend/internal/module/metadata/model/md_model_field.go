package model

import "time"

// MdModelField 模型字段模型
type MdModelField struct {
	ID              string    `json:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID        string    `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';comment:租户ID"`
	ModelID         string    `json:"model_id" gorm:"index;type:varchar(64);not null;default:'0';comment:模型ID"`
	TableSchema     string    `json:"table_schema" gorm:"size:64;default:'';comment:表所在库"`
	TableID         string    `json:"table_id" gorm:"type:varchar(64);not null;default:'0';comment:表ID"`
	TableNameStr    string    `json:"table_name" gorm:"column:table_name;size:256;not null;default:'';comment:表名"`
	TableTitle      string    `json:"table_title" gorm:"size:256;default:'';comment:表标题"`
	ColumnID        string    `json:"column_id" gorm:"type:varchar(64);not null;default:'0';comment:列ID"`
	ColumnName      string    `json:"column_name" gorm:"size:256;default:'';comment:列名"`
	ColumnTitle     string    `json:"column_title" gorm:"size:256;default:'';comment:列标题"`
	Func            string    `json:"func" gorm:"size:256;not null;default:''"`
	AggFunc         string    `json:"agg_func" gorm:"size:64;not null;default:''"`
	ColumnType      string    `json:"column_type" gorm:"size:64;default:''"`
	ColumnLength    int       `json:"column_length" gorm:"default:0"`
	IsNullable      bool      `json:"is_nullable" gorm:"not null;default:true"`
	IsPrimaryKey    bool      `json:"is_primary_key" gorm:"not null;default:false"`
	IsAutoIncrement bool      `json:"is_auto_increment" gorm:"not null;default:false"`
	DefaultValue    string    `json:"default_value" gorm:"size:256;default:''"`
	FieldType       string    `json:"field_type" gorm:"size:64;default:''"`       // 逻辑类型 (string, number, date, etc.)
	Max             float64   `json:"max" gorm:"type:decimal(20,4);default:0"`    // 最大值 (数值)
	Min             float64   `json:"min" gorm:"type:decimal(20,4);default:0"`    // 最小值 (数值)
	MaxLength       int       `json:"max_length" gorm:"default:0"`                // 最大长度 (字符串)
	ValidationRule  string    `json:"validation_rule" gorm:"size:512;default:''"` // 正则校验
	ValidationExpr  string    `json:"validation_expr" gorm:"size:512;default:''"` // 表达式校验
	ShowTitle       string    `json:"show_title" gorm:"size:128;not null;default:''"`
	ShowWidth       int       `json:"show_width" gorm:"not null;default:100"`
	IsDeleted       bool      `json:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID        string    `json:"create_id" gorm:"size:64;default:'0';comment:创建人ID"`
	CreateBy        string    `json:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt        time.Time `json:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID        string    `json:"update_id" gorm:"size:64;default:'0';comment:更新人ID"`
	UpdateBy        string    `json:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt        time.Time `json:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdModelField) TableName() string {
	return "md_model_field"
}
