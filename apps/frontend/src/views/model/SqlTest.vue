<template>
    <div class="sql-test-container container-padding">
        <div class="page-header">
            <div class="header-left">
                <el-icon :size="24" class="page-icon">
                    <Monitor />
                </el-icon>
                <h1 class="title">SQL 模型测试</h1>
            </div>
            <el-tag type="info" effect="plain">验证动态查询接口与参数配置</el-tag>
        </div>
        <div class="content-wrapper" ref="wrapperRef">
            <div class="content-body">
                <!-- 左侧：配置区域 -->
                <div class="config-panel" :style="{ width: leftWidth + 'px' }">
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
                                    <el-radio-button value="id">按模型 ID</el-radio-button>
                                    <el-radio-button value="code">按模型代码</el-radio-button>
                                </el-radio-group>
                            </el-form-item>
                            <div v-if="selectedModel" class="parameters-section">
                                <div class="section-title">
                                    <h3>查询参数</h3>
                                    <el-tag v-if="!queryParams.length" type="info" size="small">无</el-tag>
                                </div>
                                <el-table v-if="queryParams.length" :data="queryParams" border stripe size="small" style="width: 100%; margin-top: 10px">
                                    <el-table-column prop="name" label="参数" width="100" />
                                    <el-table-column prop="type" label="类型" width="60" />
                                    <el-table-column label="值">
                                        <template #default="{ row }">
                                            <el-input v-model="testForm.params[row.name]" :placeholder="`输入${row.name}`" size="small" />
                                        </template>
                                    </el-table-column>
                                </el-table>
                                <div v-else class="no-params-hint">该模型未定义动态参数</div>
                            </div>
                            <div v-if="selectedModel" class="path-preview-section">
                                <div class="section-title">
                                    <h3>请求路径预览</h3>
                                    <el-tag size="small" type="success">POST</el-tag>
                                </div>
                                <div class="path-url">{{ urlPreview }}</div>
                                <div class="path-hint">您可以使用该路径在外部系统或工具中直接调用此模型数据。</div>
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
                </div>

                <!-- 拖拽手柄 -->
                <div class="resize-handle" @mousedown="startResize"></div>

                <!-- 右侧：结果展示区域 -->
                <div class="result-panel">
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
                                <div class="header-right">
                                    <div class="result-meta" v-if="queryTime">
                                        耗时: <span class="highlight">{{ queryTime }}ms</span>
                                        &nbsp;&nbsp;
                                        共计: <span class="highlight">{{ total }}</span> 条
                                    </div>
                                    <el-divider direction="vertical" v-if="queryTime" />
                                    <el-button type="primary" link :icon="Download" @click="exportToExcel" :disabled="!results.length">
                                        导出 Excel
                                    </el-button>
                                </div>
                            </div>
                        </template>
                        <div class="result-content" v-loading="loading">
                            <template v-if="results.length">
                                <!-- 表格视图 -->
                                <el-table v-if="resultTab === 'table'" :data="results" border stripe height="100%" style="width: 100%">
                                    <el-table-column v-for="col in resultColumns" :key="col" :prop="col" :label="getColumnLabel(col)" :width="getColumnWidth(col)" show-overflow-tooltip />
                                </el-table>
                                <!-- JSON 视图 -->
                                <div v-else class="json-preview">
                                    <pre><code>{{ JSON.stringify(results, null, 2) }}</code></pre>
                                </div>
                            </template>
                            <el-empty v-else description="暂无查询结果，请点击执行查询" />
                        </div>
                    </el-card>
                </div>
            </div>
        </div>
    </div>
</template>
<script setup lang="ts">
import { getAllModels, getModelFields, getModelParams, queryDataByCode, queryDataById } from '@/api/model'
import { useAppStore } from '@/stores/app'
import type { Model } from '@/types/metadata'
import { Download, Monitor, Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'

const appStore = useAppStore()

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
const leftWidth = ref(400) // 初始左侧宽度
const isResizing = ref(false)

const startResize = (e: MouseEvent) => {
    isResizing.value = true
    document.addEventListener('mousemove', doResize)
    document.addEventListener('mouseup', stopResize)
    document.body.style.cursor = 'col-resize'
    document.body.style.userSelect = 'none'
}

const doResize = (e: MouseEvent) => {
    if (!isResizing.value) return
    const sidebarWidth = appStore.sidebar.opened ? appStore.sidebar.width : 64
    const newWidth = e.clientX - sidebarWidth - 20 // 20 为 AppMain 的 padding
    if (newWidth > 300 && newWidth < 800) {
        leftWidth.value = newWidth
    }
}

const stopResize = () => {
    isResizing.value = false
    document.removeEventListener('mousemove', doResize)
    document.removeEventListener('mouseup', stopResize)
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
}

const testForm = reactive({
    params: {} as Record<string, any>
})

const queryParams = ref<any[]>([])
const modelFields = ref<any[]>([])

const getColumnLabel = (columnName: string) => {
    const field = modelFields.value.find(f => f.column_name === columnName || f.columnName === columnName)
    return field?.show_title || field?.showTitle || field?.column_title || field?.columnTitle || columnName
}

const getColumnWidth = (columnName: string) => {
    const field = modelFields.value.find(f => f.column_name === columnName || f.columnName === columnName)
    return field?.show_width || field?.showWidth || 150
}

const exportToExcel = () => {
    if (!results.value.length) {
        ElMessage.warning('暂无数据可导出')
        return
    }

    // 1. 准备数据
    const columns = resultColumns.value

    // 2. 构建 XML Spreadsheet 头部
    let xml = '<?xml version="1.0"?>\n' +
        '<?mso-application progid="Excel.Sheet"?>\n' +
        '<Workbook xmlns="urn:schemas-microsoft-com:office:spreadsheet"\n' +
        ' xmlns:o="urn:schemas-microsoft-com:office:office"\n' +
        ' xmlns:x="urn:schemas-microsoft-com:office:excel"\n' +
        ' xmlns:ss="urn:schemas-microsoft-com:office:spreadsheet"\n' +
        ' xmlns:html="http://www.w3.org/TR/REC-html40">\n' +
        ' <Styles>\n' +
        '  <Style ss:ID="Default" ss:Name="Normal">\n' +
        '   <Alignment ss:Vertical="Center"/>\n' +
        '  </Style>\n' +
        '  <Style ss:ID="Header">\n' +
        '   <Font ss:Bold="1"/>\n' +
        '   <Interior ss:Color="#EFEFEF" ss:Pattern="Solid"/>\n' +
        '  </Style>\n' +
        ' </Styles>\n' +
        ' <Worksheet ss:Name="查询结果">\n' +
        '  <Table>\n'

    // 3. 设置列宽
    columns.forEach(col => {
        const width = getColumnWidth(col)
        xml += `   <Column ss:Width="${width || 100}"/>\n`
    })

    // 4. 添加表头
    xml += '   <Row ss:StyleID="Header">\n'
    columns.forEach(col => {
        const label = getColumnLabel(col)
        xml += `    <Cell><Data ss:Type="String">${escapeXml(label)}</Data></Cell>\n`
    })
    xml += '   </Row>\n'

    // 5. 添加数据行
    results.value.forEach(row => {
        xml += '   <Row>\n'
        columns.forEach(col => {
            let val = row[col]
            let type = 'String'
            if (val === null || val === undefined) {
                val = ''
            } else if (typeof val === 'number') {
                type = 'Number'
            } else {
                val = String(val)
            }
            xml += `    <Cell><Data ss:Type="${type}">${escapeXml(String(val))}</Data></Cell>\n`
        })
        xml += '   </Row>\n'
    })

    // 6. 结束 XML
    xml += '  </Table>\n' +
        ' </Worksheet>\n' +
        '</Workbook>'

    // 7. 触发下载 (.xls)
    const blob = new Blob([xml], { type: 'application/vnd.ms-excel' })
    const link = document.createElement('a')
    const url = URL.createObjectURL(blob)
    link.setAttribute('href', url)
    link.setAttribute('download', `查询结果_${new Date().toISOString().slice(0, 19).replace(/T|:/g, '-')}.xls`)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
}

const escapeXml = (unsafe: string) => {
    return unsafe.replace(/[<>&'"]/g, c => {
        switch (c) {
            case '<': return '&lt;'
            case '>': return '&gt;'
            case '&': return '&amp;'
            case '\'': return '&apos;'
            case '"': return '&quot;'
            default: return c
        }
    })
}

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

        // 获取字段列表以显示字段标题
        try {
            const fields: any = await getModelFields(id)
            modelFields.value = Array.isArray(fields) ? fields : (fields?.data || [])
        } catch (err) {
            console.error('获取模型字段失败:', err)
            modelFields.value = []
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
</script>
<style scoped>
.sql-test-container {
    background-color: var(--el-bg-color-page);
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;
    flex-shrink: 0;
}

.header-left,
.header-right {
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

.page-icon {
    color: #6366f1;
}

.config-card,
.result-card {
    height: 100%;
    display: flex;
    flex-direction: column;
}

.config-card :deep(.el-card__body) {
    flex: 1;
    overflow-y: auto;
}

.result-card :deep(.el-card__body) {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
}

.result-content {
    flex: 1;
    height: 0;
    /* 关键：强制让子元素适应 flex 容器 */
    display: flex;
    flex-direction: column;
}

.json-preview {
    flex: 1;
    overflow: auto;
    background-color: #f8f9fa;
    padding: 10px;
    border-radius: 4px;
}

.json-preview pre {
    margin: 0;
}

.content-wrapper {
    flex: 1;
    overflow: hidden;
    min-height: 0;
}

.content-body {
    display: flex;
    height: 100%;
    gap: 0; /* 使用 resize-handle 代替 gap */
}

.config-panel {
    flex-shrink: 0;
    min-width: 300px;
    height: 100%;
}

.result-panel {
    flex: 1;
    min-width: 400px;
    height: 100%;
    overflow: hidden;
}

.resize-handle {
    width: 8px;
    height: 100%;
    cursor: col-resize;
    background-color: transparent;
    transition: background-color 0.2s;
    flex-shrink: 0;
    position: relative;
    z-index: 10;
}

.resize-handle:hover,
.resize-handle:active {
    background-color: var(--el-color-primary-light-8);
}

.resize-handle::after {
    content: '';
    position: absolute;
    left: 4px;
    top: 50%;
    transform: translateY(-50%);
    width: 1px;
    height: 30px;
    background-color: var(--el-border-color);
}

.w-full {
    width: 100%;
}

.parameters-section {
    margin-top: 20px;
    border-top: 1px solid #f3f4f6;
    padding-top: 15px;
}

.path-preview-section {
    margin-top: 20px;
    border-top: 1px solid #f3f4f6;
    padding-top: 15px;
}

.parameters-section h3,
.path-preview-section h3 {
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
