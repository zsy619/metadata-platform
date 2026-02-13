<template>
  <div class="tenant-list container-padding">
    <div class="page-header">
      <h1 class="text-primary page-title">
        <el-icon class="title-icon">
          <OfficeBuilding />
        </el-icon>
        租户管理
      </h1>
      <div class="header-actions">
        <el-button type="danger" :icon="Delete" @click="handleBatchDelete" :disabled="selectedRows.length === 0">
          批量删除
        </el-button>
        <el-button type="primary" @click="handleCreate" :icon="Plus">
          新增租户
        </el-button>
      </div>
    </div>
    <el-card class="main-card">
      <div class="search-area">
        <el-input v-model="searchQuery" placeholder="请输入租户名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
        <el-select v-model="filterStatus" placeholder="筛选状态" style="width: 150px; margin-left: 10px" clearable @change="handleSearch">
          <el-option label="全部" value="" />
          <el-option label="有效" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
        <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">
          搜索
        </el-button>
        <el-button @click="handleReset" :icon="RefreshLeft">
          重置
        </el-button>
      </div>
      <div class="table-area">
        <el-table v-loading="loading" :element-loading-text="loadingText" :data="filteredData" border stripe style="width: 100%; height: 100%;" @selection-change="handleSelectionChange">
          <template #empty>
            <el-empty :description="searchQuery ? '未搜索到相关租户' : '暂无租户'">
              <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增租户</el-button>
            </el-empty>
          </template>
          <el-table-column type="selection" width="55" />
          <el-table-column prop="tenant_name" label="租户名称" width="180" show-overflow-tooltip />
          <el-table-column prop="tenant_code" label="租户编码" width="150" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag v-if="scope.row.status === 1" type="success">有效</el-tag>
              <el-tag v-else type="danger">禁用</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="create_at" label="创建时间" width="170">
            <template #default="scope">
              {{ formatDateTime(scope.row.create_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="scope">
              <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>
                编辑
              </el-button>
              <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg>
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px" label-position="right">
        <el-form-item label="租户名称" prop="tenant_name">
          <el-input v-model="formData.tenant_name" placeholder="请输入租户名称" />
        </el-form-item>
        <el-form-item label="租户编码" prop="tenant_code">
          <el-input v-model="formData.tenant_code" placeholder="请输入租户编码" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
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
import { Delete, Edit, OfficeBuilding, Plus, RefreshLeft, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import { getTenants, createTenant, updateTenant, deleteTenant } from '@/api/user'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')
const filterStatus = ref<number | ''>('')
const selectedRows = ref<any[]>([])

const allData = ref<any[]>([])

const filteredData = computed(() => {
  let data = allData.value
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    data = data.filter(item =>
      (item.tenant_name || '').toLowerCase().includes(query) ||
      (item.tenant_code || '').toLowerCase().includes(query)
    )
  }
  if (filterStatus.value !== '') {
    data = data.filter(item => item.status === filterStatus.value)
  }
  return data
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()
const formData = ref<any>({})
const submitLoading = ref(false)

const formRules: FormRules = {
  tenant_name: [{ required: true, message: '请输入租户名称', trigger: 'blur' }],
  tenant_code: [{ required: true, message: '请输入租户编码', trigger: 'blur' }]
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return isNaN(date.getTime()) ? '-' : date.toLocaleString('zh-CN')
}

const fetchData = async () => {
  loadingText.value = '加载中...'
  loading.value = true
  try {
    const res: any = await getTenants()
    allData.value = res.data || res
  } catch (error) {
    console.error('加载租户列表失败:', error)
    ElMessage.error('加载列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  // Computed property handles filtering
}

const handleDebouncedSearch = () => {
  // Debounced search handled by @input
}

const handleReset = () => {
  searchQuery.value = ''
  filterStatus.value = ''
}

const handleSelectionChange = (val: any[]) => {
  selectedRows.value = val
}

const handleCreate = () => {
  dialogTitle.value = '新增租户'
  formData.value = { status: 1 }
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑租户'
  formData.value = { ...row }
  dialogVisible.value = true
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(`确定要删除租户 "${row.tenant_name}" 吗？`, '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await deleteTenant(row.id)
      ElMessage.success('删除成功')
      fetchData()
    } catch (error: any) {
      ElMessage.error(error.message || '删除失败')
    }
  }).catch(() => {})
}

const handleBatchDelete = () => {
  ElMessageBox.confirm(`确定要删除选中的 ${selectedRows.value.length} 个租户吗？`, '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await Promise.all(selectedRows.value.map(row => deleteTenant(row.id)))
      ElMessage.success('批量删除成功')
      fetchData()
    } catch (error: any) {
      ElMessage.error(error.message || '删除失败')
    }
  }).catch(() => {})
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        if (formData.value.id) {
          await updateTenant(formData.value.id, formData.value)
          ElMessage.success('更新成功')
        } else {
          await createTenant(formData.value)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        fetchData()
      } catch (error: any) {
        ElMessage.error(error.message || '操作失败')
      } finally {
        submitLoading.value = false
      }
    }
  })
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.tenant-list {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.main-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

:deep(.el-card__body) {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 20px;
  overflow: hidden;
  box-sizing: border-box;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-shrink: 0;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  font-weight: 600;
}

.title-icon {
  font-size: 24px;
  color: var(--el-color-primary);
}

.header-actions {
  display: flex;
  gap: 10px;
}

.search-area {
  flex-shrink: 0;
  margin-bottom: 20px;
}

.table-area {
  flex: 1;
  overflow: hidden;
}
</style>
