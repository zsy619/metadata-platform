<template>
    <el-dialog v-model="visible" :title="`${userName} - 个人档案`" width="720px" destroy-on-close>
        <el-tabs v-model="activeTab" v-loading="loading">
            <!-- ── 基本档案 ── -->
            <el-tab-pane label="基本档案" name="profile">
                <el-form :model="profileForm" label-width="80px" style="padding: 10px 0">
                    <el-form-item label="昵&#12288;&#12288;称">
                        <el-input v-model="profileForm.nickname" placeholder="请输入昵称" />
                    </el-form-item>
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
                    <el-form-item label="所&nbsp;&nbsp;在&nbsp;&nbsp;地">
                        <el-input v-model="profileForm.location" placeholder="请输入所在地" />
                    </el-form-item>
                    <el-form-item label="头像URL">
                        <el-input v-model="profileForm.avatar" placeholder="请输入头像图片地址" />
                    </el-form-item>
                    <el-form-item label="个人简介">
                        <el-input v-model="profileForm.bio" type="textarea" :rows="3" placeholder="请输入个人简介" />
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="saveProfile" :loading="saving">保存档案</el-button>
                    </el-form-item>
                </el-form>
            </el-tab-pane>
            <!-- ── 地址簿 ── -->
            <el-tab-pane label="地址簿" name="address">
                <div style="margin-bottom:12px">
                    <el-button type="primary" size="small" :icon="Plus" @click="openAddressForm()">新增地址</el-button>
                </div>
                <el-table :data="addresses" border size="small">
                    <el-table-column prop="label" label="标签" width="80">
                        <template #default="{ row }">
                            <el-tag size="small">{{ labelMap[row.label] || row.label }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="receiver_name" label="收件人" width="90" />
                    <el-table-column prop="phone" label="电话" width="120" />
                    <el-table-column label="地址" show-overflow-tooltip>
                        <template #default="{ row }">{{ row.province }}{{ row.city }}{{ row.district }}{{ row.detail }}</template>
                    </el-table-column>
                    <el-table-column label="默认" width="60">
                        <template #default="{ row }">
                            <el-tag v-if="row.is_default" type="success" size="small">默认</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="160" fixed="right">
                        <template #default="{ row }">
                            <el-button text bg size="small" @click="openAddressForm(row)">编辑</el-button>
                            <el-button v-if="!row.is_default" text bg size="small" @click="handleSetDefault(row)">设默认</el-button>
                            <el-button text bg size="small" type="danger" @click="handleDeleteAddress(row)">删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
                <!-- 地址编辑表单 -->
                <el-dialog v-model="addrFormVisible" :title="addrForm.id ? '编辑地址' : '新增地址'" width="480px" append-to-body>
                    <el-form :model="addrForm" label-width="80px">
                        <el-form-item label="标签">
                            <el-select v-model="addrForm.label" style="width:100%">
                                <el-option label="家庭" value="home" />
                                <el-option label="公司" value="office" />
                                <el-option label="其他" value="other" />
                            </el-select>
                        </el-form-item>
                        <el-form-item label="收件人"><el-input v-model="addrForm.receiver_name" /></el-form-item>
                        <el-form-item label="联系电话"><el-input v-model="addrForm.phone" /></el-form-item>
                        <el-form-item label="省份"><el-input v-model="addrForm.province" /></el-form-item>
                        <el-form-item label="城市"><el-input v-model="addrForm.city" /></el-form-item>
                        <el-form-item label="区县"><el-input v-model="addrForm.district" /></el-form-item>
                        <el-form-item label="详细地址"><el-input v-model="addrForm.detail" type="textarea" :rows="2" /></el-form-item>
                        <el-form-item label="备注"><el-input v-model="addrForm.remark" /></el-form-item>
                    </el-form>
                    <template #footer>
                        <el-button @click="addrFormVisible = false">取消</el-button>
                        <el-button type="primary" @click="saveAddress" :loading="saving">确定</el-button>
                    </template>
                </el-dialog>
            </el-tab-pane>
            <!-- ── 联系方式 ── -->
            <el-tab-pane label="联系方式" name="contact">
                <div style="margin-bottom:12px">
                    <el-button type="primary" size="small" :icon="Plus" @click="openContactForm()">新增联系方式</el-button>
                </div>
                <el-table :data="contacts" border size="small">
                    <el-table-column prop="type" label="类型" width="100">
                        <template #default="{ row }">
                            <el-tag size="small" :type="contactTypeColor[row.type]">{{ contactTypeMap[row.type] || row.type }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="value" label="值" />
                    <el-table-column prop="is_verified" label="已验证" width="80">
                        <template #default="{ row }">
                            <el-tag :type="row.is_verified ? 'success' : 'info'" size="small">{{ row.is_verified ? '是' : '否' }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="remark" label="备注" width="120" show-overflow-tooltip />
                    <el-table-column label="操作" width="120" fixed="right">
                        <template #default="{ row }">
                            <el-button text bg size="small" @click="openContactForm(row)">编辑</el-button>
                            <el-button text bg size="small" type="danger" @click="handleDeleteContact(row)">删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
                <el-dialog v-model="contactFormVisible" :title="contactForm.id ? '编辑联系方式' : '新增联系方式'" width="420px" append-to-body>
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
                        <el-button @click="contactFormVisible = false">取消</el-button>
                        <el-button type="primary" @click="saveContact" :loading="saving">确定</el-button>
                    </template>
                </el-dialog>
            </el-tab-pane>
            <!-- ── 第三方账号 ── -->
            <el-tab-pane label="第三方账号" name="social">
                <div style="margin-bottom:12px">
                    <el-button type="primary" size="small" :icon="Plus" @click="openSocialForm()">绑定账号</el-button>
                </div>
                <el-table :data="socials" border size="small">
                    <el-table-column prop="provider" label="平台" width="100">
                        <template #default="{ row }">
                            <el-tag size="small" :type="providerColor[row.provider]">{{ providerMap[row.provider] || row.provider }}</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="nickname" label="平台昵称" width="120" />
                    <el-table-column prop="open_id" label="Open ID" show-overflow-tooltip />
                    <el-table-column label="操作" width="80" fixed="right">
                        <template #default="{ row }">
                            <el-button text bg size="small" type="danger" @click="handleUnbind(row)">解绑</el-button>
                        </template>
                    </el-table-column>
                </el-table>
                <el-dialog v-model="socialFormVisible" title="绑定第三方账号" width="420px" append-to-body>
                    <el-form :model="socialForm" label-width="80px">
                        <el-form-item label="平台">
                            <el-select v-model="socialForm.provider" style="width:100%">
                                <el-option label="微信" value="wechat" />
                                <el-option label="钉钉" value="dingtalk" />
                                <el-option label="GitHub" value="github" />
                                <el-option label="QQ" value="qq" />
                            </el-select>
                        </el-form-item>
                        <el-form-item label="Open ID"><el-input v-model="socialForm.open_id" /></el-form-item>
                        <el-form-item label="Union ID"><el-input v-model="socialForm.union_id" /></el-form-item>
                        <el-form-item label="平台昵称"><el-input v-model="socialForm.nickname" /></el-form-item>
                    </el-form>
                    <template #footer>
                        <el-button @click="socialFormVisible = false">取消</el-button>
                        <el-button type="primary" @click="saveSocial" :loading="saving">绑定</el-button>
                    </template>
                </el-dialog>
            </el-tab-pane>
        </el-tabs>
    </el-dialog>
</template>
<script setup lang="ts">
import {
    bindUserSocial, createUserAddress, createUserContact,
    deleteUserAddress, deleteUserContact,
    getUserAddresses, getUserContacts, getUserProfile, getUserSocials,
    setDefaultAddress, unbindUserSocial, updateUserAddress, updateUserContact, upsertUserProfile
} from '@/api/user-profile';
import { Plus } from '@element-plus/icons-vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { ref, watch } from 'vue';

const props = defineProps<{ modelValue: boolean; userId: string; userName: string }>()
const emit = defineEmits(['update:modelValue'])

const visible = ref(false)
const activeTab = ref('profile')
const loading = ref(false)
const saving = ref(false)

// ────── 基本档案 ──────
const profileForm = ref<any>({ nickname: '', avatar: '', gender: 'U', birthday: null, bio: '', location: '' })

// ────── 地址簿 ──────
const addresses = ref<any[]>([])
const addrFormVisible = ref(false)
const addrForm = ref<any>({})
const labelMap: Record<string, string> = { home: '家庭', office: '公司', other: '其他' }

// ────── 联系方式 ──────
const contacts = ref<any[]>([])
const contactFormVisible = ref(false)
const contactForm = ref<any>({})
const contactTypeMap: Record<string, string> = { email: '备用邮箱', phone: '备用手机', wechat: '微信', qq: 'QQ' }
const contactTypeColor: Record<string, any> = { email: 'primary', phone: 'success', wechat: 'success', qq: 'warning' }

// ────── 第三方账号 ──────
const socials = ref<any[]>([])
const socialFormVisible = ref(false)
const socialForm = ref<any>({})
const providerMap: Record<string, string> = { wechat: '微信', dingtalk: '钉钉', github: 'GitHub', qq: 'QQ' }
const providerColor: Record<string, any> = { wechat: 'success', dingtalk: 'primary', github: '', qq: 'warning' }

// ────── 数据加载 ──────
const loadAll = async () => {
    if (!props.userId) return
    loading.value = true
    try {
        const [profile, addrs, ctcs, socs] = await Promise.all([
            getUserProfile(props.userId),
            getUserAddresses(props.userId),
            getUserContacts(props.userId),
            getUserSocials(props.userId),
        ])
        profileForm.value = { nickname: '', avatar: '', gender: 'U', birthday: null, bio: '', location: '', ...(profile || {}) }
        addresses.value = Array.isArray(addrs) ? addrs : []
        contacts.value = Array.isArray(ctcs) ? ctcs : []
        socials.value = Array.isArray(socs) ? socs : []
    } catch {
        ElMessage.error('加载数据失败')
    } finally {
        loading.value = false
    }
}

watch(() => props.modelValue, (val) => {
    visible.value = val
    if (val) { activeTab.value = 'profile'; loadAll() }
})
watch(visible, (val) => emit('update:modelValue', val))

// ────── 档案保存 ──────
const saveProfile = async () => {
    saving.value = true
    try {
        await upsertUserProfile(props.userId, profileForm.value)
        ElMessage.success('档案保存成功')
    } catch { ElMessage.error('保存失败') }
    finally { saving.value = false }
}

// ────── 地址操作 ──────
const openAddressForm = (row?: any) => {
    addrForm.value = row ? { ...row } : { label: 'home', receiver_name: '', phone: '', province: '', city: '', district: '', detail: '', remark: '' }
    addrFormVisible.value = true
}
const saveAddress = async () => {
    saving.value = true
    try {
        if (addrForm.value.id) {
            await updateUserAddress(props.userId, addrForm.value.id, addrForm.value)
        } else {
            await createUserAddress(props.userId, addrForm.value)
        }
        ElMessage.success('保存成功')
        addrFormVisible.value = false
        addresses.value = await getUserAddresses(props.userId)
    } catch { ElMessage.error('保存失败') }
    finally { saving.value = false }
}
const handleSetDefault = async (row: any) => {
    try {
        await setDefaultAddress(props.userId, row.id)
        ElMessage.success('已设为默认地址')
        addresses.value = await getUserAddresses(props.userId)
    } catch { ElMessage.error('操作失败') }
}
const handleDeleteAddress = async (row: any) => {
    try {
        await ElMessageBox.confirm(`确定删除该地址吗？`, '提示', { type: 'warning' })
        await deleteUserAddress(props.userId, row.id)
        ElMessage.success('删除成功')
        addresses.value = await getUserAddresses(props.userId)
    } catch (e: any) { if (e !== 'cancel') ElMessage.error('删除失败') }
}

// ────── 联系方式操作 ──────
const openContactForm = (row?: any) => {
    contactForm.value = row ? { ...row } : { type: 'email', value: '', remark: '' }
    contactFormVisible.value = true
}
const saveContact = async () => {
    saving.value = true
    try {
        if (contactForm.value.id) {
            await updateUserContact(props.userId, contactForm.value.id, contactForm.value)
        } else {
            await createUserContact(props.userId, contactForm.value)
        }
        ElMessage.success('保存成功')
        contactFormVisible.value = false
        contacts.value = await getUserContacts(props.userId)
    } catch { ElMessage.error('保存失败') }
    finally { saving.value = false }
}
const handleDeleteContact = async (row: any) => {
    try {
        await ElMessageBox.confirm('确定删除该联系方式吗？', '提示', { type: 'warning' })
        await deleteUserContact(props.userId, row.id)
        ElMessage.success('删除成功')
        contacts.value = await getUserContacts(props.userId)
    } catch (e: any) { if (e !== 'cancel') ElMessage.error('删除失败') }
}

// ────── 第三方账号操作 ──────
const openSocialForm = () => {
    socialForm.value = { provider: 'wechat', open_id: '', union_id: '', nickname: '' }
    socialFormVisible.value = true
}
const saveSocial = async () => {
    saving.value = true
    try {
        await bindUserSocial(props.userId, socialForm.value)
        ElMessage.success('绑定成功')
        socialFormVisible.value = false
        socials.value = await getUserSocials(props.userId)
    } catch { ElMessage.error('绑定失败') }
    finally { saving.value = false }
}
const handleUnbind = async (row: any) => {
    try {
        await ElMessageBox.confirm(`确定解绑 ${providerMap[row.provider] || row.provider} 账号吗？`, '提示', { type: 'warning' })
        await unbindUserSocial(props.userId, row.id)
        ElMessage.success('解绑成功')
        socials.value = await getUserSocials(props.userId)
    } catch (e: any) { if (e !== 'cancel') ElMessage.error('解绑失败') }
}
</script>
