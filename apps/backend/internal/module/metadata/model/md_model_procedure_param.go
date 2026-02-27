package model

import "time"

// MdModelProcedureParam 存储过程/函数参数模型
type MdModelProcedureParam struct {
	ID           string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID     string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'';comment:租户ID"`
	ConnID       string    `json:"conn_id" form:"conn_id" gorm:"type:varchar(64);not null;default:'';comment:连接ID"`
	ProcID       string    `json:"proc_id" form:"proc_id" gorm:"type:varchar(64);not null;default:'';comment:存储过程/函数ID"`
	ProcName     string    `json:"proc_name" form:"proc_name" gorm:"size:256;default:'';comment:存储过程/函数名称"`
	ParamName    string    `json:"param_name" form:"param_name" gorm:"size:256;default:'';comment:参数名称"`
	ParamTitle   string    `json:"param_title" form:"param_title" gorm:"size:256;default:'';comment:参数标题"`
	ParamMode    string    `json:"param_mode" form:"param_mode" gorm:"size:32;default:'';comment:参数模式: IN, OUT, INOUT"`
	ParamType    string    `json:"param_type" form:"param_type" gorm:"size:64;default:'';comment:参数数据类型"`
	ParamLength  int       `json:"param_length" form:"param_length" gorm:"default:0;comment:参数长度"`
	ParamComment string    `json:"param_comment" form:"param_comment" gorm:"size:256;default:'';comment:参数描述"`
	DefaultValue string    `json:"default_value" form:"default_value" gorm:"size:256;default:'';comment:默认值"`
	Sort         int       `json:"sort" form:"sort" gorm:"default:0;comment:排序"`
	IsDeleted    bool      `json:"is_deleted" form:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID     string    `json:"create_id" form:"create_id" gorm:"size:64;default:'';comment:创建人ID"`
	CreateBy     string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt     time.Time `json:"create_at" form:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID     string    `json:"update_id" form:"update_id" gorm:"size:64;default:'';comment:更新人ID"`
	UpdateBy     string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt     time.Time `json:"update_at" form:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdModelProcedureParam) TableName() string {
	return "md_model_procedure_param"
}
