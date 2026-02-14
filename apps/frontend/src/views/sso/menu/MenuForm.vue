<template>
  <el-dialog v-model="visible" :title="title" width="550px" destroy-on-close @close="handleClose">
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" label-position="right">
      <el-form-item label="上级菜单" prop="parent_id">
        <el-tree-select v-model="formData.parent_id" :data="menuTreeData" check-strictly :render-after-expand="false" placeholder="请选择上级菜单" style="width: 100%" />
      </el-form-item>
      <el-form-item label="菜单类型" prop="menu_type">
        <el-radio-group v-model="formData.menu_type">
          <el-radio label="M">目录</el-radio>
          <el-radio label="C">菜单</el-radio>
          <el-radio label="F">按钮</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="菜单名称" prop="menu_name">
        <el-input v-model="formData.menu_name" placeholder="请输入菜单名称" />
      </el-form-item>
      <el-form-item label="权限标识" prop="menu_code">
        <el-input v-model="formData.menu_code" placeholder="请输入权限标识" />
      </el-form-item>
      <el-form-item v-if="formData.menu_type !== 'F'" label="路由地址" prop="url">
        <el-input v-model="formData.url" placeholder="请输入路由地址" />
      </el-form-item>
      <el-form-item v-if="formData.menu_type !== 'F'" label="显示状态" prop="visible">
        <el-switch v-model="formData.visible" :active-value="1" :inactive-value="0" />
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
      </el-form-item>
      <el-form-item label="排序" prop="sort">
        <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
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
  menuTreeData: any[]
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
  menu_type: 'C',
  menu_name: '',
  menu_code: '',
  url: '',
  visible: 1,
  status: 1,
  sort: 0
})

const formRules: FormRules = {
  menu_name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }]
}

watch(
  () => props.modelValue,
  (val) => {
    if (val && props.data) {
      formData.value = { ...props.data }
    } else if (val) {
      formData.value = { parent_id: '', menu_type: 'C', menu_name: '', menu_code: '', url: '', visible: 1, status: 1, sort: 0 }
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
