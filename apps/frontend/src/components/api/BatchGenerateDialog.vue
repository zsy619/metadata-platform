<template>
    <el-dialog v-model="visible" title="批量生成接口" width="800px" :close-on-click-modal="false" class="batch-generate-dialog">
        <div class="dialog-content">
            <el-steps :active="step" finish-status="success" simple style="margin-bottom: 20px">
                <el-step title="选择模型" />
                <el-step title="配置接口" />
                <el-step title="预览确认" />
            </el-steps>
            <div v-show="step === 0" class="step-content">
                <div class="step-tip">
                    <el-icon><InfoFilled /></el-icon>
                    <span>选择要生成接口的模型，系统将自动生成标准的CRUD接口</span>
                </div>
                <el-form-item label="选择模型" required>
                    <el-select v-model="selectedModelId" filterable placeholder="搜索模型名称" style="width: 100%" @change="handleModelChange">
                        <el-option v-for="model in models" :key="model.id" :label="model.name" :value="model.id">
                            <div class="model-option">
                                <span class="model-name">{{ model.name }}</span>
                                <el-tag size="small" type="info">{{ model.code }}</el-tag>
                            </div>
                        </el-option>
                    </el-select>
                </el-form-item>
                <div v-if="selectedModel" class="model-info">
                    <el-descriptions :column="2" border size="small">
                        <el-descriptions-item label="模型名称">{{ selectedModel.name }}</el-descriptions-item>
                        <el-descriptions-item label="模型编码">{{ selectedModel.code }}</el-descriptions-item>
                        <el-descriptions-item label="模型类型">{{ getModelTypeText(selectedModel.type) }}</el-descriptions-item>
                        <el-descriptions-item label="字段数量">{{ selectedModel.fieldCount || 0 }}</el-descriptions-item>
                    </el-descriptions>
                </div>
            </div>
            <div v-show="step === 1" class="step-content">
                <div class="step-tip">
                    <el-icon><InfoFilled /></el-icon>
                    <span>选择需要生成的接口类型，系统将自动生成对应的接口配置</span>
                </div>
                <el-form-item label="接口类型">
                    <el-checkbox-group v-model="selectedApiTypes">
                        <el-checkbox label="list">
                            <div class="api-type-item">
                                <el-tag type="success" size="small">GET</el-tag>
                                <span>列表查询</span>
                                <span class="api-desc">分页查询数据列表</span>
                            </div>
                        </el-checkbox>
                        <el-checkbox label="get">
                            <div class="api-type-item">
                                <el-tag type="success" size="small">GET</el-tag>
                                <span>详情查询</span>
                                <span class="api-desc">根据ID查询单条数据</span>
                            </div>
                        </el-checkbox>
                        <el-checkbox label="create">
                            <div class="api-type-item">
                                <el-tag type="primary" size="small">POST</el-tag>
                                <span>新增数据</span>
                                <span class="api-desc">创建新的数据记录</span>
                            </div>
                        </el-checkbox>
                        <el-checkbox label="update">
                            <div class="api-type-item">
                                <el-tag type="warning" size="small">PUT</el-tag>
                                <span>更新数据</span>
                                <span class="api-desc">更新指定ID的数据</span>
                            </div>
                        </el-checkbox>
                        <el-checkbox label="delete">
                            <div class="api-type-item">
                                <el-tag type="danger" size="small">DELETE</el-tag>
                                <span>删除数据</span>
                                <span class="api-desc">删除指定ID的数据</span>
                            </div>
                        </el-checkbox>
                    </el-checkbox-group>
                </el-form-item>
                <el-form-item label="路径前缀">
                    <el-input v-model="pathPrefix" placeholder="如: /api/data" />
                </el-form-item>
                <el-form-item label="接口前缀">
                    <el-input v-model="apiPrefix" placeholder="如: user" />
                </el-form-item>
                <el-form-item label="选项配置">
                    <el-checkbox v-model="options.withAuth">需要鉴权</el-checkbox>
                    <el-checkbox v-model="options.withAudit">需要审计</el-checkbox>
                    <el-checkbox v-model="options.withValidation">启用验证</el-checkbox>
                </el-form-item>
            </div>
            <div v-show="step === 2" class="step-content">
                <div class="step-tip">
                    <el-icon><InfoFilled /></el-icon>
                    <span>预览即将生成的接口列表，确认无误后点击确认生成</span>
                </div>
                <div v-if="conflictCount > 0" class="conflict-warning">
                    <el-alert type="warning" :closable="false">
                        <template #title>
                            <span>检测到 {{ conflictCount }} 个冲突接口，这些接口已存在，将被覆盖</span>
                        </template>
                    </el-alert>
                </div>
                <el-table :data="previewApis" border size="small" max-height="300">
                    <el-table-column type="index" label="序号" width="50" />
                    <el-table-column label="接口名称" prop="name" min-width="150" />
                    <el-table-column label="接口路径" min-width="200">
                        <template #default="{ row }">
                            <el-tag size="small" :type="getMethodTagType(row.method)">{{ row.method }}</el-tag>
                            <span style="margin-left: 8px">{{ row.path }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="接口编码" prop="code" width="150" />
                    <el-table-column label="状态" width="80">
                        <template #default="{ row }">
                            <el-tag v-if="row.isNew" type="success" size="small">新增</el-tag>
                            <el-tag v-else type="warning" size="small">覆盖</el-tag>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
        </div>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="handleClose">取消</el-button>
                <el-button v-if="step > 0" @click="handlePrevStep">上一步</el-button>
                <el-button v-if="step < 2" type="primary" :disabled="!canNextStep" @click="handleNextStep">下一步</el-button>
                <el-button v-else type="primary" :loading="generating" @click="handleGenerate">
                    {{ generating ? '生成中...' : '确认生成' }}
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>
<script setup lang="ts">
import { InfoFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, ref, watch } from 'vue'

interface ModelOption {
    id: number
    name: string
    code: string
    type: number
    fieldCount?: number
}

interface PreviewApi {
    name: string
    path: string
    method: string
    code: string
    isNew: boolean
}

const props = defineProps<{
    modelValue: boolean
}>()

const emit = defineEmits(['update:modelValue', 'success'])

const visible = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
})

const step = ref(0)
const models = ref<ModelOption[]>([])
const selectedModelId = ref<number>()
const selectedApiTypes = ref<string[]>(['list', 'get', 'create', 'update', 'delete'])
const pathPrefix = ref('/api/data')
const apiPrefix = ref('')
const generating = ref(false)
const conflictCount = ref(0)

const options = ref({
    withAuth: true,
    withAudit: false,
    withValidation: true
})

const selectedModel = computed(() => models.value.find(m => m.id === selectedModelId.value))

const canNextStep = computed(() => {
    if (step.value === 0) return !!selectedModelId.value
    if (step.value === 1) return selectedApiTypes.value.length > 0
    return true
})

const previewApis = computed<PreviewApi[]>(() => {
    if (!selectedModel.value) return []
    const modelCode = apiPrefix.value || selectedModel.value.code
    const apiList: PreviewApi[] = []
    const existingCodes = ['user_list', 'user_get', 'user_create', 'user_update', 'user_delete']
    const typeMap: Record<string, { name: string; method: string; path: string }> = {
        list: { name: '列表查询', method: 'GET', path: `${pathPrefix.value}/${modelCode}/list` },
        get: { name: '详情查询', method: 'GET', path: `${pathPrefix.value}/${modelCode}/:id` },
        create: { name: '新增数据', method: 'POST', path: `${pathPrefix.value}/${modelCode}` },
        update: { name: '更新数据', method: 'PUT', path: `${pathPrefix.value}/${modelCode}/:id` },
        delete: { name: '删除数据', method: 'DELETE', path: `${pathPrefix.value}/${modelCode}/:id` }
    }
    selectedApiTypes.value.forEach(type => {
        const config = typeMap[type]
        if (config) {
            const code = `${modelCode}_${type}`
            apiList.push({
                name: `${selectedModel.value!.name}${config.name}`,
                path: config.path,
                method: config.method,
                code: code,
                isNew: !existingCodes.includes(code)
            })
        }
    })
    conflictCount.value = apiList.filter(a => !a.isNew).length
    return apiList
})

watch(visible, (val) => {
    if (val) {
        step.value = 0
        selectedModelId.value = undefined
        loadModels()
    }
})

const loadModels = () => {
    models.value = [
        { id: 1, name: '用户模型', code: 'user', type: 1, fieldCount: 10 },
        { id: 2, name: '订单模型', code: 'order', type: 1, fieldCount: 15 },
        { id: 3, name: '产品模型', code: 'product', type: 1, fieldCount: 20 }
    ]
}

const handleModelChange = () => {
    apiPrefix.value = selectedModel.value?.code || ''
}

const getModelTypeText = (type: number) => {
    const map: Record<number, string> = { 1: 'SQL模型', 2: '视图', 3: '存储过程', 4: '关联模型' }
    return map[type] || 'SQL模型'
}

const getMethodTagType = (method: string) => {
    const map: Record<string, string> = { GET: 'success', POST: 'primary', PUT: 'warning', DELETE: 'danger' }
    return map[method] || 'info'
}

const handleNextStep = () => {
    if (step.value < 2) step.value++
}

const handlePrevStep = () => {
    if (step.value > 0) step.value--
}

const handleGenerate = async () => {
    if (previewApis.value.length === 0) {
        ElMessage.warning('请选择至少一种接口类型')
        return
    }
    generating.value = true
    try {
        await new Promise(resolve => setTimeout(resolve, 1000))
        ElMessage.success(`成功生成 ${previewApis.value.length} 个接口`)
        emit('success')
        handleClose()
    } catch (error) {
        ElMessage.error('生成失败')
    } finally {
        generating.value = false
    }
}

const handleClose = () => {
    visible.value = false
}
</script>
<style scoped>
.dialog-content {
    min-height: 300px;
}

.step-content {
    padding: 10px 0;
}

.step-tip {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 12px;
    background: #f0f9ff;
    border-radius: 6px;
    color: #409eff;
    font-size: 13px;
    margin-bottom: 16px;
}

.model-option {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.model-name {
    font-weight: 500;
}

.model-info {
    margin-top: 16px;
}

.api-type-item {
    display: flex;
    align-items: center;
    gap: 8px;
}

.api-desc {
    color: #909399;
    font-size: 12px;
    margin-left: 8px;
}

.conflict-warning {
    margin-bottom: 16px;
}

.dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
}
</style>
