import { BaseEntity } from './model-base';

/**
 * 操作日志
 */
export interface SysOperationLog extends BaseEntity {
    /** 追踪ID */
    traceID: string;
    /** 用户ID */
    userID: string;
    /** 用户名 */
    userName: string;
    /** 请求方法 */
    method: string;
    /** 请求路径 */
    path: string;
    /** 状态码 */
    status: number;
    /** 耗时 (ms) */
    latency: number;
    /** 客户端IP */
    clientIP: string;
    /** 用户代理 */
    userAgent: string;
}

/**
 * 数据变更日志
 */
export interface SysDataChangeLog extends BaseEntity {
    /** 操作日志ID */
    operationID: string;
    /** 模型ID */
    modelID: string;
    /** 记录ID */
    recordID: string;
    /** 变更动作 (Create/Update/Delete) */
    action: string;
    /** 变更前数据 */
    beforeData: string;
    /** 变更后数据 */
    afterData: string;
}
