<template>
    <div class="container-padding">
        <!-- 页面标题区 -->
        <div class="page-header">
            <h1 class="page-title">
                <el-icon class="title-icon">
                    <Monitor />
                </el-icon>
                访问日志
            </h1>
            <div class="header-actions">
                <el-button type="warning" :icon="Download" @click="handleExport">导出</el-button>
            </div>
        </div>
        <!-- 主内容卡片 -->
        <el-card class="main-card">
            <!-- 搜索区域 -->
            <div class="search-area">
                <el-input v-model="queryParams.user_id" placeholder="用户ID" clearable style="width: 120px" @keyup.enter="handleQuery" />
                <el-input v-model="queryParams.path" placeholder="请求路径" clearable style="width: 200px" @keyup.enter="handleQuery" />
                <el-select v-model="queryParams.method" placeholder="请求方法" clearable style="width: 120px" @change="handleQuery">
                    <el-option label="GET" value="GET" />
                    <el-option label="POST" value="POST" />
                    <el-option label="PUT" value="PUT" />
                    <el-option label="DELETE" value="DELETE" />
                </el-select>
                <el-select v-model="queryParams.status" placeholder="状态码" clearable style="width: 100px" @change="handleQuery">
                    <el-option label="2xx 成功" value="2" />
                    <el-option label="3xx 重定向" value="3" />
                    <el-option label="4xx 客户端错误" value="4" />
                    <el-option label="5xx 服务器错误" value="5" />
                </el-select>
                <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" value-format="YYYY-MM-DD" style="width: 300px" @change="handleQuery" />
                <el-button type="primary" :icon="Search" @click="handleQuery">搜索</el-button>
                <el-button :icon="RefreshLeft" @click="resetQuery">重置</el-button>
            </div>
            <!-- 统计卡片 -->
            <div class="stats-cards" v-if="statistics">
                <div class="stat-card">
                    <div class="stat-value">{{ statistics.total || 0 }}</div>
                    <div class="stat-label">总访问量</div>
                </div>
                <div class="stat-card">
                    <div class="stat-value">{{ formatLatency(statistics.avg_latency) }}</div>
                    <div class="stat-label">平均响应时间</div>
                </div>
                <div class="stat-card success">
                    <div class="stat-value">{{ getStatusCount(2) }}</div>
                    <div class="stat-label">成功请求</div>
                </div>
                <div class="stat-card danger">
                    <div class="stat-value">{{ getStatusCount(4) + getStatusCount(5) }}</div>
                    <div class="stat-label">异常请求</div>
                </div>
            </div>
            <!-- 可视化统计区域 -->
            <el-collapse v-model="activeCollapse" class="chart-collapse">
                <el-collapse-item title="可视化统计" name="charts">
                    <div class="charts-container">
                        <div class="chart-item">
                            <div class="chart-title">访问量趋势（按小时）</div>
                            <div ref="trendChartRef" class="chart"></div>
                        </div>
                        <div class="chart-item">
                            <div class="chart-title">状态码分布</div>
                            <div ref="statusChartRef" class="chart"></div>
                        </div>
                        <div class="chart-item wide">
                            <div class="chart-title">TOP 10 访问路径</div>
                            <div ref="pathChartRef" class="chart"></div>
                        </div>
                    </div>
                </el-collapse-item>
            </el-collapse>
            <!-- 表格区域 -->
            <div class="table-area">
                <el-table v-loading="loading" :data="tableData" border style="width: 100%" height="100%">
                    <el-table-column prop="trace_id" label="追踪ID" width="180" align="center" show-overflow-tooltip />
                    <el-table-column prop="method" label="方法" width="80" align="center">
                        <template #default="scope">
                            <el-tag :type="getMethodTag(scope.row.method)">{{ scope.row.method }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="path" label="请求路径" show-overflow-tooltip />
                    <el-table-column prop="user_id" label="用户" width="120" align="center" />
                    <el-table-column prop="client_ip" label="IP" width="140" align="center" />
                    <el-table-column prop="status" label="状态" width="80" align="center">
                        <template #default="scope">
                            <el-tag :type="getStatusTag(scope.row.status)">
                                {{ scope.row.status }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="latency" label="耗时" width="80" align="center">
                        <template #default="scope">
                            <span :class="{ 'slow-request': scope.row.latency > 5000 }">
                                {{ formatLatency(scope.row.latency) }}
                            </span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="create_at" label="访问时间" width="180" align="center">
                        <template #default="scope">
                            {{ formatDate(scope.row.create_at) }}
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="80" fixed="right" align="center">
                        <template #default="scope">
                            <el-button type="primary" link :icon="View" @click="viewDetails(scope.row)">
                                详情
                            </el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
            <!-- 分页区域 -->
            <div class="pagination-area">
                <el-pagination v-model:current-page="queryParams.page" v-model:page-size="queryParams.pageSize" :page-sizes="[10, 20, 50, 100]" background layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
        </el-card>
        <!-- 详情对话框 -->
        <el-dialog v-model="detailsVisible" title="访问日志详情" width="800px" destroy-on-close append-to-body>
            <el-descriptions :column="2" border v-if="currentLog">
                <el-descriptions-item label="追踪ID" :span="2">{{ currentLog.trace_id }}</el-descriptions-item>
                <el-descriptions-item label="用户ID">{{ currentLog.user_id || '-' }}</el-descriptions-item>
                <el-descriptions-item label="租户ID">{{ currentLog.tenant_id || '-' }}</el-descriptions-item>
                <el-descriptions-item label="请求方法">
                    <el-tag :type="getMethodTag(currentLog.method)">{{ currentLog.method }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="状态码">
                    <el-tag :type="getStatusTag(currentLog.status)">{{ currentLog.status }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="请求路径" :span="2">{{ currentLog.path }}</el-descriptions-item>
                <el-descriptions-item label="查询字符串" :span="2" v-if="currentLog.query_string">
                    <pre class="json-content">{{ currentLog.query_string }}</pre>
                </el-descriptions-item>
                <el-descriptions-item label="客户端IP">{{ currentLog.client_ip }}</el-descriptions-item>
                <el-descriptions-item label="响应时间">{{ formatLatency(currentLog.latency) }}</el-descriptions-item>
                <el-descriptions-item label="请求大小">{{ formatSize(currentLog.request_size) }}</el-descriptions-item>
                <el-descriptions-item label="响应大小">{{ formatSize(currentLog.response_size) }}</el-descriptions-item>
                <el-descriptions-item label="来源页面" :span="2">{{ currentLog.referer || '-' }}</el-descriptions-item>
                <el-descriptions-item label="User-Agent" :span="2">
                    <div class="user-agent">{{ currentLog.user_agent }}</div>
                </el-descriptions-item>
                <el-descriptions-item label="地理位置">{{ currentLog.country }} {{ currentLog.province }} {{ currentLog.city }}</el-descriptions-item>
                <el-descriptions-item label="运营商">{{ currentLog.isp || '-' }}</el-descriptions-item>
                <el-descriptions-item label="访问时间" :span="2">{{ formatDate(currentLog.create_at) }}</el-descriptions-item>
            </el-descriptions>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import * as echarts from 'echarts'
import { exportAccessLogs, getAccessLogs, getAccessStatistics } from '@/api/audit'
import { Download, Monitor, RefreshLeft, Search, View } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref, nextTick, onBeforeUnmount } from 'vue'

interface SysAccessLog {
    id: string
    trace_id: string
    user_id: string
    tenant_id?: string
    method: string
    path: string
    query_string?: string
    status: number
    latency: number
    client_ip: string
    user_agent: string
    referer?: string
    request_size?: number
    response_size?: number
    country?: string
    province?: string
    city?: string
    isp?: string
    create_at: string
}

const loading = ref(false)
const total = ref(0)
const tableData = ref<SysAccessLog[]>([])
const dateRange = ref([])
const detailsVisible = ref(false)
const currentLog = ref<SysAccessLog | null>(null)
const statistics = ref<any>(null)
const activeCollapse = ref('charts')

const trendChartRef = ref<HTMLElement>()
const statusChartRef = ref<HTMLElement>()
const pathChartRef = ref<HTMLElement>()

let trendChart: echarts.ECharts | null = null
let statusChart: echarts.ECharts | null = null
let pathChart: echarts.ECharts | null = null

const queryParams = reactive({
    page: 1,
    pageSize: 20,
    user_id: '',
    path: '',
    method: '',
    status: '',
    client_ip: '',
    start_time: '',
    end_time: ''
})

const formatDate = (date: string) => {
    return date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '-'
}

const formatLatency = (ms: number) => {
    if (!ms) return '0ms'
    if (ms < 1000) return `${ms}ms`
    return `${(ms / 1000).toFixed(2)}s`
}

const formatSize = (bytes: number) => {
    if (!bytes) return '-'
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`
    return `${(bytes / 1024 / 1024).toFixed(2)} MB`
}

const getMethodTag = (method: string) => {
    switch (method?.toUpperCase()) {
        case 'GET': return 'info'
        case 'POST': return 'success'
        case 'PUT': return 'warning'
        case 'DELETE': return 'danger'
        default: return ''
    }
}

const getStatusTag = (status: number) => {
    if (status >= 500) return 'danger'
    if (status >= 400) return 'warning'
    if (status >= 300) return 'info'
    if (status >= 200) return 'success'
    return ''
}

const getStatusCount = (category: number) => {
    if (!statistics.value?.status_counts) return 0
    const counts = statistics.value.status_counts as Record<number, number>
    if (category === 2) return (counts[200] || 0) + (counts[201] || 0) + (counts[204] || 0)
    if (category === 4) return Object.keys(counts).filter(k => k.startsWith('4')).reduce((sum, k) => sum + counts[Number(k)], 0)
    if (category === 5) return Object.keys(counts).filter(k => k.startsWith('5')).reduce((sum, k) => sum + counts[Number(k)], 0)
    return 0
}

const getStatistics = async () => {
    try {
        const res: any = await getAccessStatistics({
            start_time: queryParams.start_time,
            end_time: queryParams.end_time
        })
        if (res?.code === 200 || res?.code === 0) {
            statistics.value = res.data
        }
    } catch (error) {
        console.error('获取统计失败:', error)
    }
}

const getList = async () => {
    loading.value = true
    try {
        if (dateRange.value && dateRange.value.length === 2) {
            queryParams.start_time = dateRange.value[0] + ' 00:00:00'
            queryParams.end_time = dateRange.value[1] + ' 23:59:59'
        } else {
            queryParams.start_time = ''
            queryParams.end_time = ''
        }
        const res: any = await getAccessLogs(queryParams)
        const paginatedData = res?.data || {}
        tableData.value = paginatedData.list || []
        total.value = paginatedData.total || 0
    } catch (error) {
        console.error('获取访问日志失败:', error)
        ElMessage.error('获取访问日志失败')
    } finally {
        loading.value = false
    }
}

const handleQuery = () => {
    queryParams.page = 1
    getList()
    getStatistics()
}

const resetQuery = () => {
    queryParams.user_id = ''
    queryParams.path = ''
    queryParams.method = ''
    queryParams.status = ''
    queryParams.client_ip = ''
    dateRange.value = []
    handleQuery()
}

const handleSizeChange = (val: number) => {
    queryParams.pageSize = val
    queryParams.page = 1
    getList()
}

const handleCurrentChange = (val: number) => {
    queryParams.page = val
    getList()
}

const viewDetails = (log: SysAccessLog) => {
    currentLog.value = log
    detailsVisible.value = true
}

const handleExport = async () => {
    try {
        const res = await exportAccessLogs(queryParams)
        const blob = new Blob([JSON.stringify(res, null, 2)], { type: 'application/json' })
        const link = document.createElement('a')
        link.href = URL.createObjectURL(blob)
        link.download = `access_logs_${dayjs().format('YYYYMMDDHHmmss')}.json`
        link.click()
        URL.revokeObjectURL(link.href)
    } catch (error) {
        ElMessage.error('导出失败')
    }
}

onMounted(() => {
    getList()
    getStatistics()
    nextTick(() => {
        initCharts()
    })
})

onBeforeUnmount(() => {
    trendChart?.dispose()
    statusChart?.dispose()
    pathChart?.dispose()
})

const initCharts = () => {
    if (trendChartRef.value) {
        trendChart = echarts.init(trendChartRef.value)
        updateTrendChart()
    }
    if (statusChartRef.value) {
        statusChart = echarts.init(statusChartRef.value)
        updateStatusChart()
    }
    if (pathChartRef.value) {
        pathChart = echarts.init(pathChartRef.value)
        updatePathChart()
    }
}

const updateTrendChart = () => {
    if (!trendChart || !statistics.value?.hourly_counts) return
    const data = statistics.value.hourly_counts as Array<{ Hour: number; Count: number }>
    const hours = data.map(d => `${d.Hour}:00`)
    const counts = data.map(d => d.Count)
    trendChart.setOption({
        tooltip: { trigger: 'axis' },
        xAxis: { type: 'category', data: hours, name: '小时' },
        yAxis: { type: 'value', name: '访问量' },
        series: [{ data: counts, type: 'line', smooth: true, areaStyle: { opacity: 0.3 } }]
    })
}

const updateStatusChart = () => {
    if (!statusChart || !statistics.value?.status_counts) return
    const counts = statistics.value.status_counts as Record<number, number>
    const statusData = Object.entries(counts).map(([status, count]) => ({
        name: status,
        value: count
    }))
    statusChart.setOption({
        tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
        legend: { orient: 'vertical', left: 'left' },
        series: [{
            type: 'pie',
            radius: ['40%', '70%'],
            data: statusData,
            emphasis: { itemStyle: { shadowBlur: 10, shadowColor: 'rgba(0, 0, 0, 0.5)' } }
        }]
    })
}

const updatePathChart = () => {
    if (!pathChart || !statistics.value?.top_paths) return
    const paths = statistics.value.top_paths as Array<{ Path: string; Count: number }>
    const pathNames = paths.map(p => p.Path.length > 30 ? p.Path.substring(0, 30) + '...' : p.Path)
    const counts = paths.map(p => p.Count)
    pathChart.setOption({
        tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
        xAxis: { type: 'value', name: '访问次数' },
        yAxis: { type: 'category', data: pathNames, inverse: true },
        series: [{ type: 'bar', data: counts, itemStyle: { color: '#5470c6' } }]
    })
}
</script>
<style scoped>
.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;
    flex-shrink: 0;
}

.page-title {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 24px;
    font-weight: 600;
    color: #303133;
    margin: 0;
}

.title-icon {
    font-size: 24px;
    color: #409eff;
}

.header-actions {
    display: flex;
    gap: 10px;
}

.main-card {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

:deep(.el-card__body) {
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 20px;
    overflow: hidden;
    box-sizing: border-box;
}

.search-area {
    flex-shrink: 0;
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 10px;
}

.stats-cards {
    display: flex;
    gap: 15px;
    margin-bottom: 20px;
    flex-shrink: 0;
}

.stat-card {
    flex: 1;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 8px;
    padding: 15px 20px;
    color: white;
    text-align: center;
}

.stat-card.success {
    background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.stat-card.danger {
    background: linear-gradient(135deg, #eb3349 0%, #f45c43 100%);
}

.stat-value {
    font-size: 28px;
    font-weight: bold;
    margin-bottom: 5px;
}

.stat-label {
    font-size: 14px;
    opacity: 0.9;
}

.table-area {
    flex: 1;
    overflow: hidden;
    margin-bottom: 20px;
    position: relative;
}

.pagination-area {
    flex-shrink: 0;
    display: flex;
    justify-content: flex-end;
}

.slow-request {
    color: #f56c6c;
    font-weight: 600;
}

.json-content {
    white-space: pre-wrap;
    word-break: break-all;
    background-color: #f5f7fa;
    padding: 8px;
    border-radius: 4px;
    font-size: 12px;
    max-height: 100px;
    overflow-y: auto;
}

.user-agent {
    word-break: break-all;
    font-size: 12px;
    color: #606266;
}

:deep(.el-descriptions__label) {
    width: 120px;
    font-weight: bold;
}

.chart-collapse {
    margin-bottom: 20px;
}

.charts-container {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
}

.chart-item {
    flex: 1;
    min-width: 350px;
    background: #fff;
    border-radius: 8px;
    padding: 15px;
}

.chart-item.wide {
    flex: 100%;
    width: 100%;
}

.chart-title {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 15px;
    text-align: center;
}

.chart {
    width: 100%;
    height: 300px;
}
</style>
