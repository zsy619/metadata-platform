<template>
  <div class="model-edit">
    <div class="page-header">
      <h1>编辑模型</h1>
      <el-button @click="handleCancel" :icon="ArrowLeft">
        返回列表
      </el-button>
    </div>

    <el-card>
      <el-tabs v-model:active-name="activeTab" type="border-card">
        <!-- 基本信息 -->
        <el-tab-pane label="基本信息" name="basic">
          <el-form
            ref="modelFormRef"
            :model="modelForm"
            :rules="formRules"
            label-width="120px"
            class="edit-form"
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
                    :disabled="modelForm.isLocked"
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
                    :disabled="modelForm.isLocked"
                  >
                    <el-option label="SQL语句" :value="1" />
                    <el-option label="视图/表" :value="2" />
                    <el-option label="存储过程" :value="3" />
                    <el-option label="关联" :value="4" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="数据源">
                  <el-select
                    v-model="modelForm.connID"
                    placeholder="请选择数据源"
                    disabled
                  >
                    <el-option
                      v-for="conn in dataSources"
                      :key="conn.id"
                      :label="conn.connName"
                      :value="conn.id"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="数据源表">
                  <el-input
                    v-model="modelForm.tableName"
                    placeholder="数据源表"
                    disabled
                  />
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="是否公开">
                  <el-switch v-model="modelForm.isPublic" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="是否锁定">
                  <el-switch v-model="modelForm.isLocked" />
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-form-item label="模型图片">
              <el-upload
                class="avatar-uploader"
                action="#"
                :show-file-list="false"
                :on-success="handleLogoUpload"
                :before-upload="beforeLogoUpload"
              >
                <img v-if="modelForm.modelLogo" :src="modelForm.modelLogo" class="avatar" />
                <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
              </el-upload>
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
            
            <div class="form-actions">
              <el-button @click="handleCancel">取消</el-button>
              <el-button type="primary" @click="handleSubmit" :loading="submitting">
                {{ submitting ? '提交中...' : '提交' }}
              </el-button>
            </div>
          </el-form>
        </el-tab-pane>

        <!-- 字段配置 -->
        <el-tab-pane label="字段配置" name="fields">
          <div class="field-config">
            <div class="field-header">
              <h3>字段列表</h3>
              <el-button type="primary" @click="handleAddField" :icon="Plus" :disabled="modelForm.isLocked">
                添加字段
              </el-button>
            </div>
            
            <el-table
              v-loading="loadingFields"
              :data="modelFields"
              border
              style="width: 100%"
            >
              <el-table-column type="selection" width="55" />
              <el-table-column prop="columnName" label="字段名称" width="180" />
              <el-table-column prop="columnTitle" label="字段标题" width="180" />
              <el-table-column prop="showTitle" label="显示名称" width="180">
                <template #default="scope">
                  <el-input
                    v-model="scope.row.showTitle"
                    size="small"
                    placeholder="请输入显示名称"
                    :disabled="modelForm.isLocked"
                  />
                </template>
              </el-table-column>
              <el-table-column prop="isShow" label="是否显示" width="120">
                <template #default="scope">
                  <el-switch
                    v-model="scope.row.isShow"
                    active-value="1"
                    inactive-value="0"
                    :disabled="modelForm.isLocked"
                  />
                </template>
              </el-table-column>
              <el-table-column prop="showWidth" label="显示宽度" width="120">
                <template #default="scope">
                  <el-input-number
                    v-model="scope.row.showWidth"
                    :min="50"
                    :max="500"
                    size="small"
                    :disabled="modelForm.isLocked"
                  />
                </template>
              </el-table-column>
              <el-table-column prop="sort" label="排序" width="100">
                <template #default="scope">
                  <el-input-number
                    v-model="scope.row.sort"
                    :min="0"
                    :max="1000"
                    size="small"
                    :disabled="modelForm.isLocked"
                  />
                </template>
              </el-table-column>
              <el-table-column label="操作" width="150" fixed="right">
                <template #default="scope">
                  <el-button
                    type="primary"
                    size="small"
                    :icon="Edit"
                    @click="handleEditField(scope.row)"
                    :disabled="modelForm.isLocked"
                  >
                    编辑
                  </el-button>
                  <el-button
                    type="danger"
                    size="small"
                    :icon="Delete"
                    @click="handleDeleteField(scope.row)"
                    :disabled="modelForm.isLocked"
                  >
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
            
            <div class="field-actions">
              <el-button type="warning" @click="handleSaveFields" :loading="savingFields" :disabled="modelForm.isLocked || !hasFieldChanges">
                {{ savingFields ? '保存中...' : '保存字段配置' }}
              </el-button>
            </div>
          </div>
        </el-tab-pane>

        <!-- 高级配置 -->
        <el-tab-pane label="高级配置" name="advanced">
          <div class="advanced-config">
            <h3>高级配置</h3>
            <el-form
              :model="advancedConfig"
              label-width="120px"
              class="advanced-form"
            >
              <el-form-item label="缓存配置">
                <el-switch v-model="advancedConfig.enableCache" />
              </el-form-item>
              
              <el-form-item label="缓存过期时间（秒）" v-if="advancedConfig.enableCache">
                <el-input-number
                  v-model="advancedConfig.cacheExpire"
                  :min="1"
                  :max="86400"
                  placeholder="请输入缓存过期时间"
                />
              </el-form-item>
              
              <el-form-item label="查询超时时间（秒）">
                <el-input-number
                  v-model="advancedConfig.queryTimeout"
                  :min="1"
                  :max="60"
                  placeholder="请输入查询超时时间"
                />
              </el-form-item>
              
              <el-form-item label="是否启用审计">
                <el-switch v-model="advancedConfig.enableAudit" />
              </el-form-item>
              
              <div class="form-actions">
                <el-button type="primary" @click="handleSaveAdvanced" :loading="savingAdvanced">
                  {{ savingAdvanced ? '保存中...' : '保存高级配置' }}
                </el-button>
              </div>
            </el-form>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeft, Plus, Edit, Delete } from '@element-plus/icons-vue'
import type { DataSource } from '@/types/data-source'
import type { Model, ModelField } from '@/types/model'
import { getDataSources } from '@/api/data-source'
import { getModelById, updateModel, getModelFields, updateModelField, deleteModelField } from '@/api/model'

const router = useRouter()
const route = useRoute()
const modelFormRef = ref()

// 响应式数据
const activeTab = ref('basic')
const loading = ref(false)
const submitting = ref(false)
const savingFields = ref(false)
const savingAdvanced = ref(false)
const loadingFields = ref(false)
const hasFieldChanges = ref(false)

const modelForm = reactive<Partial<Model & { tableName?: string }>>({})
const modelFields = ref<ModelField[]>([])
const dataSources = ref<DataSource[]>([])

// 高级配置
const advancedConfig = reactive({
  enableCache: false,
  cacheExpire: 300,
  queryTimeout: 30,
  enableAudit: true
})

// 表单验证规则
const formRules = reactive({
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
  const modelId = Number(route.params.id)
  if (modelId) {
    fetchModel(modelId)
    fetchModelFields(modelId)
  }
  fetchDataSources()
})

// 监听字段变化
watch(
  () => [...modelFields.value],
  (newVal, oldVal) => {
    hasFieldChanges.value = JSON.stringify(newVal) !== JSON.stringify(oldVal)
  },
  { deep: true }
)

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

// 获取模型详情
const fetchModel = async (modelId: number) => {
  loading.value = true
  try {
    // 模拟API调用，实际需要替换为真实API
    // const data = await getModelById(modelId)
    // Object.assign(modelForm, data)
    
    // 模拟数据
    Object.assign(modelForm, {
      modelID: modelId,
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
      updateAt: '2024-01-23 10:00:00',
      tableName: 'user'
    })
  } catch (error) {
    console.error('获取模型详情失败:', error)
    ElMessage.error('获取模型详情失败')
  } finally {
    loading.value = false
  }
}

// 获取模型字段列表
const fetchModelFields = async (modelId: number) => {
  loadingFields.value = true
  try {
    // 模拟API调用，实际需要替换为真实API
    // const data = await getModelFields(modelId)
    // modelFields.value = data
    
    // 模拟数据
    modelFields.value = [
      {
        fieldID: 1,
        modelID: modelId,
        tableSchema: '',
        tableID: 1,
        tableName: 'user',
        tableTitle: '用户表',
        columnID: 1,
        columnName: 'id',
        columnTitle: 'ID',
        func: '',
        aggFunc: '',
        isShow: 1,
        showTitle: 'ID',
        showWidth: 100,
        state: 1,
        remark: '',
        sort: 0
      },
      {
        fieldID: 2,
        modelID: modelId,
        tableSchema: '',
        tableID: 1,
        tableName: 'user',
        tableTitle: '用户表',
        columnID: 2,
        columnName: 'username',
        columnTitle: '用户名',
        func: '',
        aggFunc: '',
        isShow: 1,
        showTitle: '用户名',
        showWidth: 180,
        state: 1,
        remark: '',
        sort: 1
      },
      {
        fieldID: 3,
        modelID: modelId,
        tableSchema: '',
        tableID: 1,
        tableName: 'user',
        tableTitle: '用户表',
        columnID: 3,
        columnName: 'email',
        columnTitle: '邮箱',
        func: '',
        aggFunc: '',
        isShow: 1,
        showTitle: '邮箱',
        showWidth: 200,
        state: 1,
        remark: '',
        sort: 2
      }
    ]
  } catch (error) {
    console.error('获取模型字段列表失败:', error)
    ElMessage.error('获取模型字段列表失败')
  } finally {
    loadingFields.value = false
  }
}

// 返回列表
const handleCancel = () => {
  router.push('/models')
}

// 表单提交
const handleSubmit = async () => {
  try {
    await modelFormRef.value.validate()
    submitting.value = true
    const modelId = Number(route.params.id)
    
    // 模拟更新模型
    // await updateModel(modelId, modelForm)
    setTimeout(() => {
      ElMessage.success('更新成功')
      submitting.value = false
    }, 1000)
  } catch (error) {
    console.error('表单验证失败:', error)
    submitting.value = false
  }
}

// 保存字段配置
const handleSaveFields = async () => {
  savingFields.value = true
  try {
    // 模拟保存字段配置
    // 实际需要遍历字段列表，调用updateModelField更新每个字段
    setTimeout(() => {
      ElMessage.success('字段配置保存成功')
      hasFieldChanges.value = false
      savingFields.value = false
    }, 1000)
  } catch (error) {
    console.error('保存字段配置失败:', error)
    ElMessage.error('保存字段配置失败')
  } finally {
    savingFields.value = false
  }
}

// 保存高级配置
const handleSaveAdvanced = async () => {
  savingAdvanced.value = true
  try {
    // 模拟保存高级配置
    setTimeout(() => {
      ElMessage.success('高级配置保存成功')
      savingAdvanced.value = false
    }, 1000)
  } catch (error) {
    console.error('保存高级配置失败:', error)
    ElMessage.error('保存高级配置失败')
  } finally {
    savingAdvanced.value = false
  }
}

// 添加字段
const handleAddField = () => {
  // 打开添加字段对话框
  ElMessage.info('添加字段功能开发中')
}

// 编辑字段
const handleEditField = (field: ModelField) => {
  // 打开编辑字段对话框
  ElMessage.info('编辑字段功能开发中')
}

// 删除字段
const handleDeleteField = (field: ModelField) => {
  ElMessageBox.confirm(
    `确定要删除字段 "${field.columnName}" 吗？`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      // 模拟删除字段
      // await deleteModelField(field.modelID, field.fieldID)
      const index = modelFields.value.findIndex(item => item.fieldID === field.fieldID)
      if (index > -1) {
        modelFields.value.splice(index, 1)
      }
      hasFieldChanges.value = true
      ElMessage.success('删除成功')
    } catch (error) {
      console.error('删除字段失败:', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {
    // 取消删除
  })
}

// 图片上传成功
const handleLogoUpload = (response: any, file: any) => {
  // 模拟上传成功
  modelForm.modelLogo = URL.createObjectURL(file.raw)
  ElMessage.success('上传成功')
}

// 图片上传前验证
const beforeLogoUpload = (file: any) => {
  const isImage = file.type.indexOf('image/') === 0
  const isLt2M = file.size / 1024 / 1024 < 2
  
  if (!isImage) {
    ElMessage.error('请上传图片文件')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB')
    return false
  }
  
  return true
}
</script>

<style scoped>
.model-edit {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.edit-form,
.advanced-form {
  max-width: 800px;
}

.form-actions {
  margin-top: 30px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.field-config {
  margin-top: 20px;
}

.field-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.field-actions {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.advanced-config {
  margin-top: 20px;
}

.advanced-config h3 {
  margin-bottom: 20px;
  font-size: 16px;
  font-weight: 600;
}

.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: border-color 0.3s;
}

.avatar-uploader .el-upload:hover {
  border-color: #409eff;
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  line-height: 178px;
  text-align: center;
}

.avatar {
  width: 178px;
  height: 178px;
  display: block;
}
</style>
