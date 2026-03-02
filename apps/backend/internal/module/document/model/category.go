package model

import (
	"time"
)

// DocumentCategory 文档分类
type DocumentCategory struct {
	ID          string    `gorm:"type:varchar(64);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(64);not null;unique" json:"name"`        // 分类名称
	Description string    `gorm:"type:varchar(255)" json:"description"`                // 分类描述
	Icon        string    `gorm:"type:varchar(64)" json:"icon,omitempty"`              // 分类图标
	SortOrder   int       `gorm:"default:0" json:"sortOrder"`                          // 排序顺序
	ParentID    string    `gorm:"type:varchar(64);default:''" json:"parentId"`         // 父分类 ID
	IsEnabled   bool      `gorm:"default:true" json:"isEnabled"`                       // 是否启用
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	
	// 非数据库字段
	Count int64 `gorm:"-" json:"count"` // 文档数量（仅用于 API 响应）
}

// TableName 指定表名
func (DocumentCategory) TableName() string {
	return "sys_document_category"
}
