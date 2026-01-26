<template>
  <div class="model-create">
    <div class="page-header">
      <h1>创建模型</h1>
      <el-button @click="handleCancel" :icon="ArrowLeft">
        返回列表
      </el-button>
    </div>

    <el-card>
      <el-steps :active="activeStep" space="300px">
        <el-step title="选择数据源" />
        <el-step title="选择表/视图" />
        <el-step title="配置模型信息" />
        <el-step title="完成" />
      </el-steps>

      <!-- 第一步：选择数据源 -->
      <div v-if="activeStep === 0" class="step-content">
        <h3>选择数据源</h3>
        <el-form-item label="数据源">
          <el-select
            v-model="modelForm.connID"
            placeholder="请选择数据源"
            style="width: 100%"
            @change="handleDataSourceChange"
          >
            <el-option
              v-for="conn in dataSources"
              :key="conn.id"
              :label="conn.connName"
              :value="conn.id"
            >
              <div class="option-content">
                <span>{{ conn.connName }}</span>
                <span class="option-subtitle">{{ conn.connKind }} - {{ conn.connHost }}:{{ conn.connPort }}</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
        <div class="step-actions">
          <el-button type="primary" @click="nextStep" :disabled="!modelForm.connID">
            下一步
          </el-button>
        </div>
      </div>

      <!-- 第二步：选择表/视图 -->
      <div v-if="activeStep === 1" class="step-content">
        <h3>选择表/视图</h3>
        <el-form-item label="对象类型">
          <el-radio-group v-model="objectType" @change="handleObjectTypeChange">
            <el-radio-button label="table">表</el-radio-button>
            <el-radio-button label="view">视图</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="选择{{ objectType === 'table' ? '表' : '视图' }}">
          <el-tree
            v-loading="loadingTables"
            :data="tables"
            :props="tableTreeProps"
            node-key="id"
            ref="tableTree"
            @node-click="handleTableClick"
            style="max-height: 400px; overflow-y: auto"
          >
            <template #default="{ node, data }">
              <span class="custom-tree-node">
                <span>{{ data.label }}</span>
                <span class="node-subtitle">{{ data.comment || '无描述' }}</span>
              </span>
            </template>
          </el-tree>
        </el-form-item>
        <div class="step-actions">
          <el-button @click="prevStep">上一步</el-button>
          <el-button type="primary" @click="nextStep" :disabled="!selectedTable">
            下一步
          </el-button>
        </div>
      </div>

      <!-- 第三步：配置模型信息 -->
      <div v-if="activeStep === 2" class="step-content">
        <h3>配置模型信息</h3>
        <el-form
          ref="modelFormRef"
          :model="modelForm"
          :rules="formRules"
          label-width="120px"
          class="create-form"
        >
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="模型名称" prop="modelName">
                <el-input
                  v-model="modelForm.modelName"
                  placeholder="请输入模型名称"
                  clearable
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="模型编码" prop="modelCode">
                <el-input
                  v-model="modelForm.modelCode"
                  placeholder="请输入模型编码（字母数字下划线组合）"
                  clearable
                />
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="模型版本" prop="modelVersion">
                <el-input
                  v-model="modelForm.modelVersion"
                  placeholder="请输入模型版本，如：1.0.0"
                  clearable
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="模型类型" prop="modelKind">
                <el-select
                  v-model="modelForm.modelKind"
                  placeholder="请选择模型类型"
                >
                  <el-option label="SQL语句" :value="1" />
                  <el-option label="视图/表" :value="2" />
                  <el-option label="存储过程" :value="3" />
                  <el-option label="关联" :value="4" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-form-item label="是否公开">
            <el-switch v-model="modelForm.isPublic" />
          </el-form-item>
          
          <el-form-item label="备注">
            <el-input
              v-model="modelForm.remark"
              type="textarea"
              placeholder="请输入模型描述信息"
              :rows="4"
              clearable
            />
          </el-form-item>
          
          <div class="step-actions">
            <el-button @click="prevStep">上一步</el-button>
            <el-button type="primary" @click="handleSubmit">
              提交
            </el-button>
          </div>
        </el-form>
      </div>

      <!-- 第四步：完成 -->
      <div v-if="activeStep === 3" class="step-content success-step">
        <el-result
          icon="success"
          title="模型创建成功"
          sub-title="您已成功创建了一个新的业务数据模型"
        >
          <template #extra>
            <el-button type="primary" @click="handleFinish">
              查看模型列表
            </el-button>
            <el-button @click="handleCreateAnother">
              创建另一个模型
            </el-button>
          </template>
        </el-result>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import type { DataSource } from '@/types/data-source'
import type { Model } from '@/types/model'
import { getDataSources } from '@/api/data-source'
import { createModel, buildModelFromTable, buildModelFromView } from '@/api/model'

const router = useRouter()
const modelFormRef = ref()
const tableTree = ref()

// 响应式数据
const activeStep = ref(0)
const loading = ref(false)
const loadingTables = ref(false)
const objectType = ref('table')
const modelForm = reactive<Partial<Model & { connName?: string }>>({
  connID: undefined,
  connName: '',
  modelName: '',
  modelCode: '',
  modelVersion: '1.0.0',
  modelKind: 2,
  isPublic: false,
  remark: ''
})

const dataSources = ref<DataSource[]>([])
const tables = ref<any[]>([])
const selectedTable = ref<any>(null)

// 表树配置
const tableTreeProps = {
  label: 'label',
  children: 'children',
  isLeaf: 'isLeaf'
}

// 表单验证规则
const formRules = reactive({
  connID: [
    { required: true, message: '请选择数据源', trigger: 'change' }
  ],
  modelName: [
    { required: true, message: '请输入模型名称', trigger: 'blur' },
    { min: 2, max: 100, message: '模型名称长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  modelCode: [
    { required: true, message: '请输入模型编码', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '模型编码只能包含字母、数字和下划线', trigger: 'blur' },
    { min: 2, max: 50, message: '模型编码长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  modelVersion: [
    { required: true, message: '请输入模型版本', trigger: 'blur' }
  ],
  modelKind: [
    { required: true, message: '请选择模型类型', trigger: 'change' }
  ]
})

// 生命周期钩子
onMounted(() => {
  fetchDataSources()
})

// 获取数据源列表
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

// 获取表/视图列表
const fetchTables = async () => {
  if (!modelForm.connID) return
  
  loadingTables.value = true
  try {
    // 模拟API调用，实际需要替换为真实API
    // const response = await getDataSourceTables(modelForm.connID, objectType.value)
    // tables.value = response.data
    
    // 模拟数据
    tables.value = [
      {
        id: 1,
        label: 'user',
        comment: '用户表',
        isLeaf: true
      },
      {
        id: 2,
        label: 'order',
        comment: '订单表',
        isLeaf: true
      },
      {
        id: 3,
        label: 'product',
        comment: '产品表',
        isLeaf: true
      }
    ]
  } catch (error) {
    console.error('获取表列表失败:', error)
    ElMessage.error('获取表列表失败')
  } finally {
    loadingTables.value = false
  }
}

// 数据源变化
const handleDataSourceChange = (connID: number) => {
  const conn = dataSources.value.find(item => item.id === connID)
  if (conn) {
    modelForm.connName = conn.connName
  }
  // 重置表选择
  selectedTable.value = null
  tables.value = []
  // 获取表列表
  fetchTables()
}

// 对象类型变化
const handleObjectTypeChange = () => {
  // 重置表选择
  selectedTable.value = null
  tables.value = []
  // 获取表列表
  fetchTables()
}

// 表点击
const handleTableClick = (data: any) => {
  selectedTable.value = data
}

// 下一步
const nextStep = () => {
  activeStep.value++
}

// 上一步
const prevStep = () => {
  activeStep.value--
}

// 表单提交
const handleSubmit = async () => {
  try {
    await modelFormRef.value.validate()
    loading.value = true
    
    // 模拟创建模型
    // let result
    // if (objectType.value === 'table') {
    //   result = await buildModelFromTable({
    //     connID: modelForm.connID as number,
    //     tableID: selectedTable.value.id,
    //     modelName: modelForm.modelName as string,
    //     modelCode: modelForm.modelCode as string,
    //     modelVersion: modelForm.modelVersion as string,
    //     modelKind: modelForm.modelKind as number,
    //     isPublic: modelForm.isPublic as boolean,
    //     remark: modelForm.remark as string
    //   })
    // } else {
    //   result = await buildModelFromView({
    //     connID: modelForm.connID as number,
    //     tableID: selectedTable.value.id,
    //     modelName: modelForm.modelName as string,
    //     modelCode: modelForm.modelCode as string,
    //     modelVersion: modelForm.modelVersion as string,
    //     modelKind: modelForm.modelKind as number,
    //     isPublic: modelForm.isPublic as boolean,
    //     remark: modelForm.remark as string
    //   })
    // }
    
    // 跳转到完成步骤
    nextStep()
  } catch (error) {
    console.error('创建模型失败:', error)
    ElMessage.error('创建模型失败')
  } finally {
    loading.value = false
  }
}

// 完成后查看模型列表
const handleFinish = () => {
  router.push('/models')
}

// 创建另一个模型
const handleCreateAnother = () => {
  // 重置表单
  activeStep.value = 0
  selectedTable.value = null
  Object.assign(modelForm, {
    connID: undefined,
    connName: '',
    modelName: '',
    modelCode: '',
    modelVersion: '1.0.0',
    modelKind: 2,
    isPublic: false,
    remark: ''
  })
}
</script>

<style scoped>
.model-create {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.step-content {
  margin-top: 40px;
}

.step-content h3 {
  margin-bottom: 20px;
  font-size: 18px;
  font-weight: 600;
}

.step-actions {
  margin-top: 40px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.option-content {
  display: flex;
  flex-direction: column;
}

.option-subtitle {
  font-size: 12px;
  color: #909399;
}

.custom-tree-node {
  display: flex;
  flex-direction: column;
}

.node-subtitle {
  font-size: 12px;
  color: #909399;
}

.success-step {
  text-align: center;
  padding: 40px 0;
}
</style>
