package model

import "time"

// MdTableField 数据连接表字段模型
type MdTableField struct {
	ID              string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID        string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';comment:租户ID"`
	ConnID          string    `json:"conn_id" form:"conn_id" gorm:"type:varchar(64);not null;default:'0';comment:连接ID"`
	TableID         string    `json:"table_id" form:"table_id" gorm:"type:varchar(64);not null;default:'0';comment:表ID"`
	TableNameStr    string    `json:"table_name" form:"table_name" gorm:"column:table_name;size:256;default:'';comment:物理表名"`
	TableTitle      string    `json:"table_title" form:"table_title" gorm:"size:256;default:'';comment:表标题"`
	ColumnName      string    `json:"column_name" form:"column_name" gorm:"size:256;default:'';comment:列名"`
	ColumnTitle     string    `json:"column_title" form:"column_title" gorm:"size:256;default:'';comment:列标题"`
	ColumnType      string    `json:"column_type" form:"column_type" gorm:"size:64;default:'';comment:数据类型，例如INT、VARCHAR(255)、TIMESTAMP等"`
	ColumnLength    int       `json:"column_length" form:"column_length" gorm:"default:0;comment:字段长度"`
	ColumnComment   string    `json:"column_comment" form:"column_comment" gorm:"size:256;default:'';comment:字段描述"`
	IsNullable      bool      `json:"is_nullable" form:"is_nullable" gorm:"not null;default:false;comment:是否允许为空"`
	IsPrimaryKey    bool      `json:"is_primary_key" form:"is_primary_key" gorm:"not null;default:false;comment:是否为主键"`
	IsAutoIncrement bool      `json:"is_auto_increment" form:"is_auto_increment" gorm:"not null;default:false;comment:是否自增"`
	DefaultValue    string    `json:"default_value" form:"default_value" gorm:"size:256;default:'';comment:默认值"`
	ExtraInfo       string    `json:"extra_info" form:"extra_info" gorm:"size:1024;default:'';comment:额外信息（如auto_increment, unique等）"`
	Remark          string    `json:"remark" form:"remark" gorm:"size:1024;default:'';comment:备注"`
	Sort            int       `json:"sort" form:"sort" gorm:"default:0;comment:排序"`
	IsDeleted       bool      `json:"is_deleted" form:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID        string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';comment:创建人ID"`
	CreateBy        string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt        time.Time `json:"create_at" form:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID        string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';comment:更新人ID"`
	UpdateBy        string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt        time.Time `json:"update_at" form:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdTableField) TableName() string {
	return "md_table_field"
}
