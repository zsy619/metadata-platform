<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary page-title">
                <el-icon class="title-icon">
                    <Avatar />
                </el-icon>
                用户管理
            </h1>
            <div class="header-actions">
                <el-button type="primary" @click="handleCreate" :icon="Plus">新增用户</el-button>
            </div>
        </div>
        <el-card class="main-card">
            <div class="search-area">
                <el-input v-model="searchQuery" placeholder="请输入用户名搜索" clearable :prefix-icon="Search" style="width: 300px" @input="handleDebouncedSearch" />
                <el-select v-model="filterStatus" placeholder="筛选状态" style="width: 150px; margin-left: 10px" clearable @change="handleSearch">
                    <el-option label="全部" value="" />
                    <el-option label="有效" :value="1" />
                    <el-option label="禁用" :value="0" />
                </el-select>
                <el-button type="primary" @click="handleSearch" :icon="Search" style="margin-left: 10px">搜索</el-button>
                <el-button @click="handleReset" :icon="RefreshLeft">重置</el-button>
            </div>
            <div class="table-area">
                <el-table v-loading="loading" :element-loading-text="loadingText" :data="filteredData" border stripe style="width: 100%; height: 100%;">
                    <template #empty>
                        <el-empty :description="searchQuery ? '未搜索到相关用户' : '暂无用户'">
                            <el-button v-if="!searchQuery" type="primary" @click="handleCreate">新增用户</el-button>
                        </el-empty>
                    </template>
                    <el-table-column prop="name" label="姓名" width="120" />
                    <el-table-column prop="account" label="账号" width="150" />
                    <el-table-column prop="mobile" label="手机号" width="130" />
                    <el-table-column prop="email" label="邮箱" width="180" show-overflow-tooltip />
                    <el-table-column prop="status" label="状态" width="80">
                        <template #default="scope">
                            <el-tag v-if="scope.row.status === 1" type="success">有效</el-tag>
                            <el-tag v-else type="danger">禁用</el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="create_at" label="创建时间" width="200">
                        <template #default="scope">{{ formatDateTime(scope.row.create_at) }}</template>
                    </el-table-column>
                    <el-table-column label="操作" width="280" fixed="right">
                        <template #default="scope">
                            <el-button type="primary" size="small" :icon="Edit" @click="handleEdit(scope.row)" text bg>编辑</el-button>
                            <el-button v-if="!scope.row.is_system" type="danger" size="small" :icon="Delete" @click="handleDelete(scope.row)" text bg>删除</el-button>
                            <el-tag v-else type="warning" size="small" effect="light" style="margin-left: 4px">系统内置</el-tag>
                            <el-dropdown trigger="click" @command="(cmd: string) => handleSetting(cmd, scope.row)">
                                <el-button type="info" size="small" text bg>
                                    设置<el-icon class="el-icon--right">
                                        <ArrowDown />
                                    </el-icon>
                                </el-button>
                                <template #dropdown>
                                    <el-dropdown-menu>
                                        <el-dropdown-item command="profile">个人档案</el-dropdown-item>
                                        <el-dropdown-item command="roles">设置角色</el-dropdown-item>
                                        <el-dropdown-item command="pos">设置职位</el-dropdown-item>
                                        <el-dropdown-item command="groups">设置用户组</el-dropdown-item>
                                        <el-dropdown-item command="role-groups">设置角色组</el-dropdown-item>
                                        <el-dropdown-item command="orgs">设置组织</el-dropdown-item>
                                    </el-dropdown-menu>
                                </template>
                            </el-dropdown>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
        </el-card>
        <!-- 新增/编辑弹窗 -->
        <UserForm v-model="dialogVisible" :data="formData" @success="loadData" />
        <UserSettingDialog v-model="settingDialogVisible" :user-id="currentUserId" :user-name="currentUserName" :setting-type="currentSettingType" @success="loadData" />
        <UserProfileDialog v-model="profileDialogVisible" :user-id="currentUserId" :user-name="currentUserName" />
    </div>
</template>
<script setup lang="ts">
import { deleteUser, getUsers } from '@/api/user'
import { ArrowDown, Avatar, Delete, Edit, Plus, RefreshLeft, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import UserForm from './UserForm.vue'
import UserProfileDialog from './UserProfileDialog.vue'
import UserSettingDialog from './UserSettingDialog.vue'

const loading = ref(false)
const loadingText = ref('加载中...')
const searchQuery = ref('')
const filterStatus = ref<number | ''>('')

const allData = ref<any[]>([])

const filteredData = computed(() => {
    let data = allData.value
    if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        data = data.filter(item => (item.name || '').toLowerCase().includes(query) || (item.account || '').toLowerCase().includes(query))
    }
    if (filterStatus.value !== '') data = data.filter(item => item.status === filterStatus.value)
    return data
})

const dialogVisible = ref(false)
const formData = ref<any>({})

const formatDateTime = (dateStr: string) => {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    return isNaN(date.getTime()) ? '-' : date.toLocaleString('zh-CN')
}

const loadData = async () => {
    loadingText.value = '加载中...'
    loading.value = true
    try {
        const res: any = await getUsers()
        allData.value = res.data || res
    } catch (error) {
        console.error('加载用户列表失败:', error)
        ElMessage.error('加载列表失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => { }
const handleDebouncedSearch = () => { }
const handleReset = () => { searchQuery.value = ''; filterStatus.value = '' }

const handleCreate = () => {
    formData.value = { status: 1 }
    dialogVisible.value = true
}

const handleEdit = (row: any) => {
    formData.value = { ...row }
    dialogVisible.value = true
}

const handleDelete = async (row: any) => {
    try {
        await ElMessageBox.confirm(`确定要删除用户 "${row.name}" 吗？`, '提示', { type: 'warning' })
        await deleteUser(row.id)
        ElMessage.success('删除成功')
        loadData()
    } catch (error: any) { if (error !== 'cancel') ElMessage.error(error.message || '删除失败') }
}

// 设置弹窗相关
const settingDialogVisible = ref(false)
const profileDialogVisible = ref(false)
const currentUserId = ref('')
const currentUserName = ref('')
const currentSettingType = ref<'roles' | 'pos' | 'groups' | 'role-groups' | 'orgs'>('roles')

const handleSetting = (command: string, row: any) => {
    currentUserId.value = row.id
    currentUserName.value = row.name || row.account
    if (command === 'profile') {
        profileDialogVisible.value = true
    } else {
        currentSettingType.value = command as 'roles' | 'pos' | 'groups' | 'role-groups' | 'orgs'
        settingDialogVisible.value = true
    }
}

onMounted(() => loadData())
</script>
<style scoped>
.sso-page {
    height: 100%;
    display: flex;
    flex-direction: column;
}

.main-card {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

:deep(.el-card__body) {
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 20px;
    overflow: hidden;
    box-sizing: border-box;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-shrink: 0;
}

.page-title {
    font-size: 20px;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 8px;
}

.title-icon {
    font-size: 24px;
    color: var(--el-color-primary);
}

.header-actions {
    display: flex;
    gap: 10px;
}

.search-area {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
    flex-wrap: wrap;
    gap: 10px;
}

.table-area {
    flex: 1;
    overflow: hidden;
}

.text-primary {
    color: var(--el-text-color-primary);
}
</style>
