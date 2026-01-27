<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary">创建数据源</h1>
            <el-button @click="handleCancel" :icon="ArrowLeft">
                返回列表
            </el-button>
        </div>
        <el-card>
            <DataSourceForm @submit="handleSubmit" @cancel="handleCancel" :submitting="submitting" />
        </el-card>
    </div>
</template>
<script setup lang="ts">
import { createConn } from '@/api/metadata'
import DataSourceForm from '@/components/DataSourceForm.vue'
import type { MdConn } from '@/types/metadata'
import { ArrowLeft } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const submitting = ref(false)

const handleCancel = () => {
    router.push('/metadata/datasource/list')
}

const handleSubmit = async (formData: Partial<MdConn>) => {
    try {
        submitting.value = true
        await createConn(formData)
        ElMessage.success('创建成功')
        router.push('/metadata/datasource/list')
    } catch (error: any) {
        console.error('Create failed', error)
        ElMessage.error(error.message || '创建失败')
    } finally {
        submitting.value = false
    }
}
</script>
<style scoped>
.container-padding {
    padding: 20px;
}

.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}
</style>