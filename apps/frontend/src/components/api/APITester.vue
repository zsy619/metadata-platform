<template>
    <div class="api-tester">
        <div class="tester-header">
            <h3>接口测试</h3>
            <el-button text @click="visible = false">
                <el-icon><Close /></el-icon>
            </el-button>
        </div>
        <div class="tester-content">
            <el-tabs v-model="activeTab" class="tester-tabs">
                <el-tab-pane label="请求配置" name="request">
                    <div class="request-config">
                        <el-form label-width="80px" size="small">
                            <el-form-item label="请求方法">
                                <el-select v-model="requestConfig.method" style="width: 120px">
                                    <el-option v-for="m in methods" :key="m" :label="m" :value="m" />
                                </el-select>
                            </el-form-item>
                            <el-form-item label="请求地址">
                                <el-input v-model="requestConfig.url" placeholder="请输入请求地址" />
                            </el-form-item>
                            <el-form-item label="请求头">
                                <div class="headers-editor">
                                    <div v-for="(header, idx) in requestConfig.headers" :key="idx" class="header-row">
                                        <el-input v-model="header.key" placeholder="Key" style="width: 40%" />
                                        <el-input v-model="header.value" placeholder="Value" style="width: 40%" />
                                        <el-button link type="danger" @click="removeHeader(idx)">
                                            <el-icon><Delete /></el-icon>
                                        </el-button>
                                    </div>
                                    <el-button link type="primary" size="small" @click="addHeader">+ 添加请求头</el-button>
                                </div>
                            </el-form-item>
                            <el-form-item v-if="showBody" label="请求体">
                                <el-input v-model="requestConfig.body" type="textarea" :rows="6" placeholder="JSON格式请求体" />
                            </el-form-item>
                        </el-form>
                        <div class="request-actions">
                            <el-button type="primary" :loading="sending" @click="sendRequest">
                                <el-icon><Promotion /></el-icon>
                                发送请求
                            </el-button>
                            <el-button @click="resetRequest">
                                <el-icon><RefreshLeft /></el-icon>
                                重置
                            </el-button>
                        </div>
                    </div>
                </el-tab-pane>
                <el-tab-pane label="历史记录" name="history">
                    <div class="history-list">
                        <div v-if="historyList.length === 0" class="empty-history">
                            <el-empty description="暂无历史记录" />
                        </div>
                        <div v-else class="history-items">
                            <div v-for="(item, idx) in historyList" :key="idx" class="history-item" @click="loadFromHistory(item)">
                                <div class="history-method">
                                    <el-tag size="small" :type="getMethodTagType(item.method)">{{ item.method }}</el-tag>
                                </div>
                                <div class="history-url">{{ item.url }}</div>
                                <div class="history-time">{{ item.time }}</div>
                                <div class="history-status">
                                    <el-tag v-if="item.status < 400" size="small" type="success">{{ item.status }}</el-tag>
                                    <el-tag v-else size="small" type="danger">{{ item.status }}</el-tag>
                                </div>
                            </div>
                        </div>
                    </div>
                </el-tab-pane>
            </el-tabs>
            <div class="response-panel">
                <div class="response-header">
                    <span class="response-title">响应结果</span>
                    <div class="response-info" v-if="responseData">
                        <el-tag :type="responseStatus < 400 ? 'success' : 'danger'" size="small">
                            {{ responseStatus }}
                        </el-tag>
                        <el-tag type="info" size="small">{{ responseTime }}ms</el-tag>
                        <el-tag type="info" size="small">{{ responseSize }}B</el-tag>
                    </div>
                </div>
                <div class="response-body">
                    <div v-if="loading" class="response-loading">
                        <el-icon class="is-loading"><Loading /></el-icon>
                        <span>发送请求中...</span>
                    </div>
                    <div v-else-if="errorMessage" class="response-error">
                        <el-icon><CircleCloseFilled /></el-icon>
                        <span>{{ errorMessage }}</span>
                    </div>
                    <pre v-else-if="responseData" class="response-content">{{ formatJson(responseData) }}</pre>
                    <el-empty v-else description="点击发送请求查看响应" />
                </div>
                <div class="response-actions">
                    <el-button v-if="responseData" link type="primary" size="small" @click="copyResponse">
                        <el-icon><DocumentCopy /></el-icon>
                        复制
                    </el-button>
                    <el-button v-if="responseData" link type="primary" size="small" @click="formatResponse">
                        <el-icon><MagicStick /></el-icon>
                        格式化
                    </el-button>
                </div>
            </div>
        </div>
    </div>
</template>
<script setup lang="ts">
import { Close, Delete, DocumentCopy, MagicStick, Promotion, RefreshLeft } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { computed, ref, watch } from 'vue'

interface HistoryItem {
    method: string
    url: string
    headers: { key: string; value: string }[]
    body: string
    status: number
    time: string
}

const props = defineProps<{
    modelValue: boolean
    api?: {
        id?: number
        name?: string
        path?: string
        method?: string
    }
}>()

const emit = defineEmits(['update:modelValue'])

const visible = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
})

const activeTab = ref('request')
const methods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH']
const sending = ref(false)
const loading = ref(false)
const errorMessage = ref('')
const responseData = ref<any>(null)
const responseStatus = ref(0)
const responseTime = ref(0)
const responseSize = ref(0)

const requestConfig = ref({
    method: 'GET',
    url: '',
    headers: [
        { key: 'Content-Type', value: 'application/json' }
    ],
    body: ''
})

const historyList = ref<HistoryItem[]>([
    { method: 'GET', url: '/api/data/user/list', headers: [], body: '', status: 200, time: '10:30:00' },
    { method: 'POST', url: '/api/data/user', headers: [], body: '{}', status: 201, time: '10:28:00' }
])

const showBody = computed(() => ['POST', 'PUT', 'PATCH'].includes(requestConfig.value.method))

watch(() => props.api, (api) => {
    if (api && api.path && api.method) {
        requestConfig.value.url = api.path
        requestConfig.value.method = api.method
    }
}, { immediate: true })

const addHeader = () => {
    requestConfig.value.headers.push({ key: '', value: '' })
}

const removeHeader = (idx: number) => {
    requestConfig.value.headers.splice(idx, 1)
}

const sendRequest = async () => {
    if (!requestConfig.value.url) {
        ElMessage.warning('请输入请求地址')
        return
    }
    loading.value = true
    errorMessage.value = ''
    responseData.value = null
    const startTime = Date.now()
    try {
        await new Promise(resolve => setTimeout(resolve, 500 + Math.random() * 1000))
        responseStatus.value = 200
        responseData.value = {
            code: 0,
            message: 'success',
            data: {
                list: [
                    { id: 1, name: '张三', age: 25 },
                    { id: 2, name: '李四', age: 30 }
                ],
                total: 2
            }
        }
        responseTime.value = Date.now() - startTime
        responseSize.value = JSON.stringify(responseData.value).length
        const historyItem: HistoryItem = {
            method: requestConfig.value.method,
            url: requestConfig.value.url,
            headers: [...requestConfig.value.headers],
            body: requestConfig.value.body,
            status: responseStatus.value,
            time: new Date().toLocaleTimeString()
        }
        historyList.value.unshift(historyItem)
    } catch (error: any) {
        errorMessage.value = error.message || '请求失败'
        responseStatus.value = 500
    } finally {
        loading.value = false
    }
}

const resetRequest = () => {
    requestConfig.value = {
        method: 'GET',
        url: '',
        headers: [{ key: 'Content-Type', value: 'application/json' }],
        body: ''
    }
    responseData.value = null
    errorMessage.value = ''
}

const loadFromHistory = (item: HistoryItem) => {
    requestConfig.value.method = item.method
    requestConfig.value.url = item.url
    requestConfig.value.headers = item.headers.length > 0 ? [...item.headers] : [{ key: 'Content-Type', value: 'application/json' }]
    requestConfig.value.body = item.body
    activeTab.value = 'request'
}

const formatJson = (data: any) => {
    try {
        return JSON.stringify(data, null, 2)
    } catch {
        return data
    }
}

const copyResponse = () => {
    if (responseData.value) {
        navigator.clipboard.writeText(formatJson(responseData.value))
        ElMessage.success('已复制到剪贴板')
    }
}

const formatResponse = () => {
    if (responseData.value) {
        responseData.value = JSON.parse(JSON.stringify(responseData.value))
    }
}

const getMethodTagType = (method: string) => {
    const map: Record<string, string> = { GET: 'success', POST: 'primary', PUT: 'warning', DELETE: 'danger', PATCH: 'info' }
    return map[method] || 'info'
}
</script>
<style scoped>
.api-tester {
    height: 100%;
    display: flex;
    flex-direction: column;
    background: #fff;
    border-radius: 8px;
}

.tester-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    border-bottom: 1px solid #ebeef5;
}

.tester-header h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
}

.tester-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

.tester-tabs {
    flex: 1;
    display: flex;
    flex-direction: column;
}

:deep(.el-tabs__content) {
    flex: 1;
    overflow: auto;
}

.request-config {
    padding: 10px;
}

.headers-editor {
    width: 100%;
}

.header-row {
    display: flex;
    gap: 8px;
    margin-bottom: 8px;
    align-items: center;
}

.request-actions {
    display: flex;
    gap: 10px;
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid #ebeef5;
}

.history-list {
    padding: 10px;
}

.history-items {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.history-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px;
    background: #fafafa;
    border-radius: 6px;
    cursor: pointer;
    transition: background 0.2s;
}

.history-item:hover {
    background: #f0f0f0;
}

.history-url {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 13px;
}

.history-time {
    color: #909399;
    font-size: 12px;
}

.response-panel {
    border-top: 1px solid #ebeef5;
    display: flex;
    flex-direction: column;
    min-height: 200px;
}

.response-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 16px;
    background: #fafafa;
    border-bottom: 1px solid #ebeef5;
}

.response-title {
    font-weight: 600;
    font-size: 14px;
}

.response-info {
    display: flex;
    gap: 8px;
}

.response-body {
    flex: 1;
    padding: 16px;
    overflow: auto;
    background: #fafafa;
}

.response-loading,
.response-error {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    height: 100%;
    color: #909399;
}

.response-error {
    color: #f56c6c;
}

.response-content {
    margin: 0;
    font-family: 'Monaco', 'Menlo', monospace;
    font-size: 12px;
    line-height: 1.5;
    white-space: pre-wrap;
    word-break: break-all;
    color: #303133;
}

.response-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    padding: 8px 16px;
    border-top: 1px solid #ebeef5;
}
</style>
