/**
 * 审计日志API服务
 */
import request from '@/utils/request'

/**
 * 获取操作日志列表
 * @param params 查询参数 (trace_id, user_id, status, page, pageSize)
 * @returns 日志列表
 */
export const getOperationLogs = async (params?: any): Promise<any> => {
    return request({
        url: '/api/metadata/logs/operations',
        method: 'get',
        params
    })
}

/**
 * 获取数据变更日志列表
 * @param params 查询参数 (operation_id, model_id, record_id, action, page, pageSize)
 * @returns 变更日志列表
 */
export const getDataChangeLogs = async (params?: any): Promise<any> => {
    return request({
        url: '/api/metadata/logs/data-changes',
        method: 'get',
        params
    })
}
