<template>
    <div class="app-container">
        <el-card class="filter-container" shadow="never">
            <el-form :inline="true" :model="queryParams" class="demo-form-inline">
                <el-form-item label="表名">
                    <el-input v-model="queryParams.table_name" placeholder="请输入表名" clearable />
                </el-form-item>
                <el-form-item label="变更类型">
                    <el-select v-model="queryParams.data_type" placeholder="请选择" clearable>
                        <el-option label="新增" value="INSERT" />
                        <el-option label="修改" value="UPDATE" />
                        <el-option label="删除" value="DELETE" />
                    </el-select>
                </el-form-item>
                <el-form-item label="变更时间">
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
                <el-table-column prop="model_id" label="表/模型" width="150" align="center" />
                <el-table-column prop="action" label="类型" width="100" align="center">
                    <template #default="scope">
                        <el-tag :type="getTypeTag(scope.row.action)">{{ scope.row.action }}</el-tag>
                    </template>
                </el-table-column>
                <el-table-column label="变更前数据" min-width="200">
                    <template #default="scope">
                        <el-popover placement="top-start" :width="400" trigger="hover" v-if="scope.row.before_data">
                            <template #reference>
                                <div class="json-preview">{{ scope.row.before_data }}</div>
                            </template>
                            <pre>{{ formatJson(scope.row.before_data) }}</pre>
                        </el-popover>
                        <span v-else>-</span>
                    </template>
                </el-table-column>
                <el-table-column label="变更后数据" min-width="200">
                    <template #default="scope">
                        <el-popover placement="top-start" :width="400" trigger="hover" v-if="scope.row.after_data">
                            <template #reference>
                                <div class="json-preview">{{ scope.row.after_data }}</div>
                            </template>
                            <pre>{{ formatJson(scope.row.after_data) }}</pre>
                        </el-popover>
                        <span v-else>-</span>
                    </template>
                </el-table-column>
                <el-table-column prop="create_by" label="操作人" width="120" align="center" />
                <el-table-column prop="create_at" label="变更时间" width="180" align="center">
                    <template #default="scope">
                        {{ formatDate(scope.row.create_at) }}
                    </template>
                </el-table-column>
            </el-table>
            <div class="pagination-container">
                <el-pagination v-model:currentPage="queryParams.page" v-model:pageSize="queryParams.pageSize" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import { exportDataChangeLogs, getDataChangeLogs } from '@/api/audit'
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
    table_name: '',
    data_type: '',
    start_time: '',
    end_time: ''
})

const formatDate = (date: string) => {
    return date ? dayjs(date).format('YYYY-MM-DD HH:mm:ss') : '-'
}

const getTypeTag = (type: string) => {
    switch (type) {
        case 'INSERT': return 'success'
        case 'UPDATE': return 'warning'
        case 'DELETE': return 'danger'
        default: return 'info'
    }
}

const formatJson = (jsonStr: string) => {
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
        const res = await getDataChangeLogs(queryParams)
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
    queryParams.table_name = ''
    queryParams.data_type = ''
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
        const res = await exportDataChangeLogs(queryParams)
        const blob = new Blob([res as any], { type: 'application/json' })
        const link = document.createElement('a')
        link.href = URL.createObjectURL(blob)
        link.download = `data_logs_${new Date().getTime()}.json`
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

.json-preview {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 300px;
    cursor: pointer;
    color: #409EFF;
}
</style>
