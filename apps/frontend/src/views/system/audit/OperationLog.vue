<template>
    <div class="app-container">
        <el-card class="filter-container" shadow="never">
            <el-form :inline="true" :model="queryParams" class="demo-form-inline">
                <el-form-item label="模块">
                    <el-input v-model="queryParams.module" placeholder="请输入模块名" clearable />
                </el-form-item>
                <el-form-item label="操作类型">
                    <el-input v-model="queryParams.type" placeholder="如：INSERT, QUERY" clearable />
                </el-form-item>
                <el-form-item label="操作时间">
                    <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" value-format="YYYY-MM-DD" />
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
                    <el-button icon="Refresh" @click="resetQuery">重置</el-button>
                    <el-button type="warning" icon="Download" @click="handleExport">导出</el-button>
                </el-form-item>
            </el-form>
        </el-card>
        <el-card class="box-card" shadow="never" style="margin-top: 20px;">
            <el-table v-loading="loading" :data="tableData" border style="width: 100%">
                <el-table-column prop="trace_id" label="追踪ID" width="180" align="center" />
                <el-table-column prop="source" label="模块" width="120" align="center" />
                <el-table-column prop="method" label="类型" width="100" align="center">
                    <template #default="scope">
                        <el-tag>{{ scope.row.method }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="path" label="路径" show-overflow-tooltip />
                <el-table-column prop="user_id" label="操作人" width="120" align="center" />
                <el-table-column prop="client_ip" label="IP" width="140" align="center" />
                <el-table-column prop="status" label="状态" width="80" align="center">
                    <template #default="scope">
                        <el-tag :type="scope.row.status < 400 ? 'success' : 'danger'">
                            {{ scope.row.status }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="latency" label="耗时(ms)" width="100" align="center" />
                <el-table-column prop="create_at" label="操作时间" width="180" align="center">
                    <template #default="scope">
                        {{ formatDate(scope.row.create_at) }}
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="100" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="primary" size="small" :icon="View" @click="viewDetails(scope.row)">
                            详情
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
            <div class="pagination-container">
                <el-pagination v-model:current-page="queryParams.page" v-model:page-size="queryParams.pageSize" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
        </el-card>
        <!-- 详情对话框 -->
        <el-dialog v-model="detailsVisible" title="操作日志详情" width="800px">
            <el-descriptions :column="2" border v-if="currentLog">
                <el-descriptions-item label="追踪ID" :span="2">{{ currentLog.trace_id }}</el-descriptions-item>
                <el-descriptions-item label="操作人">{{ currentLog.user_id }}</el-descriptions-item>
                <el-descriptions-item label="请求方法">{{ currentLog.method }}</el-descriptions-item>
                <el-descriptions-item label="请求路径" :span="2">{{ currentLog.path }}</el-descriptions-item>
                <el-descriptions-item label="状态码">
                    <el-tag :type="currentLog.status < 400 ? 'success' : 'danger'">
                        {{ currentLog.status }}
                    </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="耗时">{{ currentLog.latency }} ms</el-descriptions-item>
                <el-descriptions-item label="客户端IP">{{ currentLog.client_ip }}</el-descriptions-item>
                <el-descriptions-item label="客户端OS">{{ currentLog.os }} {{ currentLog.os_version }}</el-descriptions-item>
                <el-descriptions-item label="浏览器">{{ currentLog.browser }} {{ currentLog.browser_version }}</el-descriptions-item>
                <el-descriptions-item label="设备类型">{{ currentLog.device_type }}</el-descriptions-item>
                <el-descriptions-item label="操作时间">{{ formatDate(currentLog.create_at) }}</el-descriptions-item>
                <el-descriptions-item label="请求UA" :span="2">{{ currentLog.user_agent }}</el-descriptions-item>
                <el-descriptions-item label="错误信息" :span="2" v-if="currentLog.error_message">
                    <pre class="error-msg">{{ currentLog.error_message }}</pre>
                </el-descriptions-item>
            </el-descriptions>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { exportOperationLogs, getOperationLogs } from '@/api/audit'
import { View } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

interface SysOperationLog {
    trace_id: string
    user_id: string
    source: string
    method: string
    path: string
    status: number
    latency: number
    client_ip: string
    user_agent: string
    browser: string
    browser_version: string
    os: string
    os_version: string
    device_type: string
    error_message: string
    create_at: string
}

const loading = ref(false)
const total = ref(0)
const tableData = ref<SysOperationLog[]>([])
const dateRange = ref([])
const detailsVisible = ref(false)
const currentLog = ref<SysOperationLog | null>(null)

const queryParams = reactive({
    page: 1,
    pageSize: 20,
    module: '',
    type: '',
    start_time: '',
    end_time: ''
})

const formatDate = (date: string) => {
    return date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '-'
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
        const res = await getOperationLogs(queryParams)
        const data = res as any
        tableData.value = data.list
        total.value = data.total
    } catch (error) {
        console.error(error)
    } finally {
        loading.value = false
    }
}

const handleQuery = () => {
    queryParams.page = 1
    getList()
}

const resetQuery = () => {
    queryParams.module = ''
    queryParams.type = ''
    dateRange.value = []
    handleQuery()
}

const handleSizeChange = (val: number) => {
    queryParams.pageSize = val
    getList()
}

const handleCurrentChange = (val: number) => {
    queryParams.page = val
    getList()
}

const viewDetails = (log: SysOperationLog) => {
    currentLog.value = log
    detailsVisible.value = true
}

const handleExport = async () => {
    try {
        const res = await exportOperationLogs(queryParams)
        const blob = new Blob([res as any], { type: 'application/json' })
        const link = document.createElement('a')
        link.href = URL.createObjectURL(blob)
        link.download = `operation_logs_${new Date().getTime()}.json`
        link.click()
        URL.revokeObjectURL(link.href)
    } catch (error) {
        ElMessage.error('导出失败')
    }
}

onMounted(() => {
    getList()
})
</script>
<style scoped>
.filter-container {
    margin-bottom: 20px;
}

.pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
}

.error-msg {
    white-space: pre-wrap;
    word-break: break-all;
    background-color: #f5f7fa;
    padding: 10px;
    border-radius: 4px;
    color: #f56c6c;
    font-size: 12px;
}
</style>
