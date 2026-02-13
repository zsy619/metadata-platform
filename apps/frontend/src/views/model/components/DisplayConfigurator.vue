<template>
    <div class="display-config-panel" :class="{ 'is-collapsed': collapsed }">
        <div class="panel-header" @click="$emit('toggle-collapse')">
            <div class="header-content">
                <el-icon>
                    <Monitor />
                </el-icon>
                <span>显示配置</span>
            </div>
            <el-icon class="collapse-icon">
                <ArrowUp v-if="!collapsed" />
                <ArrowDown v-else />
            </el-icon>
        </div>
        <div v-show="!collapsed" class="panel-body">
            <div class="panel-tip">
                <el-icon><InfoFilled /></el-icon>
                <span>配置字段在前端表单和表格中的显示效果</span>
            </div>
            <el-table :data="displayConfig" border size="small" max-height="450">
                <el-table-column label="字段" min-width="120" show-overflow-tooltip>
                    <template #default="{ row }">
                        <div class="field-cell">
                            <el-icon><Grid /></el-icon>
                            <span>{{ row.name }}</span>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column label="显示标题" min-width="120">
                    <template #default="{ row }">
                        <el-input v-model="row.displayTitle" placeholder="显示标题" size="small" />
                    </template>
                </el-table-column>
                <el-table-column label="组件类型" width="140">
                    <template #default="{ row }">
                        <el-select v-model="row.componentType" placeholder="选择组件" size="small" @change="handleComponentChange(row)">
                            <el-option label="文本输入" value="input" />
                            <el-option label="文本域" value="textarea" />
                            <el-option label="数字输入" value="number" />
                            <el-option label="下拉选择" value="select" />
                            <el-option label="单选" value="radio" />
                            <el-option label="复选框" value="checkbox" />
                            <el-option label="开关" value="switch" />
                            <el-option label="日期选择" value="date" />
                            <el-option label="日期范围" value="daterange" />
                            <el-option label="时间选择" value="time" />
                            <el-option label="日期时间" value="datetime" />
                            <el-option label="文件上传" value="upload" />
                            <el-option label="滑块" value="slider" />
                            <el-option label="评分" value="rate" />
                            <el-option label="颜色选择" value="color" />
                        </el-select>
                    </template>
                </el-table-column>
                <el-table-column label="显示宽度" width="100">
                    <template #default="{ row }">
                        <el-input-number v-model="row.width" :min="50" :max="500" size="small" controls-position="right" />
                    </template>
                </el-table-column>
                <el-table-column label="显示顺序" width="90">
                    <template #default="{ row, $index }">
                        <div class="order-controls">
                            <el-button link type="primary" size="small" :disabled="$index === 0" @click="moveUp($index)">
                                <el-icon><ArrowUp /></el-icon>
                            </el-button>
                            <el-button link type="primary" size="small" :disabled="$index === displayConfig.length - 1" @click="moveDown($index)">
                                <el-icon><ArrowDown /></el-icon>
                            </el-button>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column label="属性" width="180">
                    <template #default="{ row }">
                        <div class="property-tags">
                            <el-tag v-if="row.sortable" size="small" type="info">可排序</el-tag>
                            <el-tag v-if="row.filterable" size="small" type="info">可筛选</el-tag>
                            <el-tag v-if="row.searchable" size="small" type="info">可搜索</el-tag>
                            <el-tag v-if="row.fixed" size="small" type="warning">固定列</el-tag>
                        </div>
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="180">
                    <template #default="{ row, $index }">
                        <div class="action-buttons">
                            <el-button link type="primary" size="small" @click="editFieldConfig(row)">配置</el-button>
                            <el-button link :type="row.sortable ? 'success' : 'info'" size="small" @click="toggleProperty(row, 'sortable')">
                                {{ row.sortable ? '排序' : '排序' }}
                            </el-button>
                            <el-button link :type="row.filterable ? 'success' : 'info'" size="small" @click="toggleProperty(row, 'filterable')">
                                {{ row.filterable ? '筛选' : '筛选' }}
                            </el-button>
                        </div>
                    </template>
                </el-table-column>
            </el-table>
            <div class="panel-actions">
                <el-button size="small" @click="resetToDefault">重置为默认</el-button>
                <el-button type="primary" size="small" @click="applyToAll">应用到所有字段</el-button>
            </div>
        </div>
        <el-dialog v-model="configDialogVisible" title="字段详细配置" width="500px">
            <el-form v-if="currentField" label-width="120px" label-position="right">
                <el-form-item label="占位符">
                    <el-input v-model="currentField.placeholder" placeholder="占位符文本" />
                </el-form-item>
                <el-form-item label="帮助文本">
                    <el-input v-model="currentField.helpText" type="textarea" :rows="2" placeholder="显示在字段下方的帮助信息" />
                </el-form-item>
                <el-form-item label="默认值">
                    <el-input v-model="currentField.defaultValue" placeholder="默认值" />
                </el-form-item>
                <template v-if="['select', 'radio', 'checkbox'].includes(currentField.componentType)">
                    <el-form-item label="选项配置">
                        <div class="options-editor">
                            <div v-for="(opt, idx) in currentField.options" :key="idx" class="option-item">
                                <el-input v-model="opt.label" placeholder="选项标签" style="width: 120px" />
                                <el-input v-model="opt.value" placeholder="选项值" style="width: 100px" />
                                <el-button link type="danger" @click="removeOption(idx)">
                                    <el-icon><Delete /></el-icon>
                                </el-button>
                            </div>
                            <el-button link type="primary" size="small" @click="addOption">
                                <el-icon><Plus /></el-icon> 添加选项
                            </el-button>
                        </div>
                    </el-form-item>
                </template>
                <template v-if="currentField.componentType === 'upload'">
                    <el-form-item label="上传限制">
                        <el-input-number v-model="currentField.uploadLimit" :min="1" :max="10" /> 文件
                    </el-form-item>
                    <el-form-item label="文件类型">
                        <el-select v-model="currentField.acceptTypes" multiple placeholder="选择文件类型">
                            <el-option label="图片" value="image/*" />
                            <el-option label="文档" value=".doc,.docx,.pdf" />
                            <el-option label="Excel" value=".xls,.xlsx" />
                            <el-option label="压缩包" value=".zip,.rar" />
                        </el-select>
                    </el-form-item>
                </template>
                <template v-if="['number', 'slider'].includes(currentField.componentType)">
                    <el-form-item label="最小值">
                        <el-input-number v-model="currentField.min" />
                    </el-form-item>
                    <el-form-item label="最大值">
                        <el-input-number v-model="currentField.max" />
                    </el-form-item>
                    <el-form-item label="步长">
                        <el-input-number v-model="currentField.step" :min="1" />
                    </el-form-item>
                </template>
            </el-form>
            <template #footer>
                <el-button @click="configDialogVisible = false">取消</el-button>
                <el-button type="primary" @click="saveFieldConfig">确定</el-button>
            </template>
        </el-dialog>
    </div>
</template>
<script setup lang="ts">
import { ArrowDown, ArrowUp, Delete, Grid, InfoFilled, Monitor, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { ref } from 'vue'

interface DisplayConfig {
    name: string
    displayTitle: string
    componentType: string
    width: number
    order: number
    sortable: boolean
    filterable: boolean
    searchable: boolean
    fixed: boolean
    placeholder?: string
    helpText?: string
    defaultValue?: string
    options?: { label: string; value: string }[]
    uploadLimit?: number
    acceptTypes?: string[]
    min?: number
    max?: number
    step?: number
}

const props = defineProps<{
    fields: any[]
    collapsed?: boolean
}>()

const emit = defineEmits(['update', 'toggle-collapse'])

const displayConfig = ref<DisplayConfig[]>(
    (props.fields || []).map((f, idx) => ({
        name: f.name || f.fieldName,
        displayTitle: f.displayTitle || f.name || f.fieldName,
        componentType: f.componentType || getDefaultComponentType(f.type || f.fieldType),
        width: f.width || 120,
        order: f.order || idx,
        sortable: f.sortable !== false,
        filterable: f.filterable !== false,
        searchable: f.searchable !== false,
        fixed: f.fixed || false,
        placeholder: f.placeholder || '',
        helpText: f.helpText || '',
        defaultValue: f.defaultValue || '',
        options: f.options || [],
        uploadLimit: f.uploadLimit || 1,
        acceptTypes: f.acceptTypes || [],
        min: f.min,
        max: f.max,
        step: f.step || 1
    }))
)

const configDialogVisible = ref(false)
const currentField = ref<DisplayConfig | null>(null)

function getDefaultComponentType(fieldType: string): string {
    const typeMap: Record<string, string> = {
        int: 'number',
        bigint: 'number',
        float: 'number',
        double: 'number',
        decimal: 'number',
        varchar: 'input',
        text: 'textarea',
        date: 'date',
        datetime: 'datetime',
        time: 'time',
        timestamp: 'datetime',
        bool: 'switch',
        json: 'input'
    }
    return typeMap[fieldType.toLowerCase()] || 'input'
}

const handleComponentChange = (row: DisplayConfig) => {
    if (row.componentType === 'select' && !row.options) {
        row.options = [
            { label: '选项1', value: 'option1' },
            { label: '选项2', value: 'option2' }
        ]
    }
    emit('update', displayConfig.value)
}

const editFieldConfig = (row: DisplayConfig) => {
    currentField.value = { ...row, options: row.options ? [...row.options] : [] }
    configDialogVisible.value = true
}

const saveFieldConfig = () => {
    if (!currentField.value) return
    const idx = displayConfig.value.findIndex(f => f.name === currentField.value?.name)
    if (idx !== -1) {
        displayConfig.value[idx] = { ...currentField.value }
    }
    configDialogVisible.value = false
    emit('update', displayConfig.value)
    ElMessage.success('配置已保存')
}

const addOption = () => {
    if (currentField.value && currentField.value.options) {
        currentField.value.options.push({ label: '', value: '' })
    }
}

const removeOption = (idx: number) => {
    if (currentField.value && currentField.value.options) {
        currentField.value.options.splice(idx, 1)
    }
}

const moveUp = (index: number) => {
    if (index > 0) {
        const temp = displayConfig.value[index]
        displayConfig.value[index] = displayConfig.value[index - 1]
        displayConfig.value[index - 1] = temp
        updateOrder()
    }
}

const moveDown = (index: number) => {
    if (index < displayConfig.value.length - 1) {
        const temp = displayConfig.value[index]
        displayConfig.value[index] = displayConfig.value[index + 1]
        displayConfig.value[index + 1] = temp
        updateOrder()
    }
}

const updateOrder = () => {
    displayConfig.value.forEach((item, idx) => {
        item.order = idx
    })
    emit('update', displayConfig.value)
}

const toggleProperty = (row: DisplayConfig, prop: keyof DisplayConfig) => {
    ;(row as any)[prop] = !(row as any)[prop]
    emit('update', displayConfig.value)
}

const resetToDefault = () => {
    displayConfig.value = (props.fields || []).map((f, idx) => ({
        name: f.name || f.fieldName,
        displayTitle: f.name || f.fieldName,
        componentType: getDefaultComponentType(f.type || f.fieldType),
        width: 120,
        order: idx,
        sortable: true,
        filterable: true,
        searchable: true,
        fixed: false
    }))
    emit('update', displayConfig.value)
    ElMessage.success('已重置为默认配置')
}

const applyToAll = () => {
    const defaultComponent = 'input'
    displayConfig.value.forEach(item => {
        item.componentType = defaultComponent
    })
    emit('update', displayConfig.value)
    ElMessage.success('已应用到所有字段')
}

defineExpose({
    displayConfig,
    getData: () => displayConfig.value
})
</script>
<style scoped>
.display-config-panel {
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

.order-controls {
    display: flex;
    gap: 4px;
}

.property-tags {
    display: flex;
    gap: 4px;
    flex-wrap: wrap;
}

.action-buttons {
    display: flex;
    gap: 4px;
}

.panel-actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 16px;
}

.options-editor {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.option-item {
    display: flex;
    gap: 8px;
    align-items: center;
}
</style>
