import { BaseEntity } from './model-base';

/**
 * 查询模板模型
 */
export interface MdQueryTemplate extends BaseEntity {
    /** 模型ID */
    modelID: string;
    /** 模板名称 */
    templateName: string;
    /** 模板编码 */
    templateCode: string;
    /** 是否默认 */
    isDefault: boolean;
    /** 备注 */
    remark: string;
    /** 查询条件 */
    conditions: MdQueryCondition[];
}

/**
 * 查询模板条件模型
 */
export interface MdQueryCondition extends BaseEntity {
    /** 模板ID */
    templateID: string;
    /** 逻辑运算符 (AND/OR) */
    operator1: string;
    /** 左括号 */
    brackets1: string;
    /** 表模式 */
    tableSchema: string;
    /** 表名称 */
    tableName: string;
    /** 字段名称 */
    columnName: string;
    /** 字段函数 */
    func: string;
    /** 比较运算符 (=, >, LIKE, etc.) */
    operator2: string;
    /** 值1 */
    value1: string;
    /** 值2 */
    value2: string;
    /** 右括号 */
    brackets2: string;
    /** 排序 */
    sort: number;
}
