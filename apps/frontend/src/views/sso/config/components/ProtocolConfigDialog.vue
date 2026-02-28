<template>
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px" destroy-on-close>
        <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
            <el-form-item label="配置名称" prop="config_name">
                <el-input v-model="formData.config_name" placeholder="请输入配置名称" />
            </el-form-item>
            <el-form-item label="协议类型" prop="protocol_type">
                <el-select v-model="formData.protocol_type" placeholder="请选择协议类型" style="width: 100%" @change="handleProtocolTypeChange">
                    <el-option label="OIDC" value="oidc" />
                    <el-option label="SAML" value="saml" />
                    <el-option label="LDAP" value="ldap" />
                    <el-option label="CAS" value="cas" />
                </el-select>
            </el-form-item>
            
            <!-- 动态协议配置表单 -->
            <div v-if="formData.protocol_type" class="protocol-form">
                <el-form-item v-for="field in getProtocolFields(formData.protocol_type!)" :key="field.key" :label="field.label" :prop="`dynamicForm.${field.key}`">
                    <el-input 
                        v-if="field.type === 'text'" 
                        v-model="dynamicForm[field.key]" 
                        :placeholder="field.placeholder" 
                        :disabled="field.disabled"
                    />
                    <el-input 
                        v-else-if="field.type === 'password'" 
                        v-model="dynamicForm[field.key]" 
                        type="password" 
                        :placeholder="field.placeholder"
                    />
                    <el-input-number 
                        v-else-if="field.type === 'number'" 
                        v-model="dynamicForm[field.key]" 
                        :min="field.min || 0" 
                        :max="field.max || 999999"
                    />
                    <el-switch 
                        v-else-if="field.type === 'boolean'" 
                        v-model="dynamicForm[field.key]"
                    />
                    <el-select 
                        v-else-if="field.type === 'select'" 
                        v-model="dynamicForm[field.key]" 
                        style="width: 100%" 
                        :placeholder="field.placeholder"
                    >
                        <el-option 
                            v-for="option in field.options" 
                            :key="option.value" 
                            :label="option.label" 
                            :value="option.value"
                        />
                    </el-select>
                </el-form-item>
            </div>
            
            <el-form-item label="排序" prop="sort">
                <el-input-number v-model="formData.sort" :min="0" />
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
import type { SsoProtocolConfig } from '@/types/sso'
import type { FormInstance, FormRules } from 'element-plus'
import { computed, reactive, ref, watch } from 'vue'

// Props
interface Props {
    modelValue: boolean
    editData?: SsoProtocolConfig | null
}

const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
    editData: null
})

// Emits
const emit = defineEmits<{
    'update:modelValue': [value: boolean]
    'submit': [data: Partial<SsoProtocolConfig> & { config_data: string }]
}>()

// 内部状态
const dialogVisible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
})

const dialogTitle = computed(() => props.editData ? '编辑协议配置' : '新增协议配置')
const submitLoading = ref(false)
const formRef = ref<FormInstance>()

const formData = reactive<Partial<SsoProtocolConfig>>({
    config_name: '',
    protocol_type: 'oidc',
    config_data: '',
    is_enabled: true,
    sort: 0,
    remark: ''
})

// 动态表单数据
const dynamicForm = reactive<Record<string, any>>({})

const formRules: FormRules = {
    config_name: [{ required: true, message: '请输入配置名称', trigger: 'blur' }],
    protocol_type: [{ required: true, message: '请选择协议类型', trigger: 'change' }]
}

// 协议字段配置
const protocolFields: Record<string, Array<{
    key: string
    label: string
    type: 'text' | 'password' | 'number' | 'boolean' | 'select'
    placeholder: string
    default?: any
    min?: number
    max?: number
    disabled?: boolean
    options?: Array<{ label: string; value: any }>
}>> = {
    oidc: [
        { key: 'client_id', label: '客户端 ID', type: 'text', placeholder: '请输入客户端 ID' },
        { key: 'client_secret', label: '客户端密钥', type: 'password', placeholder: '请输入客户端密钥' },
        { key: 'issuer', label: '发行方 URL', type: 'text', placeholder: '请输入发行方 URL' },
        { key: 'authorization_endpoint', label: '授权端点', type: 'text', placeholder: '请输入授权端点' },
        { key: 'token_endpoint', label: '令牌端点', type: 'text', placeholder: '请输入令牌端点' },
        { key: 'userinfo_endpoint', label: '用户信息端点', type: 'text', placeholder: '请输入用户信息端点' },
        { key: 'redirect_uri', label: '重定向 URI', type: 'text', placeholder: '请输入重定向 URI' },
        { key: 'scope', label: '作用域', type: 'text', placeholder: '请输入作用域，多个用空格分隔', default: 'openid profile email' }
    ],
    saml: [
        { key: 'entity_id', label: '实体 ID', type: 'text', placeholder: '请输入实体 ID' },
        { key: 'single_sign_on_url', label: '单点登录 URL', type: 'text', placeholder: '请输入单点登录 URL' },
        { key: 'single_logout_url', label: '单点登出 URL', type: 'text', placeholder: '请输入单点登出 URL' },
        { key: 'x509_certificate', label: 'X509 证书', type: 'text', placeholder: '请输入 X509 证书' },
        { key: 'acs_url', label: '断言消费服务 URL', type: 'text', placeholder: '请输入断言消费服务 URL' }
    ],
    ldap: [
        { key: 'server_url', label: '服务器 URL', type: 'text', placeholder: '请输入 LDAP 服务器 URL' },
        { key: 'bind_dn', label: '绑定 DN', type: 'text', placeholder: '请输入绑定 DN' },
        { key: 'bind_password', label: '绑定密码', type: 'password', placeholder: '请输入绑定密码' },
        { key: 'base_dn', label: '基础 DN', type: 'text', placeholder: '请输入基础 DN' },
        { key: 'user_filter', label: '用户过滤器', type: 'text', placeholder: '请输入用户过滤器', default: '(uid={username})' },
        { key: 'group_filter', label: '组过滤器', type: 'text', placeholder: '请输入组过滤器' },
        { key: 'port', label: '端口', type: 'number', placeholder: '请输入端口', default: 389 },
        { key: 'use_ssl', label: '使用 SSL', type: 'boolean', placeholder: '', default: false }
    ],
    cas: [
        { key: 'server_url', label: 'CAS 服务器 URL', type: 'text', placeholder: '请输入 CAS 服务器 URL' },
        { key: 'service_url', label: '服务 URL', type: 'text', placeholder: '请输入服务 URL' },
        { key: 'login_url', label: '登录 URL', type: 'text', placeholder: '请输入登录 URL' },
        { key: 'logout_url', label: '登出 URL', type: 'text', placeholder: '请输入登出 URL' },
        { key: 'validate_url', label: '验证 URL', type: 'text', placeholder: '请输入验证 URL' }
    ]
}

// 获取协议对应的字段配置
const getProtocolFields = (protocolType: string) => {
    return protocolFields[protocolType] || []
}

// 协议类型变更处理
const handleProtocolTypeChange = () => {
    // 重置动态表单
    Object.keys(dynamicForm).forEach(key => {
        delete dynamicForm[key]
    })
    
    // 设置默认值
    const fields = getProtocolFields(formData.protocol_type!)
    fields.forEach(field => {
        if (field.default !== undefined) {
            dynamicForm[field.key] = field.default
        }
    })
}

// 从 JSON 数据加载到动态表单
const loadConfigDataToForm = () => {
    if (formData.config_data) {
        try {
            const configData = JSON.parse(formData.config_data)
            Object.assign(dynamicForm, configData)
        } catch (error) {
            console.error('解析配置数据失败:', error)
        }
    }
}

// 重置表单
const resetForm = () => {
    formData.id = undefined
    formData.config_name = ''
    formData.protocol_type = 'oidc'
    formData.config_data = ''
    formData.is_enabled = true
    formData.sort = 0
    formData.remark = ''
    
    // 重置动态表单
    Object.keys(dynamicForm).forEach(key => {
        delete dynamicForm[key]
    })
    
    // 设置默认值
    handleProtocolTypeChange()
}

// 监听编辑数据变化
watch(() => props.editData, (newData) => {
    if (newData) {
        Object.assign(formData, newData)
        loadConfigDataToForm()
    } else {
        resetForm()
    }
}, { immediate: true })

// 监听对话框打开
watch(dialogVisible, (newVal) => {
    if (newVal && !props.editData) {
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
            // 自动生成 JSON 配置数据
            const configData = JSON.stringify(dynamicForm, null, 2)
            
            emit('submit', {
                ...formData,
                config_data: configData
            } as Partial<SsoProtocolConfig> & { config_data: string })
            
            dialogVisible.value = false
        } catch (error: any) {
            console.error('表单验证失败:', error)
        } finally {
            submitLoading.value = false
        }
    })
}
</script>

<style scoped>
.protocol-form {
    margin: 10px 0 20px 0;
    padding: 15px;
    background-color: #f9f9f9;
    border-radius: 8px;
    border: 1px solid #eaeaea;
}

.protocol-form :deep(.el-form-item) {
    margin-bottom: 12px;
}
</style>
