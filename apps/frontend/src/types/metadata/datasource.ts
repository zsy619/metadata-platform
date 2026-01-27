/**
 * 数据源类型定义
 */
export interface DataSource {
  /** 主键ID */
  id: number;
  /** 数据连接ID */
  connID: number;
  /** 父级ID */
  parentID: number;
  /** 数据源名称 */
  connName: string;
  /** 数据源类型（MySQL, PostgreSQL, SQL Server, Oracle） */
  connKind: string;
  /** 数据库版本 */
  connVersion: string;
  /** 主机地址 */
  connHost: string;
  /** 端口号 */
  connPort: number;
  /** 用户名 */
  connUser: string;
  /** 密码 */
  connPassword: string;
  /** 数据库名称 */
  connDatabase: string;
  /** 连接地址（自动生成） */
  connConn: string;
  /** 删除标识 */
  isDeleted: boolean;
  /** 状态：-1 未检测 1 有效 0 无效 */
  state: number;
  /** 备注 */
  remark: string;
  /** 排序 */
  sort: number;
  /** 创建时间 */
  createdAt: string;
  /** 更新时间 */
  updatedAt: string;
}

/**
 * 数据源查询参数
 */
export interface DataSourceQueryParams {
  /** 页码 */
  page?: number;
  /** 每页条数 */
  pageSize?: number;
  /** 搜索关键词 */
  search?: string;
  /** 数据源类型 */
  type?: string;
}

/**
 * 数据源响应结果
 */
export interface DataSourceResponse {
  /** 数据列表 */
  data: DataSource[];
  /** 总数 */
  total: number;
  /** 页码 */
  page: number;
  /** 每页条数 */
  pageSize: number;
}
