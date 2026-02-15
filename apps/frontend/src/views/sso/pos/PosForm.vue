<template>
  <el-dialog v-model="visible" :title="title" width="500px" destroy-on-close @close="handleClose">
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" label-position="right">
      <el-form-item label="上级职位" prop="parent_id">
        <el-tree-select
          v-model="formData.parent_id"
          :data="filteredTreeData"
          check-strictly
          :render-after-expand="false"
          placeholder="请选择上级职位（不选则为顶级职位）"
          clearable
          style="width: 100%"
        />
      </el-form-item>
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

/**
 * 树型选择器节点数据结构
 */
interface TreeNode {
  value: string;
  label: string;
  children?: TreeNode[];
}

const props = defineProps<{
  modelValue: boolean
  data: any
  posTreeData: TreeNode[]
  excludeIds: string[]
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

/**
 * 过滤树型数据，排除指定节点（编辑时排除自己及所有子节点）
 * 防止将父节点设置为自己的子节点，造成循环引用
 * @param data 原始树型数据
 * @param excludeIds 需要排除的节点ID列表
 * @returns 过滤后的树型数据
 */
const filterTreeData = (data: TreeNode[], excludeIds: string[]): TreeNode[] => {
  return data
    .filter(node => !excludeIds.includes(node.value))
    .map(node => ({
      ...node,
      children: node.children ? filterTreeData(node.children, excludeIds) : []
    }))
}

/**
 * 过滤后的树型选择器数据
 * 编辑时排除当前节点及其所有子节点，防止形成循环引用
 */
const filteredTreeData = computed(() => {
  if (!props.excludeIds || props.excludeIds.length === 0) {
    return props.posTreeData
  }
  return filterTreeData(props.posTreeData, props.excludeIds)
})

const formRef = ref<FormInstance>()
const loading = ref(false)
const formData = ref<any>({
  parent_id: '',
  pos_name: '',
  pos_code: '',
  data_range: DATA_RANGE.ALL,
  status: 1,
  sort: 0
})

const formRules: FormRules = {
  pos_name: [{ required: true, message: '请输入职位名称', trigger: 'blur' }],
  pos_code: [{ required: true, message: '请输入职位编码', trigger: 'blur' }]
}

watch(
  () => props.modelValue,
  (val) => {
    if (val && props.data) {
      formData.value = { ...props.data }
    } else if (val) {
      formData.value = { parent_id: '', pos_name: '', pos_code: '', data_range: DATA_RANGE.ALL, status: 1, sort: 0 }
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
