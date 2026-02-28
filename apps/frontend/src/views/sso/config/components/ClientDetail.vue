<template>
    <el-dialog v-model="dialogVisible" title="客户端配置详情" width="600px">
        <el-descriptions :column="2" border>
            <el-descriptions-item label="客户端名称">{{ data.client_name }}</el-descriptions-item>
            <el-descriptions-item label="客户端类型">
                <el-tag :type="getClientTypeType(data.client_type || '')">
                    {{ getClientTypeLabel(data.client_type || '') }}
                </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="客户端 ID" :span="2">{{ data.client_id }}</el-descriptions-item>
            <el-descriptions-item label="状态">
                <el-tag :type="data.is_enabled ? 'success' : 'danger'">
                    {{ data.is_enabled ? '启用' : '禁用' }}
                </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ data.create_at }}</el-descriptions-item>
            <el-descriptions-item label="回调地址" :span="2">
                <div v-if="data.redirect_uris" class="uri-list">
                    <div v-for="(uri, index) in data.redirect_uris.split(',')" :key="index" class="uri-item">
                        <el-icon class="uri-icon"><Link /></el-icon>
                        {{ uri }}
                    </div>
                </div>
                <el-empty v-else description="未配置回调地址" :image-size="60" />
            </el-descriptions-item>
            <el-descriptions-item label="登出地址" :span="2">
                <div v-if="data.post_logout_redirect_uris" class="uri-list">
                    <div v-for="(uri, index) in data.post_logout_redirect_uris.split(',')" :key="index" class="uri-item">
                        <el-icon class="uri-icon"><Link /></el-icon>
                        {{ uri }}
                    </div>
                </div>
                <el-empty v-else description="未配置登出地址" :image-size="60" />
            </el-descriptions-item>
            <el-descriptions-item label="描述" :span="2">
                {{ data.app_description || '无描述信息' }}
            </el-descriptions-item>
        </el-descriptions>
    </el-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Link } from '@element-plus/icons-vue'
import type { SsoClient } from '@/types/sso'

interface Props {
    visible: boolean
    data: Partial<SsoClient>
}

const props = withDefaults(defineProps<Props>(), {
    data: () => ({})
})

const emit = defineEmits<{
    (e: 'update:visible', value: boolean): void
}>()

const dialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
})

const getClientTypeLabel = (type: string) => {
    const typeMap: Record<string, string> = {
        web: 'Web 应用',
        spa: 'SPA 应用',
        mobile: '移动应用',
        backend: '后端服务'
    }
    return typeMap[type?.toLowerCase()] || type
}

const getClientTypeType = (type: string) => {
    const typeMap: Record<string, string> = {
        web: 'primary',
        spa: 'success',
        mobile: 'warning',
        backend: 'info'
    }
    return typeMap[type?.toLowerCase()] || ''
}
</script>

<style scoped>
.uri-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.uri-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 10px;
    background-color: #f5f7fa;
    border-radius: 4px;
    color: #606266;
    font-size: 13px;
    word-break: break-all;
}

.uri-icon {
    color: #409EFF;
    flex-shrink: 0;
}

:deep(.el-descriptions__label) {
    width: 120px;
    font-weight: 500;
}

:deep(.el-empty) {
    margin: 10px 0;
}
</style>
