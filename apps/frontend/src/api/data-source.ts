/**
 * 数据源API服务
 */
import type { DataSource, DataSourceQueryParams, DataSourceResponse, TableInfo, ViewInfo, TableStructure, ColumnInfo } from '@/types/metadata/datasource'
import request from '@/utils/request'

/**
 * 获取数据源列表
 * @param params 查询参数
 * @returns 数据源列表
 */
export const getDataSources = async (params?: DataSourceQueryParams): Promise<DataSourceResponse> => {
    return request({
        url: '/api/data-sources',
        method: 'get',
        params
    })
}

/**
 * 获取数据源详情
 * @param id 数据源ID
 * @returns 数据源详情
 */
export const getDataSourceById = async (id: number): Promise<DataSource> => {
    return request({
        url: `/api/data-sources/${id}`,
        method: 'get'
    })
}

/**
 * 创建数据源
 * @param data 数据源数据
 * @returns 创建结果
 */
export const createDataSource = async (data: Partial<DataSource>): Promise<DataSource> => {
    return request({
        url: '/api/data-sources',
        method: 'post',
        data
    })
}

/**
 * 更新数据源
 * @param id 数据源ID
 * @param data 数据源数据
 * @returns 更新结果
 */
export const updateDataSource = async (id: number, data: Partial<DataSource>): Promise<DataSource> => {
    return request({
        url: `/api/data-sources/${id}`,
        method: 'put',
        data
    })
}

/**
 * 删除数据源
 * @param id 数据源ID
 * @returns 删除结果
 */
export const deleteDataSource = async (id: number): Promise<void> => {
    return request({
        url: `/api/data-sources/${id}`,
        method: 'delete'
    })
}

/**
 * 测试数据源连接
 * @param data 数据源连接信息
 * @returns 测试结果
 */
export const testDataSourceConnection = async (data: Partial<DataSource>): Promise<{ success: boolean; message: string }> => {
    return request({
        url: '/api/data-sources/test-connection',
        method: 'post',
        data
    })
}

/**
 * 同步数据源元数据
 * @param id 数据源ID
 * @returns 同步结果
 */
export const syncDataSourceMetadata = async (id: number): Promise<void> => {
    return request({
        url: `/api/data-sources/${id}/sync`,
        method: 'post'
    })
}

/**
 * 获取数据源下的所有表
 * @param id 数据源ID
 * @returns 表列表
 */
export const getDataSourceTables = async (id: number): Promise<TableInfo[]> => {
    return request({
        url: `/api/data-sources/${id}/tables`,
        method: 'get'
    })
}

/**
 * 获取数据源下的所有视图
 * @param id 数据源ID
 * @returns 视图列表
 */
export const getDataSourceViews = async (id: number): Promise<ViewInfo[]> => {
    return request({
        url: `/api/data-sources/${id}/views`,
        method: 'get'
    })
}

/**
 * 获取数据源下的表和视图
 * @param id 数据源ID
 * @returns 表和视图列表
 */
export const getDataSourceTablesAndViews = async (id: number): Promise<{ tables: TableInfo[]; views: ViewInfo[] }> => {
    return request({
        url: `/api/data-sources/${id}/objects`,
        method: 'get'
    })
}

/**
 * 获取表结构信息
 * @param id 数据源ID
 * @param tableName 表名
 * @returns 表结构信息
 */
export const getTableStructure = async (id: number, tableName: string): Promise<TableStructure> => {
    return request({
        url: `/api/data-sources/${id}/tables/${tableName}/structure`,
        method: 'get'
    })
}

/**
 * 获取表字段列表
 * @param id 数据源ID
 * @param tableName 表名
 * @returns 字段列表
 */
export const getTableColumns = async (id: number, tableName: string): Promise<ColumnInfo[]> => {
    return request({
        url: `/api/data-sources/${id}/tables/${tableName}/columns`,
        method: 'get'
    })
}

/**
 * 预览表数据
 * @param id 数据源ID
 * @param tableName 表名
 * @param params 查询参数（page, pageSize, orderBy）
 * @returns 预览数据
 */
export const previewTableData = async (id: number, tableName: string, params?: { page?: number; pageSize?: number; orderBy?: string }): Promise<{ data: any[]; total: number }> => {
    return request({
        url: `/api/data-sources/${id}/tables/${tableName}/preview`,
        method: 'get',
        params
    })
}

/**
 * 执行SQL查询
 * @param id 数据源ID
 * @param data SQL查询参数
 * @returns 查询结果
 */
export const executeSQL = async (id: number, data: { sql: string; limit?: number }): Promise<{ data: any[]; columns: string[]; rows: number; executionTime: number }> => {
    return request({
        url: `/api/data-sources/${id}/query`,
        method: 'post',
        data
    })
}
