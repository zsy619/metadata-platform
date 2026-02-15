<template>
  <el-dialog v-model="visible" :title="title" width="600px" destroy-on-close @close="handleClose">
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" label-position="right">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="basic">
          <el-form-item label="上级组织" prop="parent_id">
            <div style="width: 100%; display: flex; align-items: center; gap: 8px;">
              <el-tree-select
                v-model="formData.parent_id"
                :data="processedTreeData"
                check-strictly
                :render-after-expand="false"
                placeholder="请选择上级组织（不选则为顶级组织）"
                clearable
                style="flex: 1"
              >
                <template #default="{ data }">
                  <span>{{ data.label }} - {{ data.org_name }}</span>
                </template>
              </el-tree-select>
              <span v-if="selectedParentName" style="color: var(--el-text-color-secondary); font-size: 14px; white-space: nowrap;">
                {{ selectedParentName }}
              </span>
            </div>
          </el-form-item>
          <el-form-item label="组织名称" prop="org_name">
            <el-input v-model="formData.org_name" placeholder="请输入组织名称" />
          </el-form-item>
          <el-form-item label="组织简称" prop="org_short">
            <el-input v-model="formData.org_short" placeholder="请输入组织简称" />
          </el-form-item>
          <el-form-item label="组织编码" prop="org_code">
            <el-input v-model="formData.org_code" placeholder="请输入组织编码" />
          </el-form-item>
          <el-form-item label="类型编码" prop="kind_code">
            <el-input v-model="formData.kind_code" placeholder="请输入类型编码" />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane label="联系信息" name="contact">
          <el-form-item label="联&nbsp;&nbsp;系&nbsp;&nbsp;人" prop="contact">
            <el-input v-model="formData.contact" placeholder="请输入联系人" />
          </el-form-item>
          <el-form-item label="联系电话" prop="phone">
            <el-input v-model="formData.phone" placeholder="请输入联系电话" />
          </el-form-item>
          <el-form-item label="联系地址" prop="address">
            <el-input v-model="formData.address" placeholder="请输入联系地址" />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane label="其他设置" name="other">
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
import { createUnit, updateUnit } from '@/api/user';
import type { FormInstance, FormRules } from 'element-plus';
import { ElMessage } from 'element-plus';
import { computed, ref, watch } from 'vue';

const props = defineProps<{
  modelValue: boolean
  data: any
  orgTreeData: any[]
  excludeIds?: string[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const title = computed(() => (props.data?.id ? '编辑组织' : '新增组织'))

const formRef = ref<FormInstance>()
const loading = ref(false)
const activeTab = ref('basic')
const formData = ref<any>({
  parent_id: '0',
  org_name: '',
  org_short: '',
  org_code: '',
  kind_code: '',
  contact: '',
  phone: '',
  address: '',
  status: 1,
  sort: 0,
  remark: ''
})

const formRules: FormRules = {
  org_name: [{ required: true, message: '请输入组织名称', trigger: 'blur' }],
  org_code: [{ required: true, message: '请输入组织编码', trigger: 'blur' }]
}

const processedTreeData = computed(() => {
  if (!props.excludeIds || props.excludeIds.length === 0) {
    return props.orgTreeData
  }
  const excludeSet = new Set(props.excludeIds)
  const filterTree = (nodes: any[]): any[] => {
    return nodes
      .filter(node => !excludeSet.has(node.value))
      .map(node => ({
        ...node,
        disabled: excludeSet.has(node.value),
        children: node.children ? filterTree(node.children) : []
      }))
  }
  return filterTree(props.orgTreeData)
})

const selectedParentName = computed(() => {
  if (!formData.value.parent_id) {
    return ''
  }
  const findNode = (nodes: any[]): string => {
    for (const node of nodes) {
      if (node.value === formData.value.parent_id) {
        return node.org_name || node.label
      }
      if (node.children?.length) {
        const found = findNode(node.children)
        if (found) return found
      }
    }
    return ''
  }
  return findNode(processedTreeData.value)
})

watch(
  () => props.modelValue,
  (val) => {
    if (val && props.data) {
      const data = { ...props.data }
      if (data.parent_id === '0' || data.parent_id === 0 || data.parent_id === '') {
        data.parent_id = ''
      }
      formData.value = data
    } else if (val) {
      formData.value = { parent_id: '', org_name: '', org_short: '', org_code: '', kind_code: '', contact: '', phone: '', address: '', status: 1, sort: 0, remark: '' }
    }
    if (val) {
      activeTab.value = 'basic'
    }
  }
)

const handleClose = () => {
  formRef.value?.resetFields()
  activeTab.value = 'basic'
  emit('update:modelValue', false)
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const submitData = { ...formData.value }
        if (!submitData.parent_id) {
          submitData.parent_id = ''
        }
        if (submitData.id) {
          await updateUnit(submitData.id, submitData)
          ElMessage.success('更新成功')
        } else {
          await createUnit(submitData)
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
