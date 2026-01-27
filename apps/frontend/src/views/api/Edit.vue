<template>
    <div class="api-edit">
        <div class="page-header">
            <h1>编辑接口</h1>
            <el-button @click="handleCancel" :icon="ArrowLeft">
                返回列表
            </el-button>
        </div>
        <el-card>
            <el-tabs v-model:active-name="activeTab" type="border-card">
                <!-- 基本信息 -->
                <el-tab-pane label="基本信息" name="basic">
                    <el-form ref="apiFormRef" :model="apiForm" :rules="formRules" label-width="120px" class="edit-form">
                        <el-row :gutter="20">
                            <el-col :span="12">
                                <el-form-item label="接口名称" prop="apiName">
                                    <el-input v-model="apiForm.apiName" placeholder="请输入接口名称" clearable />
                                </el-form-item>
                            </el-col>
                            <el-col :span="12">
                                <el-form-item label="接口编码" prop="apiCode">
                                    <el-input v-model="apiForm.apiCode" placeholder="请输入接口编码（字母数字下划线组合）" clearable />
                                </el-form-item>
                            </el-col>
                        </el-row>
                        <el-row :gutter="20">
                            <el-col :span="12">
                                <el-form-item label="所属模型" prop="modelID">
                                    <el-select v-model="apiForm.modelID" placeholder="请选择所属模型" @change="handleModelChange">
                                        <el-option v-for="model in models" :key="model.modelID" :label="model.modelName" :value="model.modelID">
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
                                    <el-select v-model="apiForm.apiType" placeholder="请选择接口类型">
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
                                    <el-select v-model="apiForm.apiMethod" placeholder="请选择请求方法">
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
                                    <el-input v-model="apiForm.apiPath" placeholder="请输入接口路径，如：/api/data/user_model/list" clearable />
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
                        <el-row :gutter="20">
                            <el-col :span="12">
                                <el-form-item label="接口状态">
                                    <el-switch v-model="apiForm.state" active-value="1" inactive-value="0" />
                                </el-form-item>
                            </el-col>
                        </el-row>
                        <el-form-item label="备注">
                            <el-input v-model="apiForm.remark" type="textarea" placeholder="请输入接口描述信息" :rows="4" clearable />
                        </el-form-item>
                        <div class="form-actions">
                            <el-button @click="handleCancel">取消</el-button>
                            <el-button type="primary" @click="handleSubmit" :loading="submitting">
                                {{ submitting ? '提交中...' : '提交' }}
                            </el-button>
                        </div>
                    </el-form>
                </el-tab-pane>
                <!-- 请求参数 -->
                <el-tab-pane label="请求参数" name="request">
                    <div class="param-config">
                        <div class="param-header">
                            <h3>请求参数列表</h3>
                            <el-button type="primary" @click="handleAddParam" :icon="Plus">
                                添加参数
                            </el-button>
                        </div>
                        <el-table :data="requestParams" border style="width: 100%">
                            <el-table-column prop="paramName" label="参数名称" width="150" />
                            <el-table-column prop="paramType" label="参数类型" width="120">
                                <template #default="scope">
                                    <el-tag :type="getParamTypeTagType(scope.row.paramType)">
                                        {{ getParamTypeText(scope.row.paramType) }}
                                    </el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column prop="dataType" label="数据类型" width="120" />
                            <el-table-column prop="isRequired" label="是否必填" width="120">
                                <template #default="scope">
                                    <el-switch v-model="scope.row.isRequired" active-value="1" inactive-value="0" />
                                </template>
                            </el-table-column>
                            <el-table-column prop="defaultValue" label="默认值" width="150">
                                <template #default="scope">
                                    <el-input v-model="scope.row.defaultValue" size="small" placeholder="请输入默认值" />
                                </template>
                            </el-table-column>
                            <el-table-column prop="sort" label="排序" width="100">
                                <template #default="scope">
                                    <el-input-number v-model="scope.row.sort" :min="0" :max="1000" size="small" />
                                </template>
                            </el-table-column>
                            <el-table-column prop="remark" label="描述" width="200">
                                <template #default="scope">
                                    <el-input v-model="scope.row.remark" size="small" placeholder="请输入描述" />
                                </template>
                            </el-table-column>
                            <el-table-column label="操作" width="120" fixed="right">
                                <template #default="scope">
                                    <el-button type="danger" size="small" :icon="Delete" @click="handleDeleteParam(scope.row)">
                                        删除
                                    </el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                        <div class="param-actions">
                            <el-button type="warning" @click="handleSaveParams" :loading="savingParams">
                                {{ savingParams ? '保存中...' : '保存参数配置' }}
                            </el-button>
                        </div>
                    </div>
                </el-tab-pane>
                <!-- 响应结构 -->
                <el-tab-pane label="响应结构" name="response">
                    <div class="response-config">
                        <div class="response-header">
                            <h3>响应字段列表</h3>
                            <el-button type="primary" @click="handleAddResponseField" :icon="Plus">
                                添加字段
                            </el-button>
                        </div>
                        <el-table :data="responseFields" border style="width: 100%">
                            <el-table-column prop="fieldName" label="字段名称" width="150" />
                            <el-table-column prop="dataType" label="数据类型" width="120" />
                            <el-table-column prop="isRequired" label="是否必填" width="120">
                                <template #default="scope">
                                    <el-switch v-model="scope.row.isRequired" active-value="1" inactive-value="0" />
                                </template>
                            </el-table-column>
                            <el-table-column prop="exampleValue" label="示例值" width="150">
                                <template #default="scope">
                                    <el-input v-model="scope.row.exampleValue" size="small" placeholder="请输入示例值" />
                                </template>
                            </el-table-column>
                            <el-table-column prop="sort" label="排序" width="100">
                                <template #default="scope">
                                    <el-input-number v-model="scope.row.sort" :min="0" :max="1000" size="small" />
                                </template>
                            </el-table-column>
                            <el-table-column prop="remark" label="描述" width="200">
                                <template #default="scope">
                                    <el-input v-model="scope.row.remark" size="small" placeholder="请输入描述" />
                                </template>
                            </el-table-column>
                            <el-table-column label="操作" width="120" fixed="right">
                                <template #default="scope">
                                    <el-button type="danger" size="small" :icon="Delete" @click="handleDeleteResponseField(scope.row)">
                                        删除
                                    </el-button>
                                </template>
                            </el-table-column>
                        </el-table>
                        <div class="response-actions">
                            <el-button type="warning" @click="handleSaveResponse" :loading="savingResponse">
                                {{ savingResponse ? '保存中...' : '保存响应配置' }}
                            </el-button>
                        </div>
                    </div>
                </el-tab-pane>
            </el-tabs>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import type { API, APIRequestParam, APIResponseData } from '@/types/api'
import type { Model } from '@/types/metadata'
import { ArrowLeft, Delete, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()
const apiFormRef = ref()

// 响应式数据
const activeTab = ref('basic')
const loading = ref(false)
const submitting = ref(false)
const savingParams = ref(false)
const savingResponse = ref(false)

const apiForm = reactive<Partial<API>>({})
const models = ref<Model[]>([])
const requestParams = ref<APIRequestParam[]>([])
const responseFields = ref<APIResponseData[]>([])

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
    const apiId = Number(route.params.id)
    if (apiId) {
        fetchAPI(apiId)
        fetchModels()
        fetchRequestParams(apiId)
        fetchResponseFields(apiId)
    }
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
            }
        ]
    } catch (error) {
        console.error('获取模型列表失败:', error)
        ElMessage.error('获取模型列表失败')
    }
}

// 获取接口详情
const fetchAPI = async (apiId: number) => {
    loading.value = true
    try {
        // 模拟API调用，实际需要替换为真实API
        // const data = await getAPIById(apiId)
        // Object.assign(apiForm, data)

        // 模拟数据
        Object.assign(apiForm, {
            apiID: apiId,
            apiName: '获取用户列表',
            apiCode: 'user_list',
            apiPath: '/api/data/user_model/list',
            apiMethod: 'GET',
            modelID: 1,
            modelName: '用户模型',
            apiType: 1,
            state: 1,
            needAuth: true,
            needAudit: true,
            remark: '获取用户列表接口',
            createID: 1,
            createBy: 'admin',
            createAt: '2024-01-23 10:00:00',
            updateID: 1,
            updateBy: 'admin',
            updateAt: '2024-01-23 10:00:00'
        })
    } catch (error) {
        console.error('获取接口详情失败:', error)
        ElMessage.error('获取接口详情失败')
    } finally {
        loading.value = false
    }
}

// 获取请求参数
const fetchRequestParams = async (apiId: number) => {
    try {
        // 模拟API调用，实际需要替换为真实API
        // const data = await getAPIRequestParams(apiId)
        // requestParams.value = data

        // 模拟数据
        requestParams.value = [
            {
                paramID: 1,
                apiID: apiId,
                paramName: 'page',
                paramType: 2,
                dataType: 'number',
                isRequired: false,
                defaultValue: '1',
                remark: '页码',
                sort: 0
            },
            {
                paramID: 2,
                apiID: apiId,
                paramName: 'pageSize',
                paramType: 2,
                dataType: 'number',
                isRequired: false,
                defaultValue: '10',
                remark: '每页条数',
                sort: 1
            },
            {
                paramID: 3,
                apiID: apiId,
                paramName: 'keyword',
                paramType: 2,
                dataType: 'string',
                isRequired: false,
                defaultValue: '',
                remark: '搜索关键词',
                sort: 2
            }
        ]
    } catch (error) {
        console.error('获取请求参数失败:', error)
        ElMessage.error('获取请求参数失败')
    }
}

// 获取响应字段
const fetchResponseFields = async (apiId: number) => {
    try {
        // 模拟API调用，实际需要替换为真实API
        // const data = await getAPIResponseData(apiId)
        // responseFields.value = data

        // 模拟数据
        responseFields.value = [
            {
                responseID: 1,
                apiID: apiId,
                fieldName: 'code',
                dataType: 'number',
                isRequired: true,
                exampleValue: '200',
                remark: '状态码',
                sort: 0
            },
            {
                responseID: 2,
                apiID: apiId,
                fieldName: 'message',
                dataType: 'string',
                isRequired: true,
                exampleValue: 'success',
                remark: '提示信息',
                sort: 1
            },
            {
                responseID: 3,
                apiID: apiId,
                fieldName: 'data',
                dataType: 'object',
                isRequired: true,
                exampleValue: '{"list": [], "total": 0}',
                remark: '响应数据',
                sort: 2
            }
        ]
    } catch (error) {
        console.error('获取响应字段失败:', error)
        ElMessage.error('获取响应字段失败')
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
    router.push('/api/list')
}

// 表单提交
const handleSubmit = async () => {
    try {
        await apiFormRef.value.validate()
        submitting.value = true
        const apiId = Number(route.params.id)

        // 模拟更新接口
        // await updateAPI(apiId, apiForm)
        setTimeout(() => {
            ElMessage.success('更新成功')
            submitting.value = false
        }, 1000)
    } catch (error) {
        console.error('表单验证失败:', error)
        submitting.value = false
    }
}

// 添加请求参数
const handleAddParam = () => {
    requestParams.value.push({
        paramID: Date.now(),
        apiID: apiForm.apiID as number,
        paramName: '',
        paramType: 2,
        dataType: 'string',
        isRequired: false,
        defaultValue: '',
        remark: '',
        sort: requestParams.value.length
    })
}

// 删除请求参数
const handleDeleteParam = (row: APIRequestParam) => {
    const index = requestParams.value.findIndex(item => item.paramID === row.paramID)
    if (index > -1) {
        requestParams.value.splice(index, 1)
    }
}

// 保存请求参数
const handleSaveParams = async () => {
    savingParams.value = true
    try {
        // 模拟保存请求参数
        setTimeout(() => {
            ElMessage.success('请求参数保存成功')
            savingParams.value = false
        }, 1000)
    } catch (error) {
        console.error('保存请求参数失败:', error)
        ElMessage.error('保存请求参数失败')
    } finally {
        savingParams.value = false
    }
}

// 添加响应字段
const handleAddResponseField = () => {
    responseFields.value.push({
        responseID: Date.now(),
        apiID: apiForm.apiID as number,
        fieldName: '',
        dataType: 'string',
        isRequired: false,
        exampleValue: '',
        remark: '',
        sort: responseFields.value.length
    })
}

// 删除响应字段
const handleDeleteResponseField = (row: APIResponseData) => {
    const index = responseFields.value.findIndex(item => item.responseID === row.responseID)
    if (index > -1) {
        responseFields.value.splice(index, 1)
    }
}

// 保存响应配置
const handleSaveResponse = async () => {
    savingResponse.value = true
    try {
        // 模拟保存响应配置
        setTimeout(() => {
            ElMessage.success('响应配置保存成功')
            savingResponse.value = false
        }, 1000)
    } catch (error) {
        console.error('保存响应配置失败:', error)
        ElMessage.error('保存响应配置失败')
    } finally {
        savingResponse.value = false
    }
}

// 获取参数类型文本
const getParamTypeText = (type: number): string => {
    const typeMap: Record<number, string> = {
        1: '路径参数',
        2: '查询参数',
        3: '请求体参数'
    }
    return typeMap[type] || '未知'
}

// 获取参数类型标签样式
const getParamTypeTagType = (type: number): string => {
    const typeMap: Record<number, string> = {
        1: 'primary',
        2: 'success',
        3: 'warning'
    }
    return typeMap[type] || 'info'
}
</script>
<style scoped>
.api-edit {
    padding: 10px;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.edit-form {
    max-width: 1000px;
}

.form-actions {
    margin-top: 30px;
    display: flex;
    justify-content: flex-end;
    gap: 10px;
}

.param-config,
.response-config {
    margin-top: 20px;
}

.param-header,
.response-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.param-actions,
.response-actions {
    margin-top: 20px;
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
    margin-top: 0;
    font-size: 16px;
    font-weight: 600;
}
</style>
