<template>
    <el-dialog v-model="visible" :title="dialogTitle" width="500px" destroy-on-close @close="handleClose" @open="handleOpen">
        <div v-loading="loading">
            <el-tree
                ref="treeRef"
                :data="treeData"
                :props="treeProps"
                show-checkbox
                node-key="id"
                default-expand-all
                :check-strictly="true"
                highlight-current
                class="tree-container"
            >
                <template #default="{ data }">
                    <span class="tree-node">
                        <span class="node-name">{{ data.name || data.role_name || data.pos_name || data.group_name || data.org_name }}</span>
                    </span>
                </template>
            </el-tree>
        </div>
        <template #footer>
            <el-button @click="handleClose">取消</el-button>
            <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">
import { getUserGroups as fetchUserGroups, getAllUserGroups, getPos, getRoleGroups, getRoles, getUnits, getUserOrgs, getUserPos, getUserRoleGroups, getUserRoles, updateUserGroups, updateUserOrgs, updateUserPos, updateUserRoleGroups, updateUserRoles } from '@/api/user'
import type { ElTree } from 'element-plus'
import { ElMessage } from 'element-plus'
import { computed, ref } from 'vue'

type SettingType = 'roles' | 'pos' | 'groups' | 'role-groups' | 'orgs'

const props = defineProps<{
    modelValue: boolean
    userId: string
    userName: string
    settingType: SettingType
}>()

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'success'): void
}>()

const visible = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
})

const dialogTitle = computed(() => {
    const titles: Record<SettingType, string> = {
        'roles': '设置角色',
        'pos': '设置职位',
        'groups': '设置用户组',
        'role-groups': '设置角色组',
        'orgs': '设置组织'
    }
    return `${titles[props.settingType]} - ${props.userName}`
})

const treeProps = computed(() => ({
    children: 'children',
    label: props.settingType === 'roles' ? 'role_name' : 
           props.settingType === 'pos' ? 'pos_name' :
           props.settingType === 'groups' ? 'group_name' :
           props.settingType === 'role-groups' ? 'group_name' : 'org_name'
}))

const treeRef = ref<InstanceType<typeof ElTree>>()
const loading = ref(false)
const submitLoading = ref(false)
const allData = ref<any[]>([])

const isRoot = (pid: string) => !pid || pid === '' || pid === '0'

const buildTree = (items: any[], parentId: string = ''): any[] => {
    return items
        .filter(item => {
            const itemPid = item.parent_id || ''
            if (isRoot(parentId)) {
                return isRoot(itemPid)
            }
            return itemPid === parentId
        })
        .sort((a, b) => (a.sort || 0) - (b.sort || 0))
        .map(item => ({
            ...item,
            children: buildTree(items, item.id)
        }))
}

const treeData = computed(() => buildTree(allData.value))

const loadData = async () => {
    loading.value = true
    try {
        let data: any[] = []
        switch (props.settingType) {
            case 'roles':
                data = await getRoles()
                break
            case 'pos':
                data = await getPos()
                break
            case 'groups':
                data = await getAllUserGroups()
                break
            case 'role-groups':
                data = await getRoleGroups()
                break
            case 'orgs':
                data = await getUnits()
                break
        }
        allData.value = data || []
    } catch (error: any) {
        ElMessage.error(error.message || '加载数据失败')
    } finally {
        loading.value = false
    }
}

const loadUserSelections = async () => {
    if (!props.userId) return
    try {
        let ids: string[] = []
        switch (props.settingType) {
            case 'roles':
                const rolesRes = await getUserRoles(props.userId)
                ids = rolesRes.role_ids || []
                break
            case 'pos':
                const posRes = await getUserPos(props.userId)
                ids = posRes.pos_ids || []
                break
            case 'groups':
                const groupsRes = await fetchUserGroups(props.userId)
                ids = groupsRes.group_ids || []
                break
            case 'role-groups':
                const roleGroupsRes = await getUserRoleGroups(props.userId)
                ids = roleGroupsRes.role_group_ids || []
                break
            case 'orgs':
                const orgsRes = await getUserOrgs(props.userId)
                ids = orgsRes.org_ids || []
                break
        }
        if (treeRef.value) {
            treeRef.value.setCheckedKeys(ids, false)
        }
    } catch (error: any) {
        ElMessage.error(error.message || '加载用户关联失败')
    }
}

const handleOpen = async () => {
    await loadData()
    await loadUserSelections()
}

const handleClose = () => {
    allData.value = []
    emit('update:modelValue', false)
}

const handleSubmit = async () => {
    if (!props.userId) {
        ElMessage.error('用户ID不能为空')
        return
    }

    submitLoading.value = true
    try {
        const checkedKeys = treeRef.value?.getCheckedKeys(false) || [] as string[]

        switch (props.settingType) {
            case 'roles':
                await updateUserRoles(props.userId, checkedKeys as string[])
                break
            case 'pos':
                await updateUserPos(props.userId, checkedKeys as string[])
                break
            case 'groups':
                await updateUserGroups(props.userId, checkedKeys as string[])
                break
            case 'role-groups':
                await updateUserRoleGroups(props.userId, checkedKeys as string[])
                break
            case 'orgs':
                await updateUserOrgs(props.userId, checkedKeys as string[])
                break
        }
        ElMessage.success('设置成功')
        handleClose()
        emit('success')
    } catch (error: any) {
        ElMessage.error(error.message || '设置失败')
    } finally {
        submitLoading.value = false
    }
}
</script>

<style scoped>
.tree-container {
    height: 450px;
    overflow-y: auto;
}

.tree-node {
    display: flex;
    align-items: center;
}

.node-name {
    font-size: 14px;
}
</style>
