<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon">
                    <Monitor />
                </el-icon>
                应用管理
            </h1>
            <div class="header-actions">
                <el-button type="primary" @click="handleCreate" :icon="Plus">新增应用</el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-input v-model="searchQuery" placeholder="请输入应用名称搜索" clearable :prefix-icon="Search" style="width: 300px" />
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
                <el-button @click="handleReset" :icon="RefreshLeft">重置</el-button>
            </div>
            <div class="table-area">
                <el-table class="tree-table" v-loading="loading" :element-loading-text="loadingText" :data="treeData" border stripe row-key="id" :tree-props="{ children: 'children', hasChildren: 'hasChildren' }" :indent="24" default-expand-all style="width: 100%; height: 100%;">
                    <template #empty>
                        <el-empty :description="searchQuery ? '未搜索到相关应用' : '暂无应用'">
                            <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增应用</el-button>
                        </el-empty>
                    </template>
                    <el-table-column prop="app_name" label="应用名称" min-width="220">
                        <template #default="scope">
                            <TreeNameCell :row="scope.row" name-field="app_name" />
                        </template>
                    </el-table-column>
                    <el-table-column prop="app_code" label="应用编码" width="150" />
                    <el-table-column prop="host" label="域名/IP" show-overflow-tooltip />
                    <el-table-column prop="status" label="状态" width="80">
                        <template #default="scope">
                            <el-tag v-if="scope.row.status === 1" type="success">启用</el-tag>
                            <el-tag v-else type="danger">禁用</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="sort" label="排序" width="80" />
                    <el-table-column prop="create_at" label="创建时间" width="170">
                        <template #default="scope">{{ formatDateTime(scope.row.create_at) }}</template>
                    </el-table-column>
                    <el-table-column label="操作" width="260" fixed="right">
                        <template #default="scope">
                            <el-button type="primary" size="small" :icon="Plus" @click="handleAddChild(scope.row)" text bg>新增子级</el-button>
                            <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>编辑</el-button>
                            <el-button v-if="!hasChildren(scope.row.id)" type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg>删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
        </el-card>
        <AppForm v-model="dialogVisible" :data="formData" :all-data="allData" @success="handleFormSuccess" />
    </div>
</template>
<script setup lang="ts">
import { deleteApp, getApps } from '@/api/user'
import TreeNameCell from '@/components/table/TreeNameCell.vue'
import { Delete, Edit, Monitor, Plus, RefreshLeft, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, ref } from 'vue'
import AppForm from './AppForm.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')
const allData = ref<any[]>([])
const treeData = ref<any[]>([])

const dialogVisible = ref(false)
const formData = ref<any>({})

const formatDateTime = (dateStr: string) => {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    return isNaN(date.getTime()) ? '-' : date.toLocaleString('zh-CN')
}

const buildTree = (list: any[], parentId: string = ''): any[] => {
    const items = list
        .filter(item => {
            const pid = item.parent_id || ''
            return pid === parentId
        })
        .sort((a, b) => (a.sort || 0) - (b.sort || 0))

    return items.map(item => {
        const children = buildTree(list, item.id)
        const result: any = { ...item }
        if (children.length > 0) {
            result.children = children
            result.hasChildren = true
        } else {
            result.hasChildren = false
        }
        return result
    })
}

const fetchData = async () => {
    loadingText.value = '加载中...'
    loading.value = true
    try {
        const res: any = await getApps()
        const flatData = res.data || res || []
        allData.value = flatData
        treeData.value = buildTree(flatData)
    } catch (error) {
        console.error('加载应用列表失败:', error)
        ElMessage.error('加载列表失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => {
    if (!searchQuery.value) {
        treeData.value = buildTree(allData.value)
        return
    }
    const query = searchQuery.value.toLowerCase()
    const filterList = allData.value.filter(item =>
        (item.app_name || '').toLowerCase().includes(query) ||
        (item.app_code || '').toLowerCase().includes(query)
    )
    const allParentIds = new Set<string>()
    const findAllParents = (pid: string) => {
        if (!pid || allParentIds.has(pid)) return
        allParentIds.add(pid)
        const parent = allData.value.find(item => item.id === pid)
        if (parent) {
            const parentPid = parent.parent_id || ''
            if (parentPid) findAllParents(parentPid)
        }
    }
    filterList.forEach(item => {
        const pid = item.parent_id || ''
        if (pid) findAllParents(pid)
    })
    const filterIds = new Set(filterList.map(item => item.id))
    const parentItems = allData.value.filter(item => allParentIds.has(item.id) && !filterIds.has(item.id))
    const fullList = [...filterList, ...parentItems]
    treeData.value = buildTree(fullList)
}

const handleReset = () => {
    searchQuery.value = ''
    treeData.value = buildTree(allData.value)
}

const handleCreate = () => {
    formData.value = { parent_id: '', status: 1, sort: 0 }
    dialogVisible.value = true
}

const handleAddChild = (row: any) => {
    formData.value = { parent_id: row.id, status: 1, sort: 0 }
    dialogVisible.value = true
}

const handleEdit = (row: any) => {
    formData.value = { ...row }
    dialogVisible.value = true
}

const hasChildren = (id: string) => {
    return allData.value.some(item => item.parent_id === id)
}

const handleDelete = (row: any) => {
    ElMessageBox.confirm(`确定要删除应用 "${row.app_name}" 吗？`, '提示', { type: 'warning' })
        .then(async () => {
            try {
                await deleteApp(row.id)
                ElMessage.success('删除成功')
                fetchData()
            } catch (error: any) { ElMessage.error(error.message || '删除失败') }
        }).catch(() => { })
}

const handleFormSuccess = () => {
    fetchData()
}

onMounted(() => { fetchData() })
</script>
