<template>
  <div class="data-source-create">
    <div class="page-header">
      <h1>创建数据源</h1>
      <el-button @click="handleCancel" :icon="ArrowLeft">
        返回列表
      </el-button>
    </div>

    <el-card>
      <el-form
        ref="dataSourceFormRef"
        :model="dataSourceForm"
        :rules="formRules"
        label-width="120px"
        class="create-form"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="数据源名称" prop="connName">
              <el-input
                v-model="dataSourceForm.connName"
                placeholder="请输入数据源名称"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="数据源类型" prop="connKind">
              <el-select
                v-model="dataSourceForm.connKind"
                placeholder="请选择数据源类型"
                clearable
              >
                <el-option label="MySQL" value="MySQL" />
                <el-option label="PostgreSQL" value="PostgreSQL" />
                <el-option label="SQL Server" value="SQL Server" />
                <el-option label="Oracle" value="Oracle" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="数据库版本" prop="connVersion">
              <el-input
                v-model="dataSourceForm.connVersion"
                placeholder="请输入数据库版本，如：8.0"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="主机地址" prop="connHost">
              <el-input
                v-model="dataSourceForm.connHost"
                placeholder="请输入主机地址，如：localhost"
                clearable
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="端口号" prop="connPort">
              <el-input-number
                v-model="dataSourceForm.connPort"
                :min="1"
                :max="65535"
                placeholder="请输入端口号"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="数据库名" prop="connDatabase">
              <el-input
                v-model="dataSourceForm.connDatabase"
                placeholder="请输入数据库名称"
                clearable
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="用户名" prop="connUser">
              <el-input
                v-model="dataSourceForm.connUser"
                placeholder="请输入数据库用户名"
                clearable
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="密码" prop="connPassword">
              <el-input
                v-model="dataSourceForm.connPassword"
                type="password"
                placeholder="请输入数据库密码"
                clearable
                show-password
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="备注">
          <el-input
            v-model="dataSourceForm.remark"
            type="textarea"
            placeholder="请输入备注信息"
            :rows="4"
            clearable
          />
        </el-form-item>

        <div class="form-actions">
          <el-button @click="handleCancel">取消</el-button>
          <el-button type="primary" @click="handleTestConnection" :loading="testingConnection">
            {{ testingConnection ? '测试连接中...' : '测试连接' }}
          </el-button>
          <el-button type="success" @click="handleSubmit" :loading="submitting">
            {{ submitting ? '提交中...' : '提交' }}
          </el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import type { DataSource } from '@/types/data-source'
import { createDataSource, testDataSourceConnection } from '@/api/data-source'

const router = useRouter()
const dataSourceFormRef = ref()

// 响应式数据
const testingConnection = ref(false)
const submitting = ref(false)

const dataSourceForm = reactive<Partial<DataSource>>({
  connName: '',
  connKind: 'MySQL',
  connVersion: '8.0',
  connHost: 'localhost',
  connPort: 3306,
  connUser: 'root',
  connPassword: '',
  connDatabase: '',
  remark: ''
})

// 表单验证规则
const formRules = reactive({
  connName: [
    { required: true, message: '请输入数据源名称', trigger: 'blur' },
    { min: 2, max: 100, message: '数据源名称长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  connKind: [
    { required: true, message: '请选择数据源类型', trigger: 'change' }
  ],
  connHost: [
    { required: true, message: '请输入主机地址', trigger: 'blur' }
  ],
  connPort: [
    { required: true, message: '请输入端口号', trigger: 'blur' },
    { type: 'number', message: '端口号必须是数字', trigger: 'blur' }
  ],
  connUser: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  connDatabase: [
    { required: true, message: '请输入数据库名', trigger: 'blur' }
  ]
})

// 返回列表
const handleCancel = () => {
  router.push('/data-sources')
}

// 测试连接
const handleTestConnection = async () => {
  try {
    await dataSourceFormRef.value.validate(validateFields)
    testingConnection.value = true

    // 模拟测试连接
    // await testDataSourceConnection(dataSourceForm)
    setTimeout(() => {
      ElMessage.success('连接测试成功')
      testingConnection.value = false
    }, 1000)
  } catch (error) {
    console.error('表单验证失败:', error)
    testingConnection.value = false
  }
}

// 表单提交
const handleSubmit = async () => {
  try {
    await dataSourceFormRef.value.validate()
    submitting.value = true

    // 模拟创建数据源
    // await createDataSource(dataSourceForm)
    setTimeout(() => {
      ElMessage.success('创建成功')
      submitting.value = false
      router.push('/data-sources')
    }, 1000)
  } catch (error) {
    console.error('表单验证失败:', error)
    submitting.value = false
  }
}

// 验证字段列表
const validateFields = ['connName', 'connKind', 'connHost', 'connPort', 'connUser', 'connDatabase']
</script>

<style scoped>
.data-source-create {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.create-form {
  max-width: 800px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 30px;
  gap: 10px;
}
</style>