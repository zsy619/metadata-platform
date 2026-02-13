<template>
  <div class="sso-page container-padding">
    <div class="page-header">
      <h1 class="text-primary page-title">
        <el-icon class="title-icon"><Collection /></el-icon>
        组织类型管理
      </h1>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate" :icon="Plus">新增类型</el-button>
      </div>
    </div>
    <el-card class="main-card">
      <div class="search-area">
        <el-input v-model="searchQuery" placeholder="请输入类型名称搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
        <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
      </div>
      <div class="table-area">
        <el-table v-loading="loading" :data="displayData" border stripe row-key="id" :tree-props="{ children: 'children', hasChildren: 'hasChildren' }">
          <el-table-column prop="kind_name" label="类型名称" width="200" />
          <el-table-column prop="kind_code" label="类型编码" width="150" />
          <el-table-column prop="kind_tag" label="等级标识" width="120" />
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
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="上级类型" prop="parent_id">
          <el-tree-select v-model="formData.parent_id" :data="treeData" check-strictly :render-after-expand="false" placeholder="请选择上级类型" style="width: 100%" />
        </el-form-item>
        <el-form-item label="类型名称" prop="kind_name">
          <el-input v-model="formData.kind_name" placeholder="请输入类型名称" />
        </el-form-item>
        <el-form-item label="类型编码" prop="kind_code">
          <el-input v-model="formData.kind_code" placeholder="请输入类型编码" />
        </el-form-item>
        <el-form-item label="等级标识" prop="kind_tag">
          <el-input v-model="formData.kind_tag" placeholder="请输入等级标识" />
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
import { Plus, Search, Collection } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { getOrgKinds, createOrgKind, updateOrgKind, deleteOrgKind } from '@/api/user'

const searchQuery = ref('')
const loading = ref(false)
const allData = ref<any[]>([])
const displayData = computed(() => allData.value)
const treeData = computed(() => {
  const buildTree = (items: any[], parentId = ''): any[] => {
    return items.filter(item => item.parent_id === parentId).map(item => ({
      value: item.id,
      label: item.kind_name,
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
  kind_name: [{ required: true, message: '请输入类型名称', trigger: 'blur' }],
  kind_code: [{ required: true, message: '请输入类型编码', trigger: 'blur' }]
}

const loadData = async () => {
  loading.value = true
  try { allData.value = await getOrgKinds() } catch (error: any) { ElMessage.error(error.message) }
  finally { loading.value = false }
}

const handleSearch = () => {}
const handleDebouncedSearch = () => {}
const handleCreate = () => { dialogTitle.value = '新增类型'; formData.value = { state: 1, sort: 0, parent_id: '' }; dialogVisible.value = true }
const handleEdit = (_row: any) => { dialogTitle.value = '编辑类型'; formData.value = { ..._row }; dialogVisible.value = true }
const handleDelete = async (_row: any) => {
  try { await ElMessageBox.confirm('确定要删除该类型吗？', '提示', { type: 'warning' }); await deleteOrgKind(_row.id); ElMessage.success('删除成功'); loadData() }
  catch (error: any) { if (error !== 'cancel') ElMessage.error(error.message) }
}
const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        formData.value.id ? await updateOrgKind(formData.value.id, formData.value) : await createOrgKind(formData.value)
        ElMessage.success(formData.value.id ? '更新成功' : '创建成功')
        dialogVisible.value = false; loadData()
      }
      catch (error: any) { ElMessage.error(error.message) }
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
.main-card { flex: 1; display: flex; flex-direction: column; }
.search-area { display: flex; margin-bottom: 20px; }
.table-area { flex: 1; }
.text-primary { color: var(--el-text-color-primary); }
</style>
