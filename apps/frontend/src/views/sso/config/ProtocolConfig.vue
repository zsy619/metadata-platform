<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon">
                    <Connection />
                </el-icon>
                协议配置
            </h1>
            <div class="header-actions">
                <el-button type="primary" @click="handleCreate" :icon="Plus">新增协议</el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-input v-model="searchForm.configName" placeholder="请输入配置名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
                <el-select v-model="searchForm.protocolType" placeholder="协议类型" style="width: 150px; margin-left: 10px" clearable @change="handleSearch">
                    <el-option label="全部" value="" />
                    <el-option label="OIDC" value="oidc" />
                    <el-option label="SAML" value="saml" />
                    <el-option label="LDAP" value="ldap" />
                    <el-option label="CAS" value="cas" />
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
                <el-table v-loading="loading" :element-loading-text="loadingText" :data="filteredConfigs" border stripe style="width: 100%; height: 100%;">
                    <template #empty>
                        <el-empty :description="searchForm.configName ? '未搜索到相关配置' : '暂无协议配置'">
                            <el-button v-if="!searchForm.configName" type="primary" @click="handleCreate">新增协议</el-button>
                        </el-empty>
                    </template>
                    <el-table-column type="index" label="序号" width="60" align="center" />
                    <el-table-column prop="config_name" label="配置名称" min-width="150" />
                    <el-table-column prop="protocol_type" label="协议类型" width="100" align="center">
                        <template #default="scope">
                            <el-tag :type="getProtocolTypeType(scope.row.protocol_type)" effect="plain">
                                {{ getProtocolTypeLabel(scope.row.protocol_type) }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="is_enabled" label="状态" width="80" align="center">
                        <template #default="scope">
                            <el-switch v-model="scope.row.is_enabled" @change="(val: boolean) => handleStatusChange(scope.row, val)" />
                        </template>
                    </el-table-column>
                    <el-table-column prop="remark" label="备注" min-width="200" show-overflow-tooltip />
                    <el-table-column prop="create_at" label="创建时间" width="170" />
                    <el-table-column label="操作" width="150" fixed="right" align="center">
                        <template #default="scope">
                            <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>编辑</el-button>
                            <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg>删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
        </el-card>

        <ProtocolConfigDialog
            v-model="dialogVisible"
            :edit-data="currentEditData"
            @submit="handleDialogSubmit"
        />
    </div>
</template>

<script setup lang="ts">
import { createProtocolConfig, deleteProtocolConfig, getProtocolConfigs, updateProtocolConfig } from '@/api/sso'
import type { SsoProtocolConfig } from '@/types/sso'
import { Connection, Delete, Edit, Plus, RefreshLeft, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'
import ProtocolConfigDialog from './components/ProtocolConfigDialog.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchForm = reactive({ configName: '', protocolType: '', isEnabled: '' as '' | boolean })
const configs = ref<SsoProtocolConfig[]>([])

const dialogVisible = ref(false)
const currentEditData = ref<SsoProtocolConfig | null>(null)

const filteredConfigs = computed(() => {
    let result = configs.value
    if (searchForm.configName) {
        result = result.filter(c => c.config_name?.toLowerCase().includes(searchForm.configName.toLowerCase()))
    }
    if (searchForm.protocolType) {
        result = result.filter(c => c.protocol_type === searchForm.protocolType)
    }
    if (searchForm.isEnabled !== '') {
        result = result.filter(c => c.is_enabled === searchForm.isEnabled)
    }
    return result
})

const getProtocolTypeLabel = (type: string) => {
    const typeMap: Record<string, string> = { oidc: 'OIDC', saml: 'SAML', ldap: 'LDAP', cas: 'CAS' }
    return typeMap[type?.toLowerCase()] || type
}

const getProtocolTypeType = (type: string) => {
    const typeMap: Record<string, string> = { oidc: 'primary', saml: 'success', ldap: 'warning', cas: 'info' }
    return typeMap[type?.toLowerCase()] || ''
}

const fetchData = async () => {
    loading.value = true
    try {
        const res: any = await getProtocolConfigs()
        configs.value = res.data || res || []
    } catch (error) {
        console.error('加载协议配置失败:', error)
        ElMessage.error('加载协议配置失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => { }
const handleDebouncedSearch = () => { }
const handleReset = () => {
    searchForm.configName = ''
    searchForm.protocolType = ''
    searchForm.isEnabled = ''
}

// 打开新增对话框
const handleCreate = () => {
    currentEditData.value = null
    dialogVisible.value = true
}

// 打开编辑对话框
const handleEdit = (row: SsoProtocolConfig) => {
    currentEditData.value = row
    dialogVisible.value = true
}

// 处理对话框提交
const handleDialogSubmit = async (data: Partial<SsoProtocolConfig> & { config_data: string }) => {
    try {
        if (data.id) {
            await updateProtocolConfig(data.id, data)
            ElMessage.success('更新成功')
        } else {
            await createProtocolConfig(data)
            ElMessage.success('创建成功')
        }
        fetchData()
    } catch (error: any) {
        ElMessage.error(error.message || '操作失败')
    }
}

const handleDelete = (row: SsoProtocolConfig) => {
    ElMessageBox.confirm(`确定要删除协议配置 "${row.config_name}" 吗？`, '确认删除', { type: 'warning' })
        .then(async () => {
            try {
                await deleteProtocolConfig(row.id)
                ElMessage.success('删除成功')
                fetchData()
            } catch (error: any) {
                ElMessage.error(error.message || '删除失败')
            }
        })
        .catch(() => { })
}

const handleStatusChange = async (row: SsoProtocolConfig, val: boolean) => {
    try {
        await updateProtocolConfig(row.id, { is_enabled: val })
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
