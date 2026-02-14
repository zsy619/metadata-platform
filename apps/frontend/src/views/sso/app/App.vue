<template>
  <div class="container-padding">
    <div class="page-header">
      <h1 class="text-primary page-title">
        <el-icon class="title-icon"><Monitor /></el-icon>
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
        <el-table v-loading="loading" :element-loading-text="loadingText" :data="treeData" border stripe style="width: 100%; height: 100%;" row-key="id" :tree-props="{ children: 'children', hasChildren: 'hasChildren' }">
          <template #empty>
            <el-empty :description="searchQuery ? '未搜索到相关应用' : '暂无应用'">
              <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增应用</el-button>
            </el-empty>
          </template>
          <el-table-column prop="app_name" label="应用名称" min-width="200">
            <template #default="scope">
              <div class="app-name-cell">
                <el-icon v-if="scope.row.children && scope.row.children.length > 0" class="tree-icon">
                  <FolderOpened v-if="scope.row.expanded" />
                  <Folder v-else />
                </el-icon>
                <el-icon v-else class="tree-icon file-icon"><Document /></el-icon>
                <span>{{ scope.row.app_name }}</span>
              </div>
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
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="scope">
              <el-button type="primary" size="small" :icon="Plus" @click="handleAddChild(scope.row)" text bg>新增子级</el-button>
              <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>编辑</el-button>
              <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg>删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px" label-position="right">
        <el-form-item v-if="formData.parent_id" label="父级应用" prop="parent_id">
          <el-input :value="parentAppName" disabled />
        </el-form-item>
        <el-form-item label="应用名称" prop="app_name">
          <el-input v-model="formData.app_name" placeholder="请输入应用名称" />
        </el-form-item>
        <el-form-item label="应用编码" prop="app_code">
          <el-input v-model="formData.app_code" placeholder="请输入应用编码" />
        </el-form-item>
        <el-form-item label="域&nbsp;&nbsp;名/IP" prop="host">
          <el-input v-model="formData.host" placeholder="请输入域名或IP地址" />
        </el-form-item>
        <el-form-item label="状&#12288;&#12288;态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="排&#12288;&#12288;序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" />
        </el-form-item>
        <el-form-item label="备&#12288;&#12288;注" prop="remark">
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
import { createApp, deleteApp, getApps, updateApp } from '@/api/user'
import { Delete, Document, Edit, Folder, FolderOpened, Monitor, Plus, RefreshLeft, Search } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, ref } from 'vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')
const allData = ref<any[]>([])
const treeData = ref<any[]>([])

const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()
const formData = ref<any>({})
const submitLoading = ref(false)
const parentAppName = ref('')

const formRules: FormRules = {
  app_name: [{ required: true, message: '请输入应用名称', trigger: 'blur' }],
  app_code: [{ required: true, message: '请输入应用编码', trigger: 'blur' }]
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return isNaN(date.getTime()) ? '-' : date.toLocaleString('zh-CN')
}

const buildTree = (list: any[], parentId: string = ''): any[] => {
  return list
    .filter(item => item.parent_id === parentId)
    .map(item => ({
      ...item,
      expanded: false,
      children: buildTree(list, item.id)
    }))
    .sort((a, b) => (a.sort || 0) - (b.sort || 0))
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
  const filteredIds = new Set(filterList.map(item => item.id))
  const parentIds = filterList.map(item => item.parent_id)
  const allParentIds = new Set<string>()
  const findAllParents = (pid: string) => {
    if (!pid || allParentIds.has(pid)) return
    allParentIds.add(pid)
    const parent = allData.value.find(item => item.id === pid)
    if (parent && parent.parent_id) findAllParents(parent.parent_id)
  }
  parentIds.forEach(pid => findAllParents(pid))
  const fullList = [...filterList, ...allData.value.filter(item => allParentIds.has(item.id))]
  treeData.value = buildTree(fullList)
}

const handleReset = () => {
  searchQuery.value = ''
  treeData.value = buildTree(allData.value)
}

const handleCreate = () => {
  dialogTitle.value = '新增应用'
  formData.value = { parent_id: '', status: 1, sort: 0 }
  parentAppName.value = ''
  dialogVisible.value = true
}

const handleAddChild = (row: any) => {
  dialogTitle.value = '新增子级应用'
  formData.value = { parent_id: row.id, status: 1, sort: 0 }
  parentAppName.value = row.app_name
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑应用'
  formData.value = { ...row }
  const parent = allData.value.find(item => item.id === row.parent_id)
  parentAppName.value = parent ? parent.app_name : ''
  dialogVisible.value = true
}

const handleDelete = (row: any) => {
  const hasChildren = allData.value.some(item => item.parent_id === row.id)
  const msg = hasChildren
    ? `该应用下存在子级应用，删除将同时删除所有子级应用。确定要删除应用 "${row.app_name}" 吗？`
    : `确定要删除应用 "${row.app_name}" 吗？`
  
  ElMessageBox.confirm(msg, '提示', { type: 'warning' })
    .then(async () => {
      try {
        await deleteApp(row.id)
        ElMessage.success('删除成功')
        fetchData()
      } catch (error: any) { ElMessage.error(error.message || '删除失败') }
    }).catch(() => {})
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        if (formData.value.id) {
          await updateApp(formData.value.id, formData.value)
          ElMessage.success('更新成功')
        } else {
          await createApp(formData.value)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        fetchData()
      } catch (error: any) { ElMessage.error(error.message || '操作失败') }
      finally { submitLoading.value = false }
    }
  })
}

onMounted(() => { fetchData() })
</script>

<style scoped>
.app-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}
.tree-icon {
  font-size: 16px;
  color: #409eff;
}
.file-icon {
  color: #909399;
}
.main-card { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
:deep(.el-card__body) { height: 100%; display: flex; flex-direction: column; padding: 20px; overflow: hidden; box-sizing: border-box; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; flex-shrink: 0; }
.page-title { display: flex; align-items: center; gap: 10px; font-size: 20px; font-weight: 600; }
.title-icon { font-size: 24px; color: var(--el-color-primary); }
.header-actions { display: flex; gap: 10px; }
.search-area { flex-shrink: 0; margin-bottom: 20px; }
.table-area { flex: 1; overflow: auto; }
</style>
