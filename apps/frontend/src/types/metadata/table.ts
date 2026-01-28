import type { BaseEntity } from './model-base';

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
