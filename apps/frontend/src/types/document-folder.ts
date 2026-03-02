/**
 * 文档文件夹/目录类型定义
 */

// ==================== 数据模型类型 ====================

/**
 * 文档文件夹
 */
export interface DocumentFolder {
  id: string
  name: string
  path: string
  parentId: string
  level: number
  description: string
  sortOrder: number
  isEnabled: boolean
  createdBy: string
  updatedBy: string
  createdAt: string
  updatedAt: string
  // 非数据库字段
  children?: DocumentFolder[]
  docCount?: number
  hasChildren?: boolean
}

/**
 * 文件夹树形结构节点
 */
export interface DocumentFolderTree {
  id: string
  name: string
  path: string
  level: number
  description?: string
  sortOrder?: number
  isEnabled?: boolean
  children?: DocumentFolderTree[]
  docCount?: number
  hasChildren?: boolean
}

// ==================== 请求参数类型 ====================

/**
 * 创建文件夹参数
 */
export interface CreateFolderParams {
  name: string
  parentId?: string
  description?: string
  sortOrder?: number
  isEnabled?: boolean
}

/**
 * 更新文件夹参数
 */
export interface UpdateFolderParams {
  name?: string
  description?: string
  sortOrder?: number
  isEnabled?: boolean
}

/**
 * 移动文件夹参数
 */
export interface MoveFolderParams {
  newParentId: string
}

/**
 * 复制文件夹参数
 */
export interface CopyFolderParams {
  newParentId: string
  newName?: string
}

// ==================== 响应类型 ====================

/**
 * 文件夹 API 响应
 */
export interface FolderResponse<T = any> {
  success: boolean
  message: string
  data?: T
  error?: string
}
