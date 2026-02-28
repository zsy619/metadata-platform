<template>
    <el-dialog v-model="dialogVisible" title="密钥详情" width="600px">
        <el-descriptions :column="2" border>
            <el-descriptions-item label="密钥名称">{{ viewData.key_name }}</el-descriptions-item>
            <el-descriptions-item label="密钥类型">
                <el-tag>{{ getKeyTypeLabel(viewData.key_type) }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="用途">
                <el-tag :type="viewData.key_usage === 'signing' ? 'success' : 'warning'">
                    {{ getKeyUsageLabel(viewData.key_usage) }}
                </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="算法">{{ viewData.algorithm }}</el-descriptions-item>
            <el-descriptions-item label="Key ID" :span="2">{{ viewData.key_id }}</el-descriptions-item>
            <el-descriptions-item label="状态">
                <el-tag :type="viewData.is_enabled ? 'success' : 'danger'">
                    {{ viewData.is_enabled ? '启用' : '禁用' }}
                </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ viewData.create_at }}</el-descriptions-item>
            <el-descriptions-item label="公钥" :span="2" v-if="viewData.public_key">
                <el-input v-model="viewData.public_key" type="textarea" :rows="4" readonly />
            </el-descriptions-item>
            <el-descriptions-item label="备注" :span="2">{{ viewData.remark }}</el-descriptions-item>
        </el-descriptions>
    </el-dialog>
</template>

<script setup lang="ts">
import type { SsoKey } from '@/types/sso'
import { computed } from 'vue'

// Props
interface Props {
    modelValue: boolean
    data: Partial<SsoKey> | null
}

const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
    data: null
})

// Emits
const emit = defineEmits<{
    'update:modelValue': [value: boolean]
}>()

// 内部状态
const dialogVisible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
})

const viewData = computed(() => props.data || {})

// 工具函数
const getKeyTypeLabel = (type: string | undefined) => {
    if (!type) return ''
    const typeMap: Record<string, string> = { 'rsa': 'RSA', 'ec': 'EC', 'octet': 'Octet' }
    return typeMap[type.toLowerCase()] || type
}

const getKeyUsageLabel = (usage: string | undefined) => {
    if (!usage) return ''
    const usageMap: Record<string, string> = { 'signing': '签名', 'encryption': '加密', 'both': '两者' }
    return usageMap[usage.toLowerCase()] || usage
}
</script>
