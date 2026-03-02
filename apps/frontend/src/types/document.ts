/**
 * 文档管理模块类型定义
 */

/**
 * 文档信息接口
 */
export interface DocumentInfo {
  /**
   * 文档 ID
   */
  id: string
  /**
   * 文档标题
   */
  title: string
  /**
   * 文档分类
   */
  category?: string
  /**
   * 文档路径
   */
  path?: string
  /**
   * 文档描述
   */
  description?: string
  /**
   * 创建时间
   */
  createdAt?: string
  /**
   * 更新时间
   */
  updatedAt?: string
  /**
   * 文档大小（字节）
   */
  size?: number
  /**
   * 版本号
   */
  version?: number
  /**
   * 创建者
   */
  createdBy?: string
  /**
   * 更新者
   */
  updatedBy?: string
  /**
   * 是否启用
   */
  isEnabled?: boolean
  /**
   * 排序
   */
  sortOrder?: number
  /**
   * 标签（可选）
   */
  tags?: string[]
  /**
   * 收藏状态（可选）
   */
  isFavorited?: boolean
  /**
   * 阅读进度（可选）
   */
  readProgress?: number
  /**
   * 文档内容
   */
  content?: string
}

/**
 * 文档详情接口
 */
export interface DocumentDetail extends DocumentInfo {
  /**
   * 文档内容（Markdown 格式）
   */
  content: string
  /**
   * 目录结构（可选，由后端生成）
   */
  toc?: any[]
}

/**
 * 文档分类接口
 */
export interface DocumentCategory {
  /**
   * 分类 ID
   */
  id: string
  /**
   * 分类名称
   */
  name: string
  /**
   * 分类描述
   */
  description: string
  /**
   * 该分类下的文档数量
   */
  count: number
  /**
   * 分类图标（可选）
   */
  icon?: string
}

/**
 * 文档查询参数
 */
export interface DocumentQueryParams {
  /**
   * 分类筛选
   */
  category?: string
  /**
   * 关键词搜索
   */
  keyword?: string
  /**
   * 页码
   */
  page?: number
  /**
   * 每页数量
   */
  pageSize?: number
}

/**
 * 文档搜索结果
 */
export interface DocumentSearchResult {
  /**
   * 文档列表
   */
  list: DocumentInfo[]
  /**
   * 总数
   */
  total: number
}
