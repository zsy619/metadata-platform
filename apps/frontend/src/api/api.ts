/**
 * API接口管理服务
 */
import request from '@/utils/request'
import type { API, APIQueryParams, APIResponse, APITestResult } from '@/types/api'

/**
 * 获取接口列表
 * @param params 查询参数
 * @returns 接口列表
 */
export const getAPIs = async (params?: APIQueryParams): Promise<APIResponse> => {
  return request({
    url: '/api/apis',
    method: 'get',
    params
  })
}

/**
 * 获取接口详情
 * @param id 接口ID
 * @returns 接口详情
 */
export const getAPIById = async (id: number): Promise<API> => {
  return request({
    url: `/api/apis/${id}`,
    method: 'get'
  })
}

/**
 * 创建接口
 * @param data 接口数据
 * @returns 创建结果
 */
export const createAPI = async (data: Partial<API>): Promise<API> => {
  return request({
    url: '/api/apis',
    method: 'post',
    data
  })
}

/**
 * 更新接口
 * @param id 接口ID
 * @param data 接口数据
 * @returns 更新结果
 */
export const updateAPI = async (id: number, data: Partial<API>): Promise<API> => {
  return request({
    url: `/api/apis/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除接口
 * @param id 接口ID
 * @returns 删除结果
 */
export const deleteAPI = async (id: number): Promise<void> => {
  return request({
    url: `/api/apis/${id}`,
    method: 'delete'
  })
}

/**
 * 启用接口
 * @param id 接口ID
 * @returns 启用结果
 */
export const enableAPI = async (id: number): Promise<API> => {
  return request({
    url: `/api/apis/${id}/enable`,
    method: 'post'
  })
}

/**
 * 禁用接口
 * @param id 接口ID
 * @returns 禁用结果
 */
export const disableAPI = async (id: number): Promise<API> => {
  return request({
    url: `/api/apis/${id}/disable`,
    method: 'post'
  })
}

/**
 * 测试接口
 * @param id 接口ID
 * @param data 测试数据
 * @returns 测试结果
 */
export const testAPI = async (id: number, data?: any): Promise<APITestResult> => {
  return request({
    url: `/api/apis/${id}/test`,
    method: 'post',
    data
  })
}

/**
 * 批量生成接口
 * @param modelId 模型ID
 * @param types 接口类型列表
 * @returns 生成结果
 */
export const batchGenerateAPIs = async (modelId: number, types: number[]): Promise<API[]> => {
  return request({
    url: '/api/apis/batch-generate',
    method: 'post',
    data: {
      modelId,
      types
    }
  })
}

/**
 * 获取接口文档
 * @param id 接口ID
 * @returns 接口文档
 */
export const getAPIDocument = async (id: number): Promise<string> => {
  return request({
    url: `/api/apis/${id}/document`,
    method: 'get'
  })
}

/**
 * 配置接口Mock
 * @param id 接口ID
 * @param data Mock配置
 * @returns 配置结果
 */
export const configAPIMock = async (id: number, data: any): Promise<any> => {
  return request({
    url: `/api/apis/${id}/mock`,
    method: 'post',
    data
  })
}
