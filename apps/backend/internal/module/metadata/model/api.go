package model

import "time"

// API 元数据API模型
type API struct {
	ID        string     `json:"id" gorm:"primary_key;type:varchar(64);comment:主键ID"`
	TenantID  string     `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0';comment:租户ID"`
	Name      string     `json:"name" gorm:"size:128;not null;comment:API名称"`
	Code      string     `json:"code" gorm:"size:128;not null;uniqueIndex;comment:API编码"`
	Path      string     `json:"path" gorm:"size:512;not null;comment:API路径"`
	Method    string     `json:"method" gorm:"size:16;not null;comment:请求方法"`
	IsPublic  bool       `json:"is_public" gorm:"default:false;comment:是否公开"`
	IsDeleted bool       `json:"is_deleted" gorm:"default:false;comment:是否删除"`
	State     int        `json:"state" gorm:"default:1;comment:状态"`
	Remark    string     `json:"remark" gorm:"size:512;comment:备注"`
	Sort      int        `json:"sort" gorm:"default:0;comment:排序"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index;comment:删除时间"`
}