<template>
    <el-dialog v-model="visible" :title="title" width="640px" destroy-on-close @close="handleClose">
        <el-tabs v-model="activeTab" v-loading="loading">
            <!-- Tab 1：账户信息 -->
            <el-tab-pane label="账户信息" name="account" style="height: 495px;">
                <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" style="padding:10px 0">
                    <el-form-item label="账　　号" prop="account">
                        <el-input v-model="formData.account" placeholder="请输入账号" autocomplete="off" />
                    </el-form-item>
                    <el-form-item label="密　　码" prop="password">
                        <el-input v-model="formData.password" type="password" show-password autocomplete="new-password" :placeholder="formData.id ? '留空则保持原密码' : '请输入密码'" />
                    </el-form-item>
                    <el-form-item label="姓　　名" prop="name">
                        <el-input v-model="formData.name" placeholder="请输入姓名" />
                    </el-form-item>
                    <el-form-item label="手 &nbsp;机 &nbsp;号" prop="mobile">
                        <el-input v-model="formData.mobile" placeholder="请输入手机号" />
                    </el-form-item>
                    <el-form-item label="邮　　箱" prop="email">
                        <el-input v-model="formData.email" placeholder="请输入邮箱" />
                    </el-form-item>
                    <el-form-item label="状　　态" prop="status">
                        <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
                    </el-form-item>
                    <el-form-item v-if="formData.status === 1" label="过期时间" prop="end_time">
                        <el-date-picker v-model="formData.end_time" type="datetime" placeholder="不设置则永不过期" value-format="YYYY-MM-DDTHH:mm:ss" clearable style="width:100%" />
                    </el-form-item>
                    <el-form-item label="备　　注" prop="remark">
                        <el-input v-model="formData.remark" type="textarea" :rows="2" placeholder="请输入备注" />
                    </el-form-item>
                </el-form>
            </el-tab-pane>
            <!-- Tab 2：个人档案 -->
            <el-tab-pane label="用户档案" name="profile" style="height: 495px;">
                <el-form :model="profileForm" label-width="100px" style="padding:10px 0">
                    <el-form-item label="昵&#12288;&#12288;称"><el-input v-model="profileForm.nickname" /></el-form-item>
                    <el-form-item label="性&#12288;&#12288;别">
                        <el-radio-group v-model="profileForm.gender">
                            <el-radio value="男">男</el-radio>
                            <el-radio value="女">女</el-radio>
                            <el-radio value="保密">保密</el-radio>
                        </el-radio-group>
                    </el-form-item>
                    <el-form-item label="生&#12288;&#12288;日">
                        <el-date-picker v-model="profileForm.birthday" type="date" placeholder="选择生日" value-format="YYYY-MM-DD" style="width:100%" />
                    </el-form-item>
                    <el-form-item label="所&nbsp;&nbsp;在&nbsp;&nbsp;地"><el-input v-model="profileForm.location" /></el-form-item>
                    <el-form-item label="头像URL"><el-input v-model="profileForm.avatar" placeholder="输入图片URL" /></el-form-item>
                    <el-form-item label="个人简介">
                        <el-input v-model="profileForm.bio" type="textarea" :rows="3" />
                    </el-form-item>
                </el-form>
            </el-tab-pane>
            <!-- Tab 3：地址簿 -->
            <el-tab-pane label="用户地址簿" name="address" style="height: 495px;">
                <div style="margin-bottom:10px">
                    <el-button type="primary" size="small" :icon="Plus" @click="openAddrForm()">新增地址</el-button>
                </div>
                <el-table :data="addresses" border size="small" max-height="260">
                    <el-table-column prop="label" label="标签" width="70">
                        <template #default="{ row }"><el-tag size="small">{{ addrLabelMap[row.label] || row.label }}</el-tag></template>
                    </el-table-column>
                    <el-table-column prop="receiver_name" label="收件人" width="80" />
                    <el-table-column label="地址" show-overflow-tooltip>
                        <template #default="{ row }">{{ row.province }}{{ row.city }}{{ row.district }}{{ row.detail }}</template>
                    </el-table-column>
                    <el-table-column label="默认" width="55">
                        <template #default="{ row }"><el-tag v-if="row.is_default" type="success" size="small">默认</el-tag></template>
                    </el-table-column>
                    <el-table-column label="操作" width="145" fixed="right">
                        <template #default="{ row }">
                            <el-button text bg size="small" @click="openAddrForm(row)">编辑</el-button>
                            <el-button v-if="!row.is_default" text bg size="small" @click="doSetDefault(row)">设默认</el-button>
                            <el-button text bg size="small" type="danger" @click="doDelAddr(row)">删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
                <!-- 地址子弹窗 -->
                <el-dialog v-model="addrFormVis" :title="addrForm.id ? '编辑地址' : '新增地址'" width="460px" append-to-body>
                    <el-form :model="addrForm" label-width="80px">
                        <el-form-item label="标&#12288;&#12288;签">
                            <el-select v-model="addrForm.label" style="width:100%">
                                <el-option label="家庭" value="home" />
                                <el-option label="公司" value="office" />
                                <el-option label="其他" value="other" />
                            </el-select>
                        </el-form-item>
                        <el-form-item label="收&nbsp;&nbsp;件&nbsp;&nbsp;人"><el-input v-model="addrForm.receiver_name" /></el-form-item>
                        <el-form-item label="联系电话"><el-input v-model="addrForm.phone" /></el-form-item>
                        <el-form-item label="省&#12288;&#12288;份"><el-input v-model="addrForm.province" /></el-form-item>
                        <el-form-item label="城&#12288;&#12288;市"><el-input v-model="addrForm.city" /></el-form-item>
                        <el-form-item label="区&#12288;&#12288;县"><el-input v-model="addrForm.district" /></el-form-item>
                        <el-form-item label="详细地址"><el-input v-model="addrForm.detail" type="textarea" :rows="2" /></el-form-item>
                    </el-form>
                    <template #footer>
                        <el-button @click="addrFormVis = false">取消</el-button>
                        <el-button type="primary" @click="saveAddr" :loading="extSaving">确定</el-button>
                    </template>
                </el-dialog>
            </el-tab-pane>
            <!-- Tab 4：额外联系 -->
            <el-tab-pane label="额外联系方式" name="contact" style="height: 495px;">
                <div style="margin-bottom:10px">
                    <el-button type="primary" size="small" :icon="Plus" @click="openContactForm()">新增</el-button>
                </div>
                <el-table :data="contacts" border size="small" max-height="260">
                    <el-table-column prop="type" label="类型" width="100">
                        <template #default="{ row }">
                            <el-tag size="small">{{ contactTypeMap[row.type] || row.type }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="value" label="值" />
                    <el-table-column prop="remark" label="备注" width="120" show-overflow-tooltip />
                    <el-table-column label="操作" width="100" fixed="right">
                        <template #default="{ row }">
                            <el-button text bg size="small" @click="openContactForm(row)">编辑</el-button>
                            <el-button text bg size="small" type="danger" @click="doDelContact(row)">删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
                <!-- 联系方式子弹窗 -->
                <el-dialog v-model="contactFormVis" :title="contactForm.id ? '编辑联系方式' : '新增联系方式'" width="420px" append-to-body>
                    <el-form :model="contactForm" label-width="80px">
                        <el-form-item label="类型">
                            <el-select v-model="contactForm.type" style="width:100%">
                                <el-option label="备用邮箱" value="email" />
                                <el-option label="备用手机" value="phone" />
                                <el-option label="微信" value="wechat" />
                                <el-option label="QQ" value="qq" />
                            </el-select>
                        </el-form-item>
                        <el-form-item label="值"><el-input v-model="contactForm.value" /></el-form-item>
                        <el-form-item label="备注"><el-input v-model="contactForm.remark" /></el-form-item>
                    </el-form>
                    <template #footer>
                        <el-button @click="contactFormVis = false">取消</el-button>
                        <el-button type="primary" @click="saveContact" :loading="extSaving">确定</el-button>
                    </template>
                </el-dialog>
            </el-tab-pane>
        </el-tabs>
        <template #footer>
            <el-button @click="handleClose">取消</el-button>
            <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
                {{ formData.id ? '保存' : '创建用户' }}
            </el-button>
        </template>
    </el-dialog>
</template>
<script setup lang="ts">
import { createUser, updateUser } from '@/api/user'
import {
    createUserAddress, createUserContact,
    deleteUserAddress, deleteUserContact,
    getUserAddresses, getUserContacts, getUserProfile,
    setDefaultAddress, updateUserAddress, updateUserContact, upsertUserProfile
} from '@/api/user-profile'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, ref, watch } from 'vue'

const props = defineProps<{
    modelValue: boolean
    data: any
}>()

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'success'): void
}>()

const visible = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
})

const title = computed(() => props.data?.id ? '编辑用户' : '新增用户')

// ────── 账户 Tab ──────
const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const formData = ref<any>({})
const formRules = computed<FormRules>(() => ({
    account: [{ required: true, message: '请输入账号', trigger: 'blur' }],
    password: formData.value.id ? [] : [{ required: true, message: '请输入密码', trigger: 'blur' }]
}))

// ────── 扩展 Tab ──────
const activeTab = ref('account')
const loading = ref(false)
const extSaving = ref(false)

// 档案
const profileForm = ref<any>({ nickname: '', avatar: '', gender: '保密', birthday: null, bio: '', location: '' })

// 地址
const addresses = ref<any[]>([])
const addrFormVis = ref(false)
const addrForm = ref<any>({})
const addrLabelMap: Record<string, string> = { home: '家庭', office: '公司', other: '其他' }

// 联系
const contacts = ref<any[]>([])
const contactFormVis = ref(false)
const contactForm = ref<any>({})
const contactTypeMap: Record<string, string> = { email: '备用邮箱', phone: '备用手机', wechat: '微信', qq: 'QQ' }

// ────── 初始化：弹窗打开时加载数据 ──────
watch(() => props.modelValue, async (val) => {
    if (!val) return
    activeTab.value = 'account'

    if (props.data?.id) {
        // 编辑模式：回填账户数据，异步加载扩展信息
        formData.value = { ...props.data }
        loading.value = true
        try {
            const [profile, addrs, ctcs] = await Promise.all([
                getUserProfile(props.data.id),
                getUserAddresses(props.data.id),
                getUserContacts(props.data.id),
            ])
            profileForm.value = { nickname: '', avatar: '', gender: '保密', birthday: null, bio: '', location: '', ...(profile || {}) }
            addresses.value = Array.isArray(addrs) ? addrs : []
            contacts.value = Array.isArray(ctcs) ? ctcs : []
        } catch { /* 加载失败不阻断弹窗 */ }
        finally { loading.value = false }
    } else {
        // 新增模式：重置所有数据
        formData.value = { status: 1 }
        profileForm.value = { nickname: '', avatar: '', gender: '保密', birthday: null, bio: '', location: '' }
        addresses.value = []
        contacts.value = []
    }
})

const handleClose = () => {
    formRef.value?.resetFields()
    emit('update:modelValue', false)
}

// ────── 统一保存：账户 + 档案 + 待创建地址/联系方式 ──────
const handleSubmit = async () => {
    if (!formRef.value) return
    await formRef.value.validate(async (valid) => {
        if (!valid) return
        submitLoading.value = true
        try {
            let userId = formData.value.id
            // 1. 保存账户信息
            if (userId) {
                await updateUser(userId, formData.value)
            } else {
                const res: any = await createUser(formData.value)
                // 后端直接返回 SsoUser 对象，id 字段即为用户 ID
                userId = res?.id || ''
                if (userId) formData.value = { ...formData.value, id: userId }
            }
            if (!userId) {
                ElMessage.error('获取用户 ID 失败，扩展信息保存跳过')
                return
            }
            // 2. 保存档案（如果有内容）
            const hasProfile = Object.values(profileForm.value).some(v => v !== '' && v !== null && v !== '保密')
            if (hasProfile) {
                await upsertUserProfile(userId, profileForm.value).catch(() => { })
            }
            // 3. 批量创建未保存的地址（新增模式下没有 id 的条目）
            const newAddrs = addresses.value.filter(a => !a.id)
            for (const addr of newAddrs) {
                await createUserAddress(userId, addr).catch(() => { })
            }
            // 4. 批量创建未保存的联系方式
            const newContacts = contacts.value.filter(c => !c.id)
            for (const c of newContacts) {
                await createUserContact(userId, c).catch(() => { })
            }
            ElMessage.success(formData.value.id ? '保存成功' : '创建成功')
            handleClose()
            emit('success')
        } catch (error: any) {
            ElMessage.error(error.message || '操作失败')
        } finally {
            submitLoading.value = false
        }
    })
}

// ────── 地址操作 ──────
const openAddrForm = (row?: any) => {
    addrForm.value = row
        ? { ...row }
        : {
            label: 'home',
            receiver_name: '',
            phone: '',
            province: '',
            city: '',
            district: '',
            detail: ''
        }
    addrFormVis.value = true
}
const saveAddr = async () => {
    if (!formData.value.id) {
        // 新增模式：将地址暂存到列表，等创建用户时一并提交
        if (addrForm.value._idx !== undefined) {
            addresses.value[addrForm.value._idx] = { ...addrForm.value }
        } else {
            addresses.value.push({ ...addrForm.value })
        }
        addrFormVis.value = false
        return
    }
    extSaving.value = true
    try {
        addrForm.value.id
            ? await updateUserAddress(formData.value.id, addrForm.value.id, addrForm.value)
            : await createUserAddress(formData.value.id, addrForm.value)
        ElMessage.success('保存成功')
        addrFormVis.value = false
        addresses.value = await getUserAddresses(formData.value.id)
    } catch { ElMessage.error('保存失败') }
    finally { extSaving.value = false }
}
const doSetDefault = async (row: any) => {
    if (!formData.value.id) {
        // 新增模式：本地设置默认
        addresses.value = addresses.value.map(a => ({ ...a, is_default: a === row }))
        return
    }
    try {
        await setDefaultAddress(formData.value.id, row.id)
        addresses.value = await getUserAddresses(formData.value.id)
    } catch { ElMessage.error('操作失败') }
}
const doDelAddr = async (row: any) => {
    try {
        if (!row.id) {
            // 新增模式：本地移除
            const idx = addresses.value.indexOf(row)
            if (idx > -1) addresses.value.splice(idx, 1)
            return
        }
        await deleteUserAddress(formData.value.id, row.id)
        ElMessage.success('删除成功')
        addresses.value = await getUserAddresses(formData.value.id)
    } catch (e: any) { if (e !== 'cancel') ElMessage.error('删除失败') }
}

// ────── 联系方式操作 ──────
const openContactForm = (row?: any) => {
    contactForm.value = row ? { ...row } : { type: 'email', value: '', remark: '' }
    contactFormVis.value = true
}
const saveContact = async () => {
    if (!formData.value.id) {
        // 新增模式：将联系方式暂存到列表
        if (contactForm.value._idx !== undefined) {
            contacts.value[contactForm.value._idx] = { ...contactForm.value }
        } else {
            contacts.value.push({ ...contactForm.value })
        }
        contactFormVis.value = false
        return
    }
    extSaving.value = true
    try {
        contactForm.value.id
            ? await updateUserContact(formData.value.id, contactForm.value.id, contactForm.value)
            : await createUserContact(formData.value.id, contactForm.value)
        ElMessage.success('保存成功')
        contactFormVis.value = false
        contacts.value = await getUserContacts(formData.value.id)
    } catch { ElMessage.error('保存失败') }
    finally { extSaving.value = false }
}
const doDelContact = async (row: any) => {
    try {
        await ElMessageBox.confirm('确定删除该联系方式吗？', '提示', { type: 'warning' })
        if (!row.id) {
            const idx = contacts.value.indexOf(row)
            if (idx > -1) contacts.value.splice(idx, 1)
            return
        }
        await deleteUserContact(formData.value.id, row.id)
        ElMessage.success('删除成功')
        contacts.value = await getUserContacts(formData.value.id)
    } catch (e: any) { if (e !== 'cancel') ElMessage.error('删除失败') }
}
</script>
