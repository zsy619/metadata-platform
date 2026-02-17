<template>
  <div class="container-padding">
    <div class="page-header">
      <h1 class="text-primary page-title">
        <el-icon class="title-icon"><UserFilled /></el-icon>
        角色管理
      </h1>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate" :icon="Plus">新增角色</el-button>
      </div>
    </div>
    <el-card class="main-card">
      <div class="search-area">
        <el-input v-model="searchQuery" placeholder="请输入角色名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
        <el-select v-model="filterStatus" placeholder="筛选状态" style="width: 150px; margin-left: 10px" clearable @change="handleSearch">
          <el-option label="全部" value="" />
          <el-option label="有效" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
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
          style="width: 100%; height: 100%;"
          row-key="id"
          :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
          :indent="24"
          default-expand-all
        >
          <template #empty>
            <el-empty :description="searchQuery ? '未搜索到相关角色' : '暂无角色'">
              <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增角色</el-button>
            </el-empty>
          </template>
          <el-table-column prop="role_name" label="角色名称" min-width="200">
            <template #default="scope">
              <div class="role-name-cell">
                <el-icon v-if="scope.row.children && scope.row.children.length > 0" class="tree-icon folder-icon">
                  <Folder />
                </el-icon>
                <el-icon v-else class="tree-icon leaf-icon"><User /></el-icon>
                <span class="role-name-text">{{ scope.row.role_name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="role_code" label="角色编码" width="150" />
          <el-table-column prop="data_range" label="数据范围" width="100">
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
          <el-table-column prop="remark" label="备注" show-overflow-tooltip />
          <el-table-column prop="sort" label="序号" width="80" />
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

    <RoleForm v-model="dialogVisible" :data="formData" :all-roles="allData" @success="loadData" />
  </div>
</template>

<script setup lang="ts">
import { deleteRole, getRoles } from '@/api/user'
import { DATA_RANGE, DATA_RANGE_LABELS } from '@/utils/constants'
import { Delete, Edit, Folder, Plus, RefreshLeft, Search, User, UserFilled } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, ref } from 'vue'
import RoleForm from './RoleForm.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')
const filterStatus = ref<number | ''>('')

const allData = ref<any[]>([])
const treeData = ref<any[]>([])

const dialogVisible = ref(false)
const formData = ref<any>({})

/**
 * 将扁平数据转换为树形结构
 * @param list 扁平数据列表
 * @param parentId 父ID，默认为空字符串（顶级节点）
 * @returns 树形结构数据
 */
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
    }
    return result
  })
}

/**
 * 加载角色列表数据
 */
const loadData = async () => {
  loadingText.value = '加载中...'
  loading.value = true
  try {
    const res: any = await getRoles()
    const flatData = res.data || res || []
    allData.value = flatData
    treeData.value = buildTree(flatData)
  } catch (error) {
    console.error('加载角色列表失败:', error)
    ElMessage.error('加载列表失败')
  } finally {
    loading.value = false
  }
}

/**
 * 搜索处理
 */
const handleSearch = () => {
  if (!searchQuery.value && filterStatus.value === '') {
    treeData.value = buildTree(allData.value)
    return
  }
  
  const query = searchQuery.value.toLowerCase()
  const filterList = allData.value.filter(item => {
    const nameMatch = (item.role_name || '').toLowerCase().includes(query)
    const statusMatch = filterStatus.value === '' || item.status === filterStatus.value
    return nameMatch && statusMatch
  })
  
  // 收集所有父节点ID，确保搜索结果能显示完整路径
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
  
  const fullList = [...filterList, ...allData.value.filter(item => allParentIds.has(item.id))]
  treeData.value = buildTree(fullList)
}

const handleDebouncedSearch = () => {
  handleSearch()
}

const handleReset = () => {
  searchQuery.value = ''
  filterStatus.value = ''
  treeData.value = buildTree(allData.value)
}

/**
 * 新增顶级角色
 */
const handleCreate = () => {
  formData.value = { parent_id: '', status: 1, sort: 0, data_range: DATA_RANGE.ALL }
  dialogVisible.value = true
}

/**
 * 新增子级角色
 * @param row 父级角色数据
 */
const handleAddChild = (row: any) => {
  formData.value = { parent_id: row.id, status: 1, sort: 0, data_range: DATA_RANGE.ALL }
  dialogVisible.value = true
}

/**
 * 编辑角色
 * @param row 角色数据
 */
const handleEdit = (row: any) => {
  formData.value = { ...row }
  dialogVisible.value = true
}

/**
 * 删除角色
 * @param row 角色数据
 */
const handleDelete = async (row: any) => {
  // 检查是否存在子节点
  const hasChildren = allData.value.some(item => item.parent_id === row.id)
  const msg = hasChildren
    ? `该角色下存在子级角色，删除将同时删除所有子级角色。确定要删除角色 "${row.role_name}" 吗？`
    : `确定要删除角色 "${row.role_name}" 吗？`
  
  try {
    await ElMessageBox.confirm(msg, '提示', { type: 'warning' })
    await deleteRole(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

onMounted(() => loadData())
</script>

<style scoped>
.role-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}
.role-name-text {
  font-weight: 500;
}
.tree-icon {
  font-size: 16px;
  flex-shrink: 0;
}
.folder-icon {
  color: var(--el-color-warning);
}
.leaf-icon {
  color: var(--el-text-color-secondary);
}
.main-card { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
:deep(.el-card__body) { height: 100%; display: flex; flex-direction: column; padding: 20px; overflow: hidden; box-sizing: border-box; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; flex-shrink: 0; }
.page-title { font-size: 20px; font-weight: 600; display: flex; align-items: center; gap: 8px; }
.title-icon { font-size: 24px; color: var(--el-color-primary); }
.header-actions { display: flex; gap: 10px; }
.search-area { display: flex; align-items: center; margin-bottom: 20px; flex-wrap: wrap; gap: 10px; }
.table-area { flex: 1; overflow: hidden; }
.text-primary { color: var(--el-text-color-primary); }

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
