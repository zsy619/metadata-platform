<template>
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px" destroy-on-close @close="handleClose">
        <el-form
            ref="formRef"
            :model="formData"
            :rules="formRules"
            label-width="120px"
            class="client-form"
        >
            <el-row :gutter="20">
                <el-col :span="12">
                    <el-form-item label="客户端名称" prop="client_name">
                        <el-input v-model="formData.client_name" placeholder="请输入客户端名称" />
                    </el-form-item>
                </el-col>
                <el-col :span="12">
                    <el-form-item label="客户端类型" prop="client_type">
                        <el-select v-model="formData.client_type" placeholder="请选择客户端类型" style="width: 100%">
                            <el-option label="Web 应用" value="web" />
                            <el-option label="SPA 应用" value="spa" />
                            <el-option label="移动应用" value="mobile" />
                            <el-option label="后端服务" value="backend" />
                        </el-select>
                    </el-form-item>
                </el-col>
            </el-row>
            <el-row :gutter="20">
                <el-col :span="12">
                    <el-form-item label="客户端 ID" prop="client_id">
                        <el-input v-model="formData.client_id" placeholder="系统自动生成或自定义" />
                    </el-form-item>
                </el-col>
                <el-col :span="12">
                    <el-form-item label="客户端密钥" prop="client_secret">
                        <el-input v-model="formData.client_secret" type="password" show-password placeholder="请输入客户端密钥" />
                    </el-form-item>
                </el-col>
            </el-row>
            <el-form-item label="回调地址" prop="redirect_uris">
                <el-input
                    v-model="formData.redirect_uris"
                    type="textarea"
                    :rows="3"
                    placeholder="多个回调地址用逗号分隔，如：https://app1.com/callback,https://app2.com/callback"
                />
            </el-form-item>
            <el-form-item label="登出地址" prop="post_logout_redirect_uris">
                <el-input
                    v-model="formData.post_logout_redirect_uris"
                    type="textarea"
                    :rows="2"
                    placeholder="多个登出地址用逗号分隔"
                />
            </el-form-item>
            <el-form-item label="是否启用" prop="is_enabled">
                <el-switch v-model="formData.is_enabled" />
            </el-form-item>
            <el-form-item label="描述" prop="app_description">
                <el-input
                    v-model="formData.app_description"
                    type="textarea"
                    :rows="2"
                    placeholder="请输入客户端描述"
                />
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="handleClose">取消</el-button>
            <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { SsoClient } from '@/types/sso'

interface Props {
    visible: boolean
    title: string
    data?: Partial<SsoClient>
    submitLoading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    data: undefined,
    submitLoading: false
})

const emit = defineEmits<{
    (e: 'update:visible', value: boolean): void
    (e: 'submit', data: Partial<SsoClient>): void
}>()

const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
})

const dialogTitle = computed(() => props.title)
const submitLoading = computed(() => props.submitLoading)

const formRef = ref<FormInstance>()

const formData = reactive<Partial<SsoClient>>({
    client_name: '',
    client_id: '',
    client_secret: '',
    client_type: 'web',
    redirect_uris: '',
    post_logout_redirect_uris: '',
    is_enabled: true,
    app_description: ''
})

const formRules: FormRules = {
    client_name: [{ required: true, message: '请输入客户端名称', trigger: 'blur' }],
    client_type: [{ required: true, message: '请选择客户端类型', trigger: 'change' }],
    redirect_uris: [{ required: true, message: '请输入回调地址', trigger: 'blur' }]
}

// 监听数据变化，填充表单
watch(() => props.data, (newData) => {
    if (newData) {
        Object.assign(formData, newData)
    }
}, { immediate: true })

// 监听可见性变化，重置表单
watch(() => props.visible, (newVal) => {
    if (!newVal && formRef.value) {
        formRef.value.resetFields()
    }
})

const resetForm = () => {
    formData.id = undefined
    formData.client_name = ''
    formData.client_id = ''
    formData.client_secret = ''
    formData.client_type = 'web'
    formData.redirect_uris = ''
    formData.post_logout_redirect_uris = ''
    formData.is_enabled = true
    formData.app_description = ''
}

const handleClose = () => {
    if (formRef.value) {
        formRef.value.resetFields()
    }
    emit('update:visible', false)
}

const handleSubmit = async () => {
    if (!formRef.value) return
    await formRef.value.validate(async (valid) => {
        if (!valid) return
        emit('submit', { ...formData })
    })
}

// 暴露重置方法给父组件
defineExpose({
    resetForm
})
</script>

<style scoped>
.client-form :deep(.el-form-item) {
    margin-bottom: 22px;
}
</style>
