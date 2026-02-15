<template>
  <el-dialog v-model="visible" :title="title" width="500px" destroy-on-close @close="handleClose">
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" label-position="right">
      <el-form-item label="上级类型" prop="parent_id">
        <el-tree-select
          v-model="formData.parent_id"
          :data="processedTreeData"
          check-strictly
          :render-after-expand="false"
          placeholder="请选择上级类型（不选则为顶级）"
          clearable
          style="width: 100%"
        />
      </el-form-item>
      <el-form-item label="类型名称" prop="name">
        <el-input v-model="formData.name" placeholder="请输入类型名称" />
      </el-form-item>
      <el-form-item label="类型编码" prop="code">
        <el-input v-model="formData.code" placeholder="请输入类型编码" />
      </el-form-item>
      <el-form-item label="状&#12288;&#12288;态" prop="status">
        <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
      </el-form-item>
      <el-form-item label="序&#12288;&#12288;号" prop="sort">
        <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
      </el-form-item>
      <el-form-item label="备&#12288;&#12288;注" prop="remark">
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
import { createOrgKind, updateOrgKind } from '@/api/user';
import type { FormInstance, FormRules } from 'element-plus';
import { ElMessage } from 'element-plus';
import { computed, ref, watch } from 'vue';

const props = defineProps<{
  modelValue: boolean
  data: any
  allData: any[]
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

const title = computed(() => (props.data?.id ? '编辑类型' : '新增类型'))

const formRef = ref<FormInstance>()
const loading = ref(false)
const formData = ref<any>({
  parent_id: '',
  name: '',
  code: '',
  status: 1,
  sort: 0,
  remark: ''
})

const formRules: FormRules = {
  name: [{ required: true, message: '请输入类型名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入类型编码', trigger: 'blur' }]
}

const getDisabledIds = (): string[] => {
  if (!formData.value.id) return []
  const disabledIds = [formData.value.id]
  const addChildren = (pid: string) => {
    props.allData.forEach(item => {
      const itemPid = item.parent_id || ''
      if (itemPid === pid) {
        disabledIds.push(item.id)
        addChildren(item.id)
      }
    })
  }
  addChildren(formData.value.id)
  return disabledIds
}

const buildTreeSelectData = (list: any[], parentId: string = '', excludeIds: Set<string> = new Set()): any[] => {
  return list
    .filter(item => {
      const pid = item.parent_id || ''
      return pid === parentId && !excludeIds.has(item.id)
    })
    .map(item => ({
      value: item.id,
      label: item.kind_name || item.name,
      children: buildTreeSelectData(list, item.id, excludeIds)
    }))
    .sort((a, b) => {
      const itemA = list.find(i => i.id === a.value)
      const itemB = list.find(i => i.id === b.value)
      return (itemA?.sort || 0) - (itemB?.sort || 0)
    })
}

const processedTreeData = computed(() => {
  if (!formData.value.id) {
    return buildTreeSelectData(props.allData, '', new Set())
  }
  const excludeIds = new Set(getDisabledIds())
  return buildTreeSelectData(props.allData, '', excludeIds)
})

watch(
  () => props.modelValue,
  (val) => {
    if (val && props.data) {
      const data = { ...props.data }
      if (data.kind_name && !data.name) {
        data.name = data.kind_name
      }
      if (data.kind_code && !data.code) {
        data.code = data.kind_code
      }
      formData.value = data
    } else if (val) {
      formData.value = { parent_id: '', name: '', code: '', status: 1, sort: 0, remark: '' }
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
      if (formData.value.id && formData.value.parent_id) {
        const disabledIds = getDisabledIds()
        if (disabledIds.includes(formData.value.parent_id)) {
          ElMessage.error('不能选择当前类型的下级作为父级类型')
          return
        }
      }
      loading.value = true
      try {
        if (formData.value.id) {
          await updateOrgKind(formData.value.id, formData.value)
          ElMessage.success('更新成功')
        } else {
          await createOrgKind(formData.value)
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
