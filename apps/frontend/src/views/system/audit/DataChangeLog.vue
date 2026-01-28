<template>
    <div class="container-padding">
        <!-- 页面标题区 -->
        <div class="page-header">
            <h1 class="page-title">
                <el-icon class="title-icon">
                    <Edit />
                </el-icon>
                数据变更日志
            </h1>
            <div class="header-actions">
                <el-button type="warning" :icon="Download" @click="handleExport">导出</el-button>
            </div>
        </div>
        <!-- 主内容卡片 -->
        <el-card class="main-card">
            <!-- 搜索区域 -->
            <div class="search-area">
                <el-input v-model="queryParams.trace_id" placeholder="请输入追踪ID" clearable style="width: 200px" @keyup.enter="handleQuery" />
                <el-input v-model="queryParams.user_id" placeholder="操作人" clearable style="width: 150px" @keyup.enter="handleQuery" />
                <el-select v-model="queryParams.data_type" placeholder="变更类型" clearable style="width: 120px" @change="handleQuery">
                    <el-option label="新增" value="CREATE" />
                    <el-option label="修改" value="UPDATE" />
                    <el-option label="删除" value="DELETE" />
                </el-select>
                <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" value-format="YYYY-MM-DD" style="width: 300px" @change="handleQuery" />
                <el-button type="primary" :icon="Search" @click="handleQuery">搜索</el-button>
                <el-button :icon="RefreshLeft" @click="resetQuery">重置</el-button>
            </div>
            <!-- 表格区域 -->
            <div class="table-area">
                <el-table v-loading="loading" :data="tableData" border style="width: 100%" height="100%">
                    <el-table-column prop="trace_id" label="追踪ID" width="180" align="center" show-overflow-tooltip />
                    <el-table-column prop="model_id" label="业务模型" width="150" align="center" />
                    <el-table-column prop="action" label="变更类型" width="100" align="center">
                        <template #default="scope">
                            <el-tag :type="getTypeTag(scope.row.action)">{{ scope.row.action }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="create_by" label="操作人" width="120" align="center" />
                    <el-table-column prop="create_at" label="变更时间" width="180" align="center">
                        <template #default="scope">
                            {{ formatDate(scope.row.create_at) }}
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="100" fixed="right" align="center">
                        <template #default="scope">
                            <el-button type="primary" link :icon="View" @click="viewDetails(scope.row)">
                                查看详情
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
        <el-dialog v-model="detailsVisible" title="数据变更详情" width="900px" destroy-on-close append-to-body>
            <el-descriptions :column="2" border v-if="currentLog">
                <el-descriptions-item label="追踪ID" :span="2">
                    <el-link type="primary" @click="goOperationLog(currentLog.trace_id)">{{ currentLog.trace_id }}</el-link>
                </el-descriptions-item>
                <el-descriptions-item label="模型名称">{{ currentLog.model_id }}</el-descriptions-item>
                <el-descriptions-item label="记录ID">{{ currentLog.record_id }}</el-descriptions-item>
                <el-descriptions-item label="变更类型">
                    <el-tag :type="getTypeTag(currentLog.action)">{{ currentLog.action }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="操作人">{{ currentLog.create_by }}</el-descriptions-item>
                <el-descriptions-item label="变更时间" :span="2">{{ formatDate(currentLog.create_at) }}</el-descriptions-item>
                <el-descriptions-item label="变更前数据" :span="2">
                    <div class="json-container" v-if="currentLog.before_data">
                        <pre>{{ formatJson(currentLog.before_data) }}</pre>
                    </div>
                    <span v-else>无 (新增操作)</span>
                </el-descriptions-item>
                <el-descriptions-item label="变更后数据" :span="2">
                    <div class="json-container" v-if="currentLog.after_data">
                        <pre>{{ formatJson(currentLog.after_data) }}</pre>
                    </div>
                    <span v-else>无 (删除操作)</span>
                </el-descriptions-item>
            </el-descriptions>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { exportDataChangeLogs, getDataChangeLogs } from '@/api/audit'
import { Download, Edit, RefreshLeft, Search, View } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const total = ref(0)
const tableData = ref([])
const dateRange = ref([])
const detailsVisible = ref(false)
const currentLog = ref<any>(null)

const queryParams = reactive({
    page: 1,
    pageSize: 20,
    trace_id: '',
    user_id: '',
    data_type: '',
    start_time: '',
    end_time: ''
})

const formatDate = (date: string) => {
    return date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '-'
}

const getTypeTag = (type: string) => {
    switch (type?.toUpperCase()) {
        case 'CREATE':
        case 'INSERT': return 'success'
        case 'UPDATE': return 'warning'
        case 'DELETE': return 'danger'
        default: return 'info'
    }
}

const formatJson = (jsonStr: string) => {
    if (!jsonStr) return '-'
    try {
        return JSON.stringify(JSON.parse(jsonStr), null, 2)
    } catch (e) {
        return jsonStr
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
        const res: any = await getDataChangeLogs(queryParams)
        // 关键修复：正确处理带分页的数据
        const paginatedData = res?.data || {}
        tableData.value = paginatedData.list || []
        total.value = paginatedData.total || 0
    } catch (error) {
        console.error('获取列表失败:', error)
        ElMessage.error('获取列表失败')
    } finally {
        loading.value = false
    }
}

const handleQuery = () => {
    queryParams.page = 1
    getList()
}

const resetQuery = () => {
    queryParams.trace_id = ''
    queryParams.user_id = ''
    queryParams.data_type = ''
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

const viewDetails = (log: any) => {
    currentLog.value = log
    detailsVisible.value = true
}

const goOperationLog = (traceId: string) => {
    router.push({
        path: '/system/audit/operation',
        query: { trace_id: traceId }
    })
}

const handleExport = async () => {
    try {
        const res = await exportDataChangeLogs(queryParams)
        const blob = new Blob([JSON.stringify(res, null, 2)], { type: 'application/json' })
        const link = document.createElement('a')
        link.href = URL.createObjectURL(blob)
        link.download = `data_change_logs_${dayjs().format('YYYYMMDDHHmmss')}.json`
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
/* ==================== 标准布局样式 ==================== */
.container-padding {
    padding: 20px;
    padding-bottom: 0;
    height: calc(100vh - 84px);
    display: flex;
    flex-direction: column;
    overflow: hidden;
    box-sizing: border-box;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
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

.json-container {
    background-color: #f5f7fa;
    padding: 10px;
    border-radius: 4px;
    max-height: 300px;
    overflow-y: auto;
}

.json-container pre {
    margin: 0;
    font-size: 12px;
    white-space: pre-wrap;
    word-break: break-all;
}

:deep(.el-descriptions__label) {
    width: 120px;
    font-weight: bold;
}
</style>
