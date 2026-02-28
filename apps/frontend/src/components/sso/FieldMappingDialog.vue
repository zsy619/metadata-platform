<template>
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
        <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
            <el-form-item label="映射名称" prop="mapping_name">
                <el-input v-model="formData.mapping_name" placeholder="请输入映射名称" />
            </el-form-item>
            <el-form-item label="源字段" prop="source_field">
                <el-input v-model="formData.source_field" placeholder="如: email, username" />
            </el-form-item>
            <el-form-item label="目标字段" prop="target_field">
                <el-select v-model="formData.target_field" placeholder="请选择目标字段" style="width: 100%">
                    <el-option label="用户名" value="username" />
                    <el-option label="邮箱" value="email" />
                    <el-option label="手机号" value="phone" />
                    <el-option label="昵称" value="nickname" />
                    <el-option label="姓名" value="real_name" />
                    <el-option label="部门" value="department" />
                    <el-option label="职位" value="position" />
                </el-select>
            </el-form-item>
            <el-form-item label="字段类型" prop="field_type">
                <el-select v-model="formData.field_type" style="width: 100%">
                    <el-option label="字符串" value="string" />
                    <el-option label="数字" value="number" />
                    <el-option label="布尔" value="boolean" />
                    <el-option label="日期" value="date" />
                </el-select>
            </el-form-item>
            <el-form-item label="默认值" prop="default_value">
                <el-input v-model="formData.default_value" placeholder="当源字段为空时的默认值" />
            </el-form-item>
            <el-form-item label="必填">
                <el-switch v-model="formData.is_required" />
            </el-form-item>
            <el-form-item label="启用">
                <el-switch v-model="formData.is_enabled" />
            </el-form-item>
            <el-form-item label="备注" prop="remark">
                <el-input v-model="formData.remark" type="textarea" :rows="2" />
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="handleCancel">取消</el-button>
            <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import type { SsoFieldMapping } from '@/types/sso'

// Props
const props = defineProps<{
    visible: boolean
    title: string
    formData: Partial<SsoFieldMapping>
    loading: boolean
}>()

// Emits
const emit = defineEmits<{
    (e: 'update:visible', value: boolean): void
    (e: 'submit', data: Partial<SsoFieldMapping>): void
}>()

const dialogVisible = ref(props.visible)
const dialogTitle = ref(props.title)
const submitLoading = ref(props.loading)
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive<Partial<SsoFieldMapping>>({ ...props.formData })

// 表单验证规则
const formRules: FormRules = {
    mapping_name: [{ required: true, message: '请输入映射名称', trigger: 'blur' }],
    source_field: [{ required: true, message: '请输入源字段', trigger: 'blur' }],
    target_field: [{ required: true, message: '请选择目标字段', trigger: 'change' }]
}

// 监听props变化
watch(() => props.visible, (newVal) => {
    dialogVisible.value = newVal
})

watch(() => props.title, (newVal) => {
    dialogTitle.value = newVal
})

watch(() => props.loading, (newVal) => {
    submitLoading.value = newVal
})

watch(() => props.formData, (newVal) => {
    Object.assign(formData, newVal)
}, { deep: true })

// 处理取消
const handleCancel = () => {
    emit('update:visible', false)
}

// 处理提交
const handleSubmit = async () => {
    if (!formRef.value) return
    await formRef.value.validate(async (valid) => {
        if (!valid) return
        emit('submit', formData)
    })
}
</script>

<style scoped>
/* 组件样式 */
</style>