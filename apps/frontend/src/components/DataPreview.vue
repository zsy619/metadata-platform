<template>
    <div class="data-preview" v-loading="loading">
        <div class="preview-header" v-if="columns.length > 0">
            <div class="column-info">
                <el-tag size="small" type="info">表名: {{ tableName }}</el-tag>
                <el-tag size="small" type="info">列数: {{ columns.length }}</el-tag>
                <el-tag size="small" type="success">行数: {{ total }}</el-tag>
            </div>
            <div class="column-list">
                <el-popover placement="bottom" :width="300" trigger="click">
                    <template #reference>
                        <el-button size="small" type="primary" link>
                            <el-icon><List /></el-icon>
                            查看列信息
                        </el-button>
                    </template>
                    <div class="column-detail">
                        <div class="column-detail-title">字段列表</div>
                        <el-scrollbar height="200px">
                            <ul class="column-items">
                                <li v-for="col in columns" :key="col.prop">
                                    <span class="col-name">{{ col.label }}</span>
                                    <span class="col-type">{{ col.type }}</span>
                                </li>
                            </ul>
                        </el-scrollbar>
                    </div>
                </el-popover>
            </div>
        </div>
        <el-empty v-if="!data.length && !loading" description="暂无数据" />
        <div v-else class="table-container">
            <el-table :data="data" border stripe height="100%" style="width: 100%">
                <el-table-column v-for="col in columns" :key="col.prop" :prop="col.prop" :label="col.label" min-width="120" show-overflow-tooltip />
            </el-table>
        </div>
        <div v-if="total > 0" class="pagination-wrapper">
            <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handleCurrentChange" />
        </div>
    </div>
</template>
<script setup lang="ts">
import { List } from '@element-plus/icons-vue'
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
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

const loadData = async () => {
    if (!props.connId || !props.tableName) return

    loading.value = true
    try {
        const res: any = await previewTableData(props.connId, props.tableName, props.schema, pageSize.value)

        if (res && res.data) {
            if (res.data.length > 0 && columns.value.length === 0) {
                const firstRow = res.data[0]
                columns.value = Object.keys(firstRow).map(key => ({
                    prop: key,
                    label: key,
                    type: 'string'
                }))
            }
            data.value = res.data
            total.value = res.total || res.data.length
        } else if (res && Array.isArray(res)) {
            if (res.length > 0 && columns.value.length === 0) {
                const firstRow = res[0]
                columns.value = Object.keys(firstRow).map(key => ({
                    prop: key,
                    label: key,
                    type: 'string'
                }))
            }
            data.value = res
            total.value = res.length
        } else {
            data.value = []
            total.value = 0
        }
    } catch (error: any) {
        console.error('加载数据失败', error)
        ElMessage.error(error.message || '加载预览数据失败')
    } finally {
        loading.value = false
    }
}

const handleSizeChange = (val: number) => {
    pageSize.value = val
    currentPage.value = 1
    loadData()
}

const handleCurrentChange = (val: number) => {
    currentPage.value = val
    loadData()
}

watch(() => [props.connId, props.tableName], () => {
    currentPage.value = 1
    columns.value = []
    loadData()
}, { immediate: true })
</script>
<style scoped>
.data-preview {
    height: 100%;
    width: 100%;
    display: flex;
    flex-direction: column;
}

.preview-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 15px;
    background-color: #f5f7fa;
    border-bottom: 1px solid #e4e7ed;
    flex-shrink: 0;
}

.column-info {
    display: flex;
    gap: 8px;
    align-items: center;
}

.column-list {
    display: flex;
    align-items: center;
}

.table-container {
    flex: 1;
    overflow: hidden;
}

.pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    padding: 10px 15px;
    background-color: #fff;
    border-top: 1px solid #e4e7ed;
    flex-shrink: 0;
}

.column-detail {
    padding: 0;
}

.column-detail-title {
    font-weight: 600;
    margin-bottom: 10px;
    padding-bottom: 8px;
    border-bottom: 1px solid #ebeef5;
}

.column-items {
    list-style: none;
    padding: 0;
    margin: 0;
}

.column-items li {
    display: flex;
    justify-content: space-between;
    padding: 6px 0;
    border-bottom: 1px dashed #ebeef5;
    font-size: 13px;
}

.column-items li:last-child {
    border-bottom: none;
}

.col-name {
    font-weight: 500;
    color: #303133;
}

.col-type {
    color: #909399;
    font-size: 12px;
}
</style>
