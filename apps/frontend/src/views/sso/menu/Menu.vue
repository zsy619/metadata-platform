<template>
  <div class="container-padding">
    <div class="page-header">
      <h1 class="text-primary page-title">
        <el-icon class="title-icon"><Menu /></el-icon>
        菜单管理
      </h1>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate" :icon="Plus">新增菜单</el-button>
      </div>
    </div>
    <el-card class="main-card">
      <div class="search-area">
        <el-select v-model="filterAppCode" placeholder="选择应用" style="width: 200px" clearable @change="handleSearch">
          <el-option label="全部应用" value="" />
          <el-option v-for="app in appList" :key="app.id" :label="app.app_name" :value="app.app_code" />
        </el-select>
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
        <el-table v-loading="loading" :data="displayData" border stripe row-key="id" :tree-props="{ children: 'children', hasChildren: 'hasChildren' }" default-expand-all>
          <el-table-column prop="menu_name" label="菜单名称" width="200" />
          <el-table-column prop="icon" label="图标" width="60">
            <template #default="scope">
              <el-icon v-if="scope.row.icon"><component :is="scope.row.icon" /></el-icon>
            </template>
          </el-table-column>
          <el-table-column prop="menu_type" label="类型" width="80">
            <template #default="scope">
              <el-tag v-if="scope.row.menu_type === 'M'" type="warning">目录</el-tag>
              <el-tag v-else-if="scope.row.menu_type === 'C'" type="success">菜单</el-tag>
              <el-tag v-else-if="scope.row.menu_type === 'F'" type="info">按钮</el-tag>
              <el-tag v-else type="primary">资源</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="url" label="路由地址" show-overflow-tooltip />
          <el-table-column prop="menu_code" label="权限标识" width="150" show-overflow-tooltip />
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
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="scope">
              <el-button type="primary" link @click="handleEdit(scope.row)">编辑</el-button>
              <el-button type="danger" link @click="handleDelete(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px" label-position="right">
        <el-form-item label="所属应用" prop="app_code">
          <el-select v-model="formData.app_code" placeholder="请选择所属应用" style="width: 100%" @change="handleAppChange">
            <el-option v-for="app in appList" :key="app.id" :label="app.app_name" :value="app.app_code" />
          </el-select>
        </el-form-item>
        <el-form-item label="上级菜单" prop="parent_id">
          <el-tree-select
            v-model="formData.parent_id"
            :data="filteredMenuTreeData"
            check-strictly
            :render-after-expand="false"
            placeholder="请选择上级菜单（不选则为顶级菜单）"
            style="width: 100%"
            clearable
          />
        </el-form-item>
        <el-form-item label="菜单类型" prop="menu_type">
          <el-radio-group v-model="formData.menu_type">
            <el-radio label="M">目录</el-radio>
            <el-radio label="C">菜单</el-radio>
            <el-radio label="F">按钮</el-radio>
            <el-radio label="Z">资源</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="菜单名称" prop="menu_name">
          <el-input v-model="formData.menu_name" placeholder="请输入菜单名称" />
        </el-form-item>
        <el-form-item label="权限标识" prop="menu_code">
          <el-input v-model="formData.menu_code" placeholder="请输入权限标识" />
        </el-form-item>
        <el-form-item v-if="formData.menu_type !== 'F' && formData.menu_type !== 'Z'" label="路由地址" prop="url">
          <el-input v-model="formData.url" placeholder="请输入路由地址" />
        </el-form-item>
        <el-form-item v-if="formData.menu_type !== 'F' && formData.menu_type !== 'Z'" label="图标" prop="icon">
          <el-input v-model="formData.icon" placeholder="请输入图标名称" />
        </el-form-item>
        <el-form-item v-if="formData.menu_type !== 'F' && formData.menu_type !== 'Z'" label="是否可见" prop="is_visible">
          <el-switch v-model="formData.is_visible" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="2" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { createMenu, deleteMenu, getMenus, updateMenu, getApps } from '@/api/user'
import { Menu, Plus, Search } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'

const searchQuery = ref('')
const filterType = ref('')
const filterAppCode = ref('')
const loading = ref(false)
const allData = ref<any[]>([])
const appList = ref<any[]>([])

const displayData = computed(() => {
  let result = allData.value
  if (filterAppCode.value) {
    result = result.filter(item => item.app_code === filterAppCode.value)
  }
  if (searchQuery.value) {
    result = filterTreeByName(result, searchQuery.value)
  }
  if (filterType.value) {
    result = filterTreeByType(result, filterType.value)
  }
  return buildTree(result)
})

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

const getAllDescendantIds = (parentId: string): string[] => {
  const ids: string[] = [parentId]
  const findChildren = (pid: string) => {
    allData.value.forEach(item => {
      if (item.parent_id === pid) {
        ids.push(item.id)
        findChildren(item.id)
      }
    })
  }
  findChildren(parentId)
  return ids
}

const filteredMenuTreeData = computed(() => {
  const currentId = formData.value.id
  const currentAppCode = formData.value.app_code

  let filteredItems = allData.value.filter(item => {
    if (!currentAppCode) return true
    return item.app_code === currentAppCode
  })

  if (currentId) {
    const excludeIds = getAllDescendantIds(currentId)
    filteredItems = filteredItems.filter(item => !excludeIds.includes(item.id))
  }

  const buildTreeSelect = (items: any[], parentId = ''): any[] => {
    return items
      .filter(item => (parentId === '' ? !item.parent_id || item.parent_id === '' : item.parent_id === parentId))
      .sort((a, b) => a.sort - b.sort)
      .map(item => ({
        value: item.id,
        label: item.menu_name,
        children: buildTreeSelect(items, item.id)
      }))
  }

  const treeData = buildTreeSelect(filteredItems)
  return [
    { value: '', label: '顶级菜单', children: [] },
    ...treeData
  ]
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()
const formData = ref<any>({})
const submitLoading = ref(false)

const formRules: FormRules = {
  app_code: [{ required: true, message: '请选择所属应用', trigger: 'change' }],
  menu_name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  menu_code: [{ required: true, message: '请输入权限标识', trigger: 'blur' }],
  menu_type: [{ required: true, message: '请选择菜单类型', trigger: 'change' }]
}

const loadData = async () => {
  loading.value = true
  try {
    const [menus, apps] = await Promise.all([getMenus(), getApps()])
    allData.value = menus
    appList.value = buildAppTree(apps)
  } catch (error: any) {
    ElMessage.error(error.message)
  } finally {
    loading.value = false
  }
}

const buildAppTree = (apps: any[]): any[] => {
  const buildTree = (items: any[], parentId = ''): any[] => {
    return items
      .filter(item => (parentId === '' ? !item.parent_id || item.parent_id === '' : item.parent_id === parentId))
      .sort((a, b) => a.sort - b.sort)
      .map(item => ({
        ...item,
        children: buildTree(items, item.id)
      }))
  }
  return buildTree(apps)
}

const handleSearch = () => {}

const handleDebouncedSearch = () => {}

const handleCreate = () => {
  dialogTitle.value = '新增菜单'
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

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑菜单'
  formData.value = { ...row }
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  const hasChildren = allData.value.some(item => item.parent_id === row.id)
  if (hasChildren) {
    ElMessage.warning('该菜单下存在子菜单，请先删除子菜单')
    return
  }
  try {
    await ElMessageBox.confirm('确定要删除该菜单吗？', '提示', { type: 'warning' })
    await deleteMenu(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') ElMessage.error(error.message)
  }
}

const handleAppChange = () => {
  formData.value.parent_id = ''
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        if (formData.value.id) {
          await updateMenu(formData.value.id, formData.value)
          ElMessage.success('更新成功')
        } else {
          await createMenu(formData.value)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        loadData()
      } catch (error: any) {
        ElMessage.error(error.message)
      } finally {
        submitLoading.value = false
      }
    }
  })
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
