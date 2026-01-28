package model

import "time"

// SysDataChangeLog 系统数据变更日志
type SysDataChangeLog struct {
	ID         string    `json:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TraceID    string    `json:"trace_id" gorm:"type:varchar(64);index;comment:追踪ID"` // Link to Operation Log
	ModelID    string    `json:"model_id" gorm:"type:varchar(64);index;comment:模型ID"`
	RecordID   string    `json:"record_id" gorm:"type:varchar(64);index;comment:数据ID"`
	Action     string    `json:"action" gorm:"type:varchar(32);comment:操作类型"`   // CREATE, UPDATE, DELETE
	BeforeData string    `json:"before_data" gorm:"type:longtext;comment:变更前数据"` // JSON string
	AfterData  string    `json:"after_data" gorm:"type:longtext;comment:变更后数据"`  // JSON string
	CreateBy   string    `json:"create_by" gorm:"type:varchar(64);comment:创建者"`
	Source     string    `json:"source" gorm:"type:varchar(64);comment:来源模块"` // 来源模块
	Remark       string    `json:"remark" gorm:"size:1024;default:'';comment:备注"`
	CreateAt   time.Time `json:"create_at" gorm:"autoCreateTime;comment:创建时间"`
}

// TableName 指定表名
func (SysDataChangeLog) TableName() string {
	return "sys_data_change_log"
}
