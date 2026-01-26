<template>
  <div class="api-create">
    <div class="page-header">
      <h1>创建接口</h1>
      <el-button @click="handleCancel" :icon="ArrowLeft">
        返回列表
      </el-button>
    </div>

    <el-card>
      <el-form
        ref="apiFormRef"
        :model="apiForm"
        :rules="formRules"
        label-width="120px"
        class="create-form"
      >
        <h3>基本信息</h3>
        <el-divider />
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="接口名称" prop="apiName">
              <el-input
                v-model="apiForm.apiName"
                placeholder="请输入接口名称"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="接口编码" prop="apiCode">
              <el-input
                v-model="apiForm.apiCode"
                placeholder="请输入接口编码（字母数字下划线组合）"
                clearable
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="所属模型" prop="modelID">
              <el-select
                v-model="apiForm.modelID"
                placeholder="请选择所属模型"
                @change="handleModelChange"
              >
                <el-option
                  v-for="model in models"
                  :key="model.modelID"
                  :label="model.modelName"
                  :value="model.modelID"
                >
                  <div class="option-content">
                    <span>{{ model.modelName }}</span>
                    <span class="option-subtitle">{{ model.modelCode }} - {{ model.modelKind === 2 ? '表' : '视图' }}</span>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="接口类型" prop="apiType">
              <el-select
                v-model="apiForm.apiType"
                placeholder="请选择接口类型"
              >
                <el-option label="查询" :value="1" />
                <el-option label="新增" :value="2" />
                <el-option label="更新" :value="3" />
                <el-option label="删除" :value="4" />
                <el-option label="自定义" :value="5" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="请求方法" prop="apiMethod">
              <el-select
                v-model="apiForm.apiMethod"
                placeholder="请选择请求方法"
              >
                <el-option label="GET" value="GET" />
                <el-option label="POST" value="POST" />
                <el-option label="PUT" value="PUT" />
                <el-option label="DELETE" value="DELETE" />
                <el-option label="PATCH" value="PATCH" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="接口路径" prop="apiPath">
              <el-input
                v-model="apiForm.apiPath"
                placeholder="请输入接口路径，如：/api/data/user_model/list"
                clearable
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="需要鉴权">
              <el-switch v-model="apiForm.needAuth" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="需要审计">
              <el-switch v-model="apiForm.needAudit" />
            </el-form-item>
          </el-col>
        </el-row>
        
        <h3>接口描述</h3>
        <el-divider />
        <el-form-item label="备注">
          <el-input
            v-model="apiForm.remark"
            type="textarea"
            placeholder="请输入接口描述信息"
            :rows="4"
            clearable
          />
        </el-form-item>
        
        <div class="form-actions">
          <el-button @click="handleCancel">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            {{ submitting ? '提交中...' : '提交' }}
          </el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import type { API } from '@/types/api'
import type { Model } from '@/types/model'
import { createAPI } from '@/api/api'
import { getModels } from '@/api/model'

const router = useRouter()
const apiFormRef = ref()

// 响应式数据
const submitting = ref(false)
const models = ref<Model[]>([])

const apiForm = reactive<Partial<API>>({
  apiName: '',
  apiCode: '',
  apiPath: '',
  apiMethod: 'GET',
  modelID: undefined,
  modelName: '',
  apiType: 1,
  state: 1,
  needAuth: true,
  needAudit: true,
  remark: ''
})

// 表单验证规则
const formRules = reactive({
  apiName: [
    { required: true, message: '请输入接口名称', trigger: 'blur' },
    { min: 2, max: 100, message: '接口名称长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  apiCode: [
    { required: true, message: '请输入接口编码', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '接口编码只能包含字母、数字和下划线', trigger: 'blur' },
    { min: 2, max: 50, message: '接口编码长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  modelID: [
    { required: true, message: '请选择所属模型', trigger: 'change' }
  ],
  apiPath: [
    { required: true, message: '请输入接口路径', trigger: 'blur' },
    { pattern: /^\/.*$/, message: '接口路径必须以斜杠开头', trigger: 'blur' }
  ],
  apiMethod: [
    { required: true, message: '请选择请求方法', trigger: 'change' }
  ],
  apiType: [
    { required: true, message: '请选择接口类型', trigger: 'change' }
  ]
})

// 生命周期钩子
onMounted(() => {
  fetchModels()
})

// 获取模型列表
const fetchModels = async () => {
  try {
    // 模拟API调用，实际需要替换为真实API
    // const response = await getModels()
    // models.value = response.data
    
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
  } catch (error) {
    console.error('获取模型列表失败:', error)
    ElMessage.error('获取模型列表失败')
  }
}

// 模型变化
const handleModelChange = (modelID: number) => {
  const model = models.value.find(item => item.modelID === modelID)
  if (model) {
    apiForm.modelName = model.modelName
  }
}

// 返回列表
const handleCancel = () => {
  router.push('/apis')
}

// 表单提交
const handleSubmit = async () => {
  try {
    await apiFormRef.value.validate()
    submitting.value = true
    
    // 模拟创建接口
    // await createAPI(apiForm)
    setTimeout(() => {
      ElMessage.success('创建成功')
      submitting.value = false
      router.push('/apis')
    }, 1000)
  } catch (error) {
    console.error('表单验证失败:', error)
    submitting.value = false
  }
}
</script>

<style scoped>
.api-create {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.create-form {
  max-width: 1000px;
}

.form-actions {
  margin-top: 30px;
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

h3 {
  margin-top: 20px;
  font-size: 16px;
  font-weight: 600;
}
</style>
