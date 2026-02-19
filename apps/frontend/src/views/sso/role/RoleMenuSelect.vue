<template>
    <el-dialog v-model="visible" title="配置角色菜单" width="600px" destroy-on-close @close="handleClose" @open="handleOpen">
        <div v-loading="loading">
            <el-tree
                ref="treeRef"
                :data="menuTreeData"
                :props="{ children: 'children', label: 'menu_name' }"
                show-checkbox
                node-key="id"
                default-expand-all
                :check-strictly="false"
                highlight-current
                style="max-height: 500px; overflow-y: auto"
            >
                <template #default="{ data }">
                    <span class="menu-node">
                        <span class="menu-name">{{ data.menu_name }}</span>
                        <el-tag v-if="data.menu_type === 'M'" type="warning" size="small" style="margin-left: 8px">目录</el-tag>
                        <el-tag v-else-if="data.menu_type === 'C'" type="success" size="small" style="margin-left: 8px">菜单</el-tag>
                        <el-tag v-else-if="data.menu_type === 'F'" type="info" size="small" style="margin-left: 8px">按钮</el-tag>
                        <el-tag v-else type="primary" size="small" style="margin-left: 8px">资源</el-tag>
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
import { getMenus, getRoleMenus, updateRoleMenus } from '@/api/user'
import type { FormInstance } from 'element-plus'
import { ElMessage } from 'element-plus'
import { computed, ref, watch } from 'vue'

const props = defineProps<{
    modelValue: boolean
    roleId: string
    roleName: string
}>()

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'success'): void
}>()

const visible = computed({
    get: () => props.modelValue,
    set: (val) => emit('update:modelValue', val)
})

const treeRef = ref<FormInstance>()
const loading = ref(false)
const submitLoading = ref(false)
const allMenus = ref<any[]>([])

// 判断是否为根节点（兼容 '' 和 '0' 两种情况）
const isRoot = (pid: string) => !pid || pid === '' || pid === '0'

// 构建菜单树形数据
const buildMenuTree = (items: any[], parentId: string = ''): any[] => {
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
            children: buildMenuTree(items, item.id)
        }))
}

// 菜单树形数据
const menuTreeData = computed(() => buildMenuTree(allMenus.value))

// 加载菜单数据
const loadMenus = async () => {
    loading.value = true
    try {
        const menus = await getMenus()
        allMenus.value = menus || []
    } catch (error: any) {
        ElMessage.error(error.message || '加载菜单失败')
    } finally {
        loading.value = false
    }
}

// 加载角色已配置的菜单
const loadRoleMenus = async () => {
    if (!props.roleId) return
    try {
        const res = await getRoleMenus(props.roleId)
        const menuIds = res.menu_ids || []
        // 设置选中的节点
        if (treeRef.value) {
            treeRef.value.setCheckedKeys(menuIds, false)
        }
    } catch (error: any) {
        ElMessage.error(error.message || '加载角色菜单失败')
    }
}

// 弹窗打开时加载数据
const handleOpen = async () => {
    await loadMenus()
    await loadRoleMenus()
}

const handleClose = () => {
    allMenus.value = []
    emit('update:modelValue', false)
}

const handleSubmit = async () => {
    if (!props.roleId) {
        ElMessage.error('角色ID不能为空')
        return
    }

    submitLoading.value = true
    try {
        // 获取选中的菜单ID（包括半选中的父节点）
        const checkedKeys = treeRef.value?.getCheckedKeys(false) || []
        const halfCheckedKeys = treeRef.value?.getHalfCheckedKeys() || []
        const allMenuIds = [...checkedKeys, ...halfCheckedKeys] as string[]

        await updateRoleMenus(props.roleId, allMenuIds)
        ElMessage.success('角色菜单配置成功')
        handleClose()
        emit('success')
    } catch (error: any) {
        ElMessage.error(error.message || '配置失败')
    } finally {
        submitLoading.value = false
    }
}
</script>

<style scoped>
.menu-node {
    display: flex;
    align-items: center;
}

.menu-name {
    font-size: 14px;
}
</style>
