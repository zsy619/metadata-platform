<template>
    <div class="data-source-list container-padding">
        <div class="page-header">
            <h1 class="text-primary">数据源管理</h1>
            <div class="header-actions">
                <el-button type="danger" :icon="Delete" @click="handleBatchDelete" :disabled="selectedRows.length === 0">
                    批量删除
                </el-button>
                <el-button type="primary" @click="handleCreate" :icon="Plus">
                    新增数据源
                </el-button>
            </div>
        </div>
        <el-card>
            <div class="flex-center m-b-lg">
                <el-input v-model="searchQuery" placeholder="请输入数据源名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleSearch" />
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
                <el-button type="primary" @click="fetchDataSources" :icon="Refresh" style="margin-left: 10px">
                    刷新
                </el-button>
            </div>
            <el-table v-loading="loading" :data="filteredDataSources" border stripe style="width: 100%" @selection-change="handleSelectionChange">
                <el-table-column type="selection" width="55" />
                <el-table-column prop="conn_name" label="数据源名称" width="200" show-overflow-tooltip />
                <el-table-column prop="conn_kind" label="类型" width="120" />
                <el-table-column prop="conn_host" label="主机" width="180" show-overflow-tooltip />
                <el-table-column prop="conn_port" label="端口" width="80" />
                <el-table-column prop="conn_database" label="数据库" width="150" show-overflow-tooltip />
                <el-table-column prop="state" label="状态" width="100">
                    <template #default="scope">
                        <el-tag :type="scope.row.state === 1 ? 'success' : 'danger'">
                            {{ scope.row.state === 1 ? '有效' : '未检测' }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="update_at" label="更新时间" width="170">
                    <template #default="scope">
                        {{ formatDateTime(scope.row.update_at || scope.row.create_at) }}
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="320" fixed="right">
                    <template #default="scope">
                        <el-button type="success" size="small" :icon="Connection" @click="handleTestConnection(scope.row)" :disabled="scope.row.state === 1" plain>
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
            <div class="flex-between m-t-lg">
                <span class="text-secondary">共 {{ filteredDataSources.length }} 条记录</span>
                <!-- Pagination could be added here if backend supports it -->
            </div>
        </el-card>
        <!-- Object Browser Dialog -->
        <el-dialog v-model="browserVisible" :title="`浏览数据源: ${currentConn?.conn_name}`" width="80%" top="5vh" custom-class="browser-dialog" destroy-on-close>
            <div class="browser-container">
                <div class="browser-sidebar">
                    <ObjectBrowser v-if="currentConn" :data-source-id="currentConn.id" @select-table="handleSelectTable" @select-view="handleSelectView" />
                </div>
                <div class="browser-main">
                    <DataPreview v-if="selectedTable" :conn-id="currentConn?.id" :table-name="selectedTable" />
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
import {
    Connection,
    Delete, Edit,
    Folder,
    Plus,
    Refresh,
    Search
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const dataSources = ref<MdConn[]>([])
const searchQuery = ref('')
const filterType = ref('')
const selectedRows = ref<MdConn[]>([])

// Browser Dialog
const browserVisible = ref(false)
const currentConn = ref<MdConn>()
const selectedTable = ref('')

const filteredDataSources = computed(() => {
    return dataSources.value.filter(item => {
        const matchName = item.conn_name.toLowerCase().includes(searchQuery.value.toLowerCase())
        const matchType = filterType.value ? item.conn_kind === filterType.value : true
        return matchName && matchType
    })
})

const formatDateTime = (dateStr: string | undefined) => {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    return isNaN(date.getTime()) ? '-' : date.toLocaleString()
}

const fetchDataSources = async () => {
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
}

const handleSelectionChange = (val: MdConn[]) => {
    selectedRows.value = val
}

const handleCreate = () => router.push('/metadata/datasource/create')
const handleEdit = (row: MdConn) => router.push(`/metadata/datasource/${row.id}/edit`)

const handleDelete = (row: MdConn) => {
    ElMessageBox.confirm(`确定要删除 "${row.conn_name}" 吗？`, '警告', {
        type: 'warning'
    }).then(async () => {
        try {
            await deleteConn(row.id)
            ElMessage.success('删除成功')
            fetchDataSources()
        } catch (error) {
            ElMessage.error('删除失败')
        }
    })
}

const handleBatchDelete = () => {
    ElMessageBox.confirm(`确定要删除选中的 ${selectedRows.value.length} 个数据源吗？`, '警告', {
        type: 'warning'
    }).then(async () => {
        try {
            await Promise.all(selectedRows.value.map(row => deleteConn(row.id)))
            ElMessage.success('批量删除成功')
            fetchDataSources()
        } catch (error) {
            ElMessage.error('部分删除失败，请刷新重试')
        }
    })
}

const handleTestConnection = async (row: MdConn) => {
    loading.value = true
    try {
        const res = await testConn(row.id)
        if (res && res.success) { // adjust check based on actual API response
            ElMessage.success('连接成功')
            row.state = 1
        } else {
            ElMessage.error(res?.message || '连接失败')
            row.state = 0
        }
    } catch (error: any) {
        ElMessage.error(error.message || '连接测试失败')
    } finally {
        loading.value = false
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
.container-padding {
    padding: 20px;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.browser-container {
    display: flex;
    height: 60vh;
    border: 1px solid var(--el-border-color);
}

.browser-sidebar {
    width: 300px;
    border-right: 1px solid var(--el-border-color);
}

.browser-main {
    flex: 1;
    overflow: hidden;
}

/* Deep selector for dialog body to maximize space */
:deep(.browser-dialog .el-dialog__body) {
    padding: 0;
}
</style>