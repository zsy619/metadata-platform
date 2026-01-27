<template>
    <div class="model-list">
        <div class="page-header">
            <h1>模型管理</h1>
            <el-button type="primary" @click="handleCreate" :icon="Plus">
                新增模型
            </el-button>
        </div>
        <el-card>
            <div class="card-header">
                <el-input v-model="searchQuery" placeholder="请输入模型名称或编码搜索" clearable :prefix-icon="Search" style="width: 300px" @clear="handleSearch" @keyup.enter="handleSearch" />
                <el-select v-model="filterKind" placeholder="筛选模型类型" style="width: 180px; margin-left: 10px" @change="handleSearch" clearable>
                    <el-option label="全部" value="" />
                    <el-option label="SQL语句" :value="1" />
                    <el-option label="视图/表" :value="2" />
                    <el-option label="存储过程" :value="3" />
                    <el-option label="关联" :value="4" />
                </el-select>
                <el-select v-model="filterConn" placeholder="筛选数据源" style="width: 180px; margin-left: 10px" @change="handleSearch" clearable>
                    <el-option label="全部" value="" />
                    <el-option v-for="conn in dataSources" :key="conn.id" :label="conn.connName" :value="conn.id" />
                </el-select>
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px;">
                    搜索
                </el-button>
            </div>
            <el-table v-loading="loading" :data="models" border style="width: 100%">
                <el-table-column prop="modelName" label="模型名称" width="200" />
                <el-table-column prop="modelCode" label="模型编码" width="180" />
                <el-table-column prop="modelVersion" label="版本" width="100" />
                <el-table-column prop="modelKind" label="模型类型" width="120">
                    <template #default="scope">
                        <el-tag :type="getModelKindTagType(scope.row.modelKind)">
                            {{ getModelKindText(scope.row.modelKind) }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="connName" label="数据源" width="150" />
                <el-table-column prop="isPublic" label="是否公开" width="100">
                    <template #default="scope">
                        <el-switch v-model="scope.row.isPublic" disabled active-color="#13ce66" inactive-color="#ff4d4f" />
                    </template>
                </el-table-column>
                <el-table-column prop="state" label="状态" width="100">
                    <template #default="scope">
                        <el-tag :type="scope.row.state === 1 ? 'success' : 'danger'">
                            {{ scope.row.state === 1 ? '有效' : '无效' }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="createAt" label="创建时间" width="180" />
                <el-table-column label="操作" width="200" fixed="right" class-name="action-column">
                    <template #default="scope">
                        <el-button type="primary" size="small" :icon="View" @click="handlePreview(scope.row)">
                            预览
                        </el-button>
                        <el-button type="success" size="small" :icon="Edit" @click="handleEdit(scope.row)" :disabled="scope.row.isLocked">
                            编辑
                        </el-button>
                        <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" :disabled="scope.row.isLocked">
                            删除
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
            <div class="pagination-container">
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import { getConns } from '@/api/metadata'
import { deleteModel, getModels } from '@/api/model'
import type { Model } from '@/types/metadata'
import type { DataSource } from '@/types/metadata/datasource'
import {
    Delete,
    Edit,
    Plus,
    Search,
    View
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const searchQuery = ref('')
const filterKind = ref('')
const filterConn = ref('')
const models = ref<Model[]>([])
const dataSources = ref<DataSource[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 生命周期钩子
onMounted(() => {
    fetchDataSources()
    fetchModels()
})

// 获取数据源列表（用于筛选）
const fetchDataSources = async () => {
    try {
        const response: any = await getConns()
        dataSources.value = response?.data || []
    } catch (error) {
        console.error('获取数据源列表失败:', error)
    }
}

// 获取模型列表
const fetchModels = async () => {
    loading.value = true
    try {
        const response: any = await getModels({
            page: currentPage.value,
            pageSize: pageSize.value,
            search: searchQuery.value,
            kind: filterKind.value ? Number(filterKind.value) : undefined
        })

        // 处理后端返回结构
        const data = response?.data || []
        models.value = data
        total.value = response?.total || data.length
    } catch (error) {
        console.error('获取模型列表失败:', error)
        ElMessage.error('获取模型列表失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => {
    currentPage.value = 1
    fetchModels()
}

const handleCurrentChange = (page: number) => {
    currentPage.value = page
    fetchModels()
}

const handleSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    fetchModels()
}

const handleCreate = () => router.push('/model/create')

const handleEdit = (row: any) => {
    // 适配 ID 类型
    const id = row.id || row.ID
    router.push(`/model/${id}/edit`)
}

const handleDelete = (row: any) => {
    const id = row.id || row.ID
    const name = row.modelName

    ElMessageBox.confirm(`确定要删除模型 "${name}" 吗？`, '删除确认', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(async () => {
        try {
            await deleteModel(Number(id))
            ElMessage.success('删除成功')
            fetchModels()
        } catch (error) {
            console.error('删除模型失败:', error)
        }
    })
}

const handlePreview = (row: any) => {
    const id = row.id || row.ID
    router.push(`/model/${id}/preview`)
}

// 获取模型类型文本
const getModelKindText = (kind: number): string => {
    const textMap: Record<number, string> = {
        1: 'SQL语句',
        2: '视图/表',
        3: '存储过程',
        4: '关联'
    }
    return textMap[kind] || '未知'
}

// 获取模型类型标签样式
const getModelKindTagType = (kind: number): string => {
    const typeMap: Record<number, string> = {
        1: 'primary',
        2: 'success',
        3: 'warning',
        4: 'info'
    }
    return typeMap[kind] || 'info'
}
</script>
<style scoped>
.model-list {
    padding: 10px;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.card-header {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
}

.pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
}

:deep(.action-column .cell) {
    white-space: nowrap;
}
</style>
