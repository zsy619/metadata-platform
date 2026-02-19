<template>
  <el-dialog v-model="visible" title="配置角色菜单" width="600px" destroy-on-close @close="handleClose" @open="handleOpen">
    <div v-loading="loading" class="tree-container">
      <el-tree
        ref="treeRef"
        :data="menuTreeData"
        :props="{ children: 'children', label: 'menu_name' }"
        show-checkbox
        node-key="id"
        :default-expand-all="isExpanded"
        :check-strictly="false"
        highlight-current
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
      <div style="display: flex; justify-content: space-between; width: 100%;">
        <div>
          <el-button @click="toggleExpand">{{ isExpanded ? '折叠' : '展开' }}</el-button>
          <el-button @click="handleSelectAll">全选</el-button>
          <el-button @click="handleClearAll">取消</el-button>
        </div>
        <div>
          <el-button @click="handleClose">关闭</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { getMenus, getRoleMenus, updateRoleMenus } from '@/api/user'
import type { ElTree } from 'element-plus'
import { ElMessage } from 'element-plus'
import { computed, ref, nextTick } from 'vue'

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

const treeRef = ref<InstanceType<typeof ElTree>>()
const loading = ref(false)
const submitLoading = ref(false)
const allMenus = ref<any[]>([])
const isExpanded = ref(true)

const isRoot = (pid: string) => !pid || pid === '' || pid === '0'

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

const menuTreeData = computed(() => buildMenuTree(allMenus.value))

const getAllNodeIds = (nodes: any[]): string[] => {
  let ids: string[] = []
  for (const node of nodes) {
    ids.push(node.id)
    if (node.children && node.children.length > 0) {
      ids = ids.concat(getAllNodeIds(node.children))
    }
  }
  return ids
}

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

const loadRoleMenus = async () => {
  if (!props.roleId) return
  try {
    const res = await getRoleMenus(props.roleId)
    const menuIds = res.menu_ids || []
    if (treeRef.value) {
      treeRef.value.setCheckedKeys(menuIds, false)
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载角色菜单失败')
  }
}

const handleOpen = async () => {
  isExpanded.value = true
  await loadMenus()
  await loadRoleMenus()
}

const handleClose = () => {
  allMenus.value = []
  emit('update:modelValue', false)
}

const toggleExpand = () => {
  isExpanded.value = !isExpanded.value
  const nodes = treeRef.value?.store?.nodesMap
  if (nodes) {
    Object.values(nodes).forEach((node: any) => {
      node.expanded = isExpanded.value
    })
  }
}

const handleSelectAll = () => {
  const allIds = getAllNodeIds(menuTreeData.value)
  treeRef.value?.setCheckedKeys(allIds, false)
}

const handleClearAll = () => {
  treeRef.value?.setCheckedKeys([], false)
}

const handleSubmit = async () => {
  if (!props.roleId) {
    ElMessage.error('角色ID不能为空')
    return
  }

  submitLoading.value = true
  try {
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
.tree-container {
  height: 400px;
  overflow-y: auto;
}

.menu-node {
  display: flex;
  align-items: center;
}

.menu-name {
  font-size: 14px;
}
</style>
