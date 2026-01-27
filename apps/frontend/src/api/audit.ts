import request from '@/utils/request'

export interface AuditLogParams {
    page: number
    pageSize: number
    start_time?: string
    end_time?: string
    user_id?: string
    account?: string
    module?: string
    status?: string
}

export function getLoginLogs(params: AuditLogParams) {
    return request({
        url: '/api/audit/login',
        method: 'get',
        params
    })
}

export function exportLoginLogs(params: Partial<AuditLogParams>) {
    return request({
        url: '/api/audit/login/export',
        method: 'get',
        params,
        responseType: 'blob'
    })
}

export function getOperationLogs(params: AuditLogParams) {
    return request({
        url: '/api/audit/operation',
        method: 'get',
        params
    })
}

export function exportOperationLogs(params: Partial<AuditLogParams>) {
    return request({
        url: '/api/audit/operation/export',
        method: 'get',
        params,
        responseType: 'blob'
    })
}

export function getDataChangeLogs(params: AuditLogParams) {
    return request({
        url: '/api/audit/data',
        method: 'get',
        params
    })
}

export function exportDataChangeLogs(params: Partial<AuditLogParams>) {
    return request({
        url: '/api/audit/data/export',
        method: 'get',
        params,
        responseType: 'blob'
    })
}
