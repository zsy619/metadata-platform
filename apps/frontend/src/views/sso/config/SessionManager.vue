<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon">
                    <Timer />
                </el-icon>
                会话管理
            </h1>
            <div class="header-actions">
                <el-button type="danger" @click="handleBatchOffline" :disabled="!selectedSessions.length" :icon="CircleClose">批量下线</el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-input v-model="searchForm.userId" placeholder="用户ID" clearable :prefix-icon="Search" style="width: 200px" />
                <el-select v-model="searchForm.clientId" placeholder="选择客户端" style="width: 180px; margin-left: 10px" clearable>
                    <el-option v-for="c in clients" :key="c.id" :label="c.client_name" :value="c.id" />
                </el-select>
                <el-select v-model="searchForm.status" placeholder="状态" style="width: 120px; margin-left: 10px" clearable>
                    <el-option label="全部" value="" />
                    <el-option label="活跃" value="active" />
                    <el-option label="已过期" value="expired" />
                    <el-option label="已注销" value="revoked" />
                </el-select>
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
                <el-button @click="handleReset" :icon="RefreshLeft">重置</el-button>
            </div>
            <div class="table-area">
                <el-table v-loading="loading" :element-loading-text="loadingText" :data="filteredSessions" border stripe @selection-change="handleSelectionChange" style="width: 100%; height: 100%;">
                    <template #empty>
                        <el-empty :description="searchForm.userId ? '未搜索到相关会话' : '暂无会话'" />
                    </template>
                    <el-table-column type="selection" width="55" align="center" />
                    <el-table-column type="index" label="序号" width="60" align="center" />
                    <el-table-column prop="user_id" label="用户ID" min-width="120" />
                    <el-table-column prop="client_id" label="客户端ID" min-width="150" show-overflow-tooltip />
                    <el-table-column prop="protocol_type" label="协议" width="100" align="center">
                        <template #default="scope">
                            <el-tag size="small">{{ scope.row.protocol_type?.toUpperCase() }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="ip_address" label="IP地址" width="130" />
                    <el-table-column prop="status" label="状态" width="90" align="center">
                        <template #default="scope">
                            <el-tag :type="getStatusType(scope.row.status)" size="small">{{ getStatusLabel(scope.row.status) }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="create_at" label="登录时间" width="170" />
                    <el-table-column prop="expires_at" label="过期时间" width="170" />
                    <el-table-column label="操作" width="120" fixed="right" align="center">
                        <template #default="scope">
                            <el-button type="danger" size="small" :icon="CircleClose" @click="handleOffline(scope.row)" :disabled="scope.row.status !== 'active'" text bg>下线</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
        </el-card>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, RefreshLeft, CircleClose, Timer } from '@element-plus/icons-vue'
import type { SsoSession, SsoClient } from '@/types/sso'
import { getSessions, revokeSession, getClients } from '@/api/sso'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchForm = reactive({ userId: '', clientId: '', status: '' })
const sessions = ref<SsoSession[]>([])
const clients = ref<SsoClient[]>([])
const selectedSessions = ref<SsoSession[]>([])

const filteredSessions = computed(() => {
    let result = sessions.value
    if (searchForm.userId) {
        result = result.filter(s => s.user_id?.toLowerCase().includes(searchForm.userId.toLowerCase()))
    }
    if (searchForm.clientId) {
        result = result.filter(s => s.client_id === searchForm.clientId)
    }
    if (searchForm.status) {
        result = result.filter(s => s.status === searchForm.status)
    }
    return result
})

const getStatusType = (status: string) => {
    const typeMap: Record<string, string> = { active: 'success', expired: 'info', revoked: 'danger', inactive: 'warning' }
    return typeMap[status] || 'info'
}

const getStatusLabel = (status: string) => {
    const labelMap: Record<string, string> = { active: '活跃', expired: '已过期', revoked: '已注销', inactive: '不活跃' }
    return labelMap[status] || status
}

const fetchClients = async () => {
    try {
        const res: any = await getClients()
        clients.value = res.data || res || []
    } catch (error) {
        console.error('加载客户端失败:', error)
    }
}

const fetchData = async () => {
    loading.value = true
    try {
        const res: any = await getSessions()
        sessions.value = res.data || res || []
    } catch (error) {
        console.error('加载会话失败:', error)
        ElMessage.error('加载会话失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => { }
const handleReset = () => { searchForm.userId = ''; searchForm.clientId = ''; searchForm.status = '' }
const handleSelectionChange = (selection: SsoSession[]) => { selectedSessions.value = selection }

const handleOffline = (row: SsoSession) => {
    ElMessageBox.confirm(`确定要强制下线用户 "${row.user_id}" 吗？`, '确认下线', { type: 'warning' })
        .then(async () => {
            try {
                await revokeSession(row.id)
                ElMessage.success('强制下线成功')
                fetchData()
            } catch (error: any) {
                ElMessage.error(error.message || '操作失败')
            }
        })
        .catch(() => { })
}

const handleBatchOffline = () => {
    const count = selectedSessions.value.length
    ElMessageBox.confirm(`确定要批量下线选中的 ${count} 个会话吗？`, '确认批量下线', { type: 'warning' })
        .then(async () => {
            try {
                const promises = selectedSessions.value.map(s => revokeSession(s.id))
                await Promise.all(promises)
                ElMessage.success(`成功下线 ${count} 个会话`)
                fetchData()
            } catch (error: any) {
                ElMessage.error(error.message || '批量下线失败')
            }
        })
        .catch(() => { })
}

onMounted(() => { fetchClients(); fetchData() })
</script>

<style scoped>
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
}

.text-primary {
    color: var(--el-text-color-primary);
}
</style>
