<template>
    <div class="container-padding">
        <!-- 页面标题区 -->
        <div class="page-header">
            <h1 class="page-title">
                <el-icon class="title-icon">
                    <List />
                </el-icon>
                字段列表
            </h1>
        </div>
        <!-- 主内容卡片 -->
        <el-card class="main-card">
            <!-- 搜索区域 -->
            <div class="search-area">
                <el-select v-model="selectedTable" placeholder="选择表" style="width: 240px" @change="handleTableChange" clearable>
                    <el-option v-for="table in tables" :key="table.id" :label="table.table_name" :value="table.id" />
                </el-select>
                <el-input v-model="searchQuery" placeholder="搜索字段名称" clearable :prefix-icon="Search" style="width: 300px; margin-left: 10px" />
                <el-button type="primary" :icon="Search" style="margin-left: 10px" @click="handleSearch">搜索</el-button>
                <el-button :icon="RefreshLeft" @click="handleReset">重置</el-button>
            </div>
            <!-- 表格区域 -->
            <div class="table-area">
                <el-table v-loading="loading" :data="pagedFields" border stripe style="width: 100%; height: 100%">
                    <el-table-column prop="column_name" label="名称" width="200" sortable />
                    <el-table-column prop="column_type" label="类型" width="150" />
                    <el-table-column prop="column_length" label="长度" width="100" />
                    <el-table-column prop="is_nullable" label="可为空" width="100">
                        <template #default="scope">
                            <el-tag :type="scope.row.is_nullable ? 'info' : 'danger'" size="small">
                                {{ scope.row.is_nullable ? 'Yes' : 'No' }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="is_primary_key" label="主键" width="80" align="center">
                        <template #default="scope">
                            <el-icon v-if="scope.row.is_primary_key" color="#E6A23C">
                                <Key />
                            </el-icon>
                        </template>
                    </el-table-column>
                    <el-table-column prop="column_comment" label="备注" min-width="200" />
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
import { getFieldsByTableId, getTables } from '@/api/metadata'
import type { MdTable, MdTableField } from '@/types/metadata'
import { Key, List, RefreshLeft, Search } from '@element-plus/icons-vue'
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// 响应式数据
const loading = ref(false)
const selectedTable = ref('')
const tables = ref<MdTable[]>([])
const allFields = ref<MdTableField[]>([])
const searchQuery = ref('')

// 分页状态
const currentPage = ref(1)
const pageSize = ref(10)
const total = computed(() => filteredFields.value.length)

// 计算属性 - 筛选
const filteredFields = computed(() => {
    if (!searchQuery.value) return allFields.value
    const query = searchQuery.value.toLowerCase()
    return allFields.value.filter(f =>
        f.column_name.toLowerCase().includes(query) ||
        (f.column_comment && f.column_comment.toLowerCase().includes(query))
    )
})

// 计算属性 - 分页数据
const pagedFields = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    return filteredFields.value.slice(start, end)
})

// 分页事件处理
const handleSizeChange = (val: number) => {
    pageSize.value = val
    currentPage.value = 1
}

const handleCurrentChange = (val: number) => {
    currentPage.value = val
}

// 搜索/重置处理
const handleSearch = () => {
    currentPage.value = 1
}

const handleReset = () => {
    searchQuery.value = ''
    selectedTable.value = ''
    currentPage.value = 1
    allFields.value = []
}

// 生命周期
onMounted(async () => {
    await fetchTables()

    // 如果 URL 中有 tableId 参数，则默认选中并加载
    if (route.query.tableId) {
        selectedTable.value = route.query.tableId as string
        fetchFields()
    }
})

// 获取表列表
const fetchTables = async () => {
    try {
        const res: any = await getTables('')
        tables.value = res?.data || []
    } catch (error) {
        console.error('获取表列表失败:', error)
    }
}

// 获取字段列表
const fetchFields = async () => {
    if (!selectedTable.value) {
        allFields.value = []
        return
    }
    loading.value = true
    try {
        const res: any = await getFieldsByTableId(selectedTable.value)
        if (Array.isArray(res)) {
            allFields.value = res
        } else {
            allFields.value = res?.data || []
        }
        currentPage.value = 1 // 重置页码
    } catch (error) {
        console.error('获取字段列表失败:', error)
    } finally {
        loading.value = false
    }
}

const handleTableChange = () => {
    fetchFields()
}

</script>
<style scoped>
/* ==================== 标准布局样式 ==================== */
.container-padding {
    padding-top: 20px;
    padding-bottom: 0;
    padding-left: 0;
    padding-right: 0;
    height: calc(100vh - 70px);
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
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
}

.pagination-area {
    flex-shrink: 0;
    display: flex;
    justify-content: flex-end;
}
</style>
