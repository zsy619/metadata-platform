<template>
  <div class="container-padding">
    <div class="page-header">
      <h1 class="text-primary page-title">
        <el-icon class="title-icon"><User /></el-icon>
        用户组管理
      </h1>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate" :icon="Plus">新增用户组</el-button>
      </div>
    </div>
    <el-card class="main-card">
      <div class="search-area">
        <el-input v-model="searchQuery" placeholder="请输入用户组名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
        <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
        <el-button @click="handleReset" :icon="RefreshLeft">重置</el-button>
      </div>
      <div class="table-area">
        <el-table
          class="tree-table"
          v-loading="loading"
          :element-loading-text="loadingText"
          :data="treeData"
          border
          stripe
          row-key="id"
          :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
          :indent="24"
          default-expand-all
          style="width: 100%; height: 100%;"
        >
          <template #empty>
            <el-empty :description="searchQuery ? '未搜索到相关用户组' : '暂无用户组'">
              <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增用户组</el-button>
            </el-empty>
          </template>
          <el-table-column prop="group_name" label="用户组名称" min-width="200">
            <template #default="scope">
              <TreeNameCell :row="scope.row" name-field="group_name" />
            </template>
          </el-table-column>
          <el-table-column prop="group_code" label="用户组编码" width="150" />
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
              <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg v-if="!scope.row.hasChildren">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
    <UserGroupForm v-model="dialogVisible" :data="formData" :group-tree-data="groupTreeSelectData" :exclude-ids="excludeIds" @success="loadData" />
    <UserGroupRoleDialog v-model="roleDialogVisible" :group-id="currentGroupId" :group-name="currentGroupName" @success="loadData" />
  </div>
</template>

<script setup lang="ts">
import { deleteUserGroup, getAllUserGroups } from '@/api/user'
import TreeNameCell from '@/components/table/TreeNameCell.vue'
import { Delete, Edit, Plus, RefreshLeft, Search, User } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import UserGroupForm from './UserGroupForm.vue'
import UserGroupRoleDialog from './UserGroupRoleDialog.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')

const allData = ref<any[]>([])

const dialogVisible = ref(false)
const formData = ref<any>({})
const excludeIds = ref<string[]>([])

// 用户组关联角色弹窗相关
const roleDialogVisible = ref(false)
const currentGroupId = ref('')
const currentGroupName = ref('')

const isRoot = (pid: string) => !pid || pid === '' || pid === '0'

/**
 * 构建树型表格数据
 */
const buildTableTree = (items: any[], parentId = '', level = 0): any[] => {
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

/**
 * 树型表格数据
 */
const treeData = computed(() => {
  if (!Array.isArray(allData.value)) return []
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    const filtered = allData.value.filter(item => (item.group_name || '').toLowerCase().includes(query))
    
    // 收集所有父节点ID
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
    filtered.forEach(item => {
      const pid = item.parent_id || ''
      if (pid && pid !== '' && pid !== '0') findAllParents(pid)
    })
    const filteredIds = new Set(filtered.map(item => item.id))
    const fullList = [...filtered, ...allData.value.filter(item => allParentIds.has(item.id) && !filteredIds.has(item.id))]
    return buildTableTree(fullList, '')
  }
  return buildTableTree(allData.value, '')
})

/**
 * 构建树型选择器数据结构
 */
const buildTreeSelectData = (items: any[], parentId = '', excludeIdsList: string[] = []): any[] => {
  if (!Array.isArray(items)) return []
  return items
    .filter(item => (item.parent_id || '') === parentId && !excludeIdsList.includes(item.id))
    .map(item => ({
      value: item.id,
      label: item.group_name,
      children: buildTreeSelectData(items, item.id, excludeIdsList)
    }))
}

/**
 * 获取指定节点的所有子节点ID（递归）
 */
const getAllChildrenIds = (items: any[], parentId: string): string[] => {
  if (!Array.isArray(items)) return []
  const children = items.filter(item => (item.parent_id || '') === parentId)
  let ids: string[] = []
  for (const child of children) {
    ids.push(child.id)
    ids = ids.concat(getAllChildrenIds(items, child.id))
  }
  return ids
}

/**
 * 树型选择器数据
 */
const groupTreeSelectData = computed(() => buildTreeSelectData(allData.value))

const loadData = async () => {
  loadingText.value = '加载中...'
  loading.value = true
  try {
    const res: any = await getAllUserGroups()
    const data = res.data || res
    allData.value = Array.isArray(data) ? data : []
  } catch (error) {
    console.error('加载用户组列表失败:', error)
    ElMessage.error('加载列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {}
const handleDebouncedSearch = () => {}
const handleReset = () => { searchQuery.value = '' }

const handleCreate = () => {
  formData.value = { status: 1, sort: 0, parent_id: '', group_name: '', group_code: '', remark: '' }
  excludeIds.value = []
  dialogVisible.value = true
}

const handleCreateChild = (row: any) => {
  formData.value = { status: 1, sort: 0, parent_id: row.id, group_name: '', group_code: '', remark: '' }
  excludeIds.value = []
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  formData.value = { ...row }
  excludeIds.value = [row.id, ...getAllChildrenIds(allData.value, row.id)]
  dialogVisible.value = true
}

/**
 * 管理用户组关联角色
 * @param row 用户组数据
 */
const handleManageRoles = (row: any) => {
  currentGroupId.value = row.id
  currentGroupName.value = row.group_name
  roleDialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除用户组 "${row.group_name}" 吗？`, '提示', { type: 'warning' })
    await deleteUserGroup(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error: any) { if (error !== 'cancel') ElMessage.error(error.message || '删除失败') }
}

onMounted(() => loadData())
</script>

