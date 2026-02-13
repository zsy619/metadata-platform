<template>
  <div class="sso-page container-padding">
    <div class="page-header">
      <h1 class="text-primary page-title">
        <el-icon class="title-icon">
          <OfficeBuilding />
        </el-icon>
        租户管理
      </h1>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate" :icon="Plus">
          新增租户
        </el-button>
      </div>
    </div>
    <el-card class="main-card">
      <div class="search-area">
        <el-input v-model="searchQuery" placeholder="请输入租户名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
        <el-select v-model="filterState" placeholder="筛选状态" style="width: 150px; margin-left: 10px" clearable @change="handleSearch">
          <el-option label="全部" value="" />
          <el-option label="有效" :value="1" />
          <el-option label="禁用" :value="0" />
        </el-select>
        <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
        <el-button @click="handleReset" :icon="RefreshLeft">重置</el-button>
      </div>
      <div class="table-area">
        <el-table v-loading="loading" :data="displayData" border stripe style="width: 100%">
          <el-table-column prop="tenant_name" label="租户名称" width="180" show-overflow-tooltip />
          <el-table-column prop="tenant_code" label="租户编码" width="150" />
          <el-table-column prop="linkman" label="联系人" width="120" />
          <el-table-column prop="contact" label="联系电话" width="150" />
          <el-table-column prop="address" label="联系地址" show-overflow-tooltip />
          <el-table-column prop="state" label="状态" width="80">
            <template #default="scope">
              <el-tag v-if="scope.row.state === 1" type="success">有效</el-tag>
              <el-tag v-else type="danger">禁用</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="create_at" label="创建时间" width="170">
            <template #default="scope">
              {{ formatDateTime(scope.row.create_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="scope">
              <el-button type="primary" link @click="handleEdit(scope.row)">编辑</el-button>
              <el-button type="danger" link @click="handleDelete(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="pagination-wrapper">
          <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]" :total="total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
        </div>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="租户名称" prop="tenant_name">
          <el-input v-model="formData.tenant_name" placeholder="请输入租户名称" />
        </el-form-item>
        <el-form-item label="租户编码" prop="tenant_code">
          <el-input v-model="formData.tenant_code" placeholder="请输入租户编码" />
        </el-form-item>
        <el-form-item label="联系人" prop="linkman">
          <el-input v-model="formData.linkman" placeholder="请输入联系人" />
        </el-form-item>
        <el-form-item label="联系电话" prop="contact">
          <el-input v-model="formData.contact" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="联系地址" prop="address">
          <el-input v-model="formData.address" placeholder="请输入联系地址" />
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
import { Plus, Search, RefreshLeft, OfficeBuilding } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getTenants, createTenant, updateTenant, deleteTenant } from '@/api/user'

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
    data = data.filter(item =>
      (item.tenant_name || '').toLowerCase().includes(query) ||
      (item.tenant_code || '').toLowerCase().includes(query)
    )
  }
  if (filterState.value !== '') {
    data = data.filter(item => item.state === filterState.value)
  }
  total.value = data.length
  const start = (currentPage.value - 1) * pageSize.value
  return data.slice(start, start + pageSize.value)
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
  return date.toLocaleString('zh-CN')
}

const loadData = async () => {
  loading.value = true
  try {
    allData.value = await getTenants()
    total.value = allData.value.length
  } catch (error: any) {
    ElMessage.error(error.message || '加载数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { currentPage.value = 1 }
const handleDebouncedSearch = () => { handleSearch() }
const handleReset = () => { searchQuery.value = ''; filterState.value = ''; handleSearch() }
const handleSizeChange = () => { currentPage.value = 1 }
const handleCurrentChange = () => {}

const handleCreate = () => {
  dialogTitle.value = '新增租户'
  formData.value = { state: 1 }
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = '编辑租户'
  formData.value = { ...row }
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除该租户吗？', '提示', { type: 'warning' })
    await deleteTenant(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') ElMessage.error(error.message || '删除失败')
  }
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
        loadData()
      } catch (error: any) {
        ElMessage.error(error.message || '操作失败')
      } finally {
        submitLoading.value = false
      }
    }
  })
}

onMounted(() => { loadData() })
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
