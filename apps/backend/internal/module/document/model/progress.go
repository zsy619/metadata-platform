package model

import (
	"time"
)

// DocumentReadProgress 文档阅读进度
type DocumentReadProgress struct {
	ID           string    `gorm:"type:varchar(64);primaryKey" json:"id"`
	DocumentID   string    `gorm:"type:varchar(64);not null;uniqueIndex:idx_user_doc" json:"documentId"`
	UserID       string    `gorm:"type:varchar(64);not null;uniqueIndex:idx_user_doc" json:"userId"`
	LastReadAt   time.Time `gorm:"autoCreateTime" json:"lastReadAt"`
	ReadPosition int       `gorm:"default:0" json:"readPosition"` // 阅读位置（百分比）
}

// TableName 指定表名
func (DocumentReadProgress) TableName() string {
	return "sys_document_read_progress"
}
