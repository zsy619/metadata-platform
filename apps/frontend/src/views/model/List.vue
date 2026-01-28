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
                <el-button type="primary" :icon="Plus" @click="handleCreate">新建模型</el-button>
            </div>
        </div>
        <!-- 主内容卡片 -->
        <el-card class="main-card">
            <!-- 搜索区域 -->
            <div class="search-area">
                <el-input v-model="searchQuery" placeholder="搜索模型名称或编码" clearable style="width: 300px" @keyup.enter="handleSearch" />
                <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
                <el-button :icon="RefreshLeft" @click="handleReset">重置</el-button>
            </div>
            <!-- 表格区域 -->
            <div class="table-area">
                <el-table v-loading="loading" :data="models" border style="width: 100%" height="100%">
                    <el-table-column prop="modelName" label="模型名称" min-width="150" sortable />
                    <el-table-column prop="modelCode" label="模型代码" min-width="150" sortable />
                    <el-table-column prop="remark" label="描述" min-width="200" show-overflow-tooltip />
                    <el-table-column prop="createAt" label="创建时间" width="180" />
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
import { Box, Delete, Plus, RefreshLeft, Search, View } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const models = ref<Model[]>([])
const searchQuery = ref('')
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 生命周期钩子
onMounted(() => {
    fetchModels()
})

// 获取模型列表
const fetchModels = async () => {
    loading.value = true
    try {
        const response: any = await getModels({
            page: currentPage.value,
            pageSize: pageSize.value,
            search: searchQuery.value
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

// 新建模型
const handleCreate = () => {
    router.push('/model/create')
}

// 编辑模型
const handleEdit = (row: Model) => {
    router.push(`/model/${row.id}/edit`)
}

// 删除模型
const handleDelete = (row: Model) => {
    ElMessageBox.confirm(`确定要删除模型 "${row.modelName}" 吗？`, '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(async () => {
        try {
            await deleteModel(row.id)
            ElMessage.success('删除成功')
            fetchModels()
        } catch (error) {
            console.error('删除模型失败:', error)
        }
    }).catch(() => { })
}
</script>
<style scoped>
/* ==================== 标准布局样式 ==================== */
.container-padding {
    padding: 20px;
    padding-bottom: 0;
    height: calc(100vh - 84px);
    display: flex;
    flex-direction: column;
    overflow: hidden;
    box-sizing: border-box;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
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
