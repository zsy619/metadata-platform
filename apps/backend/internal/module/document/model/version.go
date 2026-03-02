package model

import (
	"time"
)

// DocumentVersion 文档版本历史
type DocumentVersion struct {
	ID          string    `gorm:"type:varchar(64);primaryKey" json:"id"`
	DocumentID  string    `gorm:"type:varchar(64);not null;index" json:"documentId"` // 文档 ID
	Version     int       `gorm:"not null;index" json:"version"`                     // 版本号
	Content     string    `gorm:"type:text;not null" json:"content"`                 // 版本内容
	ChangeLog   string    `gorm:"type:text" json:"changeLog"`                        // 变更说明
	Size        int64     `gorm:"not null" json:"size"`                              // 版本大小
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	CreatedBy   string    `gorm:"type:varchar(64)" json:"createdBy"`
}

// TableName 指定表名
func (DocumentVersion) TableName() string {
	return "sys_document_version"
}
