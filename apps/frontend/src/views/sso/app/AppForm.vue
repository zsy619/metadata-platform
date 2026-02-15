<template>
  <el-dialog v-model="visible" :title="title" width="550px" destroy-on-close @close="handleClose">
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" label-position="right">
      <el-form-item label="父级应用" prop="parent_id">
        <el-tree-select
          v-model="formData.parent_id"
          :data="processedTreeSelectData"
          :render-after-expand="false"
          check-strictly
          placeholder="请选择父级应用（不选则为顶级）"
          clearable
          style="width: 100%"
          :filter-method="filterApp"
        />
      </el-form-item>
      <el-form-item label="应用名称" prop="app_name">
        <el-input v-model="formData.app_name" placeholder="请输入应用名称" />
      </el-form-item>
      <el-form-item label="应用编码" prop="app_code">
        <el-input v-model="formData.app_code" placeholder="请输入应用编码" />
      </el-form-item>
      <el-form-item label="域&nbsp;名&nbsp;/&nbsp;IP" prop="host">
        <el-input v-model="formData.host" placeholder="请输入域名或IP地址" />
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
import { createApp, updateApp } from '@/api/user';
import type { FormInstance, FormRules } from 'element-plus';
import { ElMessage } from 'element-plus';
import { computed, ref, watch } from 'vue';

const props = defineProps<{
  modelValue: boolean
  data: any
  allData: any[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const title = computed(() => (props.data?.id ? '编辑应用' : '新增应用'))

const formRef = ref<FormInstance>()
const loading = ref(false)
const formData = ref<any>({
  parent_id: '',
  app_name: '',
  app_code: '',
  host: '',
  status: 1,
  sort: 0,
  remark: ''
})

const formRules: FormRules = {
  app_name: [{ required: true, message: '请输入应用名称', trigger: 'blur' }],
  app_code: [{ required: true, message: '请输入应用编码', trigger: 'blur' }]
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

const buildTreeSelectData = (list: any[], parentId: string = '', excludeIds: string[] = []): any[] => {
  return list
    .filter(item => {
      const pid = item.parent_id || ''
      return pid === parentId && !excludeIds.includes(item.id)
    })
    .map(item => ({
      value: item.id,
      label: item.app_name,
      children: buildTreeSelectData(list, item.id, excludeIds)
    }))
    .sort((a, b) => {
      const itemA = list.find(i => i.id === a.value)
      const itemB = list.find(i => i.id === b.value)
      return (itemA?.sort || 0) - (itemB?.sort || 0)
    })
}

const processedTreeSelectData = computed(() => {
  const excludeIds = getDisabledIds()
  return buildTreeSelectData(props.allData, '', excludeIds)
})

const filterApp = (node: any, keyword: string) => {
  return node.label.toLowerCase().includes(keyword.toLowerCase())
}

const getMaxSort = (): number => {
  if (!props.allData || props.allData.length === 0) return 1
  const maxSort = Math.max(...props.allData.map(item => item.sort || 0))
  return maxSort + 1
}

watch(
  () => props.modelValue,
  (val) => {
    if (val && props.data) {
      const data = { ...props.data }
      if (data.parent_id === '0' || data.parent_id === 0) {
        data.parent_id = ''
      }
      formData.value = data
    } else if (val) {
      formData.value = { parent_id: '', app_name: '', app_code: '', host: '', status: 1, sort: getMaxSort(), remark: '' }
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
          ElMessage.error('不能选择当前应用的下级作为父级应用')
          return
        }
      }
      loading.value = true
      try {
        const submitData = { ...formData.value }
        if (!submitData.parent_id) {
          submitData.parent_id = ''
        }
        if (formData.value.id) {
          await updateApp(formData.value.id, submitData)
          ElMessage.success('更新成功')
        } else {
          await createApp(submitData)
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
