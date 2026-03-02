/**
 * 文档文件夹管理模块 API 服务
 * 用于管理文档目录/文件夹的 CRUD 操作
 */
import request from '@/utils/request'
import type { DocumentFolder, DocumentFolderTree, CreateFolderParams, UpdateFolderParams, MoveFolderParams, CopyFolderParams } from '@/types/document-folder'

// ==================== 文件夹 CRUD API ====================

/**
 * 创建文件夹
 * @param data 创建参数
 */
export const createFolder = async (data: CreateFolderParams): Promise<DocumentFolder> => {
  return request({
    url: '/api/documents/folders',
    method: 'post',
    data
  })
}

/**
 * 更新文件夹
 * @param id 文件夹 ID
 * @param data 更新参数
 */
export const updateFolder = async (id: string, data: UpdateFolderParams): Promise<DocumentFolder> => {
  return request({
    url: `/api/documents/folders/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除文件夹
 * @param id 文件夹 ID
 */
export const deleteFolder = async (id: string): Promise<void> => {
  return request({
    url: `/api/documents/folders/${id}`,
    method: 'delete'
  })
}

/**
 * 根据 ID 获取文件夹详情
 * @param id 文件夹 ID
 */
export const getFolderById = async (id: string): Promise<DocumentFolder> => {
  return request({
    url: `/api/documents/folders/${id}`,
    method: 'get'
  })
}

/**
 * 根据路径获取文件夹
 * @param path 文件夹路径
 */
export const getFolderByPath = async (path: string): Promise<DocumentFolder> => {
  return request({
    url: '/api/documents/folders/by-path',
    method: 'get',
    params: { path }
  })
}

/**
 * 获取文件夹列表
 * @param params 查询参数
 */
export const getFolderList = async (params?: {
  parentId?: string
  isEnabled?: boolean
}): Promise<DocumentFolder[]> => {
  return request({
    url: '/api/documents/folders',
    method: 'get',
    params
  })
}

/**
 * 获取文件夹树形结构
 * @param parentId 父文件夹 ID（可选，不传则从根节点开始）
 */
export const getFolderTree = async (parentId?: string): Promise<DocumentFolderTree[]> => {
  const res: any = await request({
    url: '/api/documents/folders/tree',
    method: 'get',
    params: { parentId }
  })
  // 后端返回格式：{ success: true, data: [...] }
  return res.data || res || []
}

// ==================== 文件夹操作 API ====================

/**
 * 移动文件夹
 * @param id 文件夹 ID
 * @param data 移动参数
 */
export const moveFolder = async (id: string, data: MoveFolderParams): Promise<DocumentFolder> => {
  return request({
    url: `/api/documents/folders/${id}/move`,
    method: 'post',
    data
  })
}

/**
 * 复制文件夹
 * @param id 文件夹 ID
 * @param data 复制参数
 */
export const copyFolder = async (id: string, data: CopyFolderParams): Promise<DocumentFolder> => {
  return request({
    url: `/api/documents/folders/${id}/copy`,
    method: 'post',
    data
  })
}
