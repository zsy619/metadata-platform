package model

import "time"

// API 元数据API模型
type API struct {
	ID        string     `json:"id" gorm:"primary_key;type:varchar(64)"`
	TenantID  string     `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	Name      string     `json:"name" gorm:"size:128;not null"`
	Code      string     `json:"code" gorm:"size:128;not null;uniqueIndex"`
	Path      string     `json:"path" gorm:"size:512;not null"`
	Method    string     `json:"method" gorm:"size:16;not null"`
	IsPublic  bool       `json:"is_public" gorm:"default:false"`
	IsDeleted bool       `json:"is_deleted" gorm:"default:false"`
	State     int        `json:"state" gorm:"default:1"`
	Remark    string     `json:"remark" gorm:"size:512"`
	Sort      int        `json:"sort" gorm:"default:0"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}