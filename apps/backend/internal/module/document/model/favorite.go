package model

import (
	"time"
)

// DocumentFavorite 文档收藏
type DocumentFavorite struct {
	ID         string    `gorm:"type:varchar(64);primaryKey" json:"id"`
	DocumentID string    `gorm:"type:varchar(64);not null;uniqueIndex:idx_user_doc" json:"documentId"`
	UserID     string    `gorm:"type:varchar(64);not null;uniqueIndex:idx_user_doc" json:"userId"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

// TableName 指定表名
func (DocumentFavorite) TableName() string {
	return "sys_document_favorite"
}
