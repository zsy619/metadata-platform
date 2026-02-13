<template>
    <div class="validation-config-panel" :class="{ 'is-collapsed': collapsed }">
        <div class="panel-header" @click="$emit('toggle-collapse')">
            <div class="header-content">
                <el-icon>
                    <CircleCheck />
                </el-icon>
                <span>验证规则配置</span>
            </div>
            <el-icon class="collapse-icon">
                <ArrowUp v-if="!collapsed" />
                <ArrowDown v-else />
            </el-icon>
        </div>
        <div v-show="!collapsed" class="panel-body">
            <div class="panel-tip">
                <el-icon><InfoFilled /></el-icon>
                <span>为每个字段配置验证规则，支持必填、长度、范围、格式等多种验证方式</span>
            </div>
            <el-tabs v-model="activeTab" class="validation-tabs">
                <el-tab-pane label="字段验证" name="field">
                    <div class="field-list">
                        <el-table :data="fieldsWithValidation" border size="small" max-height="400">
                            <el-table-column label="字段" min-width="120" show-overflow-tooltip>
                                <template #default="{ row }">
                                    <div class="field-cell">
                                        <el-icon><Grid /></el-icon>
                                        <span>{{ row.name }}</span>
                                        <el-tag v-if="row.isRequired" size="small" type="danger">必填</el-tag>
                                    </div>
                                </template>
                            </el-table-column>
                            <el-table-column label="数据类型" width="100">
                                <template #default="{ row }">
                                    <span class="type-tag">{{ row.type }}</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="验证规则" min-width="200">
                                <template #default="{ row, $index }">
                                    <el-select v-model="row.validations" multiple placeholder="选择验证规则" size="small" @change="handleValidationChange($index)">
                                        <el-option label="必填" value="required" :disabled="row.isPrimaryKey" />
                                        <el-option label="邮箱" value="email" />
                                        <el-option label="手机号" value="mobile" />
                                        <el-option label="身份证" value="idCard" />
                                        <el-option label="URL" value="url" />
                                        <el-option label="数字" value="number" />
                                        <el-option label="整数" value="integer" />
                                        <el-option label="日期" value="date" />
                                        <el-option label="JSON" value="json" />
                                        <el-option label="字母" value="alpha" />
                                        <el-option label="字母数字" value="alphanumeric" />
                                    </el-select>
                                </template>
                            </el-table-column>
                            <el-table-column label="长度限制" width="180">
                                <template #default="{ row, $index }">
                                    <div class="length-inputs" v-if="showLengthLimit(row.validations)">
                                        <el-input-number v-model="row.minLength" :min="0" :max="row.maxLength" size="small" placeholder="最小" controls-position="right" style="width: 70px" @change="handleValidationChange($index)" />
                                        <span>-</span>
                                        <el-input-number v-model="row.maxLength" :min="row.minLength || 0" size="small" placeholder="最大" controls-position="right" style="width: 70px" @change="handleValidationChange($index)" />
                                    </div>
                                    <span v-else class="no-limit">-</span>
                                </template>
                            </el-table-column>
                            <el-table-column label="自定义错误提示" min-width="150">
                                <template #default="{ row }">
                                    <el-input v-model="row.customMessage" placeholder="自定义提示" size="small" />
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>
                </el-tab-pane>
                <el-tab-pane label="预设模板" name="template">
                    <div class="template-list">
                        <div class="template-category">
                            <div class="category-title">常用验证模板</div>
                            <div class="template-items">
                                <div class="template-item" @click="applyTemplate('email')">
                                    <el-icon><Message /></el-icon>
                                    <span>邮箱验证</span>
                                </div>
                                <div class="template-item" @click="applyTemplate('mobile')">
                                    <el-icon><Phone /></el-icon>
                                    <span>手机号验证</span>
                                </div>
                                <div class="template-item" @click="applyTemplate('idCard')">
                                    <el-icon><IdCard /></el-icon>
                                    <span>身份证验证</span>
                                </div>
                                <div class="template-item" @click="applyTemplate('url')">
                                    <el-icon><Link /></el-icon>
                                    <span>URL验证</span>
                                </div>
                                <div class="template-item" @click="applyTemplate('username')">
                                    <el-icon><User /></el-icon>
                                    <span>用户名验证</span>
                                </div>
                                <div class="template-item" @click="applyTemplate('password')">
                                    <el-icon><Lock /></el-icon>
                                    <span>密码强度验证</span>
                                </div>
                            </div>
                        </div>
                        <div class="template-category">
                            <div class="category-title">自定义模板</div>
                            <div class="template-actions">
                                <el-button type="primary" size="small" :icon="Plus" @click="showSaveTemplateDialog = true">保存当前配置为模板</el-button>
                                <el-button size="small" :icon="Upload" @click="importTemplate">导入模板</el-button>
                            </div>
                        </div>
                    </div>
                </el-tab-pane>
                <el-tab-pane label="测试" name="test">
                    <div class="test-panel">
                        <el-form label-width="100px">
                            <el-form-item label="选择字段">
                                <el-select v-model="testField" placeholder="选择字段" size="small">
                                    <el-option v-for="field in fieldsWithValidation" :key="field.name" :label="field.name" :value="field.name" />
                                </el-select>
                            </el-form-item>
                            <el-form-item label="测试值">
                                <el-input v-model="testValue" placeholder="输入测试值" size="small" />
                            </el-form-item>
                            <el-form-item>
                                <el-button type="primary" size="small" @click="runTest">验证</el-button>
                                <el-button size="small" @click="testValue = ''">重置</el-button>
                            </el-form-item>
                        </el-form>
                        <div v-if="testResult !== null" class="test-result" :class="{ success: testResult, error: !testResult }">
                            <el-icon v-if="testResult"><CircleCheckFilled /></el-icon>
                            <el-icon v-else><CircleCloseFilled /></el-icon>
                            <span>{{ testResult ? '验证通过' : '验证失败' }}</span>
                        </div>
                    </div>
                </el-tab-pane>
            </el-tabs>
        </div>
    </div>
</template>
<script setup lang="ts">
import { ArrowDown, ArrowUp, CircleCheck, CircleCheckFilled, CircleCloseFilled, Grid, IdCard, InfoFilled, Link, Lock, Message, Phone, Plus, Upload, User } from '@element-plus/icons-vue'
import { ref } from 'vue'

interface ValidationRule {
    name: string
    type: string
    isRequired: boolean
    isPrimaryKey: boolean
    validations: string[]
    minLength?: number
    maxLength?: number
    customMessage: string
}

const props = defineProps<{
    fields: any[]
    collapsed?: boolean
}>()

const emit = defineEmits(['update', 'toggle-collapse'])

const activeTab = ref('field')
const testField = ref('')
const testValue = ref('')
const testResult = ref<boolean | null>(null)
const showSaveTemplateDialog = ref(false)

const fieldsWithValidation = ref<ValidationRule[]>(
    (props.fields || []).map(f => ({
        name: f.name || f.fieldName,
        type: f.type || f.fieldType,
        isRequired: f.isRequired || false,
        isPrimaryKey: f.isPrimaryKey || false,
        validations: f.validations || [],
        minLength: f.minLength,
        maxLength: f.maxLength,
        customMessage: f.customMessage || ''
    }))
)

const showLengthLimit = (validations: string[]) => {
    return validations.some(v => ['required', 'string', 'username', 'password'].includes(v))
}

const handleValidationChange = (index: number) => {
    emit('update', fieldsWithValidation.value)
}

const applyTemplate = (template: string) => {
    const templateRules: Record<string, string[]> = {
        email: ['required', 'email'],
        mobile: ['required', 'mobile'],
        idCard: ['required', 'idCard'],
        url: ['url'],
        username: ['required', 'alphanumeric', 'minLength:4', 'maxLength:20'],
        password: ['required', 'minLength:6', 'maxLength:32']
    }
    const rules = templateRules[template] || []
    fieldsWithValidation.value.forEach(field => {
        if (template === 'username' || template === 'password') {
            field.validations = rules.filter(r => !r.includes(':'))
            rules.forEach(r => {
                if (r.includes('minLength')) {
                    field.minLength = parseInt(r.split(':')[1])
                }
                if (r.includes('maxLength')) {
                    field.maxLength = parseInt(r.split(':')[1])
                }
            })
        } else {
            field.validations = rules
        }
    })
    emit('update', fieldsWithValidation.value)
}

const importTemplate = () => {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = '.json'
    input.onchange = (e: any) => {
        const file = e.target.files[0]
        if (file) {
            const reader = new FileReader()
            reader.onload = (event) => {
                try {
                    const template = JSON.parse(event.target?.result as string)
                    fieldsWithValidation.value = template
                    emit('update', fieldsWithValidation.value)
                } catch {
                    ElMessage.error('模板格式错误')
                }
            }
            reader.readAsText(file)
        }
    }
    input.click()
}

const runTest = () => {
    if (!testField.value || !testValue.value) {
        ElMessage.warning('请选择字段并输入测试值')
        return
    }
    const field = fieldsWithValidation.value.find(f => f.name === testField.value)
    if (!field) return

    let valid = true
    const value = testValue.value

    if (field.validations.includes('required') && !value) {
        valid = false
    } else if (field.validations.includes('email') && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
        valid = false
    } else if (field.validations.includes('mobile') && !/^1[3-9]\d{9}$/.test(value)) {
        valid = false
    } else if (field.validations.includes('idCard') && !/^\d{17}[\dXx]$/.test(value)) {
        valid = false
    } else if (field.validations.includes('number') && isNaN(Number(value))) {
        valid = false
    } else if (field.validations.includes('integer') && !/^\d+$/.test(value)) {
        valid = false
    } else if (field.validations.includes('alpha') && !/^[a-zA-Z]+$/.test(value)) {
        valid = false
    } else if (field.validations.includes('alphanumeric') && !/^[a-zA-Z0-9]+$/.test(value)) {
        valid = false
    } else if (field.minLength !== undefined && value.length < field.minLength) {
        valid = false
    } else if (field.maxLength !== undefined && value.length > field.maxLength) {
        valid = false
    }

    testResult.value = valid
    if (!valid) {
        ElMessage.warning(field.customMessage || '验证失败')
    } else {
        ElMessage.success('验证通过')
    }
}

import { ElMessage } from 'element-plus'

defineExpose({
    fieldsWithValidation,
    getData: () => fieldsWithValidation.value
})
</script>
<style scoped>
.validation-config-panel {
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.panel-header {
    padding: 12px 16px;
    border-bottom: 1px solid #ebeef5;
    font-weight: 600;
    color: #303133;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
}

.header-content {
    display: flex;
    align-items: center;
    gap: 8px;
}

.panel-body {
    padding: 16px;
}

.panel-tip {
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

.field-cell {
    display: flex;
    align-items: center;
    gap: 6px;
}

.type-tag {
    font-size: 12px;
    color: #909399;
}

.length-inputs {
    display: flex;
    align-items: center;
    gap: 4px;
}

.no-limit {
    color: #c0c4cc;
}

.template-list {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.template-category {
    background: #fafafa;
    border-radius: 8px;
    padding: 16px;
}

.category-title {
    font-weight: 600;
    margin-bottom: 12px;
    color: #303133;
}

.template-items {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 10px;
}

.template-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 12px;
    background: #fff;
    border: 1px solid #dcdfe6;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
}

.template-item:hover {
    border-color: #409eff;
    color: #409eff;
}

.template-actions {
    display: flex;
    gap: 10px;
}

.test-panel {
    padding: 20px;
    background: #fafafa;
    border-radius: 8px;
}

.test-result {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-top: 20px;
    padding: 12px 16px;
    border-radius: 6px;
    font-weight: 500;
}

.test-result.success {
    background: #f0f9ff;
    color: #67c23a;
}

.test-result.error {
    background: #fef0f0;
    color: #f56c6c;
}

:deep(.el-tabs__content) {
    overflow: visible;
}
</style>
