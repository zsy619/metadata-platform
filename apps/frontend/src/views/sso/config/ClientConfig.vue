<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon">
                    <Monitor />
                </el-icon>
                客户端配置
            </h1>
            <div class="header-actions">
                <el-button type="primary" @click="handleCreate" :icon="Plus">新增客户端</el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-input v-model="searchForm.clientName" placeholder="请输入客户端名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
                <el-select v-model="searchForm.clientType" placeholder="客户端类型" style="width: 150px; margin-left: 10px" clearable @change="handleSearch">
                    <el-option label="全部" value="" />
                    <el-option label="Web 应用" value="web" />
                    <el-option label="SPA 应用" value="spa" />
                    <el-option label="移动应用" value="mobile" />
                    <el-option label="后端服务" value="backend" />
                </el-select>
                <el-select v-model="searchForm.isEnabled" placeholder="状态" style="width: 150px; margin-left: 10px" clearable @change="handleSearch">
                    <el-option label="全部" value="" />
                    <el-option label="启用" :value="true" />
                    <el-option label="禁用" :value="false" />
                </el-select>
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
                <el-button @click="handleReset" :icon="RefreshLeft">重置</el-button>
            </div>
            <div class="table-area">
                <el-table v-loading="loading" :element-loading-text="loadingText" :data="filteredClients" border stripe style="width: 100%; height: 100%;">
                    <template #empty>
                        <el-empty :description="searchForm.clientName ? '未搜索到相关客户端' : '暂无客户端配置'">
                            <el-button v-if="!searchForm.clientName" type="primary" @click="handleCreate">新增客户端</el-button>
                        </el-empty>
                    </template>
                    <el-table-column type="index" label="序号" width="60" align="center" />
                    <el-table-column prop="client_name" label="客户端名称" min-width="150" show-overflow-tooltip />
                    <el-table-column prop="client_id" label="客户端 ID" min-width="200" show-overflow-tooltip />
                    <el-table-column prop="client_type" label="客户端类型" width="110" align="center">
                        <template #default="scope">
                            <el-tag :type="getClientTypeType(scope.row.client_type)" effect="plain">
                                {{ getClientTypeLabel(scope.row.client_type) }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="redirect_uris" label="回调地址" min-width="200" show-overflow-tooltip>
                        <template #default="scope">
                            {{ scope.row.redirect_uris?.split(',').length || 0 }} 个地址
                        </template>
                    </el-table-column>
                    <el-table-column prop="is_enabled" label="状态" width="80" align="center">
                        <template #default="scope">
                            <el-switch v-model="scope.row.is_enabled" @change="(val: boolean) => handleStatusChange(scope.row, val)" />
                        </template>
                    </el-table-column>
                    <el-table-column prop="create_at" label="创建时间" width="170" />
                    <el-table-column label="操作" width="200" fixed="right" align="center">
                        <template #default="scope">
                            <el-button type="primary" size="small" :icon="View" @click="handleView(scope.row)" text bg>查看</el-button>
                            <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>编辑</el-button>
                            <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg>删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
        </el-card>

        <!-- 新增/编辑对话框 -->
        <ClientForm
            v-model:visible="dialogVisible"
            :title="dialogTitle"
            :data="currentClient"
            :submit-loading="submitLoading"
            @submit="handleSubmit"
        />

        <!-- 查看详情对话框 -->
        <ClientDetail
            v-model:visible="viewDialogVisible"
            :data="viewData"
        />
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, RefreshLeft, View, Edit, Delete, Monitor } from '@element-plus/icons-vue'
import type { SsoClient } from '@/types/sso'
import { getClients, createClient, updateClient, deleteClient } from '@/api/sso'
import ClientForm from './components/ClientForm.vue'
import ClientDetail from './components/ClientDetail.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchForm = reactive({ clientName: '', clientType: '', isEnabled: '' as '' | boolean })
const clients = ref<SsoClient[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitLoading = ref(false)
const currentClient = ref<Partial<SsoClient>>({})

const viewDialogVisible = ref(false)
const viewData = ref<Partial<SsoClient>>({})

const filteredClients = computed(() => {
    let result = clients.value
    if (searchForm.clientName) {
        result = result.filter(c => c.client_name?.toLowerCase().includes(searchForm.clientName.toLowerCase()))
    }
    if (searchForm.clientType) {
        result = result.filter(c => c.client_type === searchForm.clientType)
    }
    if (searchForm.isEnabled !== '') {
        result = result.filter(c => c.is_enabled === searchForm.isEnabled)
    }
    return result
})

const getClientTypeLabel = (type: string) => {
    const typeMap: Record<string, string> = { web: 'Web 应用', spa: 'SPA 应用', mobile: '移动应用', backend: '后端服务' }
    return typeMap[type?.toLowerCase()] || type
}

const getClientTypeType = (type: string) => {
    const typeMap: Record<string, string> = { web: 'primary', spa: 'success', mobile: 'warning', backend: 'info' }
    return typeMap[type?.toLowerCase()] || ''
}

const fetchData = async () => {
    loading.value = true
    try {
        const res: any = await getClients()
        clients.value = res.data || res || []
    } catch (error) {
        console.error('加载客户端配置失败:', error)
        ElMessage.error('加载客户端配置失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => { }
const handleDebouncedSearch = () => { }
const handleReset = () => {
    searchForm.clientName = ''
    searchForm.clientType = ''
    searchForm.isEnabled = ''
}

const handleCreate = () => {
    dialogTitle.value = '新增客户端'
    currentClient.value = {}
    dialogVisible.value = true
}

const handleEdit = (row: SsoClient) => {
    dialogTitle.value = '编辑客户端'
    currentClient.value = { ...row }
    dialogVisible.value = true
}

const handleView = (row: SsoClient) => {
    viewData.value = { ...row }
    viewDialogVisible.value = true
}

const handleSubmit = async (data: Partial<SsoClient>) => {
    submitLoading.value = true
    try {
        if (data.id) {
            await updateClient(data.id, data)
            ElMessage.success('更新成功')
        } else {
            await createClient(data)
            ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        fetchData()
    } catch (error: any) {
        ElMessage.error(error.message || '操作失败')
    } finally {
        submitLoading.value = false
    }
}

const handleDelete = (row: SsoClient) => {
    ElMessageBox.confirm(`确定要删除客户端 "${row.client_name}" 吗？删除后无法恢复！`, '确认删除', { type: 'warning', confirmButtonClass: 'el-button--danger' })
        .then(async () => {
            try {
                await deleteClient(row.id)
                ElMessage.success('删除成功')
                fetchData()
            } catch (error: any) {
                ElMessage.error(error.message || '删除失败')
            }
        })
        .catch(() => { })
}

const handleStatusChange = async (row: SsoClient, val: boolean) => {
    try {
        await updateClient(row.id, { is_enabled: val })
        ElMessage.success(val ? '已启用' : '已禁用')
    } catch (error: any) {
        ElMessage.error(error.message || '状态更新失败')
        row.is_enabled = !val
    }
}

onMounted(() => fetchData())
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
