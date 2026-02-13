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

/**
 * 查询模板 (前端使用)
 */
export interface QueryTemplate {
    /** 模板ID */
    id: number;
    /** 模型ID */
    modelId: number;
    /** 模板名称 */
    templateName: string;
    /** 模板编码 */
    templateCode: string;
    /** 是否默认 */
    isDefault: boolean;
    /** 备注 */
    remark: string;
    /** 创建时间 */
    createdAt: string;
    /** 更新时间 */
    updatedAt: string;
}

/**
 * 查询条件
 */
export interface QueryCondition {
    /** 条件ID */
    id: number;
    /** 模板ID */
    templateId: number;
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

/**
 * 查询模板查询参数
 */
export interface QueryTemplateQueryParams {
    /** 页码 */
    page?: number;
    /** 每页条数 */
    pageSize?: number;
    /** 搜索关键词 */
    search?: string;
    /** 模型ID */
    modelId?: number;
    /** 是否默认 */
    isDefault?: boolean;
}

/**
 * 查询模板响应结果
 */
export interface QueryTemplateResponse {
    /** 数据列表 */
    data: QueryTemplate[];
    /** 总数 */
    total: number;
    /** 页码 */
    page: number;
    /** 每页条数 */
    pageSize: number;
}
