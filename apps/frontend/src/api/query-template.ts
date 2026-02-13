/**
 * 查询模板API服务
 */
import type { QueryTemplate, QueryTemplateQueryParams, QueryTemplateResponse, QueryCondition } from '@/types/metadata/query-template'
import request from '@/utils/request'

/**
 * 获取查询模板列表
 * @param params 查询参数
 * @returns 查询模板列表
 */
export const getQueryTemplates = async (params?: QueryTemplateQueryParams): Promise<QueryTemplateResponse> => {
    return request({
        url: '/api/query-templates',
        method: 'get',
        params
    })
}

/**
 * 获取查询模板详情
 * @param id 模板ID
 * @returns 模板详情
 */
export const getQueryTemplateById = async (id: number): Promise<QueryTemplate> => {
    return request({
        url: `/api/query-templates/${id}`,
        method: 'get'
    })
}

/**
 * 创建查询模板
 * @param data 模板数据
 * @returns 创建结果
 */
export const createQueryTemplate = async (data: Partial<QueryTemplate>): Promise<QueryTemplate> => {
    return request({
        url: '/api/query-templates',
        method: 'post',
        data
    })
}

/**
 * 更新查询模板
 * @param id 模板ID
 * @param data 模板数据
 * @returns 更新结果
 */
export const updateQueryTemplate = async (id: number, data: Partial<QueryTemplate>): Promise<QueryTemplate> => {
    return request({
        url: `/api/query-templates/${id}`,
        method: 'put',
        data
    })
}

/**
 * 删除查询模板
 * @param id 模板ID
 * @returns 删除结果
 */
export const deleteQueryTemplate = async (id: number): Promise<void> => {
    return request({
        url: `/api/query-templates/${id}`,
        method: 'delete'
    })
}

/**
 * 设置默认模板
 * @param id 模板ID
 * @param modelId 模型ID
 * @returns 设置结果
 */
export const setDefaultQueryTemplate = async (id: number, modelId: number): Promise<QueryTemplate> => {
    return request({
        url: `/api/query-templates/${id}/set-default`,
        method: 'post',
        data: { modelId }
    })
}

/**
 * 复制查询模板
 * @param id 模板ID
 * @returns 复制结果
 */
export const duplicateQueryTemplate = async (id: number): Promise<QueryTemplate> => {
    return request({
        url: `/api/query-templates/${id}/duplicate`,
        method: 'post'
    })
}

/**
 * 预览查询模板结果
 * @param id 模板ID
 * @param params 查询参数
 * @returns 预览结果
 */
export const previewQueryTemplate = async (id: number, params?: Record<string, any>): Promise<any> => {
    return request({
        url: `/api/query-templates/${id}/preview`,
        method: 'post',
        data: params
    })
}

/**
 * 获取查询模板的条件列表
 * @param id 模板ID
 * @returns 条件列表
 */
export const getQueryConditions = async (id: number): Promise<QueryCondition[]> => {
    return request({
        url: `/api/query-templates/${id}/conditions`,
        method: 'get'
    })
}

/**
 * 保存查询模板的条件配置
 * @param id 模板ID
 * @param conditions 条件列表
 * @returns 保存结果
 */
export const saveQueryConditions = async (id: number, conditions: QueryCondition[]): Promise<void> => {
    return request({
        url: `/api/query-templates/${id}/conditions`,
        method: 'put',
        data: { conditions }
    })
}
