<template>
    <el-dialog v-model="dialogVisible" title="生成新密钥" width="500px" destroy-on-close>
        <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
            <el-form-item label="密钥类型" prop="key_type">
                <el-select v-model="formData.key_type" placeholder="请选择密钥类型" style="width: 100%" @change="handleKeyTypeChange">
                    <el-option label="RSA" value="rsa" />
                    <el-option label="Octet (对称密钥)" value="octet" />
                </el-select>
            </el-form-item>
            <el-form-item label="算法" prop="algorithm">
                <el-select v-model="formData.algorithm" placeholder="请选择算法" style="width: 100%">
                    <el-option v-for="alg in availableAlgorithms" :key="alg.value" :label="alg.label" :value="alg.value" />
                </el-select>
            </el-form-item>
            <el-form-item label="备注" prop="remark">
                <el-input v-model="formData.remark" type="textarea" :rows="2" placeholder="请输入备注" />
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="handleCancel">取消</el-button>
            <el-button type="primary" @click="handleSubmit" :loading="submitLoading">生成</el-button>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus'
import { computed, reactive, ref, watch } from 'vue'

// Props
interface Props {
    modelValue: boolean
}

const props = withDefaults(defineProps<Props>(), {
    modelValue: false
})

// Emits
const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: { key_type: 'rsa' | 'octet'; algorithm: string; remark: string }]
}>()

// 内部状态
const dialogVisible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
})

const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive({ key_type: 'rsa' as 'rsa' | 'octet', algorithm: 'RS256', remark: '' })

const formRules: FormRules = {
    key_type: [{ required: true, message: '请选择密钥类型', trigger: 'change' }],
    algorithm: [{ required: true, message: '请选择算法', trigger: 'change' }]
}

const availableAlgorithms = computed(() => {
    const algMap: Record<string, Array<{ label: string, value: string }>> = {
        'rsa': [{ label: 'RS256', value: 'RS256' }, { label: 'RS384', value: 'RS384' }, { label: 'RS512', value: 'RS512' }],
        'octet': [{ label: 'HS256', value: 'HS256' }, { label: 'HS384', value: 'HS384' }, { label: 'HS512', value: 'HS512' }]
    }
    return algMap[formData.key_type] || []
})

// 密钥类型变更处理
const handleKeyTypeChange = () => {
    // 重置算法为默认值
    formData.algorithm = formData.key_type === 'rsa' ? 'RS256' : 'HS256'
}

// 重置表单
const resetForm = () => {
    formData.key_type = 'rsa'
    formData.algorithm = 'RS256'
    formData.remark = ''
}

// 监听对话框打开
watch(dialogVisible, (newVal) => {
    if (newVal) {
        resetForm()
    }
})

// 取消按钮
const handleCancel = () => {
    dialogVisible.value = false
}

// 提交按钮
const handleSubmit = async () => {
    if (!formRef.value) return
    await formRef.value.validate(async (valid) => {
        if (!valid) return
        submitLoading.value = true
        try {
            emit('submit', {
                key_type: formData.key_type,
                algorithm: formData.algorithm,
                remark: formData.remark
            })
            dialogVisible.value = false
        } catch (error: any) {
            console.error('表单验证失败:', error)
        } finally {
            submitLoading.value = false
        }
    })
}
</script>
