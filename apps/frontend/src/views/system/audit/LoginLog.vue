<template>
    <div class="container-padding">
        <!-- 页面标题区 -->
        <div class="page-header">
            <h1 class="page-title">
                <el-icon class="title-icon">
                    <User />
                </el-icon>
                登录日志
            </h1>
            <div class="header-actions">
                <el-button type="warning" :icon="Download" @click="handleExport">导出</el-button>
            </div>
        </div>
        <!-- 主内容卡片 -->
        <el-card class="main-card">
            <!-- 搜索区域 -->
            <div class="search-area">
                <el-input v-model="queryParams.account" placeholder="请输入账号" clearable style="width: 200px" @keyup.enter="handleQuery" />
                <el-input v-model="queryParams.client_ip" placeholder="请输入IP" clearable style="width: 200px" @keyup.enter="handleQuery" />
                <el-select v-model="queryParams.login_status" placeholder="状态" clearable style="width: 120px" @change="handleQuery">
                    <el-option label="成功" :value="1" />
                    <el-option label="失败" :value="0" />
                    <el-option label="退出" :value="2" />
                </el-select>
                <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" value-format="YYYY-MM-DD" style="width: 300px" @change="handleQuery" />
                <el-button type="primary" :icon="Search" @click="handleQuery">搜索</el-button>
                <el-button :icon="RefreshLeft" @click="resetQuery">重置</el-button>
            </div>
            <!-- 表格区域 -->
            <div class="table-area">
                <el-table v-loading="loading" :data="tableData" border style="width: 100%" height="100%">
                    <el-table-column prop="create_at" label="登录时间" width="180" align="center">
                        <template #default="scope">
                            {{ formatDate(scope.row.create_at) }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="account" label="账号" width="120" align="center" />
                    <el-table-column prop="client_ip" label="登录IP" width="140" align="center" />
                    <el-table-column prop="login_status" label="状态" width="100" align="center">
                        <template #default="scope">
                            <el-tag :type="scope.row.login_status === 1 ? 'success' : (scope.row.login_status === 2 ? 'info' : 'danger')">
                                {{ scope.row.login_status === 1 ? '成功' : (scope.row.login_status === 2 ? '退出' : '失败') }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="browser" label="浏览器" width="120" align="center" show-overflow-tooltip />
                    <el-table-column prop="os" label="操作系统" width="120" align="center" show-overflow-tooltip />
                    <el-table-column prop="ip_location" label="归属地" align="center" />
                    <el-table-column label="操作" width="100" fixed="right" align="center">
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
        <el-dialog v-model="detailsVisible" title="登录日志详情" width="700px" append-to-body>
            <el-descriptions :column="2" border v-if="currentLog">
                <el-descriptions-item label="登录账号">{{ currentLog.account }}</el-descriptions-item>
                <el-descriptions-item label="用户ID">{{ currentLog.user_id }}</el-descriptions-item>
                <el-descriptions-item label="登录IP">{{ currentLog.client_ip }}</el-descriptions-item>
                <el-descriptions-item label="登录位置">{{ currentLog.ip_location || '-' }}</el-descriptions-item>
                <el-descriptions-item label="登录时间">{{ formatDate(currentLog.create_at) }}</el-descriptions-item>
                <el-descriptions-item label="登录状态">
                    <el-tag :type="currentLog.login_status === 1 ? 'success' : (currentLog.login_status === 2 ? 'info' : 'danger')">
                        {{ currentLog.login_status === 1 ? '成功' : (currentLog.login_status === 2 ? '退出' : '失败') }}
                    </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="浏览器">{{ currentLog.browser }} {{ currentLog.browser_version }}</el-descriptions-item>
                <el-descriptions-item label="内核">{{ currentLog.browser_engine || '-' }}</el-descriptions-item>
                <el-descriptions-item label="操作系统">{{ currentLog.os }} {{ currentLog.os_version }}</el-descriptions-item>
                <el-descriptions-item label="架构">{{ currentLog.os_arch || '-' }}</el-descriptions-item>
                <el-descriptions-item label="设备类型">{{ currentLog.device_type }}</el-descriptions-item>
                <el-descriptions-item label="设备型号">{{ currentLog.device_model || '-' }}</el-descriptions-item>
                <el-descriptions-item label="语言">{{ currentLog.language || '-' }}</el-descriptions-item>
                <el-descriptions-item label="平台">{{ currentLog.platform || '-' }}</el-descriptions-item>
                <el-descriptions-item label="时区">{{ currentLog.timezone || '-' }}</el-descriptions-item>
                <el-descriptions-item label="分辨率">{{ currentLog.screen_resolution || '-' }}</el-descriptions-item>
                <el-descriptions-item label="User-Agent" :span="2">{{ currentLog.user_agent }}</el-descriptions-item>
                <el-descriptions-item label="错误信息" :span="2" v-if="currentLog.error_message">
                    <span style="color: #f56c6c">{{ currentLog.error_message }}</span>
                </el-descriptions-item>
            </el-descriptions>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { exportLoginLogs, getLoginLogs } from '@/api/audit'
import { Download, RefreshLeft, Search, User, View } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

const loading = ref(false)
const total = ref(0)
const tableData = ref([])
const dateRange = ref([])
const detailsVisible = ref(false)
const currentLog = ref<any>(null)

const queryParams = reactive({
    page: 1,
    pageSize: 20,
    account: '',
    client_ip: '',
    login_status: undefined,
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

        const res: any = await getLoginLogs(queryParams)
        // 关键修复：后端 SuccessWithPagination 将数据放在 .data.list 中
        const paginatedData = res?.data || {}
        tableData.value = paginatedData.list || []
        total.value = paginatedData.total || 0
    } catch (error) {
        console.error('获取登录日志失败:', error)
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
    queryParams.account = ''
    queryParams.client_ip = ''
    queryParams.login_status = undefined
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

const handleExport = async () => {
    try {
        const res = await exportLoginLogs(queryParams)
        const blob = new Blob([JSON.stringify(res, null, 2)], { type: 'application/json' })
        const link = document.createElement('a')
        link.href = URL.createObjectURL(blob)
        link.download = `login_logs_${dayjs().format('YYYYMMDDHHmmss')}.json`
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
/* ==================== 标准布局样式 ==================== */
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

:deep(.el-descriptions__label) {
    width: 120px;
    font-weight: bold;
}
</style>
