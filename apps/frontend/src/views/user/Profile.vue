<template>
    <div class="user-profile">
        <el-row :gutter="20" class="profile-row">
            <!-- 左侧概览 -->
            <el-col :span="8">
                <el-card class="box-card full-height-card">
                    <div class="user-header">
                        <el-avatar :size="100" :src="userInfo.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" />
                        <div class="user-name">{{ userInfo.name }}</div>
                        <div class="user-role">
                            <el-tag :type="userInfo.kind === 1 ? 'danger' : 'success'" size="small">
                                {{ formatKind(userInfo.kind) }}
                            </el-tag>
                        </div>
                    </div>
                    <div class="user-bio">
                        <div class="bio-item">
                            <div class="label">最后登录时间:</div>
                            <div class="value">{{ formatDate(userInfo.last_login_time) }}</div>
                        </div>
                        <div class="bio-item">
                            <div class="label">最后登录IP:</div>
                            <div class="value">{{ userInfo.last_ip || '-' }}</div>
                        </div>
                        <div class="bio-item">
                            <div class="label">注册时间:</div>
                            <div class="value">{{ formatDate(userInfo.create_at) }}</div>
                        </div>
                        <div class="bio-item">
                            <div class="label">账号状态:</div>
                            <div class="value">
                                <el-tag :type="userInfo.state === 1 ? 'success' : 'danger'" size="small">
                                    {{ userInfo.state === 1 ? '启用' : '禁用' }}
                                </el-tag>
                            </div>
                        </div>
                    </div>
                </el-card>
            </el-col>
            <!-- 右侧详情 -->
            <el-col :span="16">
                <el-card class="full-height-card">
                    <template #header>
                        <div class="card-header">
                            <span>资料详情</span>
                        </div>
                    </template>
                    <el-tabs v-model="activeTab">
                        <!-- 基本资料 -->
                        <el-tab-pane label="基本资料" name="info">
                            <el-descriptions :column="1" border class="user-info-desc">
                                <el-descriptions-item label="用户ID">{{ userInfo.id }}</el-descriptions-item>
                                <el-descriptions-item label="登录账号">{{ userInfo.account }}</el-descriptions-item>
                                <el-descriptions-item label="用户编号">{{ userInfo.code || '-' }}</el-descriptions-item>
                                <el-descriptions-item label="性别">{{ userInfo.sex }}</el-descriptions-item>
                                <el-descriptions-item label="手机号码">{{ userInfo.mobile || '-' }}</el-descriptions-item>
                                <el-descriptions-item label="电子邮箱">{{ userInfo.email || '-' }}</el-descriptions-item>
                                <el-descriptions-item label="身份证号">{{ userInfo.id_card || '-' }}</el-descriptions-item>
                                <el-descriptions-item label="备注信息">{{ userInfo.remark || '-' }}</el-descriptions-item>
                            </el-descriptions>
                        </el-tab-pane>
                        <!-- 拥有角色 (新增) -->
                        <el-tab-pane label="拥有角色" name="role">
                            <el-table :data="userInfo.roles" style="width: 100%" border>
                                <el-table-column prop="role_name" label="角色名称" />
                                <el-table-column prop="role_code" label="角色编码" width="180" />
                                <el-table-column prop="remark" label="备注" />
                            </el-table>
                        </el-tab-pane>
                        <!-- 所属组织 -->
                        <el-tab-pane label="所属组织" name="organization">
                            <el-table :data="userInfo.organizations" style="width: 100%" border>
                                <el-table-column prop="unit_name" label="组织名称" />
                                <el-table-column prop="unit_code" label="组织编码" width="180" />
                                <el-table-column prop="remark" label="备注" />
                            </el-table>
                        </el-tab-pane>
                        <!-- 担任职位 -->
                        <el-tab-pane label="担任职位" name="position">
                            <el-table :data="userInfo.positions" style="width: 100%" border>
                                <el-table-column prop="pos_name" label="职位名称" />
                                <el-table-column prop="pos_code" label="职位编码" width="180" />
                                <el-table-column prop="grade" label="职级" width="100" align="center" />
                                <el-table-column prop="remark" label="备注" />
                            </el-table>
                        </el-tab-pane>
                        <!-- 修改密码 -->
                        <el-tab-pane label="修改密码" name="password">
                            <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="100px">
                                <el-form-item label="原密码" prop="oldPassword">
                                    <el-input v-model="passwordForm.oldPassword" type="password" show-password placeholder="请输入原密码" />
                                </el-form-item>
                                <el-form-item label="新密码" prop="newPassword">
                                    <el-input v-model="passwordForm.newPassword" type="password" show-password placeholder="请输入新密码" />
                                </el-form-item>
                                <el-form-item label="确认密码" prop="confirmPassword">
                                    <el-input v-model="passwordForm.confirmPassword" type="password" show-password placeholder="请再次输入新密码" />
                                </el-form-item>
                                <el-form-item>
                                    <el-button type="primary" :loading="loading" @click="handleUpdatePassword">
                                        修改密码
                                    </el-button>
                                    <el-button @click="resetForm(passwordFormRef)">重置</el-button>
                                </el-form-item>
                            </el-form>
                            <div class="password-rule-info">
                                <h4>密码设置要求：</h4>
                                <p>1. 长度在 8 到 20 个字符之间</p>
                                <p>2. 必须包含大写字母、小写字母、数字和特殊字符</p>
                            </div>
                        </el-tab-pane>
                    </el-tabs>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>
<script setup lang="ts">
import { getUserProfile, updateUserPassword } from '@/api/auth'
import type { User } from '@/types/user'
import { ElMessage, type FormInstance } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

const activeTab = ref('info')
const loading = ref(false)

const userInfo = reactive<User>({
    id: '',
    tenant_id: '0',
    account: '',
    name: '',
    kind: 99,
    state: 1,
    code: '',
    sex: '',
    mobile: '',
    email: '',
    id_card: '',
    create_at: '',
    update_at: '', // Added required field
    remark: '',
    avatar: '',
    last_login_time: '',
    last_ip: '',
    roles: [],
    organizations: [],
    positions: [],
    // Add missing required fields from User/BaseEntity to satisfy type
    account_id: '',
    svc_code: '',
    salt: '',
    unit_id: '',
    school: '',
    class: '',
    endTime: '', // types/user.ts uses camelCase 'endTime'? Check Step 1581. Yes 'endTime'. Wait, I fixed User to snake_case?
    // Let's re-check Step 1571/1581 diffs.
    // Step 1571 updated User interface.
    // It kept 'endTime: string;' at line 58?
    // Let me check my Step 1571 diff again. 
    // It updated 'firstLogin' -> 'first_login'.
    // Did I update 'endTime'? 
    // Step 1571 Replacement: 
    //   StartLine: 66 -> Target: firstLogin... 
    //   It didn't touch line 58.
    // So 'endTime' is still 'endTime' in types/user.ts?
    // Let's check Step 1581 view: "58:     endTime: string;"
    // Yes. But backend probably sends snake_case `end_time`?
    // If backend sends `end_time`, and type expects `endTime`, we have a mismatch.
    // I should fix `endTime` in `types/user.ts` too if I want perfection.
    // But for now, let's just initialize what's needed.
    // Warning: `userInfo` initialization might fail type check if I miss fields.
    // But `Index.vue` uses `reactive<UserInfo>` previously.
    // If I use `reactive<User>`, I need all non-optional fields.
    // `User` extends `BaseEntity`. `BaseEntity` has `id, tenant_id, create_at, update_at`.
    // My replacement content below attempts to satisfy `User`.
    sort: 0,
    first_login: 0,
    login_error_count: 0,
    create_by: '',
    update_by: ''
})

const passwordFormRef = ref<FormInstance>()
const passwordForm = reactive({
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
})

const formatDate = (dateStr: string | undefined) => {
    if (!dateStr || dateStr === '' || dateStr === '0001-01-01 00:00:00') return '-'
    const date = new Date(dateStr)
    if (isNaN(date.getTime())) return dateStr

    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    const seconds = String(date.getSeconds()).padStart(2, '0')

    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

const validatePass = (_rule: any, value: string, callback: any) => {
    if (value === '') {
        callback(new Error('请输入密码'))
    } else {
        // 强制复杂度检查：8-20位，包含大小写字母、数字和特殊字符
        const regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,20}$/
        if (!regex.test(value)) {
            callback(new Error('密码长度需8-20位，且包含大小写字母、数字和特殊字符'))
        } else if (passwordForm.confirmPassword !== '') {
            if (!passwordFormRef.value) return
            passwordFormRef.value.validateField('confirmPassword', () => { })
        }
        callback()
    }
}

const validatePass2 = (_rule: any, value: string, callback: any) => {
    if (value === '') {
        callback(new Error('请再次输入密码'))
    } else if (value !== passwordForm.newPassword) {
        callback(new Error('两次输入密码不一致!'))
    } else {
        callback()
    }
}

const passwordRules = reactive({
    oldPassword: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
    newPassword: [{ validator: validatePass, trigger: 'blur' }],
    confirmPassword: [{ validator: validatePass2, trigger: 'blur' }]
})

const formatKind = (kind: number) => {
    const map: Record<number, string> = {
        1: '超级管理员',
        2: '子管理员',
        99: '普通用户'
    }
    return map[kind] || '其他'
}

const fetchUserInfo = async () => {
    try {
        const res = await getUserProfile()
        const data = res.data

        // 基础信息
        userInfo.id = data.id
        userInfo.account = data.account
        userInfo.name = data.name
        userInfo.kind = data.kind
        userInfo.state = data.state

        // 详细信息
        userInfo.code = data.code
        userInfo.sex = data.sex || '男'
        userInfo.mobile = data.mobile
        userInfo.email = data.email
        userInfo.id_card = data.id_card
        userInfo.create_at = data.create_at
        userInfo.remark = data.remark
        userInfo.avatar = data.avatar
        userInfo.last_login_time = data.last_login_time
        userInfo.last_ip = data.last_ip

        // 关联信息
        userInfo.roles = data.roles || []
        userInfo.organizations = data.organizations || []
        userInfo.positions = data.positions || []

    } catch (error) {
        console.error('获取用户信息失败', error)
    }
}

const handleUpdatePassword = async () => {
    if (!passwordFormRef.value) return

    await passwordFormRef.value.validate(async (valid) => {
        if (valid) {
            loading.value = true
            try {
                await updateUserPassword({
                    old_password: passwordForm.oldPassword,
                    new_password: passwordForm.newPassword
                })
                ElMessage.success('密码修改成功，请重新登录')
                resetForm(passwordFormRef.value)
            } catch (error) {
                console.error(error)
            } finally {
                loading.value = false
            }
        }
    })
}

const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
}

onMounted(() => {
    fetchUserInfo()
})
</script>
<style scoped>
.user-profile {
    padding: 10px;
}

.box-card {
    text-align: center;
}

.user-header {
    padding-bottom: 20px;
    border-bottom: 1px solid #EBEEF5;
    margin-bottom: 20px;
}

.user-name {
    font-size: 24px;
    font-weight: bold;
    margin: 10px 0;
}

.user-role {
    margin-top: 5px;
}

.user-bio {
    text-align: left;
}

.bio-item {
    display: flex;
    justify-content: space-between;
    margin-bottom: 10px;
    font-size: 14px;
    color: #606266;
}

.bio-item .value {
    color: #303133;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.profile-row {
    display: flex;
    align-items: stretch;
}

.full-height-card {
    height: 100%;
}

:deep(.user-info-desc .el-descriptions__label) {
    width: 25%;
}

.password-rule-info {
    margin-top: 20px;
    padding: 15px;
    background-color: #fef0f0;
    border-radius: 4px;
    border: 1px solid #fde2e2;
    color: #f56c6c;
    font-size: 13px;
}

.password-rule-info h4 {
    margin: 0 0 8px 0;
    font-weight: bold;
}

.password-rule-info p {
    margin: 4px 0;
}
</style>
