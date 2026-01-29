<template>
    <div class="container-padding">
        <!-- 慢查询 TOP 10 -->
        <el-card title="慢查询 TOP 10 (模拟数据)" class="mb-4">
            <el-table :data="slowQueries" stripe style="width: 100%">
                <el-table-column prop="query" label="SQL语句" show-overflow-tooltip />
                <el-table-column prop="duration" label="耗时(ms)" width="120" sortable />
                <el-table-column prop="timestamp" label="发生时间" width="180" />
                <el-table-column prop="user" label="用户" width="120" />
            </el-table>
        </el-card>
        <!-- 错误接口 TOP 10 -->
        <el-card title="错误接口 TOP 10 (模拟数据)" class="mb-4">
            <el-table :data="errorInterfaces" stripe style="width: 100%">
                <el-table-column prop="path" label="接口路径" show-overflow-tooltip />
                <el-table-column prop="method" label="请求方法" width="100" />
                <el-table-column prop="errorCount" label="错误次数" width="100" sortable />
                <el-table-column prop="latestError" label="最新错误信息" show-overflow-tooltip />
                <el-table-column prop="timestamp" label="最近发生时间" width="180" />
            </el-table>
        </el-card>
        <!-- 响应时间分布 -->
        <el-card title="响应时间分布">
            <div ref="barChartRef" style="height: 400px;"></div>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import * as echarts from 'echarts'
import { onMounted, onUnmounted, ref } from 'vue'

const barChartRef = ref<HTMLElement>()
let barChart: echarts.EChartsType | null = null

const slowQueries = [
    { query: 'SELECT * FROM huge_table WHERE name LIKE "%test%"', duration: 1250, timestamp: '2024-03-20 10:00:00', user: 'admin' },
    { query: 'SELECT COUNT(*) FROM logs', duration: 890, timestamp: '2024-03-20 10:05:00', user: 'user1' },
    { query: 'UPDATE settings SET value = "x" WHERE id = 1', duration: 500, timestamp: '2024-03-20 10:10:00', user: 'admin' },
    // ... more mock data
]

const errorInterfaces = [
    { path: '/api/v1/users/123', method: 'GET', errorCount: 45, latestError: '500 Internal Server Error', timestamp: '2024-03-20 10:15:00' },
    { path: '/api/v1/login', method: 'POST', errorCount: 23, latestError: '401 Unauthorized', timestamp: '2024-03-20 10:20:00' },
    { path: '/api/v1/data/sync', method: 'POST', errorCount: 12, latestError: '504 Gateway Timeout', timestamp: '2024-03-20 10:25:00' },
    // ... more mock data
]

const initChart = () => {
    if (barChartRef.value) {
        barChart = echarts.init(barChartRef.value)
        barChart.setOption({
            tooltip: { trigger: 'axis' },
            xAxis: { type: 'category', data: ['<10ms', '10-50ms', '50-100ms', '100-500ms', '500-1s', '>1s'] },
            yAxis: { type: 'value' },
            series: [{
                data: [12000, 5000, 3000, 1000, 200, 50],
                type: 'bar',
                showBackground: true,
                backgroundStyle: { color: 'rgba(180, 180, 180, 0.2)' }
            }]
        })
    }
}

const resizeHandler = () => {
    barChart?.resize()
}

onMounted(() => {
    initChart()
    window.addEventListener('resize', resizeHandler)
})

onUnmounted(() => {
    window.removeEventListener('resize', resizeHandler)
    barChart?.dispose()
})
</script>
<style scoped>
.mb-4 {
    margin-bottom: 20px;
}
</style>
