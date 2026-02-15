<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon">
                    <Collection />
                </el-icon>
                组织类型管理
            </h1>
            <div class="header-actions">
                <el-button type="primary" @click="handleCreate" :icon="Plus">新增类型</el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-input v-model="searchQuery" placeholder="请输入类型名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
                <el-button @click="handleReset" :icon="RefreshLeft">重置</el-button>
            </div>
            <div class="table-area">
                <el-table class="tree-table" v-loading="loading" :element-loading-text="loadingText" :data="tableData" border stripe row-key="id" :tree-props="{ children: 'children', hasChildren: 'hasChildren' }" :indent="24" default-expand-all style="width: 100%; height: 100%;">
                    <template #empty>
                        <el-empty :description="searchQuery ? '未搜索到相关类型' : '暂无类型'">
                            <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增类型</el-button>
                        </el-empty>
                    </template>
                    <el-table-column prop="kind_name" label="类型名称" min-width="220">
                        <template #default="scope">
                            <TreeNameCell :row="scope.row" name-field="kind_name" />
                        </template>
                    </el-table-column>
                    <el-table-column prop="kind_code" label="类型编码" width="150" />
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
        <OrgKindForm v-model="dialogVisible" :data="formData" :all-data="allData" @success="handleFormSuccess" />
    </div>
</template>
<script setup lang="ts">
import { deleteOrgKind, getOrgKinds } from '@/api/user'
import TreeNameCell from '@/components/table/TreeNameCell.vue'
import { Collection, Delete, Edit, Plus, RefreshLeft, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import OrgKindForm from './OrgKindForm.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')

const allData = ref<any[]>([])

const dialogVisible = ref(false)
const formData = ref<any>({})

const getAllDescendantIds = (parentId: string): string[] => {
    const ids: string[] = [parentId]
    const findChildren = (pid: string) => {
        allData.value
            .filter(item => (item.parent_id || '') === pid)
            .forEach(item => {
                ids.push(item.id)
                findChildren(item.id)
            })
    }
    findChildren(parentId)
    return ids
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
        data = data.filter(item => (item.kind_name || '').toLowerCase().includes(query))
    }

    if (!searchQuery.value) {
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

const loadData = async () => {
    loadingText.value = '加载中...'
    loading.value = true
    try {
        const res: any = await getOrgKinds()
        allData.value = res.data || res
    } catch (error) {
        console.error('加载组织类型列表失败:', error)
        ElMessage.error('加载列表失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => { }
const handleDebouncedSearch = () => { }
const handleReset = () => { searchQuery.value = '' }

const handleCreate = () => {
    formData.value = { status: 1, sort: 0, parent_id: '' }
    dialogVisible.value = true
}

const handleAddChild = (row: any) => {
    formData.value = { status: 1, sort: 0, parent_id: row.id }
    dialogVisible.value = true
}

const handleEdit = (row: any) => {
    const data = { ...row }
    if (data.kind_name && !data.name) {
        data.name = data.kind_name
    }
    if (data.kind_code && !data.code) {
        data.code = data.kind_code
    }
    if (data.parent_id === '0' || data.parent_id === 0) {
        data.parent_id = ''
    }
    formData.value = data
    dialogVisible.value = true
}

const handleDelete = async (row: any) => {
    try {
        await ElMessageBox.confirm(`确定要删除类型 "${row.kind_name}" 吗？`, '提示', { type: 'warning' })
        await deleteOrgKind(row.id)
        ElMessage.success('删除成功')
        loadData()
    } catch (error: any) { if (error !== 'cancel') ElMessage.error(error.message || '删除失败') }
}

const handleFormSuccess = () => {
    loadData()
}

onMounted(() => loadData())
</script>
