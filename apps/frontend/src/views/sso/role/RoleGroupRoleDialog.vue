<template>
    <el-dialog v-model="visible" title="关联角色" width="500px" :close-on-click-modal="false" @close="handleClose" @open="handleOpen">
        <div v-loading="loading" class="role-tree-container">
            <el-tree
                ref="treeRef"
                :data="roleTreeData"
                :props="{ children: 'children', label: 'label' }"
                show-checkbox
                node-key="id"
                default-expand-all
                :check-strictly="true"
                highlight-current
            >
                <template #default="{ node, data }">
                    <span class="tree-node">
                        <el-icon v-if="data.children && data.children.length > 0" class="folder-icon"><Folder /></el-icon>
                        <el-icon v-else class="leaf-icon"><User /></el-icon>
                        <span class="node-label">{{ node.label }}</span>
                    </span>
                </template>
            </el-tree>
        </div>
        <template #footer>
            <el-button @click="handleClose">取消</el-button>
            <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">
import { getRoleGroupRoles, getRoles, updateRoleGroupRoles } from '@/api/user'
import { Folder, User } from '@element-plus/icons-vue'
import type { ElTree } from 'element-plus'
import { ElMessage } from 'element-plus'
import { computed, nextTick, ref } from 'vue'

const props = defineProps<{
    modelValue: boolean
    groupId: string
    groupName: string
}>()

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'success'): void
}>()

const visible = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
})

const loading = ref(false)
const submitting = ref(false)
const allRoles = ref<any[]>([])
const treeRef = ref<InstanceType<typeof ElTree>>()

const isRoot = (pid: string) => !pid || pid === '' || pid === '0'

const roleTreeData = computed(() => {
    const buildTree = (items: any[], parentId = ''): any[] => {
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
                id: String(item.id),
                label: item.role_name,
                children: buildTree(items, item.id)
            }))
    }
    return buildTree(allRoles.value)
})

const handleOpen = async () => {
    if (!props.groupId) return

    loading.value = true
    try {
        const [rolesRes, groupRoles] = await Promise.all([
            getRoles(),
            getRoleGroupRoles(props.groupId)
        ])

        allRoles.value = Array.isArray(rolesRes) ? rolesRes : (rolesRes as any).data || []

        await nextTick()
        if (treeRef.value && groupRoles && groupRoles.length > 0) {
            treeRef.value.setCheckedKeys(groupRoles, false)
        }
    } catch (error: any) {
        console.error('加载数据失败:', error)
        ElMessage.error(error.message || '加载数据失败')
    } finally {
        loading.value = false
    }
}

const handleClose = () => {
    allRoles.value = []
    emit('update:modelValue', false)
}

const handleSubmit = async () => {
    if (!treeRef.value) {
        ElMessage.error('树组件未初始化')
        return
    }

    if (!props.groupId) {
        ElMessage.error('角色组ID不能为空')
        return
    }

    const checkedKeys = treeRef.value.getCheckedKeys(false)
    const roleIds = (checkedKeys || []).map(key => String(key))

    submitting.value = true
    try {
        await updateRoleGroupRoles(props.groupId, roleIds)
        ElMessage.success('更新成功')
        handleClose()
        emit('success')
    } catch (error: any) {
        console.error('更新失败:', error)
        ElMessage.error(error.message || '更新失败')
    } finally {
        submitting.value = false
    }
}
</script>

<style scoped>
.role-tree-container {
    height: 400px;
    overflow-y: auto;
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 4px;
    padding: 10px;
}

.tree-node {
    display: flex;
    align-items: center;
    gap: 6px;
}

.folder-icon {
    color: var(--el-color-warning);
    font-size: 16px;
}

.leaf-icon {
    color: var(--el-text-color-secondary);
    font-size: 16px;
}

.node-label {
    font-size: 14px;
}
</style>
