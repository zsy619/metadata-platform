<template>
  <el-dialog v-model="visible" :title="title" width="500px" destroy-on-close @close="handleClose">
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" label-position="right">
      <!-- 上级角色选择器 -->
      <el-form-item label="上级角色" prop="parent_id">
        <el-tree-select
          v-model="formData.parent_id"
          :data="treeSelectData"
          :render-after-expand="false"
          check-strictly
          placeholder="请选择上级角色（不选则为顶级）"
          clearable
          style="width: 100%"
          :filter-method="filterRole"
        />
      </el-form-item>
      <el-form-item label="角色名称" prop="role_name">
        <el-input v-model="formData.role_name" placeholder="请输入角色名称" />
      </el-form-item>
      <el-form-item label="角色编码" prop="role_code">
        <el-input v-model="formData.role_code" placeholder="请输入角色编码" />
      </el-form-item>
      <el-form-item label="数据范围" prop="data_range">
        <el-select v-model="formData.data_range" style="width: 100%">
          <el-option v-for="item in DATA_RANGE_OPTIONS" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="数据权限" prop="data_scope" v-if="formData.data_range === '2'">
        <el-tree-select v-model="selectedOrgIds" :data="orgTreeData" multiple check-strictly default-expand-all :render-after-expand="true" placeholder="请选择数据范围（组织）" clearable collapse-tags collapse-tags-tooltip style="width: 100%" />
      </el-form-item>
      <el-form-item label="状　　态" prop="status">
        <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
      </el-form-item>
      <el-form-item label="序　　号" prop="sort">
        <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
      </el-form-item>
      <el-form-item label="备　　注" prop="remark">
        <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入备注" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="loading">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { createRole, updateRole } from '@/api/user'
import { DATA_RANGE, DATA_RANGE_OPTIONS } from '@/utils/constants'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { computed, ref, watch } from 'vue'

const props = defineProps<{
  modelValue: boolean
  data: any
  allRoles: any[]
  excludeIds?: string[]
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

const title = computed(() => (props.data?.id ? '编辑角色' : '新增角色'))

const formRef = ref<FormInstance>()
const loading = ref(false)
const formData = ref<any>({
  parent_id: '',
  role_name: '',
  role_code: '',
  data_range: DATA_RANGE.ALL,
  status: 1,
  sort: 0,
  remark: ''
})

const formRules: FormRules = {
  role_name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  role_code: [{ required: true, message: '请输入角色编码', trigger: 'blur' }]
}

/**
 * 构建树形选择器数据
 * 排除当前节点及其所有子节点，防止循环引用
 * @param list 扁平数据列表
 * @param parentId 父ID
 * @param excludeIds 需要排除的ID列表
 * @returns 树形选择器数据
 */
const buildTreeSelectData = (list: any[], parentId: string = '', excludeIds: string[] = []): any[] => {
  return list
    .filter(item => {
      const pid = item.parent_id || ''
      return pid === parentId && !excludeIds.includes(item.id)
    })
    .map(item => ({
      value: item.id,
      label: item.role_name,
      children: buildTreeSelectData(list, item.id, excludeIds)
    }))
    .sort((a, b) => {
      const itemA = list.find(i => i.id === a.value)
      const itemB = list.find(i => i.id === b.value)
      return (itemA?.sort || 0) - (itemB?.sort || 0)
    })
}

/**
 * 树形选择器数据（计算属性）
 * 使用 props 传入的 excludeIds 排除当前编辑节点及其所有子节点
 */
const treeSelectData = computed(() => {
  const excludeIds = props.excludeIds || []
  return buildTreeSelectData(props.allRoles, '', excludeIds)
})

/**
 * 构建组织树形数据
 * @param list 组织列表
 * @param parentId 父ID
 * @returns 树形数据
 */
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

/**
 * 组织树形数据（用于数据权限选择）
 */
const orgTreeData = computed(() => {
  if (!props.orgList) return []
  return buildOrgTreeData(props.orgList)
})

/**
 * 选中的组织ID（用于数据权限）
 */
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

/**
 * 树形选择器过滤方法
 * @param node 节点数据
 * @param keyword 搜索关键词
 */
const filterRole = (node: any, keyword: string) => {
  return node.label.toLowerCase().includes(keyword.toLowerCase())
}

watch(
  () => props.modelValue,
  (val) => {
    if (val && props.data) {
      formData.value = { ...props.data }
    } else if (val) {
      formData.value = { 
        parent_id: '', 
        role_name: '', 
        role_code: '', 
        data_range: DATA_RANGE.ALL, 
        status: 1, 
        sort: 0, 
        remark: '' 
      }
    }
  }
)

const handleClose = () => {
  formRef.value?.resetFields()
  emit('update:modelValue', false)
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      // 验证上级角色选择是否合法（不能选择自己或自己的子节点作为上级）
      if (formData.value.id && formData.value.parent_id && props.excludeIds) {
        if (props.excludeIds.includes(formData.value.parent_id)) {
          ElMessage.error('不能选择当前角色或其下级角色作为上级角色')
          return
        }
      }
      
      loading.value = true
      try {
        // 构造提交数据
        const submitData = {
          parent_id: formData.value.parent_id || '',
          role_name: formData.value.role_name || '',
          role_code: formData.value.role_code || '',
          data_range: formData.value.data_range,
          data_scope: formData.value.data_scope || '',
          status: formData.value.status,
          sort: formData.value.sort,
          remark: formData.value.remark || ''
        }
        
        if (formData.value.id) {
          await updateRole(formData.value.id, submitData)
          ElMessage.success('更新成功')
        } else {
          await createRole(submitData)
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
