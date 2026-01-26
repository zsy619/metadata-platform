package model

import "time"

// SysOperationLog 系统操作日志
type SysOperationLog struct {
	TraceID      string    `json:"trace_id" gorm:"primary_key;type:varchar(64)"`
	UserID       string    `json:"user_id" gorm:"type:varchar(64);index"`
	TenantID     string    `json:"tenant_id" gorm:"type:varchar(64);index"`
	Method       string    `json:"method" gorm:"type:varchar(10)"`
	Path         string    `json:"path" gorm:"type:varchar(255)"`
	Status       int       `json:"status" gorm:"type:int"`
	Latency      int64     `json:"latency" gorm:"type:bigint"` // milliseconds
	ClientIP     string    `json:"client_ip" gorm:"type:varchar(50)"`
	UserAgent    string    `json:"user_agent" gorm:"type:varchar(255)"`
	ErrorMessage string    `json:"error_message" gorm:"type:text"`
	CreateAt     time.Time `json:"create_at" gorm:"index"`
}

// TableName 指定表名
func (SysOperationLog) TableName() string {
	return "sys_operation_log"
}

// SysDataChangeLog 系统数据变更日志
type SysDataChangeLog struct {
	ID          string    `json:"id" gorm:"primary_key;type:varchar(64)"`
	TraceID     string    `json:"trace_id" gorm:"type:varchar(64);index"` // Link to Operation Log
	ModelID     string    `json:"model_id" gorm:"type:varchar(64);index"`
	RecordID    string    `json:"record_id" gorm:"type:varchar(64);index"`
	Action      string    `json:"action" gorm:"type:varchar(20)"` // CREATE, UPDATE, DELETE
	BeforeData  string    `json:"before_data" gorm:"type:longtext"` // JSON string
	AfterData   string    `json:"after_data" gorm:"type:longtext"`  // JSON string
	CreateBy    string    `json:"create_by" gorm:"type:varchar(64)"`
	CreateAt    time.Time `json:"create_at"`
}

// TableName 指定表名
func (SysDataChangeLog) TableName() string {
	return "sys_data_change_log"
}
