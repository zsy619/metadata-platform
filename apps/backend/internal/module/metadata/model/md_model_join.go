package model

import "time"

// MdModelJoin 模型-关联模型
type MdModelJoin struct {
	ID               string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID         string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';comment:租户ID"`
	ParentID         string    `json:"parent_id" form:"parent_id" gorm:"type:varchar(64);not null;default:'0';comment:父ID"`
	ModelID          string    `json:"model_id" form:"model_id" gorm:"index;type:varchar(64);not null;default:'0';comment:模型ID"`
	JoinType         string    `json:"join_type" form:"join_type" gorm:"size:64;not null;default:'';comment:关联类型：Left Join/Right Join/Inner Join"`
	TableSchema      string    `json:"table_schema" form:"table_schema" gorm:"size:64;default:'';comment:表模式"`
	TableID          string    `json:"table_id" form:"table_id" gorm:"type:varchar(64);not null;default:'0';comment:表ID"`
	TableNameStr     string    `json:"table_name" form:"table_name" gorm:"column:table_name;size:256;not null;default:'';comment:表名称"`
	TableTitle       string    `json:"table_title" form:"table_title" gorm:"size:256;default:'';comment:表标题"`
	JoinTableSchema  string    `json:"join_table_schema" form:"join_table_schema" gorm:"size:64;default:'';comment:表模式"`
	JoinTableID      string    `json:"join_table_id" form:"join_table_id" gorm:"type:varchar(64);not null;default:'0';comment:关联表ID"`
	JoinTableNameStr string    `json:"join_table_name" form:"join_table_name" gorm:"column:join_table_name;size:256;not null;default:'';comment:关联表"`
	JoinTableTitle   string    `json:"join_table_title" form:"join_table_title" gorm:"size:256;default:'';comment:表标题"`
	Remark           string    `json:"remark" form:"remark" gorm:"size:1024;default:'';comment:备注"`
	IsDeleted        bool      `json:"is_deleted" form:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID         string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';comment:创建人ID"`
	CreateBy         string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt         time.Time `json:"create_at" form:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID         string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';comment:更新人ID"`
	UpdateBy         string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt         time.Time `json:"update_at" form:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdModelJoin) TableName() string {
	return "md_model_join"
}

// MdModelJoinField 模型-关联字段
type MdModelJoinField struct {
	ID              string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID        string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';comment:租户ID"`
	JoinID          string    `json:"join_id" form:"join_id" gorm:"type:varchar(64);not null;default:'0';comment:关联ID"`
	Operator1       string    `json:"operator1" form:"operator1" gorm:"size:64;not null;default:'';comment:操作符：and/or"`
	Brackets1       string    `json:"brackets1" form:"brackets1" gorm:"size:64;not null;default:'';comment:括号：("`
	ColumnID        string    `json:"column_id" form:"column_id" gorm:"type:varchar(64);not null;default:'0';comment:字段ID"`
	ColumnName      string    `json:"column_name" form:"column_name" gorm:"size:256;default:'';comment:字段名称"`
	ColumnTitle     string    `json:"column_title" form:"column_title" gorm:"size:256;default:'';comment:字段标题"`
	Func            string    `json:"func" form:"func" gorm:"size:256;not null;default:'';comment:字段函数"`
	Operator2       string    `json:"operator2" form:"operator2" gorm:"size:64;not null;default:'=';comment:运算符：=/</>/<=/>=/like/between"`
	JoinColumnID    string    `json:"join_column_id" form:"join_column_id" gorm:"type:varchar(64);not null;default:'0';comment:关联字段ID"`
	JoinColumnName  string    `json:"join_column_name" form:"join_column_name" gorm:"size:256;not null;default:'';comment:关联字段名称"`
	JoinColumnTitle string    `json:"join_column_title" form:"join_column_title" gorm:"size:256;not null;default:'';comment:关联字段标题"`
	JoinFunc        string    `json:"join_func" form:"join_func" gorm:"size:256;not null;default:'';comment:关联字段函数"`
	Value1          string    `json:"value1" form:"value1" gorm:"size:128;not null;default:'';comment:值1"`
	Value2          string    `json:"value2" form:"value2" gorm:"size:128;not null;default:'';comment:值2"`
	Brackets2       string    `json:"brackets2" form:"brackets2" gorm:"size:64;not null;default:'';comment:括号：)"`
	Order           int       `json:"order" form:"order" gorm:"not null;default:0;comment:排序"`
	IsDeleted       bool      `json:"is_deleted" form:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID        string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';comment:创建人ID"`
	CreateBy        string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt        time.Time `json:"create_at" form:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID        string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';comment:更新人ID"`
	UpdateBy        string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt        time.Time `json:"update_at" form:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdModelJoinField) TableName() string {
	return "md_model_join_field"
}
