package model

import "time"

// MdModelWhere 模型-where条件模型
type MdModelWhere struct {
	ID                string    `json:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID          string    `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';comment:租户ID"`
	ModelID           string    `json:"model_id" gorm:"index;type:varchar(64);not null;default:'0';comment:模型ID"`
	Operator1         string    `json:"operator1" gorm:"size:64;not null;default:''"`
	Brackets1         string    `json:"brackets1" gorm:"size:64;not null;default:''"`
	TableSchema       string    `json:"table_schema" gorm:"size:64;default:''"`
	TableID           string    `json:"table_id" gorm:"type:varchar(64);not null;default:'0'"`
	TableNameStr      string    `json:"table_name" gorm:"column:table_name;size:256;not null;default:''"`
	TableTitle        string    `json:"table_title" gorm:"size:256;default:''"`
	ColumnID          string    `json:"column_id" gorm:"type:varchar(64);not null;default:'0'"`
	ColumnName        string    `json:"column_name" gorm:"size:256;default:''"`
	ColumnTitle       string    `json:"column_title" gorm:"size:256;default:''"`
	Func              string    `json:"func" gorm:"size:256;not null;default:''"`
	Operator2         string    `json:"operator2" gorm:"size:64;not null;default:''"`
	WhereTableSchema  string    `json:"where_table_schema" gorm:"size:64;default:''"`
	WhereTableID      string    `json:"where_table_id" gorm:"type:varchar(64);not null;default:'0'"`
	WhereTableNameStr string    `json:"where_table_name" gorm:"column:where_table_name;size:256;not null;default:''"`
	WhereTableTitle   string    `json:"where_table_title" gorm:"size:256;default:''"`
	WhereColumnID     string    `json:"where_column_id" gorm:"type:varchar(64);not null;default:'0'"`
	WhereColumnName   string    `json:"where_column_name" gorm:"size:256;not null;default:''"`
	WhereColumnTitle  string    `json:"where_column_title" gorm:"size:256;not null;default:''"`
	WhereFunc         string    `json:"where_func" gorm:"size:64;not null;default:''"`
	Value1            string    `json:"value1" gorm:"size:128;not null;default:''"`
	Value2            string    `json:"value2" gorm:"size:128;not null;default:''"`
	ParamKey          string    `json:"param_key" gorm:"size:128;not null;default:''"`
	Brackets2         string    `json:"brackets2" gorm:"size:64;not null;default:''"`
	Remark            string    `json:"remark" gorm:"size:1024;default:'';comment:备注"`
	IsDeleted         bool      `json:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID          string    `json:"create_id" gorm:"size:64;default:'0';comment:创建人ID"`
	CreateBy          string    `json:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt          time.Time `json:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID          string    `json:"update_id" gorm:"size:64;default:'0';comment:更新人ID"`
	UpdateBy          string    `json:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt          time.Time `json:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdModelWhere) TableName() string {
	return "md_model_where"
}
