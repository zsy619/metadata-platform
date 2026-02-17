<template>
  <el-dialog v-model="visible" :title="title" width="600px" destroy-on-close @close="handleClose">
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" label-position="right">
      <el-tabs v-model="activeTab">
        <!-- 基本信息 -->
        <el-tab-pane label="基本信息" name="basic" style="height: 300px;">
          <el-form-item label="所属应用" prop="app_code">
            <el-tree-select
              v-model="formData.app_code"
              :data="appTreeData"
              check-strictly
              :render-after-expand="false"
              placeholder="请选择所属应用（不选则不属于任何应用）"
              style="width: 100%"
              clearable
              @change="handleAppChange"
            />
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
          <el-form-item label="菜单名称" prop="menu_name">
            <el-input v-model="formData.menu_name" placeholder="请输入菜单名称" />
          </el-form-item>
          <el-form-item label="菜单标识" prop="menu_code">
            <el-input v-model="formData.menu_code" placeholder="请输入菜单标识" />
          </el-form-item>
          <el-form-item label="菜单类型" prop="menu_type">
            <el-radio-group v-model="formData.menu_type">
              <el-radio value="M">目录</el-radio>
              <el-radio value="C">菜单</el-radio>
              <el-radio value="F">按钮</el-radio>
              <el-radio value="Z">资源</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="图&#12288;&#12288;标" prop="icon">
            <IconPicker v-model="formData.icon" />
          </el-form-item>
        </el-tab-pane>
        <!-- 路由配置 -->
        <el-tab-pane label="路由配置" name="route" style="height: 300px;" v-if="formData.menu_type !== 'F' && formData.menu_type !== 'Z'">
          <el-form-item label="路由地址" prop="url">
            <el-input v-model="formData.url" placeholder="请输入路由地址" />
          </el-form-item>
          <el-form-item label="是否可见" prop="is_visible">
            <el-switch v-model="formData.is_visible" />
          </el-form-item>
        </el-tab-pane>
        <!-- 数据范围 -->
        <el-tab-pane label="数据范围" name="dataScope" style="height: 300px;">
          <el-form-item label="权限范围" prop="data_range">
            <el-select v-model="formData.data_range" style="width: 100%">
              <el-option v-for="item in DATA_RANGE_OPTIONS" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="自定义范围" prop="data_scope" v-if="formData.data_range === '2'">
            <el-tree-select
              v-model="selectedOrgIds"
              :data="orgTreeData"
              multiple
              check-strictly
              default-expand-all
              :render-after-expand="true"
              placeholder="请选择数据范围（组织）"
              clearable
              collapse-tags
              collapse-tags-tooltip
              style="width: 100%"
              :filter-method="filterOrg"
            />
          </el-form-item>
        </el-tab-pane>
        <!-- 其他设置 -->
        <el-tab-pane label="其他设置" name="other" style="height: 300px;">
          <el-form-item label="状&#12288;&#12288;态" prop="status">
            <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
          </el-form-item>
          <el-form-item label="序&#12288;&#12288;号" prop="sort">
            <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
          </el-form-item>
          <el-form-item label="备&#12288;&#12288;注" prop="remark">
            <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入备注" />
          </el-form-item>
        </el-tab-pane>
      </el-tabs>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="loading">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { createMenu, updateMenu } from '@/api/user';
import IconPicker from '@/components/icon/IconPicker.vue';
import { DATA_RANGE, DATA_RANGE_OPTIONS } from '@/utils/constants';
import type { FormInstance, FormRules } from 'element-plus';
import { ElMessage } from 'element-plus';
import { computed, ref, watch } from 'vue';

const props = defineProps<{
  modelValue: boolean
  data: any
  menuList: any[]
  appList: any[]
  orgList?: any[]
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
const activeTab = ref('basic')
const formData = ref<any>({
  parent_id: '',
  app_code: '',
  menu_type: 'C',
  menu_name: '',
  menu_code: '',
  url: '',
  icon: '',
  is_visible: true,
  data_range: DATA_RANGE.ALL,
  data_scope: '',
  status: 1,
  sort: 0,
  remark: ''
})

const formRules: FormRules = {
  menu_name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  menu_code: [{ required: true, message: '请输入权限标识', trigger: 'blur' }],
  menu_type: [{ required: true, message: '请选择菜单类型', trigger: 'change' }]
}

const appTreeData = computed(() => {
  const buildTree = (items: any[], parentId = ''): any[] => {
    return items
      .filter(item => (parentId === '' ? !item.parent_id || item.parent_id === '' : item.parent_id === parentId))
      .sort((a, b) => a.sort - b.sort)
      .map(item => ({
        value: item.app_code,
        label: item.app_name,
        children: buildTree(items, item.id)
      }))
  }
  return buildTree(props.appList)
})

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

const buildOrgTreeData = (list: any[], parentId: string = ''): any[] => {
  return list
    .filter(item => (item.parent_id || '') === parentId)
    .map(item => ({
      value: item.id,
      label: item.org_name || item.org_code,
      children: buildOrgTreeData(list, item.id)
    }))
    .sort((a, b) => {
      const itemA = list.find(i => i.id === a.value)
      const itemB = list.find(i => i.id === b.value)
      return (itemA?.sort || 0) - (itemB?.sort || 0)
    })
}

const orgTreeData = computed(() => {
  if (!props.orgList) return []
  return buildOrgTreeData(props.orgList)
})

const selectedOrgIds = computed({
  get: () => {
    if (!formData.value.data_scope) return []
    if (Array.isArray(formData.value.data_scope)) return formData.value.data_scope
    return formData.value.data_scope.split(',').filter(Boolean)
  },
  set: (val: string[]) => {
    formData.value.data_scope = val.join(',')
  }
})

const filterOrg = (node: any, keyword: string) => {
  return node.label.toLowerCase().includes(keyword.toLowerCase())
}

watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      activeTab.value = 'basic'
      if (props.data) {
        formData.value = { ...props.data }
      } else {
        formData.value = {
          parent_id: '',
          app_code: '',
          menu_type: 'C',
          menu_name: '',
          menu_code: '',
          url: '',
          icon: '',
          is_visible: true,
          data_range: DATA_RANGE.ALL,
          data_scope: '',
          status: 1,
          sort: 0,
          remark: ''
        }
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
