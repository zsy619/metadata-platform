import type { BaseEntity } from './model-base';

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
