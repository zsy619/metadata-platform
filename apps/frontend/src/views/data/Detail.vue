<template>
    <div class="data-detail container-padding">
        <div class="page-header">
            <div class="header-left">
                <el-button :icon="ArrowLeft" @click="handleBack">返回</el-button>
                <h1 class="page-title">
                    <el-icon class="title-icon">
                        <Document />
                    </el-icon>
                    数据详情
                </h1>
            </div>
            <div class="header-actions">
                <el-button type="primary" :icon="Edit" @click="handleEdit">编辑</el-button>
                <el-button type="danger" :icon="Delete" @click="handleDelete">删除</el-button>
            </div>
        </div>
        <div class="detail-content">
            <el-card class="info-card">
                <template #header>
                    <div class="card-header">
                        <span>基本信息</span>
                    </div>
                </template>
                <div class="info-grid">
                    <div v-for="field in displayFields" :key="field.prop" class="info-item">
                        <div class="info-label">{{ field.label }}</div>
                        <div class="info-value">
                            <template v-if="field.type === 'image'">
                                <el-image v-if="detailData[field.prop]" :src="detailData[field.prop]" fit="cover" style="width: 80px; height: 80px" preview-src-list [src]="detailData[field.prop]" />
                                <span v-else class="empty-value">-</span>
                            </template>
                            <template v-else-if="field.type === 'switch'">
                                <el-tag :type="detailData[field.prop] ? 'success' : 'info'">
                                    {{ detailData[field.prop] ? '启用' : '禁用' }}
                                </el-tag>
                            </template>
                            <template v-else-if="field.type === 'tag'">
                                <el-tag :type="getTagType(detailData[field.prop], field.tagTypes)">
                                    {{ getTagText(detailData[field.prop], field.tagOptions) }}
                                </el-tag>
                            </template>
                            <template v-else-if="field.type === 'datetime'">
                                {{ formatDate(detailData[field.prop]) }}
                            </template>
                            <template v-else-if="field.type === 'json'">
                                <pre class="json-value">{{ formatJson(detailData[field.prop]) }}</pre>
                            </template>
                            <template v-else>
                                <span :class="{ 'empty-value': !detailData[field.prop] }">
                                    {{ detailData[field.prop] ?? '-' }}
                                </span>
                            </template>
                        </div>
                    </div>
                </div>
            </el-card>
            <el-card class="timeline-card">
                <template #header>
                    <div class="card-header">
                        <span>操作日志</span>
                    </div>
                </template>
                <el-timeline>
                    <el-timeline-item v-for="(log, index) in operationLogs" :key="index" :timestamp="log.timestamp" :type="log.type" :icon="log.icon">
                        <div class="timeline-content">
                            <div class="timeline-title">{{ log.title }}</div>
                            <div class="timeline-desc">{{ log.description }}</div>
                            <div class="timeline-user">{{ log.user }}</div>
                        </div>
                    </el-timeline-item>
                </el-timeline>
                <el-empty v-if="operationLogs.length === 0" description="暂无操作日志" />
            </el-card>
        </div>
    </div>
</template>
<script setup lang="ts">
import { ArrowLeft, Delete, Document, Edit } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

interface Field {
    prop: string
    label: string
    type?: string
    tagTypes?: Record<string, string>
    tagOptions?: Record<string, string>
}

interface OperationLog {
    timestamp: string
    title: string
    description: string
    user: string
    type?: string
    icon?: string
}

const route = useRoute()
const router = useRouter()

const modelId = computed(() => Number(route.params.modelId) || 1)
const dataId = computed(() => Number(route.params.id) || 1)

const displayFields = ref<Field[]>([
    { prop: 'id', label: 'ID' },
    { prop: 'username', label: '用户名' },
    { prop: 'email', label: '邮箱' },
    { prop: 'phone', label: '手机号' },
    { prop: 'status', label: '状态', type: 'switch' },
    { prop: 'createTime', label: '创建时间', type: 'datetime' },
    { prop: 'updateTime', label: '更新时间', type: 'datetime' }
])

const detailData = ref<any>({})

const operationLogs = ref<OperationLog[]>([
    { timestamp: '2024-01-23 10:00:00', title: '创建记录', description: '新增用户数据', user: '系统', type: 'success' },
    { timestamp: '2024-01-23 14:30:00', title: '更新记录', description: '修改用户邮箱', user: 'admin', type: 'primary' },
    { timestamp: '2024-01-24 09:15:00', title: '状态变更', description: '启用用户账号', user: 'admin', type: 'warning' }
])

onMounted(() => {
    loadDetail()
})

const loadDetail = async () => {
    detailData.value = {
        id: dataId.value,
        username: 'user1',
        email: 'user1@example.com',
        phone: '13800138000',
        status: 1,
        createTime: '2024-01-23 10:00:00',
        updateTime: '2024-01-24 09:15:00'
    }
}

const handleBack = () => {
    router.back()
}

const handleEdit = () => {
    router.push(`/data/${modelId.value}/edit/${dataId.value}`)
}

const handleDelete = () => {
    ElMessageBox.confirm('确定要删除这条数据吗？', '提示', { type: 'warning' }).then(async () => {
        ElMessage.success('删除成功')
        router.push(`/data/${modelId.value}/list`)
    })
}

const getTagType = (value: any, tagTypes?: Record<string, string>) => tagTypes?.[value] || 'info'
const getTagText = (value: any, tagOptions?: Record<string, string>) => tagOptions?.[value] ?? value
const formatDate = (value: string) => value ? new Date(value).toLocaleString() : '-'
const formatJson = (value: any) => {
    try {
        return typeof value === 'string' ? JSON.stringify(JSON.parse(value), null, 2) : JSON.stringify(value, null, 2)
    } catch {
        return value
    }
}
</script>
<style scoped>
.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.header-left {
    display: flex;
    align-items: center;
    gap: 16px;
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

.detail-content {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.info-card,
.timeline-card {
    width: 100%;
}

.card-header {
    font-weight: 600;
    font-size: 16px;
}

.info-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 20px;
}

.info-item {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.info-label {
    color: #909399;
    font-size: 14px;
}

.info-value {
    color: #303133;
    font-size: 14px;
    font-weight: 500;
    word-break: break-all;
}

.empty-value {
    color: #c0c4cc;
}

.json-value {
    background: #f5f7fa;
    padding: 12px;
    border-radius: 6px;
    font-size: 12px;
    font-family: 'Monaco', 'Menlo', monospace;
    white-space: pre-wrap;
    word-break: break-all;
    max-height: 200px;
    overflow: auto;
}

.timeline-content {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.timeline-title {
    font-weight: 500;
    color: #303133;
}

.timeline-desc {
    color: #606266;
    font-size: 13px;
}

.timeline-user {
    color: #909399;
    font-size: 12px;
}
</style>
