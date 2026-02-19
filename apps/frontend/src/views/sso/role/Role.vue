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
          :data="tableData"
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
              <TreeNameCell :row="scope.row" name-field="role_name" />
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
          <el-table-column label="操作" width="360" fixed="right">
            <template #default="scope">
              <el-button type="primary" size="small" :icon="Plus" @click="handleAddChild(scope.row)" text bg>新增子级</el-button>
              <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>编辑</el-button>
              <el-button type="warning" size="small" @click="handleConfigMenu(scope.row)" text bg>配置菜单</el-button>
              <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg v-if="!scope.row.is_system && !scope.row.hasChildren">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <RoleForm v-model="dialogVisible" :data="formData" :all-roles="allData" :exclude-ids="excludeIds" :org-list="orgList" @success="loadData" />
    <RoleMenuSelect v-model="menuDialogVisible" :role-id="currentRoleId" :role-name="currentRoleName" @success="loadData" />
  </div>
</template>

<script setup lang="ts">
import { deleteRole, getRoles, getUnits } from '@/api/user'
import TreeNameCell from '@/components/table/TreeNameCell.vue'
import { DATA_RANGE, DATA_RANGE_LABELS } from '@/utils/constants'
import { Delete, Edit, Plus, RefreshLeft, Search, UserFilled } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import RoleForm from './RoleForm.vue'
import RoleMenuSelect from './RoleMenuSelect.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')
const filterStatus = ref<number | ''>('')

const allData = ref<any[]>([])
const orgList = ref<any[]>([])

const dialogVisible = ref(false)
const formData = ref<any>({})

/**
 * 编辑时需要排除的节点ID列表（包含自己及所有子节点）
 * 用于上级角色选择器，防止循环引用
 */
const excludeIds = ref<string[]>([])

/**
 * 获取指定节点的所有子节点ID（递归）
 * @param parentId 父节点ID
 * @returns 所有子节点ID列表
 */
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

/**
 * 树型表格数据 - 用于表格展示
 * 支持搜索时保持树形结构完整
 */
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
        data = data.filter(item => (item.role_name || '').toLowerCase().includes(query))
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

/**
 * 加载角色列表数据
 */
const loadData = async () => {
    loadingText.value = '加载中...'
    loading.value = true
    try {
        const [roleRes, orgRes]: [any, any] = await Promise.all([
            getRoles(),
            getUnits()
        ])
        allData.value = roleRes.data || roleRes || []
        orgList.value = orgRes.data || orgRes
    } catch (error) {
        console.error('加载角色列表失败:', error)
        ElMessage.error('加载列表失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => { }
const handleDebouncedSearch = () => { }
const handleReset = () => { searchQuery.value = ''; filterStatus.value = '' }

/**
 * 新增顶级角色
 */
const handleCreate = () => {
    formData.value = { parent_id: '', status: 1, sort: 0, data_range: DATA_RANGE.ALL }
    excludeIds.value = []
    dialogVisible.value = true
}

/**
 * 新增子级角色
 * @param row 父级角色数据
 */
const handleAddChild = (row: any) => {
    formData.value = { parent_id: row.id, status: 1, sort: 0, data_range: DATA_RANGE.ALL }
    excludeIds.value = []
    dialogVisible.value = true
}

/**
 * 编辑角色
 * @param row 角色数据
 */
const handleEdit = (row: any) => {
    formData.value = { ...row }
    excludeIds.value = getAllDescendantIds(row.id)
    dialogVisible.value = true
}

// 菜单配置弹窗相关
const menuDialogVisible = ref(false)
const currentRoleId = ref('')
const currentRoleName = ref('')

/**
 * 配置角色菜单
 * @param row 角色数据
 */
const handleConfigMenu = (row: any) => {
    currentRoleId.value = row.id
    currentRoleName.value = row.role_name
    menuDialogVisible.value = true
}

/**
 * 删除角色
 * @param row 角色数据
 */
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除角色 "${row.role_name}" 吗？`, '提示', { type: 'warning' })
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
