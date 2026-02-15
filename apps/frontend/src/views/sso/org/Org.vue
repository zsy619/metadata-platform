<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon">
                    <School />
                </el-icon>
                组织管理
            </h1>
            <div class="header-actions">
                <el-button type="primary" @click="handleCreate" :icon="Plus">新增组织</el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-input v-model="searchQuery" placeholder="请输入组织名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
                <el-select v-model="filterStatus" placeholder="筛选状态" style="width: 150px; margin-left: 10px" clearable @change="handleSearch">
                    <el-option label="全部" value="" />
                    <el-option label="有效" :value="1" />
                    <el-option label="禁用" :value="0" />
                </el-select>
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
                <el-button @click="handleReset" :icon="RefreshLeft">重置</el-button>
            </div>
            <div class="table-area">
                <el-table class="tree-table" v-loading="loading" :element-loading-text="loadingText" :data="tableData" border stripe row-key="id" :tree-props="{ children: 'children', hasChildren: 'hasChildren' }" :indent="24" default-expand-all style="width: 100%; height: 100%;">
                    <template #empty>
                        <el-empty :description="searchQuery ? '未搜索到相关组织' : '暂无组织'">
                            <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增组织</el-button>
                        </el-empty>
                    </template>
                    <el-table-column prop="org_name" label="组织名称" min-width="260">
                        <template #default="scope">
                            <TreeNameCell :row="scope.row" name-field="org_name" />
                        </template>
                    </el-table-column>
                    <el-table-column prop="org_short" label="简称" width="120" show-overflow-tooltip />
                    <el-table-column prop="org_code" label="组织编码" width="150" />
                    <el-table-column prop="kind_code" label="类型编码" width="120" />
                    <el-table-column prop="contact" label="联系人" width="120" />
                    <el-table-column prop="phone" label="联系电话" width="150" />
                    <el-table-column prop="status" label="状态" width="80">
                        <template #default="scope">
                            <el-tag v-if="scope.row.status === 1" type="success">有效</el-tag>
                            <el-tag v-else type="danger">禁用</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="sort" label="排序" width="80" />
                    <el-table-column label="操作" width="260" fixed="right">
                        <template #default="scope">
                            <el-button type="primary" size="small" :icon="Plus" @click="handleAddChild(scope.row)" text bg>新增子级</el-button>
                            <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>编辑</el-button>
                            <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg>删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
        </el-card>
        <OrgForm v-model="dialogVisible" :data="formData" :org-tree-data="filteredOrgTreeData" :exclude-ids="excludeIds" @success="handleFormSuccess" />
    </div>
</template>
<script setup lang="ts">
import { deleteUnit, getUnits } from '@/api/user'
import TreeNameCell from '@/components/table/TreeNameCell.vue'
import { Delete, Edit, Plus, RefreshLeft, School, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import OrgForm from './OrgForm.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')
const filterStatus = ref<number | ''>('')

const allData = ref<any[]>([])

const dialogVisible = ref(false)
const formData = ref<any>({})
const excludeIds = ref<string[]>([])

const getAllDescendantIds = (parentId: string): string[] => {
    const ids: string[] = [parentId]
    const findChildren = (pid: string) => {
        allData.value
            .filter(item => (item.parent_id || '0') === pid)
            .forEach(item => {
                ids.push(item.id)
                findChildren(item.id)
            })
    }
    findChildren(parentId)
    return ids
}

const buildTreeData = (items: any[], parentId = '0', excludeIdSet: Set<string> = new Set()): any[] => {
    return items
        .filter(item => (item.parent_id || '0') === parentId && !excludeIdSet.has(item.id))
        .map(item => ({
            value: item.id,
            label: item.org_code,
            org_name: item.org_name,
            disabled: excludeIdSet.has(item.id),
            children: buildTreeData(items, item.id, excludeIdSet)
        }))
}

const tableData = computed(() => {
    const buildTableTree = (items: any[], parentId = '', level = 0): any[] => {
        return items
            .filter(item => (item.parent_id || '') === parentId)
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
        data = data.filter(item => (item.org_name || '').toLowerCase().includes(query))
    }
    if (filterStatus.value !== '') {
        data = data.filter(item => item.status === filterStatus.value)
    }

    if (!searchQuery.value && filterStatus.value === '') {
        return buildTableTree(allData.value)
    }

    const allParentIds = new Set<string>()
    const findAllParents = (pid: string) => {
        if (!pid || pid === '' || allParentIds.has(pid)) return
        allParentIds.add(pid)
        const parent = allData.value.find(item => item.id === pid)
        if (parent) {
            const parentPid = parent.parent_id || ''
            if (parentPid && parentPid !== '') findAllParents(parentPid)
        }
    }
    data.forEach(item => {
        const pid = item.parent_id || ''
        if (pid && pid !== '') findAllParents(pid)
    })
    const fullList = [...data, ...allData.value.filter(item => allParentIds.has(item.id))]
    return buildTableTree(fullList)
})

const orgTreeData = computed(() => {
    return buildTreeData(allData.value)
})

const filteredOrgTreeData = computed(() => {
    if (excludeIds.value.length === 0) {
        return orgTreeData.value
    }
    const excludeSet = new Set(excludeIds.value)
    return buildTreeData(allData.value, '', excludeSet)
})

const loadData = async () => {
    loadingText.value = '加载中...'
    loading.value = true
    try {
        const res: any = await getUnits()
        allData.value = res.data || res
    } catch (error) {
        console.error('加载组织列表失败:', error)
        ElMessage.error('加载列表失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => { }
const handleDebouncedSearch = () => { }
const handleReset = () => { searchQuery.value = ''; filterStatus.value = '' }

const handleCreate = () => {
    formData.value = { status: 1, sort: 0, parent_id: '' }
    excludeIds.value = []
    dialogVisible.value = true
}

const handleAddChild = (row: any) => {
    formData.value = { status: 1, sort: 0, parent_id: row.id }
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
        await ElMessageBox.confirm(`确定要删除组织 "${row.org_name}" 吗？`, '提示', { type: 'warning' })
        await deleteUnit(row.id)
        ElMessage.success('删除成功')
        loadData()
    } catch (error: any) { if (error !== 'cancel') ElMessage.error(error.message || '删除失败') }
}

const handleFormSuccess = () => {
    loadData()
}

onMounted(() => loadData())
</script>
