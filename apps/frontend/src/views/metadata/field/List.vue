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
                <el-select v-model="selectedConn" placeholder="所有数据源" style="width: 200px" clearable>
                    <el-option label="所有数据源" value="" />
                    <el-option v-for="conn in connections" :key="conn.id" :label="conn.conn_name" :value="conn.id" />
                </el-select>
                <el-select v-model="selectedTable" placeholder="选择表" style="width: 240px; margin-left: 10px" clearable>
                    <el-option v-for="table in filteredTables" :key="table.id" :label="table.table_name" :value="table.id" />
                </el-select>
                <el-input v-model="searchQuery" placeholder="搜索字段名称或表名" clearable :prefix-icon="Search" style="width: 300px; margin-left: 10px" />
                <el-button type="primary" :icon="Search" style="margin-left: 10px" @click="handleSearch">搜索</el-button>
                <el-button :icon="RefreshLeft" @click="handleReset">重置</el-button>
            </div>
            <!-- 表格区域 -->
            <div class="table-area">
                <el-table v-loading="loading" :element-loading-text="loadingText" :data="pagedFields" border stripe style="width: 100%; height: 100%">
                    <template #empty>
                        <el-empty :description="searchQuery ? '未搜索到相关字段' : '暂无字段数据'" />
                    </template>
                    <el-table-column prop="table_name" label="表名称" width="180" sortable show-overflow-tooltip />
                    <el-table-column prop="column_name" label="字段名称" width="180" sortable show-overflow-tooltip />
                    <el-table-column prop="column_type" label="类型" width="130" />
                    <el-table-column prop="column_length" label="长度" width="90" />
                    <el-table-column prop="is_nullable" label="可为空" width="90">
                        <template #default="scope">
                            <el-tag :type="scope.row.is_nullable ? 'info' : 'danger'" size="small">
                                {{ scope.row.is_nullable ? '是' : '否' }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="is_primary_key" label="主键" width="90" align="center">
                        <template #default="scope">
                            <el-tag :type="scope.row.is_primary_key ? 'warning' : 'info'" size="small" effect="plain">
                                {{ scope.row.is_primary_key ? '是' : '否' }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="column_comment" label="备注" min-width="200" show-overflow-tooltip />
                    <el-table-column label="操作" width="150" fixed="right">
                        <template #default="scope">
                            <el-button type="primary" size="small" :icon="Edit" link @click="handleEdit(scope.row)">
                                编辑
                            </el-button>
                            <el-button type="danger" size="small" :icon="Delete" link @click="handleDelete(scope.row)">
                                删除
                            </el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
            <!-- 分页区域 -->
            <div class="pagination-area">
                <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]" background layout="total, sizes, prev, pager, next, jumper" :total="total" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
        </el-card>
        <!-- 编辑备注对话框 -->
        <el-dialog v-model="editDialogVisible" title="编辑字段备注" width="500px">
            <el-form :model="editForm" label-width="80px">
                <el-form-item label="字段名称">
                    <el-input v-model="editForm.column_name" disabled />
                </el-form-item>
                <el-form-item label="备注信息">
                    <el-input v-model="editForm.column_comment" type="textarea" :rows="4" placeholder="请输入备注信息" />
                </el-form-item>
            </el-form>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="editDialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="handleEditSubmit" :loading="submitting">确定</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { deleteField, getConns, getFields, getTables, updateField } from '@/api/metadata'
import type { MdTable, MdTableField } from '@/types/metadata'
import { showDeleteConfirm } from '@/utils/confirm'
import { Delete, Edit, List, RefreshLeft, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// 响应式数据
const loading = ref(false)
const loadingText = ref('加载中...')
const selectedConn = ref('')
const selectedTable = ref('')
const connections = ref<any[]>([])
const allTables = ref<MdTable[]>([])
const allFields = ref<MdTableField[]>([])
const searchQuery = ref('')

// 编辑相关
const editDialogVisible = ref(false)
const submitting = ref(false)
const editForm = ref({
    id: '',
    column_name: '',
    column_comment: ''
})

// 分页状态
const currentPage = ref(1)
const pageSize = ref(20)
const total = computed(() => filteredFields.value.length)

// 计算属性 - 联动表列表
const filteredTables = computed(() => {
    if (!selectedConn.value) return allTables.value
    return allTables.value.filter(t => t.conn_id == selectedConn.value)
})

// 计算属性 - 本地关键字筛选
const filteredFields = computed(() => {
    if (!searchQuery.value) return allFields.value
    const query = searchQuery.value.toLowerCase()
    return allFields.value.filter(f =>
        f.column_name.toLowerCase().includes(query) ||
        (f.table_name && f.table_name.toLowerCase().includes(query)) ||
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

// 核心搜索/加载函数
const doSearch = async () => {
    loadingText.value = '加载中...'
    loading.value = true
    try {
        console.log(`执行搜索: 数据源=${selectedConn.value}, 表=${selectedTable.value}`)
        const res: any = await getFields(selectedConn.value, selectedTable.value)
        allFields.value = Array.isArray(res) ? res : (res?.data || [])
        console.log(`搜索结果: 获取到 ${allFields.value.length} 个字段`)
        currentPage.value = 1
    } catch (error) {
        console.error('搜索字段列表失败:', error)
    } finally {
        loading.value = false
    }
}

// 搜索按钮处理
const handleSearch = () => {
    doSearch()
}

const handleReset = () => {
    console.log('执行搜索重置')
    selectedConn.value = ''
    selectedTable.value = ''
    searchQuery.value = ''
    currentPage.value = 1
    doSearch()
}

// 编辑处理
const handleEdit = (row: MdTableField) => {
    editForm.value = {
        id: row.id,
        column_name: row.column_name,
        column_comment: row.column_comment || ''
    }
    editDialogVisible.value = true
}

const handleEditSubmit = async () => {
    if (!editForm.value.id) return
    submitting.value = true
    try {
        await updateField(editForm.value.id, {
            column_comment: editForm.value.column_comment
        })
        ElMessage.success('更新成功')
        editDialogVisible.value = false
        doSearch() // 刷新列表
    } catch (error) {
        console.error('更新字段失败:', error)
        ElMessage.error('操作失败')
    } finally {
        submitting.value = false
    }
}

// 删除处理
const handleDelete = (row: MdTableField) => {
    showDeleteConfirm(`确定要物理删除字段 "${row.column_name}" 的元数据吗？此操作不可恢复。`).then(async () => {
        loadingText.value = '正在删除...'
        loading.value = true
        try {
            await deleteField(row.id)
            ElMessage.success('删除成功')
            // doSearch will reset loading text
            doSearch()
        } catch (error) {
            console.error('删除字段失败:', error)
            ElMessage.error('操作失败')
            loading.value = false
        }
    })
}

// 侦听器 - 联动逻辑
watch(selectedConn, (newVal, oldVal) => {
    console.log(`数据源已切换: ${oldVal} -> ${newVal}`)
    // 只有在手动切换（非初始化恢复）且新旧值不一致时才重置表选择
    if (newVal !== undefined && oldVal !== undefined && newVal !== oldVal) {
        selectedTable.value = ''
    }
    doSearch()
})

watch(selectedTable, (newVal, oldVal) => {
    console.log(`目标表已切换: ${oldVal} -> ${newVal}`)
    if (newVal !== oldVal) {
        doSearch()
    }
})

// 生命周期
onMounted(async () => {
    loading.value = true
    try {
        await Promise.all([
            fetchConnections(),
            fetchTables()
        ])

        if (route.query.tableId) {
            const tId = route.query.tableId as string
            const table = allTables.value.find(t => t.id === tId)
            if (table) {
                // 先锁住 conn，避免重置逻辑触发
                selectedConn.value = table.conn_id
                // 延迟一个小 tick 确保 watch(selectedConn) 处理完后再设 table
                setTimeout(() => {
                    selectedTable.value = tId
                }, 100)
            } else {
                selectedTable.value = tId
                await doSearch()
            }
        } else {
            await doSearch()
        }
    } finally {
        loading.value = false
    }
})

// 获取数据源
const fetchConnections = async () => {
    try {
        const res: any = await getConns()
        connections.value = res?.data || res || []
    } catch (error) {
        console.error('获取数据源失败:', error)
    }
}

// 获取表列表
const fetchTables = async () => {
    try {
        const res: any = await getTables()
        allTables.value = res?.data || res || []
    } catch (error) {
        console.error('获取表列表失败:', error)
    }
}
</script>
<style scoped>
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
