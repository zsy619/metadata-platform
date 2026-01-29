<template>
    <div class="container-padding">
        <div class="page-header">
            <h1 class="text-primary">编辑数据源</h1>
            <el-button @click="handleCancel" :icon="ArrowLeft">
                返回列表
            </el-button>
        </div>
        <el-card>
            <DataSourceForm v-if="initialData" :initial-data="initialData" :is-edit="true" @submit="handleSubmit" @cancel="handleCancel" :submitting="submitting" />
            <div v-else v-loading="loading" style="height: 200px;"></div>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import { getConnById, updateConn } from '@/api/metadata'
import DataSourceForm from '@/components/DataSourceForm.vue'
import type { MdConn } from '@/types/metadata'
import { ArrowLeft } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const route = useRoute()
const submitting = ref(false)
const loading = ref(true)
const initialData = ref<Partial<MdConn>>()

onMounted(async () => {
    const id = route.params.id as string
    if (id) {
        try {
            const res: any = await getConnById(id)
            initialData.value = res.data || res
        } catch (error) {
            console.error('Fetch failed', error)
            ElMessage.error('加载数据失败')
        } finally {
            loading.value = false
        }
    }
})

const handleCancel = () => {
    router.push('/metadata/datasource/list')
}

const handleSubmit = async (formData: Partial<MdConn>) => {
    try {
        submitting.value = true
        await updateConn(route.params.id as string, formData)
        ElMessage.success('更新成功')
        router.push('/metadata/datasource/list')
    } catch (error: any) {
        console.error('Update failed', error)
        ElMessage.error(error.message || '更新失败')
    } finally {
        submitting.value = false
    }
}
</script>
<style scoped>
/* ==================== 标准布局样式 ==================== */
.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;
}
</style>
