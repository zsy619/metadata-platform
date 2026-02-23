package model

import "time"

// MdModelGroup 模型-选择字段模型
type MdModelGroup struct {
	ID           string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID     string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'';comment:租户ID"`
	ModelID      string    `json:"model_id" form:"model_id" gorm:"index;type:varchar(64);not null;default:'';comment:模型ID"`
	TableSchema  string    `json:"table_schema" form:"table_schema" gorm:"size:64;default:'';comment:表模式"`
	TableID      string    `json:"table_id" form:"table_id" gorm:"type:varchar(64);not null;default:'';comment:表ID"`
	TableNameStr string    `json:"table_name" form:"table_name" gorm:"column:table_name;size:256;not null;default:'';comment:表名称"`
	TableTitle   string    `json:"table_title" form:"table_title" gorm:"size:256;default:'';comment:表标题"`
	ColumnID     string    `json:"column_id" form:"column_id" gorm:"type:varchar(64);not null;default:'';comment:字段ID"`
	ColumnName   string    `json:"column_name" form:"column_name" gorm:"size:256;default:'';comment:字段名称"`
	ColumnTitle  string    `json:"column_title" form:"column_title" gorm:"size:256;default:'';comment:字段标题"`
	Func         string    `json:"func" form:"func" gorm:"size:256;not null;default:'';comment:字段函数"`
	AggFunc      string    `json:"agg_func" form:"agg_func" gorm:"size:64;not null;default:'';comment:聚合函数：sum/count/avg/max/min"`
	IsShow       int       `json:"is_show" form:"is_show" gorm:"not null;default:0;comment:展示: 1 有效 0 无效"`
	ShowTitle    string    `json:"show_title" form:"show_title" gorm:"size:128;not null;default:'';comment:字段显示名称"`
	ShowWidth    int       `json:"show_width" form:"show_width" gorm:"not null;default:100;comment:字段显示宽度"`
	IsDeleted    bool      `json:"is_deleted" form:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID     string    `json:"create_id" form:"create_id" gorm:"size:64;default:'';comment:创建人ID"`
	CreateBy     string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt     time.Time `json:"create_at" form:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID     string    `json:"update_id" form:"update_id" gorm:"size:64;default:'';comment:更新人ID"`
	UpdateBy     string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt     time.Time `json:"update_at" form:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdModelGroup) TableName() string {
	return "md_model_group"
}
