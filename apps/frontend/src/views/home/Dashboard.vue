<template>
    <div class="dashboard-container">
        <!-- 统计卡片行 -->
        <div class="stats-row">
            <el-card class="stat-card" shadow="hover">
                <div class="stat-content">
                    <div class="stat-icon user-icon">
                        <el-icon><User /></el-icon>
                    </div>
                    <div class="stat-info">
                        <div class="stat-value">{{ stats.user_count }}</div>
                        <div class="stat-label">用户总数</div>
                    </div>
                </div>
            </el-card>
            <el-card class="stat-card" shadow="hover">
                <div class="stat-content">
                    <div class="stat-icon role-icon">
                        <el-icon><Key /></el-icon>
                    </div>
                    <div class="stat-info">
                        <div class="stat-value">{{ stats.role_count }}</div>
                        <div class="stat-label">角色总数</div>
                    </div>
                </div>
            </el-card>
            <el-card class="stat-card" shadow="hover">
                <div class="stat-content">
                    <div class="stat-icon org-icon">
                        <el-icon><OfficeBuilding /></el-icon>
                    </div>
                    <div class="stat-info">
                        <div class="stat-value">{{ stats.org_count }}</div>
                        <div class="stat-label">组织总数</div>
                    </div>
                </div>
            </el-card>
            <el-card class="stat-card" shadow="hover">
                <div class="stat-content">
                    <div class="stat-icon menu-icon">
                        <el-icon><Menu /></el-icon>
                    </div>
                    <div class="stat-info">
                        <div class="stat-value">{{ stats.menu_count }}</div>
                        <div class="stat-label">菜单总数</div>
                    </div>
                </div>
            </el-card>
            <el-card class="stat-card" shadow="hover">
                <div class="stat-content">
                    <div class="stat-icon pos-icon">
                        <el-icon><Briefcase /></el-icon>
                    </div>
                    <div class="stat-info">
                        <div class="stat-value">{{ stats.pos_count }}</div>
                        <div class="stat-label">职位总数</div>
                    </div>
                </div>
            </el-card>
            <el-card class="stat-card" shadow="hover">
                <div class="stat-content">
                    <div class="stat-icon group-icon">
                        <el-icon><Folder /></el-icon>
                    </div>
                    <div class="stat-info">
                        <div class="stat-value">{{ stats.user_group_count }}</div>
                        <div class="stat-label">用户组</div>
                    </div>
                </div>
            </el-card>
            <el-card class="stat-card" shadow="hover">
                <div class="stat-content">
                    <div class="stat-icon role-group-icon">
                        <el-icon><Collection /></el-icon>
                    </div>
                    <div class="stat-info">
                        <div class="stat-value">{{ stats.role_group_count }}</div>
                        <div class="stat-label">角色组</div>
                    </div>
                </div>
            </el-card>
        </div>

        <!-- 图表区域 -->
        <el-row :gutter="20" class="chart-row">
            <el-col :span="12">
                <el-card class="chart-card" shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <span>登录趋势（最近7天）</span>
                            <el-tag type="success" effect="plain">Login Trend</el-tag>
                        </div>
                    </template>
                    <div ref="loginTrendChartRef" class="chart-container" v-loading="loading.chart"></div>
                </el-card>
            </el-col>
            <el-col :span="12">
                <el-card class="chart-card" shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <span>角色分布</span>
                            <el-tag type="warning" effect="plain">Role Distribution</el-tag>
                        </div>
                    </template>
                    <div ref="rolePieChartRef" class="chart-container" v-loading="loading.chart"></div>
                </el-card>
            </el-col>
        </el-row>

        <el-row :gutter="20" class="chart-row">
            <el-col :span="8">
                <el-card class="chart-card" shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <span>组织架构分布</span>
                            <el-tag type="info" effect="plain">Org Structure</el-tag>
                        </div>
                    </template>
                    <div ref="orgTreeChartRef" class="chart-container-small" v-loading="loading.chart"></div>
                </el-card>
            </el-col>
            <el-col :span="8">
                <el-card class="chart-card" shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <span>用户状态分布</span>
                            <el-tag type="primary" effect="plain">User Status</el-tag>
                        </div>
                    </template>
                    <div ref="userStatusChartRef" class="chart-container-small" v-loading="loading.chart"></div>
                </el-card>
            </el-col>
            <el-col :span="8">
                <el-card class="chart-card" shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <span>操作日志统计</span>
                            <el-tag type="danger" effect="plain">Operation Log</el-tag>
                        </div>
                    </template>
                    <div ref="operationChartRef" class="chart-container-small" v-loading="loading.chart"></div>
                </el-card>
            </el-col>
        </el-row>

        <!-- 日志区域 -->
        <div class="logs-row">
            <el-card class="log-card" shadow="hover">
                <template #header>
                    <div class="card-header">
                        <span>最近登录日志</span>
                        <el-button type="primary" link @click="goToLoginLogs">查看更多</el-button>
                    </div>
                </template>
                <el-table :data="loginLogs" style="width: 100%" size="small" max-height="280">
                    <el-table-column prop="user_account" label="用户账号" width="120" />
                    <el-table-column prop="user_name" label="用户名称" width="100" />
                    <el-table-column prop="login_ip" label="登录IP" width="130" />
                    <el-table-column prop="login_status" label="状态" width="80">
                        <template #default="scope">
                            <el-tag :type="scope.row.login_status === 1 ? 'success' : 'danger'" size="small">
                                {{ scope.row.login_status === 1 ? '成功' : '失败' }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="create_at" label="登录时间" show-overflow-tooltip>
                        <template #default="scope">
                            {{ formatTime(scope.row.create_at) }}
                        </template>
                    </el-table-column>
                </el-table>
            </el-card>

            <el-card class="log-card" shadow="hover">
                <template #header>
                    <div class="card-header">
                        <span>最近操作日志</span>
                        <el-button type="primary" link @click="goToOperationLogs">查看更多</el-button>
                    </div>
                </template>
                <el-table :data="operationLogs" style="width: 100%" size="small" max-height="280">
                    <el-table-column prop="user_account" label="操作用户" width="100" />
                    <el-table-column prop="module" label="模块" width="80" />
                    <el-table-column prop="action" label="操作" width="80" />
                    <el-table-column prop="description" label="描述" show-overflow-tooltip />
                    <el-table-column prop="create_at" label="操作时间" width="160">
                        <template #default="scope">
                            {{ formatTime(scope.row.create_at) }}
                        </template>
                    </el-table-column>
                </el-table>
            </el-card>
        </div>
    </div>
</template>

<script setup lang="ts">
import {
    getDashboardStats,
    getRecentLoginLogs,
    getRecentOperationLogs,
    getLoginTrend,
    getUserStatusDistribution,
    getOperationStats,
    getOrgDistribution,
    type LoginTrendItem,
    type UserStatusDistribution,
    type OperationStats,
    type OrgDistribution
} from '@/api/user'
import { Briefcase, Collection, Folder, Key, Menu, OfficeBuilding, User } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import { onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const stats = ref({
    user_count: 0,
    role_count: 0,
    org_count: 0,
    menu_count: 0,
    pos_count: 0,
    user_group_count: 0,
    role_group_count: 0
})

const loginLogs = ref<any[]>([])
const operationLogs = ref<any[]>([])

// 图表数据
const loginTrendData = ref<LoginTrendItem[]>([])
const userStatusData = ref<UserStatusDistribution>({ active: 0, inactive: 0, locked: 0, pending: 0 })
const operationStatsData = ref<OperationStats>({ create: 0, update: 0, delete: 0, query: 0, export: 0 })
const orgDistributionData = ref<OrgDistribution[]>([])

const loading = ref({
    chart: false
})

const formatTime = (time: string) => {
    if (!time) return '-'
    const date = new Date(time)
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
    })
}

const loadStats = async () => {
    try {
        const res = await getDashboardStats()
        stats.value = res
    } catch (error: any) {
        ElMessage.error(error.message || '加载统计数据失败')
    }
}

const loadLoginLogs = async () => {
    try {
        const res = await getRecentLoginLogs()
        loginLogs.value = res || []
    } catch (error: any) {
        console.error('加载登录日志失败:', error)
    }
}

const loadOperationLogs = async () => {
    try {
        const res = await getRecentOperationLogs()
        operationLogs.value = res || []
    } catch (error: any) {
        console.error('加载操作日志失败:', error)
    }
}

// 加载图表数据
const loadChartData = async () => {
    loading.value.chart = true
    try {
        // 并行加载所有图表数据
        const [trend, status, operation, org] = await Promise.all([
            getLoginTrend(),
            getUserStatusDistribution(),
            getOperationStats(),
            getOrgDistribution()
        ])
        loginTrendData.value = trend || []
        userStatusData.value = status || { active: 0, inactive: 0, locked: 0, pending: 0 }
        operationStatsData.value = operation || { create: 0, update: 0, delete: 0, query: 0, export: 0 }
        orgDistributionData.value = org || []
    } catch (error: any) {
        console.error('加载图表数据失败:', error)
    } finally {
        loading.value.chart = false
    }
}

const goToLoginLogs = () => {
    router.push('/audit/login-log')
}

const goToOperationLogs = () => {
    router.push('/audit/operation-log')
}

// 图表引用
const loginTrendChartRef = ref<HTMLElement>()
const rolePieChartRef = ref<HTMLElement>()
const orgTreeChartRef = ref<HTMLElement>()
const userStatusChartRef = ref<HTMLElement>()
const operationChartRef = ref<HTMLElement>()

// 图表实例
let loginTrendChart: echarts.ECharts | null = null
let rolePieChart: echarts.ECharts | null = null
let orgTreeChart: echarts.ECharts | null = null
let userStatusChart: echarts.ECharts | null = null
let operationChart: echarts.ECharts | null = null

// 初始化登录趋势图
const initLoginTrendChart = () => {
    if (!loginTrendChartRef.value) return
    loginTrendChart = echarts.init(loginTrendChartRef.value)

    const dates = loginTrendData.value.map(item => item.date)
    const successData = loginTrendData.value.map(item => item.success)
    const failData = loginTrendData.value.map(item => item.fail)

    loginTrendChart.setOption({
        tooltip: {
            trigger: 'axis',
            axisPointer: { type: 'shadow' }
        },
        legend: {
            data: ['登录成功', '登录失败'],
            bottom: 0
        },
        grid: {
            left: '3%',
            right: '4%',
            bottom: '15%',
            top: '10%',
            containLabel: true
        },
        xAxis: {
            type: 'category',
            data: dates,
            axisLabel: { color: '#606266' }
        },
        yAxis: {
            type: 'value',
            name: '次数',
            axisLabel: { color: '#606266' }
        },
        series: [
            {
                name: '登录成功',
                type: 'bar',
                stack: 'total',
                data: successData,
                itemStyle: { color: '#67C23A' },
                barWidth: '40%'
            },
            {
                name: '登录失败',
                type: 'bar',
                stack: 'total',
                data: failData,
                itemStyle: { color: '#F56C6C' },
                barWidth: '40%'
            }
        ]
    })
}

// 初始化角色分布饼图
const initRolePieChart = () => {
    if (!rolePieChartRef.value) return
    rolePieChart = echarts.init(rolePieChartRef.value)

    rolePieChart.setOption({
        tooltip: {
            trigger: 'item',
            formatter: '{b}: {c} ({d}%)'
        },
        legend: {
            orient: 'vertical',
            right: '5%',
            top: 'center'
        },
        series: [
            {
                type: 'pie',
                radius: ['40%', '70%'],
                center: ['40%', '50%'],
                avoidLabelOverlap: false,
                itemStyle: {
                    borderRadius: 10,
                    borderColor: '#fff',
                    borderWidth: 2
                },
                label: {
                    show: false,
                    position: 'center'
                },
                emphasis: {
                    label: {
                        show: true,
                        fontSize: 16,
                        fontWeight: 'bold'
                    }
                },
                labelLine: { show: false },
                data: [
                    { value: stats.value.role_count || 0, name: '系统角色', itemStyle: { color: '#409EFF' } },
                    { value: stats.value.role_group_count || 0, name: '角色组', itemStyle: { color: '#67C23A' } },
                    { value: stats.value.user_group_count || 0, name: '用户组', itemStyle: { color: '#E6A23C' } },
                    { value: stats.value.pos_count || 0, name: '职位角色', itemStyle: { color: '#F56C6C' } }
                ]
            }
        ]
    })
}

// 初始化组织架构图
const initOrgTreeChart = () => {
    if (!orgTreeChartRef.value) return
    orgTreeChart = echarts.init(orgTreeChartRef.value)

    // 使用真实数据或默认数据
    const chartData = orgDistributionData.value.length > 0
        ? orgDistributionData.value.map((item, index) => ({
            name: item.name,
            value: item.value,
            itemStyle: { color: ['#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399'][index % 5] }
        }))
        : [
            { name: '暂无数据', value: 1, itemStyle: { color: '#909399' } }
        ]

    orgTreeChart.setOption({
        tooltip: { trigger: 'item' },
        series: [
            {
                type: 'treemap',
                data: chartData,
                breadcrumb: { show: false },
                label: {
                    show: true,
                    formatter: '{b}',
                    fontSize: 12
                },
                itemStyle: {
                    borderColor: '#fff',
                    borderWidth: 2,
                    gapWidth: 2
                }
            }
        ]
    })
}

// 初始化用户状态图
const initUserStatusChart = () => {
    if (!userStatusChartRef.value) return
    userStatusChart = echarts.init(userStatusChartRef.value)

    const { active, inactive, locked, pending } = userStatusData.value

    userStatusChart.setOption({
        tooltip: { trigger: 'item' },
        series: [
            {
                type: 'pie',
                radius: '70%',
                center: ['50%', '50%'],
                data: [
                    { value: active, name: '正常', itemStyle: { color: '#67C23A' } },
                    { value: inactive, name: '禁用', itemStyle: { color: '#F56C6C' } },
                    { value: locked, name: '锁定', itemStyle: { color: '#E6A23C' } },
                    { value: pending, name: '待审核', itemStyle: { color: '#909399' } }
                ],
                label: {
                    formatter: '{b}\n{d}%',
                    fontSize: 11
                },
                emphasis: {
                    itemStyle: {
                        shadowBlur: 10,
                        shadowOffsetX: 0,
                        shadowColor: 'rgba(0, 0, 0, 0.5)'
                    }
                }
            }
        ]
    })
}

// 初始化操作日志统计图
const initOperationChart = () => {
    if (!operationChartRef.value) return
    operationChart = echarts.init(operationChartRef.value)

    const { create, update, delete: del, query, export: exp } = operationStatsData.value

    // 计算最大值
    const maxVal = Math.max(create, update, del, query, exp, 100)

    operationChart.setOption({
        tooltip: { trigger: 'axis' },
        radar: {
            indicator: [
                { name: '新增', max: maxVal },
                { name: '修改', max: maxVal },
                { name: '删除', max: maxVal },
                { name: '查询', max: maxVal },
                { name: '导出', max: maxVal }
            ],
            center: ['50%', '55%'],
            radius: '60%'
        },
        series: [
            {
                type: 'radar',
                data: [
                    {
                        value: [create, update, del, query, exp],
                        name: '操作统计',
                        areaStyle: {
                            color: 'rgba(64, 158, 255, 0.3)'
                        },
                        lineStyle: { color: '#409EFF' },
                        itemStyle: { color: '#409EFF' }
                    }
                ]
            }
        ]
    })
}

// 窗口大小变化时重绘图表
const resizeHandler = () => {
    loginTrendChart?.resize()
    rolePieChart?.resize()
    orgTreeChart?.resize()
    userStatusChart?.resize()
    operationChart?.resize()
}

onMounted(async () => {
    // 加载基础数据
    await loadStats()
    loadLoginLogs()
    loadOperationLogs()

    // 加载图表数据
    await loadChartData()

    // 初始化图表
    initLoginTrendChart()
    initRolePieChart()
    initOrgTreeChart()
    initUserStatusChart()
    initOperationChart()

    // 监听窗口大小变化
    window.addEventListener('resize', resizeHandler)
})

onUnmounted(() => {
    window.removeEventListener('resize', resizeHandler)
    loginTrendChart?.dispose()
    rolePieChart?.dispose()
    orgTreeChart?.dispose()
    userStatusChart?.dispose()
    operationChart?.dispose()
})
</script>

<style scoped>
.dashboard-container {
    padding: 20px;
    background-color: #f5f7fa;
    min-height: calc(100vh - 60px);
}

.stats-row {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    margin-bottom: 20px;
}

.stat-card {
    flex: 1;
    min-width: 180px;
    border-radius: 8px;
    transition: transform 0.3s ease;
}

.stat-card:hover {
    transform: translateY(-2px);
}

.stat-content {
    display: flex;
    align-items: center;
    gap: 16px;
}

.stat-icon {
    width: 56px;
    height: 56px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 28px;
    color: white;
}

.user-icon {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.role-icon {
    background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.org-icon {
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.menu-icon {
    background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.pos-icon {
    background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
}

.group-icon {
    background: linear-gradient(135deg, #a18cd1 0%, #fbc2eb 100%);
}

.role-group-icon {
    background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%);
}

.stat-info {
    flex: 1;
}

.stat-value {
    font-size: 28px;
    font-weight: 600;
    color: #303133;
    line-height: 1.2;
}

.stat-label {
    font-size: 14px;
    color: #909399;
    margin-top: 4px;
}

.chart-row {
    margin-bottom: 20px;
}

.chart-card {
    border-radius: 8px;
    height: 100%;
}

.chart-container {
    height: 280px;
    width: 100%;
}

.chart-container-small {
    height: 220px;
    width: 100%;
}

.logs-row {
    display: flex;
    gap: 20px;
}

.log-card {
    flex: 1;
    border-radius: 8px;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 600;
    font-size: 16px;
}

:deep(.el-card__header) {
    padding: 12px 20px;
    border-bottom: 1px solid #ebeef5;
}

:deep(.el-card__body) {
    padding: 16px 20px;
}
</style>
