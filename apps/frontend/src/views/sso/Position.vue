<template>
  <div class="sso-page container-padding">
    <div class="page-header">
      <h1 class="text-primary page-title">
        <el-icon class="title-icon"><Briefcase /></el-icon>
        职位管理
      </h1>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate" :icon="Plus">新增职位</el-button>
      </div>
    </div>
    <el-card class="main-card">
      <div class="search-area">
        <el-input v-model="searchQuery" placeholder="请输入职位名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
        <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
      </div>
      <div class="table-area">
        <el-table v-loading="loading" :data="displayData" border stripe>
          <el-table-column prop="pos_name" label="职位名称" width="180" />
          <el-table-column prop="pos_code" label="职位编码" width="150" />
          <el-table-column prop="organization_id" label="所属组织" width="150" />
          <el-table-column prop="state" label="状态" width="80">
            <template #default="scope">
              <el-tag v-if="scope.row.state === 1" type="success">有效</el-tag>
              <el-tag v-else type="danger">禁用</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="sort" label="排序" width="80" />
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="scope">
              <el-button type="primary" link @click="handleEdit(scope.row)">编辑</el-button>
              <el-button type="danger" link @click="handleDelete(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="formData" label-width="100px">
        <el-form-item label="职位名称" prop="pos_name">
          <el-input v-model="formData.pos_name" placeholder="请输入职位名称" />
        </el-form-item>
        <el-form-item label="职位编码" prop="pos_code">
          <el-input v-model="formData.pos_code" placeholder="请输入职位编码" />
        </el-form-item>
        <el-form-item label="状态" prop="state">
          <el-switch v-model="formData.state" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" />
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
import { Plus, Search, Briefcase } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { getPos, createPos, updatePos, deletePos } from '@/api/user'

const searchQuery = ref('')
const loading = ref(false)
const allData = ref<any[]>([])
const displayData = computed(() => {
  if (!searchQuery.value) return allData.value
  const query = searchQuery.value.toLowerCase()
  return allData.value.filter(item => (item.pos_name || '').toLowerCase().includes(query))
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()
const formData = ref<any>({})
const submitLoading = ref(false)

const loadData = async () => {
  loading.value = true
  try { allData.value = await getPos() } catch (error: any) { ElMessage.error(error.message) }
  finally { loading.value = false }
}
const handleSearch = () => {}
const handleDebouncedSearch = () => {}
const handleCreate = () => { dialogTitle.value = '新增职位'; formData.value = { state: 1, sort: 0 }; dialogVisible.value = true }
const handleEdit = (row: any) => { dialogTitle.value = '编辑职位'; formData.value = { ...row }; dialogVisible.value = true }
const handleDelete = async (row: any) => {
  try { await ElMessageBox.confirm('确定要删除该职位吗？', '提示', { type: 'warning' }); await deletePos(row.id); ElMessage.success('删除成功'); loadData() }
  catch (error: any) { if (error !== 'cancel') ElMessage.error(error.message) }
}
const handleSubmit = async () => {
  submitLoading.value = true
  try { formData.value.id ? await updatePos(formData.value.id, formData.value) : await createPos(formData.value); ElMessage.success(formData.value.id ? '更新成功' : '创建成功'); dialogVisible.value = false; loadData() }
  catch (error: any) { ElMessage.error(error.message) }
  finally { submitLoading.value = false }
}
onMounted(() => loadData())
</script>

<style scoped>
.sso-page { height: 100%; display: flex; flex-direction: column; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.page-title { font-size: 20px; font-weight: 600; display: flex; align-items: center; gap: 8px; }
.title-icon { font-size: 24px; }
.main-card { flex: 1; display: flex; flex-direction: column; }
.search-area { display: flex; margin-bottom: 20px; }
.table-area { flex: 1; }
.text-primary { color: var(--el-text-color-primary); }
</style>
