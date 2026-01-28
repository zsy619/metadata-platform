package model

import "time"

// SysOperationLog 系统操作日志
type SysOperationLog struct {
	TraceID  string `json:"trace_id" gorm:"primary_key;type:varchar(64);comment:追踪ID"`
	UserID   string `json:"user_id" gorm:"type:varchar(64);index;comment:用户ID"`
	TenantID string `json:"tenant_id" gorm:"type:varchar(64);index;comment:租户ID"`

	// 请求信息
	Method  string `json:"method" gorm:"type:varchar(16);comment:请求方法"`
	Path    string `json:"path" gorm:"type:varchar(256);comment:请求路径"`
	Status  int    `json:"status" gorm:"type:int;comment:响应状态码"`
	Latency int64  `json:"latency" gorm:"type:bigint;comment:延迟(毫秒)"` // milliseconds

	// 客户端信息 (同步 LoginLog 字段结构，便于统一分析)
	ClientIP       string `json:"client_ip" gorm:"type:varchar(64);comment:客户端IP"`
	UserAgent      string `json:"user_agent" gorm:"type:text;comment:UA字符串"`
	Browser        string `json:"browser" gorm:"type:varchar(128);comment:浏览器"`
	BrowserVersion string `json:"browser_version" gorm:"type:varchar(128);comment:浏览器版本"`
	OS             string `json:"os" gorm:"type:varchar(128);comment:操作系统"`
	OSVersion      string `json:"os_version" gorm:"type:varchar(128);comment:操作系统版本"`
	DeviceType     string `json:"device_type" gorm:"type:varchar(64);comment:设备类型"`
	Language       string `json:"language" gorm:"type:varchar(64);comment:语言"`
	Platform       string `json:"platform" gorm:"type:varchar(64);comment:平台"`

	// 业务信息
	ErrorMessage string    `json:"error_message" gorm:"type:text;comment:错误消息"`
	Source       string    `json:"source" gorm:"type:varchar(64);comment:来源模块"` // 来源模块
	Remark       string    `json:"remark" gorm:"size:1024;default:'';comment:备注"`
	CreateAt     time.Time `json:"create_at" gorm:"index;autoCreateTime;comment:创建时间"`
}

// TableName 指定表名
func (SysOperationLog) TableName() string {
	return "sys_operation_log"
}
