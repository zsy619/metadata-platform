package model

import (
	"time"
)

// Document 文档模型
type Document struct {
	ID          string    `gorm:"type:varchar(64);primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null;index:idx_title_category" json:"title"`   // 文档标题
	Category    string    `gorm:"type:varchar(64);not null;index:idx_title_category" json:"category"` // 文档分类
	Path        string    `gorm:"type:varchar(512);not null;uniqueIndex" json:"path"`                 // 文档路径
	Description string    `gorm:"type:text" json:"description"`                                       // 文档描述
	Content     string    `gorm:"type:text;not null" json:"content"`                                  // 文档内容（Markdown）
	Size        int64     `gorm:"not null" json:"size"`                                               // 文档大小（字节）
	Tags        string    `gorm:"type:varchar(512)" json:"tags,omitempty"`                            // 标签（JSON 数组）
	TOC         string    `gorm:"type:text" json:"toc,omitempty"`                                     // 目录结构（JSON）
	CreatedAt   time.Time `gorm:"autoCreateTime;index" json:"createdAt"`                              // 创建时间
	UpdatedAt   time.Time `gorm:"autoUpdateTime;index" json:"updatedAt"`                              // 更新时间
	CreatedBy   string    `gorm:"type:varchar(64)" json:"createdBy,omitempty"`                        // 创建人
	UpdatedBy   string    `gorm:"type:varchar(64)" json:"updatedBy,omitempty"`                        // 更新人
	Version     int       `gorm:"default:1" json:"version"`                                           // 版本号
	IsPublished bool      `gorm:"default:true" json:"isPublished"`                                    // 是否已发布
	ViewCount   int64     `gorm:"default:0" json:"viewCount"`                                         // 阅读次数
}

// TableName 指定表名
func (Document) TableName() string {
	return "sys_document"
}
