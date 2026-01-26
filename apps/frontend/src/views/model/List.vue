<template>
  <div class="model-list">
    <div class="page-header">
      <h1>模型管理</h1>
      <el-button type="primary" @click="handleCreate" :icon="Plus">
        新增模型
      </el-button>
    </div>

    <el-card>
      <div class="card-header">
        <el-input
          v-model="searchQuery"
          placeholder="请输入模型名称或编码搜索"
          clearable
          prefix-icon="Search"
          style="width: 300px"
        />
        <el-select
          v-model="filterKind"
          placeholder="筛选模型类型"
          style="width: 180px; margin-left: 10px"
        >
          <el-option label="全部" value="" />
          <el-option label="SQL语句" :value="1" />
          <el-option label="视图/表" :value="2" />
          <el-option label="存储过程" :value="3" />
          <el-option label="关联" :value="4" />
        </el-select>
        <el-select
          v-model="filterConn"
          placeholder="筛选数据源"
          style="width: 180px; margin-left: 10px"
        >
          <el-option label="全部" value="" />
          <el-option
            v-for="conn in dataSources"
            :key="conn.id"
            :label="conn.connName"
            :value="conn.id"
          />
        </el-select>
        <el-button type="primary" @click="handleSearch" :icon="Search">
          搜索
        </el-button>
      </div>

      <el-table
        v-loading="loading"
        :data="models"
        border
        style="width: 100%"
      >
        <el-table-column prop="modelName" label="模型名称" width="200" />
        <el-table-column prop="modelCode" label="模型编码" width="180" />
        <el-table-column prop="modelVersion" label="版本" width="100" />
        <el-table-column prop="modelKind" label="模型类型" width="120">
          <template #default="scope">
            <el-tag
              :type="getModelKindTagType(scope.row.modelKind)"
            >
              {{ getModelKindText(scope.row.modelKind) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="connName" label="数据源" width="150" />
        <el-table-column prop="isPublic" label="是否公开" width="100">
          <template #default="scope">
            <el-switch
              v-model="scope.row.isPublic"
              disabled
              active-color="#13ce66"
              inactive-color="#ff4d4f"
            />
          </template>
        </el-table-column>
        <el-table-column prop="state" label="状态" width="100">
          <template #default="scope">
            <el-tag
              :type="scope.row.state === 1 ? 'success' : 'danger'"
            >
              {{ scope.row.state === 1 ? '有效' : '无效' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              :icon="View"
              @click="handlePreview(scope.row)"
            >
              预览
            </el-button>
            <el-button
              type="success"
              size="small"
              :icon="Edit"
              @click="handleEdit(scope.row)"
              :disabled="scope.row.isLocked"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              size="small"
              :icon="Delete"
              @click="handleDelete(scope.row)"
              :disabled="scope.row.isLocked"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  Plus, 
  Search, 
  View, 
  Edit, 
  Delete 
} from '@element-plus/icons-vue'
import type { Model } from '@/types/model'
import type { DataSource } from '@/types/data-source'
import { getModels } from '@/api/model'
import { getDataSources } from '@/api/data-source'

const router = useRouter()

// 响应式数据
const loading = ref(false)
const searchQuery = ref('')
const filterKind = ref('')
const filterConn = ref('')
const models = ref<Model[]>([])
const dataSources = ref<DataSource[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 生命周期钩子
onMounted(() => {
  fetchDataSources()
  fetchModels()
})

// 获取数据源列表（用于筛选）
const fetchDataSources = async () => {
  try {
    // 模拟API调用，实际需要替换为真实API
    // const response = await getDataSources()
    // dataSources.value = response.data
    
    // 模拟数据
    dataSources.value = [
      { id: 1, connName: '测试MySQL数据源', connKind: 'MySQL', connVersion: '8.0', connHost: 'localhost', connPort: 3306, connUser: 'root', connPassword: 'password', connDatabase: 'test_db', connConn: '', isDeleted: false, state: 1, remark: '', sort: 0, createdAt: '', updatedAt: '' }
    ]
  } catch (error) {
    console.error('获取数据源列表失败:', error)
    ElMessage.error('获取数据源列表失败')
  }
}

// 获取模型列表
const fetchModels = async () => {
  loading.value = true
  try {
    // 模拟API调用，实际需要替换为真实API
    // const response = await getModels({
    //   page: currentPage.value,
    //   pageSize: pageSize.value,
    //   search: searchQuery.value,
    //   kind: filterKind.value ? Number(filterKind.value) : undefined,
    //   connID: filterConn.value ? Number(filterConn.value) : undefined
    // })
    // models.value = response.data
    // total.value = response.total
    
    // 模拟数据
    models.value = [
      {
        modelID: 1,
        parentID: 0,
        connID: 1,
        connName: '测试MySQL数据源',
        modelName: '用户模型',
        modelCode: 'user_model',
        modelVersion: '1.0.0',
        modelLogo: '',
        modelKind: 2,
        isPublic: true,
        isLocked: false,
        isDeleted: false,
        state: 1,
        remark: '用户信息模型',
        sort: 0,
        createID: 1,
        createBy: 'admin',
        createAt: '2024-01-23 10:00:00',
        updateID: 1,
        updateBy: 'admin',
        updateAt: '2024-01-23 10:00:00'
      },
      {
        modelID: 2,
        parentID: 0,
        connID: 1,
        connName: '测试MySQL数据源',
        modelName: '订单模型',
        modelCode: 'order_model',
        modelVersion: '1.0.0',
        modelLogo: '',
        modelKind: 2,
        isPublic: false,
        isLocked: false,
        isDeleted: false,
        state: 1,
        remark: '订单信息模型',
        sort: 0,
        createID: 1,
        createBy: 'admin',
        createAt: '2024-01-23 11:00:00',
        updateID: 1,
        updateBy: 'admin',
        updateAt: '2024-01-23 11:00:00'
      }
    ]
    total.value = models.value.length
  } catch (error) {
    console.error('获取模型列表失败:', error)
    ElMessage.error('获取模型列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  fetchModels()
}

// 页码变化
const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchModels()
}

// 每页条数变化
const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  fetchModels()
}

// 新增模型
const handleCreate = () => {
  router.push('/models/create')
}

// 编辑模型
const handleEdit = (row: Model) => {
  router.push(`/models/${row.modelID}/edit`)
}

// 删除模型
const handleDelete = (row: Model) => {
  ElMessageBox.confirm(
    `确定要删除模型 "${row.modelName}" 吗？`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      // 模拟删除操作
      // await deleteModel(row.modelID)
      ElMessage.success('删除成功')
      fetchModels()
    } catch (error) {
      console.error('删除模型失败:', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {
    // 取消删除
  })
}

// 预览模型
const handlePreview = (row: Model) => {
  router.push(`/models/${row.modelID}/preview`)
}

// 获取模型类型文本
const getModelKindText = (kind: number): string => {
  const kindMap: Record<number, string> = {
    1: 'SQL语句',
    2: '视图/表',
    3: '存储过程',
    4: '关联'
  }
  return kindMap[kind] || '未知'
}

// 获取模型类型标签样式
const getModelKindTagType = (kind: number): string => {
  const typeMap: Record<number, string> = {
    1: 'info',
    2: 'success',
    3: 'warning',
    4: 'danger'
  }
  return typeMap[kind] || 'info'
}
</script>

<style scoped>
.model-list {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}
</style>
