/**
 * 模型API服务
 */
import type { Model, ModelBuildParams, ModelField, ModelQueryParams, ModelResponse, SQLModelBuildParams } from '@/types/metadata'
import request from '@/utils/request'

/**
 * 获取模型列表
 * @param params 查询参数
 * @returns 模型列表
 */
export const getModels = async (params?: ModelQueryParams): Promise<ModelResponse> => {
    const { pageSize, ...rest } = params || {}
    return request({
        url: '/api/metadata/models',
        method: 'get',
        params: {
            ...rest,
            page_size: pageSize
        }
    })
}

/**
 * 获取所有模型列表 (不分页)
 * @returns 模型列表
 */
export const getAllModels = async (): Promise<Model[]> => {
    return request({
        url: '/api/metadata/models/all',
        method: 'get'
    })
}

/**
 * 自动生成模型编码
 * @returns 生成的编码
 */
export const generateModelCode = async (): Promise<{ code: string }> => {
    return request({
        url: '/api/metadata/utils/generate-model-code-16',
        method: 'get'
    })
}

/**
 * 获取模型详情
 * @param id 模型ID
 * @returns 模型详情
 */
export const getModelById = async (id: string): Promise<Model> => {
    return request({
        url: `/api/metadata/models/${id}`,
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
        url: '/api/metadata/models',
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
export const updateModel = async (id: string, data: Partial<Model>): Promise<Model> => {
    return request({
        url: `/api/metadata/models/${id}`,
        method: 'put',
        data
    })
}

/**
 * 删除模型
 * @param id 模型ID
 * @returns 删除结果
 */
export const deleteModel = async (id: string): Promise<void> => {
    return request({
        url: `/api/metadata/models/${id}`,
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
        url: '/api/metadata/models/build-from-table',
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
        url: '/api/metadata/models/build-from-view',
        method: 'post',
        data: params
    })
}

/**
 * 创建SQL模型
 * @param params 构建参数
 * @returns 构建结果
 */
export const createModelSql = async (params: SQLModelBuildParams): Promise<Model> => {
    return request({
        url: '/api/metadata/models/build-from-sql',
        method: 'post',
        data: params
    })
}

/**
 * 更新SQL模型
 * @param params 构建参数
 * @returns 构建结果
 */
export const updateModelSql = async (params: SQLModelBuildParams): Promise<Model> => {
    return request({
        url: '/api/metadata/models/update-sql',
        method: 'put',
        data: params
    })
}

/**
 * 测试SQL
 * @param params 测试参数
 * @returns 测试结果
 */
export const testSQL = async (params: any): Promise<any> => {
    return request({
        url: '/api/metadata/models/test-sql',
        method: 'post',
        data: params
    })
}

/**
 * 获取模型字段列表
 * @param modelId 模型ID
 * @returns 字段列表
 */
export const getModelFields = async (modelId: string): Promise<ModelField[]> => {
    return request({
        url: `/api/metadata/models/${modelId}/fields`,
        method: 'get'
    })
}

/**
 * 获取模型参数列表
 * @param modelId 模型ID
 * @returns 参数列表
 */
export const getModelParams = async (modelId: string): Promise<any[]> => {
    return request({
        url: `/api/metadata/models/${modelId}/params`,
        method: 'get'
    })
}

/**
 * 获取模型SQL内容
 * @param modelId 模型ID
 * @returns SQL内容
 */
export const getModelSql = async (modelId: string): Promise<any> => {
    return request({
        url: `/api/metadata/models/${modelId}/sql`,
        method: 'get'
    })
}

/**
 * 添加模型字段
 * @param modelId 模型ID
 * @param data 字段数据
 * @returns 添加结果
 */
export const addModelField = async (modelId: string, data: Partial<ModelField>): Promise<ModelField> => {
    return request({
        url: `/api/metadata/models/${modelId}/fields`,
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
export const updateModelField = async (modelId: string, fieldId: string, data: Partial<ModelField>): Promise<ModelField> => {
    return request({
        url: `/api/metadata/models/${modelId}/fields/${fieldId}`,
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
export const deleteModelField = async (modelId: string, fieldId: string): Promise<void> => {
    return request({
        url: `/api/metadata/models/${modelId}/fields/${fieldId}`,
        method: 'delete'
    })
}

/**
 * 预览模型数据
 * @param modelId 模型ID
 * @param params 查询参数
 * @returns 预览数据
 */
export const previewModelData = async (modelId: string, params?: any): Promise<any> => {
    // 兼容 limit (旧参数) 和 pageSize (新参数)
    const { pageSize, limit, ...rest } = params || {}
    const size = pageSize || limit

    return request({
        url: `/api/metadata/models/${modelId}/preview`,
        method: 'post',
        data: {
            ...rest,
            page_size: size
        }
    })
}

/**
 * 验证模型配置
 * @param modelId 模型ID
 * @returns 验证结果
 */
export const validateModel = async (modelId: string): Promise<{ success: boolean; message: string }> => {
    return request({
        url: `/api/metadata/models/${modelId}/validate`,
        method: 'post'
    })
}

// ==================== 增强功能相关API ====================

/**
 * 获取模型字段增强配置
 * @param modelId 模型ID
 * @returns 增强配置列表
 */
export const getFieldEnhancements = async (modelId: string): Promise<any[]> => {
    return request({
        url: `/api/metadata/models/${modelId}/fields/enhancements`,
        method: 'get'
    })
}

/**
 * 批量更新模型字段增强配置
 * @param modelId 模型ID
 * @param data 增强配置数据
 * @returns 更新结果
 */
export const batchUpdateEnhancements = async (modelId: string, data: any[]): Promise<any> => {
    return request({
        url: `/api/metadata/models/${modelId}/fields/batch-enhancements`,
        method: 'post',
        data
    })
}

/**
 * 复制查询模板
 * @param modelId 模型ID
 * @param templateId 模板ID
 * @returns 新模板详情
 */
export const duplicateQueryTemplate = async (modelId: string, templateId: string): Promise<any> => {
    return request({
        url: `/api/metadata/models/${modelId}/query-templates/${templateId}/duplicate`,
        method: 'post'
    })
}

/**
 * 预览查询模板
 * @param modelId 模型ID
 * @param templateId 模板ID
 * @returns 预览信息
 */
export const previewQueryTemplate = async (modelId: string, templateId: string): Promise<any> => {
    return request({
        url: `/api/metadata/models/${modelId}/query-templates/${templateId}/preview`,
        method: 'get'
    })
}

/**
 * 获取模型关联的API列表
 * @param modelId 模型ID
 * @returns API列表
 */
export const getAPIsByModelId = async (modelId: string): Promise<any[]> => {
    return request({
        url: '/api/metadata/apis',
        method: 'get',
        params: { model_id: modelId }
    })
}

/**
 * 为模型批量生成CRUD接口
 * @param modelId 模型ID
 * @returns 生成结果
 */
export const batchGenerateAPIs = async (modelId: string): Promise<any> => {
    return request({
        url: '/api/metadata/apis/batch-generate',
        method: 'post',
        data: { model_id: modelId }
    })
}
/**
 * 统一数据查询 (按ID)
 * @param id 模型ID
 * @param data 查询条件
 */
export const queryDataById = async (id: string, data: any): Promise<any> => {
    return request({
        url: `/api/metadata/models/query/by-id/${id}`,
        method: 'post',
        data
    })
}

/**
 * 统一数据查询 (按代码)
 * @param code 模型代码
 * @param data 查询条件
 */
export const queryDataByCode = async (code: string, data: any): Promise<any> => {
    return request({
        url: `/api/metadata/models/query/by-code/${code}`,
        method: 'post',
        data
    })
}
