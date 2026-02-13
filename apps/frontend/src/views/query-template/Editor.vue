<template>
    <div class="query-template-editor container-padding">
        <div class="page-header">
            <h1 class="page-title">
                <el-icon class="title-icon">
                    <Document />
                </el-icon>
                {{ isEdit ? '编辑查询模板' : '新建查询模板' }}
            </h1>
            <div class="header-actions">
                <el-button @click="handleCancel">取消</el-button>
                <el-button type="info" @click="handleTest">测试</el-button>
                <el-button type="primary" :loading="saving" @click="handleSave">
                    {{ saving ? '保存中...' : '保存' }}
                </el-button>
            </div>
        </div>
        <div class="editor-content">
            <el-row :gutter="20">
                <el-col :span="14">
                    <el-card class="config-card">
                        <template #header>
                            <div class="card-header">
                                <span>模板配置</span>
                            </div>
                        </template>
                        <el-form ref="formRef" :model="formData" :rules="formRules" label-width="140px" label-position="right">
                            <el-form-item label="模板名称" prop="templateName">
                                <el-input v-model="formData.templateName" placeholder="请输入模板名称" />
                            </el-form-item>
                            <el-form-item label="模板编码" prop="templateCode">
                                <el-input v-model="formData.templateCode" placeholder="请输入模板编码" :disabled="isEdit" />
                            </el-form-item>
                            <el-form-item label="所属模型" prop="modelId">
                                <el-select v-model="formData.modelId" placeholder="请选择模型" @change="handleModelChange" style="width: 100%">
                                    <el-option v-for="model in models" :key="model.id" :label="model.name" :value="model.id">
                                        <div class="model-option">
                                            <span>{{ model.name }}</span>
                                            <el-tag size="small" type="info">{{ model.code }}</el-tag>
                                        </div>
                                    </el-option>
                                </el-select>
                            </el-form-item>
                            <el-form-item label="设为默认">
                                <el-switch v-model="formData.isDefault" />
                            </el-form-item>
                            <el-form-item label="备注">
                                <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入备注信息" />
                            </el-form-item>
                        </el-form>
                    </el-card>
                    <el-card class="condition-card">
                        <template #header>
                            <div class="card-header">
                                <span>查询条件配置</span>
                            </div>
                        </template>
                        <QueryBuilder ref="queryBuilderRef" :fields="modelFields" v-model="formData.conditions" :show-sql-preview="true" @change="handleConditionsChange" />
                    </el-card>
                </el-col>
                <el-col :span="10">
                    <el-card class="preview-card">
                        <template #header>
                            <div class="card-header">
                                <span>结果预览</span>
                                <el-button size="small" type="primary" :loading="previewLoading" @click="handlePreview">
                                    刷新预览
                                </el-button>
                            </div>
                        </template>
                        <div class="preview-content">
                            <div v-if="previewData.length === 0 && !previewLoading" class="empty-preview">
                                <el-empty description="暂无数据" />
                            </div>
                            <el-table v-else :data="previewData" border size="small" max-height="400">
                                <el-table-column v-for="col in previewColumns" :key="col.prop" :prop="col.prop" :label="col.label" min-width="100" show-overflow-tooltip />
                            </el-table>
                            <div v-if="previewTotal > 0" class="preview-pagination">
                                <el-pagination v-model:current-page="previewPage" v-model:page-size="previewPageSize" :total="previewTotal" small layout="prev, pager, next" @current-change="handlePreview" />
                            </div>
                        </div>
                    </el-card>
                </el-col>
            </el-row>
        </div>
    </div>
</template>
<script setup lang="ts">
import { Document } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import QueryBuilder from '@/components/query/QueryBuilder.vue'

interface Model {
    id: number
    name: string
    code: string
}

interface Field {
    name: string
    label: string
    type: string
    options?: { label: string; value: any }[]
}

const route = useRoute()
const router = useRouter()

const isEdit = ref(!!route.params.id)
const saving = ref(false)
const formRef = ref<FormInstance>()
const queryBuilderRef = ref()

const models = ref<Model[]>([
    { id: 1, name: '用户模型', code: 'user' },
    { id: 2, name: '订单模型', code: 'order' },
    { id: 3, name: '产品模型', code: 'product' }
])

const modelFields = ref<Field[]>([])

const formData = ref({
    templateName: '',
    templateCode: '',
    modelId: undefined as number | undefined,
    isDefault: false,
    remark: '',
    conditions: [] as any[],
    generatedSql: ''
})

const formRules: FormRules = {
    templateName: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
    templateCode: [{ required: true, message: '请输入模板编码', trigger: 'blur' }],
    modelId: [{ required: true, message: '请选择模型', trigger: 'change' }]
}

const previewLoading = ref(false)
const previewData = ref<any[]>([])
const previewColumns = ref<any[]>([])
const previewTotal = ref(0)
const previewPage = ref(1)
const previewPageSize = ref(10)

onMounted(() => {
    if (isEdit.value) {
        loadTemplate()
    }
})

const loadTemplate = () => {
    formData.value = {
        templateName: '用户查询模板',
        templateCode: 'user_query',
        modelId: 1,
        isDefault: true,
        remark: '用户查询模板',
        conditions: [],
        generatedSql: ''
    }
    loadModelFields(1)
}

const handleModelChange = (modelId: number) => {
    loadModelFields(modelId)
    formData.value.conditions = []
    formData.value.generatedSql = ''
}

const loadModelFields = (modelId: number) => {
    const fieldMap: Record<number, Field[]> = {
        1: [
            { name: 'id', label: '用户ID', type: 'int' },
            { name: 'username', label: '用户名', type: 'string' },
            { name: 'email', label: '邮箱', type: 'string' },
            { name: 'phone', label: '手机号', type: 'string' },
            { name: 'status', label: '状态', type: 'int', options: [{ label: '启用', value: 1 }, { label: '禁用', value: 0 }] },
            { name: 'createTime', label: '创建时间', type: 'date' }
        ],
        2: [
            { name: 'id', label: '订单ID', type: 'int' },
            { name: 'orderNo', label: '订单号', type: 'string' },
            { name: 'amount', label: '金额', type: 'decimal' },
            { name: 'status', label: '状态', type: 'int' }
        ],
        3: [
            { name: 'id', label: '产品ID', type: 'int' },
            { name: 'name', label: '产品名称', type: 'string' },
            { name: 'price', label: '价格', type: 'decimal' },
            { name: 'stock', label: '库存', type: 'int' }
        ]
    }
    modelFields.value = fieldMap[modelId] || []
}

const handleConditionsChange = (conditions: any[], sql: string) => {
    formData.value.generatedSql = sql
}

const handlePreview = async () => {
    if (!formData.value.modelId) {
        ElMessage.warning('请先选择模型')
        return
    }
    previewLoading.value = true
    try {
        await new Promise(resolve => setTimeout(resolve, 500))
        previewData.value = [
            { id: 1, username: '张三', email: 'zhangsan@example.com', phone: '13800138000', status: 1 },
            { id: 2, username: '李四', email: 'lisi@example.com', phone: '13800138001', status: 1 },
            { id: 3, username: '王五', email: 'wangwu@example.com', phone: '13800138002', status: 0 }
        ]
        previewColumns.value = Object.keys(previewData.value[0] || {}).map(key => ({
            prop: key,
            label: key
        }))
        previewTotal.value = 100
    } finally {
        previewLoading.value = false
    }
}

const handleTest = () => {
    const sql = queryBuilderRef.value?.getSql()
    ElMessage.success(`生成的SQL: ${sql || '无'}`)
}

const handleSave = async () => {
    if (!formRef.value) return
    await formRef.value.validate(async (valid) => {
        if (valid) {
            saving.value = true
            try {
                await new Promise(resolve => setTimeout(resolve, 1000))
                ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
                router.push('/query-template/list')
            } catch (error) {
                ElMessage.error('保存失败')
            } finally {
                saving.value = false
            }
        }
    })
}

const handleCancel = () => {
    router.back()
}
</script>
<style scoped>
.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.page-title {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 24px;
    font-weight: 600;
    color: #303133;
    margin: 0;
}

.title-icon {
    font-size: 24px;
    color: #409eff;
}

.header-actions {
    display: flex;
    gap: 10px;
}

.editor-content {
    min-height: calc(100vh - 150px);
}

.config-card,
.condition-card,
.preview-card {
    margin-bottom: 20px;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 600;
}

.model-option {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.preview-content {
    min-height: 300px;
}

.empty-preview {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 300px;
}

.preview-pagination {
    display: flex;
    justify-content: center;
    margin-top: 12px;
}
</style>
