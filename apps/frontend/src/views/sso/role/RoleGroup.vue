<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon"><Folder /></el-icon>
                角色组管理
            </h1>
            <div class="header-actions">
                <el-button type="primary" @click="handleCreate" :icon="Plus">新增角色组</el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-input v-model="searchQuery" placeholder="请输入角色组名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
                <el-select v-model="filterStatus" placeholder="状态筛选" clearable style="width: 120px; margin-left: 10px">
                    <el-option label="有效" :value="1" />
                    <el-option label="禁用" :value="0" />
                </el-select>
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
                <el-button @click="handleReset" :icon="RefreshLeft">重置</el-button>
            </div>
            <div class="table-area">
                <el-table
                    class="tree-table"
                    v-loading="loading"
                    :element-loading-text="loadingText"
                    :data="tableData"
                    border
                    stripe
                    row-key="id"
                    :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
                    :indent="24"
                    default-expand-all
                    style="width: 100%; height: 100%"
                >
                    <template #empty>
                        <el-empty :description="searchQuery ? '未搜索到相关角色组' : '暂无角色组'">
                            <el-button v-if="!searchQuery && filterStatus === ''" type="primary" @click="handleCreate">新增角色组</el-button>
                        </el-empty>
                    </template>
                    <el-table-column label="角色组名称" min-width="200">
                        <template #default="scope">
                            <TreeNameCell :row="scope.row" name-field="group_name" />
                        </template>
                    </el-table-column>
                    <el-table-column prop="group_code" label="角色组编码" width="150" />
                    <el-table-column prop="status" label="状态" width="80">
                        <template #default="scope">
                            <el-tag v-if="scope.row.status === 1" type="success">有效</el-tag>
                            <el-tag v-else type="danger">禁用</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="sort" label="序号" width="80" />
                    <el-table-column prop="remark" label="备注" show-overflow-tooltip />
                    <el-table-column label="操作" width="320" fixed="right">
                        <template #default="scope">
                            <el-button type="primary" size="small" :icon="Plus" @click="handleAddChild(scope.row)" text bg>新增子级</el-button>
                            <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>编辑</el-button>
                            <el-button type="warning" size="small" @click="handleManageRoles(scope.row)" text bg>关联角色</el-button>
                            <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg v-if="!scope.row.is_system && !scope.row.hasChildren">删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
        </el-card>
        <RoleGroupForm v-model="dialogVisible" :data="formData" :all-groups="allData" :exclude-ids="excludeIds" @success="loadData" />
        <RoleGroupRoleDialog v-model="roleDialogVisible" :group-id="currentGroupId" :group-name="currentGroupName" @success="loadData" />
    </div>
</template>

<script setup lang="ts">
import { deleteRoleGroup, getRoleGroups } from '@/api/user'
import TreeNameCell from '@/components/table/TreeNameCell.vue'
import { Delete, Edit, Folder, Plus, RefreshLeft, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import RoleGroupForm from './RoleGroupForm.vue'
import RoleGroupRoleDialog from './RoleGroupRoleDialog.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')
const filterStatus = ref<number | ''>('')

const allData = ref<any[]>([])

const dialogVisible = ref(false)
const formData = ref<any>({})

const excludeIds = ref<string[]>([])

// 角色组关联角色弹窗相关
const roleDialogVisible = ref(false)
const currentGroupId = ref('')
const currentGroupName = ref('')

const isRoot = (pid: string) => !pid || pid === '' || pid === '0'

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

const tableData = computed(() => {
    const buildTableTree = (items: any[], parentId: string, level = 0): any[] => {
        return items
            .filter(item => {
                const itemPid = item.parent_id || ''
                if (isRoot(parentId)) {
                    return isRoot(itemPid)
                }
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
        data = data.filter(item => (item.group_name || '').toLowerCase().includes(query))
    }
    if (filterStatus.value !== '') {
        data = data.filter(item => item.status === filterStatus.value)
    }

    if (!searchQuery.value && filterStatus.value === '') {
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

const loadData = async () => {
    loadingText.value = '加载中...'
    loading.value = true
    try {
        const res: any = await getRoleGroups()
        allData.value = res.data || res || []
    } catch (error) {
        console.error('加载角色组列表失败:', error)
        ElMessage.error('加载列表失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => { }
const handleDebouncedSearch = () => { }
const handleReset = () => { searchQuery.value = ''; filterStatus.value = '' }

const handleCreate = () => {
    formData.value = { parent_id: '', status: 1, sort: 0 }
    excludeIds.value = []
    dialogVisible.value = true
}

const handleAddChild = (row: any) => {
    formData.value = { parent_id: row.id, status: 1, sort: 0 }
    excludeIds.value = []
    dialogVisible.value = true
}

const handleEdit = (row: any) => {
    formData.value = { ...row }
    excludeIds.value = getAllDescendantIds(row.id)
    dialogVisible.value = true
}

/**
 * 管理角色组关联角色
 * @param row 角色组数据
 */
const handleManageRoles = (row: any) => {
    currentGroupId.value = row.id
    currentGroupName.value = row.group_name
    roleDialogVisible.value = true
}

const handleDelete = async (row: any) => {
    try {
        await ElMessageBox.confirm(`确定要删除角色组 "${row.group_name}" 吗？`, '提示', { type: 'warning' })
        await deleteRoleGroup(row.id)
        ElMessage.success('删除成功')
        loadData()
    } catch (error: any) { if (error !== 'cancel') ElMessage.error(error.message || '删除失败') }
}

onMounted(() => loadData())
</script>


