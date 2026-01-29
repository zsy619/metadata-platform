<template>
    <div class="sql-test-container">
        <div class="page-header">
            <div class="header-left">
                <el-button @click="goBack" :icon="ArrowLeft" circle />
                <h1 class="title">SQL 模型测试</h1>
            </div>
            <el-tag type="info" effect="plain">验证动态查询接口与参数配置</el-tag>
        </div>
        <el-row :gutter="20">
            <!-- 左侧：配置区域 -->
            <el-col :span="8">
                <el-card class="config-card" shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <span>查询配置</span>
                            <el-button type="primary" link :icon="Refresh" @click="initModels">刷新列表</el-button>
                        </div>
                    </template>
                    <el-form :model="testForm" label-position="top">
                        <el-form-item label="选择模型">
                            <el-select v-model="selectedModelId" placeholder="请选择要测试的模型" filterable class="w-full" @change="handleModelChange">
                                <el-option v-for="item in models" :key="item.id" :label="`${item.model_name || item.modelName} (${item.model_code || item.modelCode})`" :value="item.id" />
                            </el-select>
                        </el-form-item>
                        <el-form-item label="查询模式">
                            <el-radio-group v-model="queryMode">
                                <el-radio-button label="id">按模型 ID</el-radio-button>
                                <el-radio-button label="code">按模型代码</el-radio-button>
                            </el-radio-group>
                        </el-form-item>
                        <div v-if="selectedModel" class="parameters-section">
                            <div class="section-title">
                                <h3>查询参数</h3>
                                <el-tag v-if="!queryParams.length" type="info" size="small">无</el-tag>
                            </div>
                            <el-table v-if="queryParams.length" :data="queryParams" border stripe size="small" style="width: 100%; margin-top: 10px">
                                <el-table-column prop="name" label="参数名称" width="120" />
                                <el-table-column prop="type" label="类型" width="100" />
                                <el-table-column label="值">
                                    <template #default="{ row }">
                                        <el-input v-model="testForm.params[row.name]" :placeholder="`输入${row.name}`" size="small" />
                                    </template>
                                </el-table-column>
                            </el-table>
                            <div v-else class="no-params-hint">该模型未定义动态参数</div>
                        </div>
                        <div class="action-bar">
                            <el-button type="primary" :loading="loading" :disabled="!selectedModelId" class="w-full" size="large" @click="performQuery">
                                {{ loading ? '查询中...' : '执行查询' }}
                            </el-button>
                            <div v-if="loading" class="progress-info">
                                <el-progress :percentage="queryProgress" :status="queryProgress === 100 ? 'success' : ''" :indeterminate="queryProgress === 0" />
                                <div class="progress-text">正在从数据源获取结果...</div>
                            </div>
                        </div>
                    </el-form>
                </el-card>
                <el-card class="path-card" shadow="never" v-if="selectedModel">
                    <template #header>
                        <div class="card-header">
                            <span>请求路径预览</span>
                            <el-tag size="small" type="success">POST</el-tag>
                        </div>
                    </template>
                    <div class="path-url">{{ urlPreview }}</div>
                    <div class="path-hint">您可以使用该路径在外部系统或工具中直接调用此模型数据。</div>
                </el-card>
            </el-col>
            <!-- 右侧：结果展示区域 -->
            <el-col :span="16">
                <el-card class="result-card" shadow="hover">
                    <template #header>
                        <div class="card-header">
                            <div class="tab-triggers">
                                <span class="tab-trigger" :class="{ active: resultTab === 'table' }" @click="resultTab = 'table'">
                                    表格视图
                                </span>
                                <span class="tab-trigger" :class="{ active: resultTab === 'json' }" @click="resultTab = 'json'">
                                    JSON 视图
                                </span>
                            </div>
                            <div class="result-meta" v-if="queryTime">
                                耗时: <span class="highlight">{{ queryTime }}ms</span>
                                &nbsp;&nbsp;
                                共计: <span class="highlight">{{ total }}</span> 条
                            </div>
                        </div>
                    </template>
                    <div class="result-content" v-loading="loading">
                        <template v-if="results.length">
                            <!-- 表格视图 -->
                            <el-table v-if="resultTab === 'table'" :data="results" border stripe height="650" style="width: 100%">
                                <el-table-column v-for="col in resultColumns" :key="col" :prop="col" :label="col" min-width="150" show-overflow-tooltip />
                            </el-table>
                            <!-- JSON 视图 -->
                            <div v-else class="json-preview">
                                <pre><code>{{ JSON.stringify(results, null, 2) }}</code></pre>
                            </div>
                        </template>
                        <el-empty v-else description="暂无查询结果，请点击执行查询" />
                    </div>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>
<script setup lang="ts">
import { getAllModels, getModelParams, queryDataByCode, queryDataById } from '@/api/model'
import type { Model } from '@/types/metadata'
import { ArrowLeft, Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const models = ref<Model[]>([])
const selectedModelId = ref('')
const selectedModel = ref<Model | null>(null)
const queryMode = ref('code') // 'id' or 'code'
const resultTab = ref('table')
const loading = ref(false)
const results = ref<any[]>([])
const total = ref(0)
const queryTime = ref(0)
const queryProgress = ref(0)

const testForm = reactive({
    params: {} as Record<string, any>
})

const queryParams = ref<any[]>([])

const urlPreview = computed(() => {
    if (!selectedModel.value) return ''
    const base = '/api/metadata/models/query'
    if (queryMode.value === 'id') {
        return `${base}/by-id/${selectedModel.value.id}`
    }
    const code = selectedModel.value.model_code || selectedModel.value.modelCode
    return `${base}/by-code/${code}`
})

const resultColumns = computed(() => {
    if (!results.value.length) return []
    return Object.keys(results.value[0])
})

const initModels = async () => {
    try {
        const res: any = await getAllModels()
        models.value = Array.isArray(res) ? res : (res.data || [])
    } catch (err) {
        ElMessage.error('获取模型列表失败')
    }
}

onMounted(() => {
    initModels()
})

const handleModelChange = async (id: string) => {
    const model = models.value.find(m => m.id === id)
    if (model) {
        selectedModel.value = model
        testForm.params = {}

        // 从 API 获取参数列表
        try {
            const params: any = await getModelParams(id)
            queryParams.value = Array.isArray(params) ? params : (params?.data || [])

            // 初始化默认值
            queryParams.value.forEach(p => {
                testForm.params[p.name] = p.default || ''
            })
        } catch (err) {
            console.error('获取模型参数失败:', err)
            queryParams.value = []
        }
    }
}

const performQuery = async () => {
    if (!selectedModel.value) return

    loading.value = true
    queryProgress.value = 0
    results.value = []
    const startTime = Date.now()

    // 模拟进度
    const progressInterval = setInterval(() => {
        if (queryProgress.value < 90) {
            queryProgress.value += Math.floor(Math.random() * 10)
        }
    }, 200)

    try {
        let res: any
        if (queryMode.value === 'id') {
            res = await queryDataById(selectedModel.value.id!, testForm.params)
        } else {
            const code = selectedModel.value.model_code || selectedModel.value.modelCode
            res = await queryDataByCode(code!, testForm.params)
        }

        clearInterval(progressInterval)
        queryProgress.value = 100

        const data = res.data || res
        results.value = data.list || (Array.isArray(data) ? data : [])
        total.value = data.total || results.value.length
        queryTime.value = Date.now() - startTime

        if (results.value.length === 0) {
            ElMessage.info('查询完成，但未找到匹配数据')
        }
    } catch (err: any) {
        clearInterval(progressInterval)
        queryProgress.value = 0
        ElMessage.error(err.message || '查询失败')
    } finally {
        loading.value = false
    }
}

const goBack = () => {
    router.push('/metadata/model/list')
}
</script>
<style scoped>
.sql-test-container {
    padding: 20px;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.header-left {
    display: flex;
    align-items: center;
    gap: 15px;
}

.title {
    margin: 0;
    font-size: 24px;
    font-weight: 600;
    color: #1f2937;
}

.config-card,
.result-card {
    height: 100%;
}

.w-full {
    width: 100%;
}

.parameters-section {
    margin-top: 20px;
    border-top: 1px solid #f3f4f6;
    padding-top: 15px;
}

.parameters-section h3 {
    font-size: 16px;
    margin-bottom: 15px;
    color: #4b5563;
}

.action-bar {
    margin-top: 30px;
}

.section-title {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
}

.no-params-hint {
    padding: 15px;
    background-color: #f9fafb;
    border-radius: 4px;
    color: #9ca3af;
    text-align: center;
    font-size: 13px;
    border: 1px dashed #e5e7eb;
}

.progress-info {
    margin-top: 20px;
}

.progress-text {
    margin-top: 8px;
    font-size: 12px;
    color: #6b7280;
    text-align: center;
}

.path-card {
    margin-top: 20px;
    background-color: #f9fafb;
}

.path-url {
    font-family: monospace;
    background: #1f2937;
    color: #d1d5db;
    padding: 10px;
    border-radius: 4px;
    word-break: break-all;
    margin: 10px 0;
    font-size: 13px;
}

.path-hint {
    font-size: 12px;
    color: #9ca3af;
}

.tab-triggers {
    display: flex;
    gap: 20px;
}

.tab-trigger {
    cursor: pointer;
    color: #9ca3af;
    font-weight: 500;
    padding: 5px 0;
    border-bottom: 2px solid transparent;
    transition: all 0.3s;
}

.tab-trigger.active {
    color: #4f46e5;
    border-bottom-color: #4f46e5;
}

.result-meta {
    font-size: 13px;
    color: #6b7280;
}

.highlight {
    color: #4f46e5;
    font-weight: 600;
}

.result-content {
    min-height: 400px;
}

.json-preview {
    background: #f3f4f6;
    padding: 15px;
    border-radius: 8px;
    max-height: 650px;
    overflow: auto;
}

.json-preview pre {
    margin: 0;
    font-family: 'Fira Code', 'Roboto Mono', monospace;
    font-size: 14px;
    line-height: 1.5;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
</style>
