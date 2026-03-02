package service

import (
	"errors"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	"metadata-platform/internal/utils"
	"metadata-platform/internal/module/document/model"
	"metadata-platform/internal/module/document/repository"
)

// DocumentFolderService 文档目录服务接口
type DocumentFolderService interface {
	// 文件夹基本操作
	CreateFolder(ctx *app.RequestContext, folder *model.DocumentFolder) (*model.DocumentFolder, error)
	UpdateFolder(ctx *app.RequestContext, id string, folder *model.DocumentFolder) (*model.DocumentFolder, error)
	DeleteFolder(id string) error
	GetFolderByID(id string) (*model.DocumentFolder, error)
	GetFolderByPath(path string) (*model.DocumentFolder, error)
	
	// 文件夹列表
	GetFolderList(ctx *app.RequestContext, parentID string, enabledOnly bool) ([]*model.DocumentFolder, error)
	GetFolderTree() ([]*model.DocumentFolderTree, error)
	
	// 移动和复制
	MoveFolder(folderID, newParentID string) error
	CopyFolder(folderID, newParentID string) (*model.DocumentFolder, error)
	
	// 验证
	ValidateFolder(folder *model.DocumentFolder) error
}

// documentFolderService 文档目录服务实现
type documentFolderService struct {
	repo repository.DocumentFolderRepository
}

// NewDocumentFolderService 创建文档目录服务实例
func NewDocumentFolderService(repo repository.DocumentFolderRepository) DocumentFolderService {
	return &documentFolderService{
		repo: repo,
	}
}

// CreateFolder 创建文件夹
func (s *documentFolderService) CreateFolder(ctx *app.RequestContext, folder *model.DocumentFolder) (*model.DocumentFolder, error) {
	// 验证文件夹
	if err := s.ValidateFolder(folder); err != nil {
		return nil, err
	}
	
	// 先生成 ID（用于路径生成）
	folder.ID = utils.GetSnowflake().GenerateIDString()
	
	// 生成路径（处理空字符串的情况）
	if folder.ParentID != "" && folder.ParentID != "null" {
		parent, err := s.repo.GetFolderByID(folder.ParentID)
		if err != nil {
			return nil, errors.New("父文件夹不存在")
		}
		folder.Path = parent.Path + "/" + folder.ID
		folder.Level = parent.Level + 1
	} else {
		folder.Path = "/" + folder.ID
		folder.Level = 0
	}
	
	// 设置创建人
	folder.CreatedBy = s.getCurrentUserID(ctx)
	folder.UpdatedBy = s.getCurrentUserID(ctx)
	
	// 创建文件夹
	if err := s.repo.CreateFolder(folder); err != nil {
		return nil, err
	}
	
	return folder, nil
}

// UpdateFolder 更新文件夹
func (s *documentFolderService) UpdateFolder(ctx *app.RequestContext, id string, folder *model.DocumentFolder) (*model.DocumentFolder, error) {
	// 获取原文件夹
	oldFolder, err := s.repo.GetFolderByID(id)
	if err != nil {
		return nil, errors.New("文件夹不存在")
	}
	
	// 保留原 ID 和路径，用于验证
	folder.ID = id
	folder.ParentID = oldFolder.ParentID
	
	// 验证文件夹（包括同级重名检查）
	if err := s.ValidateFolder(folder); err != nil {
		return nil, err
	}
	
	// 更新字段
	oldFolder.Name = folder.Name
	oldFolder.Description = folder.Description
	oldFolder.SortOrder = folder.SortOrder
	oldFolder.IsEnabled = folder.IsEnabled
	oldFolder.UpdatedBy = s.getCurrentUserID(ctx)
	
	// 保存更新
	if err := s.repo.UpdateFolder(oldFolder); err != nil {
		return nil, err
	}
	
	return oldFolder, nil
}

// DeleteFolder 删除文件夹
func (s *documentFolderService) DeleteFolder(id string) error {
	// 检查是否有子文件夹
	hasChildren, err := s.repo.HasSubFolders(id)
	if err != nil {
		return err
	}
	
	if hasChildren {
		return errors.New("无法删除包含子文件夹的文件夹")
	}
	
	// 检查是否有文档
	docCount, err := s.repo.GetFolderDocCount(id)
	if err != nil {
		return err
	}
	
	if docCount > 0 {
		return errors.New("无法删除包含文档的文件夹")
	}
	
	return s.repo.DeleteFolder(id)
}

// GetFolderByID 根据 ID 获取文件夹
func (s *documentFolderService) GetFolderByID(id string) (*model.DocumentFolder, error) {
	return s.repo.GetFolderByID(id)
}

// GetFolderByPath 根据路径获取文件夹
func (s *documentFolderService) GetFolderByPath(path string) (*model.DocumentFolder, error) {
	return s.repo.GetFolderByPath(path)
}

// GetFolderList 获取文件夹列表
func (s *documentFolderService) GetFolderList(ctx *app.RequestContext, parentID string, enabledOnly bool) ([]*model.DocumentFolder, error) {
	return s.repo.GetFolderList(parentID, enabledOnly)
}

// GetFolderTree 获取文件夹树
func (s *documentFolderService) GetFolderTree() ([]*model.DocumentFolderTree, error) {
	return s.repo.GetFolderTree()
}

// MoveFolder 移动文件夹
func (s *documentFolderService) MoveFolder(folderID, newParentID string) error {
	return s.repo.MoveFolder(folderID, newParentID)
}

// CopyFolder 复制文件夹
func (s *documentFolderService) CopyFolder(folderID, newParentID string) (*model.DocumentFolder, error) {
	return s.repo.CopyFolder(folderID, newParentID)
}

// ValidateFolder 验证文件夹
func (s *documentFolderService) ValidateFolder(folder *model.DocumentFolder) error {
	if folder.Name == "" {
		return errors.New("文件夹名称不能为空")
	}
	
	if len(folder.Name) > 255 {
		return errors.New("文件夹名称不能超过 255 个字符")
	}
	
	if folder.Description != "" && len(folder.Description) > 512 {
		return errors.New("描述不能超过 512 个字符")
	}
	
	// 检查名称是否合法（不能包含特殊字符）
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(folder.Name, char) {
			return errors.New("文件夹名称不能包含特殊字符：" + strings.Join(invalidChars, ", "))
		}
	}
	
	// 检查同级目录下是否有重名（排除自身）
	siblings, err := s.repo.GetFolderList(folder.ParentID, false)
	if err == nil && siblings != nil {
		for _, sibling := range siblings {
			// 如果是更新操作，排除自身
			if folder.ID != "" && sibling.ID == folder.ID {
				continue
			}
			// 检查名称是否相同（忽略大小写）
			if strings.EqualFold(sibling.Name, folder.Name) {
				return errors.New("同级目录下已存在同名文件夹：" + folder.Name)
			}
		}
	}
	
	return nil
}

// getCurrentUserID 获取当前用户 ID
func (s *documentFolderService) getCurrentUserID(ctx *app.RequestContext) string {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return "system"
	}
	
	if uid, ok := userID.(string); ok {
		return uid
	}
	
	return "system"
}
