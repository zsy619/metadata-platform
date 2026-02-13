<template>
    <div class="dynamic-form">
        <el-form ref="formRef" :model="formData" :rules="formRules" :label-width="labelWidth" :label-position="labelPosition" :size="size" :disabled="disabled" :show-message="showMessage">
            <el-row :gutter="gutter">
                <el-col v-for="field in visibleFields" :key="field.prop" :span="field.span || 24" :xs="field.xs || 24" :sm="field.sm || 24" :md="field.md || 12" :lg="field.lg || 8">
                    <el-form-item :label="field.label" :prop="field.prop" :required="field.required" :error="field.error">
                        <template v-if="field.type === 'input' || field.type === 'text'">
                            <el-input v-model="formData[field.prop]" :type="field.inputType || 'text'" :placeholder="field.placeholder" :disabled="field.disabled || disabled" :clearable="field.clearable !== false" :readonly="field.readonly" :maxlength="field.maxlength" :minlength="field.minlength" :show-word-limit="field.showWordLimit" :prefix-icon="field.prefixIcon" :suffix-icon="field.suffixIcon" @change="handleChange(field)">
                                <template #prefix v-if="field.prefix">{{ field.prefix }}</template>
                                <template #suffix v-if="field.suffix">{{ field.suffix }}</template>
                            </el-input>
                        </template>
                        <template v-else-if="field.type === 'textarea'">
                            <el-input v-model="formData[field.prop]" type="textarea" :placeholder="field.placeholder" :disabled="field.disabled || disabled" :clearable="field.clearable !== false" :readonly="field.readonly" :rows="field.rows || 3" :maxlength="field.maxlength" :show-word-limit="field.showWordLimit" @change="handleChange(field)" />
                        </template>
                        <template v-else-if="field.type === 'number'">
                            <el-input-number v-model="formData[field.prop]" :placeholder="field.placeholder" :disabled="field.disabled || disabled" :min="field.min" :max="field.max" :step="field.step" :precision="field.precision" :controls="field.controls !== false" :controls-position="field.controlsPosition || 'right'" @change="handleChange(field)" />
                        </template>
                        <template v-else-if="field.type === 'select'">
                            <el-select v-model="formData[field.prop]" :placeholder="field.placeholder" :disabled="field.disabled || disabled" :clearable="field.clearable !== false" :multiple="field.multiple" :collapse-tags="field.collapseTags" :filterable="field.filterable" :remote="field.remote" :remote-method="field.remoteMethod" :loading="field.loading" @change="handleChange(field)">
                                <el-option v-for="opt in field.options" :key="opt.value" :label="opt.label" :value="opt.value" :disabled="opt.disabled" />
                            </el-select>
                        </template>
                        <template v-else-if="field.type === 'radio'">
                            <el-radio-group v-model="formData[field.prop]" :disabled="field.disabled || disabled" @change="handleChange(field)">
                                <el-radio v-for="opt in field.options" :key="opt.value" :label="opt.value" :disabled="opt.disabled">{{ opt.label }}</el-radio>
                            </el-radio-group>
                        </template>
                        <template v-else-if="field.type === 'checkbox'">
                            <el-checkbox-group v-model="formData[field.prop]" :disabled="field.disabled || disabled" @change="handleChange(field)">
                                <el-checkbox v-for="opt in field.options" :key="opt.value" :label="opt.value" :disabled="opt.disabled">{{ opt.label }}</el-checkbox>
                            </el-checkbox-group>
                        </template>
                        <template v-else-if="field.type === 'switch'">
                            <el-switch v-model="formData[field.prop]" :disabled="field.disabled || disabled" :active-text="field.activeText" :inactive-text="field.inactiveText" :active-value="field.activeValue" :inactive-value="field.inactiveValue" @change="handleChange(field)" />
                        </template>
                        <template v-else-if="field.type === 'date'">
                            <el-date-picker v-model="formData[field.prop]" type="date" :placeholder="field.placeholder" :disabled="field.disabled || disabled" :clearable="field.clearable !== false" :format="field.format" :value-format="field.valueFormat" :picker-options="field.pickerOptions" @change="handleChange(field)" />
                        </template>
                        <template v-else-if="field.type === 'daterange'">
                            <el-date-picker v-model="formData[field.prop]" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" :disabled="field.disabled || disabled" :clearable="field.clearable !== false" :format="field.format" :value-format="field.valueFormat" @change="handleChange(field)" />
                        </template>
                        <template v-else-if="field.type === 'datetime'">
                            <el-date-picker v-model="formData[field.prop]" type="datetime" :placeholder="field.placeholder" :disabled="field.disabled || disabled" :clearable="field.clearable !== false" :format="field.format" :value-format="field.valueFormat" :picker-options="field.pickerOptions" @change="handleChange(field)" />
                        </template>
                        <template v-else-if="field.type === 'datetimerange'">
                            <el-date-picker v-model="formData[field.prop]" type="datetimerange" range-separator="至" start-placeholder="开始时间" end-placeholder="结束时间" :disabled="field.disabled || disabled" :clearable="field.clearable !== false" :format="field.format" :value-format="field.valueFormat" @change="handleChange(field)" />
                        </template>
                        <template v-else-if="field.type === 'time'">
                            <el-time-picker v-model="formData[field.prop]" :placeholder="field.placeholder" :disabled="field.disabled || disabled" :clearable="field.clearable !== false" :format="field.format" :value-format="field.valueFormat" @change="handleChange(field)" />
                        </template>
                        <template v-else-if="field.type === 'timeRange'">
                            <el-time-picker v-model="formData[field.prop]" is-range range-separator="至" start-placeholder="开始时间" end-placeholder="结束时间" :disabled="field.disabled || disabled" :clearable="field.clearable !== false" :format="field.format" :value-format="field.valueFormat" @change="handleChange(field)" />
                        </template>
                        <template v-else-if="field.type === 'cascader'">
                            <el-cascader v-model="formData[field.prop]" :options="field.options" :placeholder="field.placeholder" :disabled="field.disabled || disabled" :clearable="field.clearable !== false" :filterable="field.filterable" :props="field.cascaderProps" @change="handleChange(field)" />
                        </template>
                        <template v-else-if="field.type === 'slider'">
                            <el-slider v-model="formData[field.prop]" :disabled="field.disabled || disabled" :min="field.min" :max="field.max" :step="field.step" :show-stops="field.showStops" :range="field.range" @change="handleChange(field)" />
                        </template>
                        <template v-else-if="field.type === 'rate'">
                            <el-rate v-model="formData[field.prop]" :disabled="field.disabled || disabled" :max="field.max" :allow-half="field.allowHalf" :show-text="field.showText" :texts="field.texts" @change="handleChange(field)" />
                        </template>
                        <template v-else-if="field.type === 'color'">
                            <el-color-picker v-model="formData[field.prop]" :disabled="field.disabled || disabled" :show-alpha="field.showAlpha" @change="handleChange(field)" />
                        </template>
                        <template v-else-if="field.type === 'upload'">
                            <el-upload v-model:file-list="formData[field.prop]" :action="field.action" :headers="field.headers" :data="field.data" :accept="field.accept" :limit="field.limit" :multiple="field.multiple" :drag="field.drag" :disabled="field.disabled || disabled" :list-type="field.listType" :auto-upload="field.autoUpload !== false" :before-upload="field.beforeUpload" :before-remove="field.beforeRemove" :on-success="(res, file, fileList) => handleUploadSuccess(res, file, fileList, field)" :on-error="(err, file, fileList) => handleUploadError(err, file, fileList, field)" :on-progress="(event, file, fileList) => handleUploadProgress(event, file, fileList, field)" :on-preview="field.onPreview" :on-remove="(file, fileList) => handleUploadRemove(file, fileList, field)">
                                <el-button v-if="field.listType !== 'picture-card'" type="primary" :disabled="field.disabled">
                                    <el-icon v-if="field.uploadIcon"><Upload /></el-icon>{{ field.uploadText || '上传' }}
                                </el-button>
                                <el-icon v-else class="el-icon--upload"><Upload /></el-icon>
                            </el-upload>
                        </template>
                        <template v-else-if="field.type === 'slot'">
                            <slot :name="field.slotName" :field="field" :model="formData"></slot>
                        </template>
                        <template v-if="field.suffixText">
                            <span class="suffix-text">{{ field.suffixText }}</span>
                        </template>
                    </el-form-item>
                </el-col>
            </el-row>
            <el-form-item v-if="showButtons">
                <el-button type="primary" :loading="loading" @click="handleSubmit">{{ submitText }}</el-button>
                <el-button @click="handleReset">{{ resetText }}</el-button>
                <slot name="buttons"></slot>
            </el-form-item>
        </el-form>
    </div>
</template>
<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'

export interface FormField {
    prop: string
    label: string
    type: 'input' | 'textarea' | 'number' | 'select' | 'radio' | 'checkbox' | 'switch' | 'date' | 'daterange' | 'datetime' | 'datetimerange' | 'time' | 'timeRange' | 'cascader' | 'slider' | 'rate' | 'color' | 'upload' | 'slot'
    value?: any
    defaultValue?: any
    placeholder?: string
    disabled?: boolean
    readonly?: boolean
    required?: boolean
    rules?: any[]
    span?: number
    xs?: number
    sm?: number
    md?: number
    lg?: number
    inputType?: string
    maxlength?: number
    minlength?: number
    showWordLimit?: boolean
    prefix?: string
    prefixIcon?: string
    suffix?: string
    suffixIcon?: string
    suffixText?: string
    rows?: number
    min?: number
    max?: number
    step?: number
    precision?: number
    controls?: boolean
    controlsPosition?: 'right' | 'default'
    options?: { label: string; value: any; disabled?: boolean }[]
    multiple?: boolean
    collapseTags?: boolean
    filterable?: boolean
    remote?: boolean
    remoteMethod?: (query: string) => void
    loading?: boolean
    activeText?: string
    inactiveText?: string
    activeValue?: any
    inactiveValue?: any
    format?: string
    valueFormat?: string
    pickerOptions?: any
    cascaderProps?: any
    showStops?: boolean
    range?: boolean
    allowHalf?: boolean
    showText?: boolean
    texts?: string[]
    showAlpha?: boolean
    action?: string
    headers?: Record<string, string>
    data?: Record<string, any>
    accept?: string
    limit?: number
    drag?: boolean
    listType?: 'text' | 'picture' | 'picture-card'
    autoUpload?: boolean
    beforeUpload?: (file: File) => boolean | Promise<File>
    beforeRemove?: (file: File, fileList: File[]) => boolean | Promise<boolean>
    onSuccess?: (response: any, file: File, fileList: File[]) => void
    onError?: (error: Error, file: File, fileList: File[]) => void
    onProgress?: (event: any, file: File, fileList: File[]) => void
    onPreview?: (file: File) => void
    onRemove?: (file: File, fileList: File[]) => void
    slotName?: string
    uploadIcon?: boolean
    uploadText?: string
    clearable?: boolean
    error?: string
    vif?: boolean
}

interface Props {
    modelValue?: Record<string, any>
    fields?: FormField[]
    rules?: FormRules
    labelWidth?: string | number
    labelPosition?: 'left' | 'right' | 'top'
    size?: 'large' | 'default' | 'small'
    disabled?: boolean
    showMessage?: boolean
    gutter?: number
    showButtons?: boolean
    submitText?: string
    resetText?: string
    loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    modelValue: () => ({}),
    fields: () => [],
    rules: () => ({}),
    labelWidth: '120px',
    labelPosition: 'right',
    size: 'default',
    disabled: false,
    showMessage: true,
    gutter: 20,
    showButtons: true,
    submitText: '提交',
    resetText: '重置',
    loading: false
})

const emit = defineEmits(['update:modelValue', 'submit', 'reset', 'change'])

const formRef = ref<FormInstance>()
const formData = ref<Record<string, any>>({})

const visibleFields = computed(() => {
    return props.fields.filter(field => field.vif === undefined || field.vif)
})

watch(() => props.modelValue, (val) => {
    formData.value = { ...val }
}, { immediate: true, deep: true })

watch(formData, (val) => {
    emit('update:modelValue', val)
}, { deep: true })

const formRules = computed(() => {
    const rules: FormRules = { ...props.rules }
    props.fields.forEach(field => {
        if (field.required && !rules[field.prop]) {
            rules[field.prop] = [{ required: true, message: `${field.label}不能为空`, trigger: ['blur', 'change'] }]
        }
    })
    return rules
})

const handleChange = (field: FormField) => {
    emit('change', field, formData.value[field.prop])
}

const handleSubmit = async () => {
    if (!formRef.value) return
    await formRef.value.validate((valid) => {
        if (valid) {
            emit('submit', formData.value)
        }
    })
}

const handleReset = () => {
    formRef.value?.resetFields()
    emit('reset')
}

const validate = async () => {
    return await formRef.value?.validate()
}

const validateField = async (prop: string) => {
    return await formRef.value?.validateField(prop)
}

const clearValidate = (props?: string | string[]) => {
    formRef.value?.clearValidate(props)
}

const resetFields = () => {
    formRef.value?.resetFields()
}

const setFieldValue = (prop: string, value: any) => {
    formData.value[prop] = value
}

const getFieldValue = (prop: string) => {
    return formData.value[prop]
}

const getValues = () => {
    return { ...formData.value }
}

const setValues = (values: Record<string, any>) => {
    formData.value = { ...formData.value, ...values }
}

const handleUploadSuccess = (response: any, file: File, fileList: File[], field: FormField) => {
    if (field.onSuccess) {
        field.onSuccess(response, file, fileList)
    }
}

const handleUploadError = (error: Error, file: File, fileList: File[], field: FormField) => {
    if (field.onError) {
        field.onError(error, file, fileList)
    }
}

const handleUploadProgress = (event: any, file: File, fileList: File[], field: FormField) => {
    if (field.onProgress) {
        field.onProgress(event, file, fileList)
    }
}

const handleUploadRemove = (file: File, fileList: File[], field: FormField) => {
    if (field.onRemove) {
        field.onRemove(file, fileList)
    }
}

defineExpose({
    formRef,
    validate,
    validateField,
    clearValidate,
    resetFields,
    setFieldValue,
    getFieldValue,
    getValues,
    setValues
})
</script>
<style scoped>
.dynamic-form {
    width: 100%;
}

.suffix-text {
    margin-left: 8px;
    color: #909399;
    font-size: 14px;
}

:deep(.el-upload-list) {
    margin-top: 8px;
}

:deep(.el-upload-dragger) {
    padding: 20px;
}
</style>
