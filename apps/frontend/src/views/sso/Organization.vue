<template>
  <div class="sso-page container-padding">
    <div class="page-header">
      <h1 class="text-primary page-title">
        <el-icon class="title-icon"><School /></el-icon>
        组织管理
      </h1>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate" :icon="Plus">新增组织</el-button>
      </div>
    </div>
    <el-card class="main-card">
      <div class="search-area">
        <el-input v-model="searchQuery" placeholder="请输入组织名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
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
          <el-table-column prop="unit_name" label="组织名称" width="200" show-overflow-tooltip />
          <el-table-column prop="unit_short" label="简称" width="150" />
          <el-table-column prop="unit_code" label="组织编码" width="150" />
          <el-table-column prop="kind_code" label="类型编码" width="120" />
          <el-table-column prop="contact" label="联系人" width="120" />
          <el-table-column prop="phone" label="联系电话" width="150" />
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
        <div class="pagination-wrapper">
          <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :total="total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
        </div>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="550px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="上级组织" prop="parent_id">
          <el-tree-select v-model="formData.parent_id" :data="orgTreeData" check-strictly :render-after-expand="false" placeholder="请选择上级组织" style="width: 100%" />
        </el-form-item>
        <el-form-item label="组织名称" prop="unit_name">
          <el-input v-model="formData.unit_name" placeholder="请输入组织名称" />
        </el-form-item>
        <el-form-item label="组织简称" prop="unit_short">
          <el-input v-model="formData.unit_short" placeholder="请输入组织简称" />
        </el-form-item>
        <el-form-item label="组织编码" prop="unit_code">
          <el-input v-model="formData.unit_code" placeholder="请输入组织编码" />
        </el-form-item>
        <el-form-item label="类型编码" prop="kind_code">
          <el-input v-model="formData.kind_code" placeholder="请输入类型编码" />
        </el-form-item>
        <el-form-item label="联系人" prop="contact">
          <el-input v-model="formData.contact" placeholder="请输入联系人" />
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="formData.phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="联系地址" prop="address">
          <el-input v-model="formData.address" placeholder="请输入联系地址" />
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
import { Plus, Search, RefreshLeft, School } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getUnits, createUnit, updateUnit, deleteUnit } from '@/api/user'

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
    data = data.filter(item => (item.unit_name || '').toLowerCase().includes(query))
  }
  if (filterState.value !== '') data = data.filter(item => item.state === filterState.value)
  total.value = data.length
  return data.slice((currentPage.value - 1) * pageSize.value, currentPage.value * pageSize.value)
})

const orgTreeData = computed(() => {
  const buildTree = (items: any[], parentId = ''): any[] => {
    return items.filter(item => item.parent_id === parentId).map(item => ({
      value: item.id,
      label: item.unit_name,
      children: buildTree(items, item.id)
    }))
  }
  return buildTree(allData.value)
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const formRef = ref<FormInstance>()
const formData = ref<any>({})
const submitLoading = ref(false)
const formRules: FormRules = {
  unit_name: [{ required: true, message: '请输入组织名称', trigger: 'blur' }],
  unit_code: [{ required: true, message: '请输入组织编码', trigger: 'blur' }]
}

const loadData = async () => {
  loading.value = true
  try { allData.value = await getUnits() } catch (error: any) { ElMessage.error(error.message) }
  finally { loading.value = false }
}

const handleSearch = () => { currentPage.value = 1 }
const handleDebouncedSearch = () => handleSearch()
const handleReset = () => { searchQuery.value = ''; filterState.value = ''; handleSearch() }
const handleSizeChange = () => { currentPage.value = 1 }
const handleCurrentChange = () => {}

const handleCreate = () => {
  dialogTitle.value = '新增组织'
  formData.value = { state: 1, sort: 0, parent_id: '' }
  dialogVisible.value = true
}
const handleEdit = (row: any) => { dialogTitle.value = '编辑组织'; formData.value = { ...row }; dialogVisible.value = true }
const handleDelete = async (row: any) => {
  try { await ElMessageBox.confirm('确定要删除该组织吗？', '提示', { type: 'warning' }); await deleteUnit(row.id); ElMessage.success('删除成功'); loadData() }
  catch (error: any) { if (error !== 'cancel') ElMessage.error(error.message) }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        formData.value.id ? await updateUnit(formData.value.id, formData.value) : await createUnit(formData.value)
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
