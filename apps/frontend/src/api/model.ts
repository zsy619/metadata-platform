/**
 * 模型API服务
 */
import request from '@/utils/request'
import type { Model, ModelField, ModelQueryParams, ModelResponse, ModelBuildParams } from '@/types/model'

/**
 * 获取模型列表
 * @param params 查询参数
 * @returns 模型列表
 */
export const getModels = async (params?: ModelQueryParams): Promise<ModelResponse> => {
  return request({
    url: '/api/models',
    method: 'get',
    params
  })
}

/**
 * 获取模型详情
 * @param id 模型ID
 * @returns 模型详情
 */
export const getModelById = async (id: number): Promise<Model> => {
  return request({
    url: `/api/models/${id}`,
    method: 'get'
  })
}

/**
 * 创建模型
 * @param data 模型数据
 * @returns 创建结果
 */
export const createModel = async (data: Partial<Model>): Promise<Model> => {
  return request({
    url: '/api/models',
    method: 'post',
    data
  })
}

/**
 * 更新模型
 * @param id 模型ID
 * @param data 模型数据
 * @returns 更新结果
 */
export const updateModel = async (id: number, data: Partial<Model>): Promise<Model> => {
  return request({
    url: `/api/models/${id}`,
    method: 'put',
    data
  })
}

/**
 * 删除模型
 * @param id 模型ID
 * @returns 删除结果
 */
export const deleteModel = async (id: number): Promise<void> => {
  return request({
    url: `/api/models/${id}`,
    method: 'delete'
  })
}

/**
 * 从表构建模型
 * @param params 构建参数
 * @returns 构建结果
 */
export const buildModelFromTable = async (params: ModelBuildParams): Promise<Model> => {
  return request({
    url: '/api/models/build-from-table',
    method: 'post',
    data: params
  })
}

/**
 * 从视图构建模型
 * @param params 构建参数
 * @returns 构建结果
 */
export const buildModelFromView = async (params: ModelBuildParams): Promise<Model> => {
  return request({
    url: '/api/models/build-from-view',
    method: 'post',
    data: params
  })
}

/**
 * 获取模型字段列表
 * @param modelId 模型ID
 * @returns 字段列表
 */
export const getModelFields = async (modelId: number): Promise<ModelField[]> => {
  return request({
    url: `/api/models/${modelId}/fields`,
    method: 'get'
  })
}

/**
 * 添加模型字段
 * @param modelId 模型ID
 * @param data 字段数据
 * @returns 添加结果
 */
export const addModelField = async (modelId: number, data: Partial<ModelField>): Promise<ModelField> => {
  return request({
    url: `/api/models/${modelId}/fields`,
    method: 'post',
    data
  })
}

/**
 * 更新模型字段
 * @param modelId 模型ID
 * @param fieldId 字段ID
 * @param data 字段数据
 * @returns 更新结果
 */
export const updateModelField = async (modelId: number, fieldId: number, data: Partial<ModelField>): Promise<ModelField> => {
  return request({
    url: `/api/models/${modelId}/fields/${fieldId}`,
    method: 'put',
    data
  })
}

/**
 * 删除模型字段
 * @param modelId 模型ID
 * @param fieldId 字段ID
 * @returns 删除结果
 */
export const deleteModelField = async (modelId: number, fieldId: number): Promise<void> => {
  return request({
    url: `/api/models/${modelId}/fields/${fieldId}`,
    method: 'delete'
  })
}

/**
 * 预览模型数据
 * @param modelId 模型ID
 * @param params 查询参数
 * @returns 预览数据
 */
export const previewModelData = async (modelId: number, params?: any): Promise<any> => {
  return request({
    url: `/api/models/${modelId}/preview`,
    method: 'post',
    data: params
  })
}

/**
 * 验证模型配置
 * @param modelId 模型ID
 * @returns 验证结果
 */
export const validateModel = async (modelId: number): Promise<{ success: boolean; message: string }> => {
  return request({
    url: `/api/models/${modelId}/validate`,
    method: 'post'
  })
}
