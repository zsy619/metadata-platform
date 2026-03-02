package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"metadata-platform/internal/module/document/model"
)

// DocumentRepository 文档数据访问接口
type DocumentRepository interface {
	// 文档基本操作
	Create(doc *model.Document) error
	Update(doc *model.Document) error
	Delete(id string) error
	GetByID(id string) (*model.Document, error)
	GetByPath(path string) (*model.Document, error)
	GetList(page, pageSize int, category, keyword string) ([]*model.Document, int64, error)

	// 分类管理
	GetCategories() ([]*model.DocumentCategory, error)

	// 搜索
	Search(keyword string, category string, limit int) ([]*model.Document, error)

	// 下载
	GetContent(id string) (string, error)

	// 统计
	IncrementViewCount(id string) error

	// 收藏管理
	AddFavorite(documentID, userID string) error
	RemoveFavorite(documentID, userID string) error
	IsFavorited(documentID, userID string) (bool, error)
	GetFavorites(userID string) ([]*model.Document, error)

	// 阅读进度
	UpdateReadProgress(documentID, userID string, position int) error
	GetReadProgress(documentID, userID string) (*model.DocumentReadProgress, error)

	// 版本管理
	CreateVersion(version *model.DocumentVersion) error
	GetVersions(documentID string) ([]*model.DocumentVersion, error)
}

// documentRepository 文档数据访问实现
type documentRepository struct {
	db *gorm.DB
}

// NewDocumentRepository 创建文档数据访问实例
func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{
		db: db,
	}
}

// Create 创建文档
func (r *documentRepository) Create(doc *model.Document) error {
	return r.db.Create(doc).Error
}

// Update 更新文档
func (r *documentRepository) Update(doc *model.Document) error {
	return r.db.Save(doc).Error
}

// Delete 删除文档
func (r *documentRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除文档
		if err := tx.Where("id = ?", id).Delete(&model.Document{}).Error; err != nil {
			return err
		}
		// 删除相关版本
		if err := tx.Where("document_id = ?", id).Delete(&model.DocumentVersion{}).Error; err != nil {
			return err
		}
		// 删除收藏
		if err := tx.Where("document_id = ?", id).Delete(&model.DocumentFavorite{}).Error; err != nil {
			return err
		}
		// 删除阅读进度
		if err := tx.Where("document_id = ?", id).Delete(&model.DocumentReadProgress{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetByID 根据 ID 获取文档
func (r *documentRepository) GetByID(id string) (*model.Document, error) {
	var doc model.Document
	err := r.db.Where("id = ?", id).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// GetByPath 根据路径获取文档
func (r *documentRepository) GetByPath(path string) (*model.Document, error) {
	var doc model.Document
	err := r.db.Where("path = ?", path).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// GetList 获取文档列表（支持分页、筛选、搜索）
func (r *documentRepository) GetList(page, pageSize int, category, keyword string) ([]*model.Document, int64, error) {
	var docs []*model.Document
	var total int64

	query := r.db.Model(&model.Document{})

	// 分类筛选
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 关键词搜索
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR description LIKE ? OR content LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("updated_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&docs).Error

	return docs, total, err
}

// GetCategories 获取所有分类
func (r *documentRepository) GetCategories() ([]*model.DocumentCategory, error) {
	var categories []*model.DocumentCategory
	err := r.db.Order("sort_order ASC, name ASC").Find(&categories).Error

	// 统计每个分类的文档数量（使用单独的变量接收）
	for _, cat := range categories {
		var docCount int64
		r.db.Model(&model.Document{}).Where("category = ? AND is_published = ?", cat.Name, true).Count(&docCount)
		// 注意：DocumentCategory 模型没有 Count 字段，这里不赋值
		_ = docCount
	}

	return categories, err
}

// Search 搜索文档
func (r *documentRepository) Search(keyword string, category string, limit int) ([]*model.Document, error) {
	var docs []*model.Document

	query := r.db.Model(&model.Document{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	searchPattern := "%" + keyword + "%"
	err := query.Where("title LIKE ? OR description LIKE ? OR content LIKE ?",
		searchPattern, searchPattern, searchPattern).
		Order("updated_at DESC").
		Limit(limit).
		Find(&docs).Error

	return docs, err
}

// GetContent 获取文档内容
func (r *documentRepository) GetContent(id string) (string, error) {
	var doc model.Document
	err := r.db.Select("content").Where("id = ?", id).First(&doc).Error
	return doc.Content, err
}

// IncrementViewCount 增加阅读次数
func (r *documentRepository) IncrementViewCount(id string) error {
	return r.db.Exec("UPDATE sys_document SET view_count = view_count + 1 WHERE id = ?", id).Error
}

// AddFavorite 添加收藏
func (r *documentRepository) AddFavorite(documentID, userID string) error {
	favorite := &model.DocumentFavorite{
		ID:         uuid.New().String(),
		DocumentID: documentID,
		UserID:     userID,
	}
	return r.db.Create(favorite).Error
}

// RemoveFavorite 取消收藏
func (r *documentRepository) RemoveFavorite(documentID, userID string) error {
	return r.db.Where("document_id = ? AND user_id = ?", documentID, userID).Delete(&model.DocumentFavorite{}).Error
}

// IsFavorited 检查是否已收藏
func (r *documentRepository) IsFavorited(documentID, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.DocumentFavorite{}).
		Where("document_id = ? AND user_id = ?", documentID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetFavorites 获取用户的收藏列表
func (r *documentRepository) GetFavorites(userID string) ([]*model.Document, error) {
	var docs []*model.Document
	err := r.db.Table("sys_document_favorite").
		Select("sys_document.*").
		Joins("INNER JOIN sys_document ON sys_document.id = sys_document_favorite.document_id").
		Where("sys_document_favorite.user_id = ?", userID).
		Find(&docs).Error
	return docs, err
}

// UpdateReadProgress 更新阅读进度
func (r *documentRepository) UpdateReadProgress(documentID, userID string, position int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var progress model.DocumentReadProgress
		err := tx.Where("document_id = ? AND user_id = ?", documentID, userID).First(&progress).Error

		if err == gorm.ErrRecordNotFound {
			// 创建新记录
			progress = model.DocumentReadProgress{
				ID:           uuid.New().String(),
				DocumentID:   documentID,
				UserID:       userID,
				ReadPosition: position,
			}
			return tx.Create(&progress).Error
		} else if err != nil {
			return err
		}

		// 更新记录
		progress.ReadPosition = position
		progress.LastReadAt = time.Now()
		return tx.Save(&progress).Error
	})
}

// GetReadProgress 获取阅读进度
func (r *documentRepository) GetReadProgress(documentID, userID string) (*model.DocumentReadProgress, error) {
	var progress model.DocumentReadProgress
	err := r.db.Where("document_id = ? AND user_id = ?", documentID, userID).First(&progress).Error
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

// CreateVersion 创建文档版本
func (r *documentRepository) CreateVersion(version *model.DocumentVersion) error {
	return r.db.Create(version).Error
}

// GetVersions 获取文档版本历史
func (r *documentRepository) GetVersions(documentID string) ([]*model.DocumentVersion, error) {
	var versions []*model.DocumentVersion
	err := r.db.Where("document_id = ?", documentID).
		Order("version DESC").
		Find(&versions).Error
	return versions, err
}

// Count 辅助方法：统计文档数量
func (r *documentRepository) Count(category string) (int64, error) {
	var count int64
	query := r.db.Model(&model.Document{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Count(&count).Error
	return count, err
}
