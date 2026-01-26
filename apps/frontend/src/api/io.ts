/**
 * 数据导入导出API服务
 */
import request from '@/utils/request'

/**
 * 导出数据
 * @param modelName 模型名称
 * @param format 格式: excel/json
 * @param params 查询参数
 * @returns 导出文件流
 */
export const exportData = async (modelName: string, format: 'excel' | 'json', params?: any): Promise<any> => {
    return request({
        url: `/api/metadata/data/${modelName}/export`,
        method: 'post',
        params: { ...params, format },
        responseType: 'blob'
    })
}

/**
 * 获取导入模板
 * @param modelName 模型名称
 * @returns 模板文件流
 */
export const getImportTemplate = async (modelName: string): Promise<any> => {
    return request({
        url: `/api/metadata/data/${modelName}/import-template`,
        method: 'get',
        responseType: 'blob'
    })
}

/**
 * 导入数据
 * @param modelName 模型名称
 * @param file 文件对象
 * @returns 导入结果 (successCount, errors)
 */
export const importData = async (modelName: string, file: File): Promise<any> => {
    const formData = new FormData()
    formData.append('file', file)
    return request({
        url: `/api/metadata/data/${modelName}/import`,
        method: 'post',
        data: formData,
        headers: { 'Content-Type': 'multipart/form-data' }
    })
}
