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
                <el-table-column prop="id" label="ID" width="80" align="center" />
                <el-table-column prop="module" label="模块" width="120" align="center" />
                <el-table-column prop="type" label="类型" width="100" align="center">
                    <template #default="scope">
                        <el-tag>{{ scope.row.type }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="content" label="操作内容" show-overflow-tooltip />
                <el-table-column prop="operator" label="操作人" width="120" align="center" />
                <el-table-column prop="ip" label="IP" width="140" align="center" />
                <el-table-column prop="status" label="状态" width="80" align="center">
                    <template #default="scope">
                        <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
                            {{ scope.row.status === 1 ? '成功' : '失败' }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="cost_time" label="耗时(ms)" width="100" align="center" />
                <el-table-column prop="created_at" label="操作时间" width="180" align="center">
                    <template #default="scope">
                        {{ formatDate(scope.row.created_at) }}
                    </template>
                </el-table-column>
            </el-table>
            <div class="pagination-container">
                <el-pagination v-model:current-page="queryParams.page" v-model:page-size="queryParams.pageSize" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import { exportOperationLogs, getOperationLogs } from '@/api/audit'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

const loading = ref(false)
const total = ref(0)
const tableData = ref([])
const dateRange = ref([])

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
</style>
