<template>
    <div class="query-builder">
        <div class="builder-header">
            <div class="header-left">
                <el-icon><Filter /></el-icon>
                <span>查询条件</span>
            </div>
            <div class="header-actions">
                <el-button size="small" type="primary" :icon="Plus" @click="addCondition">添加条件</el-button>
                <el-button size="small" :icon="RefreshLeft" @click="resetConditions">重置</el-button>
            </div>
        </div>
        <div class="builder-content">
            <div v-if="conditions.length === 0" class="empty-conditions">
                <el-empty description="暂无查询条件，点击上方按钮添加" :image-size="80" />
            </div>
            <div v-else class="conditions-list">
                <div v-for="(group, groupIndex) in conditionGroups" :key="groupIndex" class="condition-group">
                    <div v-if="groupIndex > 0" class="group-operator">
                        <el-select v-model="group.logicOperator" size="small" style="width: 80px">
                            <el-option label="并且" value="AND" />
                            <el-option label="或者" value="OR" />
                        </el-select>
                    </div>
                    <div class="group-conditions">
                        <div v-for="(condition, condIndex) in group.conditions" :key="condIndex" class="condition-row">
                            <div v-if="condIndex > 0" class="condition-operator">
                                <el-select v-model="condition.logicOperator" size="small" style="width: 70px">
                                    <el-option label="并且" value="AND" />
                                    <el-option label="或者" value="OR" />
                                </el-select>
                            </div>
                            <div v-if="condIndex === 0 && groupIndex > 0" class="condition-brackets">
                                <el-checkbox v-model="condition.leftBracket" size="small">(</el-checkbox>
                            </div>
                            <div class="condition-fields">
                                <el-select v-model="condition.field" placeholder="选择字段" size="small" style="width: 150px" filterable @change="handleFieldChange(condition)">
                                    <el-option v-for="field in fields" :key="field.name" :label="field.label" :value="field.name">
                                        <div class="field-option">
                                            <span>{{ field.label }}</span>
                                            <el-tag size="small" type="info">{{ field.type }}</el-tag>
                                        </div>
                                    </el-option>
                                </el-select>
                                <el-select v-model="condition.operator" placeholder="操作符" size="small" style="width: 120px" @change="handleOperatorChange(condition)">
                                    <el-option v-for="op in getOperators(condition.field)" :key="op.value" :label="op.label" :value="op.value" />
                                </el-select>
                                <div class="condition-values">
                                    <template v-if="condition.operator === 'between'">
                                        <el-input v-model="condition.value1" placeholder="最小值" size="small" style="width: 100px" />
                                        <span class="between-separator">-</span>
                                        <el-input v-model="condition.value2" placeholder="最大值" size="small" style="width: 100px" />
                                    </template>
                                    <template v-else-if="condition.operator === 'in' || condition.operator === 'not_in'">
                                        <el-select v-model="condition.value1" placeholder="选择值" size="small" style="width: 200px" multiple filterable allow-create>
                                            <el-option v-for="opt in getFieldOptions(condition.field)" :key="opt.value" :label="opt.label" :value="opt.value" />
                                        </el-select>
                                    </template>
                                    <template v-else-if="condition.operator === 'is_null' || condition.operator === 'is_not_null'">
                                        <span class="no-value">无需填写</span>
                                    </template>
                                    <template v-else-if="isDateField(condition.field)">
                                        <el-date-picker v-model="condition.value1" type="date" placeholder="选择日期" size="small" style="width: 150px" value-format="YYYY-MM-DD" />
                                    </template>
                                    <template v-else-if="isNumberField(condition.field)">
                                        <el-input-number v-model="condition.value1" placeholder="值" size="small" controls-position="right" style="width: 150px" />
                                    </template>
                                    <template v-else>
                                        <el-input v-model="condition.value1" placeholder="请输入值" size="small" style="width: 150px" />
                                    </template>
                                </div>
                            </div>
                            <div v-if="condIndex === 0 && groupIndex > 0" class="condition-brackets">
                                <el-checkbox v-model="condition.rightBracket" size="small">)</el-checkbox>
                            </div>
                            <div class="condition-actions">
                                <el-button link type="primary" size="small" :icon="Plus" @click="addConditionInGroup(group)" />
                                <el-button link type="danger" size="small" :icon="Delete" :disabled="group.conditions.length <= 1" @click="removeCondition(groupIndex, condIndex)" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div v-if="showSqlPreview" class="builder-footer">
            <div class="sql-preview">
                <div class="sql-label">
                    <el-icon><Document /></el-icon>
                    <span>SQL预览</span>
                </div>
                <div class="sql-content">
                    <code>{{ generatedSql }}</code>
                </div>
            </div>
        </div>
    </div>
</template>
<script setup lang="ts">
import { Delete, Document, Filter, Plus, RefreshLeft } from '@element-plus/icons-vue'
import { computed, ref, watch } from 'vue'

interface Field {
    name: string
    label: string
    type: string
    options?: { label: string; value: any }[]
}

interface Condition {
    field: string
    operator: string
    value1: any
    value2: any
    logicOperator: string
    leftBracket: boolean
    rightBracket: boolean
}

interface ConditionGroup {
    logicOperator: string
    conditions: Condition[]
}

const props = defineProps<{
    fields: Field[]
    modelValue?: Condition[]
    showSqlPreview?: boolean
}>()

const emit = defineEmits(['update:modelValue', 'change'])

const conditions = ref<Condition[]>(props.modelValue || [])

const conditionGroups = computed<ConditionGroup[]>(() => {
    if (conditions.value.length === 0) return []
    const groups: ConditionGroup[] = []
    let currentGroup: ConditionGroup = { logicOperator: 'AND', conditions: [] }
    conditions.value.forEach((cond, idx) => {
        if (idx === 0 || (cond.logicOperator && cond.logicOperator !== 'AND')) {
            if (currentGroup.conditions.length > 0) {
                groups.push(currentGroup)
            }
            currentGroup = { logicOperator: cond.logicOperator || 'AND', conditions: [cond] }
        } else {
            currentGroup.conditions.push(cond)
        }
    })
    if (currentGroup.conditions.length > 0) {
        groups.push(currentGroup)
    }
    return groups
})

const generatedSql = computed(() => {
    if (conditions.value.length === 0) return ''
    const sqlParts: string[] = []
    conditionGroups.value.forEach((group, groupIdx) => {
        if (groupIdx > 0) {
            sqlParts.push(group.logicOperator)
        }
        const condParts: string[] = []
        group.conditions.forEach((cond, condIdx) => {
            if (condIdx > 0) {
                condParts.push(cond.logicOperator || 'AND')
            }
            const field = props.fields.find(f => f.name === cond.field)
            const fieldName = field ? field.label : cond.field
            const operatorSql = getOperatorSql(cond.operator, fieldName, cond.value1, cond.value2)
            let part = operatorSql
            if (cond.leftBracket && groupIdx > 0) part = '(' + part
            if (cond.rightBracket && groupIdx > 0) part = part + ')'
            condParts.push(part)
        })
        sqlParts.push(condParts.join(' '))
    })
    return 'WHERE ' + sqlParts.join(' ')
})

const operators = [
    { value: '=', label: '等于', types: ['all'] },
    { value: '!=', label: '不等于', types: ['all'] },
    { value: '>', label: '大于', types: ['number', 'date', 'datetime'] },
    { value: '>=', label: '大于等于', types: ['number', 'date', 'datetime'] },
    { value: '<', label: '小于', types: ['number', 'date', 'datetime'] },
    { value: '<=', label: '小于等于', types: ['number', 'date', 'datetime'] },
    { value: 'like', label: '包含', types: ['string'] },
    { value: 'not_like', label: '不包含', types: ['string'] },
    { value: 'starts_with', label: '开头是', types: ['string'] },
    { value: 'ends_with', label: '结尾是', types: ['string'] },
    { value: 'between', label: '介于', types: ['number', 'date', 'datetime'] },
    { value: 'in', label: '在范围内', types: ['all'] },
    { value: 'not_in', label: '不在范围内', types: ['all'] },
    { value: 'is_null', label: '为空', types: ['all'] },
    { value: 'is_not_null', label: '不为空', types: ['all'] }
]

const getFieldType = (fieldName: string) => {
    const field = props.fields.find(f => f.name === fieldName)
    return field?.type || 'string'
}

const isDateField = (fieldName: string) => {
    const type = getFieldType(fieldName)
    return ['date', 'datetime', 'time'].includes(type)
}

const isNumberField = (fieldName: string) => {
    const type = getFieldType(fieldName)
    return ['int', 'bigint', 'float', 'double', 'decimal', 'number'].includes(type)
}

const getOperators = (fieldName: string) => {
    const type = getFieldType(fieldName)
    return operators.filter(op => op.types.includes('all') || op.types.includes(type))
}

const getFieldOptions = (fieldName: string) => {
    const field = props.fields.find(f => f.name === fieldName)
    return field?.options || []
}

const getOperatorSql = (operator: string, fieldName: string, value1: any, value2: any) => {
    const field = props.fields.find(f => f.label === fieldName || f.name === fieldName)
    const safeField = field ? `\`${field.name}\`` : `\`${fieldName}\``
    switch (operator) {
        case '=': return `${safeField} = '${value1}'`
        case '!=': return `${safeField} != '${value1}'`
        case '>': return `${safeField} > '${value1}'`
        case '>=': return `${safeField} >= '${value1}'`
        case '<': return `${safeField} < '${value1}'`
        case '<=': return `${safeField} <= '${value1}'`
        case 'like': return `${safeField} LIKE '%${value1}%'`
        case 'not_like': return `${safeField} NOT LIKE '%${value1}%'`
        case 'starts_with': return `${safeField} LIKE '${value1}%'`
        case 'ends_with': return `${safeField} LIKE '%${value1}'`
        case 'between': return `${safeField} BETWEEN '${value1}' AND '${value2}'`
        case 'in': return `${safeField} IN (${Array.isArray(value1) ? value1.map(v => `'${v}'`).join(', ') : value1})`
        case 'not_in': return `${safeField} NOT IN (${Array.isArray(value1) ? value1.map(v => `'${v}'`).join(', ') : value1})`
        case 'is_null': return `${safeField} IS NULL`
        case 'is_not_null': return `${safeField} IS NOT NULL`
        default: return `${safeField} = '${value1}'`
    }
}

const addCondition = () => {
    if (conditions.value.length > 0 && conditions.value[conditions.value.length - 1].logicOperator === undefined) {
        conditions.value[conditions.value.length - 1].logicOperator = 'AND'
    }
    conditions.value.push({
        field: props.fields[0]?.name || '',
        operator: '=',
        value1: '',
        value2: '',
        logicOperator: 'AND',
        leftBracket: false,
        rightBracket: false
    })
    emitChange()
}

const addConditionInGroup = (group: ConditionGroup) => {
    group.conditions.push({
        field: props.fields[0]?.name || '',
        operator: '=',
        value1: '',
        value2: '',
        logicOperator: 'AND',
        leftBracket: false,
        rightBracket: false
    })
    syncConditions()
    emitChange()
}

const removeCondition = (groupIndex: number, condIndex: number) => {
    const group = conditionGroups.value[groupIndex]
    if (group.conditions.length > 1) {
        group.conditions.splice(condIndex, 1)
    } else if (conditionGroups.value.length > 1) {
        conditionGroups.value.splice(groupIndex, 1)
    } else {
        conditions.value = []
    }
    syncConditions()
    emitChange()
}

const resetConditions = () => {
    conditions.value = []
    emitChange()
}

const handleFieldChange = (condition: Condition) => {
    condition.operator = '='
    condition.value1 = ''
    condition.value2 = ''
    emitChange()
}

const handleOperatorChange = (condition: Condition) => {
    emitChange()
}

const syncConditions = () => {
    conditions.value = []
    conditionGroups.value.forEach((group, groupIdx) => {
        group.conditions.forEach((cond, condIdx) => {
            if (groupIdx > 0 && condIdx === 0) {
                cond.logicOperator = group.logicOperator
            }
            conditions.value.push({ ...cond })
        })
    })
}

const emitChange = () => {
    emit('update:modelValue', conditions.value)
    emit('change', conditions.value, generatedSql.value)
}

watch(() => props.modelValue, (val) => {
    conditions.value = val || []
}, { deep: true })

defineExpose({
    getConditions: () => conditions.value,
    getSql: () => generatedSql.value
})
</script>
<style scoped>
.query-builder {
    border: 1px solid #ebeef5;
    border-radius: 8px;
    background: #fff;
}

.builder-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    border-bottom: 1px solid #ebeef5;
    background: #fafafa;
    border-radius: 8px 8px 0 0;
}

.header-left {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
    color: #303133;
}

.header-actions {
    display: flex;
    gap: 8px;
}

.builder-content {
    padding: 16px;
    min-height: 150px;
    max-height: 400px;
    overflow-y: auto;
}

.empty-conditions {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 120px;
}

.conditions-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.condition-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.group-operator {
    padding-left: 20px;
}

.group-conditions {
    display: flex;
    flex-direction: column;
    gap: 8px;
    padding-left: 20px;
}

.condition-row {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px;
    background: #f5f7fa;
    border-radius: 6px;
}

.condition-operator {
    min-width: 70px;
}

.condition-brackets {
    display: flex;
    align-items: center;
}

.condition-fields {
    display: flex;
    align-items: center;
    gap: 8px;
}

.condition-values {
    display: flex;
    align-items: center;
    gap: 4px;
}

.between-separator {
    color: #909399;
    padding: 0 4px;
}

.no-value {
    color: #909399;
    font-size: 12px;
}

.condition-actions {
    display: flex;
    gap: 4px;
    margin-left: auto;
}

.field-option {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.builder-footer {
    border-top: 1px solid #ebeef5;
    padding: 12px 16px;
    background: #fafafa;
    border-radius: 0 0 8px 8px;
}

.sql-preview {
    display: flex;
    align-items: flex-start;
    gap: 8px;
}

.sql-label {
    display: flex;
    align-items: center;
    gap: 4px;
    color: #606266;
    font-size: 13px;
    font-weight: 500;
}

.sql-content {
    flex: 1;
    background: #2d2d2d;
    border-radius: 6px;
    padding: 10px 12px;
    overflow-x: auto;
}

.sql-content code {
    color: #50fa7b;
    font-family: 'Monaco', 'Menlo', monospace;
    font-size: 12px;
    line-height: 1.5;
}
</style>
