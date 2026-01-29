<template>
    <div class="container-padding">
        <!-- 统计卡片 -->
        <el-row :gutter="20" class="mb-4">
            <el-col :span="6" v-for="item in cards" :key="item.title">
                <el-card shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <span>{{ item.title }}</span>
                            <el-tag :type="item.type" effect="plain">{{ item.tag }}</el-tag>
                        </div>
                    </template>
                    <div class="card-value">{{ item.value }}</div>
                    <div class="card-desc">
                        {{ item.desc }}
                        <span :class="item.trend > 0 ? 'up' : 'down'">
                            {{ Math.abs(item.trend) }}%
                            <el-icon>
                                <Top v-if="item.trend > 0" />
                                <Bottom v-else />
                            </el-icon>
                        </span>
                    </div>
                </el-card>
            </el-col>
        </el-row>
        <!-- 图表区域 -->
        <el-row :gutter="20">
            <el-col :span="16">
                <el-card title="流量趋势">
                    <div ref="lineChartRef" style="height: 350px;"></div>
                </el-card>
            </el-col>
            <el-col :span="8">
                <el-card title="错误分布">
                    <div ref="pieChartRef" style="height: 350px;"></div>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>
<script setup lang="ts">
import { Bottom, Top } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { onMounted, onUnmounted, ref } from 'vue'

const lineChartRef = ref<HTMLElement>()
const pieChartRef = ref<HTMLElement>()
let lineChart: echarts.EChartsType | null = null
let pieChart: echarts.EChartsType | null = null

const cards = ref([
    { title: '总请求数', value: '0', desc: '较昨日', trend: 0, type: 'primary', tag: 'Total', key: 'requests' },
    { title: '平均响应时间', value: '0ms', desc: '较昨日', trend: 0, type: 'success', tag: 'Avg', key: 'latency' },
    { title: '错误率', value: '0%', desc: '较昨日', trend: 0, type: 'danger', tag: 'Error', key: 'error_rate' },
    { title: 'QPS', value: '0', desc: '较昨日', trend: 0, type: 'warning', tag: 'Realtime', key: 'qps' }
])

// WebSocket Logic
let ws: WebSocket | null = null
const xAxisData = ref<string[]>([])
const seriesData = ref<number[]>([])

const initWebSocket = () => {
    // Determine WS protocol based on current protocol
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    // Use localhost:8888 for dev, or relative path for prod
    const wsUrl = `${protocol}//localhost:8888/api/monitor/ws`

    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
        console.log('Monitor WebSocket Connected')
    }

    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data)
            if (data.type === 'realtime') {
                // Update Cards
                cards.value[0].value = data.requests.toLocaleString()
                cards.value[1].value = `${data.latency}ms`
                cards.value[2].value = `${(data.error_rate * 100).toFixed(2)}%`
                cards.value[3].value = data.qps.toString()

                // Update Chart
                const timeStr = data.timestamp
                updateLineChart(timeStr, data.qps)
            }
        } catch (e) {
            console.error('WS Data Parse Error', e)
        }
    }

    ws.onclose = () => {
        console.log('Monitor WebSocket Closed, retrying in 3s...')
        setTimeout(initWebSocket, 3000)
    }
}

const updateLineChart = (time: string, value: number) => {
    if (!lineChart) return

    xAxisData.value.push(time)
    seriesData.value.push(value)

    if (xAxisData.value.length > 20) {
        xAxisData.value.shift()
        seriesData.value.shift()
    }

    lineChart.setOption({
        xAxis: { data: xAxisData.value },
        series: [{ data: seriesData.value }]
    })
}

const initCharts = () => {
    if (lineChartRef.value) {
        lineChart = echarts.init(lineChartRef.value)
        lineChart.setOption({
            tooltip: { trigger: 'axis' },
            xAxis: { type: 'category', data: [] },
            yAxis: { type: 'value' },
            series: [{ data: [], type: 'line', smooth: true, areaStyle: {} }]
        })
    }
    // ... Pie chart init (keep mock for now or update if backend sends it)
    if (pieChartRef.value) {
        pieChart = echarts.init(pieChartRef.value)
        pieChart.setOption({
            tooltip: { trigger: 'item' },
            series: [
                {
                    type: 'pie',
                    radius: ['40%', '70%'],
                    data: [
                        { value: 1048, name: 'Success' },
                        { value: 735, name: 'Validation Error' },
                        { value: 580, name: 'Timeout' },
                        { value: 484, name: 'Auth Error' },
                        { value: 300, name: 'Server Error' }
                    ]
                }
            ]
        })
    }
}

const resizeHandler = () => {
    lineChart?.resize()
    pieChart?.resize()
}

onMounted(() => {
    initCharts()
    initWebSocket()
    window.addEventListener('resize', resizeHandler)
})

onUnmounted(() => {
    ws?.close()
    window.removeEventListener('resize', resizeHandler)
    lineChart?.dispose()
    pieChart?.dispose()
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

.up {
    color: #f56c6c;
    margin-left: 5px;
}

.down {
    color: #67c23a;
    margin-left: 5px;
}

.mb-4 {
    margin-bottom: 20px;
}
</style>
