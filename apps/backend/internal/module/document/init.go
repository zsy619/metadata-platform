package document

import (
	"fmt"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/utils"

	"gorm.io/gorm"
)

// InitMenus 初始化文档管理模块的菜单数据
// 在应用启动时调用，将文档管理相关的菜单添加到系统中
func InitMenus(db *gorm.DB) error {
	// 检查是否已经初始化过
	var count int64
	if err := db.Model(&model.SsoMenu{}).Where("app_code = ?", "document").Count(&count).Error; err != nil {
		return err
	}

	// 如果已经存在菜单，则跳过初始化
	if count > 0 {
		fmt.Println("文档管理菜单已存在，跳过初始化")
		return nil
	}

	fmt.Println("开始初始化文档管理菜单...")

	// 生成所有 ID
	rootID := utils.GetSnowflake().GenerateIDString()
	listID := utils.GetSnowflake().GenerateIDString()
	createFolderID := utils.GetSnowflake().GenerateIDString()
	editFolderID := utils.GetSnowflake().GenerateIDString()
	deleteFolderID := utils.GetSnowflake().GenerateIDString()
	moveFolderID := utils.GetSnowflake().GenerateIDString()
	copyFolderID := utils.GetSnowflake().GenerateIDString()
	viewDocID := utils.GetSnowflake().GenerateIDString()
	createDocID := utils.GetSnowflake().GenerateIDString()
	editDocID := utils.GetSnowflake().GenerateIDString()
	deleteDocID := utils.GetSnowflake().GenerateIDString()

	// 文档管理模块的菜单数据
	menus := []model.SsoMenu{
		{
			ID:        rootID,
			ParentID:  "",
			AppCode:   "document",
			MenuName:  "文档管理",
			MenuCode:  "document",
			Status:    1,
			IsVisible: true,
			MenuType:  "M", // 目录
			Icon:      "fa-file-lines",
			URL:       "/documents",
			Sort:      100,
			Tier:      1,
			Remark:    "文档管理模块",
			CreateID:  "system",
			CreateBy:  "system",
			TenantID:  "1",
		},
		{
			ID:        listID,
			ParentID:  rootID, // 父菜单为"文档管理"
			AppCode:   "document",
			MenuName:  "文档目录",
			MenuCode:  "document_list",
			Status:    1,
			IsVisible: true,
			MenuType:  "C", // 菜单
			Icon:      "fa-folder-tree",
			URL:       "/documents/list",
			Method:    "GET",
			Sort:      1,
			Tier:      2,
			Remark:    "文档目录管理",
			CreateID:  "system",
			CreateBy:  "system",
			TenantID:  "1",
		},
		// 文件夹管理按钮
		{
			ID:        createFolderID,
			ParentID:  rootID,
			AppCode:   "document",
			MenuName:  "新建文件夹",
			MenuCode:  "document_folder_create",
			Status:    1,
			IsVisible: true,
			MenuType:  "F", // 按钮
			URL:       "/api/documents/folders",
			Method:    "POST",
			Sort:      2,
			Tier:      2,
			Remark:    "创建文件夹权限",
			CreateID:  "system",
			CreateBy:  "system",
			TenantID:  "1",
		},
		{
			ID:        editFolderID,
			ParentID:  rootID,
			AppCode:   "document",
			MenuName:  "编辑文件夹",
			MenuCode:  "document_folder_edit",
			Status:    1,
			IsVisible: true,
			MenuType:  "F",
			URL:       "/api/documents/folders/*",
			Method:    "PUT",
			Sort:      3,
			Tier:      2,
			Remark:    "编辑文件夹权限",
			CreateID:  "system",
			CreateBy:  "system",
			TenantID:  "1",
		},
		{
			ID:        deleteFolderID,
			ParentID:  rootID,
			AppCode:   "document",
			MenuName:  "删除文件夹",
			MenuCode:  "document_folder_delete",
			Status:    1,
			IsVisible: true,
			MenuType:  "F",
			URL:       "/api/documents/folders/*",
			Method:    "DELETE",
			Sort:      4,
			Tier:      2,
			Remark:    "删除文件夹权限",
			CreateID:  "system",
			CreateBy:  "system",
			TenantID:  "1",
		},
		{
			ID:        moveFolderID,
			ParentID:  rootID,
			AppCode:   "document",
			MenuName:  "移动文件夹",
			MenuCode:  "document_folder_move",
			Status:    1,
			IsVisible: true,
			MenuType:  "F",
			URL:       "/api/documents/folders/*/move",
			Method:    "POST",
			Sort:      5,
			Tier:      2,
			Remark:    "移动文件夹权限",
			CreateID:  "system",
			CreateBy:  "system",
			TenantID:  "1",
		},
		{
			ID:        copyFolderID,
			ParentID:  rootID,
			AppCode:   "document",
			MenuName:  "复制文件夹",
			MenuCode:  "document_folder_copy",
			Status:    1,
			IsVisible: true,
			MenuType:  "F",
			URL:       "/api/documents/folders/*/copy",
			Method:    "POST",
			Sort:      6,
			Tier:      2,
			Remark:    "复制文件夹权限",
			CreateID:  "system",
			CreateBy:  "system",
			TenantID:  "1",
		},
		// 文档管理按钮
		{
			ID:        viewDocID,
			ParentID:  rootID,
			AppCode:   "document",
			MenuName:  "查看文档",
			MenuCode:  "document_view",
			Status:    1,
			IsVisible: true,
			MenuType:  "F",
			URL:       "/api/documents/*",
			Method:    "GET",
			Sort:      7,
			Tier:      2,
			Remark:    "查看文档权限",
			CreateID:  "system",
			CreateBy:  "system",
			TenantID:  "1",
		},
		{
			ID:        createDocID,
			ParentID:  rootID,
			AppCode:   "document",
			MenuName:  "创建文档",
			MenuCode:  "document_create",
			Status:    1,
			IsVisible: true,
			MenuType:  "F",
			URL:       "/api/documents",
			Method:    "POST",
			Sort:      8,
			Tier:      2,
			Remark:    "创建文档权限",
			CreateID:  "system",
			CreateBy:  "system",
			TenantID:  "1",
		},
		{
			ID:        editDocID,
			ParentID:  rootID,
			AppCode:   "document",
			MenuName:  "编辑文档",
			MenuCode:  "document_edit",
			Status:    1,
			IsVisible: true,
			MenuType:  "F",
			URL:       "/api/documents/*",
			Method:    "PUT",
			Sort:      9,
			Tier:      2,
			Remark:    "编辑文档权限",
			CreateID:  "system",
			CreateBy:  "system",
			TenantID:  "1",
		},
		{
			ID:        deleteDocID,
			ParentID:  rootID,
			AppCode:   "document",
			MenuName:  "删除文档",
			MenuCode:  "document_delete",
			Status:    1,
			IsVisible: true,
			MenuType:  "F",
			URL:       "/api/documents/*",
			Method:    "DELETE",
			Sort:      10,
			Tier:      2,
			Remark:    "删除文档权限",
			CreateID:  "system",
			CreateBy:  "system",
			TenantID:  "1",
		},
	}

	// 批量插入菜单数据
	if err := db.Create(&menus).Error; err != nil {
		return fmt.Errorf("初始化文档管理菜单失败：%w", err)
	}

	fmt.Printf("文档管理菜单初始化完成，共创建 %d 个菜单\n", len(menus))
	return nil
}
