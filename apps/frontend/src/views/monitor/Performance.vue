<template>
    <div class="container-padding">
        <!-- 慢查询 TOP 10 -->
        <el-card class="mb-4">
            <template #header>
                <div class="card-header">
                    <span>慢请求 TOP 10</span>
                    <el-button type="primary" size="small" @click="loadSlowQueries" :loading="loading.slowQueries">
                        <el-icon><Refresh /></el-icon>
                        刷新
                    </el-button>
                </div>
            </template>
            <el-table :data="slowQueries" stripe style="width: 100%" v-loading="loading.slowQueries" empty-text="暂无慢请求数据">
                <el-table-column prop="query" label="请求路径" show-overflow-tooltip min-width="300" />
                <el-table-column prop="duration" label="耗时(ms)" width="120" sortable>
                    <template #default="{ row }">
                        <el-tag :type="getDurationTagType(row.duration)">
                            {{ row.duration }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="timestamp" label="发生时间" width="180" />
                <el-table-column prop="user" label="用户" width="120" />
            </el-table>
        </el-card>

        <!-- 错误接口 TOP 10 -->
        <el-card class="mb-4">
            <template #header>
                <div class="card-header">
                    <span>错误接口 TOP 10</span>
                    <el-button type="primary" size="small" @click="loadErrorInterfaces" :loading="loading.errorInterfaces">
                        <el-icon><Refresh /></el-icon>
                        刷新
                    </el-button>
                </div>
            </template>
            <el-table :data="errorInterfaces" stripe style="width: 100%" v-loading="loading.errorInterfaces" empty-text="暂无错误接口数据">
                <el-table-column prop="path" label="接口路径" show-overflow-tooltip min-width="250" />
                <el-table-column prop="method" label="请求方法" width="100">
                    <template #default="{ row }">
                        <el-tag :type="getMethodTagType(row.method)" size="small">
                            {{ row.method }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="error_count" label="错误次数" width="100" sortable>
                    <template #default="{ row }">
                        <el-tag type="danger">{{ row.error_count }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="latest_error" label="最新错误信息" show-overflow-tooltip min-width="200" />
                <el-table-column prop="timestamp" label="最近发生时间" width="180" />
            </el-table>
        </el-card>

        <!-- 访问量TOP接口 -->
        <el-card class="mb-4">
            <template #header>
                <div class="card-header">
                    <span>访问量 TOP 10 接口</span>
                    <el-button type="primary" size="small" @click="loadTopPaths" :loading="loading.topPaths">
                        <el-icon><Refresh /></el-icon>
                        刷新
                    </el-button>
                </div>
            </template>
            <el-table :data="topPaths" stripe style="width: 100%" v-loading="loading.topPaths" empty-text="暂无访问数据">
                <el-table-column type="index" label="排名" width="60" />
                <el-table-column prop="path" label="接口路径" show-overflow-tooltip min-width="400" />
                <el-table-column prop="count" label="访问次数" width="150" sortable>
                    <template #default="{ row }">
                        <div class="progress-cell">
                            <span class="count-value">{{ formatNumber(row.count) }}</span>
                            <el-progress 
                                :percentage="getPercentage(row.count)" 
                                :stroke-width="10"
                                :show-text="false"
                            />
                        </div>
                    </template>
                </el-table-column>
            </el-table>
        </el-card>

        <!-- 响应时间分布 -->
        <el-row :gutter="20">
            <el-col :span="12">
                <el-card>
                    <template #header>
                        <div class="card-header">
                            <span>响应时间分布</span>
                            <el-button type="primary" size="small" @click="loadLatencyDistribution" :loading="loading.latencyDistribution">
                                <el-icon><Refresh /></el-icon>
                                刷新
                            </el-button>
                        </div>
                    </template>
                    <div ref="barChartRef" style="height: 400px;" v-loading="loading.latencyDistribution"></div>
                </el-card>
            </el-col>
            <el-col :span="12">
                <el-card>
                    <template #header>
                        <div class="card-header">
                            <span>按小时访问趋势</span>
                            <el-button type="primary" size="small" @click="loadHourlyTrend" :loading="loading.hourlyTrend">
                                <el-icon><Refresh /></el-icon>
                                刷新
                            </el-button>
                        </div>
                    </template>
                    <div ref="hourlyChartRef" style="height: 400px;" v-loading="loading.hourlyTrend"></div>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<script setup lang="ts">
import {
    getErrorInterfaces,
    getHourlyTrend,
    getLatencyDistribution,
    getSlowQueries,
    getTopPaths,
    type ErrorInterface,
    type HourlyStats,
    type LatencyDistribution,
    type PathStats,
    type SlowQuery
} from '@/api/monitor'
import { Refresh } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import { onMounted, onUnmounted, ref } from 'vue'

// 图表引用
const barChartRef = ref<HTMLElement>()
const hourlyChartRef = ref<HTMLElement>()

// 图表实例
let barChart: echarts.EChartsType | null = null
let hourlyChart: echarts.EChartsType | null = null

// 数据
const slowQueries = ref<SlowQuery[]>([])
const errorInterfaces = ref<ErrorInterface[]>([])
const topPaths = ref<PathStats[]>([])

// 加载状态
const loading = ref({
    slowQueries: false,
    errorInterfaces: false,
    topPaths: false,
    latencyDistribution: false,
    hourlyTrend: false
})

// 最大访问量（用于计算百分比）
const maxCount = ref(0)

// 格式化数字
const formatNumber = (num: number): string => {
    if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
    if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
    return num.toString()
}

// 加载慢查询数据
const loadSlowQueries = async () => {
    loading.value.slowQueries = true
    try {
        const data = await getSlowQueries(10)
        slowQueries.value = data || []
    } catch (error: any) {
        console.error('Failed to load slow queries:', error)
        ElMessage.error(error.message || '加载慢查询数据失败')
        slowQueries.value = []
    } finally {
        loading.value.slowQueries = false
    }
}

// 加载错误接口数据
const loadErrorInterfaces = async () => {
    loading.value.errorInterfaces = true
    try {
        const data = await getErrorInterfaces(10)
        errorInterfaces.value = data || []
    } catch (error: any) {
        console.error('Failed to load error interfaces:', error)
        ElMessage.error(error.message || '加载错误接口数据失败')
        errorInterfaces.value = []
    } finally {
        loading.value.errorInterfaces = false
    }
}

// 加载TOP路径数据
const loadTopPaths = async () => {
    loading.value.topPaths = true
    try {
        const data = await getTopPaths(10)
        topPaths.value = data || []
        if (topPaths.value.length > 0) {
            maxCount.value = topPaths.value[0].count
        } else {
            maxCount.value = 0
        }
    } catch (error: any) {
        console.error('Failed to load top paths:', error)
        ElMessage.error(error.message || '加载TOP路径数据失败')
        topPaths.value = []
        maxCount.value = 0
    } finally {
        loading.value.topPaths = false
    }
}

// 加载响应时间分布
const loadLatencyDistribution = async () => {
    loading.value.latencyDistribution = true
    try {
        const data = await getLatencyDistribution()
        updateBarChart(data || [])
    } catch (error: any) {
        console.error('Failed to load latency distribution:', error)
        ElMessage.error(error.message || '加载响应时间分布数据失败')
        updateBarChart([])
    } finally {
        loading.value.latencyDistribution = false
    }
}

// 加载按小时趋势
const loadHourlyTrend = async () => {
    loading.value.hourlyTrend = true
    try {
        const data = await getHourlyTrend()
        updateHourlyChart(data || [])
    } catch (error: any) {
        console.error('Failed to load hourly trend:', error)
        ElMessage.error(error.message || '加载小时趋势数据失败')
        updateHourlyChart([])
    } finally {
        loading.value.hourlyTrend = false
    }
}

// 初始化柱状图
const initBarChart = () => {
    if (barChartRef.value) {
        barChart = echarts.init(barChartRef.value)
        barChart.setOption({
            tooltip: { 
                trigger: 'axis',
                axisPointer: { type: 'shadow' }
            },
            xAxis: { 
                type: 'category', 
                data: [],
                axisLabel: { rotate: 0 }
            },
            yAxis: { 
                type: 'value', 
                name: '请求数',
                axisLine: { show: true }
            },
            series: [{
                data: [],
                type: 'bar',
                showBackground: true,
                backgroundStyle: { color: 'rgba(180, 180, 180, 0.1)' },
                itemStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        { offset: 0, color: '#409EFF' },
                        { offset: 1, color: '#79bbff' }
                    ])
                },
                label: {
                    show: true,
                    position: 'top',
                    formatter: '{c}'
                }
            }],
            grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true }
        })
    }
}

// 更新柱状图
const updateBarChart = (data: LatencyDistribution[]) => {
    if (!barChart) return
    barChart.setOption({
        xAxis: { data: data.map(d => d.range) },
        series: [{ data: data.map(d => d.count) }]
    })
}

// 初始化小时趋势图
const initHourlyChart = () => {
    if (hourlyChartRef.value) {
        hourlyChart = echarts.init(hourlyChartRef.value)
        hourlyChart.setOption({
            tooltip: { trigger: 'axis' },
            xAxis: { 
                type: 'category', 
                data: [],
                axisLabel: { 
                    formatter: (value: number) => `${value}:00`
                }
            },
            yAxis: { 
                type: 'value', 
                name: '请求数'
            },
            series: [{
                data: [],
                type: 'line',
                smooth: true,
                areaStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        { offset: 0, color: 'rgba(103, 194, 58, 0.5)' },
                        { offset: 1, color: 'rgba(103, 194, 58, 0.1)' }
                    ])
                },
                lineStyle: { color: '#67C23A' },
                itemStyle: { color: '#67C23A' },
                symbol: 'circle',
                symbolSize: 6
            }],
            grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true }
        })
    }
}

// 更新小时趋势图
const updateHourlyChart = (data: HourlyStats[]) => {
    if (!hourlyChart) return
    
    // 补全24小时数据
    const fullData: number[] = []
    const dataMap = new Map(data.map(d => [d.hour, d.count]))
    for (let i = 0; i < 24; i++) {
        fullData.push(dataMap.get(i) || 0)
    }

    hourlyChart.setOption({
        xAxis: { data: Array.from({ length: 24 }, (_, i) => i) },
        series: [{ data: fullData }]
    })
}

// 获取耗时标签类型
const getDurationTagType = (duration: number): 'success' | 'warning' | 'danger' | 'info' => {
    if (duration < 500) return 'success'
    if (duration < 1000) return 'warning'
    if (duration < 3000) return 'danger'
    return 'info'
}

// 获取请求方法标签类型
const getMethodTagType = (method: string): 'primary' | 'success' | 'warning' | 'danger' | 'info' => {
    const methodMap: Record<string, 'primary' | 'success' | 'warning' | 'danger' | 'info'> = {
        'GET': 'primary',
        'POST': 'success',
        'PUT': 'warning',
        'DELETE': 'danger',
        'PATCH': 'info'
    }
    return methodMap[method] || 'info'
}

// 计算百分比
const getPercentage = (count: number): number => {
    if (maxCount.value === 0) return 0
    return Math.round((count / maxCount.value) * 100)
}

const resizeHandler = () => {
    barChart?.resize()
    hourlyChart?.resize()
}

// 加载所有数据
const loadAllData = () => {
    loadSlowQueries()
    loadErrorInterfaces()
    loadTopPaths()
    loadLatencyDistribution()
    loadHourlyTrend()
}

onMounted(() => {
    initBarChart()
    initHourlyChart()
    loadAllData()
    window.addEventListener('resize', resizeHandler)
})

onUnmounted(() => {
    window.removeEventListener('resize', resizeHandler)
    barChart?.dispose()
    hourlyChart?.dispose()
})
</script>

<style scoped>
.mb-4 {
    margin-bottom: 20px;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.progress-cell {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.count-value {
    font-weight: 500;
    color: #303133;
}
</style>
