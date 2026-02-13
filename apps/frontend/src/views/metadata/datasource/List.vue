<template>
    <div class="data-source-list container-padding">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon">
                    <DataLine />
                </el-icon>
                数据源管理
            </h1>
            <div class="header-actions">
                <el-button type="danger" :icon="Delete" @click="handleBatchDelete" :disabled="selectedRows.length === 0">
                    批量删除
                </el-button>
                <el-button type="primary" @click="handleCreate" :icon="Plus">
                    新增数据源
                </el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-input v-model="searchQuery" placeholder="请输入数据源名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
                <el-select v-model="filterType" placeholder="筛选数据源类型" style="width: 180px; margin-left: 10px" clearable @change="handleSearch">
                    <el-option label="全部" value="" />
                    <el-option-group label="关系型">
                        <el-option label="MySQL" value="MySQL" />
                        <el-option label="PostgreSQL" value="PostgreSQL" />
                        <el-option label="SQL Server" value="SQL Server" />
                        <el-option label="Oracle" value="Oracle" />
                        <el-option label="TiDB" value="TiDB" />
                        <el-option label="OceanBase" value="OceanBase" />
                        <el-option label="SQLite" value="SQLite" />
                    </el-option-group>
                    <el-option-group label="大数据/分析">
                        <el-option label="ClickHouse" value="ClickHouse" />
                        <el-option label="Doris" value="Doris" />
                        <el-option label="StarRocks" value="StarRocks" />
                    </el-option-group>
                    <el-option-group label="国产化">
                        <el-option label="Dameng (DM)" value="DM" />
                        <el-option label="Kingbase" value="Kingbase" />
                        <el-option label="OpenGauss" value="OpenGauss" />
                    </el-option-group>
                    <el-option-group label="NoSQL">
                        <el-option label="MongoDB" value="MongoDB" />
                        <el-option label="Redis" value="Redis" />
                    </el-option-group>
                </el-select>
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">
                    搜索
                </el-button>
                <el-button @click="handleReset" :icon="RefreshLeft">
                    重置
                </el-button>
            </div>
            <div class="table-area">
                <el-table v-loading="loading" :element-loading-text="loadingText" :data="pagedDataSources" border stripe style="width: 100%; height: 100%;" @selection-change="handleSelectionChange">
                    <template #empty>
                        <el-empty :description="searchQuery ? '未搜索到相关数据源' : '暂无数据源'">
                            <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增数据源</el-button>
                        </el-empty>
                    </template>
                    <el-table-column type="selection" width="55" />
                    <el-table-column prop="conn_name" label="数据源名称" width="200" show-overflow-tooltip />
                    <el-table-column prop="conn_kind" label="类型" width="120" />
                    <el-table-column prop="conn_host" label="主机" width="180" show-overflow-tooltip />
                    <el-table-column prop="conn_port" label="端口" width="80" />
                    <el-table-column prop="conn_database" label="数据库" width="150" show-overflow-tooltip />
                    <el-table-column prop="state" label="状态" width="100">
                        <template #default="scope">
                            <el-tag v-if="scope.row.status === 1" type="success">有效</el-tag>
                            <el-tag v-else-if="scope.row.status === 2" type="danger">连接错误</el-tag>
                            <el-tag v-else type="info">未检测</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="update_at" label="更新时间" width="170">
                        <template #default="scope">
                            {{ formatDateTime(scope.row.update_at || scope.row.create_at) }}
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="320" fixed="right">
                        <template #default="scope">
                            <el-button type="success" size="small" :icon="Connection" @click="handleTestConnection(scope.row)" :disabled="scope.row.status === 1" plain>
                                测试
                            </el-button>
                            <el-button type="warning" size="small" :icon="Folder" @click="handleBrowse(scope.row)" plain>
                                浏览
                            </el-button>
                            <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>
                                编辑
                            </el-button>
                            <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg>
                                删除
                            </el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
            <div class="pagination-area">
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]" background layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
        </el-card>
        <!-- Object Browser Dialog -->
        <el-dialog v-model="browserVisible" :title="`浏览数据源: ${currentConn?.conn_name}`" width="90%" top="5vh" custom-class="browser-dialog" destroy-on-close>
            <div class="browser-container">
                <div class="browser-sidebar" :style="{ width: browserSidebarWidth + 'px' }">
                    <ObjectBrowser v-if="currentConn" :data-source-id="currentConn.id" @select-table="handleSelectTable" @select-view="handleSelectView" />
                </div>
                <div class="resize-handle" @mousedown="startBrowserResize"></div>
                <div class="browser-main">
                    <DataPreview v-if="selectedTable && currentConn" :conn-id="currentConn.id" :table-name="selectedTable" />
                    <el-empty v-else description="请从左侧选择表或视图查看数据" />
                </div>
            </div>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { deleteConn, getConns, testConn } from '@/api/metadata'
import DataPreview from '@/components/DataPreview.vue'
import ObjectBrowser from '@/components/ObjectBrowser.vue'
import type { MdConn } from '@/types/metadata'
import { showDeleteConfirm } from '@/utils/confirm'
import { debounce } from '@/utils/debounce'
import {
    Connection,
    DataLine,
    Delete,
    Edit,
    Folder,
    Plus,
    RefreshLeft,
    Search
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const loadingText = ref('加载中...')
const dataSources = ref<MdConn[]>([])
const searchQuery = ref('')
const filterType = ref('')
const selectedRows = ref<MdConn[]>([])

// Browser Dialog
const browserVisible = ref(false)
const currentConn = ref<MdConn>()
const selectedTable = ref('')
const browserSidebarWidth = ref(300)
const isBrowserResizing = ref(false)

const startBrowserResize = (e: MouseEvent) => {
    isBrowserResizing.value = true
    document.addEventListener('mousemove', doBrowserResize)
    document.addEventListener('mouseup', stopBrowserResize)
    document.body.style.cursor = 'col-resize'
    document.body.style.userSelect = 'none'
}

const doBrowserResize = (e: MouseEvent) => {
    if (!isBrowserResizing.value) return
    // 对于 Dialog，我们简化计算，因为它是居中的
    // 但更稳妥的是基于 Offset
    const container = document.querySelector('.browser-container')
    if (container) {
        const rect = container.getBoundingClientRect()
        const newWidth = e.clientX - rect.left
        if (newWidth > 200 && newWidth < 600) {
            browserSidebarWidth.value = newWidth
        }
    }
}

const stopBrowserResize = () => {
    isBrowserResizing.value = false
    document.removeEventListener('mousemove', doBrowserResize)
    document.removeEventListener('mouseup', stopBrowserResize)
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
}

// Pagination
const currentPage = ref(1)
const pageSize = ref(20)
const total = computed(() => filteredDataSources.value.length)

const filteredDataSources = computed(() => {
    return dataSources.value.filter(item => {
        const matchName = item.conn_name.toLowerCase().includes(searchQuery.value.toLowerCase())
        const matchType = filterType.value ? item.conn_kind === filterType.value : true
        return matchName && matchType
    })
})

const pagedDataSources = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    return filteredDataSources.value.slice(start, end)
})

const handleSizeChange = (val: number) => {
    pageSize.value = val
    currentPage.value = 1
}

const handleCurrentChange = (val: number) => {
    currentPage.value = val
}

const formatDateTime = (dateStr: string | undefined) => {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    return isNaN(date.getTime()) ? '-' : date.toLocaleString()
}

const fetchDataSources = async () => {
    loadingText.value = '加载中...'
    loading.value = true
    try {
        const res: any = await getConns()
        dataSources.value = res.data || res
    } catch (error) {
        console.error('Fetch error', error)
        ElMessage.error('加载列表失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => {
    // Computed property handles filtering
    currentPage.value = 1
}

const handleDebouncedSearch = debounce(handleSearch, 300)

const handleReset = () => {
    searchQuery.value = ''
    filterType.value = ''
    currentPage.value = 1
}

const handleSelectionChange = (val: MdConn[]) => {
    selectedRows.value = val
}

const handleCreate = () => router.push('/metadata/datasource/create')
const handleEdit = (row: MdConn) => router.push(`/metadata/datasource/${row.id}/edit`)

const handleDelete = (row: MdConn) => {
    showDeleteConfirm(`确定要删除 "${row.conn_name}" 吗？`).then(async () => {
        try {
            await deleteConn(row.id)
            ElMessage.success('删除成功')
            fetchDataSources()
        } catch (error) {
            console.error(error)
        }
    })
}

const handleBatchDelete = () => {
    showDeleteConfirm(`确定要删除选中的 ${selectedRows.value.length} 个数据源吗？`).then(async () => {
        try {
            await Promise.all(selectedRows.value.map(row => deleteConn(row.id)))
            ElMessage.success('批量删除成功')
            fetchDataSources()
        } catch (error) {
            console.error(error)
        }
    })
}

const handleTestConnection = async (row: MdConn) => {
    loadingText.value = `正在测试与 "${row.conn_name}" 的连接...`
    loading.value = true
    try {
        const res = await testConn(row.id)
        // 后端返回格式: { code: 200, message: "success", data: "连接成功" }
        if (res && (res as any).code === 200) {
            ElMessage.success('连接成功')
        } else {
            ElMessage.error((res as any)?.message || '连接失败')
        }
    } catch (error: any) {
        ElMessage.error(error.message || '连接测试失败')
    } finally {
        loading.value = false
        // 无论成功还是失败，都刷新列表以更新状态
        await fetchDataSources()
    }
}

const handleBrowse = (row: MdConn) => {
    currentConn.value = row
    selectedTable.value = ''
    browserVisible.value = true
}

const handleSelectTable = ({ tableName }: { tableName: string }) => {
    selectedTable.value = tableName
}
const handleSelectView = ({ viewName }: { viewName: string }) => {
    selectedTable.value = viewName // Reuse generic table viewer for now
}

onMounted(() => {
    fetchDataSources()
})
</script>
<style scoped>
/* ==================== 标准布局样式 ==================== */
.main-card {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    /* Ensure card body also uses flex */
}

/* Deep selector for card body to allow flex content */
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
}

.table-area {
    flex: 1;
    overflow: hidden;
    margin-bottom: 20px;
    /* Ensure table takes full internal height */
}

.pagination-area {
    flex-shrink: 0;
    display: flex;
    justify-content: flex-end;
    margin-top: auto;
    /* Push to bottom if space allows, though flex direction with table flex 1 handles it */
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
    /* Space between icon and text */
}

.title-icon {
    font-size: 24px;
    /* Slightly larger icon */
}

.browser-container {
    display: flex;
    height: 60vh;
    border: 1px solid var(--el-border-color);
}

.browser-sidebar {
    flex-shrink: 0;
    border-right: 1px solid var(--el-border-color);
    background-color: #f8f9fa;
}

.resize-handle {
    width: 6px;
    height: 100%;
    cursor: col-resize;
    background-color: transparent;
    transition: background-color 0.2s;
    flex-shrink: 0;
    position: relative;
    z-index: 10;
}

.resize-handle:hover,
.resize-handle:active {
    background-color: var(--el-color-primary-light-8);
}

.resize-handle::after {
    content: '';
    position: absolute;
    left: 2px;
    top: 50%;
    transform: translateY(-50%);
    width: 1px;
    height: 30px;
    background-color: var(--el-border-color);
}

.browser-main {
    flex: 1;
    overflow: hidden;
    background-color: #fff;
}

/* Deep selector for dialog body to maximize space */
:deep(.browser-dialog .el-dialog__body) {
    padding: 0;
}
</style>