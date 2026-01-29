/**
 * 元数据模块API服务
 */
import type { MdConn, MdTable, MdTableField } from '@/types/metadata'
import request from '@/utils/request'

// ==================== 数据连接相关API ====================

/**
 * 获取数据连接列表
 * @returns 数据连接列表
 */
export const getConns = async (): Promise<MdConn[]> => {
    return request({
        url: '/api/metadata/conns',
        method: 'get'
    })
}

/**
 * 根据ID获取数据连接
 * @param id 数据连接ID
 * @returns 数据连接详情
 */
export const getConnById = async (id: string): Promise<MdConn> => {
    return request({
        url: `/api/metadata/conns/${id}`,
        method: 'get'
    })
}

/**
 * 创建数据连接
 * @param data 数据连接数据
 * @returns 创建结果
 */
export const createConn = async (data: Partial<MdConn>): Promise<MdConn> => {
    return request({
        url: '/api/metadata/conns',
        method: 'post',
        data
    })
}

/**
 * 更新数据连接
 * @param id 数据连接ID
 * @param data 数据连接数据
 * @returns 更新结果
 */
export const updateConn = async (id: string, data: Partial<MdConn>): Promise<MdConn> => {
    return request({
        url: `/api/metadata/conns/${id}`,
        method: 'put',
        data
    })
}

/**
 * 删除数据连接
 * @param id 数据连接ID
 * @returns 删除结果
 */
export const deleteConn = async (id: string): Promise<void> => {
    return request({
        url: `/api/metadata/conns/${id}`,
        method: 'delete'
    })
}

/**
 * 根据父ID获取数据连接列表
 * @param parentID 父级ID
 * @returns 数据连接列表
 */
export const getConnsByParentId = async (parentID: string): Promise<MdConn[]> => {
    return request({
        url: `/api/metadata/conns/parent/${parentID}`,
        method: 'get'
    })
}

// ==================== 数据连接表相关API ====================

/**
 * 获取所有表列表
 * @param tenantID 租户ID
 * @returns 表列表
 */
export const getTables = async (): Promise<MdTable[]> => {
    return request({
        url: '/api/metadata/tables',
        method: 'get'
    })
}

/**
 * 根据ID获取表
 * @param id 表ID
 * @returns 表详情
 */
export const getTableById = async (id: string): Promise<MdTable> => {
    return request({
        url: `/api/metadata/tables/${id}`,
        method: 'get'
    })
}

/**
 * 根据连接ID获取表列表
 * @param connID 连接ID
 * @returns 表列表
 */
export const getTablesByConnId = async (connID: string): Promise<MdTable[]> => {
    return request({
        url: `/api/metadata/tables/conn/${connID}`,
        method: 'get'
    })
}

/**
 * 创建表
 * @param data 表数据
 * @returns 创建结果
 */
export const createTable = async (data: Partial<MdTable>): Promise<MdTable> => {
    return request({
        url: '/api/metadata/tables',
        method: 'post',
        data
    })
}

/**
 * 更新表
 * @param id 表ID
 * @param data 表数据
 * @returns 更新结果
 */
export const updateTable = async (id: string, data: Partial<MdTable>): Promise<MdTable> => {
    return request({
        url: `/api/metadata/tables/${id}`,
        method: 'put',
        data
    })
}

/**
 * 删除表
 * @param id 表ID
 * @returns 删除结果
 */
export const deleteTable = async (id: string): Promise<void> => {
    return request({
        url: `/api/metadata/tables/${id}`,
        method: 'delete'
    })
}

// ==================== 数据连接表字段相关API ====================

/**
 * 获取所有字段列表
 * @param tenantID 租户ID
 * @returns 字段列表
 */
export const getFields = async (connID?: string, tableID?: string): Promise<MdTableField[]> => {
    return request({
        url: '/api/metadata/fields',
        method: 'get',
        params: {
            conn_id: connID,
            table_id: tableID
        }
    })
}

/**
 * 根据ID获取字段
 * @param id 字段ID
 * @returns 字段详情
 */
export const getFieldById = async (id: string): Promise<MdTableField> => {
    return request({
        url: `/api/metadata/fields/${id}`,
        method: 'get'
    })
}

/**
 * 根据表ID获取字段列表
 * @param tableID 表ID
 * @returns 字段列表
 */
export const getFieldsByTableId = async (tableID: string): Promise<MdTableField[]> => {
    return request({
        url: `/api/metadata/fields/table/${tableID}`,
        method: 'get'
    })
}

/**
 * 创建字段
 * @param data 字段数据
 * @returns 创建结果
 */
export const createField = async (data: Partial<MdTableField>): Promise<MdTableField> => {
    return request({
        url: '/api/metadata/fields',
        method: 'post',
        data
    })
}

/**
 * 更新字段
 * @param id 字段ID
 * @param data 字段数据
 * @returns 更新结果
 */
export const updateField = async (id: string, data: Partial<MdTableField>): Promise<MdTableField> => {
    return request({
        url: `/api/metadata/fields/${id}`,
        method: 'put',
        data
    })
}

/**
 * 删除字段
 * @param id 字段ID
 * @returns 删除结果
 */
export const deleteField = async (id: string): Promise<void> => {
    return request({
        url: `/api/metadata/fields/${id}`,
        method: 'delete'
    })
}

/**
 * 根据表ID删除所有字段
 * @param tableID 表ID
 * @returns 删除结果
 */
export const deleteFieldsByTableId = async (tableID: string): Promise<void> => {
    return request({
        url: `/api/metadata/fields/table/${tableID}`,
        method: 'delete'
    })
}

// ==================== 增强功能相关API ====================

/**
 * 测试数据连接
 * @param id 数据连接ID
 * @returns 测试结果
 */
export const testConn = async (id: string): Promise<{ success: boolean; message: string }> => {
    return request({
        url: `/api/metadata/conns/${id}/test`,
        method: 'post'
    })
}

/**
 * 测试原始数据连接（创建前测试）
 * @param data 数据连接及其配置
 * @returns 测试结果
 */
export const testRawConn = async (data: Partial<MdConn>): Promise<{ success: boolean; message: string }> => {
    return request({
        url: '/api/metadata/conns/test-raw',
        method: 'post',
        data
    })
}

/**
 * 获取数据库所有模式(Schema)
 * @param id 数据连接ID
 * @returns Schema列表
 */
export const getSchemas = async (id: string): Promise<string[]> => {
    return request({
        url: `/api/metadata/conns/${id}/schemas`,
        method: 'get'
    })
}

/**
 * 获取数据库所有表名称
 * @param id 数据连接ID
 * @returns 表名列表
 */
export const getDBTables = async (id: string, schema?: string): Promise<any> => {
    return request({
        url: `/api/metadata/conns/${id}/tables`,
        method: 'get',
        params: { schema }
    })
}

/**
 * 获取数据库所有视图名称
 * @param id 数据连接ID
 * @returns 视图名列表
 */
export const getDBViews = async (id: string, schema?: string): Promise<any> => {
    return request({
        url: `/api/metadata/conns/${id}/views`,
        method: 'get',
        params: { schema }
    })
}

/**
 * 获取表结构信息
 * @param id 数据连接ID
 * @param tableName 表名
 * @returns 表结构详情
 */
export const getTableStructureFromDB = async (id: string, tableName: string, schema?: string): Promise<any> => {
    return request({
        url: `/api/metadata/conns/${id}/tables/${tableName}/structure`,
        method: 'get',
        params: { schema }
    })
}

/**
 * 预览表数据
 * @param id 数据连接ID
 * @param tableName 表名
 * @param schema Schema
 * @param limit 限制条数
 * @returns 表数据预览
 */
export const previewTableData = async (id: string, tableName: string, schema?: string, limit: number = 100): Promise<any> => {
    return request({
        url: `/api/metadata/conns/${id}/tables/${tableName}/preview`,
        method: 'get',
        params: { schema, limit }
    })
}

/**
 * 获取数据库所有字段（跨表）
 * @param connID 数据连接ID
 * @returns 字段列表
 */
export const getDBFields = async (connID: string): Promise<MdTableField[]> => {
    return request({
        url: '/api/metadata/fields',
        method: 'get',
        params: {
            conn_id: connID
        }
    })
}

