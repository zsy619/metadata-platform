/**
 * 元数据模块类型定义
 */
import { BaseEntity } from './model-base';

/**
 * 数据连接模型
 */
export interface MdConn extends BaseEntity {
    /** 父级ID */
    parent_id: string;
    /** 数据连接名称 */
    conn_name: string;
    /** 数据连接类型（例如MySQL, Oracle, SQLServer, DB2, DM, KingbaseES） */
    conn_kind: string;
    /** 数据库版本（例如8.0, 12c, 2019） */
    conn_version: string;
    /** 数据连接主机地址 */
    conn_host: string;
    /** 数据连接端口号 */
    conn_port: number;
    /** 用户名 */
    conn_user: string;
    /** 密码 */
    conn_password: string;
    /** 数据库 */
    conn_database: string;
    /** 链接地址:自动生成 */
    conn_conn: string;
    /** 删除标识 */
    is_deleted: boolean;
    /** 检测是否有效:-1 未检测 1 有效 0 无效 */
    state: number;
    /** 描述信息 */
    remark: string;
    /** 排序 */
    sort: number;
    /** 创建人id */
    create_id: string;
    /** 创建人 */
    create_by: string;
    /** 修改人id */
    update_id: string;
    /** 修改人 */
    update_by: string;
}

/**
 * 数据连接表模型
 */
export interface MdTable extends BaseEntity {
    /** 数据连接ID */
    conn_id: string;
    /** 数据连接名称 */
    conn_name: string;
    /** 表模式 */
    table_schema: string;
    /** 表名称 */
    table_name: string;
    /** 表标题 */
    table_title: string;
    /** 表类型 */
    table_type: string;
    /** 表描述 */
    table_comment: string;
    /** 删除标识 */
    is_deleted: boolean;
    /** 状态:-1 未检测 1 有效 0 无效 */
    state: number;
    /** 描述信息 */
    remark: string;
    /** 排序 */
    sort: number;
    /** 创建人id */
    create_id: string;
    /** 创建人 */
    create_by: string;
    /** 修改人id */
    update_id: string;
    /** 修改人 */
    update_by: string;
}

/**
 * 数据连接表字段模型
 */
export interface MdTableField extends BaseEntity {
    /** 数据连接ID */
    conn_id: string;
    /** 表ID */
    table_id: string;
    /** 表标题 */
    table_title: string;
    /** 字段名称 */
    column_name: string;
    /** 字段标题 */
    column_title: string;
    /** 数据类型，例如INT、VARCHAR(255)、TIMESTAMP等 */
    column_type: string;
    /** 字段长度 */
    column_length: number;
    /** 字段描述 */
    column_comment: string;
    /** 是否允许为空 */
    is_nullable: boolean;
    /** 是否为主键 */
    is_primary_key: boolean;
    /** 是否自增 */
    is_auto_increment: boolean;
    /** 默认值 */
    default_value: string;
    /** 额外信息（如auto_increment, unique等） */
    extra_info: string;
    /** 删除标识 */
    is_deleted: boolean;
    /** 状态:1 有效 0 无效 */
    state: number;
    /** 描述信息 */
    remark: string;
    /** 排序 */
    sort: number;
    /** 创建人id */
    create_id: string;
    /** 创建人 */
    create_by: string;
    /** 修改人id */
    update_id: string;
    /** 修改人 */
    update_by: string;
}
