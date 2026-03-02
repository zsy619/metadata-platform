package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"metadata-platform/internal/module/document/model"
	"metadata-platform/internal/module/document/repository"
	"metadata-platform/internal/utils"
	"strings"
)

// DocumentService 文档服务接口
type DocumentService interface {
	// 文档管理
	CreateDocument(ctx context.Context, doc *model.Document) (*model.Document, error)
	UpdateDocument(ctx context.Context, id string, doc *model.Document) (*model.Document, error)
	DeleteDocument(id string) error
	GetDocumentByID(id string) (*model.Document, error)
	GetDocumentByPath(path string) (*model.Document, error)
	GetDocumentList(page, pageSize int, category, keyword string) ([]*model.Document, int64, error)

	// 分类管理
	GetCategories() ([]*model.DocumentCategory, error)

	// 搜索
	SearchDocuments(keyword string, category string, limit int) ([]*model.Document, error)

	// 下载
	GetDocumentContent(id string) (string, error)

	// 收藏管理
	ToggleFavorite(ctx context.Context, documentID string) (bool, error)
	GetMyFavorites(ctx context.Context) ([]*model.Document, error)

	// 阅读进度
	UpdateReadProgress(ctx context.Context, documentID string, position int) error
	GetMyReadProgress(ctx context.Context, documentID string) (*model.DocumentReadProgress, error)

	// 版本管理
	GetDocumentVersions(documentID string) ([]*model.DocumentVersion, error)
	RestoreVersion(documentID string, version int) (*model.Document, error)
}

// documentService 文档服务实现
type documentService struct {
	repo repository.DocumentRepository
}

// NewDocumentService 创建文档服务实例
func NewDocumentService(repo repository.DocumentRepository) DocumentService {
	return &documentService{
		repo: repo,
	}
}

// CreateDocument 创建文档
func (s *documentService) CreateDocument(ctx context.Context, doc *model.Document) (*model.Document, error) {
	// 生成 ID
	doc.ID = utils.GetSnowflake().GenerateIDString()

	// 计算文档大小
	doc.Size = int64(len(doc.Content))

	// 解析标签
	if doc.Tags != "" {
		if !strings.HasPrefix(doc.Tags, "[") {
			// 如果是逗号分隔的字符串，转换为 JSON 数组
			tags := strings.Split(doc.Tags, ",")
			tagsJSON, _ := json.Marshal(tags)
			doc.Tags = string(tagsJSON)
		}
	}

	// 生成目录
	doc.TOC = s.generateTOC(doc.Content)

	// 创建文档
	if err := s.repo.Create(doc); err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	// 创建初始版本
	version := &model.DocumentVersion{
		ID:         utils.GetSnowflake().GenerateIDString(),
		DocumentID: doc.ID,
		Version:    1,
		Content:    doc.Content,
		Size:       doc.Size,
		CreatedBy:  s.getCurrentUserID(ctx),
	}
	if err := s.repo.CreateVersion(version); err != nil {
		// 记录错误但不中断流程
		fmt.Printf("Failed to create document version: %v\n", err)
	}

	return doc, nil
}

// UpdateDocument 更新文档
func (s *documentService) UpdateDocument(ctx context.Context, id string, doc *model.Document) (*model.Document, error) {
	// 获取原文档
	oldDoc, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	// 更新字段
	oldDoc.Title = doc.Title
	oldDoc.Category = doc.Category
	oldDoc.Path = doc.Path
	oldDoc.Description = doc.Description
	oldDoc.Content = doc.Content
	oldDoc.Tags = doc.Tags
	oldDoc.Version = oldDoc.Version + 1
	oldDoc.UpdatedBy = s.getCurrentUserID(ctx)

	// 计算大小
	oldDoc.Size = int64(len(doc.Content))

	// 生成目录
	oldDoc.TOC = s.generateTOC(doc.Content)

	// 保存更新
	if err := s.repo.Update(oldDoc); err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	// 创建新版本
	version := &model.DocumentVersion{
		ID:         utils.GetSnowflake().GenerateIDString(),
		DocumentID: id,
		Version:    oldDoc.Version,
		Content:    doc.Content,
		Size:       oldDoc.Size,
		ChangeLog:  doc.Description,
		CreatedBy:  s.getCurrentUserID(ctx),
	}
	if err := s.repo.CreateVersion(version); err != nil {
		fmt.Printf("Failed to create document version: %v\n", err)
	}

	return oldDoc, nil
}

// DeleteDocument 删除文档
func (s *documentService) DeleteDocument(id string) error {
	return s.repo.Delete(id)
}

// GetDocumentByID 根据 ID 获取文档
func (s *documentService) GetDocumentByID(id string) (*model.Document, error) {
	doc, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 增加阅读次数
	go func() {
		if err := s.repo.IncrementViewCount(id); err != nil {
			fmt.Printf("Failed to increment view count: %v\n", err)
		}
	}()

	return doc, nil
}

// GetDocumentByPath 根据路径获取文档
func (s *documentService) GetDocumentByPath(path string) (*model.Document, error) {
	return s.repo.GetByPath(path)
}

// GetDocumentList 获取文档列表
func (s *documentService) GetDocumentList(page, pageSize int, category, keyword string) ([]*model.Document, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return s.repo.GetList(page, pageSize, category, keyword)
}

// GetCategories 获取分类列表
func (s *documentService) GetCategories() ([]*model.DocumentCategory, error) {
	return s.repo.GetCategories()
}

// SearchDocuments 搜索文档
func (s *documentService) SearchDocuments(keyword string, category string, limit int) ([]*model.Document, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.Search(keyword, category, limit)
}

// GetDocumentContent 获取文档内容
func (s *documentService) GetDocumentContent(id string) (string, error) {
	return s.repo.GetContent(id)
}

// ToggleFavorite 切换收藏状态
func (s *documentService) ToggleFavorite(ctx context.Context, documentID string) (bool, error) {
	userID := s.getCurrentUserID(ctx)

	// 检查是否已收藏
	favorited, err := s.repo.IsFavorited(documentID, userID)
	if err != nil {
		return false, err
	}

	if favorited {
		// 取消收藏
		if err := s.repo.RemoveFavorite(documentID, userID); err != nil {
			return false, err
		}
		return false, nil
	} else {
		// 添加收藏
		if err := s.repo.AddFavorite(documentID, userID); err != nil {
			return false, err
		}
		return true, nil
	}
}

// GetMyFavorites 获取我的收藏
func (s *documentService) GetMyFavorites(ctx context.Context) ([]*model.Document, error) {
	userID := s.getCurrentUserID(ctx)
	return s.repo.GetFavorites(userID)
}

// UpdateReadProgress 更新阅读进度
func (s *documentService) UpdateReadProgress(ctx context.Context, documentID string, position int) error {
	userID := s.getCurrentUserID(ctx)

	if position < 0 {
		position = 0
	}
	if position > 100 {
		position = 100
	}

	return s.repo.UpdateReadProgress(documentID, userID, position)
}

// GetMyReadProgress 获取我的阅读进度
func (s *documentService) GetMyReadProgress(ctx context.Context, documentID string) (*model.DocumentReadProgress, error) {
	userID := s.getCurrentUserID(ctx)
	return s.repo.GetReadProgress(documentID, userID)
}

// GetDocumentVersions 获取文档版本历史
func (s *documentService) GetDocumentVersions(documentID string) ([]*model.DocumentVersion, error) {
	return s.repo.GetVersions(documentID)
}

// RestoreVersion 恢复文档版本
func (s *documentService) RestoreVersion(documentID string, version int) (*model.Document, error) {
	// 获取指定版本
	versions, err := s.repo.GetVersions(documentID)
	if err != nil {
		return nil, err
	}

	var targetVersion *model.DocumentVersion
	for _, v := range versions {
		if v.Version == version {
			targetVersion = v
			break
		}
	}

	if targetVersion == nil {
		return nil, errors.New("version not found")
	}

	// 获取原文档
	doc, err := s.repo.GetByID(documentID)
	if err != nil {
		return nil, err
	}

	// 恢复内容
	doc.Content = targetVersion.Content
	doc.Size = targetVersion.Size
	doc.Version = doc.Version + 1

	// 保存
	if err := s.repo.Update(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// generateTOC 生成目录结构（简单实现）
func (s *documentService) generateTOC(content string) string {
	lines := strings.Split(content, "\n")
	toc := []map[string]interface{}{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			// 解析标题级别
			level := 0
			for _, ch := range line {
				if ch == '#' {
					level++
				} else {
					break
				}
			}

			if level <= 3 { // 只处理 1-3 级标题
				title := strings.TrimSpace(strings.TrimPrefix(line, strings.Repeat("#", level)))
				toc = append(toc, map[string]interface{}{
					"level": level,
					"title": title,
				})
			}
		}
	}

	tocJSON, _ := json.Marshal(toc)
	return string(tocJSON)
}

// getCurrentUserID 获取当前用户 ID
func (s *documentService) getCurrentUserID(ctx context.Context) string {
	// 从上下文中获取用户 ID
	// 注意：context.Context 本身不支持 Get 方法，需要通过 value 传递
	// 在实际使用中，用户 ID 应该通过 context.WithValue 传递
	userID := ctx.Value("user_id")
	if userID == nil {
		return "system"
	}

	if uid, ok := userID.(string); ok {
		return uid
	}

	return "system"
}

// ValidateDocument 验证文档
func (s *documentService) ValidateDocument(doc *model.Document) error {
	if doc.Title == "" {
		return errors.New("title is required")
	}

	if doc.Category == "" {
		return errors.New("category is required")
	}

	if doc.Path == "" {
		return errors.New("path is required")
	}

	if doc.Content == "" {
		return errors.New("content is required")
	}

	return nil
}
