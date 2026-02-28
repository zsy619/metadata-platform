<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon">
                    <Switch />
                </el-icon>
                字段映射
            </h1>
            <div class="header-actions">
                <el-button type="primary" @click="handleCreate" :icon="Plus">新增映射</el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-select v-model="searchForm.protocolId" placeholder="选择协议" style="width: 200px" clearable @change="handleSearch">
                    <el-option v-for="p in protocols" :key="p.id" :label="p.config_name" :value="p.id" />
                </el-select>
                <el-input v-model="searchForm.sourceField" placeholder="请输入源字段搜索" clearable :prefix-icon="Search" style="width: 250px; margin-left: 10px" @input="handleDebouncedSearch" />
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
                <el-button @click="handleReset" :icon="RefreshLeft">重置</el-button>
            </div>
            <div class="table-area">
                <el-table v-loading="loading" :element-loading-text="loadingText" :data="filteredMappings" border stripe style="width: 100%; height: 100%;">
                    <template #empty>
                        <el-empty :description="searchForm.sourceField ? '未搜索到相关映射' : '暂无字段映射'">
                            <el-button v-if="!searchForm.sourceField" type="primary" @click="handleCreate">新增映射</el-button>
                        </el-empty>
                    </template>
                    <el-table-column type="index" label="序号" width="60" align="center" />
                    <el-table-column prop="mapping_name" label="映射名称" min-width="150" />
                    <el-table-column prop="source_field" label="源字段" min-width="150" />
                    <el-table-column prop="target_field" label="目标字段" min-width="150" />
                    <el-table-column prop="field_type" label="字段类型" width="100" align="center">
                        <template #default="scope">
                            <el-tag size="small">{{ scope.row.field_type }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="is_required" label="必填" width="80" align="center">
                        <template #default="scope">
                            <el-tag :type="scope.row.is_required ? 'danger' : 'info'" size="small">{{ scope.row.is_required ? '是' : '否' }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="is_enabled" label="状态" width="80" align="center">
                        <template #default="scope">
                            <el-switch v-model="scope.row.is_enabled" @change="(val: boolean) => handleStatusChange(scope.row, val)" />
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="150" fixed="right" align="center">
                        <template #default="scope">
                            <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>编辑</el-button>
                            <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg>删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
        </el-card>

        <!-- 字段映射对话框组件 -->
        <FieldMappingDialog
            v-model:visible="dialogVisible"
            :title="dialogTitle"
            :form-data="formData"
            :loading="submitLoading"
            @submit="handleDialogSubmit"
        />

    </div>
</template>

<script setup lang="ts">
import { createFieldMapping, deleteFieldMapping, getFieldMappings, getProtocolConfigs, updateFieldMapping } from '@/api/sso'
import type { SsoFieldMapping, SsoProtocolConfig } from '@/types/sso'
import { Delete, Edit, Plus, RefreshLeft, Search, Switch } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'
import FieldMappingDialog from '@/components/sso/FieldMappingDialog.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchForm = reactive({ protocolId: '', sourceField: '' })
const mappings = ref<SsoFieldMapping[]>([])
const protocols = ref<SsoProtocolConfig[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitLoading = ref(false)

const formData = reactive<Partial<SsoFieldMapping>>({
    mapping_name: '',
    source_field: '',
    target_field: '',
    field_type: 'string',
    default_value: '',
    is_required: false,
    is_enabled: true,
    remark: ''
})



const filteredMappings = computed(() => {
    let result = mappings.value
    if (searchForm.protocolId) {
        result = result.filter(m => m.protocol_config_id === searchForm.protocolId)
    }
    if (searchForm.sourceField) {
        result = result.filter(m => m.source_field?.toLowerCase().includes(searchForm.sourceField.toLowerCase()))
    }
    return result
})

const fetchProtocols = async () => {
    try {
        const res: any = await getProtocolConfigs()
        protocols.value = res.data || res || []
    } catch (error) {
        console.error('加载协议失败:', error)
    }
}

const fetchData = async () => {
    loading.value = true
    try {
        const res: any = await getFieldMappings()
        mappings.value = res.data || res || []
    } catch (error) {
        console.error('加载字段映射失败:', error)
        ElMessage.error('加载字段映射失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => { }
const handleDebouncedSearch = () => { }
const handleReset = () => {
    searchForm.protocolId = ''
    searchForm.sourceField = ''
}

const resetForm = () => {
    formData.id = undefined
    formData.mapping_name = ''
    formData.source_field = ''
    formData.target_field = ''
    formData.field_type = 'string'
    formData.default_value = ''
    formData.is_required = false
    formData.is_enabled = true
    formData.remark = ''
}

const handleCreate = () => {
    dialogTitle.value = '新增字段映射'
    resetForm()
    dialogVisible.value = true
}

const handleEdit = (row: SsoFieldMapping) => {
    dialogTitle.value = '编辑字段映射'
    Object.assign(formData, row)
    dialogVisible.value = true
}

const handleDialogSubmit = async (data: Partial<SsoFieldMapping>) => {
    submitLoading.value = true
    try {
        if (data.id) {
            await updateFieldMapping(data.id, data)
            ElMessage.success('更新成功')
        } else {
            await createFieldMapping(data)
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

const handleDelete = (row: SsoFieldMapping) => {
    ElMessageBox.confirm(`确定要删除映射 "${row.source_field} → ${row.target_field}" 吗？`, '确认删除', { type: 'warning' })
        .then(async () => {
            try {
                await deleteFieldMapping(row.id)
                ElMessage.success('删除成功')
                fetchData()
            } catch (error: any) {
                ElMessage.error(error.message || '删除失败')
            }
        })
        .catch(() => { })
}

const handleStatusChange = async (row: SsoFieldMapping, val: boolean) => {
    try {
        await updateFieldMapping(row.id, { is_enabled: val })
        ElMessage.success(val ? '已启用' : '已禁用')
    } catch (error: any) {
        ElMessage.error(error.message || '状态更新失败')
        row.is_enabled = !val
    }
}

onMounted(() => {
    fetchProtocols()
    fetchData()
})
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
