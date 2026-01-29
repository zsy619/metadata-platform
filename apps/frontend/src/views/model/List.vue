<template>
    <div class="container-padding">
        <!-- 页面标题区 -->
        <div class="page-header">
            <h1 class="page-title">
                <el-icon class="title-icon">
                    <Box />
                </el-icon>
                模型管理
            </h1>
            <div class="header-actions">
                <el-dropdown @command="handleCreateCommand" trigger="click">
                    <el-button type="primary">
                        <el-icon>
                            <Plus />
                        </el-icon>
                        新建模型
                        <el-icon class="el-icon--right"><arrow-down /></el-icon>
                    </el-button>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item :command="1">SQL 语句</el-dropdown-item>
                            <el-dropdown-item :command="2">视图 / 表</el-dropdown-item>
                            <el-dropdown-item :command="3">存储过程</el-dropdown-item>
                            <el-dropdown-item :command="4">关联</el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </div>
        </div>
        <!-- 主内容卡片 -->
        <el-card class="main-card">
            <!-- 搜索区域 -->
            <div class="search-area">
                <el-input v-model="searchQuery" placeholder="搜索模型名称或编码" clearable style="width: 260px" @input="handleDebouncedSearch" @keyup.enter="handleSearch" />
                <el-select v-model="filterKind" placeholder="筛选模型类型" style="width: 160px" clearable @change="handleSearch">
                    <el-option label="全部类型" :value="0" />
                    <el-option label="SQL 语句" :value="1" />
                    <el-option label="视图 / 表" :value="2" />
                    <el-option label="存储过程" :value="3" />
                    <el-option label="模型关联" :value="4" />
                </el-select>
                <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
                <el-button :icon="RefreshLeft" @click="handleReset">重置</el-button>
            </div>
            <!-- 表格区域 -->
            <div class="table-area">
                <el-table v-loading="loading" :element-loading-text="loadingText" :data="models" border style="width: 100%" height="100%">
                    <template #empty>
                        <el-empty :description="searchQuery || filterKind ? '未搜索到相关模型' : '暂无模型数据'" />
                    </template>
                    <el-table-column prop="model_name" label="模型名称" min-width="150" sortable />
                    <el-table-column prop="model_code" label="模型代码" min-width="180" sortable />
                    <el-table-column prop="model_kind" label="模型类型" width="120">
                        <template #default="scope">
                            <el-tag v-if="scope.row.model_kind === 1">SQL 语句</el-tag>
                            <el-tag v-else-if="scope.row.model_kind === 2" type="success">视图 / 表</el-tag>
                            <el-tag v-else-if="scope.row.model_kind === 3" type="warning">存储过程</el-tag>
                            <el-tag v-else-if="scope.row.model_kind === 4" type="info">关联</el-tag>
                            <el-tag v-else type="info">未知</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="remark" label="描述" min-width="200" show-overflow-tooltip />
                    <el-table-column prop="create_at" label="创建时间" width="160">
                        <template #default="scope">
                            {{ formatDateTime(scope.row.create_at) }}
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="200" fixed="right" class-name="action-column">
                        <template #default="scope">
                            <el-button type="primary" link :icon="View" @click="handleEdit(scope.row)">编辑</el-button>
                            <el-button type="danger" link :icon="Delete" @click="handleDelete(scope.row)">删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
            <!-- 分页区域 -->
            <div class="pagination-area">
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]" background layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import { deleteModel, getModels } from '@/api/model'
import type { Model } from '@/types/metadata'
import { showDeleteConfirm } from '@/utils/confirm'
import { debounce } from '@/utils/debounce'
import { ArrowDown, Box, Delete, Plus, RefreshLeft, Search, View } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const loadingText = ref('加载中...')
const models = ref<Model[]>([])
const searchQuery = ref('')
const filterKind = ref(0)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 格式化日期
const formatDateTime = (dateStr: string) => {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return dateStr
    const Y = date.getFullYear()
    const M = String(date.getMonth() + 1).padStart(2, '0')
    const D = String(date.getDate()).padStart(2, '0')
    const h = String(date.getHours()).padStart(2, '0')
    const m = String(date.getMinutes()).padStart(2, '0')
    return `${Y}-${M}-${D} ${h}:${m}`
}

// 生命周期钩子
onMounted(() => {
    fetchModels()
})

// 获取模型列表
const fetchModels = async () => {
    loadingText.value = '加载中...'
    loading.value = true
    try {
        const response: any = await getModels({
            page: currentPage.value,
            pageSize: pageSize.value,
            search: searchQuery.value,
            model_kind: filterKind.value
        })

        // 处理后端 SuccessWithPagination 返回结构 { code: 200, data: { list: [], total: 0 } }
        const paginatedData = response?.data || {}
        if (Array.isArray(paginatedData)) {
            // 兼容性处理：如果 data 直接是数组
            models.value = paginatedData
            total.value = paginatedData.length
        } else {
            models.value = paginatedData.list || []
            total.value = paginatedData.total || 0
        }
    } catch (error) {
        console.error('获取模型列表失败:', error)
        ElMessage.error('获取模型列表失败')
    } finally {
        loading.value = false
    }
}

// 搜索
const handleSearch = () => {
    currentPage.value = 1
    fetchModels()
}

// 重置
const handleReset = () => {
    searchQuery.value = ''
    filterKind.value = 0
    currentPage.value = 1
    fetchModels()
}

// 页码变化
const handleCurrentChange = (val: number) => {
    currentPage.value = val
    fetchModels()
}

// 每页条数变化
const handleSizeChange = (val: number) => {
    pageSize.value = val
    currentPage.value = 1
    fetchModels()
}

const handleDebouncedSearch = debounce(handleSearch, 500)

// 新建模型
const handleCreateCommand = (kind: number) => {
    if (kind === 1) {
        router.push('/metadata/model/create-sql')
    } else if (kind === 2) {
        router.push('/metadata/model/visual-create')
    } else {
        router.push({
            path: '/metadata/model/create',
            query: { kind }
        })
    }
}

// 编辑模型
const handleEdit = (row: Model) => {
    const kind = row.modelKind ?? row.model_kind
    if (kind === 1) {
        router.push(`/metadata/model/edit-sql/${row.id}`)
    } else {
        router.push(`/metadata/model/${row.id}/edit`)
    }
}

// 删除模型
const handleDelete = (row: Model) => {
    showDeleteConfirm(`确定要删除模型 "${row.model_name}" 吗？`).then(async () => {
        loadingText.value = '正在删除...'
        loading.value = true
        try {
            await deleteModel(row.id)
            ElMessage.success('删除成功')
            // fetchModels will reset loading text to '加载中...'
            fetchModels()
        } catch (error) {
            console.error('删除模型失败:', error)
            loading.value = false // only reset if error, success handled by fetch
        }
    })
}
</script>
<style scoped>
/* ==================== 标准布局样式 ==================== */
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
    margin-bottom: 20px;
    position: relative;
}

.pagination-area {
    flex-shrink: 0;
    display: flex;
    justify-content: flex-end;
}

:deep(.action-column .cell) {
    white-space: nowrap;
}
</style>
