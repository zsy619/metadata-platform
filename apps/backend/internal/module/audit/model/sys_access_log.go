package model

import "time"

// SysAccessLog 系统访问日志
type SysAccessLog struct {
	ID          string    `json:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TraceID     string    `json:"trace_id" gorm:"type:varchar(64);index;comment:追踪ID"`
	UserID      string    `json:"user_id" gorm:"type:varchar(64);index;comment:用户ID"`
	TenantID    string    `json:"tenant_id" gorm:"type:varchar(64);index;comment:租户ID"`
	Method      string    `json:"method" gorm:"type:varchar(16);comment:请求方法"`
	Path        string    `json:"path" gorm:"type:varchar(512);comment:请求路径"`
	QueryString string    `json:"query_string" gorm:"type:varchar(1024);comment:查询字符串"`
	Status      int       `json:"status" gorm:"type:int;comment:响应状态码"`
	Latency     int64     `json:"latency" gorm:"type:bigint;comment:延迟(毫秒)"`
	ClientIP    string    `json:"client_ip" gorm:"type:varchar(64);comment:客户端IP"`
	UserAgent   string    `json:"user_agent" gorm:"type:text;comment:UA字符串"`
	Referer     string    `json:"referer" gorm:"type:varchar(512);comment:来源页面"`
	RequestSize int64     `json:"request_size" gorm:"type:bigint;comment:请求大小"`
	ResponseSize int64    `json:"response_size" gorm:"type:bigint;comment:响应大小"`
	Country     string    `json:"country" gorm:"type:varchar(64);comment:国家地区"`
	Province    string    `json:"province" gorm:"type:varchar(64);comment:省份"`
	City        string    `json:"city" gorm:"type:varchar(64);comment:城市"`
	ISP         string    `json:"isp" gorm:"type:varchar(64);comment:运营商"`
	CreateAt    time.Time `json:"create_at" gorm:"index;autoCreateTime;comment:创建时间"`
}

// TableName 指定表名
func (SysAccessLog) TableName() string {
	return "sys_access_log"
}
