<template>
  <el-dialog v-model="visible" :title="title" width="600px" destroy-on-close @close="handleClose">
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" label-position="right">
      <el-form-item label="所属应用" prop="app_code">
        <el-select v-model="formData.app_code" placeholder="请选择所属应用" style="width: 100%" @change="handleAppChange">
          <el-option v-for="app in appList" :key="app.id" :label="app.app_name" :value="app.app_code" />
        </el-select>
      </el-form-item>
      <el-form-item label="上级菜单" prop="parent_id">
        <el-tree-select
          v-model="formData.parent_id"
          :data="filteredMenuTreeData"
          check-strictly
          :render-after-expand="false"
          placeholder="请选择上级菜单（不选则为顶级菜单）"
          style="width: 100%"
          clearable
        />
      </el-form-item>
      <el-form-item label="菜单类型" prop="menu_type">
        <el-radio-group v-model="formData.menu_type">
          <el-radio label="M">目录</el-radio>
          <el-radio label="C">菜单</el-radio>
          <el-radio label="F">按钮</el-radio>
          <el-radio label="Z">资源</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="菜单名称" prop="menu_name">
        <el-input v-model="formData.menu_name" placeholder="请输入菜单名称" />
      </el-form-item>
      <el-form-item label="权限标识" prop="menu_code">
        <el-input v-model="formData.menu_code" placeholder="请输入权限标识" />
      </el-form-item>
      <el-form-item v-if="formData.menu_type !== 'F' && formData.menu_type !== 'Z'" label="路由地址" prop="url">
        <el-input v-model="formData.url" placeholder="请输入路由地址" />
      </el-form-item>
      <el-form-item v-if="formData.menu_type !== 'F' && formData.menu_type !== 'Z'" label="图标" prop="icon">
        <el-input v-model="formData.icon" placeholder="请输入图标名称" />
      </el-form-item>
      <el-form-item v-if="formData.menu_type !== 'F' && formData.menu_type !== 'Z'" label="是否可见" prop="is_visible">
        <el-switch v-model="formData.is_visible" />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
      </el-form-item>
      <el-form-item label="排序" prop="sort">
        <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
      </el-form-item>
      <el-form-item label="备注" prop="remark">
        <el-input v-model="formData.remark" type="textarea" :rows="2" placeholder="请输入备注" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="loading">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { createMenu, updateMenu } from '@/api/user'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { computed, ref, watch } from 'vue'

const props = defineProps<{
  modelValue: boolean
  data: any
  menuList: any[]
  appList: any[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const title = computed(() => (props.data?.id ? '编辑菜单' : '新增菜单'))

const formRef = ref<FormInstance>()
const loading = ref(false)
const formData = ref<any>({
  parent_id: '',
  app_code: '',
  menu_type: 'C',
  menu_name: '',
  menu_code: '',
  url: '',
  icon: '',
  is_visible: true,
  status: 1,
  sort: 0,
  remark: ''
})

const formRules: FormRules = {
  app_code: [{ required: true, message: '请选择所属应用', trigger: 'change' }],
  menu_name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  menu_code: [{ required: true, message: '请输入权限标识', trigger: 'blur' }],
  menu_type: [{ required: true, message: '请选择菜单类型', trigger: 'change' }]
}

const getAllDescendantIds = (parentId: string): string[] => {
  const ids: string[] = [parentId]
  const findChildren = (pid: string) => {
    props.menuList.forEach(item => {
      if (item.parent_id === pid) {
        ids.push(item.id)
        findChildren(item.id)
      }
    })
  }
  findChildren(parentId)
  return ids
}

const filteredMenuTreeData = computed(() => {
  const currentId = formData.value.id
  const currentAppCode = formData.value.app_code

  let filteredItems = props.menuList.filter(item => {
    if (!currentAppCode) return true
    return item.app_code === currentAppCode
  })

  if (currentId) {
    const excludeIds = getAllDescendantIds(currentId)
    filteredItems = filteredItems.filter(item => !excludeIds.includes(item.id))
  }

  const buildTreeSelect = (items: any[], parentId = ''): any[] => {
    return items
      .filter(item => (parentId === '' ? !item.parent_id || item.parent_id === '' : item.parent_id === parentId))
      .sort((a, b) => a.sort - b.sort)
      .map(item => ({
        value: item.id,
        label: item.menu_name,
        children: buildTreeSelect(items, item.id)
      }))
  }

  const treeData = buildTreeSelect(filteredItems)
  return [
    { value: '', label: '顶级菜单', children: [] },
    ...treeData
  ]
})

watch(
  () => props.modelValue,
  (val) => {
    if (val && props.data) {
      formData.value = { ...props.data }
    } else if (val) {
      formData.value = {
        parent_id: '',
        app_code: '',
        menu_type: 'C',
        menu_name: '',
        menu_code: '',
        url: '',
        icon: '',
        is_visible: true,
        status: 1,
        sort: 0,
        remark: ''
      }
    }
  }
)

const handleAppChange = () => {
  formData.value.parent_id = ''
}

const handleClose = () => {
  formRef.value?.resetFields()
  emit('update:modelValue', false)
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        if (formData.value.id) {
          await updateMenu(formData.value.id, formData.value)
          ElMessage.success('更新成功')
        } else {
          await createMenu(formData.value)
          ElMessage.success('创建成功')
        }
        handleClose()
        emit('success')
      } catch (error: any) {
        ElMessage.error(error.message || '操作失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>
