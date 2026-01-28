package model

import "time"

// SysLoginLog 系统登录日志
type SysLoginLog struct {
	ID          string `json:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	UserID      string `json:"user_id" gorm:"type:varchar(64);index;comment:用户ID"`
	Account     string `json:"account" gorm:"type:varchar(128);index;comment:账号名"`
	LoginStatus int    `json:"login_status" gorm:"type:tinyint;comment:1成功 0失败 2退出"` // 1:Success, 0:Fail, 2:Logout

	// 浏览器信息
	Browser        string `json:"browser" gorm:"type:varchar(128);index;comment:浏览器名称"`
	BrowserVersion string `json:"browser_version" gorm:"type:varchar(128);comment:浏览器版本"`
	BrowserEngine  string `json:"browser_engine" gorm:"type:varchar(128);comment:浏览器引擎"` // Layout Engine (e.g. Blink, Gecko)
	Language       string `json:"language" gorm:"type:varchar(64);comment:语言"`        // Accepted Language
	UserAgent      string `json:"user_agent" gorm:"type:text;comment:UA字符串"`             // Raw UA

	// 操作系统信息
	OS        string `json:"os" gorm:"type:varchar(128);index;comment:操作系统"`
	OSVersion string `json:"os_version" gorm:"type:varchar(128);comment:操作系统版本"`
	OSArch    string `json:"os_arch" gorm:"type:varchar(64);comment:操作系统架构"` // x86, arm64

	// 设备信息
	DeviceType       string `json:"device_type" gorm:"type:varchar(64);comment:设备类型"` // Desktop, Mobile, Tablet
	DeviceModel      string `json:"device_model" gorm:"type:varchar(128);comment:设备型号"`
	ScreenResolution string `json:"screen_resolution" gorm:"type:varchar(64);comment:屏幕分辨率"` // e.g. 1920x1080

	// 网络与区域
	ClientIP   string `json:"client_ip" gorm:"type:varchar(64);index;comment:客户端IP"`
	IPLocation string `json:"ip_location" gorm:"type:varchar(128);comment:IP归属地"` // Country/City
	Timezone   string `json:"timezone" gorm:"type:varchar(64);comment:时区"`     // e.g. Asia/Shanghai

	Platform     string    `json:"platform" gorm:"type:varchar(64);comment:来源平台"` // Web, App, etc.
	CreateAt     time.Time `json:"create_at" gorm:"index;autoCreateTime;comment:创建时间"`
	ErrorMessage string    `json:"error_message" gorm:"type:text;comment:错误消息"`
	Remark       string    `json:"remark" gorm:"size:1024;default:'';comment:备注"`
}

// TableName 指定表名
func (SysLoginLog) TableName() string {
	return "sys_login_log"
}
