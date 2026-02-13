<template>
    <div class="build-preview">
        <div class="preview-header">
            <el-icon size="20"><Document /></el-icon>
            <span>配置预览</span>
        </div>
        <el-scrollbar class="preview-content">
            <div class="preview-section">
                <div class="section-title">
                    <el-icon><InfoFilled /></el-icon>
                    <span>基本信息</span>
                </div>
                <div class="info-grid">
                    <div class="info-item">
                        <span class="label">模型名称:</span>
                        <span class="value">{{ modelInfo.name || '-' }}</span>
                    </div>
                    <div class="info-item">
                        <span class="label">模型编码:</span>
                        <span class="value">{{ modelInfo.code || '-' }}</span>
                    </div>
                    <div class="info-item">
                        <span class="label">数据源:</span>
                        <span class="value">{{ modelInfo.dataSourceName || '-' }}</span>
                    </div>
                    <div class="info-item">
                        <span class="label">模型类型:</span>
                        <span class="value">{{ getModelTypeText(modelInfo.type) }}</span>
                    </div>
                    <div class="info-item">
                        <span class="label">版本:</span>
                        <span class="value">v{{ modelInfo.version || '1.0' }}</span>
                    </div>
                    <div class="info-item">
                        <span class="label">备注:</span>
                        <span class="value">{{ modelInfo.remark || '-' }}</span>
                    </div>
                </div>
            </div>
            <div class="preview-section">
                <div class="section-title">
                    <el-icon><Grid /></el-icon>
                    <span>字段列表 ({{ fields.length }} 个字段)</span>
                </div>
                <el-table :data="fields" border size="small" max-height="250">
                    <el-table-column type="index" label="序号" width="50" align="center" />
                    <el-table-column label="字段名" prop="name" min-width="100" show-overflow-tooltip />
                    <el-table-column label="显示名" prop="displayTitle" min-width="100" show-overflow-tooltip />
                    <el-table-column label="类型" prop="type" width="80" />
                    <el-table-column label="组件" prop="componentType" width="80" />
                    <el-table-column label="主键" width="50" align="center">
                        <template #default="{ row }">
                            <el-icon v-if="row.isPrimaryKey" color="#67c23a"><CircleCheckFilled /></el-icon>
                            <span v-else>-</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="必填" width="50" align="center">
                        <template #default="{ row }">
                            <el-icon v-if="row.isRequired" color="#f56c6c"><CircleCheckFilled /></el-icon>
                            <span v-else>-</span>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
            <div v-if="validationRules && validationRules.length > 0" class="preview-section">
                <div class="section-title">
                    <el-icon><CircleCheck /></el-icon>
                    <span>验证规则 ({{ getValidationCount() }} 条)</span>
                </div>
                <div class="validation-list">
                    <div v-for="(rule, idx) in validationRules" :key="idx" class="validation-item">
                        <span class="field-name">{{ rule.name }}</span>
                        <el-tag v-for="v in rule.validations" :key="v" size="small" type="info" style="margin-right: 4px">
                            {{ getValidationText(v) }}
                        </el-tag>
                    </div>
                </div>
            </div>
            <div class="preview-section">
                <div class="section-title">
                    <el-icon><Monitor /></el-icon>
                    <span>显示配置</span>
                </div>
                <div class="display-summary">
                    <div class="summary-item">
                        <span class="label">可排序字段:</span>
                        <span class="value">{{ sortableCount }} 个</span>
                    </div>
                    <div class="summary-item">
                        <span class="label">可筛选字段:</span>
                        <span class="value">{{ filterableCount }} 个</span>
                    </div>
                    <div class="summary-item">
                        <span class="label">可搜索字段:</span>
                        <span class="value">{{ searchableCount }} 个</span>
                    </div>
                </div>
            </div>
            <div class="preview-section">
                <div class="section-title">
                    <el-icon><Connection /></el-icon>
                    <span>关联配置 ({{ joins.length }} 个关联)</span>
                </div>
                <div v-if="joins.length === 0" class="empty-text">暂无关联配置</div>
                <div v-else class="join-list">
                    <div v-for="(join, idx) in joins" :key="idx" class="join-item">
                        <span class="join-type">{{ getJoinTypeText(join.type) }}</span>
                        <span class="join-from">{{ join.fromTable }}.{{ join.fromField }}</span>
                        <span class="join-arrow">→</span>
                        <span class="join-to">{{ join.toTable }}.{{ join.toField }}</span>
                    </div>
                </div>
            </div>
        </el-scrollbar>
        <div class="preview-footer">
            <el-checkbox v-model="confirmed" class="confirm-check">
                我已确认以上配置正确
            </el-checkbox>
            <div class="footer-actions">
                <el-button @click="$emit('prev')">上一步</el-button>
                <el-button type="primary" :disabled="!confirmed" :loading="loading" @click="handleCreate">
                    {{ loading ? '创建中...' : '确认创建' }}
                </el-button>
            </div>
        </div>
    </div>
</template>
<script setup lang="ts">
import { CircleCheck, CircleCheckFilled, Connection, Document, Grid, InfoFilled, Monitor } from '@element-plus/icons-vue'
import { computed, ref } from 'vue'

interface ModelInfo {
    name: string
    code: string
    dataSourceId?: number
    dataSourceName?: string
    type?: number
    version?: string
    remark?: string
}

interface FieldInfo {
    name: string
    displayTitle: string
    type: string
    componentType: string
    isPrimaryKey: boolean
    isRequired: boolean
}

interface ValidationRule {
    name: string
    validations: string[]
}

interface JoinInfo {
    type: string
    fromTable: string
    fromField: string
    toTable: string
    toField: string
}

interface DisplayConfig {
    sortable: boolean
    filterable: boolean
    searchable: boolean
}

const props = defineProps<{
    modelInfo: ModelInfo
    fields: FieldInfo[]
    validationRules?: ValidationRule[]
    displayConfig?: DisplayConfig[]
    joins?: JoinInfo[]
    loading?: boolean
}>()

const emit = defineEmits(['prev', 'confirm'])

const confirmed = ref(false)

const sortableCount = computed(() => props.displayConfig?.filter(d => d.sortable).length || 0)
const filterableCount = computed(() => props.displayConfig?.filter(d => d.filterable).length || 0)
const searchableCount = computed(() => props.displayConfig?.filter(d => d.searchable).length || 0)

const getModelTypeText = (type?: number) => {
    const typeMap: Record<number, string> = {
        1: 'SQL模型',
        2: '视图模型',
        3: '存储过程',
        4: '关联模型'
    }
    return typeMap[type || 1] || 'SQL模型'
}

const getValidationText = (v: string) => {
    const textMap: Record<string, string> = {
        required: '必填',
        email: '邮箱',
        mobile: '手机号',
        idCard: '身份证',
        url: 'URL',
        number: '数字',
        integer: '整数',
        date: '日期',
        json: 'JSON',
        alpha: '字母',
        alphanumeric: '字母数字'
    }
    return textMap[v] || v
}

const getValidationCount = () => {
    return props.validationRules?.reduce((acc, rule) => acc + (rule.validations?.length || 0), 0) || 0
}

const getJoinTypeText = (type?: string) => {
    const typeMap: Record<string, string> = {
        inner: 'INNER JOIN',
        left: 'LEFT JOIN',
        right: 'RIGHT JOIN',
        full: 'FULL JOIN'
    }
    return typeMap[type || 'left'] || 'LEFT JOIN'
}

const handleCreate = () => {
    if (!confirmed.value) return
    emit('confirm')
}
</script>
<style scoped>
.build-preview {
    display: flex;
    flex-direction: column;
    height: 100%;
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.preview-header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 16px;
    border-bottom: 1px solid #ebeef5;
    font-weight: 600;
    color: #303133;
}

.preview-content {
    flex: 1;
    padding: 16px;
    overflow: hidden;
}

.preview-section {
    margin-bottom: 20px;
    padding: 16px;
    background: #fafafa;
    border-radius: 8px;
}

.preview-section:last-child {
    margin-bottom: 0;
}

.section-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 12px;
    padding-bottom: 8px;
    border-bottom: 1px solid #ebeef5;
}

.info-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
}

.info-item {
    display: flex;
    gap: 8px;
}

.info-item .label {
    color: #909399;
    white-space: nowrap;
}

.info-item .value {
    color: #303133;
    font-weight: 500;
}

.validation-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.validation-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px;
    background: #fff;
    border-radius: 4px;
}

.field-name {
    font-weight: 500;
    min-width: 100px;
}

.display-summary {
    display: flex;
    gap: 20px;
}

.summary-item {
    display: flex;
    align-items: center;
    gap: 8px;
}

.summary-item .label {
    color: #909399;
}

.summary-item .value {
    font-weight: 500;
    color: #409eff;
}

.empty-text {
    color: #909399;
    font-size: 13px;
    padding: 12px 0;
}

.join-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.join-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px;
    background: #fff;
    border-radius: 4px;
    font-size: 13px;
}

.join-type {
    padding: 2px 6px;
    background: #409eff;
    color: #fff;
    border-radius: 3px;
    font-size: 11px;
}

.join-from,
.join-to {
    color: #303133;
}

.join-arrow {
    color: #909399;
}

.preview-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px;
    border-top: 1px solid #ebeef5;
}

.confirm-check {
    margin: 0;
}

.footer-actions {
    display: flex;
    gap: 10px;
}
</style>
