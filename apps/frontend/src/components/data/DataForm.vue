<template>
    <el-dialog :model-value="visible" title="数据表单" width="600px" :close-on-click-modal="false" @update:model-value="$emit('update:visible', $event)">
        <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
            <el-form-item v-for="field in fields" :key="field.prop" :label="field.label" :prop="field.prop">
                <template v-if="field.type === 'input' || !field.type">
                    <el-input v-model="formData[field.prop]" :placeholder="field.placeholder" :disabled="field.disabled" />
                </template>
                <template v-else-if="field.type === 'textarea'">
                    <el-input v-model="formData[field.prop]" type="textarea" :rows="3" :placeholder="field.placeholder" :disabled="field.disabled" />
                </template>
                <template v-else-if="field.type === 'number'">
                    <el-input-number v-model="formData[field.prop]" :min="field.min" :max="field.max" :disabled="field.disabled" controls-position="right" />
                </template>
                <template v-else-if="field.type === 'switch'">
                    <el-switch v-model="formData[field.prop]" :disabled="field.disabled" />
                </template>
                <template v-else-if="field.type === 'select'">
                    <el-select v-model="formData[field.prop]" :placeholder="field.placeholder" :disabled="field.disabled" style="width: 100%">
                        <el-option v-for="opt in field.options" :key="opt.value" :label="opt.label" :value="opt.value" />
                    </el-select>
                </template>
                <template v-else-if="field.type === 'date'">
                    <el-date-picker v-model="formData[field.prop]" type="date" :placeholder="field.placeholder" :disabled="field.disabled" value-format="YYYY-MM-DD" style="width: 100%" />
                </template>
                <template v-else-if="field.type === 'datetime'">
                    <el-date-picker v-model="formData[field.prop]" type="datetime" :placeholder="field.placeholder" :disabled="field.disabled" value-format="YYYY-MM-DD HH:mm:ss" style="width: 100%" />
                </template>
                <template v-else-if="field.type === 'radio'">
                    <el-radio-group v-model="formData[field.prop]" :disabled="field.disabled">
                        <el-radio v-for="opt in field.options" :key="opt.value" :label="opt.value">{{ opt.label }}</el-radio>
                    </el-radio-group>
                </template>
                <template v-else-if="field.type === 'checkbox'">
                    <el-checkbox-group v-model="formData[field.prop]" :disabled="field.disabled">
                        <el-checkbox v-for="opt in field.options" :key="opt.value" :label="opt.value">{{ opt.label }}</el-checkbox>
                    </el-checkbox-group>
                </template>
                <template v-if="field.helpText">
                    <div class="help-text">{{ field.helpText }}</div>
                </template>
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="$emit('update:visible', false)">取消</el-button>
            <el-button type="primary" :loading="submitting" @click="handleSubmit">{{ submitting ? '提交中...' : '提交' }}</el-button>
        </template>
    </el-dialog>
</template>
<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { computed, nextTick, ref, watch } from 'vue'

interface Field {
    prop: string
    label: string
    type?: string
    placeholder?: string
    disabled?: boolean
    required?: boolean
    options?: { label: string; value: any }[]
    helpText?: string
    min?: number
    max?: number
    defaultValue?: any
    validator?: (rule: any, value: any, callback: any) => void
}

const props = defineProps<{
    visible: boolean
    modelId?: number
    fields: Field[]
    data?: Record<string, any>
}>()

const emit = defineEmits(['update:visible', 'success'])

const formRef = ref<FormInstance>()
const submitting = ref(false)
const isEdit = computed(() => !!props.data && Object.keys(props.data).length > 0)

const formData = ref<Record<string, any>>({})

watch(() => props.visible, (val) => {
    if (val) {
        initFormData()
    }
})

watch(() => props.data, () => {
    initFormData()
}, { immediate: true })

const initFormData = () => {
    nextTick(() => {
        formRef.value?.resetFields()
        const data: Record<string, any> = {}
        props.fields.forEach(field => {
            if (props.data && props.data[field.prop] !== undefined) {
                data[field.prop] = props.data[field.prop]
            } else if (field.defaultValue !== undefined) {
                data[field.prop] = field.defaultValue
            } else {
                data[field.prop] = ''
            }
        })
        formData.value = data
    })
}

const formRules = computed<FormRules>(() => {
    const rules: FormRules = {}
    props.fields.forEach(field => {
        if (field.required) {
            rules[field.prop] = [
                { required: true, message: `${field.label}不能为空`, trigger: ['blur', 'change'] }
            ]
        }
        if (field.validator) {
            rules[field.prop] = rules[field.prop] || []
            rules[field.prop].push({ validator: field.validator, trigger: 'blur' })
        }
    })
    return rules
})

const handleSubmit = async () => {
    if (!formRef.value) return
    await formRef.value.validate(async (valid) => {
        if (valid) {
            submitting.value = true
            try {
                await new Promise(resolve => setTimeout(resolve, 500))
                ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
                emit('success', formData.value)
            } catch (error) {
                ElMessage.error('操作失败')
            } finally {
                submitting.value = false
            }
        }
    })
}

const setData = (data: Record<string, any>) => {
    formData.value = { ...data }
}

const resetForm = () => {
    formRef.value?.resetFields()
}

defineExpose({
    setData,
    resetForm
})
</script>
<style scoped>
.help-text {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
}
</style>
