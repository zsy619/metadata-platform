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

const cards = [
    { title: '总请求数', value: '1,234,567', desc: '较昨日', trend: 12.5, type: 'primary', tag: 'Total' },
    { title: '平均响应时间', value: '128ms', desc: '较昨日', trend: -5.2, type: 'success', tag: 'Avg' },
    { title: '错误率', value: '0.45%', desc: '较昨日', trend: -0.1, type: 'danger', tag: 'Error' },
    { title: 'QPS', value: '450', desc: '较昨日', trend: 8.4, type: 'warning', tag: 'Realtime' }
]

const initCharts = () => {
    if (lineChartRef.value) {
        lineChart = echarts.init(lineChartRef.value)
        lineChart.setOption({
            tooltip: { trigger: 'axis' },
            xAxis: { type: 'category', data: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00'] },
            yAxis: { type: 'value' },
            series: [{ data: [120, 132, 101, 134, 90, 230], type: 'line', smooth: true, areaStyle: {} }]
        })
    }
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
    window.addEventListener('resize', resizeHandler)
})

onUnmounted(() => {
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
