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
        <el-table v-loading="loading" :element-loading-text="loadingText" :data="filteredData" border stripe style="width: 100%; height: 100%;">
          <template #empty>
            <el-empty :description="searchQuery ? '未搜索到相关角色' : '暂无角色'">
              <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增角色</el-button>
            </el-empty>
          </template>
          <el-table-column prop="role_name" label="角色名称" width="180" show-overflow-tooltip />
          <el-table-column prop="role_code" label="角色编码" width="150" />
          <el-table-column prop="data_scope" label="数据范围" width="100">
            <template #default="scope">
              <el-tag v-if="scope.row.data_scope === '1'" type="success">全部</el-tag>
              <el-tag v-else-if="scope.row.data_scope === '2'" type="warning">自定义</el-tag>
              <el-tag v-else-if="scope.row.data_scope === '3'" type="info">本部门</el-tag>
              <el-tag v-else type="info">本部门及以下</el-tag>
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
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="scope">
              <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>编辑</el-button>
              <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg>删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px" label-position="right">
        <el-form-item label="角色名称" prop="role_name">
          <el-input v-model="formData.role_name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色编码" prop="role_code">
          <el-input v-model="formData.role_code" placeholder="请输入角色编码" />
        </el-form-item>
        <el-form-item label="数据范围" prop="data_scope">
          <el-select v-model="formData.data_scope" style="width: 100%">
            <el-option label="全部数据权限" value="1" />
            <el-option label="自定数据权限" value="2" />
            <el-option label="本部门数据权限" value="3" />
            <el-option label="本部门及以下" value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="状&#12288;&#12288;态" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="序&#12288;&#12288;号" prop="sort">
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
import { createRole, deleteRole, getRoles, updateRole } from '@/api/user'
import { Delete, Edit, Plus, RefreshLeft, Search, UserFilled } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')
const filterStatus = ref<number | ''>('')

const allData = ref<any[]>([])

const filteredData = computed(() => {
  let data = allData.value
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    data = data.filter(item => (item.role_name || '').toLowerCase().includes(query))
  }
  if (filterStatus.value !== '') data = data.filter(item => item.status === filterStatus.value)
  return data
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()
const formData = ref<any>({})
const submitLoading = ref(false)
const formRules: FormRules = {
  role_name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  role_code: [{ required: true, message: '请输入角色编码', trigger: 'blur' }]
}

const loadData = async () => {
  loadingText.value = '加载中...'
  loading.value = true
  try {
    const res: any = await getRoles()
    allData.value = res.data || res
  } catch (error) {
    console.error('加载角色列表失败:', error)
    ElMessage.error('加载列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {}
const handleDebouncedSearch = () => {}
const handleReset = () => { searchQuery.value = ''; filterStatus.value = '' }

const handleCreate = () => {
  dialogTitle.value = '新增角色'
  formData.value = { status: 1, sort: 0, data_scope: '1' }
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑角色'
  formData.value = { ...row }
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除角色 "${row.role_name}" 吗？`, '提示', { type: 'warning' })
    await deleteRole(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error: any) { if (error !== 'cancel') ElMessage.error(error.message || '删除失败') }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        formData.value.id ? await updateRole(formData.value.id, formData.value) : await createRole(formData.value)
        ElMessage.success(formData.value.id ? '更新成功' : '创建成功')
        dialogVisible.value = false
        loadData()
      } catch (error: any) { ElMessage.error(error.message || '操作失败') }
      finally { submitLoading.value = false }
    }
  })
}

onMounted(() => loadData())
</script>

<style scoped>
.sso-page { height: 100%; display: flex; flex-direction: column; }
.main-card { flex: 1; display: flex; flex-direction: column; overflow: hidden; }
:deep(.el-card__body) { height: 100%; display: flex; flex-direction: column; padding: 20px; overflow: hidden; box-sizing: border-box; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; flex-shrink: 0; }
.page-title { font-size: 20px; font-weight: 600; display: flex; align-items: center; gap: 8px; }
.title-icon { font-size: 24px; color: var(--el-color-primary); }
.header-actions { display: flex; gap: 10px; }
.search-area { display: flex; align-items: center; margin-bottom: 20px; flex-wrap: wrap; gap: 10px; }
.table-area { flex: 1; overflow: hidden; }
.text-primary { color: var(--el-text-color-primary); }
</style>
