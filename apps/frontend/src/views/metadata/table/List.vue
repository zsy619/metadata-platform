<template>
    <div class="container-padding">
        <!-- 页面标题区 -->
        <div class="page-header">
            <h1 class="page-title">
                <el-icon class="title-icon">
                    <Grid />
                </el-icon>
                表与视图
            </h1>
        </div>
        <!-- 主内容卡片 -->
        <el-card class="main-card">
            <!-- 搜索区域 -->
            <div class="search-area">
                <el-select v-model="selectedConn" placeholder="选择数据源" style="width: 240px" @change="handleConnChange">
                    <el-option v-for="conn in connections" :key="conn.id" :label="conn.conn_name" :value="conn.id" />
                </el-select>
                <el-input v-model="searchQuery" placeholder="搜索表或视图名称" clearable :prefix-icon="Search" style="width: 300px; margin-left: 10px" />
                <el-button type="primary" :icon="Search" style="margin-left: 10px" @click="handleSearch">搜索</el-button>
                <el-button :icon="RefreshLeft" @click="handleReset">重置</el-button>
                <el-dropdown trigger="click" @command="openSelectDialog" style="margin-left: 10px">
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
            </div>
            <!-- 表格区域 -->
            <div class="table-area">
                <el-table v-loading="loading" :element-loading-text="loadingText" :data="pagedTables" border stripe style="width: 100%; height: 100%">
                    <template #empty>
                        <el-empty :description="searchQuery ? '未搜索到相关表/视图' : '暂无表/视图数据'" />
                    </template>
                    <el-table-column prop="table_name" label="名称" min-width="150" sortable />
                    <el-table-column prop="table_title" label="标题" min-width="150" />
                    <el-table-column prop="table_comment" label="备注" min-width="200" />
                    <el-table-column prop="table_type" label="类型" width="100">
                        <template #default="scope">
                            <el-tag :type="scope.row.table_type === 'VIEW' ? 'warning' : 'success'">
                                {{ scope.row.table_type }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="table_schema" label="模式" width="120" show-overflow-tooltip />
                    <el-table-column prop="create_at" label="导入时间" width="160">
                        <template #default="scope">
                            {{ formatDateTime(scope.row.create_at) }}
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="320" fixed="right" class-name="action-column">
                        <template #default="scope">
                            <el-button type="info" size="small" :icon="View" @click="handleViewDetail(scope.row)">
                                详情
                            </el-button>
                            <el-button type="primary" size="small" :icon="Refresh" @click="handleRefreshTable(scope.row)">
                                刷新
                            </el-button>
                            <el-button type="success" size="small" :icon="Edit" @click="handleEdit(scope.row)">
                                修改
                            </el-button>
                            <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)">
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
        <!-- 选择弹窗 -->
        <el-dialog v-model="dialogVisible" :title="dialogTitle" width="1000px" destroy-on-close class="custom-dialog transfer-dialog" :close-on-click-modal="false">
            <div style="margin-bottom: 15px;" v-if="schemas.length > 0">
                <span class="m-r-xs">模式(Schema): </span>
                <el-select v-model="selectedSchema" placeholder="选择模式" style="width: 200px" @change="handleSchemaChange">
                    <el-option v-for="schema in schemas" :key="schema" :label="schema" :value="schema" />
                </el-select>
            </div>
            <el-transfer v-model="selectedValues" v-loading="dialogLoading" :data="transferData" :titles="['未入库', '已选择']" filterable :props="{
                key: 'name',
                label: 'name'
            }" class="full-width-transfer">
                <template #default="{ option }">
                    <span :title="option.comment">{{ option.name }}</span>
                </template>
            </el-transfer>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="dialogVisible = false" size="large">取消</el-button>
                    <el-button type="primary" @click="handleConfirmSelect" :loading="saving" size="large">确认导入</el-button>
                </div>
            </template>
        </el-dialog>
        <!-- 修改表信息对话框 -->
        <el-dialog v-model="editDialogVisible" title="修改表信息" width="600px" class="custom-dialog">
            <el-form :model="editForm" label-width="120px" label-position="right">
                <el-form-item label="表标题">
                    <el-input v-model="editForm.table_title" placeholder="请输入表标题" clearable />
                </el-form-item>
                <el-form-item label="表备注">
                    <el-input v-model="editForm.table_comment" type="textarea" :rows="4" placeholder="请输入表备注" />
                </el-form-item>
            </el-form>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="editDialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="handleEditSubmit">确定</el-button>
                </div>
            </template>
        </el-dialog>
        <!-- 详情弹窗 -->
        <el-dialog v-model="detailDialogVisible" title="表详情" width="900px" class="custom-dialog detail-dialog">
            <div v-loading="detailLoading">
                <el-descriptions title="基本信息" :column="2" border>
                    <el-descriptions-item label="名称">{{ currentTable.table_name }}</el-descriptions-item>
                    <el-descriptions-item label="标题">{{ currentTable.table_title }}</el-descriptions-item>
                    <el-descriptions-item label="类型">
                        <el-tag :type="currentTable.table_type === 'VIEW' ? 'warning' : 'success'" size="small">
                            {{ currentTable.table_type }}
                        </el-tag>
                    </el-descriptions-item>
                    <el-descriptions-item label="模式">{{ currentTable.table_schema || '-' }}</el-descriptions-item>
                    <el-descriptions-item label="备注" :span="2">{{ currentTable.table_comment || '-' }}</el-descriptions-item>
                </el-descriptions>
                <div class="m-t-lg">
                    <div class="field-header">
                        <div class="field-title">字段列表 ({{ currentTableFields.length }})</div>
                        <el-button type="primary" size="small" plain :icon="Refresh" @click="handleDetailSync" :loading="detailSyncing">
                            同步数据库结构
                        </el-button>
                    </div>
                    <el-table :data="currentTableFields" border stripe size="small" height="400px">
                        <el-table-column prop="sort" label="排序" width="70" align="center" />
                        <el-table-column prop="column_name" label="字段名称" min-width="150" />
                        <el-table-column prop="column_title" label="字段标题" min-width="150" />
                        <el-table-column prop="column_type" label="类型" width="120" />
                        <el-table-column prop="column_length" label="长度" width="80" />
                        <el-table-column label="约束" width="120">
                            <template #default="{ row }">
                                <el-tag v-if="row.is_primary_key" type="danger" size="small">PK</el-tag>
                                <el-tag v-if="!row.is_nullable" type="warning" size="small" class="m-l-xs">Not Null</el-tag>
                            </template>
                        </el-table-column>
                        <el-table-column prop="column_comment" label="备注" min-width="200" show-overflow-tooltip />
                    </el-table>
                </div>
            </div>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="detailDialogVisible = false">关闭</el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { createField, createTable, deleteFieldsByTableId, deleteTable, getConns, getDBTables, getDBViews, getFieldsByTableId, getSchemas, getTablesByConnId, getTableStructureFromDB, updateTable } from '@/api/metadata'
import type { MdTable, MdTableField } from '@/types/metadata'
import { showConfirm, showDeleteConfirm } from '@/utils/confirm'
import { ArrowDown, Delete, Edit, Grid, Refresh, RefreshLeft, Search, View } from '@element-plus/icons-vue'
import { ElLoading, ElMessage } from 'element-plus'
import { computed, onMounted, ref } from 'vue'

// 响应式数据
const loading = ref(false)
const loadingText = ref('加载中...')
const selectedConn = ref('')
const connections = ref<any[]>([])
const allTables = ref<MdTable[]>([])
const searchQuery = ref('')

// 分页状态
const currentPage = ref(1)
const pageSize = ref(20)
const total = computed(() => filteredTables.value.length)

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
    const existingNames = new Set(allTables.value.map(t => t.table_name))
    return dbObjects.value
        .filter(obj => !existingNames.has(obj.name))
        .map(obj => ({
            name: obj.name,
            comment: obj.comment
        }))
})

// ... existing fetch functions ...

const schemas = ref<string[]>([])
const selectedSchema = ref('')

const openSelectDialog = async (type: 'TABLE' | 'VIEW') => {
    if (!selectedConn.value) {
        ElMessage.warning('请先选择数据源')
        return
    }
    dialogType.value = type
    dialogVisible.value = true
    selectedValues.value = []
    schemas.value = []
    selectedSchema.value = ''

    // 先获取Schema列表
    try {
        const res: any = await getSchemas(selectedConn.value)
        const schemaList = res?.data || res || []
        schemas.value = Array.isArray(schemaList) ? schemaList : []
        // 如果有Schema，默认选中第一个
        if (schemas.value.length > 0) {
            selectedSchema.value = schemas.value[0]
        }
    } catch (error) {
        console.warn('获取Schema列表失败:', error)
    }

    await fetchDBObjects()
}

const handleSchemaChange = () => {
    fetchDBObjects()
}

const fetchDBObjects = async () => {
    dialogLoading.value = true
    try {
        const res: any = dialogType.value === 'TABLE'
            ? await getDBTables(selectedConn.value, selectedSchema.value)
            : await getDBViews(selectedConn.value, selectedSchema.value)

        const list = res?.data || res
        dbObjects.value = Array.isArray(list) ? list : []
    } catch (error) {
        console.error('获取数据库对象失败:', error)
        ElMessage.error('获取列表失败')
        dbObjects.value = []
    } finally {
        dialogLoading.value = false
    }
}

const handleConfirmSelect = async () => {
    if (selectedValues.value.length === 0) {
        ElMessage.warning('请选择至少一项')
        return
    }

    try {
        await showConfirm(`确定要导入选中的 ${selectedValues.value.length} 个${dialogType.value === 'TABLE' ? '表' : '视图'}吗？`, '导入确认', 'info')
    } catch {
        return
    }

    saving.value = true
    let successCount = 0
    let failCount = 0

    try {
        // 获取选中的对象详情
        const selectedDetail = dbObjects.value.filter(obj => selectedValues.value.includes(obj.name))

        ElMessage.info(`开始导入 ${selectedDetail.length} 个${dialogType.value === 'TABLE' ? '表' : '视图'}...`)

        // 逐个处理表/视图(串行处理,确保数据一致性)
        for (const obj of selectedDetail) {
            try {
                // 1. 从数据库获取表结构信息
                const structureRes = await getTableStructureFromDB(selectedConn.value, obj.name, selectedSchema.value)
                const tableStructure = structureRes?.data || structureRes

                if (!tableStructure) {
                    console.error(`获取表 ${obj.name} 的结构信息失败`)
                    failCount++
                    continue
                }

                // 2. 创建表记录
                const tableData: Partial<MdTable> = {
                    conn_id: selectedConn.value,
                    table_name: obj.name,
                    table_title: obj.name,
                    table_comment: tableStructure.comment || obj.comment || '',
                    table_type: dialogType.value,
                    table_schema: tableStructure.schema || selectedSchema.value || '',
                    state: 1
                }

                const resTable: any = await createTable(tableData)
                const createdTable = resTable?.data || resTable

                if (!createdTable || !createdTable.id) {
                    console.error(`创建表 ${obj.name} 记录失败, 返回原始数据:`, resTable)
                    failCount++
                    continue
                }

                // 3. 获取并保存字段信息
                let columns = []
                if (Array.isArray(tableStructure)) {
                    columns = tableStructure
                } else {
                    columns = tableStructure.columns || tableStructure.fields || []
                }

                if (columns.length > 0) {
                    // 批量创建字段记录
                    const fieldPromises = columns.map((col: any, index: number) => {
                        const fieldData: any = {
                            conn_id: selectedConn.value,
                            table_id: createdTable.id,
                            table_name: obj.name,
                            table_title: obj.name,
                            column_name: col.name || col.column_name,
                            column_title: col.name || col.column_name,
                            column_type: col.type || col.data_type || '',
                            column_length: col.length || col.character_maximum_length || 0,
                            column_comment: col.comment || col.column_comment || '',
                            is_nullable: col.is_nullable === true || col.nullable === true,
                            is_primary_key: col.is_primary_key === true,
                            is_auto_increment: col.is_auto_increment === true,
                            default_value: col.default_value || col.column_default || '',
                            extra_info: col.extra || '',
                            state: 1,
                            sort: col.sort || index + 1
                        }
                        return createField(fieldData)
                    })

                    await Promise.all(fieldPromises)
                }

                successCount++
                ElMessage.success(`成功导入: ${obj.name}`)
            } catch (error) {
                console.error(`导入 ${obj.name} 失败:`, error)
                failCount++
                ElMessage.error(`导入失败: ${obj.name}`)
            }
        }

        // 显示最终结果
        if (successCount > 0) {
            ElMessage.success(`成功导入 ${successCount} 个对象${failCount > 0 ? `, ${failCount} 个失败` : ''}`)
        } else {
            ElMessage.error('全部导入失败')
        }

        // 关闭对话框并刷新列表
        if (successCount > 0) {
            dialogVisible.value = false
            fetchTables()
        }
    } catch (error) {
        console.error('导入过程出错:', error)
        ElMessage.error('导入过程出错,请查看控制台')
    } finally {
        saving.value = false
    }
}

// 计算属性 - 筛选
const filteredTables = computed(() => {
    if (!searchQuery.value) return allTables.value
    const query = searchQuery.value.toLowerCase()
    return allTables.value.filter(t =>
        t.table_name.toLowerCase().includes(query) ||
        (t.table_comment && t.table_comment.toLowerCase().includes(query))
    )
})

// 计算属性 - 分页数据
const pagedTables = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    return filteredTables.value.slice(start, end)
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
    selectedConn.value = connections.value.length > 0 ? connections.value[0].id as string : ''
    currentPage.value = 1
    fetchTables()
}

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
    loadingText.value = '加载中...'
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

// 时间格式化函数
const formatDateTime = (dateStr: string | undefined) => {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return '-'

    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')

    return `${year}-${month}-${day} ${hours}:${minutes}`
}

// 查看详情
const handleViewDetail = async (row: MdTable) => {
    console.log('正在查看表详情, 参数 row:', JSON.parse(JSON.stringify(row)))
    currentTable.value = row
    detailDialogVisible.value = true
    detailLoading.value = true
    currentTableFields.value = [] // 重置字段列表

    if (!row.id) {
        console.error('查看详情失败: row.id 为空')
        ElMessage.error('该项目数据异常 (缺少ID)')
        detailLoading.value = false
        return
    }

    try {
        const res: any = await getFieldsByTableId(row.id as string)
        console.log('获取字段列表成功:', res)
        // 兼容后端统一响应格式 {code: 200, data: [...], message: "success"}
        currentTableFields.value = Array.isArray(res) ? res : (res?.data || [])
    } catch (error: any) {
        console.error('获取字段列表失败, ID:', row.id, '错误详情:', error)
        // 允许继续显示基本信息, 但提示字段获取失败
        ElMessage.warning('无法获取字段列表, 仅显示基本信息')
    } finally {
        detailLoading.value = false
    }
}

// 在详情中同步字段信息 (解决 500 报错或数据损坏)
const detailSyncing = ref(false)
const handleDetailSync = async () => {
    if (!currentTable.value.id) return

    try {
        await showConfirm(
            '同步将重新从数据库读取字段结构并更新当前元数据, 确认继续?',
            '同步确认',
            'info'
        )

        detailSyncing.value = true
        // 复用刷新逻辑的核心
        await doRefreshTable(currentTable.value)

        // 重新加载显示
        const res: any = await getFieldsByTableId(currentTable.value.id as string)
        currentTableFields.value = Array.isArray(res) ? res : (res?.data || [])
        ElMessage.success('字段同步成功')
    } catch (error: any) {
        if (error !== 'cancel') {
            console.error('详情同步失败:', error)
            ElMessage.error(error?.message || '同步失败')
        }
    } finally {
        detailSyncing.value = false
    }
}

// 提取公共的刷新逻辑
const doRefreshTable = async (table: any) => {
    // 1. 删除现有字段
    await deleteFieldsByTableId(table.id as string)

    // 2. 从数据库获取最新表结构
    const structureRes = await getTableStructureFromDB(table.conn_id, table.table_name)
    const tableStructure = structureRes?.data || structureRes?.result || structureRes

    // 兼容后端直接返回字段数组的情况
    let columns = []
    if (Array.isArray(tableStructure)) {
        columns = tableStructure
    } else if (tableStructure && (tableStructure.columns || tableStructure.fields)) {
        columns = tableStructure.columns || tableStructure.fields
    } else {
        throw new Error('获取表结构失败 (数据库返回格式无法识别)')
    }

    // 3. 批量创建字段记录
    if (columns.length > 0) {
        const fieldPromises = columns.map((col: any, index: number) => {
            const fieldData: any = {
                conn_id: table.conn_id,
                table_id: table.id,
                table_name: table.table_name,
                table_title: table.table_title || table.table_name,
                column_name: col.name || col.column_name,
                column_title: col.name || col.column_name,
                column_type: col.type || col.data_type || '',
                column_length: col.length || col.character_maximum_length || 0,
                column_comment: col.comment || col.column_comment || '',
                is_nullable: col.is_nullable === true || col.nullable === true,
                is_primary_key: col.is_primary_key === true,
                is_auto_increment: col.is_auto_increment === true,
                default_value: col.default_value || col.column_default || '',
                extra_info: col.extra || '',
                state: 1,
                sort: col.sort || index + 1
            }
            return createField(fieldData)
        })

        await Promise.all(fieldPromises)
    }
}

// 刷新表字段 (主页面调用)
const handleRefreshTable = async (row: MdTable) => {
    showConfirm(
        `刷新将删除现有字段并重新从数据库同步,确定要刷新表 "${row.table_title || row.table_name}" 吗?`,
        '刷新确认',
        'warning'
    ).then(async () => {
        const loading = ElLoading.service({ text: '正在刷新...' })
        try {
            await doRefreshTable(row)
            ElMessage.success('刷新成功')
            // 如果刷新的是当前正在查看的详情, 则也更新详情列表
            if (detailDialogVisible.value && currentTable.value.id === row.id) {
                const res: any = await getFieldsByTableId(row.id as string)
                currentTableFields.value = Array.isArray(res) ? res : (res?.data || [])
            }
        } catch (error: any) {
            console.error('刷新失败:', error)
            ElMessage.error(error?.message || '刷新失败')
        } finally {
            loading.close()
        }
    }).catch(() => {
        // 用户取消
    })
}

// 修改表(只修改标题和备注)
const editDialogVisible = ref(false)
const editForm = ref({
    id: '',
    table_title: '',
    table_comment: ''
})

const detailDialogVisible = ref(false)
const currentTable = ref<Partial<MdTable>>({})
const currentTableFields = ref<MdTableField[]>([])
const detailLoading = ref(false)

const handleEdit = (row: MdTable) => {
    editForm.value = {
        id: row.id as string,
        table_title: row.table_title || '',
        table_comment: row.table_comment || ''
    }
    editDialogVisible.value = true
}

const handleEditSubmit = async () => {
    if (!editForm.value.id) return
    try {
        await updateTable(editForm.value.id, {
            table_title: editForm.value.table_title,
            table_comment: editForm.value.table_comment
        })
        ElMessage.success('修改成功')
        editDialogVisible.value = false
        fetchTables()
    } catch (error) {
        console.error('修改失败:', error)
        ElMessage.error('修改失败')
    }
}

// 删除表
const handleDelete = async (row: MdTable) => {
    if (!row.id) {
        ElMessage.error('该表数据异常,缺少ID')
        return
    }

    let fieldCount = 0
    let checkFailed = false

    try {
        const res: any = await getFieldsByTableId(row.id as string)
        const fields = Array.isArray(res) ? res : (res?.data || [])
        fieldCount = fields.length

        if (fieldCount > 0) {
            ElMessage.warning(`该表有 ${fieldCount} 个字段,不可删除`)
            return
        }
    } catch (error: any) {
        console.warn('获取字段列表失败,可能表数据异常:', error)
        checkFailed = true
    }

    // 执行删除确认
    const confirmMsg = checkFailed
        ? `无法获取该表的字段列表（可能数据异常），是否确定要强制删除表 "${row.table_title || row.table_name}"？`
        : `确定要删除表 "${row.table_title || row.table_name}" 吗?`

    showDeleteConfirm(confirmMsg).then(async () => {
        try {
            await deleteTable(row.id as string)
            ElMessage.success('删除成功')
            fetchTables()
        } catch (delError: any) {
            console.error('删除操作失败:', delError)
            ElMessage.error(delError?.message || '删除失败')
        }
    }).catch(() => {
        // 用户取消
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

/* ==================== 详情弹窗样式 ==================== */
.m-t-lg {
    margin-top: 24px;
}

.field-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    padding-right: 4px;
}

.field-title {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    padding-left: 10px;
    border-left: 4px solid #409eff;
}

:deep(.detail-dialog .el-descriptions__title) {
    font-size: 16px;
    font-weight: 600;
}


/* ==================== 穿梭框弹窗特有处理 ==================== */
:deep(.transfer-dialog .el-dialog__body) {
    overflow: hidden;
}

.dialog-footer {
    display: flex;
    justify-content: center;
    gap: 16px;
}

/* ==================== 穿梭框完美布局 ==================== */
/* 穿梭框占满整个对话框宽度 */
.full-width-transfer {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

/* 左右面板等宽,完美对称 */
:deep(.full-width-transfer .el-transfer-panel) {
    width: calc(50% - 40px);
    height: 500px;
    border: none;
    border-radius: 8px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
}

/* 面板头部样式 - 修复选择框与标题对齐 */
:deep(.full-width-transfer .el-transfer-panel__header) {
    background: #f5f7fa;
    padding: 12px 16px;
    border-bottom: 1px solid #dcdfe6;
    display: flex;
    align-items: center;
    min-height: 48px;
}

:deep(.full-width-transfer .el-transfer-panel__header .el-checkbox) {
    margin-right: 12px;
    display: flex;
    align-items: center;
}

/* 确保checkbox和文字在同一行 */
:deep(.full-width-transfer .el-transfer-panel__header .el-checkbox__input) {
    display: flex;
    align-items: center;
}

:deep(.full-width-transfer .el-transfer-panel__header .el-checkbox__label) {
    display: flex;
    align-items: center;
    padding-left: 8px;
}

/* 面板标题 */
:deep(.full-width-transfer .el-transfer-panel__header .el-transfer-panel__header-title) {
    font-weight: 600;
    color: #303133;
    line-height: 1;
    display: flex;
    align-items: center;
}

/* 搜索框 */
:deep(.full-width-transfer .el-transfer-panel__filter) {
    padding: 12px 16px;
    border-bottom: 1px solid #ebeef5;
}

:deep(.full-width-transfer .el-transfer-panel__filter .el-input__inner) {
    border-radius: 6px;
}

/* 列表区域 */
:deep(.full-width-transfer .el-transfer-panel__body) {
    flex: 1;
    overflow: auto;
    padding: 8px 0;
}

:deep(.full-width-transfer .el-transfer-panel__list) {
    height: 100%;
}

/* 列表项样式 - 精确对齐选择框与文字 */
/* Checkbox作为flex容器 - 确保垂直居中 */
:deep(.full-width-transfer .el-transfer-panel__item.el-checkbox) {
    display: flex !important;
    flex-direction: row !important;
    align-items: center !important;
    justify-content: flex-start !important;
    width: calc(100% - 16px) !important;
    min-height: 36px !important;
    padding: 0 16px !important;
    margin: 0 8px 4px !important;
    border-radius: 4px !important;
    cursor: pointer !important;
    transition: background-color 0.2s !important;
    box-sizing: border-box !important;
}

:deep(.full-width-transfer .el-transfer-panel__item.el-checkbox:hover) {
    background-color: #f5f7fa !important;
}

/* Checkbox输入区域 - 精确尺寸和位置 */
:deep(.full-width-transfer .el-transfer-panel__item.el-checkbox .el-checkbox__input) {
    display: inline-flex !important;
    align-items: center !important;
    justify-content: center !important;
    width: 16px !important;
    height: 16px !important;
    margin: 0 !important;
    padding: 0 !important;
    flex-shrink: 0 !important;
    position: relative !important;
    top: 0 !important;
}

/* Checkbox方框样式 */
:deep(.full-width-transfer .el-transfer-panel__item.el-checkbox .el-checkbox__inner) {
    display: block !important;
    width: 16px !important;
    height: 16px !important;
    border: 1px solid #dcdfe6 !important;
    border-radius: 3px !important;
    background-color: #fff !important;
    box-sizing: border-box !important;
    transition: all 0.3s !important;
    margin: 0 !important;
    padding: 0 !important;
    position: relative !important;
    top: 0 !important;
}

/* 选中状态 */
:deep(.full-width-transfer .el-transfer-panel__item.el-checkbox.is-checked .el-checkbox__inner) {
    background-color: #409eff !important;
    border-color: #409eff !important;
}

/* Checkbox文字标签 - 间距2px,垂直居中 */
:deep(.full-width-transfer .el-transfer-panel__item.el-checkbox .el-checkbox__label) {
    display: inline-block !important;
    margin-left: 2px !important;
    padding: 0 !important;
    font-size: 14px !important;
    color: #606266 !important;
    line-height: 16px !important;
    vertical-align: top !important;
    white-space: normal !important;
    word-wrap: break-word !important;
    word-break: break-all !important;
    flex: 1 !important;
    position: relative !important;
    top: 0 !important;
}

/* 隐藏原始input */
:deep(.full-width-transfer .el-transfer-panel__item.el-checkbox .el-checkbox__original) {
    opacity: 0 !important;
    outline: none !important;
    position: absolute !important;
    left: -9999px !important;
    margin: 0 !important;
    width: 0 !important;
    height: 0 !important;
    z-index: -1 !important;
}

/* 禁用状态 */
:deep(.full-width-transfer .el-transfer-panel__item.el-checkbox.is-disabled) {
    cursor: not-allowed !important;
    opacity: 0.5 !important;
}

:deep(.full-width-transfer .el-transfer-panel__item.el-checkbox.is-disabled .el-checkbox__label) {
    color: #c0c4cc !important;
    cursor: not-allowed !important;
}

/* 空状态 */
:deep(.full-width-transfer .el-transfer-panel__empty) {
    color: #909399;
    font-size: 14px;
}

/* ==================== 中间按钮区域 - 垂直居中排列 ==================== */
:deep(.full-width-transfer .el-transfer__buttons) {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    padding: 0 10px;
    gap: 12px;
}

/* 按钮样式 - 圆形,完美对称 */
:deep(.full-width-transfer .el-transfer__buttons .el-button) {
    margin: 0;
    width: 48px;
    height: 48px;
    padding: 0;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    transition: all 0.3s;
}

:deep(.full-width-transfer .el-transfer__buttons .el-button:hover) {
    transform: scale(1.1);
    box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

:deep(.full-width-transfer .el-transfer__buttons .el-button .el-icon) {
    font-size: 18px;
}

/* 禁用状态 */
:deep(.full-width-transfer .el-transfer__buttons .el-button:disabled) {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
}

/* 滚动条美化 - 默认隐藏,悬停时显示 */
:deep(.full-width-transfer .el-transfer-panel__body) {
    scrollbar-width: thin;
    scrollbar-color: transparent transparent;
}

:deep(.full-width-transfer .el-transfer-panel__body:hover) {
    scrollbar-color: #dcdfe6 #f5f7fa;
}

:deep(.full-width-transfer .el-transfer-panel__body::-webkit-scrollbar) {
    width: 6px;
}

:deep(.full-width-transfer .el-transfer-panel__body::-webkit-scrollbar-thumb) {
    background-color: transparent;
    border-radius: 3px;
    transition: background-color 0.3s;
}

:deep(.full-width-transfer .el-transfer-panel__body:hover::-webkit-scrollbar-thumb) {
    background-color: #dcdfe6;
}

:deep(.full-width-transfer .el-transfer-panel__body::-webkit-scrollbar-thumb:hover) {
    background-color: #c0c4cc;
}

:deep(.full-width-transfer .el-transfer-panel__body::-webkit-scrollbar-track) {
    background-color: transparent;
}

:deep(.full-width-transfer .el-transfer-panel__body:hover::-webkit-scrollbar-track) {
    background-color: #f5f7fa;
}

:deep(.action-column .cell) {
    white-space: nowrap;
}
</style>
