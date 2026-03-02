package repository

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"metadata-platform/internal/module/document/model"
)

// DocumentFolderRepository 文档目录数据访问接口
type DocumentFolderRepository interface {
	// 文件夹基本操作
	CreateFolder(folder *model.DocumentFolder) error
	UpdateFolder(folder *model.DocumentFolder) error
	DeleteFolder(id string) error
	GetFolderByID(id string) (*model.DocumentFolder, error)
	GetFolderByPath(path string) (*model.DocumentFolder, error)

	// 文件夹列表
	GetFolderList(parentID string, enabledOnly bool) ([]*model.DocumentFolder, error)
	GetAllFolders() ([]*model.DocumentFolder, error)

	// 树形结构
	GetFolderTree() ([]*model.DocumentFolderTree, error)
	buildFolderTree(folders []*model.DocumentFolder, parentID string) []*model.DocumentFolderTree

	// 移动和复制
	MoveFolder(folderID, newParentID string) error
	CopyFolder(folderID, newParentID string) (*model.DocumentFolder, error)

	// 统计
	GetFolderDocCount(folderID string) (int64, error)
	HasSubFolders(folderID string) (bool, error)
}

// documentFolderRepository 文档目录数据访问实现
type documentFolderRepository struct {
	db *gorm.DB
}

// NewDocumentFolderRepository 创建文档目录数据访问实例
func NewDocumentFolderRepository(db *gorm.DB) DocumentFolderRepository {
	return &documentFolderRepository{
		db: db,
	}
}

// CreateFolder 创建文件夹
func (r *documentFolderRepository) CreateFolder(folder *model.DocumentFolder) error {
	// 如果 ID 为空，则生成（兼容 Service 层已生成 ID 的情况）
	if folder.ID == "" {
		folder.ID = uuid.New().String()
	}
	return r.db.Create(folder).Error
}

// UpdateFolder 更新文件夹
func (r *documentFolderRepository) UpdateFolder(folder *model.DocumentFolder) error {
	return r.db.Save(folder).Error
}

// DeleteFolder 删除文件夹
func (r *documentFolderRepository) DeleteFolder(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除文件夹
		if err := tx.Where("id = ?", id).Delete(&model.DocumentFolder{}).Error; err != nil {
			return err
		}

		// 递归删除子文件夹
		var subFolders []*model.DocumentFolder
		if err := tx.Where("parent_id = ?", id).Find(&subFolders).Error; err != nil {
			return err
		}

		for _, subFolder := range subFolders {
			if err := r.DeleteFolder(subFolder.ID); err != nil {
				return err
			}
		}

		return nil
	})
}

// GetFolderByID 根据 ID 获取文件夹
func (r *documentFolderRepository) GetFolderByID(id string) (*model.DocumentFolder, error) {
	var folder model.DocumentFolder
	err := r.db.Where("id = ?", id).First(&folder).Error
	if err != nil {
		return nil, err
	}

	// 统计文档数量
	folder.DocCount, _ = r.GetFolderDocCount(id)
	folder.HasChildren, _ = r.HasSubFolders(id)

	return &folder, nil
}

// GetFolderByPath 根据路径获取文件夹
func (r *documentFolderRepository) GetFolderByPath(path string) (*model.DocumentFolder, error) {
	var folder model.DocumentFolder
	err := r.db.Where("path = ?", path).First(&folder).Error
	if err != nil {
		return nil, err
	}

	// 统计文档数量
	folder.DocCount, _ = r.GetFolderDocCount(folder.ID)
	folder.HasChildren, _ = r.HasSubFolders(folder.ID)

	return &folder, nil
}

// GetFolderList 获取文件夹列表
func (r *documentFolderRepository) GetFolderList(parentID string, enabledOnly bool) ([]*model.DocumentFolder, error) {
	var folders []*model.DocumentFolder
	query := r.db.Where("parent_id = ?", parentID)

	if enabledOnly {
		query = query.Where("is_enabled = ?", true)
	}

	err := query.Order("sort_order ASC, name ASC").Find(&folders).Error
	if err != nil {
		return nil, err
	}

	// 为每个文件夹统计文档数量和子文件夹
	for _, folder := range folders {
		folder.DocCount, _ = r.GetFolderDocCount(folder.ID)
		folder.HasChildren, _ = r.HasSubFolders(folder.ID)
	}

	return folders, nil
}

// GetAllFolders 获取所有文件夹
func (r *documentFolderRepository) GetAllFolders() ([]*model.DocumentFolder, error) {
	var folders []*model.DocumentFolder
	err := r.db.Order("level ASC, sort_order ASC, name ASC").Find(&folders).Error
	return folders, err
}

// GetFolderTree 获取文件夹树形结构
func (r *documentFolderRepository) GetFolderTree() ([]*model.DocumentFolderTree, error) {
	folders, err := r.GetAllFolders()
	if err != nil {
		return nil, err
	}

	return r.buildFolderTree(folders, ""), nil
}

// buildFolderTree 递归构建文件夹树
func (r *documentFolderRepository) buildFolderTree(folders []*model.DocumentFolder, parentID string) []*model.DocumentFolderTree {
	var tree []*model.DocumentFolderTree

	for _, folder := range folders {
		if folder.ParentID == parentID {
			node := &model.DocumentFolderTree{
				ID:          folder.ID,
				Name:        folder.Name,
				Path:        folder.Path,
				Level:       folder.Level,
				DocCount:    folder.DocCount,
				HasChildren: folder.HasChildren,
			}

			// 递归查找子文件夹
			node.Children = r.buildFolderTree(folders, folder.ID)

			tree = append(tree, node)
		}
	}

	return tree
}

// MoveFolder 移动文件夹
func (r *documentFolderRepository) MoveFolder(folderID, newParentID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 获取原文件夹
		var folder model.DocumentFolder
		if err := tx.Where("id = ?", folderID).First(&folder).Error; err != nil {
			return err
		}

		// 获取新父文件夹
		var newParent model.DocumentFolder
		if newParentID != "" {
			if err := tx.Where("id = ?", newParentID).First(&newParent).Error; err != nil {
				return err
			}

			// 检查是否移动到子文件夹（避免循环）
			if strings.Contains(newParent.Path, "/"+folder.ID+"/") || newParent.ID == folder.ID {
				return gorm.ErrInvalidData
			}

			// 更新父 ID 和路径
			folder.ParentID = newParentID
			folder.Path = newParent.Path + "/" + folder.ID
			folder.Level = newParent.Level + 1
		} else {
			// 移动到根目录
			folder.ParentID = ""
			folder.Path = "/" + folder.ID
			folder.Level = 0
		}

		// 更新文件夹
		if err := tx.Save(&folder).Error; err != nil {
			return err
		}

		// 递归更新子文件夹路径
		return r.updateSubFolderPaths(tx, folder.ID, folder.Path, folder.Level)
	})
}

// updateSubFolderPaths 递归更新子文件夹路径
func (r *documentFolderRepository) updateSubFolderPaths(tx *gorm.DB, parentID, parentPath string, parentLevel int) error {
	var subFolders []*model.DocumentFolder
	if err := tx.Where("parent_id = ?", parentID).Find(&subFolders).Error; err != nil {
		return err
	}

	for _, subFolder := range subFolders {
		subFolder.Path = parentPath + "/" + subFolder.ID
		subFolder.Level = parentLevel + 1

		if err := tx.Save(&subFolder).Error; err != nil {
			return err
		}

		// 递归更新
		if err := r.updateSubFolderPaths(tx, subFolder.ID, subFolder.Path, subFolder.Level); err != nil {
			return err
		}
	}

	return nil
}

// CopyFolder 复制文件夹
func (r *documentFolderRepository) CopyFolder(folderID, newParentID string) (*model.DocumentFolder, error) {
	// 获取原文件夹
	var folder model.DocumentFolder
	if err := r.db.Where("id = ?", folderID).First(&folder).Error; err != nil {
		return nil, err
	}

	// 创建新文件夹
	newFolder := &model.DocumentFolder{
		ID:          uuid.New().String(),
		Name:        folder.Name,
		Description: folder.Description,
		SortOrder:   folder.SortOrder,
		IsEnabled:   folder.IsEnabled,
		CreatedBy:   "system",
		UpdatedBy:   "system",
	}

	// 设置父节点和路径
	if newParentID != "" {
		// 获取新父文件夹
		var newParent model.DocumentFolder
		if err := r.db.Where("id = ?", newParentID).First(&newParent).Error; err != nil {
			return nil, err
		}

		newFolder.ParentID = newParentID
		newFolder.Path = newParent.Path + "/" + newFolder.ID
		newFolder.Level = newParent.Level + 1
	} else {
		// 复制到根目录
		newFolder.ParentID = ""
		newFolder.Path = "/" + newFolder.ID
		newFolder.Level = 0
	}

	// 保存新文件夹
	if err := r.db.Create(&newFolder).Error; err != nil {
		return nil, err
	}

	// 递归复制子文件夹
	if err := r.copySubFolders(r.db, folder.ID, newFolder.ID, newFolder.Path, newFolder.Level); err != nil {
		return nil, err
	}

	// 统计文档数量（暂时返回 0）
	newFolder.DocCount = 0
	newFolder.HasChildren, _ = r.HasSubFolders(newFolder.ID)

	return newFolder, nil
}

// copySubFolders 递归复制子文件夹
func (r *documentFolderRepository) copySubFolders(tx *gorm.DB, parentID, newParentID, newParentPath string, newParentLevel int) error {
	var subFolders []*model.DocumentFolder
	if err := tx.Where("parent_id = ?", parentID).Find(&subFolders).Error; err != nil {
		return err
	}

	for _, subFolder := range subFolders {
		newSubFolderID := uuid.New().String()
		newSubFolder := &model.DocumentFolder{
			ID:          newSubFolderID,
			Name:        subFolder.Name,
			Description: subFolder.Description,
			SortOrder:   subFolder.SortOrder,
			IsEnabled:   subFolder.IsEnabled,
			ParentID:    newParentID,
			Path:        newParentPath + "/" + newSubFolderID,
			Level:       newParentLevel + 1,
			CreatedBy:   "system",
			UpdatedBy:   "system",
		}

		if err := tx.Create(&newSubFolder).Error; err != nil {
			return err
		}

		// 递归复制
		if err := r.copySubFolders(tx, subFolder.ID, newSubFolder.ID, newSubFolder.Path, newSubFolder.Level); err != nil {
			return err
		}
	}

	return nil
}

// GetFolderDocCount 获取文件夹下的文档数量
func (r *documentFolderRepository) GetFolderDocCount(folderID string) (int64, error) {
	var count int64
	// 这里假设文档有一个 folder_id 字段，如果没有可以暂时返回 0
	// err := r.db.Model(&model.Document{}).Where("folder_id = ?", folderID).Count(&count).Error
	return count, nil
}

// HasSubFolders 检查是否有子文件夹
func (r *documentFolderRepository) HasSubFolders(folderID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.DocumentFolder{}).Where("parent_id = ?", folderID).Count(&count).Error
	return count > 0, err
}
