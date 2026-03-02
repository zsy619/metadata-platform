package model

import (
	"time"
)

// DocumentFolder 文档目录/文件夹
type DocumentFolder struct {
	ID          string    `gorm:"type:varchar(64);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`                // 文件夹名称
	Path        string    `gorm:"type:varchar(512);not null;uniqueIndex" json:"path"`    // 完整路径（如：/root/level1/level2）
	ParentID    string    `gorm:"type:varchar(64);default:'';index" json:"parentId"`     // 父文件夹 ID
	Level       int       `gorm:"default:0" json:"level"`                                // 层级（从 0 开始）
	Description string    `gorm:"type:varchar(512)" json:"description"`                  // 描述
	SortOrder   int       `gorm:"default:0" json:"sortOrder"`                          // 排序顺序
	IsEnabled   bool      `gorm:"default:true" json:"isEnabled"`                       // 是否启用
	CreatedBy   string    `gorm:"type:varchar(64)" json:"createdBy"`
	UpdatedBy   string    `gorm:"type:varchar(64)" json:"updatedBy"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	
	// 非数据库字段
	Children   []*DocumentFolder `gorm:"-" json:"children,omitempty"`   // 子文件夹
	DocCount   int64             `gorm:"-" json:"docCount"`             // 该文件夹下的文档数量
	HasChildren bool            `gorm:"-" json:"hasChildren"`          // 是否有子文件夹
}

// TableName 指定表名
func (DocumentFolder) TableName() string {
	return "sys_document_folder"
}

// DocumentFolderTree 文件夹树形结构
type DocumentFolderTree struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Path        string              `json:"path"`
	Level       int                 `json:"level"`
	Children    []*DocumentFolderTree `json:"children,omitempty"`
	DocCount    int64               `json:"docCount"`
	HasChildren bool                `json:"hasChildren"`
}
