package model

import "time"

// MdConn 数据连接模型
type MdConn struct {
	ID           string    `json:"id" form:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID     string    `json:"tenant_id" form:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';comment:租户ID"`
	ParentID     string    `json:"parent_id" form:"parent_id" gorm:"type:varchar(64);not null;default:'0';comment:父ID"`
	ConnName     string    `json:"conn_name" form:"conn_name" gorm:"size:256;default:'';comment:连接名称"`
	ConnKind     string    `json:"conn_kind" form:"conn_kind" gorm:"size:64;default:'';comment:数据连接类型（例如MySQL, Oracle, SQLServer, DB2, DM, KingbaseES）"`
	ConnVersion  string    `json:"conn_version" form:"conn_version" gorm:"size:64;default:'';comment:数据库版本（例如8.0, 12c, 2019）"`
	ConnHost     string    `json:"conn_host" form:"conn_host" gorm:"size:128;default:'';comment:数据连接主机地址"`
	ConnPort     int       `json:"conn_port" form:"conn_port" gorm:"not null;default:0;comment:数据连接端口号"`
	ConnUser     string    `json:"conn_user" form:"conn_user" gorm:"size:128;default:'';comment:用户名"`
	ConnPassword string    `json:"conn_password" form:"conn_password" gorm:"size:128;default:'';comment:密码"`
	ConnDatabase string    `json:"conn_database" form:"conn_database" gorm:"size:128;default:'';comment:数据库"`
	ConnConn     string    `json:"conn_conn" form:"conn_conn" gorm:"size:1024;default:'';comment:链接地址：自动生成"`
	State        int       `json:"state" form:"state" gorm:"not null;default:0;comment:连接状态: 0=未检测, 1=有效, 2=连接失败"`
	Remark       string    `json:"remark" form:"remark" gorm:"size:512;default:'';comment:备注"`
	IsDeleted    bool      `json:"is_deleted" form:"is_deleted" gorm:"default:false;comment:是否删除"`
	CreateID     string    `json:"create_id" form:"create_id" gorm:"size:64;default:'0';comment:创建人ID"`
	CreateBy     string    `json:"create_by" form:"create_by" gorm:"size:64;default:'';comment:创建人"`
	CreateAt     time.Time `json:"create_at" form:"create_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdateID     string    `json:"update_id" form:"update_id" gorm:"size:64;default:'0';comment:更新人ID"`
	UpdateBy     string    `json:"update_by" form:"update_by" gorm:"size:64;default:'';comment:更新人"`
	UpdateAt     time.Time `json:"update_at" form:"update_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (MdConn) TableName() string {
	return "md_conn"
}
