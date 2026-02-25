<template>
    <div class="container-padding">
        <!-- 访问统计卡片 -->
        <el-row :gutter="20" class="mb-4">
            <el-col :span="6" v-for="item in accessCards" :key="item.title">
                <el-card shadow="hover" v-loading="loading.dashboard">
                    <template #header>
                        <div class="card-header">
                            <span>{{ item.title }}</span>
                            <el-tag :type="item.type" effect="plain">{{ item.tag }}</el-tag>
                        </div>
                    </template>
                    <div class="card-value">{{ item.value }}</div>
                    <div class="card-desc">{{ item.desc }}</div>
                </el-card>
            </el-col>
        </el-row>

        <!-- 用户统计卡片 -->
        <el-row :gutter="20" class="mb-4">
            <el-col :span="6" v-for="item in userCards" :key="item.title">
                <el-card shadow="hover" v-loading="loading.dashboard">
                    <template #header>
                        <div class="card-header">
                            <span>{{ item.title }}</span>
                            <el-icon :size="20" :color="item.color"><component :is="item.icon" /></el-icon>
                        </div>
                    </template>
                    <div class="card-value">{{ item.value }}</div>
                    <div class="card-desc">{{ item.desc }}</div>
                </el-card>
            </el-col>
        </el-row>

        <!-- 业务统计卡片 -->
        <el-row :gutter="20" class="mb-4">
            <el-col :span="4" v-for="item in businessCards" :key="item.title">
                <el-card shadow="hover" class="business-card" v-loading="loading.dashboard">
                    <div class="business-icon" :style="{ backgroundColor: item.color }">
                        <el-icon :size="24"><component :is="item.icon" /></el-icon>
                    </div>
                    <div class="business-info">
                        <div class="business-value">{{ item.value }}</div>
                        <div class="business-title">{{ item.title }}</div>
                    </div>
                </el-card>
            </el-col>
        </el-row>

        <!-- 图表区域 -->
        <el-row :gutter="20">
            <el-col :span="16">
                <el-card title="流量趋势">
                    <template #extra>
                        <el-button-group size="small">
                            <el-button :type="trendType === 'daily' ? 'primary' : ''" @click="trendType = 'daily'">日</el-button>
                            <el-button :type="trendType === 'hourly' ? 'primary' : ''" @click="trendType = 'hourly'">时</el-button>
                        </el-button-group>
                    </template>
                    <div ref="lineChartRef" style="height: 350px;" v-loading="loading.dashboard"></div>
                </el-card>
            </el-col>
            <el-col :span="8">
                <el-card title="状态码分布">
                    <div ref="pieChartRef" style="height: 350px;" v-loading="loading.dashboard"></div>
                </el-card>
            </el-col>
        </el-row>

        <!-- 系统资源 -->
        <el-row :gutter="20" class="mt-4">
            <el-col :span="8">
                <el-card>
                    <template #header>
                        <div class="card-header">
                            <span>CPU使用率</span>
                            <el-tag v-if="systemStats?.cpu" :type="getTagTypeByPercent(systemStats.cpu.used_percent)">
                                {{ systemStats.cpu.used_percent.toFixed(1) }}%
                            </el-tag>
                        </div>
                    </template>
                    <div ref="cpuGaugeRef" style="height: 200px;" v-loading="loading.system"></div>
                </el-card>
            </el-col>
            <el-col :span="8">
                <el-card>
                    <template #header>
                        <div class="card-header">
                            <span>内存使用率</span>
                            <el-tag v-if="systemStats?.memory" :type="getTagTypeByPercent(systemStats.memory.used_percent)">
                                {{ systemStats.memory.used_gb.toFixed(1) }} / {{ systemStats.memory.total_gb.toFixed(1) }} GB
                            </el-tag>
                        </div>
                    </template>
                    <div ref="memoryGaugeRef" style="height: 200px;" v-loading="loading.system"></div>
                </el-card>
            </el-col>
            <el-col :span="8">
                <el-card>
                    <template #header>
                        <div class="card-header">
                            <span>磁盘使用率</span>
                            <el-tag v-if="systemStats?.disk" :type="getTagTypeByPercent(systemStats.disk.used_percent)">
                                {{ systemStats.disk.used_gb.toFixed(1) }} / {{ systemStats.disk.total_gb.toFixed(1) }} GB
                            </el-tag>
                        </div>
                    </template>
                    <div ref="diskGaugeRef" style="height: 200px;" v-loading="loading.system"></div>
                </el-card>
            </el-col>
        </el-row>

        <!-- 系统信息 -->
        <el-row :gutter="20" class="mt-4" v-if="systemStats?.system">
            <el-col :span="24">
                <el-card title="系统信息">
                    <el-descriptions :column="6" border>
                        <el-descriptions-item label="主机名">{{ systemStats.system.hostname }}</el-descriptions-item>
                        <el-descriptions-item label="操作系统">{{ systemStats.system.os }}</el-descriptions-item>
                        <el-descriptions-item label="平台">{{ systemStats.system.platform }}</el-descriptions-item>
                        <el-descriptions-item label="内核版本">{{ systemStats.system.kernel_version }}</el-descriptions-item>
                        <el-descriptions-item label="运行时间">{{ systemStats.system.uptime_desc }}</el-descriptions-item>
                        <el-descriptions-item label="Go版本">{{ systemStats.runtime?.go_version }}</el-descriptions-item>
                    </el-descriptions>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<script setup lang="ts">
import {
    getDashboardSummary,
    getHourlyTrend,
    getSystemStats,
    type DashboardSummary,
    type HourlyStats,
    type SystemStats,
    type TrendData
} from '@/api/monitor'
import { Briefcase, CircleCheck, Document, Key, Menu, OfficeBuilding, Timer, User, UserFilled } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import { markRaw, onMounted, onUnmounted, ref, shallowRef, watch } from 'vue'

// 图表引用
const lineChartRef = ref<HTMLElement>()
const pieChartRef = ref<HTMLElement>()
const cpuGaugeRef = ref<HTMLElement>()
const memoryGaugeRef = ref<HTMLElement>()
const diskGaugeRef = ref<HTMLElement>()

// 图表实例
let lineChart: echarts.EChartsType | null = null
let pieChart: echarts.EChartsType | null = null
let cpuGauge: echarts.EChartsType | null = null
let memoryGauge: echarts.EChartsType | null = null
let diskGauge: echarts.EChartsType | null = null

// 趋势类型
const trendType = ref<'daily' | 'hourly'>('daily')

// 加载状态
const loading = ref({
    dashboard: false,
    system: false
})

// 访问统计卡片
const accessCards = ref([
    { title: '总请求数', value: '-', desc: '历史累计', type: 'primary' as const, tag: 'Total' },
    { title: '今日请求', value: '-', desc: '今日累计', type: 'success' as const, tag: 'Today' },
    { title: '当前QPS', value: '-', desc: '实时统计', type: 'warning' as const, tag: 'QPS' },
    { title: '错误率', value: '-', desc: '历史统计', type: 'danger' as const, tag: 'Error' }
])

// 用户统计卡片
const userCards = ref([
    { title: '总用户数', value: '-', desc: '注册用户总数', icon: markRaw(User), color: '#409EFF' },
    { title: '活跃用户', value: '-', desc: '最近7天登录', icon: markRaw(CircleCheck), color: '#67C23A' },
    { title: '在线用户', value: '-', desc: '30分钟内活跃', icon: markRaw(UserFilled), color: '#E6A23C' },
    { title: '今日新增', value: '-', desc: '今日注册用户', icon: markRaw(Timer), color: '#F56C6C' }
])

// 业务统计卡片
const businessCards = ref([
    { title: '角色数', value: '-', icon: markRaw(Key), color: '#409EFF' },
    { title: '组织数', value: '-', icon: markRaw(OfficeBuilding), color: '#67C23A' },
    { title: '应用数', value: '-', icon: markRaw(Document), color: '#E6A23C' },
    { title: '菜单数', value: '-', icon: markRaw(Menu), color: '#F56C6C' },
    { title: '职位数', value: '-', icon: markRaw(Briefcase), color: '#00d4aa' },
    { title: '租户数', value: '-', icon: markRaw(OfficeBuilding), color: '#9b59b6' }
])

// 系统统计数据
const systemStats = shallowRef<SystemStats | null>(null)

// 趋势数据缓存
const trendData = ref<{
    daily: TrendData[]
    hourly: HourlyStats[]
}>({
    daily: [],
    hourly: []
})

// 加载仪表盘数据
const loadDashboardData = async () => {
    loading.value.dashboard = true
    try {
        const data: DashboardSummary = await getDashboardSummary()
        
        // 更新访问统计卡片
        if (data?.access) {
            accessCards.value[0].value = formatNumber(data.access.total_requests)
            accessCards.value[1].value = formatNumber(data.access.today_requests)
            accessCards.value[2].value = data.access.qps.toFixed(1)
            accessCards.value[3].value = `${data.access.error_rate.toFixed(2)}%`
        }

        // 更新用户统计卡片
        if (data?.business) {
            userCards.value[0].value = formatNumber(data.business.total_users)
            userCards.value[1].value = formatNumber(data.business.active_users)
            userCards.value[2].value = formatNumber(data.business.online_users)
            userCards.value[3].value = formatNumber(data.business.today_new_users)
            
            // 更新业务统计卡片
            businessCards.value[0].value = formatNumber(data.business.total_roles)
            businessCards.value[1].value = formatNumber(data.business.total_orgs)
            businessCards.value[2].value = formatNumber(data.business.total_apps)
            businessCards.value[3].value = formatNumber(data.business.total_menus)
            businessCards.value[4].value = formatNumber(data.business.total_positions)
            businessCards.value[5].value = formatNumber(data.business.total_tenants)
        }

        // 更新系统统计
        if (data?.system) {
            systemStats.value = data.system
            updateGauges(data.system)
        }

        // 缓存趋势数据
        if (data?.daily_trend && data.daily_trend.length > 0) {
            trendData.value.daily = data.daily_trend
            if (trendType.value === 'daily') {
                updateLineChart(data.daily_trend)
            }
        }

        // 更新状态码分布图
        if (data?.status_distribution && data.status_distribution.length > 0) {
            updatePieChart(data.status_distribution)
        }
    } catch (error: any) {
        console.error('Failed to load dashboard data:', error)
        ElMessage.error(error.message || '加载仪表盘数据失败')
    } finally {
        loading.value.dashboard = false
    }
}

// 加载小时趋势数据
const loadHourlyTrendData = async () => {
    try {
        const data = await getHourlyTrend()
        trendData.value.hourly = data || []
        if (trendType.value === 'hourly' && data && data.length > 0) {
            updateLineChart(data.map(d => ({ time: `${d.hour}:00`, value: d.count })))
        }
    } catch (error: any) {
        console.error('Failed to load hourly trend:', error)
        trendData.value.hourly = []
    }
}

// 监听趋势类型变化
watch(trendType, (newType) => {
    if (newType === 'daily' && trendData.value.daily.length > 0) {
        updateLineChart(trendData.value.daily)
    } else if (newType === 'hourly') {
        if (trendData.value.hourly.length > 0) {
            updateLineChart(trendData.value.hourly.map(d => ({ time: `${d.hour}:00`, value: d.count })))
        } else {
            loadHourlyTrendData()
        }
    }
})

// 格式化数字
const formatNumber = (num: number): string => {
    if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
    if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
    return num.toString()
}

// 根据百分比获取标签类型
const getTagTypeByPercent = (percent: number): 'success' | 'warning' | 'danger' | 'info' => {
    if (percent < 50) return 'success'
    if (percent < 80) return 'warning'
    return 'danger'
}

// 初始化折线图
const initLineChart = () => {
    if (lineChartRef.value) {
        lineChart = echarts.init(lineChartRef.value)
        lineChart.setOption({
            tooltip: { trigger: 'axis' },
            xAxis: { 
                type: 'category', 
                data: [],
                axisLabel: { rotate: 45 }
            },
            yAxis: { type: 'value', name: '请求数' },
            series: [{
                data: [],
                type: 'line',
                smooth: true,
                areaStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        { offset: 0, color: 'rgba(64, 158, 255, 0.5)' },
                        { offset: 1, color: 'rgba(64, 158, 255, 0.1)' }
                    ])
                },
                lineStyle: { color: '#409EFF' },
                itemStyle: { color: '#409EFF' }
            }],
            grid: { left: '3%', right: '4%', bottom: '15%', containLabel: true }
        })
    }
}

// 更新折线图
const updateLineChart = (data: { time: string; value: number }[]) => {
    if (!lineChart) return
    lineChart.setOption({
        xAxis: { data: data.map(d => d.time) },
        series: [{ data: data.map(d => d.value) }]
    })
}

// 初始化饼图
const initPieChart = () => {
    if (pieChartRef.value) {
        pieChart = echarts.init(pieChartRef.value)
        pieChart.setOption({
            tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
            legend: { orient: 'vertical', left: 'left' },
            series: [{
                type: 'pie',
                radius: ['40%', '70%'],
                center: ['60%', '50%'],
                data: [],
                label: { show: false },
                emphasis: {
                    label: { show: true, fontSize: 14, fontWeight: 'bold' }
                }
            }]
        })
    }
}

// 更新饼图
const updatePieChart = (data: { status: number; count: number }[]) => {
    if (!pieChart) return
    
    const statusMap: Record<string, { name: string; color: string }> = {
        '2xx': { name: '2xx 成功', color: '#67C23A' },
        '3xx': { name: '3xx 重定向', color: '#409EFF' },
        '4xx': { name: '4xx 客户端错误', color: '#E6A23C' },
        '5xx': { name: '5xx 服务端错误', color: '#F56C6C' }
    }

    const groupedData: Record<string, number> = { '2xx': 0, '3xx': 0, '4xx': 0, '5xx': 0 }

    data.forEach(item => {
        if (item.status >= 200 && item.status < 300) groupedData['2xx'] += item.count
        else if (item.status >= 300 && item.status < 400) groupedData['3xx'] += item.count
        else if (item.status >= 400 && item.status < 500) groupedData['4xx'] += item.count
        else if (item.status >= 500) groupedData['5xx'] += item.count
    })

    const pieData = Object.entries(groupedData)
        .filter(([_, count]) => count > 0)
        .map(([key, count]) => ({
            name: statusMap[key].name,
            value: count,
            itemStyle: { color: statusMap[key].color }
        }))

    pieChart.setOption({ series: [{ data: pieData }] })
}

// 初始化仪表盘
const initGauges = () => {
    const gaugeOption = {
        series: [{
            type: 'gauge',
            center: ['50%', '60%'],
            radius: '80%',
            startAngle: 200,
            endAngle: -20,
            min: 0,
            max: 100,
            splitNumber: 10,
            itemStyle: { color: '#409EFF' },
            progress: { show: true, width: 20 },
            pointer: { show: false },
            axisLine: { lineStyle: { width: 20, color: [[1, '#E5E5E5']] } },
            axisTick: { show: false },
            splitLine: { show: false },
            axisLabel: { show: false },
            title: { show: false },
            detail: {
                valueAnimation: true,
                width: '60%',
                lineHeight: 30,
                borderRadius: 8,
                offsetCenter: [0, '10%'],
                fontSize: 24,
                fontWeight: 'bold',
                formatter: '{value}%',
                color: 'inherit'
            },
            data: [{ value: 0 }]
        }]
    }

    if (cpuGaugeRef.value) {
        cpuGauge = echarts.init(cpuGaugeRef.value)
        cpuGauge.setOption(gaugeOption)
    }
    if (memoryGaugeRef.value) {
        memoryGauge = echarts.init(memoryGaugeRef.value)
        memoryGauge.setOption(gaugeOption)
    }
    if (diskGaugeRef.value) {
        diskGauge = echarts.init(diskGaugeRef.value)
        diskGauge.setOption(gaugeOption)
    }
}

// 更新仪表盘
const updateGauges = (stats: SystemStats) => {
    if (cpuGauge && stats.cpu) {
        cpuGauge.setOption({
            series: [{ 
                data: [{ value: stats.cpu.used_percent.toFixed(1) }],
                itemStyle: { color: getColorByPercent(stats.cpu.used_percent) }
            }]
        })
    }
    if (memoryGauge && stats.memory) {
        memoryGauge.setOption({
            series: [{ 
                data: [{ value: stats.memory.used_percent.toFixed(1) }],
                itemStyle: { color: getColorByPercent(stats.memory.used_percent) }
            }]
        })
    }
    if (diskGauge && stats.disk) {
        diskGauge.setOption({
            series: [{ 
                data: [{ value: stats.disk.used_percent.toFixed(1) }],
                itemStyle: { color: getColorByPercent(stats.disk.used_percent) }
            }]
        })
    }
}

// 根据百分比获取颜色
const getColorByPercent = (percent: number): string => {
    if (percent < 50) return '#67C23A'
    if (percent < 80) return '#E6A23C'
    return '#F56C6C'
}

// 刷新系统统计
const refreshSystemStats = async () => {
    loading.value.system = true
    try {
        const stats = await getSystemStats()
        systemStats.value = stats
        updateGauges(stats)
    } catch (error: any) {
        console.error('Failed to refresh system stats:', error)
    } finally {
        loading.value.system = false
    }
}

// 定时刷新
let refreshTimer: number | null = null

const startRefresh = () => {
    refreshTimer = window.setInterval(refreshSystemStats, 30000)
}

const resizeHandler = () => {
    lineChart?.resize()
    pieChart?.resize()
    cpuGauge?.resize()
    memoryGauge?.resize()
    diskGauge?.resize()
}

onMounted(() => {
    initLineChart()
    initPieChart()
    initGauges()
    loadDashboardData()
    startRefresh()
    window.addEventListener('resize', resizeHandler)
})

onUnmounted(() => {
    if (refreshTimer) clearInterval(refreshTimer)
    window.removeEventListener('resize', resizeHandler)
    lineChart?.dispose()
    pieChart?.dispose()
    cpuGauge?.dispose()
    memoryGauge?.dispose()
    diskGauge?.dispose()
})
</script>

<style scoped>
.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.card-value {
    font-size: 24px;
    font-weight: bold;
    margin: 10px 0;
}

.card-desc {
    font-size: 14px;
    color: #909399;
}

.mb-4 {
    margin-bottom: 20px;
}

.mt-4 {
    margin-top: 20px;
}

.business-card :deep(.el-card__body) {
    display: flex;
    align-items: center;
    padding: 15px;
}

.business-icon {
    width: 48px;
    height: 48px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    margin-right: 12px;
}

.business-info {
    flex: 1;
}

.business-value {
    font-size: 24px;
    font-weight: bold;
    color: #303133;
}

.business-title {
    font-size: 14px;
    color: #909399;
    margin-top: 4px;
}
</style>
