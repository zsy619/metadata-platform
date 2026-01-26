package model

import "time"

// MdModel 模型定义模型
type MdModel struct {
	ID          string     `json:"id" gorm:"primary_key;type:varchar(64)"`
	TenantID    string     `json:"tenant_id" gorm:"index;type:varchar(64);not null;default:'0'"`
	ParentID    string     `json:"parent_id" gorm:"type:varchar(64);not null;default:'0'"`
	ConnID      string     `json:"conn_id" gorm:"type:varchar(64);not null;default:'0'"`
	ConnName    string     `json:"conn_name" gorm:"size:256;default:''"`
	ModelName   string     `json:"model_name" gorm:"size:128;not null;default:''"`
	ModelCode   string     `json:"model_code" gorm:"size:128;not null;default:'';uniqueIndex:uix_md_model_title_creator"`
	ModelVersion string    `json:"model_version" gorm:"size:64;not null;default:'1.0.0'"`
	ModelLogo   string     `json:"model_logo" gorm:"size:512;not null;default:''"`
	ModelKind   int        `json:"model_kind" gorm:"not null;default:0"`
	IsPublic    bool       `json:"is_public" gorm:"not null;default:false"`
	IsLocked    bool       `json:"is_locked" gorm:"default:false"`
	IsDeleted   bool       `json:"is_deleted" gorm:"default:false"`
	CreateID    string     `json:"create_id" gorm:"size:64;default:'0'"`
	CreateBy    string     `json:"create_by" gorm:"size:64;default:'';uniqueIndex:uix_md_model_title_creator"`
	CreateAt    time.Time  `json:"create_at"`
	UpdateID    string     `json:"update_id" gorm:"size:64;default:'0'"`
	UpdateBy    string     `json:"update_by" gorm:"size:64;default:''"`
	UpdateAt    time.Time  `json:"update_at"`
}

// TableName 指定表名
func (MdModel) TableName() string {
	return "md_model"
}
