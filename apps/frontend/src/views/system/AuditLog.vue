<template>
    <div class="audit-log">
        <div class="page-header">
            <h1>审计日志</h1>
        </div>
        <el-card>
            <el-tabs v-model="activeTab">
                <el-tab-pane label="操作日志" name="operation">
                    <el-table :data="operationLogs" border v-loading="loadingOp">
                        <el-table-column prop="createdAt" label="时间" width="180" />
                        <el-table-column prop="userName" label="用户" width="120" />
                        <el-table-column prop="method" label="方法" width="80" />
                        <el-table-column prop="path" label="路径" />
                        <el-table-column prop="status" label="状态" width="100">
                            <template #default="scope">
                                <el-tag :type="scope.row.status < 400 ? 'success' : 'danger'">
                                    {{ scope.row.status }}
                                </el-tag>
                            </template>
                        </el-table-column>
                        <el-table-column prop="latency" label="耗时(ms)" width="100" />
                        <el-table-column prop="clientIP" label="IP" width="140" />
                        <el-table-column label="操作" width="120">
                            <template #default="scope">
                                <el-button type="primary" link @click="viewDetails(scope.row)">详情</el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-tab-pane>
                <el-tab-pane label="数据变更" name="dataChange">
                    <el-table :data="dataChangeLogs" border v-loading="loadingData">
                        <el-table-column prop="createdAt" label="时间" width="180" />
                        <el-table-column prop="action" label="动作" width="100">
                            <template #default="scope">
                                <el-tag :type="getActionType(scope.row.action)">{{ scope.row.action }}</el-tag>
                            </template>
                        </el-table-column>
                        <el-table-column prop="modelID" label="模型ID" width="100" />
                        <el-table-column prop="recordID" label="记录ID" width="100" />
                        <el-table-column label="详情">
                            <template #default="scope">
                                <el-button type="primary" link @click="viewDiff(scope.row)">查看变更</el-button>
                            </template>
                        </el-table-column>
                    </el-table>
                </el-tab-pane>
            </el-tabs>
        </el-card>
        <!-- 详情对话框 -->
        <el-dialog v-model="detailsVisible" title="详情" width="600px" class="custom-dialog">
            <pre v-if="currentLog">{{ JSON.stringify(currentLog, null, 2) }}</pre>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { getDataChangeLogs, getOperationLogs } from '@/api/audit'
import { onMounted, ref } from 'vue'

const activeTab = ref('operation')
const operationLogs = ref([])
const dataChangeLogs = ref([])
const loadingOp = ref(false)
const loadingData = ref(false)
const detailsVisible = ref(false)
const currentLog = ref(null)

onMounted(() => {
    fetchOpLogs()
    fetchDataChangeLogs()
})

const fetchOpLogs = async () => {
    loadingOp.value = true
    try {
        const res = await getOperationLogs()
        operationLogs.value = res.data || []
    } finally {
        loadingOp.value = false
    }
}

const fetchDataChangeLogs = async () => {
    loadingData.value = true
    try {
        const res = await getDataChangeLogs()
        dataChangeLogs.value = res.data || []
    } finally {
        loadingData.value = false
    }
}

const getActionType = (action: string) => {
    if (action === 'Create') return 'success'
    if (action === 'Update') return 'warning'
    if (action === 'Delete') return 'danger'
    return 'info'
}

const viewDetails = (log: any) => {
    currentLog.value = log
    detailsVisible.value = true
}

const viewDiff = (log: any) => {
    currentLog.value = log
    detailsVisible.value = true
}
</script>
<style scoped>
.audit-log {
    padding: 20px;
}

.page-header {
    margin-bottom: 20px;
}
</style>
