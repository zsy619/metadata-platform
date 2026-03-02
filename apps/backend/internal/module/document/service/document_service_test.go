package service

import (
	"context"
	"metadata-platform/internal/module/document/model"
	"metadata-platform/internal/module/document/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 创建测试数据库连接
func setupTestDB() (*gorm.DB, error) {
	// 使用 SQLite 内存数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	// 自动迁移表结构（忽略错误，因为内存数据库每次都是新的）
	_ = db.AutoMigrate(
		&model.Document{},
		&model.DocumentCategory{},
		&model.DocumentVersion{},
		&model.DocumentFavorite{},
		&model.DocumentReadProgress{},
	)

	return db, nil
}

// TestDocumentService_GetDocumentList 测试获取文档列表
func TestDocumentService_GetDocumentList(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := repository.NewDocumentRepository(db)
	svc := NewDocumentService(repo)

	// 创建测试数据
	doc := &model.Document{
		ID:          "test-1",
		Title:       "测试文档",
		Category:    "系统概述",
		Path:        "/test/doc1",
		Description: "测试文档描述",
		Content:     "测试文档内容",
		Size:        1024,
	}
	err = repo.Create(doc)
	assert.NoError(t, err)

	// 测试获取列表
	docs, total, err := svc.GetDocumentList(1, 10, "", "")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, docs, 1)
	assert.Equal(t, "测试文档", docs[0].Title)
}

// TestDocumentService_GetCategories 测试获取分类
func TestDocumentService_GetCategories(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := repository.NewDocumentRepository(db)
	svc := NewDocumentService(repo)

	// 创建测试分类
	category := &model.DocumentCategory{
		ID:          "cat-1",
		Name:        "系统概述",
		Description: "系统概述分类",
		SortOrder:   1,
	}
	err = db.Create(category).Error
	assert.NoError(t, err)

	// 测试获取分类
	categories, err := svc.GetCategories()
	assert.NoError(t, err)
	assert.Len(t, categories, 1)
	assert.Equal(t, "系统概述", categories[0].Name)
}

// TestDocumentService_CreateDocument 测试创建文档
func TestDocumentService_CreateDocument(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := repository.NewDocumentRepository(db)
	svc := NewDocumentService(repo)

	// 创建测试文档
	ctx := context.Background()
	doc := &model.Document{
		ID:          "test-1",
		Title:       "测试文档",
		Category:    "系统概述",
		Path:        "/test/doc1",
		Description: "测试文档描述",
		Content:     "测试文档内容",
		Size:        1024,
	}

	createdDoc, err := svc.CreateDocument(ctx, doc)
	assert.NoError(t, err)
	assert.NotNil(t, createdDoc)
	assert.Equal(t, "测试文档", createdDoc.Title)
}

// TestDocumentService_UpdateDocument 测试更新文档
func TestDocumentService_UpdateDocument(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := repository.NewDocumentRepository(db)
	svc := NewDocumentService(repo)

	ctx := context.Background()

	// 创建测试文档
	doc := &model.Document{
		ID:          "test-1",
		Title:       "测试文档",
		Category:    "系统概述",
		Path:        "/test/doc1",
		Description: "测试文档描述",
		Content:     "测试文档内容",
		Size:        1024,
	}
	err = repo.Create(doc)
	assert.NoError(t, err)

	// 测试更新文档
	updateDoc := &model.Document{
		Title:       "更新后的文档",
		Description: "更新后的描述",
	}
	updatedDoc, err := svc.UpdateDocument(ctx, "test-1", updateDoc)
	assert.NoError(t, err)
	assert.Equal(t, "更新后的文档", updatedDoc.Title)
	assert.Equal(t, "更新后的描述", updatedDoc.Description)
}

// TestDocumentService_DeleteDocument 测试删除文档
func TestDocumentService_DeleteDocument(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := repository.NewDocumentRepository(db)
	svc := NewDocumentService(repo)

	// 创建测试文档
	doc := &model.Document{
		ID:          "test-1",
		Title:       "测试文档",
		Category:    "系统概述",
		Path:        "/test/doc1",
		Description: "测试文档描述",
		Content:     "测试文档内容",
		Size:        1024,
	}
	err = repo.Create(doc)
	assert.NoError(t, err)

	// 测试删除文档
	err = svc.DeleteDocument("test-1")
	assert.NoError(t, err)

	// 验证文档已删除
	_, err = repo.GetByID("test-1")
	assert.Error(t, err)
}

// TestDocumentService_SearchDocuments 测试搜索文档
func TestDocumentService_SearchDocuments(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := repository.NewDocumentRepository(db)
	svc := NewDocumentService(repo)

	// 创建测试数据
	docs := []*model.Document{
		{
			ID:          "test-1",
			Title:       "Go 语言入门",
			Category:    "技术架构",
			Path:        "/go/intro",
			Description: "Go 语言基础教程",
			Content:     "Go 语言是一种静态类型编程语言",
			Size:        2048,
		},
		{
			ID:          "test-2",
			Title:       "Python 进阶",
			Category:    "技术架构",
			Path:        "/python/advanced",
			Description: "Python 高级编程技巧",
			Content:     "Python 是一种动态类型编程语言",
			Size:        3072,
		},
	}

	for _, doc := range docs {
		err = repo.Create(doc)
		assert.NoError(t, err)
	}

	// 测试搜索
	results, err := svc.SearchDocuments("Go", "", 10)
	assert.NoError(t, err)
	assert.Greater(t, len(results), 0)
}

// TestDocumentService_Favorite 测试收藏功能
func TestDocumentService_Favorite(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := repository.NewDocumentRepository(db)
	svc := NewDocumentService(repo)

	ctx := context.Background()

	// 创建测试文档
	doc := &model.Document{
		ID:          "test-1",
		Title:       "测试文档",
		Category:    "系统概述",
		Path:        "/test/doc1",
		Description: "测试文档描述",
		Content:     "测试文档内容",
		Size:        1024,
	}
	err = repo.Create(doc)
	assert.NoError(t, err)

	// 测试添加收藏
	favorited, err := svc.ToggleFavorite(ctx, "test-1")
	assert.NoError(t, err)
	assert.True(t, favorited)

	// 测试取消收藏
	favorited, err = svc.ToggleFavorite(ctx, "test-1")
	assert.NoError(t, err)
	assert.False(t, favorited)
}

// TestDocumentService_ReadProgress 测试阅读进度
func TestDocumentService_ReadProgress(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := repository.NewDocumentRepository(db)
	svc := NewDocumentService(repo)

	ctx := context.Background()

	// 创建测试文档
	doc := &model.Document{
		ID:          "test-1",
		Title:       "测试文档",
		Category:    "系统概述",
		Path:        "/test/doc1",
		Description: "测试文档描述",
		Content:     "测试文档内容",
		Size:        1024,
	}
	err = repo.Create(doc)
	assert.NoError(t, err)

	// 测试更新阅读进度
	err = svc.UpdateReadProgress(ctx, "test-1", 50)
	assert.NoError(t, err)

	// 测试获取阅读进度
	progress, err := svc.GetMyReadProgress(ctx, "test-1")
	assert.NoError(t, err)
	assert.NotNil(t, progress)
	assert.Equal(t, 50, progress.ReadPosition)
}

// TestDocumentService_Versions 测试版本管理
func TestDocumentService_Versions(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	repo := repository.NewDocumentRepository(db)
	svc := NewDocumentService(repo)

	// 创建测试文档
	doc := &model.Document{
		ID:          "test-1",
		Title:       "测试文档",
		Category:    "系统概述",
		Path:        "/test/doc1",
		Description: "测试文档描述",
		Content:     "测试文档内容",
		Size:        1024,
		Version:     1,
	}
	err = repo.Create(doc)
	assert.NoError(t, err)

	// 测试获取版本历史
	versions, err := svc.GetDocumentVersions("test-1")
	assert.NoError(t, err)
	assert.Len(t, versions, 0) // 初始没有版本历史
}
