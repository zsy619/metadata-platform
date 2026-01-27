package model

import "time"

// SysDataChangeLog 系统数据变更日志
type SysDataChangeLog struct {
	ID         string    `json:"id" gorm:"primary_key;type:varchar(64)"`
	TraceID    string    `json:"trace_id" gorm:"type:varchar(64);index"` // Link to Operation Log
	ModelID    string    `json:"model_id" gorm:"type:varchar(64);index"`
	RecordID   string    `json:"record_id" gorm:"type:varchar(64);index"`
	Action     string    `json:"action" gorm:"type:varchar(32)"`   // CREATE, UPDATE, DELETE
	BeforeData string    `json:"before_data" gorm:"type:longtext"` // JSON string
	AfterData  string    `json:"after_data" gorm:"type:longtext"`  // JSON string
	CreateBy   string    `json:"create_by" gorm:"type:varchar(64)"`
	Source     string    `json:"source" gorm:"type:varchar(64);comment:来源模块"` // 来源模块
	CreateAt   time.Time `json:"create_at"`
}

// TableName 指定表名
func (SysDataChangeLog) TableName() string {
	return "sys_data_change_log"
}
