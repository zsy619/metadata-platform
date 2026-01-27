package model

import "time"

// SysOperationLog 系统操作日志
type SysOperationLog struct {
	TraceID  string `json:"trace_id" gorm:"primary_key;type:varchar(64)"`
	UserID   string `json:"user_id" gorm:"type:varchar(64);index"`
	TenantID string `json:"tenant_id" gorm:"type:varchar(64);index"`
	
	// 请求信息
	Method  string `json:"method" gorm:"type:varchar(16)"`
	Path    string `json:"path" gorm:"type:varchar(256)"`
	Status  int    `json:"status" gorm:"type:int"`
	Latency int64  `json:"latency" gorm:"type:bigint"` // milliseconds
	
	// 客户端信息 (同步 LoginLog 字段结构，便于统一分析)
	ClientIP       string `json:"client_ip" gorm:"type:varchar(64)"`
	UserAgent      string `json:"user_agent" gorm:"type:text"`
	Browser        string `json:"browser" gorm:"type:varchar(128)"`
	BrowserVersion string `json:"browser_version" gorm:"type:varchar(128)"`
	OS             string `json:"os" gorm:"type:varchar(128)"`
	OSVersion      string `json:"os_version" gorm:"type:varchar(128)"`
	DeviceType     string `json:"device_type" gorm:"type:varchar(64)"`

	// 业务信息
	ErrorMessage string    `json:"error_message" gorm:"type:text"`
	Source       string    `json:"source" gorm:"type:varchar(64);comment:来源模块"` // 来源模块
	CreateAt     time.Time `json:"create_at" gorm:"index"`
}

// TableName 指定表名
func (SysOperationLog) TableName() string {
	return "sys_operation_log"
}
