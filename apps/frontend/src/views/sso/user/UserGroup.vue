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
          v-loading="loading"
          :element-loading-text="loadingText"
          :data="treeData"
          border
          stripe
          row-key="id"
          :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
          default-expand-all
          style="width: 100%; height: 100%;"
        >
          <template #empty>
            <el-empty :description="searchQuery ? '未搜索到相关用户组' : '暂无用户组'">
              <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增用户组</el-button>
            </el-empty>
          </template>
          <el-table-column prop="group_name" label="用户组名称" width="220" show-overflow-tooltip />
          <el-table-column prop="group_code" label="用户组编码" width="150" />
          <el-table-column prop="status" label="状态" width="80">
            <template #default="scope">
              <el-tag v-if="scope.row.status === 1" type="success">有效</el-tag>
              <el-tag v-else type="danger">禁用</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="sort" label="排序" width="80" />
          <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
          <el-table-column label="操作" width="250" fixed="right">
            <template #default="scope">
              <el-button type="primary" size="small" :icon="Plus" @click="handleCreateChild(scope.row)" text bg>新增子级</el-button>
              <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>编辑</el-button>
              <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg v-if="!hasChildren(scope.row.id)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
    <UserGroupForm v-model="dialogVisible" :data="formData" :group-tree-data="groupTreeSelectData" :exclude-ids="excludeIds" @success="loadData" />
  </div>
</template>

<script setup lang="ts">
import { deleteUserGroup, getUserGroups } from '@/api/user'
import { Delete, Edit, Plus, RefreshLeft, Search, User } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import UserGroupForm from './UserGroupForm.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')

const allData = ref<any[]>([])

/**
 * 构建树型数据结构
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
 * 树型表格数据
 */
const treeData = computed(() => {
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    return allData.value.filter(item => (item.group_name || '').toLowerCase().includes(query))
  }
  return buildTree(allData.value)
})

/**
 * 构建树型选择器数据结构
 */
const buildTreeSelectData = (items: any[], parentId = '', excludeIds: string[] = []): any[] => {
  return items
    .filter(item => (item.parent_id || '') === parentId && !excludeIds.includes(item.id))
    .map(item => ({
      value: item.id,
      label: item.group_name,
      children: buildTreeSelectData(items, item.id, excludeIds)
    }))
}

/**
 * 获取指定节点的所有子节点ID（递归）
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
 * 树型选择器数据
 */
const groupTreeSelectData = computed(() => buildTreeSelectData(allData.value))

/**
 * 编辑时需要排除的节点ID列表
 */
const excludeIds = ref<string[]>([])

const dialogVisible = ref(false)
const formData = ref<any>({})

const loadData = async () => {
  loadingText.value = '加载中...'
  loading.value = true
  try {
    const res: any = await getUserGroups()
    allData.value = res.data || res
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
  formData.value = { status: 1, sort: 0, parent_id: '' }
  excludeIds.value = []
  dialogVisible.value = true
}

const handleCreateChild = (row: any) => {
  formData.value = { status: 1, sort: 0, parent_id: row.id }
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
    await ElMessageBox.confirm(`确定要删除用户组 "${row.group_name}" 吗？`, '提示', { type: 'warning' })
    await deleteUserGroup(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error: any) { if (error !== 'cancel') ElMessage.error(error.message || '删除失败') }
}

onMounted(() => loadData())
</script>

<style scoped>
.main-card { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
:deep(.el-card__body) { height: 100%; display: flex; flex-direction: column; padding: 20px; overflow: hidden; box-sizing: border-box; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; flex-shrink: 0; }
.page-title { font-size: 20px; font-weight: 600; display: flex; align-items: center; gap: 8px; }
.title-icon { font-size: 24px; color: var(--el-color-primary); }
.header-actions { display: flex; gap: 10px; }
.search-area { display: flex; margin-bottom: 20px; flex-shrink: 0; }
.table-area { flex: 1; overflow: hidden; }
.text-primary { color: var(--el-text-color-primary); }
</style>
