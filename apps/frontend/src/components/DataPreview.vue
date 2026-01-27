<template>
    <div class="data-preview" v-loading="loading">
        <el-empty v-if="!data.length && !loading" description="暂无数据" />
        <el-table v-else :data="data" border stripe height="100%" style="width: 100%">
            <el-table-column v-for="col in columns" :key="col.prop" :prop="col.prop" :label="col.label" min-width="120" show-overflow-tooltip />
        </el-table>
    </div>
</template>
<script setup lang="ts">
import { previewTableData } from '@/api/metadata'
import { ElMessage } from 'element-plus'
import { ref, watch } from 'vue'

const props = defineProps({
    connId: {
        type: String,
        required: true
    },
    tableName: {
        type: String,
        required: true
    },
    schema: {
        type: String,
        default: ''
    }
})

const loading = ref(false)
const data = ref<any[]>([])
const columns = ref<any[]>([])

const loadData = async () => {
    if (!props.connId || !props.tableName) return

    loading.value = true
    try {
        const res: any = await previewTableData(props.connId, props.tableName, props.schema)
        const result = res.data || res

        if (result && result.length > 0) {
            // Generate columns from first row
            const firstRow = result[0]
            columns.value = Object.keys(firstRow).map(key => ({
                prop: key,
                label: key
            }))
            data.value = result
        } else {
            data.value = []
            columns.value = []
        }
    } catch (error: any) {
        console.error('加载数据失败', error)
        ElMessage.error(error.message || '加载预览数据失败')
    } finally {
        loading.value = false
    }
}

watch(() => [props.connId, props.tableName], () => {
    loadData()
}, { immediate: true })
</script>
<style scoped>
.data-preview {
    height: 100%;
    width: 100%;
}
</style>
