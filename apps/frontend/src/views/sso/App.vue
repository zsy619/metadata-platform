<template>
  <div class="sso-page container-padding">
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
        <el-input v-model="searchQuery" placeholder="请输入应用名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
        <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
        <el-button @click="handleReset" :icon="RefreshLeft">重置</el-button>
      </div>
      <div class="table-area">
        <el-table v-loading="loading" :data="displayData" border stripe>
          <el-table-column prop="application_name" label="应用名称" width="180" show-overflow-tooltip />
          <el-table-column prop="application_code" label="应用编码" width="150" />
          <el-table-column prop="host" label="域名/IP" show-overflow-tooltip />
          <el-table-column prop="state" label="状态" width="80">
            <template #default="scope">
              <el-tag v-if="scope.row.state === 1" type="success">启用</el-tag>
              <el-tag v-else type="danger">禁用</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="sort" label="排序" width="80" />
          <el-table-column prop="create_at" label="创建时间" width="170">
            <template #default="scope">{{ formatDateTime(scope.row.create_at) }}</template>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="应用名称" prop="application_name">
          <el-input v-model="formData.application_name" placeholder="请输入应用名称" />
        </el-form-item>
        <el-form-item label="应用编码" prop="application_code">
          <el-input v-model="formData.application_code" placeholder="请输入应用编码" />
        </el-form-item>
        <el-form-item label="域名/IP" prop="host">
          <el-input v-model="formData.host" placeholder="请输入域名或IP地址" />
        </el-form-item>
        <el-form-item label="状态" prop="state">
          <el-switch v-model="formData.state" :active-value="1" :inactive-value="0" />
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
import { ref, computed, onMounted } from 'vue'
import { Plus, Search, RefreshLeft, Monitor } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getApps, createApp, updateApp, deleteApp } from '@/api/user'

const searchQuery = ref('')
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const allData = ref<any[]>([])
const displayData = computed(() => {
  let data = allData.value
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    data = data.filter(item => (item.application_name || '').toLowerCase().includes(query))
  }
  total.value = data.length
  return data.slice((currentPage.value - 1) * pageSize.value, currentPage.value * pageSize.value)
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()
const formData = ref<any>({})
const submitLoading = ref(false)
const formRules: FormRules = {
  application_name: [{ required: true, message: '请输入应用名称', trigger: 'blur' }],
  application_code: [{ required: true, message: '请输入应用编码', trigger: 'blur' }]
}

const formatDateTime = (dateStr: string) => dateStr ? new Date(dateStr).toLocaleString('zh-CN') : '-'

const loadData = async () => {
  loading.value = true
  try { allData.value = await getApps() } catch (error: any) { ElMessage.error(error.message) }
  finally { loading.value = false }
}

const handleSearch = () => { currentPage.value = 1 }
const handleDebouncedSearch = () => handleSearch()
const handleReset = () => { searchQuery.value = ''; handleSearch() }
const handleSizeChange = () => { currentPage.value = 1 }
const handleCurrentChange = () => {}

const handleCreate = () => { dialogTitle.value = '新增应用'; formData.value = { state: 1, sort: 0 }; dialogVisible.value = true }
const handleEdit = (row: any) => { dialogTitle.value = '编辑应用'; formData.value = { ...row }; dialogVisible.value = true }
const handleDelete = async (row: any) => {
  try { await ElMessageBox.confirm('确定要删除该应用吗？', '提示', { type: 'warning' }); await deleteApp(row.id); ElMessage.success('删除成功'); loadData() }
  catch (error: any) { if (error !== 'cancel') ElMessage.error(error.message) }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        formData.value.id ? await updateApp(formData.value.id, formData.value) : await createApp(formData.value)
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
