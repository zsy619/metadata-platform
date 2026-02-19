<template>
    <el-dialog v-model="visible" :title="title" width="500px" destroy-on-close @close="handleClose">
        <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" label-position="right">
            <el-form-item label="上级角色组" prop="parent_id">
                <el-tree-select
                    v-model="formData.parent_id"
                    :data="groupTreeData"
                    :props="{ children: 'children', label: 'group_name', value: 'id' }"
                    check-strictly
                    clearable
                    placeholder="请选择上级角色组（不选则为顶级）"
                    style="width: 100%"
                />
            </el-form-item>
            <el-form-item label="角色组名称" prop="group_name">
                <el-input v-model="formData.group_name" placeholder="请输入角色组名称" />
            </el-form-item>
            <el-form-item label="角色组编码" prop="group_code">
                <el-input v-model="formData.group_code" placeholder="请输入角色组编码" />
            </el-form-item>
            <el-form-item label="状&#12288;&#12288;&#12288;态" prop="status">
                <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
            </el-form-item>
            <el-form-item label="序&#12288;&#12288;&#12288;号" prop="sort">
                <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
            </el-form-item>
            <el-form-item label="备&#12288;&#12288;&#12288;注" prop="remark">
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
import { createRoleGroup, updateRoleGroup } from '@/api/user';
import type { FormInstance, FormRules } from 'element-plus';
import { ElMessage } from 'element-plus';
import { computed, ref, watch } from 'vue';

const props = defineProps<{
    modelValue: boolean
    data: any
    allGroups: any[]
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

const title = computed(() => (props.data?.id ? '编辑角色组' : '新增角色组'))

const formRef = ref<FormInstance>()
const loading = ref(false)
const formData = ref<any>({
    parent_id: '',
    group_name: '',
    group_code: '',
    status: 1,
    sort: 0,
    remark: ''
})

const formRules: FormRules = {
    group_name: [{ required: true, message: '请输入角色组名称', trigger: 'blur' }],
    group_code: [{ required: true, message: '请输入角色组编码', trigger: 'blur' }]
}

const isRoot = (pid: string) => !pid || pid === '' || pid === '0'

const buildGroupTree = (items: any[], parentId: string = ''): any[] => {
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
            children: buildGroupTree(items, item.id)
        }))
}

const groupTreeData = computed(() => {
    const filteredGroups = props.allGroups.filter(g => !props.excludeIds.includes(g.id))
    return buildGroupTree(filteredGroups)
})

watch(
    () => props.modelValue,
    (val) => {
        if (val && props.data) {
            formData.value = { ...props.data }
        } else if (val) {
            formData.value = { parent_id: '', group_name: '', group_code: '', status: 1, sort: 0, remark: '' }
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
                    await updateRoleGroup(formData.value.id, formData.value)
                    ElMessage.success('更新成功')
                } else {
                    await createRoleGroup(formData.value)
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
