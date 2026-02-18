<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon">
                    <Menu />
                </el-icon>
                菜单管理
            </h1>
            <div class="header-actions">
                <el-button type="primary" @click="handleCreate" :icon="Plus">新增菜单</el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-tree-select v-model="filterAppCode" :data="appTreeData" check-strictly :render-after-expand="false" placeholder="选择应用" clearable style="width: 200px" />
                <el-input v-model="searchQuery" placeholder="请输入菜单名称搜索" clearable :prefix-icon="Search" style="width: 300px; margin-left: 10px" @input="handleDebouncedSearch" />
                <el-select v-model="filterType" placeholder="菜单类型" style="width: 150px; margin-left: 10px" clearable @change="handleSearch">
                    <el-option label="全部" value="" />
                    <el-option label="目录" value="M" />
                    <el-option label="菜单" value="C" />
                    <el-option label="按钮" value="F" />
                    <el-option label="资源" value="Z" />
                </el-select>
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
            </div>
            <div class="table-area">
                <el-table class="tree-table" v-loading="loading" :data="displayData" border stripe row-key="id" :tree-props="{ children: 'children', hasChildren: 'hasChildren' }" :indent="24" default-expand-all style="width: 100%; height: 100%;">
                    <el-table-column prop="menu_name" label="菜单名称" min-width="260">
                        <template #default="scope">
                            <TreeNameCell :row="scope.row" name-field="menu_name" />
                        </template>
                    </el-table-column>
                    <el-table-column prop="menu_name" label="菜单名称" width="200" />
                    <el-table-column prop="icon" label="图标" width="60">
                        <template #default="scope">
                            <font-awesome-icon :icon="parseIcon(scope.row.icon)" v-if="scope.row.icon && parseIcon(scope.row.icon)" />
                        </template>
                    </el-table-column>
                    <el-table-column prop="menu_code" label="菜单标识" width="260" show-overflow-tooltip />
                    <el-table-column prop="menu_type" label="类型" width="80">
                        <template #default="scope">
                            <el-tag v-if="scope.row.menu_type === 'M'" type="warning">目录</el-tag>
                            <el-tag v-else-if="scope.row.menu_type === 'C'" type="success">菜单</el-tag>
                            <el-tag v-else-if="scope.row.menu_type === 'F'" type="info">按钮</el-tag>
                            <el-tag v-else type="primary">资源</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="url" label="路由地址" min-width="200" show-overflow-tooltip />
                    <el-table-column prop="is_visible" label="显示" width="60">
                        <template #default="scope">
                            <el-tag v-if="scope.row.is_visible" type="success">显示</el-tag>
                            <el-tag v-else type="info">隐藏</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="status" label="状态" width="60">
                        <template #default="scope">
                            <el-tag v-if="scope.row.status === 1" type="success">启用</el-tag>
                            <el-tag v-else type="danger">禁用</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="sort" label="排序" width="60" />
                    <el-table-column label="操作" width="300" fixed="right">
                        <template #default="scope">
                            <el-button type="primary" link @click="handleCreateChild(scope.row)">新增子级</el-button>
                            <el-button type="primary" link @click="handleEdit(scope.row)">编辑</el-button>
                            <el-button type="danger" link @click="handleDelete(scope.row)" v-if="!hasChildren(scope.row.id)">删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
        </el-card>
        <MenuForm v-model="dialogVisible" :data="formData" :menu-list="allData" :app-list="appList" :org-list="orgList" @success="loadData" />
    </div>
</template>
<script setup lang="ts">
import { deleteMenu, getApps, getMenus, getUnits } from '@/api/user'
import TreeNameCell from '@/components/table/TreeNameCell.vue'
import { Menu, Plus, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import MenuForm from './MenuForm.vue'

const parseIcon = (iconName: string) => {
    if (!iconName) return null

    if (iconName.startsWith('fa-solid ') || iconName.startsWith('fa-regular ') || iconName.startsWith('fa-brands ')) {
        return iconName
    }

    if (iconName.startsWith('fa-')) {
        return `fa-solid ${iconName}`
    }

    return null
}

const searchQuery = ref('')
const filterType = ref('')
const filterAppCode = ref('')
const loading = ref(false)
const allData = ref<any[]>([])
const appList = ref<any[]>([])
const orgList = ref<any[]>([])

const appTreeData = computed(() => {
    const buildTree = (items: any[], parentId = ''): any[] => {
        return items
            .filter(item => (parentId === '' ? !item.parent_id || item.parent_id === '' : item.parent_id === parentId))
            .sort((a, b) => (a.sort || 0) - (b.sort || 0))
            .map(item => ({
                value: item.app_code,
                label: item.app_name,
                children: buildTree(items, item.id)
            }))
    }
    return buildTree(appList.value)
})

const displayData = computed(() => {
    let result = allData.value
    if (filterAppCode.value) {
        result = filterByAppCode(result, filterAppCode.value)
    }
    if (searchQuery.value) {
        result = filterTreeByName(result, searchQuery.value)
    }
    if (filterType.value) {
        result = filterTreeByType(result, filterType.value)
    }
    return buildTree(result)
})

const filterByAppCode = (items: any[], targetAppCode: string): any[] => {
    return items.filter(item => item.app_code === targetAppCode)
}

const buildTree = (items: any[], parentId = ''): any[] => {
    return items
        .filter(item => (parentId === '' ? !item.parent_id || item.parent_id === '' : item.parent_id === parentId))
        .sort((a, b) => a.sort - b.sort)
        .map(item => ({
            ...item,
            children: buildTree(items, item.id)
        }))
}

const filterTreeByName = (items: any[], query: string): any[] => {
    const lowerQuery = query.toLowerCase()
    return items.filter(item => item.menu_name?.toLowerCase().includes(lowerQuery))
}

const filterTreeByType = (items: any[], type: string): any[] => {
    return items.filter(item => item.menu_type === type)
}

const dialogVisible = ref(false)
const formData = ref<any>({})

const loadData = async () => {
    loading.value = true
    try {
        const [menus, apps, orgs] = await Promise.all([getMenus(), getApps(), getUnits()])
        allData.value = menus
        appList.value = apps
        orgList.value = orgs
    } catch (error: any) {
        ElMessage.error(error.message)
    } finally {
        loading.value = false
    }
}

const handleSearch = () => { }

const handleDebouncedSearch = () => { }

const handleCreate = () => {
    formData.value = {
        status: 1,
        sort: 0,
        is_visible: true,
        menu_type: 'C',
        parent_id: '',
        app_code: filterAppCode.value || ''
    }
    dialogVisible.value = true
}

const handleCreateChild = (row: any) => {
    formData.value = {
        status: 1,
        sort: 0,
        is_visible: true,
        menu_type: 'C',
        parent_id: row.id,
        app_code: row.app_code
    }
    dialogVisible.value = true
}

const handleEdit = (row: any) => {
    formData.value = { ...row }
    dialogVisible.value = true
}

const hasChildren = (id: string) => {
    return allData.value.some(item => item.parent_id === id)
}

const handleDelete = async (row: any) => {
    try {
        await ElMessageBox.confirm('确定要删除该菜单吗？', '提示', { type: 'warning' })
        await deleteMenu(row.id)
        ElMessage.success('删除成功')
        loadData()
    } catch (error: any) {
        if (error !== 'cancel') ElMessage.error(error.message)
    }
}

onMounted(() => loadData())
</script>
<style scoped>
.sso-page {
    height: 100%;
    display: flex;
    flex-direction: column;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
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
}

.main-card {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.search-area {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
    flex-wrap: wrap;
    gap: 10px;
}

.table-area {
    flex: 1;
}

.text-primary {
    color: var(--el-text-color-primary);
}
</style>
