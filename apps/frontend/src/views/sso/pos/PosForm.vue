<template>
  <el-dialog v-model="visible" :title="title" width="500px" destroy-on-close @close="handleClose">
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" label-position="right">
      <el-form-item label="职位名称" prop="pos_name">
        <el-input v-model="formData.pos_name" placeholder="请输入职位名称" />
      </el-form-item>
      <el-form-item label="职位编码" prop="pos_code">
        <el-input v-model="formData.pos_code" placeholder="请输入职位编码" />
      </el-form-item>
      <el-form-item label="数据范围" prop="data_range">
        <el-select v-model="formData.data_range" style="width: 100%">
          <el-option v-for="item in DATA_RANGE_OPTIONS" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
      </el-form-item>
      <el-form-item label="序号" prop="sort">
        <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
      </el-form-item>
      <el-form-item label="备注" prop="remark">
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
import { createPos, updatePos } from '@/api/user';
import { DATA_RANGE, DATA_RANGE_OPTIONS } from '@/utils/constants';
import type { FormInstance, FormRules } from 'element-plus';
import { ElMessage } from 'element-plus';
import { computed, ref, watch } from 'vue';

const props = defineProps<{
  modelValue: boolean
  data: any
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const title = computed(() => (props.data?.id ? '编辑职位' : '新增职位'))

const formRef = ref<FormInstance>()
const loading = ref(false)
const formData = ref<any>({
  pos_name: '',
  pos_code: '',
  data_range: DATA_RANGE.ALL,
  status: 1,
  sort: 0
})

const formRules: FormRules = {
  pos_name: [{ required: true, message: '请输入职位名称', trigger: 'blur' }]
}

watch(
  () => props.modelValue,
  (val) => {
    if (val && props.data) {
      formData.value = { ...props.data }
    } else if (val) {
      formData.value = { pos_name: '', pos_code: '', data_range: DATA_RANGE.ALL, status: 1, sort: 0 }
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
      loading.value = true
      try {
        if (formData.value.id) {
          await updatePos(formData.value.id, formData.value)
          ElMessage.success('更新成功')
        } else {
          await createPos(formData.value)
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
