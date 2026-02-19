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
                <el-table class="tree-table" v-loading="loading" :element-loading-text="loadingText" :data="tableData" border stripe row-key="id" :tree-props="{ children: 'children', hasChildren: 'hasChildren' }" :indent="24" default-expand-all style="width: 100%; height: 100%;">
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
                    <el-table-column label="操作" width="320" fixed="right">
                        <template #default="scope">
                            <el-button type="primary" size="small" :icon="Plus" @click="handleCreateChild(scope.row)" text bg>新增子级</el-button>
                            <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>编辑</el-button>
                            <el-button type="warning" size="small" @click="handleManageRoles(scope.row)" text bg>关联角色</el-button>
                            <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg v-if="!scope.row.is_system && !scope.row.hasChildren">删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
        </el-card>
        <PosForm v-model="dialogVisible" :data="formData" :pos-tree-data="posTreeSelectData" :exclude-ids="excludeIds" :org-list="orgList" @success="loadData" />
        <PosRoleDialog v-model="roleDialogVisible" :pos-id="currentPosId" :pos-name="currentPosName" @success="loadData" />
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
import PosRoleDialog from './PosRoleDialog.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')

const allData = ref<any[]>([])
const orgList = ref<any[]>([])

const dialogVisible = ref(false)
const formData = ref<any>({})
const excludeIds = ref<string[]>([])

// 职位角色弹窗相关
const roleDialogVisible = ref(false)
const currentPosId = ref('')
const currentPosName = ref('')

const getAllDescendantIds = (parentId: string): string[] => {
    const ids: string[] = [parentId]
    const findChildren = (pid: string) => {
        allData.value
            .filter(item => {
                const itemPid = item.parent_id || ''
                return itemPid === pid
            })
            .forEach(item => {
                ids.push(item.id)
                findChildren(item.id)
            })
    }
    findChildren(parentId)
    return ids
}

const buildTreeData = (items: any[], parentId = '', excludeIdSet: Set<string> = new Set()): any[] => {
    // 判断是否为根节点（兼容 '' 和 '0' 两种情况）
    const isRoot = (pid: string) => !pid || pid === '' || pid === '0'
    
    return items
        .filter(item => {
            if (excludeIdSet.has(item.id)) return false
            const itemPid = item.parent_id || ''
            if (isRoot(parentId)) {
                return isRoot(itemPid)
            }
            return itemPid === parentId
        })
        .map(item => ({
            value: item.id,
            label: item.pos_name,
            disabled: excludeIdSet.has(item.id),
            children: buildTreeData(items, item.id, excludeIdSet)
        }))
}

const tableData = computed(() => {
    // 判断是否为根节点（兼容 '' 和 '0' 两种情况）
    const isRoot = (pid: string) => !pid || pid === '' || pid === '0'
    
    const buildTableTree = (items: any[], parentId: string, level = 0): any[] => {
        return items
            .filter(item => {
                const itemPid = item.parent_id || ''
                // 根节点匹配
                if (isRoot(parentId)) {
                    return isRoot(itemPid)
                }
                // 子节点精确匹配
                return itemPid === parentId
            })
            .sort((a, b) => (a.sort || 0) - (b.sort || 0))
            .map(item => {
                const children = buildTableTree(items, item.id, level + 1)
                const result: any = { ...item, level }
                if (children.length > 0) {
                    result.children = children
                    result.hasChildren = true
                } else {
                    result.hasChildren = false
                }
                return result
            })
    }

    let data = allData.value
    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        data = data.filter(item => (item.pos_name || '').toLowerCase().includes(query))
    }

    if (!searchQuery.value) {
        return buildTableTree(allData.value, '')
    }

    const allParentIds = new Set<string>()
    const findAllParents = (pid: string) => {
        if (!pid || pid === '' || pid === '0' || allParentIds.has(pid)) return
        allParentIds.add(pid)
        const parent = allData.value.find(item => item.id === pid)
        if (parent) {
            const parentPid = parent.parent_id || ''
            if (parentPid && parentPid !== '' && parentPid !== '0') findAllParents(parentPid)
        }
    }
    data.forEach(item => {
        const pid = item.parent_id || ''
        if (pid && pid !== '' && pid !== '0') findAllParents(pid)
    })
    const dataIds = new Set(data.map(item => item.id))
    const fullList = [...data, ...allData.value.filter(item => allParentIds.has(item.id) && !dataIds.has(item.id))]
    return buildTableTree(fullList, '')
})

const posTreeSelectData = computed(() => buildTreeData(allData.value))

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
    excludeIds.value = getAllDescendantIds(row.id)
    dialogVisible.value = true
}

const handleDelete = async (row: any) => {
    try {
        await ElMessageBox.confirm(`确定要删除职位 "${row.pos_name}" 吗？`, '提示', { type: 'warning' })
        await deletePos(row.id)
        ElMessage.success('删除成功')
        loadData()
    } catch (error: any) { if (error !== 'cancel') ElMessage.error(error.message || '删除失败') }
}

/**
 * 管理职位角色
 * @param row 职位数据
 */
const handleManageRoles = (row: any) => {
    currentPosId.value = row.id
    currentPosName.value = row.pos_name
    roleDialogVisible.value = true
}

onMounted(() => loadData())
</script>

