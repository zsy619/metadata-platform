<template>
  <div class="sso-page container-padding">
    <div class="page-header">
      <h1 class="text-primary page-title">
        <el-icon class="title-icon">
          <component :is="icon" />
        </el-icon>
        {{ title }}
      </h1>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate" :icon="Plus">
          新增{{ title.replace('管理', '') }}
        </el-button>
      </div>
    </div>
    <el-card class="main-card">
      <div class="search-area">
        <el-input
          v-model="searchQuery"
          :placeholder="`请输入${searchPlaceholder}搜索`"
          clearable
          :prefix-icon="Search"
          style="width: 300px"
          @input="handleDebouncedSearch"
        />
        <el-select
          v-model="filterStatus"
          placeholder="筛选状态"
          style="width: 150px; margin-left: 10px"
          clearable
          @change="handleSearch"
        >
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
        <el-table
          v-loading="loading"
          :element-loading-text="loadingText"
          :data="displayData"
          border
          stripe
          style="width: 100%"
          @selection-change="handleSelectionChange"
        >
          <template #empty>
            <el-empty :description="searchQuery ? `未搜索到相关${searchPlaceholder}` : '暂无数据'">
              <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增</el-button>
            </el-empty>
          </template>
          <el-table-column v-if="showSelection" type="selection" width="55" />
          <el-table-column
            v-for="col in columns"
            :key="col.prop"
            :prop="col.prop"
            :label="col.label"
            :width="col.width"
            :show-overflow-tooltip="col.showOverflowTooltip !== false"
          >
            <template v-if="col.slot" #default="scope">
              <slot :name="col.slot" :row="scope.row" />
            </template>
          </el-table-column>
          <el-table-column label="操作" :width="actionWidth" fixed="right">
            <template #default="scope">
              <el-button type="primary" link @click="handleEdit(scope.row)">
                编辑
              </el-button>
              <el-button type="danger" link @click="handleDelete(scope.row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="140px" label-position="right">
        <el-form-item
          v-for="field in formFields"
          :key="field.prop"
          :label="field.label"
          :prop="field.prop"
        >
          <el-input
            v-if="field.type === 'input'"
            v-model="formData[field.prop]"
            :placeholder="field.placeholder"
          />
          <el-select
            v-else-if="field.type === 'select'"
            v-model="formData[field.prop]"
            :placeholder="field.placeholder"
            style="width: 100%"
          >
            <el-option
              v-for="opt in field.options"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </el-select>
          <el-switch
            v-else-if="field.type === 'switch'"
            v-model="formData[field.prop]"
          />
          <el-input
            v-else-if="field.type === 'textarea'"
            v-model="formData[field.prop]"
            type="textarea"
            :rows="3"
            :placeholder="field.placeholder"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Plus, Search, RefreshLeft } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

interface Column {
  prop: string
  label: string
  width?: number
  showOverflowTooltip?: boolean
  slot?: string
}

interface FormField {
  prop: string
  label: string
  type: 'input' | 'select' | 'switch' | 'textarea'
  placeholder?: string
  options?: Array<{ label: string; value: any }>
}

const props = defineProps<{
  title: string
  icon: any
  searchPlaceholder: string
  columns: Column[]
  formFields: FormField[]
  formRules?: FormRules
  api: {
    list: () => Promise<any[]>
    create?: (data: any) => Promise<any>
    update?: (id: string, data: any) => Promise<any>
    delete?: (id: string) => Promise<void>
  }
  showSelection?: boolean
}>()

const emit = defineEmits(['refresh', 'create', 'edit', 'delete'])

const searchQuery = ref('')
const filterStatus = ref<number | ''>('')
const loading = ref(false)
const loadingText = ref('加载中...')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const selectedRows = ref<any[]>([])

const allData = ref<any[]>([])
const displayData = computed(() => {
  let data = allData.value
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    data = data.filter(item =>
      Object.values(item).some(val =>
        String(val).toLowerCase().includes(query)
      )
    )
  }
  if (filterStatus.value !== '') {
    data = data.filter(item => item.state === filterStatus.value)
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
const actionWidth = computed(() => props.showSelection ? 200 : 150)

const loadData = async () => {
  loading.value = true
  loadingText.value = '加载中...'
  try {
    allData.value = await props.api.list()
    total.value = allData.value.length
  } catch (error: any) {
    ElMessage.error(error.message || '加载数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
}

const handleDebouncedSearch = () => {
  handleSearch()
}

const handleReset = () => {
  searchQuery.value = ''
  filterStatus.value = ''
  handleSearch()
}

const handleSelectionChange = (rows: any[]) => {
  selectedRows.value = rows
}

const handleSizeChange = () => {
  currentPage.value = 1
}

const handleCurrentChange = () => {}

const handleCreate = () => {
  dialogTitle.value = `新增${props.title.replace('管理', '')}`
  formData.value = {}
  props.formFields.forEach(field => {
    formData.value[field.prop] = field.type === 'switch' ? true : ''
  })
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  dialogTitle.value = `编辑${props.title.replace('管理', '')}`
  formData.value = { ...row }
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除该${props.title.replace('管理', '')}吗？`,
      '提示',
      { type: 'warning' }
    )
    if (props.api.delete) {
      await props.api.delete(row.id)
      ElMessage.success('删除成功')
      loadData()
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        if (formData.value.id && props.api.update) {
          await props.api.update(formData.value.id, formData.value)
          ElMessage.success('更新成功')
        } else if (props.api.create) {
          await props.api.create(formData.value)
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

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.sso-page {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-icon {
  font-size: 24px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.main-card {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.search-area {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 10px;
}

.table-area {
  flex: 1;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.text-primary {
  color: var(--el-text-color-primary);
}
</style>
