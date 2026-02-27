<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="page-title">
                <el-icon class="title-icon">
                    <Document />
                </el-icon>
                {{ pageTitle }}
            </h1>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-select v-model="selectedConn" placeholder="选择数据源" style="width: 240px" @change="handleConnChange">
                    <el-option v-for="conn in connections" :key="conn.id" :label="conn.conn_name" :value="conn.id" />
                </el-select>
                <div style="margin-left: 10px;" v-if="schemas.length > 0">
                    <span class="m-r-xs">模式(Schema): </span>
                    <el-select v-model="selectedSchema" placeholder="选择模式" style="width: 200px" @change="handleSchemaChange">
                        <el-option v-for="schema in schemas" :key="schema" :label="schema" :value="schema" />
                    </el-select>
                </div>
                <el-input v-model="searchQuery" :placeholder="searchPlaceholder" clearable :prefix-icon="Search" style="width: 300px; margin-left: 10px" />
                <el-button type="primary" :icon="Search" style="margin-left: 10px" @click="handleSearch">搜索</el-button>
                <el-button :icon="RefreshLeft" @click="handleReset">重置</el-button>
                <el-dropdown trigger="click" @command="openSelectDialog" style="margin-left: 10px">
                    <el-button type="primary" plain>
                        元数据导入<el-icon class="el-icon--right"><arrow-down /></el-icon>
                    </el-button>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item command="PROCEDURE" v-if="isProcedurePage">选择存储过程</el-dropdown-item>
                            <el-dropdown-item command="FUNCTION" v-if="isFunctionPage">选择函数</el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </div>
            <div class="table-area">
                <el-table v-loading="loading" :element-loading-text="loadingText" :data="pagedProcedures" border stripe style="width: 100%; height: 100%">
                    <template #empty>
                        <el-empty :description="searchQuery ? '未搜索到相关内容' : '暂无数据'" />
                    </template>
                    <el-table-column prop="proc_name" label="名称" min-width="180" sortable />
                    <el-table-column prop="proc_type" label="类型" width="120" v-if="showTypeColumn">
                        <template #default="scope">
                            <el-tag :type="scope.row.proc_type === 'FUNCTION' ? 'warning' : 'success'">
                                {{ scope.row.proc_type === 'FUNCTION' ? '函数' : '存储过程' }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="return_type" label="返回类型" width="150" show-overflow-tooltip />
                    <el-table-column prop="proc_comment" label="注释" min-width="200" show-overflow-tooltip />
                    <el-table-column prop="create_at" label="导入时间" width="160">
                        <template #default="scope">
                            {{ formatDateTime(scope.row.create_at) }}
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="200" fixed="right" class-name="action-column">
                        <template #default="scope">
                            <el-button type="info" size="small" :icon="View" @click="handleViewDetail(scope.row)">
                                详情
                            </el-button>
                            <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)">
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
        <el-dialog v-model="dialogVisible" :title="dialogTitle" width="1000px" destroy-on-close class="custom-dialog transfer-dialog" :close-on-click-modal="false">
            <div style="margin-bottom: 15px;" v-if="schemas.length > 0">
                <span class="m-r-xs">模式(Schema): </span>
                <el-select v-model="selectedSchema" placeholder="选择模式" style="width: 200px" @change="handleSchemaChangeForDialog">
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
        <el-dialog v-model="detailDialogVisible" title="存储过程/函数详情" width="900px" class="custom-dialog detail-dialog">
            <div v-loading="detailLoading">
                <el-descriptions title="基本信息" :column="2" border>
                    <el-descriptions-item label="名称">{{ currentProcedure.proc_name }}</el-descriptions-item>
                    <el-descriptions-item label="类型" v-if="showTypeColumn">
                        <el-tag :type="currentProcedure.proc_type === 'FUNCTION' ? 'warning' : 'success'" size="small">
                            {{ currentProcedure.proc_type === 'FUNCTION' ? '函数' : '存储过程' }}
                        </el-tag>
                    </el-descriptions-item>
                    <el-descriptions-item label="返回类型" v-if="currentProcedure.return_type">{{ currentProcedure.return_type }}</el-descriptions-item>
                    <el-descriptions-item label="模式">{{ currentProcedure.proc_schema || '-' }}</el-descriptions-item>
                    <el-descriptions-item label="注释" :span="2">{{ currentProcedure.proc_comment || '-' }}</el-descriptions-item>
                </el-descriptions>
                <div class="m-t-lg" v-if="currentProcedure.definition">
                    <div class="field-header">
                        <div class="field-title">定义代码</div>
                    </div>
                    <el-input v-model="currentProcedure.definition" type="textarea" :rows="20" readonly style="width: 100%" />
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
import { getConns, getSchemas, getDBProcedures, getDBFunctions, getProceduresByConnId, createProcedure, deleteProcedure, getParamsByProcId } from '@/api/metadata'
import { ArrowDown, Delete, Document, RefreshLeft, Search, View } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { showConfirm, showDeleteConfirm } from '@/utils/confirm'
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const loading = ref(false)
const loadingText = ref('加载中...')
const selectedConn = ref('')
const connections = ref<any[]>([])
const allProcedures = ref<any[]>([])
const searchQuery = ref('')
const schemas = ref<string[]>([])
const selectedSchema = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = computed(() => filteredProcedures.value.length)
const detailDialogVisible = ref(false)
const currentProcedure = ref<any>({})
const detailLoading = ref(false)

const dialogVisible = ref(false)
const dialogType = ref<'PROCEDURE' | 'FUNCTION'>('PROCEDURE')
const dialogLoading = ref(false)
const saving = ref(false)
const selectedValues = ref<string[]>([])
const dbObjects = ref<any[]>([])

const isProcedurePage = computed(() => route.path.includes('/procedure'))
const isFunctionPage = computed(() => route.path.includes('/function'))
const pageTitle = computed(() => isProcedurePage.value ? '存储过程' : '函数')
const searchPlaceholder = computed(() => isProcedurePage.value ? '搜索存储过程名称' : '搜索函数名称')
const showTypeColumn = computed(() => !isProcedurePage.value && !isFunctionPage.value)

const dialogTitle = computed(() => dialogType.value === 'PROCEDURE' ? '选择存储过程入库' : '选择函数入库')
const transferData = computed(() => {
    const existingNames = new Set(allProcedures.value.map(p => p.proc_name || p.name))
    return dbObjects.value
        .filter(obj => !existingNames.has(obj.name))
        .map(obj => ({
            name: obj.name,
            comment: obj.comment
        }))
})

const openSelectDialog = async (type: 'PROCEDURE' | 'FUNCTION') => {
    if (!selectedConn.value) {
        ElMessage.warning('请先选择数据源')
        return
    }
    dialogType.value = type
    dialogVisible.value = true
    selectedValues.value = []
    
    try {
        const res: any = await getSchemas(selectedConn.value)
        const schemaList = res?.data || res || []
        schemas.value = Array.isArray(schemaList) ? schemaList : []
        if (schemas.value.length > 0) {
            selectedSchema.value = schemas.value[0]
        }
    } catch (error) {
        console.warn('获取Schema列表失败:', error)
    }

    await fetchDBObjects()
}

const fetchDBObjects = async () => {
    dialogLoading.value = true
    try {
        const res: any = dialogType.value === 'PROCEDURE'
            ? await getDBProcedures(selectedConn.value, selectedSchema.value)
            : await getDBFunctions(selectedConn.value, selectedSchema.value)

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

const handleSchemaChangeForDialog = () => {
    fetchDBObjects()
}

const handleConfirmSelect = async () => {
    if (selectedValues.value.length === 0) {
        ElMessage.warning('请选择至少一项')
        return
    }

    try {
        await showConfirm(`确定要导入选中的 ${selectedValues.value.length} 个${dialogType.value === 'PROCEDURE' ? '存储过程' : '函数'}吗？`, '导入确认', 'info')
    } catch {
        return
    }

    saving.value = true
    let successCount = 0
    let failCount = 0

    try {
        const selectedDetail = dbObjects.value.filter(obj => selectedValues.value.includes(obj.name))
        ElMessage.info(`开始导入 ${selectedDetail.length} 个${dialogType.value === 'PROCEDURE' ? '存储过程' : '函数'}...`)

        const selectedConnData = connections.value.find(c => c.id === selectedConn.value)
        const connName = selectedConnData?.conn_name || ''

        for (const obj of selectedDetail) {
            try {
                const procData: any = {
                    conn_id: selectedConn.value,
                    conn_name: connName,
                    proc_schema: selectedSchema.value || '',
                    proc_name: obj.name,
                    proc_title: obj.name,
                    proc_type: dialogType.value,
                    proc_comment: obj.comment || '',
                    definition: obj.definition || '',
                    return_type: obj.return_type || '',
                    language: obj.language || ''
                }

                await createProcedure(procData)
                successCount++
                ElMessage.success(`成功导入: ${obj.name}`)
            } catch (error) {
                console.error(`导入 ${obj.name} 失败:`, error)
                failCount++
                ElMessage.error(`导入失败: ${obj.name}`)
            }
        }

        if (successCount > 0) {
            ElMessage.success(`成功导入 ${successCount} 个对象${failCount > 0 ? `, ${failCount} 个失败` : ''}`)
        } else {
            ElMessage.error('全部导入失败')
        }

        if (successCount > 0) {
            dialogVisible.value = false
            fetchProcedures()
        }
    } catch (error) {
        console.error('导入过程出错:', error)
        ElMessage.error('导入过程出错,请查看控制台')
    } finally {
        saving.value = false
    }
}

const filteredProcedures = computed(() => {
    let result = allProcedures.value
    
    if (isProcedurePage.value) {
        result = result.filter(p => p.proc_type === 'PROCEDURE')
    } else if (isFunctionPage.value) {
        result = result.filter(p => p.proc_type === 'FUNCTION')
    }
    
    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        result = result.filter(p =>
            (p.proc_name && p.proc_name.toLowerCase().includes(query)) ||
            (p.proc_comment && p.proc_comment.toLowerCase().includes(query))
        )
    }
    
    return result
})

const pagedProcedures = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    return filteredProcedures.value.slice(start, end)
})

const handleSizeChange = (val: number) => {
    pageSize.value = val
    currentPage.value = 1
}

const handleCurrentChange = (val: number) => {
    currentPage.value = val
}

const handleSearch = () => {
    currentPage.value = 1
}

const handleReset = () => {
    searchQuery.value = ''
    selectedSchema.value = schemas.value.length > 0 ? schemas.value[0] : ''
    currentPage.value = 1
    fetchProcedures()
}

onMounted(async () => {
    await fetchConnections()
})

const fetchConnections = async () => {
    try {
        const res: any = await getConns()
        connections.value = res?.data || []
        if (connections.value.length > 0) {
            selectedConn.value = connections.value[0].id as string
            await fetchSchemas()
            fetchProcedures()
        }
    } catch (error) {
        console.error('获取数据源失败:', error)
    }
}

const fetchSchemas = async () => {
    if (!selectedConn.value) return
    try {
        const res: any = await getSchemas(selectedConn.value)
        const schemaList = res?.data || res || []
        schemas.value = Array.isArray(schemaList) ? schemaList : []
        if (schemas.value.length > 0) {
            selectedSchema.value = schemas.value[0]
        }
    } catch (error) {
        console.warn('获取Schema列表失败:', error)
        schemas.value = []
    }
}

const fetchProcedures = async () => {
    if (!selectedConn.value) return
    loadingText.value = '加载中...'
    loading.value = true
    try {
        const res: any = await getProceduresByConnId(selectedConn.value)
        allProcedures.value = res?.data || []
    } catch (error) {
        console.error('获取存储过程和函数列表失败:', error)
        ElMessage.error('获取列表失败')
        allProcedures.value = []
    } finally {
        loading.value = false
    }
}

const handleConnChange = async () => {
    await fetchSchemas()
    fetchProcedures()
}

const handleSchemaChange = () => {
    fetchProcedures()
}

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

const handleViewDetail = (row: any) => {
    currentProcedure.value = row
    detailDialogVisible.value = true
}

const handleDelete = async (row: any) => {
    if (!row.id) {
        ElMessage.error('该存储过程/函数数据异常,缺少ID')
        return
    }

    let paramCount = 0

    try {
        const res: any = await getParamsByProcId(row.id as string)
        const params = Array.isArray(res) ? res : (res?.data || [])
        paramCount = params.length
    } catch (error: any) {
        console.warn('获取参数列表失败:', error)
    }

    const confirmMsg = paramCount > 0
        ? `确定要删除 "${row.proc_title || row.proc_name}" 吗？\n\n该${row.proc_type === 'FUNCTION' ? '函数' : '存储过程'}包含 ${paramCount} 个参数，将一起被删除。`
        : `确定要删除 "${row.proc_title || row.proc_name}" 吗?`

    showDeleteConfirm(confirmMsg).then(async () => {
        try {
            await deleteProcedure(row.id as string)
            ElMessage.success('删除成功')
            fetchProcedures()
        } catch (delError: any) {
            console.error('删除操作失败:', delError)
            ElMessage.error(delError?.message || '删除失败')
        }
    }).catch(() => {
    })
}
</script>

<style scoped>
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

.dialog-footer {
    display: flex;
    justify-content: center;
    gap: 16px;
}

.full-width-transfer {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

:deep(.full-width-transfer .el-transfer-panel) {
    width: calc(50% - 40px);
    height: 500px;
    border: none;
    border-radius: 8px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
}

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

:deep(.full-width-transfer .el-transfer-panel__header .el-checkbox__input) {
    display: flex;
    align-items: center;
}

:deep(.full-width-transfer .el-transfer-panel__header .el-checkbox__label) {
    display: flex;
    align-items: center;
    padding-left: 8px;
}

:deep(.full-width-transfer .el-transfer-panel__header .el-transfer-panel__header-title) {
    font-weight: 600;
    color: #303133;
    line-height: 1;
    display: flex;
    align-items: center;
}

:deep(.full-width-transfer .el-transfer-panel__filter) {
    padding: 12px 16px;
    border-bottom: 1px solid #ebeef5;
}

:deep(.full-width-transfer .el-transfer-panel__filter .el-input__inner) {
    border-radius: 6px;
}

:deep(.full-width-transfer .el-transfer-panel__body) {
    flex: 1;
    overflow: auto;
    padding: 8px 0;
}

:deep(.full-width-transfer .el-transfer-panel__list) {
    height: 100%;
}

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

:deep(.full-width-transfer .el-transfer-panel__item.el-checkbox.is-checked .el-checkbox__inner) {
    background-color: #409eff !important;
    border-color: #409eff !important;
}

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

:deep(.full-width-transfer .el-transfer-panel__item.el-checkbox.is-disabled) {
    cursor: not-allowed !important;
    opacity: 0.5 !important;
}

:deep(.full-width-transfer .el-transfer-panel__item.el-checkbox.is-disabled .el-checkbox__label) {
    color: #c0c4cc !important;
    cursor: not-allowed !important;
}

:deep(.full-width-transfer .el-transfer-panel__empty) {
    color: #909399;
    font-size: 14px;
}

:deep(.full-width-transfer .el-transfer__buttons) {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    padding: 0 10px;
    gap: 12px;
}

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

:deep(.full-width-transfer .el-transfer__buttons .el-button:disabled) {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
}

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

.m-r-xs {
    margin-right: 10px;
}
</style>
