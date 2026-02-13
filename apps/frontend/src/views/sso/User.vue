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
        <el-input v-model="searchQuery" placeholder="请输入用户名或账号搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
        <el-select v-model="filterState" placeholder="筛选状态" style="width: 150px; margin-left: 10px" clearable @change="handleSearch">
          <el-option label="全部" value="" />
          <el-option label="有效" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
        <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
        <el-button @click="handleReset" :icon="RefreshLeft">重置</el-button>
      </div>
      <div class="table-area">
        <el-table v-loading="loading" :data="displayData" border stripe>
          <el-table-column prop="account" label="账号" width="150" />
          <el-table-column prop="name" label="姓名" width="120" />
          <el-table-column prop="mobile" label="手机号" width="140" />
          <el-table-column prop="email" label="邮箱" show-overflow-tooltip />
          <el-table-column prop="sex" label="性别" width="80" />
          <el-table-column prop="state" label="状态" width="80">
            <template #default="scope">
              <el-tag v-if="scope.row.state === 1" type="success">有效</el-tag>
              <el-tag v-else type="danger">禁用</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="kind" label="类型" width="100">
            <template #default="scope">
              <el-tag v-if="scope.row.kind === 1" type="danger">超级管理员</el-tag>
              <el-tag v-else-if="scope.row.kind === 2" type="warning">子管理员</el-tag>
              <el-tag v-else type="info">其他</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="last_login_time" label="最后登录" width="170">
            <template #default="scope">{{ scope.row.last_login_time ? formatDateTime(scope.row.last_login_time) : '-' }}</template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="scope">
              <el-button type="primary" link @click="handleEdit(scope.row)">编辑</el-button>
              <el-button type="danger" link @click="handleDelete(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="pagination-wrapper">
          <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
        </div>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="550px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="账号" prop="account">
          <el-input v-model="formData.account" placeholder="请输入账号" />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="formData.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="formData.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="手机号" prop="mobile">
          <el-input v-model="formData.mobile" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="性别" prop="sex">
          <el-radio-group v-model="formData.sex">
            <el-radio label="男">男</el-radio>
            <el-radio label="女">女</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="类型" prop="kind">
          <el-select v-model="formData.kind" style="width: 100%">
            <el-option label="超级管理员" :value="1" />
            <el-option label="子管理员" :value="2" />
            <el-option label="其他" :value="99" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="state">
          <el-switch v-model="formData.state" :active-value="1" :inactive-value="0" />
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
import { ref, computed, onMounted } from 'vue'
import { Plus, Search, RefreshLeft, Avatar } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getUsers, createUser, updateUser, deleteUser } from '@/api/user'

const searchQuery = ref('')
const filterState = ref<number | ''>('')
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const allData = ref<any[]>([])

const displayData = computed(() => {
  let data = allData.value
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    data = data.filter(item => (item.account || '').toLowerCase().includes(query) || (item.name || '').toLowerCase().includes(query))
  }
  if (filterState.value !== '') data = data.filter(item => item.state === filterState.value)
  total.value = data.length
  return data.slice((currentPage.value - 1) * pageSize.value, currentPage.value * pageSize.value)
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()
const formData = ref<any>({})
const submitLoading = ref(false)
const formRules: FormRules = {
  account: [{ required: true, message: '请输入账号', trigger: 'blur' }],
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }]
}

const formatDateTime = (dateStr: string) => dateStr ? new Date(dateStr).toLocaleString('zh-CN') : '-'

const loadData = async () => {
  loading.value = true
  try { allData.value = await getUsers() } catch (error: any) { ElMessage.error(error.message) }
  finally { loading.value = false }
}

const handleSearch = () => { currentPage.value = 1 }
const handleDebouncedSearch = () => handleSearch()
const handleReset = () => { searchQuery.value = ''; filterState.value = ''; handleSearch() }
const handleSizeChange = () => { currentPage.value = 1 }
const handleCurrentChange = () => {}

const handleCreate = () => { dialogTitle.value = '新增用户'; formData.value = { state: 1, kind: 99, sex: '男' }; dialogVisible.value = true }
const handleEdit = (row: any) => { dialogTitle.value = '编辑用户'; formData.value = { ...row }; delete formData.value.password; dialogVisible.value = true }
const handleDelete = async (row: any) => {
  try { await ElMessageBox.confirm('确定要删除该用户吗？', '提示', { type: 'warning' }); await deleteUser(row.id); ElMessage.success('删除成功'); loadData() }
  catch (error: any) { if (error !== 'cancel') ElMessage.error(error.message) }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        formData.value.id ? await updateUser(formData.value.id, formData.value) : await createUser(formData.value)
        ElMessage.success(formData.value.id ? '更新成功' : '创建成功')
        dialogVisible.value = false; loadData()
      } catch (error: any) { ElMessage.error(error.message) }
      finally { submitLoading.value = false }
    }
  })
}

onMounted(() => loadData())
</script>

<style scoped>
.sso-page { height: 100%; display: flex; flex-direction: column; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.page-title { font-size: 20px; font-weight: 600; display: flex; align-items: center; gap: 8px; }
.title-icon { font-size: 24px; }
.header-actions { display: flex; gap: 10px; }
.main-card { flex: 1; display: flex; flex-direction: column; }
.search-area { display: flex; align-items: center; margin-bottom: 20px; flex-wrap: wrap; gap: 10px; }
.table-area { flex: 1; }
.pagination-wrapper { display: flex; justify-content: flex-end; margin-top: 20px; }
.text-primary { color: var(--el-text-color-primary); }
</style>
