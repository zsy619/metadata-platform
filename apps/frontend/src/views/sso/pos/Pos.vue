<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon">
                    <Briefcase />
                </el-icon>
                职位管理
            </h1>
            <div class="header-actions">
                <el-button type="primary" @click="handleCreate" :icon="Plus">新增职位</el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-input v-model="searchQuery" placeholder="请输入职位名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
                <el-button @click="handleReset" :icon="RefreshLeft">重置</el-button>
            </div>
            <div class="table-area">
                <el-table class="tree-table" v-loading="loading" :element-loading-text="loadingText" :data="treeData" border stripe row-key="id" :tree-props="{ children: 'children', hasChildren: 'hasChildren' }" :indent="24" default-expand-all style="width: 100%; height: 100%;">
                    <template #empty>
                        <el-empty :description="searchQuery ? '未搜索到相关职位' : '暂无职位'">
                            <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增职位</el-button>
                        </el-empty>
                    </template>
                    <el-table-column prop="pos_name" label="职位名称" min-width="260">
                        <template #default="scope">
                            <TreeNameCell :row="scope.row" name-field="pos_name" />
                        </template>
                    </el-table-column>
                    <el-table-column prop="pos_code" label="职位编码" width="150" />
                    <el-table-column prop="data_range" label="数据范围" width="120">
                        <template #default="scope">
                            <el-tag v-if="scope.row.data_range === DATA_RANGE.ALL" type="success">{{ DATA_RANGE_LABELS[DATA_RANGE.ALL] }}</el-tag>
                            <el-tag v-else-if="scope.row.data_range === DATA_RANGE.CUSTOM" type="warning">{{ DATA_RANGE_LABELS[DATA_RANGE.CUSTOM] }}</el-tag>
                            <el-tag v-else-if="scope.row.data_range === DATA_RANGE.DEPARTMENT" type="info">{{ DATA_RANGE_LABELS[DATA_RANGE.DEPARTMENT] }}</el-tag>
                            <el-tag v-else type="info">{{ DATA_RANGE_LABELS[DATA_RANGE.DEPT_AND_BELOW] }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="status" label="状态" width="80">
                        <template #default="scope">
                            <el-tag v-if="scope.row.status === 1" type="success">有效</el-tag>
                            <el-tag v-else type="danger">禁用</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="sort" label="排序" width="80" />
                    <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
                    <el-table-column label="操作" width="260" fixed="right">
                        <template #default="scope">
                            <el-button type="primary" size="small" :icon="Plus" @click="handleCreateChild(scope.row)" text bg>新增子级</el-button>
                            <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>编辑</el-button>
                            <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg v-if="!scope.row.is_system && !hasChildren(scope.row.id)">删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
        </el-card>
        <PosForm v-model="dialogVisible" :data="formData" :pos-tree-data="posTreeSelectData" :exclude-ids="excludeIds" :org-list="orgList" @success="loadData" />
    </div>
</template>
<script setup lang="ts">
import { deletePos, getPos, getUnits } from '@/api/user'
import TreeNameCell from '@/components/table/TreeNameCell.vue'
import { DATA_RANGE, DATA_RANGE_LABELS } from '@/utils/constants'
import { Briefcase, Delete, Edit, Plus, RefreshLeft, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import PosForm from './PosForm.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')

const allData = ref<any[]>([])
const orgList = ref<any[]>([])

/**
 * 构建树型数据结构
 * @param items 扁平数据列表
 * @param parentId 父节点ID，默认为空字符串表示根节点
 * @returns 树型结构数据
 */
const buildTree = (items: any[], parentId = ''): any[] => {
    return items
        .filter(item => (item.parent_id || '') === parentId)
        .map(item => ({
            ...item,
            children: buildTree(items, item.id)
        }))
}

/**
 * 树型表格数据 - 用于表格展示
 * 支持搜索时返回扁平化结果
 */
const treeData = computed(() => {
    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        return allData.value.filter(item => (item.pos_name || '').toLowerCase().includes(query))
    }
    return buildTree(allData.value)
})

/**
 * 构建树型选择器数据结构
 * @param items 扁平数据列表
 * @param parentId 父节点ID
 * @param excludeIds 需要排除的节点ID列表（用于编辑时排除自己及子节点）
 * @returns 树型选择器数据
 */
const buildTreeSelectData = (items: any[], parentId = '', excludeIds: string[] = []): any[] => {
    return items
        .filter(item => (item.parent_id || '') === parentId && !excludeIds.includes(item.id))
        .map(item => ({
            value: item.id,
            label: item.pos_name,
            children: buildTreeSelectData(items, item.id, excludeIds)
        }))
}

/**
 * 获取指定节点的所有子节点ID（递归）
 * @param items 所有数据列表
 * @param parentId 父节点ID
 * @returns 所有子节点ID列表
 */
const getAllChildrenIds = (items: any[], parentId: string): string[] => {
    const children = items.filter(item => (item.parent_id || '') === parentId)
    let ids: string[] = []
    for (const child of children) {
        ids.push(child.id)
        ids = ids.concat(getAllChildrenIds(items, child.id))
    }
    return ids
}

/**
 * 树型选择器数据 - 用于上级职位选择
 */
const posTreeSelectData = computed(() => buildTreeSelectData(allData.value))

/**
 * 编辑时需要排除的节点ID列表（包含自己及所有子节点）
 */
const excludeIds = ref<string[]>([])

const dialogVisible = ref(false)
const formData = ref<any>({})

const loadData = async () => {
    loadingText.value = '加载中...'
    loading.value = true
    try {
        const [posRes, orgRes]: [any, any] = await Promise.all([
            getPos(),
            getUnits()
        ])
        allData.value = posRes.data || posRes
        orgList.value = orgRes.data || orgRes
    } catch (error) {
        console.error('加载职位列表失败:', error)
        ElMessage.error('加载列表失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => { }
const handleDebouncedSearch = () => { }
const handleReset = () => { searchQuery.value = '' }

const handleCreate = () => {
    formData.value = { status: 1, sort: 0, data_range: DATA_RANGE.ALL, parent_id: '' }
    excludeIds.value = []
    dialogVisible.value = true
}

const handleCreateChild = (row: any) => {
    formData.value = { status: 1, sort: 0, data_range: DATA_RANGE.ALL, parent_id: row.id }
    excludeIds.value = []
    dialogVisible.value = true
}

const handleEdit = (row: any) => {
    formData.value = { ...row }
    excludeIds.value = [row.id, ...getAllChildrenIds(allData.value, row.id)]
    dialogVisible.value = true
}

const hasChildren = (id: string): boolean => {
    return allData.value.some(item => (item.parent_id || '') === id)
}

const handleDelete = async (row: any) => {
    try {
        await ElMessageBox.confirm(`确定要删除职位 "${row.pos_name}" 吗？`, '提示', { type: 'warning' })
        await deletePos(row.id)
        ElMessage.success('删除成功')
        loadData()
    } catch (error: any) { if (error !== 'cancel') ElMessage.error(error.message || '删除失败') }
}

onMounted(() => loadData())
</script>
<style scoped>
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

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-shrink: 0;
}

.page-title {
    font-size: 20px;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 8px;
}

.title-icon {
    font-size: 24px;
    color: var(--el-color-primary);
}

.header-actions {
    display: flex;
    gap: 10px;
}

.search-area {
    display: flex;
    margin-bottom: 20px;
    flex-shrink: 0;
}

.table-area {
    flex: 1;
    overflow: hidden;
}

.text-primary {
    color: var(--el-text-color-primary);
}

:deep(.el-table__row) {
    transition: background-color 0.15s ease;
}

:deep(.el-table__row--level-0) {
    background-color: var(--el-fill-color-lighter);
}

:deep(.el-table__row--level-0 > td:first-child) {
    font-weight: 600;
}

:deep(.el-table__indent) {
    padding-left: 0 !important;
}

:deep(.el-table__expand-icon) {
    color: var(--el-color-primary);
    font-size: 14px;
}

:deep(.el-table__expand-icon--expanded) {
    transform: rotate(90deg);
}

:deep(.el-table .el-table__cell) {
    padding: 12px 0;
}

:deep(.el-table .el-table__body tr:hover > td) {
    background-color: var(--el-color-primary-light-9) !important;
}
</style>
