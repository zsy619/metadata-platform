package model

import "time"

// MdConn 数据连接模型
type MdConn struct {
	ID          string     `json:"id" gorm:"primary_key;type:varchar(64)"`
	TenantID    string     `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	ParentID    string     `json:"parent_id" gorm:"type:varchar(64);not null;default:'0'"`
	ConnName    string     `json:"conn_name" gorm:"size:256;default:''"`
	ConnKind    string     `json:"conn_kind" gorm:"size:64;default:''"`
	ConnVersion string     `json:"conn_version" gorm:"size:64;default:''"`
	ConnHost    string     `json:"conn_host" gorm:"size:128;default:''"`
	ConnPort    int        `json:"conn_port" gorm:"not null;default:0"`
	ConnUser    string     `json:"conn_user" gorm:"size:128;default:''"`
	ConnPassword string     `json:"conn_password" gorm:"size:128;default:''"`
	ConnDatabase string     `json:"conn_database" gorm:"size:128;default:''"`
	ConnConn    string     `json:"conn_conn" gorm:"size:1024;default:''"`
	IsDeleted   bool       `json:"is_deleted" gorm:"default:false"`
	CreateID    string     `json:"create_id" gorm:"size:64;default:'0'"`
	CreateBy    string     `json:"create_by" gorm:"size:64;default:''"`
	CreateAt    time.Time  `json:"create_at"`
	UpdateID    string     `json:"update_id" gorm:"size:64;default:'0'"`
	UpdateBy    string     `json:"update_by" gorm:"size:64;default:''"`
	UpdateAt    time.Time  `json:"update_at"`
}

// TableName 指定表名
func (MdConn) TableName() string {
	return "md_conn"
}
