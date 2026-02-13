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

/**
 * 表信息
 */
export interface TableInfo {
  /** 表名 */
  tableName: string;
  /** 表备注 */
  tableComment?: string;
  /** 表类型 */
  tableType?: string;
  /** 记录数 */
  tableRows?: number;
  /** 数据大小 */
  dataLength?: number;
  /** 索引大小 */
  indexLength?: number;
  /** 创建时间 */
  createTime?: string;
  /** 更新时间 */
  updateTime?: string;
}

/**
 * 视图信息
 */
export interface ViewInfo {
  /** 视图名 */
  viewName: string;
  /** 视图备注 */
  viewComment?: string;
  /** 视图定义 */
  definition?: string;
}

/**
 * 表结构信息
 */
export interface TableStructure {
  /** 表名 */
  tableName: string;
  /** 表备注 */
  tableComment?: string;
  /** 引擎 */
  engine?: string;
  /** 字符集 */
  charset?: string;
  /** 字段列表 */
  columns: ColumnInfo[];
}

/**
 * 字段信息
 */
export interface ColumnInfo {
  /** 字段名 */
  columnName: string;
  /** 字段类型 */
  columnType: string;
  /** 数据类型 */
  dataType: string;
  /** 字段备注 */
  columnComment?: string;
  /** 是否主键 */
  isPrimaryKey: boolean;
  /** 是否自增 */
  isAutoIncrement: boolean;
  /** 是否允许空值 */
  isNullable: boolean;
  /** 默认值 */
  defaultValue?: string;
  /** 字符长度 */
  characterMaximumLength?: number;
  /** 数值精度 */
  numericPrecision?: number;
  /** 数值小数位 */
  numericScale?: number;
  /** 排序 */
  ordinalPosition: number;
}
