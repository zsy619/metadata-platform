<template>
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="900px" destroy-on-close class="protocol-config-dialog">
        <el-tabs type="border-card" class="protocol-tabs">
            <!-- 标签页 1：基础信息 -->
            <el-tab-pane label="基础信息">
                <div class="tab-content">
                    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="120px">
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
                        <el-form-item label="排序" prop="sort">
                            <el-input-number v-model="formData.sort" :min="0" />
                        </el-form-item>
                        <el-form-item label="启用">
                            <el-switch v-model="formData.is_enabled" />
                        </el-form-item>
                        <el-form-item label="备注" prop="remark">
                            <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入备注信息" />
                        </el-form-item>
                    </el-form>
                </div>
            </el-tab-pane>
            
            <!-- 标签页 2：协议配置 -->
            <el-tab-pane label="协议配置">
                <div class="tab-content protocol-config-content">
                    <el-alert 
                        v-if="!formData.protocol_type" 
                        title="请先在基础信息中选择协议类型" 
                        type="warning" 
                        show-icon
                        class="mb-4"
                    />
                    
                    <template v-else>
                        <el-tabs type="card" class="protocol-sub-tabs">
                            <!-- 基础配置 -->
                            <el-tab-pane label="基础配置">
                                <el-form 
                                    v-if="formData.protocol_type" 
                                    ref="protocolFormRef" 
                                    :model="dynamicForm" 
                                    :rules="getProtocolRules()" 
                                    label-width="150px"
                                >
                                    <el-form-item 
                                        v-for="field in getBasicFields(formData.protocol_type)" 
                                        :key="field.key" 
                                        :label="field.label" 
                                        :prop="field.key"
                                    >
                                        <el-input 
                                            v-if="field.type === 'text'" 
                                            v-model="dynamicForm[field.key]" 
                                            :placeholder="field.placeholder" 
                                        />
                                        <el-input 
                                            v-else-if="field.type === 'password'" 
                                            v-model="dynamicForm[field.key]" 
                                            type="password" 
                                            :placeholder="field.placeholder"
                                            show-password
                                        />
                                        <el-input-number 
                                            v-else-if="field.type === 'number'" 
                                            v-model="dynamicForm[field.key]" 
                                            :min="field.min || 0" 
                                            :max="field.max || 999999"
                                            style="width: 100%"
                                        />
                                    </el-form-item>
                                </el-form>
                            </el-tab-pane>
                            
                            <!-- 端点配置 -->
                            <el-tab-pane label="端点配置">
                                <el-form 
                                    v-if="formData.protocol_type" 
                                    ref="endpointFormRef" 
                                    :model="dynamicForm" 
                                    :rules="getProtocolRules()" 
                                    label-width="150px"
                                >
                                    <el-form-item 
                                        v-for="field in getEndpointFields(formData.protocol_type)" 
                                        :key="field.key" 
                                        :label="field.label" 
                                        :prop="field.key"
                                    >
                                        <el-input v-model="dynamicForm[field.key]" :placeholder="field.placeholder">
                                            <template #append v-if="field.default">
                                                <el-button @click="setDefaultValue(field.key, field.default)">默认</el-button>
                                            </template>
                                        </el-input>
                                    </el-form-item>
                                </el-form>
                            </el-tab-pane>
                            
                            <!-- 高级配置 -->
                            <el-tab-pane label="高级配置">
                                <el-form 
                                    v-if="formData.protocol_type" 
                                    ref="advancedFormRef" 
                                    :model="dynamicForm" 
                                    label-width="150px"
                                >
                                    <el-form-item 
                                        v-for="field in getAdvancedFields(formData.protocol_type)" 
                                        :key="field.key" 
                                        :label="field.label" 
                                        :prop="field.key"
                                    >
                                        <el-input 
                                            v-if="field.type === 'text'" 
                                            v-model="dynamicForm[field.key]" 
                                            :placeholder="field.placeholder" 
                                        />
                                        <el-input 
                                            v-else-if="field.type === 'password'" 
                                            v-model="dynamicForm[field.key]" 
                                            type="password" 
                                            :placeholder="field.placeholder"
                                            show-password
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
                                </el-form>
                            </el-tab-pane>
                        </el-tabs>
                    </template>
                </div>
            </el-tab-pane>
        </el-tabs>
        
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="handleCancel">取消</el-button>
                <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
            </div>
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
const protocolFormRef = ref<FormInstance>()
const endpointFormRef = ref<FormInstance>()
const advancedFormRef = ref<FormInstance>()

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

// 表单验证规则
const formRules: FormRules = {
    config_name: [
        { required: true, message: '请输入配置名称', trigger: 'blur' }
    ],
    protocol_type: [
        { required: true, message: '请选择协议类型', trigger: 'change' }
    ],
    sort: [
        { type: 'number', message: '排序必须为数字', trigger: 'blur' }
    ]
}

// 协议字段配置
interface ProtocolField {
    key: string
    label: string
    type: 'text' | 'password' | 'number' | 'boolean' | 'select'
    placeholder: string
    default?: any
    required?: boolean
    min?: number
    max?: number
    disabled?: boolean
    options?: Array<{ label: string; value: any }>
    category?: 'basic' | 'endpoint' | 'advanced'
}

const protocolFields: Record<string, ProtocolField[]> = {
    oidc: [
        { key: 'client_id', label: '客户端 ID', type: 'text', placeholder: '请输入客户端 ID', required: true, category: 'basic' },
        { key: 'client_secret', label: '客户端密钥', type: 'password', placeholder: '请输入客户端密钥', required: true, category: 'basic' },
        { key: 'issuer', label: '发行方 URL', type: 'text', placeholder: '请输入发行方 URL', required: true, default: 'http://localhost:8080', category: 'basic' },
        { key: 'scope', label: '作用域', type: 'text', placeholder: '请输入作用域，多个用空格分隔', default: 'openid profile email', category: 'basic' },
        
        { key: 'authorization_endpoint', label: '授权端点', type: 'text', placeholder: '请输入授权端点', default: 'http://localhost:8080/api/sso/auth/oidc/authorize', category: 'endpoint' },
        { key: 'token_endpoint', label: '令牌端点', type: 'text', placeholder: '请输入令牌端点', default: 'http://localhost:8080/api/sso/auth/oidc/token', category: 'endpoint' },
        { key: 'userinfo_endpoint', label: '用户信息端点', type: 'text', placeholder: '请输入用户信息端点', default: 'http://localhost:8080/api/sso/auth/oidc/userinfo', category: 'endpoint' },
        { key: 'redirect_uri', label: '重定向 URI', type: 'text', placeholder: '请输入重定向 URI', required: true, default: 'http://localhost:3000/callback', category: 'endpoint' },
        
        { key: 'end_session_endpoint', label: '登出端点', type: 'text', placeholder: '请输入登出端点', default: 'http://localhost:8080/api/sso/auth/oidc/logout', category: 'advanced' }
    ],
    
    saml: [
        { key: 'entity_id', label: '实体 ID', type: 'text', placeholder: '请输入实体 ID', required: true, default: 'http://localhost:8080/saml/metadata', category: 'basic' },
        { key: 'x509_certificate', label: 'X509 证书', type: 'text', placeholder: '请输入 X509 证书', required: true, category: 'basic' },
        
        { key: 'single_sign_on_url', label: '单点登录 URL', type: 'text', placeholder: '请输入单点登录 URL', required: true, default: 'http://localhost:8080/api/sso/auth/saml/sso', category: 'endpoint' },
        { key: 'acs_url', label: '断言消费服务 URL', type: 'text', placeholder: '请输入断言消费服务 URL', required: true, default: 'http://localhost:8080/api/sso/auth/saml/acs', category: 'endpoint' },
        
        { key: 'single_logout_url', label: '单点登出 URL', type: 'text', placeholder: '请输入单点登出 URL', default: 'http://localhost:8080/api/sso/auth/saml/slo', category: 'advanced' }
    ],
    
    ldap: [
        { key: 'server_url', label: '服务器 URL', type: 'text', placeholder: '请输入 LDAP 服务器 URL', required: true, default: 'ldap://localhost', category: 'basic' },
        { key: 'bind_dn', label: '绑定 DN', type: 'text', placeholder: '请输入绑定 DN', required: true, category: 'basic' },
        { key: 'bind_password', label: '绑定密码', type: 'password', placeholder: '请输入绑定密码', required: true, category: 'basic' },
        { key: 'base_dn', label: '基础 DN', type: 'text', placeholder: '请输入基础 DN', required: true, category: 'basic' },
        { key: 'port', label: '端口', type: 'number', placeholder: '请输入端口', default: 389, category: 'basic' },
        
        { key: 'user_filter', label: '用户过滤器', type: 'text', placeholder: '请输入用户过滤器', default: '(uid={username})', category: 'advanced' },
        { key: 'group_filter', label: '组过滤器', type: 'text', placeholder: '请输入组过滤器', category: 'advanced' },
        { key: 'use_ssl', label: '使用 SSL', type: 'boolean', placeholder: '', default: false, category: 'advanced' }
    ],
    
    cas: [
        { key: 'server_url', label: 'CAS 服务器 URL', type: 'text', placeholder: '请输入 CAS 服务器 URL', required: true, default: 'http://localhost:8080', category: 'basic' },
        { key: 'service_url', label: '服务 URL', type: 'text', placeholder: '请输入服务 URL', required: true, default: 'http://localhost:3000', category: 'basic' },
        
        { key: 'login_url', label: '登录 URL', type: 'text', placeholder: '请输入登录 URL', default: 'http://localhost:8080/api/sso/auth/cas/login', category: 'endpoint' },
        { key: 'logout_url', label: '登出 URL', type: 'text', placeholder: '请输入登出 URL', default: 'http://localhost:8080/api/sso/auth/cas/logout', category: 'endpoint' },
        { key: 'validate_url', label: '验证 URL', type: 'text', placeholder: '请输入验证 URL', default: 'http://localhost:8080/api/sso/auth/cas/validate', category: 'endpoint' }
    ]
}

// 获取协议对应的字段配置
const getProtocolFields = (protocolType: string): ProtocolField[] => {
    return protocolFields[protocolType] || []
}

// 按分类获取字段
const getFieldsByCategory = (protocolType: string, category: 'basic' | 'endpoint' | 'advanced'): ProtocolField[] => {
    const fields = getProtocolFields(protocolType)
    return fields.filter(f => f.category === category || (!f.category && category === 'basic'))
}

const getBasicFields = (protocolType: string) => getFieldsByCategory(protocolType, 'basic')
const getEndpointFields = (protocolType: string) => getFieldsByCategory(protocolType, 'endpoint')
const getAdvancedFields = (protocolType: string) => getFieldsByCategory(protocolType, 'advanced')

// 获取协议验证规则
const getProtocolRules = (): FormRules => {
    const rules: FormRules = {}
    const fields = getProtocolFields(formData.protocol_type!)
    
    fields.forEach(field => {
        if (field.required) {
            rules[field.key] = [
                { required: true, message: `${field.label}不能为空`, trigger: 'blur' }
            ]
        }
    })
    
    return rules
}

// 设置默认值
const setDefaultValue = (key: string, value: any) => {
    dynamicForm[key] = value
}

// 协议类型变更处理
const handleProtocolTypeChange = () => {
    Object.keys(dynamicForm).forEach(key => {
        delete dynamicForm[key]
    })
    
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
    
    Object.keys(dynamicForm).forEach(key => {
        delete dynamicForm[key]
    })
    
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

// 提交按钮 - 验证所有表单
const handleSubmit = async () => {
    // 验证基础信息表单
    if (formRef.value) {
        try {
            await formRef.value.validate()
        } catch (error) {
            return
        }
    }
    
    // 验证协议配置表单（如果已选择协议类型）
    if (formData.protocol_type) {
        if (protocolFormRef.value) {
            try {
                await protocolFormRef.value.validate()
            } catch (error) {
                // 切换到基础配置标签页
                return
            }
        }
        
        if (endpointFormRef.value) {
            try {
                await endpointFormRef.value.validate()
            } catch (error) {
                // 切换到端点配置标签页
                return
            }
        }
    }
    
    submitLoading.value = true
    try {
        const configData = JSON.stringify(dynamicForm, null, 2)
        
        emit('submit', {
            ...formData,
            config_data: configData
        } as Partial<SsoProtocolConfig> & { config_data: string })
        
        dialogVisible.value = false
    } catch (error: any) {
        console.error('提交失败:', error)
    } finally {
        submitLoading.value = false
    }
}
</script>

<style scoped>
.protocol-config-dialog :deep(.el-dialog__body) {
    padding: 0;
    max-height: 600px;
    overflow: hidden;
}

.protocol-config-dialog :deep(.el-dialog__footer) {
    padding: 15px 20px;
    border-top: 1px solid #e4e7ed;
}

.protocol-tabs {
    border: none;
}

.protocol-tabs :deep(.el-tabs__header) {
    margin: 0;
}

.protocol-tabs :deep(.el-tabs__content) {
    padding: 20px;
    height: 480px;
    overflow-y: auto;
}

.tab-content {
    padding: 10px;
}

.protocol-config-content {
    padding: 0;
}

.protocol-sub-tabs {
    border: none;
}

.protocol-sub-tabs :deep(.el-tabs__header) {
    margin-bottom: 15px;
}

.protocol-sub-tabs :deep(.el-tabs__content) {
    padding: 15px;
    background-color: #f9f9f9;
    border-radius: 4px;
}

.mb-4 {
    margin-bottom: 16px;
}

.dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
}

/* 滚动条样式 */
.protocol-tabs :deep(.el-tabs__content::-webkit-scrollbar) {
    width: 6px;
}

.protocol-tabs :deep(.el-tabs__content::-webkit-scrollbar-thumb) {
    background-color: #c1c1c1;
    border-radius: 3px;
}

.protocol-tabs :deep(.el-tabs__content::-webkit-scrollbar-thumb:hover) {
    background-color: #a8a8a8;
}
</style>
