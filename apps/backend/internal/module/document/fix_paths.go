package document

import (
	"fmt"
	"metadata-platform/internal/module/document/model"
	"gorm.io/gorm"
)

// FixFolderPaths 修复文件夹路径
// 用于修复旧数据的路径问题
func FixFolderPaths(db *gorm.DB) error {
	var folders []*model.DocumentFolder
	if err := db.Find(&folders).Error; err != nil {
		return err
	}
	
	for _, folder := range folders {
		// 如果路径不正确（如 "/" 或空字符串），则修复它
		if folder.Path == "/" || folder.Path == "" {
			if folder.ParentID == "" {
				// 根目录
				folder.Path = "/" + folder.ID
				folder.Level = 0
			} else {
				// 子目录 - 需要找到父文件夹
				var parent model.DocumentFolder
				if err := db.First(&parent, "id = ?", folder.ParentID).Error; err == nil {
					folder.Path = parent.Path + "/" + folder.ID
					folder.Level = parent.Level + 1
				}
			}
			
			// 保存修复后的文件夹
			if err := db.Save(folder).Error; err != nil {
				return fmt.Errorf("修复文件夹 %s 失败：%w", folder.ID, err)
			}
			fmt.Printf("修复文件夹：%s -> %s\n", folder.ID, folder.Path)
		}
	}
	
	return nil
}
