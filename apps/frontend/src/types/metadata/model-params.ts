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
    model_kind?: number;
    /** 数据连接ID */
    conn_id?: string;
}

/**
 * 模型构建参数 (Table/View)
 */
export interface ModelBuildParams {
    /** 数据连接ID */
    connID: number;
    /** 表ID */
    tableID: number;
    /** Schema */
    schema?: string;
    /** 表名 */
    tableName?: string;
    /** 视图名 */
    viewName?: string;
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

export interface SQLParameter {
    name: string
    type: string
    required: boolean
    default: string
}

export interface FieldMapping {
    column_name: string
    show_title: string
    show_width: number
    format: string
}

/**
 * SQL模型构建参数
 */
export interface SQLModelBuildParams {
    model_id?: string
    conn_id: string
    model_name: string
    model_code: string
    sql_content: string
    parameters: SQLParameter[]
    field_mappings: FieldMapping[]
    is_public?: boolean
    remark?: string
}