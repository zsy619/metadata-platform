/**
 * 模型查询参数
 */
export interface ModelQueryParams {
  /** 页码 */
  page?: number;
  /** 每页条数 */
  pageSize?: number;
  /** 搜索关键词 */
  search?: string;
  /** 模型类型 */
  kind?: number;
  /** 数据连接ID */
  connID?: number;
}

/**
 * 模型构建参数
 */
export interface ModelBuildParams {
  /** 数据连接ID */
  connID: number;
  /** 表ID */
  tableID: number;
  /** 模型名称 */
  modelName: string;
  /** 模型编码 */
  modelCode: string;
  /** 模型版本 */
  modelVersion?: string;
  /** 模型类型 */
  modelKind: number;
  /** 是否公开 */
  isPublic?: boolean;
  /** 备注 */
  remark?: string;
}