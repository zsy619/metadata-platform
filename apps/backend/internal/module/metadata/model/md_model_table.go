package model

import "time"

// MdModelTable 模型-表模型
type MdModelTable struct {
	ID           string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID     string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';comment:租户ID"`
	ModelID      string    `json:"model_id" form:"model_id" gorm:"index;type:varchar(64);not null;default:'0';comment:模型ID"`
	ConnID       string    `json:"conn_id" form:"conn_id" gorm:"type:varchar(64);not null;default:'0';comment:连接ID"`
	TableSchema  string    `json:"table_schema" form:"table_schema" gorm:"size:64;default:'';comment:表模式"`
	TableID      string    `json:"table_id" form:"table_id" gorm:"type:varchar(64);not null;default:'0';comment:表ID"`
	TableNameStr string    `json:"table_name" form:"table_name" gorm:"column:table_name;size:256;not null;default:'';comment:表名称"`
	TableTitle   string    `json:"table_title" form:"table_title" gorm:"size:256;default:'';comment:表标题"`
	IsMain          bool      `json:"is_main" form:"is_main" gorm:"not null;default:false;comment:是否主表"`
	Remark          string    `json:"remark" form:"remark" gorm:"size:1024;default:'';comment:备注"`
	IsDeleted    bool      `json:"is_deleted" form:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID     string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';comment:创建人ID"`
	CreateBy     string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt     time.Time `json:"create_at" form:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID     string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';comment:更新人ID"`
	UpdateBy     string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt     time.Time `json:"update_at" form:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdModelTable) TableName() string {
	return "md_model_table"
}
