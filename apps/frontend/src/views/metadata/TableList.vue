<template>
    <div class="table-list">
        <div class="page-header">
            <h1 class="text-primary">表与视图</h1>
        </div>
        <el-card>
            <div class="filter-bar m-b-lg">
                <div class="filter-left">
                    <el-select v-model="selectedConn" placeholder="选择数据源" style="width: 240px" @change="handleConnChange">
                        <el-option v-for="conn in connections" :key="conn.id" :label="conn.connName" :value="conn.id" />
                    </el-select>
                    <el-input v-model="searchQuery" placeholder="搜索表或视图名称" clearable :prefix-icon="Search" style="width: 300px; margin-left: 10px" @clear="fetchTables" @keyup.enter="fetchTables" />
                </div>
                <div class="filter-right">
                    <el-dropdown trigger="click" @command="openSelectDialog" style="margin-right: 12px">
                        <el-button type="primary" plain>
                            元数据导入<el-icon class="el-icon--right"><arrow-down /></el-icon>
                        </el-button>
                        <template #dropdown>
                            <el-dropdown-menu>
                                <el-dropdown-item command="TABLE">选择表</el-dropdown-item>
                                <el-dropdown-item command="VIEW">选择视图</el-dropdown-item>
                            </el-dropdown-menu>
                        </template>
                    </el-dropdown>
                    <el-button type="primary" :icon="Refresh" @click="fetchTables">刷新</el-button>
                </div>
            </div>
            <el-table v-loading="loading" :data="filteredTables" border stripe style="width: 100%">
                <el-table-column prop="tableName" label="名称" min-width="200" sortable />
                <el-table-column prop="tableComment" label="备注" min-width="250" />
                <el-table-column prop="tableType" label="类型" width="120">
                    <template #default="scope">
                        <el-tag :type="scope.row.tableType === 'VIEW' ? 'warning' : 'success'">
                            {{ scope.row.tableType }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="180" fixed="right" align="center">
                    <template #default="scope">
                        <el-button type="primary" size="small" text @click="handleViewFields(scope.row)">
                            查看字段
                        </el-button>
                    </template>
                </el-table-column>
            </el-table>
        </el-card>
        <!-- 选择弹窗 -->
        <el-dialog v-model="dialogVisible" :title="dialogTitle" width="800px" destroy-on-close>
            <div class="selection-container">
                <el-transfer v-model="selectedValues" v-loading="dialogLoading" :data="transferData" :titles="['未入库', '已选择']" filterable :props="{
                    key: 'name',
                    label: 'name'
                }">
                    <template #default="{ option }">
                        <span :title="option.comment">{{ option.name }}</span>
                    </template>
                </el-transfer>
            </div>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="dialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="handleConfirmSelect" :loading="saving">确认导入</el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { createTable, getConns, getDBTables, getDBViews, getTablesByConnId } from '@/api/metadata'
import type { MdConn, MdTable } from '@/types/metadata'
import { ArrowDown, Refresh, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const selectedConn = ref('')
const connections = ref<MdConn[]>([])
const allTables = ref<MdTable[]>([])
const searchQuery = ref('')

// 弹窗相关
const dialogVisible = ref(false)
const dialogType = ref<'TABLE' | 'VIEW'>('TABLE')
const dialogLoading = ref(false)
const saving = ref(false)
const selectedValues = ref<string[]>([])
const dbObjects = ref<any[]>([])

// 计算属性
const dialogTitle = computed(() => dialogType.value === 'TABLE' ? '选择表入库' : '选择视图入库')
const transferData = computed(() => {
    // 排除已经在入库列表中的
    const existingNames = new Set(allTables.value.map(t => t.tableName))
    return dbObjects.value
        .filter(obj => !existingNames.has(obj.name))
        .map(obj => ({
            name: obj.name,
            comment: obj.comment
        }))
})

// ... existing fetch functions ...

const openSelectDialog = async (type: 'TABLE' | 'VIEW') => {
    if (!selectedConn.value) {
        ElMessage.warning('请先选择数据源')
        return
    }
    dialogType.value = type
    dialogVisible.value = true
    selectedValues.value = []
    await fetchDBObjects()
}

const fetchDBObjects = async () => {
    dialogLoading.value = true
    try {
        const res = dialogType.value === 'TABLE'
            ? await getDBTables(selectedConn.value)
            : await getDBViews(selectedConn.value)
        dbObjects.value = res?.data || res || []
    } catch (error) {
        console.error('获取数据库对象失败:', error)
        ElMessage.error('获取列表失败')
    } finally {
        dialogLoading.value = false
    }
}

const handleConfirmSelect = async () => {
    if (selectedValues.value.length === 0) {
        ElMessage.warning('请选择至少一项')
        return
    }

    saving.value = true
    try {
        // 获取选中的对象详情
        const selectedDetail = dbObjects.value.filter(obj => selectedValues.value.includes(obj.name))

        // 批量创建（目前采用串行调用，后续可考虑后端支持批量）
        const promises = selectedDetail.map(obj => {
            const tableData: Partial<MdTable> = {
                connID: selectedConn.value,
                tableName: obj.name,
                tableTitle: obj.name,
                tableComment: obj.comment || '',
                tableType: dialogType.value,
                state: 1
            }
            return createTable(tableData)
        })

        await Promise.all(promises)
        ElMessage.success(`成功导入 ${selectedValues.value.length} 个对象`)
        dialogVisible.value = false
        fetchTables()
    } catch (error) {
        console.error('导入失败:', error)
        ElMessage.error('部分或全部导入失败')
    } finally {
        saving.value = false
    }
}

// 计算属性
const filteredTables = computed(() => {
    if (!searchQuery.value) return allTables.value
    const query = searchQuery.value.toLowerCase()
    return allTables.value.filter(t =>
        t.tableName.toLowerCase().includes(query) ||
        (t.tableComment && t.tableComment.toLowerCase().includes(query))
    )
})

// 生命周期
onMounted(async () => {
    await fetchConnections()
})

// 获取数据源
const fetchConnections = async () => {
    try {
        const res: any = await getConns()
        connections.value = res?.data || []
        if (connections.value.length > 0) {
            selectedConn.value = connections.value[0].id as string
            fetchTables()
        }
    } catch (error) {
        console.error('获取数据源失败:', error)
    }
}

// 获取表列表
const fetchTables = async () => {
    if (!selectedConn.value) return
    loading.value = true
    try {
        const res: any = await getTablesByConnId(selectedConn.value)
        allTables.value = res?.data || []
    } catch (error) {
        console.error('获取表列表失败:', error)
    } finally {
        loading.value = false
    }
}

const handleConnChange = () => {
    fetchTables()
}

const handleViewFields = (row: MdTable) => {
    router.push({
        path: '/metadata/fields',
        query: { tableId: row.id }
    })
}
</script>
<style scoped>
.table-list {
    padding: 20px;
}

.page-header {
    margin-bottom: 20px;
}

.filter-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.filter-left {
    display: flex;
    align-items: center;
}
</style>
