package model

import "time"

// MdQueryCondition 查询模板条件模型
type MdQueryCondition struct {
	ID           string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID     string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'';comment:租户ID"`
	TemplateID   string    `json:"template_id" form:"template_id" gorm:"index;type:varchar(64);not null;default:'';comment:模板ID"`
	Operator1    string    `json:"operator1" form:"operator1" gorm:"size:64;not null;default:''"` // AND/OR
	Brackets1    string    `json:"brackets1" form:"brackets1" gorm:"size:64;not null;default:''"` // (
	TableSchema  string    `json:"table_schema" form:"table_schema" gorm:"size:64;default:''"`
	TableNameStr string    `json:"table_name" form:"table_name" gorm:"column:table_name;size:256;not null;default:''"`
	ColumnName   string    `json:"column_name" form:"column_name" gorm:"size:256;default:''"`
	Func         string    `json:"func" form:"func" gorm:"size:256;not null;default:''"`
	Operator2    string    `json:"operator2" form:"operator2" gorm:"size:64;not null;default:''"` // =, >, <, LIKE, etc.
	Value1       string    `json:"value1" form:"value1" gorm:"size:128;not null;default:''"`
	Value2       string    `json:"value2" form:"value2" gorm:"size:128;not null;default:''"`
	Brackets2    string    `json:"brackets2" form:"brackets2" gorm:"size:64;not null;default:''"` // )
	Sort         int       `json:"sort" form:"sort" gorm:"default:0;comment:排序"`
	Remark       string    `json:"remark" form:"remark" gorm:"size:1024;default:'';comment:备注"`
	IsDeleted    bool      `json:"is_deleted" form:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID     string    `json:"create_id" form:"create_id" gorm:"size:64;default:'';comment:创建人ID"`
	CreateBy     string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt     time.Time `json:"create_at" form:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID     string    `json:"update_id" form:"update_id" gorm:"size:64;default:'';comment:更新人ID"`
	UpdateBy     string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt     time.Time `json:"update_at" form:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdQueryCondition) TableName() string {
	return "md_query_condition"
}
