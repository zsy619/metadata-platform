/**
 * 监控模块API服务
 */
import request from '@/utils/request'

// ==================== 类型定义 ====================

/** 系统信息 */
export interface SystemInfo {
  hostname: string
  os: string
  platform: string
  kernel_version: string
  uptime: number
  uptime_desc: string
}

/** CPU信息 */
export interface CPUInfo {
  count: number
  used_percent: number
  model_name: string
}

/** 内存信息 */
export interface MemoryInfo {
  total: number
  used: number
  available: number
  used_percent: number
  total_gb: number
  used_gb: number
}

/** 磁盘信息 */
export interface DiskInfo {
  total: number
  used: number
  free: number
  used_percent: number
  total_gb: number
  used_gb: number
}

/** 系统负载 */
export interface LoadInfo {
  load1: number
  load5: number
  load15: number
}

/** 运行时信息 */
export interface RuntimeInfo {
  go_version: string
  goroutine_num: number
  heap_alloc: number
  heap_sys: number
  heap_idle: number
  heap_inuse: number
  stack_inuse: number
  stack_sys: number
  num_gc: number
}

/** 系统统计数据 */
export interface SystemStats {
  timestamp: number
  system: SystemInfo
  cpu: CPUInfo
  memory: MemoryInfo
  disk: DiskInfo
  load: LoadInfo
  runtime: RuntimeInfo
}

/** 业务统计数据 */
export interface BusinessStats {
  total_users: number
  active_users: number
  online_users: number
  today_new_users: number
  total_roles: number
  total_orgs: number
  total_apps: number
  total_tenants: number
  total_menus: number
  total_positions: number
  total_role_groups: number
  total_user_groups: number
}

/** 访问统计数据 */
export interface AccessStats {
  total_requests: number
  today_requests: number
  hour_requests: number
  avg_latency: number
  error_count: number
  error_rate: number
  success_rate: number
  total_data_in: number
  total_data_out: number
  qps: number
  peak_qps: number
}

/** 趋势数据 */
export interface TrendData {
  time: string
  value: number
}

/** 按小时统计 */
export interface HourlyStats {
  hour: number
  count: number
}

/** 路径统计 */
export interface PathStats {
  path: string
  count: number
  avg_latency?: number
}

/** 状态码统计 */
export interface StatusStats {
  status: number
  count: number
}

/** 慢查询 */
export interface SlowQuery {
  query: string
  duration: number
  timestamp: string
  user: string
}

/** 错误接口 */
export interface ErrorInterface {
  path: string
  method: string
  error_count: number
  latest_error: string
  timestamp: string
}

/** 响应时间分布 */
export interface LatencyDistribution {
  range: string
  count: number
}

/** 登录统计 */
export interface LoginStats {
  today_logins: number
  today_success: number
  today_failed: number
  online_count: number
  login_trend: TrendData[]
}

/** 仪表盘汇总数据 */
export interface DashboardSummary {
  system: SystemStats
  business: BusinessStats
  access: AccessStats
  login: LoginStats
  daily_trend: TrendData[]
  status_distribution: StatusStats[]
}

/** API响应包装 */
interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

// ==================== API 接口 ====================

/**
 * 获取仪表盘汇总数据
 * @returns 仪表盘汇总数据
 */
export const getDashboardSummary = async (): Promise<DashboardSummary> => {
  const res: ApiResponse<DashboardSummary> = await request({
    url: '/api/monitor/dashboard',
    method: 'get'
  })
  return res.data
}

/**
 * 获取系统统计数据
 * @returns 系统统计数据
 */
export const getSystemStats = async (): Promise<SystemStats> => {
  const res: ApiResponse<SystemStats> = await request({
    url: '/api/monitor/system',
    method: 'get'
  })
  return res.data
}

/**
 * 获取业务统计数据
 * @returns 业务统计数据
 */
export const getBusinessStats = async (): Promise<BusinessStats> => {
  const res: ApiResponse<BusinessStats> = await request({
    url: '/api/monitor/business',
    method: 'get'
  })
  return res.data
}

/**
 * 获取访问统计数据
 * @returns 访问统计数据
 */
export const getAccessStats = async (): Promise<AccessStats> => {
  const res: ApiResponse<AccessStats> = await request({
    url: '/api/monitor/access',
    method: 'get'
  })
  return res.data
}

/**
 * 获取按小时趋势数据
 * @param date 日期 (YYYY-MM-DD)
 * @returns 按小时统计数据
 */
export const getHourlyTrend = async (date?: string): Promise<HourlyStats[]> => {
  const params: Record<string, string> = {}
  if (date) {
    params.date = date
  }
  const res: ApiResponse<HourlyStats[]> = await request({
    url: '/api/monitor/trend/hourly',
    method: 'get',
    params
  })
  return res.data
}

/**
 * 获取按日趋势数据（最近7天）
 * @returns 按日趋势数据
 */
export const getDailyTrend = async (): Promise<TrendData[]> => {
  const res: ApiResponse<TrendData[]> = await request({
    url: '/api/monitor/trend/daily',
    method: 'get'
  })
  return res.data
}

/**
 * 获取访问量TOP路径
 * @param limit 返回数量限制
 * @returns TOP路径列表
 */
export const getTopPaths = async (limit: number = 10): Promise<PathStats[]> => {
  const res: ApiResponse<PathStats[]> = await request({
    url: '/api/monitor/top-paths',
    method: 'get',
    params: { limit }
  })
  return res.data
}

/**
 * 获取状态码分布
 * @returns 状态码分布数据
 */
export const getStatusDistribution = async (): Promise<StatusStats[]> => {
  const res: ApiResponse<StatusStats[]> = await request({
    url: '/api/monitor/status-distribution',
    method: 'get'
  })
  return res.data
}

/**
 * 获取慢查询列表
 * @param limit 返回数量限制
 * @returns 慢查询列表
 */
export const getSlowQueries = async (limit: number = 10): Promise<SlowQuery[]> => {
  const res: ApiResponse<SlowQuery[]> = await request({
    url: '/api/monitor/slow-queries',
    method: 'get',
    params: { limit }
  })
  return res.data
}

/**
 * 获取错误接口列表
 * @param limit 返回数量限制
 * @returns 错误接口列表
 */
export const getErrorInterfaces = async (limit: number = 10): Promise<ErrorInterface[]> => {
  const res: ApiResponse<ErrorInterface[]> = await request({
    url: '/api/monitor/error-interfaces',
    method: 'get',
    params: { limit }
  })
  return res.data
}

/**
 * 获取响应时间分布
 * @returns 响应时间分布数据
 */
export const getLatencyDistribution = async (): Promise<LatencyDistribution[]> => {
  const res: ApiResponse<LatencyDistribution[]> = await request({
    url: '/api/monitor/latency-distribution',
    method: 'get'
  })
  return res.data
}

/**
 * 获取登录统计
 * @returns 登录统计数据
 */
export const getLoginStats = async (): Promise<LoginStats> => {
  const res: ApiResponse<LoginStats> = await request({
    url: '/api/monitor/login-stats',
    method: 'get'
  })
  return res.data
}
