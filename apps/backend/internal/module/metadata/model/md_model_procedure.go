package model

import "time"

// MdModelProcedure 存储过程/函数模型
type MdModelProcedure struct {
	ID          string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID    string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'';comment:租户ID"`
	ConnID      string    `json:"conn_id" form:"conn_id" gorm:"type:varchar(64);not null;default:'';comment:连接ID"`
	ConnName    string    `json:"conn_name" form:"conn_name" gorm:"size:256;default:'';comment:连接名称"`
	ProcSchema  string    `json:"proc_schema" form:"proc_schema" gorm:"size:64;default:'';comment:模式/命名空间"`
	ProcName    string    `json:"proc_name" form:"proc_name" gorm:"size:256;default:'';comment:存储过程/函数名称"`
	ProcTitle   string    `json:"proc_title" form:"proc_title" gorm:"size:256;default:'';comment:存储过程/函数标题"`
	ProcType    string    `json:"proc_type" form:"proc_type" gorm:"size:64;default:'';comment:类型: PROCEDURE 或 FUNCTION"`
	ProcComment string    `json:"proc_comment" form:"proc_comment" gorm:"size:1024;default:'';comment:存储过程/函数描述"`
	Definition  string    `json:"definition" form:"definition" gorm:"type:text;comment:完整定义/代码"`
	ReturnType  string    `json:"return_type" form:"return_type" gorm:"size:256;default:'';comment:返回类型(仅函数有)"`
	Language    string    `json:"language" form:"language" gorm:"size:64;default:'';comment:语言(如SQL, PL/pgSQL等)"`
	Sort        int       `json:"sort" form:"sort" gorm:"default:0;comment:排序"`
	IsDeleted   bool      `json:"is_deleted" form:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID    string    `json:"create_id" form:"create_id" gorm:"size:64;default:'';comment:创建人ID"`
	CreateBy    string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt    time.Time `json:"create_at" form:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID    string    `json:"update_id" form:"update_id" gorm:"size:64;default:'';comment:更新人ID"`
	UpdateBy    string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt    time.Time `json:"update_at" form:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdModelProcedure) TableName() string {
	return "md_model_procedure"
}
