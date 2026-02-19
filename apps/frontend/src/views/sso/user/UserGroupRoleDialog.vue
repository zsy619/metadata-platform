<template>
  <el-dialog v-model="visible" title="关联角色" width="500px" :close-on-click-modal="false" @close="handleClose">
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
      <span class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { getUserGroupRoles, getRoles, updateUserGroupRoles } from '@/api/user'
import { Folder, User } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { ElTree } from 'element-plus'
import { computed, nextTick, ref, watch } from 'vue'

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

/**
 * 构建角色树形数据用于树组件
 */
const roleTreeData = computed(() => {
  const buildTree = (items: any[], parentId = ''): any[] => {
    return items
      .filter(item => {
        const itemPid = item.parent_id || ''
        return parentId === '' ? !itemPid || itemPid === '' || itemPid === '0' : itemPid === parentId
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

/**
 * 加载角色列表和当前用户组已关联的角色
 */
const loadData = async () => {
  if (!props.groupId) return
  
  loading.value = true
  try {
    // 并行加载角色列表和已关联角色
    const [rolesRes, groupRoles] = await Promise.all([
      getRoles(),
      getUserGroupRoles(props.groupId)
    ])
    
    // 处理角色列表数据
    allRoles.value = Array.isArray(rolesRes) ? rolesRes : (rolesRes as any).data || []
    
    // 等待 DOM 更新后设置选中状态
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

/**
 * 监听弹窗显示状态
 */
watch(visible, (val) => {
  if (val && props.groupId) {
    loadData()
  }
})

/**
 * 提交更新
 */
const handleSubmit = async () => {
  if (!treeRef.value) {
    ElMessage.error('树组件未初始化')
    return
  }
  
  const checkedKeys = treeRef.value.getCheckedKeys(false)
  
  submitting.value = true
  try {
    await updateUserGroupRoles(props.groupId, checkedKeys as string[])
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

/**
 * 关闭弹窗
 */
const handleClose = () => {
  visible.value = false
  allRoles.value = []
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
