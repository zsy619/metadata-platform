<template>
    <div class="data-list container-padding">
        <div class="page-header">
            <h1 class="page-title">
                <el-icon class="title-icon">
                    <Grid />
                </el-icon>
                {{ pageTitle }}
            </h1>
            <div class="header-actions">
                <el-button type="primary" :icon="Plus" @click="handleCreate">新增数据</el-button>
                <el-button :icon="Download" @click="handleExport">导出</el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="filter-area">
                <div class="filter-row">
                    <div class="quick-filters">
                        <span class="filter-label">快速筛选:</span>
                        <el-radio-group v-model="quickFilter" size="small" @change="handleQuickFilter">
                            <el-radio-button label="">全部</el-radio-button>
                            <el-radio-button label="1">启用</el-radio-button>
                            <el-radio-button label="0">禁用</el-radio-button>
                        </el-radio-group>
                    </div>
                    <div class="search-box">
                        <el-input v-model="searchQuery" placeholder="搜索..." clearable :prefix-icon="Search" style="width: 250px" @clear="handleSearch" @keyup.enter="handleSearch" />
                        <el-button type="primary" @click="handleSearch">搜索</el-button>
                        <el-button :icon="Filter" @click="showAdvancedFilter = !showAdvancedFilter">高级筛选</el-button>
                    </div>
                </div>
                <div v-show="showAdvancedFilter" class="advanced-filter">
                    <QueryBuilder :fields="filterFields" v-model="advancedConditions" :show-sql-preview="false" @change="handleAdvancedFilter" />
                </div>
            </div>
            <div class="toolbar-area">
                <div class="column-selector">
                    <el-popover placement="bottom" :width="300" trigger="click">
                        <template #reference>
                            <el-button size="small">
                                <el-icon><Setting /></el-icon>
                                列设置
                            </el-button>
                        </template>
                        <div class="column-list">
                            <el-checkbox v-model="checkAll" :indeterminate="isIndeterminate" @change="handleCheckAll">全选</el-checkbox>
                            <el-divider style="margin: 10px 0" />
                            <el-checkbox-group v-model="visibleColumns" @change="handleVisibleColumnsChange">
                                <el-checkbox v-for="col in allColumns" :key="col.prop" :label="col.prop">{{ col.label }}</el-checkbox>
                            </el-checkbox-group>
                        </div>
                    </el-popover>
                </div>
                <div class="template-selector">
                    <el-select v-model="selectedTemplate" placeholder="应用查询模板" clearable size="small" @change="handleTemplateChange">
                        <el-option v-for="tpl in queryTemplates" :key="tpl.id" :label="tpl.name" :value="tpl.id" />
                    </el-select>
                </div>
            </div>
            <div class="table-area">
                <el-table v-loading="loading" :data="tableData" border stripe height="100%" @selection-change="handleSelectionChange" @sort-change="handleSortChange" @row-click="handleRowClick">
                    <el-table-column v-if="showSelection" type="selection" width="50" />
                    <el-table-column v-if="showIndex" type="index" label="序号" width="60" align="center" />
                    <el-table-column v-for="col in displayColumns" :key="col.prop" :prop="col.prop" :label="col.label" :width="col.width" :min-width="col.minWidth" :sortable="col.sortable ? 'custom' : false" show-overflow-tooltip>
                        <template #default="{ row }">
                            <template v-if="col.type === 'switch'">
                                <el-switch v-model="row[col.prop]" :active-value="1" :inactive-value="0" @change="handleFieldChange(row, col)" />
                            </template>
                            <template v-else-if="col.type === 'tag'">
                                <el-tag :type="getTagType(row[col.prop], col.tagTypes)">{{ getTagText(row[col.prop], col.tagOptions) }}</el-tag>
                            </template>
                            <template v-else-if="col.type === 'image'">
                                <el-image v-if="row[col.prop]" :src="row[col.prop]" fit="cover" style="width: 40px; height: 40px" preview-src-list [src]="row[col.prop]" />
                            </template>
                            <template v-else-if="col.type === 'datetime'">
                                {{ formatDate(row[col.prop]) }}
                            </template>
                            <template v-else>
                                {{ row[col.prop] ?? '-' }}
                            </template>
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" :width="actionWidth" fixed="right" align="center">
                        <template #default="{ row }">
                            <el-button type="primary" link size="small" :icon="View" @click="handleView(row)">查看</el-button>
                            <el-button type="success" link size="small" :icon="Edit" @click="handleEdit(row)">编辑</el-button>
                            <el-button type="danger" link size="small" :icon="Delete" @click="handleDelete(row)">删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
            <div class="pagination-area">
                <div class="selection-info" v-if="selectedRows.length > 0">
                    <span>已选择 {{ selectedRows.length }} 项</span>
                    <el-button type="danger" size="small" @click="handleBatchDelete">批量删除</el-button>
                </div>
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
        </el-card>
        <DataForm ref="dataFormRef" :model-id="modelId" :fields="formFields" :visible="formVisible" @update:visible="formVisible = $event" @success="handleFormSuccess" />
    </div>
</template>
<script setup lang="ts">
import { Delete, Download, Edit, Filter, Grid, Plus, Search, Setting, View } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import QueryBuilder from '@/components/query/QueryBuilder.vue'
import DataForm from '@/components/data/DataForm.vue'

interface Column {
    prop: string
    label: string
    width?: number | string
    minWidth?: number | string
    sortable?: boolean
    type?: string
    tagTypes?: Record<string, string>
    tagOptions?: Record<string, string>
}

interface QueryTemplate {
    id: number
    name: string
}

const route = useRoute()
const router = useRouter()

const modelId = computed(() => Number(route.params.modelId) || 1)
const pageTitle = computed(() => route.query.title as string || '数据查询')

const loading = ref(false)
const searchQuery = ref('')
const quickFilter = ref('')
const showAdvancedFilter = ref(false)
const advancedConditions = ref([])
const selectedTemplate = ref<number>()
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const tableData = ref<any[]>([])
const selectedRows = ref<any[]>([])
const formVisible = ref(false)
const dataFormRef = ref()

const showSelection = ref(true)
const showIndex = ref(true)
const actionWidth = ref(200)

const allColumns = ref<Column[]>([
    { prop: 'id', label: 'ID', width: 80, sortable: true },
    { prop: 'username', label: '用户名', minWidth: 120, sortable: true },
    { prop: 'email', label: '邮箱', minWidth: 150 },
    { prop: 'phone', label: '手机号', minWidth: 120 },
    { prop: 'status', label: '状态', width: 80, type: 'switch', tagTypes: { 1: 'success', 0: 'info' }, tagOptions: { 1: '启用', 0: '禁用' } },
    { prop: 'createTime', label: '创建时间', width: 180, type: 'datetime' }
])

const visibleColumns = ref(allColumns.value.map(c => c.prop))

const displayColumns = computed(() => allColumns.value.filter(c => visibleColumns.value.includes(c.prop)))

const filterFields = computed(() => allColumns.value.map(c => ({
    name: c.prop,
    label: c.label,
    type: c.prop === 'status' ? 'int' : 'string'
})))

const formFields = computed(() => allColumns.value.filter(c => c.prop !== 'id').map(c => ({
    prop: c.prop,
    label: c.label,
    type: c.prop === 'status' ? 'switch' : 'input',
    required: c.prop === 'username'
})))

const checkAll = computed({
    get: () => visibleColumns.value.length === allColumns.value.length,
    set: (val) => visibleColumns.value = val ? allColumns.value.map(c => c.prop) : []
})

const isIndeterminate = computed(() => visibleColumns.value.length > 0 && visibleColumns.value.length < allColumns.value.length)

const queryTemplates = ref<QueryTemplate[]>([
    { id: 1, name: '用户查询模板' },
    { id: 2, name: 'VIP用户模板' }
])

const loadData = async () => {
    loading.value = true
    try {
        await new Promise(resolve => setTimeout(resolve, 500))
        tableData.value = Array.from({ length: pageSize.value }, (_, i) => ({
            id: (currentPage.value - 1) * pageSize.value + i + 1,
            username: `user${i + 1}`,
            email: `user${i + 1}@example.com`,
            phone: `138${String(i + 1).padStart(8, '0')}`,
            status: i % 2,
            createTime: new Date().toISOString()
        }))
        total.value = 100
    } finally {
        loading.value = false
    }
}

watch([modelId], () => loadData(), { immediate: true })

const handleSearch = () => {
    currentPage.value = 1
    loadData()
}

const handleQuickFilter = () => {
    currentPage.value = 1
    loadData()
}

const handleAdvancedFilter = () => {
    currentPage.value = 1
    loadData()
}

const handleTemplateChange = () => {
    loadData()
}

const handleCheckAll = (val: boolean) => {
    visibleColumns.value = val ? allColumns.value.map(c => c.prop) : []
}

const handleVisibleColumnsChange = () => {}

const handleSelectionChange = (rows: any[]) => {
    selectedRows.value = rows
}

const handleSortChange = ({ prop, order }: any) => {
    loadData()
}

const handleRowClick = (row: any) => {}

const handleFieldChange = (row: any, col: Column) => {
    ElMessage.success(`字段 ${col.label} 已更新`)
}

const handleCreate = () => {
    formVisible.value = true
}

const handleView = (row: any) => {
    router.push(`/data/${modelId.value}/${row.id}`)
}

const handleEdit = (row: any) => {
    dataFormRef.value?.setData(row)
    formVisible.value = true
}

const handleDelete = (row: any) => {
    ElMessageBox.confirm(`确定要删除这条数据吗？`, '提示', { type: 'warning' }).then(async () => {
        ElMessage.success('删除成功')
        loadData()
    })
}

const handleBatchDelete = () => {
    ElMessageBox.confirm(`确定要删除选中的 ${selectedRows.value.length} 条数据吗？`, '提示', { type: 'warning' }).then(async () => {
        ElMessage.success('批量删除成功')
        loadData()
    })
}

const handleExport = () => {
    ElMessage.info('导出功能开发中...')
}

const handleSizeChange = (size: number) => {
    pageSize.value = size
    loadData()
}

const handleCurrentChange = (page: number) => {
    currentPage.value = page
    loadData()
}

const handleFormSuccess = () => {
    formVisible.value = false
    loadData()
}

const getTagType = (value: any, tagTypes?: Record<string, string>) => tagTypes?.[value] || 'info'
const getTagText = (value: any, tagOptions?: Record<string, string>) => tagOptions?.[value] ?? value
const formatDate = (value: string) => value ? new Date(value).toLocaleString() : '-'
</script>
<style scoped>
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

.header-actions {
    display: flex;
    gap: 10px;
}

.filter-area {
    margin-bottom: 16px;
}

.filter-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
}

.quick-filters {
    display: flex;
    align-items: center;
    gap: 12px;
}

.filter-label {
    color: #606266;
    font-size: 14px;
}

.search-box {
    display: flex;
    gap: 8px;
}

.advanced-filter {
    padding: 16px;
    background: #fafafa;
    border-radius: 8px;
    margin-bottom: 12px;
}

.toolbar-area {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
}

.column-list {
    max-height: 300px;
    overflow-y: auto;
}

.table-area {
    margin-bottom: 16px;
}

.pagination-area {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.selection-info {
    display: flex;
    align-items: center;
    gap: 12px;
    color: #606266;
}
</style>
