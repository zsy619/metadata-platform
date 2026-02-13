<template>
  <div class="sso-page container-padding">
    <div class="page-header">
      <h1 class="text-primary page-title">
        <el-icon class="title-icon"><Avatar /></el-icon>
        用户管理
      </h1>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate" :icon="Plus">新增用户</el-button>
      </div>
    </div>
    <el-card class="main-card">
      <div class="search-area">
        <el-input v-model="searchQuery" placeholder="请输入用户名搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
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
            <el-empty :description="searchQuery ? '未搜索到相关用户' : '暂无用户'">
              <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增用户</el-button>
            </el-empty>
          </template>
          <el-table-column prop="name" label="姓名" width="120" />
          <el-table-column prop="account" label="账号" width="150" />
          <el-table-column prop="mobile" label="手机号" width="130" />
          <el-table-column prop="email" label="邮箱" width="180" show-overflow-tooltip />
          <el-table-column prop="status" label="状态" width="80">
            <template #default="scope">
              <el-tag v-if="scope.row.status === 1" type="success">有效</el-tag>
              <el-tag v-else type="danger">禁用</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="create_at" label="创建时间" width="170">
            <template #default="scope">{{ formatDateTime(scope.row.create_at) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="scope">
              <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>编辑</el-button>
              <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg>删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="550px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px" label-position="right">
        <el-form-item label="账号" prop="account">
          <el-input v-model="formData.account" placeholder="请输入账号" />
        </el-form-item>
        <el-form-item v-if="!formData.id" label="密码" prop="password">
          <el-input v-model="formData.password" type="password" show-password placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="formData.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="mobile">
          <el-input v-model="formData.mobile" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
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
import { Avatar, Delete, Edit, Plus, RefreshLeft, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import { getUsers, createUser, updateUser, deleteUser } from '@/api/user'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')
const filterStatus = ref<number | ''>('')

const allData = ref<any[]>([])

const filteredData = computed(() => {
  let data = allData.value
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    data = data.filter(item => (item.name || '').toLowerCase().includes(query) || (item.account || '').toLowerCase().includes(query))
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
  account: [{ required: true, message: '请输入账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return isNaN(date.getTime()) ? '-' : date.toLocaleString('zh-CN')
}

const loadData = async () => {
  loadingText.value = '加载中...'
  loading.value = true
  try {
    const res: any = await getUsers()
    allData.value = res.data || res
  } catch (error) {
    console.error('加载用户列表失败:', error)
    ElMessage.error('加载列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {}
const handleDebouncedSearch = () => {}
const handleReset = () => { searchQuery.value = ''; filterStatus.value = '' }

const handleCreate = () => {
  dialogTitle.value = '新增用户'
  formData.value = { status: 1 }
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑用户'
  formData.value = { ...row }
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除用户 "${row.name}" 吗？`, '提示', { type: 'warning' })
    await deleteUser(row.id)
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
        formData.value.id ? await updateUser(formData.value.id, formData.value) : await createUser(formData.value)
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
