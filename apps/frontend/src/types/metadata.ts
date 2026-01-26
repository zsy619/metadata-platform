/**
 * 元数据模块类型定义
 */
import { BaseEntity } from './model-base';

/**
 * 数据连接模型
 */
export interface MdConn extends BaseEntity {
  /** 父级ID */
  parentID: string;
  /** 数据连接名称 */
  connName: string;
  /** 数据连接类型（例如MySQL, Oracle, SQLServer, DB2, DM, KingbaseES） */
  connKind: string;
  /** 数据库版本（例如8.0, 12c, 2019） */
  connVersion: string;
  /** 数据连接主机地址 */
  connHost: string;
  /** 数据连接端口号 */
  connPort: number;
  /** 用户名 */
  connUser: string;
  /** 密码 */
  connPassword: string;
  /** 数据库 */
  connDatabase: string;
  /** 链接地址：自动生成 */
  connConn: string;
  /** 删除标识 */
  isDeleted: boolean;
  /** 检测是否有效:-1 未检测 1 有效 0 无效 */
  state: number;
  /** 描述信息 */
  remark: string;
  /** 排序 */
  sort: number;
  /** 创建人id */
  createID: string;
  /** 创建人 */
  createBy: string;
  /** 修改人id */
  updateID: string;
  /** 修改人 */
  updateBy: string;
}

/**
 * 数据连接表模型
 */
export interface MdTable extends BaseEntity {
  /** 数据连接ID */
  connID: string;
  /** 数据连接名称 */
  connName: string;
  /** 表模式 */
  tableSchema: string;
  /** 表名称 */
  tableName: string;
  /** 表标题 */
  tableTitle: string;
  /** 表类型 */
  tableType: string;
  /** 表描述 */
  tableComment: string;
  /** 删除标识 */
  isDeleted: boolean;
  /** 状态:-1 未检测 1 有效 0 无效 */
  state: number;
  /** 描述信息 */
  remark: string;
  /** 排序 */
  sort: number;
  /** 创建人id */
  createID: string;
  /** 创建人 */
  createBy: string;
  /** 修改人id */
  updateID: string;
  /** 修改人 */
  updateBy: string;
}

/**
 * 数据连接表字段模型
 */
export interface MdTableField extends BaseEntity {
  /** 数据连接ID */
  connID: string;
  /** 表ID */
  tableID: string;
  /** 表标题 */
  tableTitle: string;
  /** 字段名称 */
  columnName: string;
  /** 字段标题 */
  columnTitle: string;
  /** 数据类型，例如INT、VARCHAR(255)、TIMESTAMP等 */
  columnType: string;
  /** 字段长度 */
  columnLength: number;
  /** 字段描述 */
  columnComment: string;
  /** 是否允许为空 */
  isNullable: boolean;
  /** 是否为主键 */
  isPrimaryKey: boolean;
  /** 是否自增 */
  isAutoIncrement: boolean;
  /** 默认值 */
  defaultValue: string;
  /** 额外信息（如auto_increment, unique等） */
  extraInfo: string;
  /** 删除标识 */
  isDeleted: boolean;
  /** 状态:1 有效 0 无效 */
  state: number;
  /** 描述信息 */
  remark: string;
  /** 排序 */
  sort: number;
  /** 创建人id */
  createID: string;
  /** 创建人 */
  createBy: string;
  /** 修改人id */
  updateID: string;
  /** 修改人 */
  updateBy: string;
}
