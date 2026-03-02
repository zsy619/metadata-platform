/**
 * 文档管理模块 API 服务
 * 用于获取系统文档、技术文档、帮助文档等
 */
import request from '@/utils/request'
import type { DocumentInfo, DocumentDetail, DocumentQueryParams } from '@/types/document'

// ==================== 文档列表 API ====================

/**
 * 获取文档列表
 * @param params 查询参数
 * @param params.category 分类筛选
 * @param params.keyword 关键词搜索
 * @param params.page 页码
 * @param params.pageSize 每页数量
 */
export const getDocumentList = async (params?: DocumentQueryParams): Promise<any> => {
  return request({
    url: '/api/documents',
    method: 'get',
    params
  })
}

/**
 * 获取所有文档分类
 */
export const getDocumentCategories = async (): Promise<any> => {
  return request({
    url: '/api/documents/categories',
    method: 'get'
  })
}

// ==================== 文档详情 API ====================

/**
 * 根据 ID 获取文档详情
 * @param id 文档 ID
 */
export const getDocumentById = async (id: string): Promise<DocumentDetail> => {
  return request({
    url: `/api/documents/${id}`,
    method: 'get'
  })
}

/**
 * 根据路径获取文档详情
 * @param path 文档路径
 */
export const getDocumentByPath = async (path: string): Promise<DocumentDetail> => {
  return request({
    url: '/api/documents/by-path',
    method: 'get',
    params: { path }
  })
}

/**
 * 获取文档内容（仅内容，用于懒加载）
 * @param id 文档 ID
 */
export const getDocumentContent = async (id: string): Promise<string> => {
  return request({
    url: `/api/documents/${id}/content`,
    method: 'get'
  })
}

// ==================== 文档搜索 API ====================

/**
 * 搜索文档
 * @param keyword 搜索关键词
 * @param options 搜索选项
 */
export const searchDocuments = async (
  keyword: string,
  options?: {
    category?: string
    limit?: number
  }
): Promise<DocumentInfo[]> => {
  return request({
    url: '/api/documents/search',
    method: 'get',
    params: {
      keyword,
      ...options
    }
  })
}

// ==================== 文档下载 API ====================

/**
 * 下载文档
 * @param id 文档 ID
 * @param format 下载格式
 */
export const downloadDocument = async (id: string, format: 'md' | 'html' | 'pdf'): Promise<Blob> => {
  return request({
    url: `/api/documents/${id}/download`,
    method: 'get',
    params: { format },
    responseType: 'blob'
  })
}

// ==================== 文档操作 API ====================

/**
 * 创建文档
 * @param data 创建参数
 */
export const createDocument = async (data: {
  title: string
  category: string
  path: string
  description?: string
  content: string
  tags?: string[]
  isPublished?: boolean
}): Promise<DocumentInfo> => {
  return request({
    url: '/api/documents',
    method: 'post',
    data
  })
}

/**
 * 更新文档
 * @param id 文档 ID
 * @param data 更新参数
 */
export const updateDocument = async (id: string, data: {
  title?: string
  category?: string
  description?: string
  content?: string
  tags?: string[]
  isPublished?: boolean
}): Promise<DocumentInfo> => {
  return request({
    url: `/api/documents/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除文档
 * @param id 文档 ID
 */
export const deleteDocument = async (id: string): Promise<void> => {
  return request({
    url: `/api/documents/${id}`,
    method: 'delete'
  })
}
