/**
 * 数据源API服务
 */
import request from '@/utils/request'
import type { DataSource, DataSourceQueryParams, DataSourceResponse } from '@/types/data-source'

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
