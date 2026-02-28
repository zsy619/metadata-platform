<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon">
                    <Key />
                </el-icon>
                密钥管理
            </h1>
            <div class="header-actions">
                <el-button type="primary" @click="handleCreate" :icon="Plus">生成密钥</el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-input v-model="searchForm.keyName" placeholder="请输入密钥名称搜索" clearable :prefix-icon="Search" style="width: 300px" />
                <el-select v-model="searchForm.keyType" placeholder="密钥类型" style="width: 150px; margin-left: 10px" clearable>
                    <el-option label="全部" value="" />
                    <el-option label="RSA" value="rsa" />
                    <el-option label="EC" value="ec" />
                    <el-option label="Octet" value="octet" />
                </el-select>
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
                <el-button @click="handleReset" :icon="RefreshLeft">重置</el-button>
            </div>
            <div class="table-area">
                <el-table v-loading="loading" :element-loading-text="loadingText" :data="filteredKeys" border stripe style="width: 100%; height: 100%;">
                    <template #empty>
                        <el-empty :description="searchForm.keyName ? '未搜索到相关密钥' : '暂无密钥'">
                            <el-button v-if="!searchForm.keyName" type="primary" @click="handleCreate">生成密钥</el-button>
                        </el-empty>
                    </template>
                    <el-table-column type="index" label="序号" width="60" align="center" />
                    <el-table-column prop="key_name" label="密钥名称" min-width="150" show-overflow-tooltip />
                    <el-table-column prop="key_type" label="密钥类型" width="100" align="center">
                        <template #default="scope">
                            <el-tag :type="getKeyTypeType(scope.row.key_type)" effect="plain">{{ getKeyTypeLabel(scope.row.key_type) }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="key_usage" label="用途" width="80" align="center">
                        <template #default="scope">
                            <el-tag :type="scope.row.key_usage === 'signing' ? 'success' : 'warning'" size="small">{{ getKeyUsageLabel(scope.row.key_usage) }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="algorithm" label="算法" width="100" align="center" />
                    <el-table-column prop="key_id" label="Key ID" min-width="200" show-overflow-tooltip />
                    <el-table-column prop="is_enabled" label="状态" width="80" align="center">
                        <template #default="scope">
                            <el-switch v-model="scope.row.is_enabled" @change="(val: boolean) => handleStatusChange(scope.row, val)" />
                        </template>
                    </el-table-column>
                    <el-table-column prop="create_at" label="创建时间" width="170" />
                    <el-table-column label="操作" width="150" fixed="right" align="center">
                        <template #default="scope">
                            <el-button type="primary" size="small" :icon="View" @click="handleView(scope.row)" text bg>查看</el-button>
                            <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg>删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
        </el-card>

        <KeyGenerateDialog
            v-model="dialogVisible"
            @submit="handleGenerateSubmit"
        />

        <KeyViewDialog
            v-model="viewDialogVisible"
            :data="viewData"
        />
    </div>
</template>

<script setup lang="ts">
import { deleteKey, generateKeyPair, getKeys, updateKey } from '@/api/sso'
import type { SsoKey } from '@/types/sso'
import { Delete, Key, Plus, RefreshLeft, Search, View } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'
import KeyGenerateDialog from './components/KeyGenerateDialog.vue'
import KeyViewDialog from './components/KeyViewDialog.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchForm = reactive({ keyName: '', keyType: '' })
const keys = ref<SsoKey[]>([])

const dialogVisible = ref(false)
const viewDialogVisible = ref(false)
const viewData = ref<Partial<SsoKey>>({})

const filteredKeys = computed(() => {
    let result = keys.value
    if (searchForm.keyName) {
        result = result.filter(k => k.key_name?.toLowerCase().includes(searchForm.keyName.toLowerCase()))
    }
    if (searchForm.keyType) {
        result = result.filter(k => k.key_type?.toLowerCase() === searchForm.keyType.toLowerCase())
    }
    return result
})

// 工具函数
const getKeyTypeLabel = (type: string) => {
    const typeMap: Record<string, string> = { 'rsa': 'RSA', 'ec': 'EC', 'octet': 'Octet' }
    return typeMap[type?.toLowerCase()] || type
}

const getKeyTypeType = (type: string) => {
    const typeMap: Record<string, string> = { 'rsa': 'primary', 'ec': 'success', 'octet': 'warning' }
    return typeMap[type?.toLowerCase()] || ''
}

const getKeyUsageLabel = (usage: string) => {
    const usageMap: Record<string, string> = { 'signing': '签名', 'encryption': '加密', 'both': '两者' }
    return usageMap[usage] || usage
}

const fetchData = async () => {
    loading.value = true
    try {
        const res: any = await getKeys()
        keys.value = res.data || res || []
    } catch (error) {
        console.error('加载密钥失败:', error)
        ElMessage.error('加载密钥失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => { }
const handleReset = () => { searchForm.keyName = ''; searchForm.keyType = '' }

const handleCreate = () => {
    dialogVisible.value = true
}

// 处理生成密钥提交
const handleGenerateSubmit = async (data: { key_type: 'rsa' | 'octet'; algorithm: string; remark: string }) => {
    try {
        await generateKeyPair(data.key_type, data.algorithm)
        ElMessage.success('密钥生成成功')
        fetchData()
    } catch (error: any) {
        ElMessage.error(error.message || '密钥生成失败')
    }
}

const handleView = (row: SsoKey) => {
    viewData.value = { ...row }
    viewDialogVisible.value = true
}

const handleDelete = (row: SsoKey) => {
    ElMessageBox.confirm(`确定要删除密钥 "${row.key_name}" 吗？`, '确认删除', { type: 'warning' })
        .then(async () => {
            try {
                await deleteKey(row.id)
                ElMessage.success('删除成功')
                fetchData()
            } catch (error: any) {
                ElMessage.error(error.message || '删除失败')
            }
        })
        .catch(() => { })
}

const handleStatusChange = async (row: SsoKey, val: boolean) => {
    try {
        await updateKey(row.id, { is_enabled: val })
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
