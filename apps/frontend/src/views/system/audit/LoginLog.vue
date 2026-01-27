<template>
    <div class="app-container">
        <!-- Search Bar -->
        <el-card class="filter-container" shadow="never">
            <el-form :inline="true" :model="queryParams" class="demo-form-inline">
                <el-form-item label="用户名">
                    <el-input v-model="queryParams.account" placeholder="请输入用户名" clearable @keyup.enter="handleQuery" />
                </el-form-item>
                <el-form-item label="登录时间">
                    <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" value-format="YYYY-MM-DD" />
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
                    <el-button icon="Refresh" @click="resetQuery">重置</el-button>
                    <el-button type="warning" icon="Download" @click="handleExport">导出</el-button>
                </el-form-item>
            </el-form>
        </el-card>
        <!-- Table -->
        <el-card class="box-card" shadow="never" style="margin-top: 20px;">
            <el-table v-loading="loading" :data="tableData" border style="width: 100%">
                <el-table-column prop="create_at" label="登录时间" width="180" align="center">
                    <template #default="scope">
                        {{ formatDate(scope.row.create_at) }}
                    </template>
                </el-table-column>
                <el-table-column prop="user_id" label="用户ID" width="120" align="center" />
                <el-table-column prop="account" label="账号" width="120" align="center" />
                <el-table-column prop="client_ip" label="登录IP" width="140" align="center" />
                <el-table-column prop="login_status" label="状态" width="100" align="center">
                    <template #default="scope">
                        <el-tag :type="scope.row.login_status === 1 ? 'success' : (scope.row.login_status === 2 ? 'info' : 'danger')">
                            {{ scope.row.login_status === 1 ? '成功' : (scope.row.login_status === 2 ? '退出' : '失败') }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="browser" label="浏览器" width="120" align="center" />
                <el-table-column prop="ip_location" label="归属地" align="center" />
                <el-table-column prop="error_message" label="消息" show-overflow-tooltip />
                <el-table-column label="操作" width="100" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="primary" size="small" :icon="View" @click="viewDetails(scope.row)">
                            详情
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
            <!-- Pagination -->
            <div class="pagination-container">
                <el-pagination v-model:currentPage="queryParams.page" v-model:pageSize="queryParams.pageSize" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
        </el-card>
        <!-- 详情对话框 -->
        <el-dialog v-model="detailsVisible" title="登录详情" width="600px">
            <pre v-if="currentLog">{{ JSON.stringify(currentLog, null, 2) }}</pre>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { exportLoginLogs, getLoginLogs } from '@/api/audit'
import { View } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

const loading = ref(false)
const total = ref(0)
const tableData = ref([])
const dateRange = ref([])
const detailsVisible = ref(false)
const currentLog = ref(null)

const queryParams = reactive({
    page: 1,
    pageSize: 20,
    account: '',
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

        const res = await getLoginLogs(queryParams)
        // Adjust based on your actual API response structure (e.g. res.data.list)
        // Assuming backend utils.Success returns data directly
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
    queryParams.account = ''
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

const viewDetails = (log: any) => {
    currentLog.value = log
    detailsVisible.value = true
}

const handleExport = async () => {
    try {
        const res = await exportLoginLogs(queryParams)
        const blob = new Blob([res as any], { type: 'application/json' }) // Update type to xlsx if real export
        // For now backend returns JSON, so we just download it as json to verify
        const link = document.createElement('a')
        link.href = URL.createObjectURL(blob)
        link.download = `login_logs_${new Date().getTime()}.json`
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
</style>
