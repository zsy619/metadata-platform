<template>
    <el-dialog v-model="visible" title="选择表与视图" width="600px" @closed="handleClosed" append-to-body>
        <div class="dialog-content">
            <div class="search-bar">
                <el-input v-model="searchQuery" placeholder="搜索表名或视图名" :prefix-icon="Search" clearable @input="handleSearch" />
            </div>
            <el-table ref="tableRef" v-loading="loading" :data="filteredTables" height="400" @selection-change="handleSelectionChange" @row-click="handleRowClick">
                <el-table-column type="selection" width="55" />
                <el-table-column prop="table_name" label="表名" />
                <el-table-column prop="table_comment" label="描述" />
            </el-table>
        </div>
        <template #footer>
            <div class="dialog-footer">
                <el-button @click="visible = false">取消</el-button>
                <el-button type="primary" :disabled="selectedTables.length === 0" @click="handleConfirm">
                    确认 ({{ selectedTables.length }})
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>
<script setup lang="ts">
import { getTablesByConnId } from '@/api/metadata';
import type { MdTable } from '@/types/metadata';
import { Search } from '@element-plus/icons-vue';
import { ElTable } from 'element-plus';
import { ref } from 'vue';

const props = defineProps<{
    connId?: string
}>()

const emit = defineEmits(['confirm'])

const visible = ref(false)
const loading = ref(false)
const tableRef = ref<InstanceType<typeof ElTable>>()
const tables = ref<MdTable[]>([])
const filteredTables = ref<MdTable[]>([])
const selectedTables = ref<MdTable[]>([])
const searchQuery = ref('')

const show = () => {
    visible.value = true
    if (props.connId) {
        fetchTables()
    }
}

const fetchTables = async () => {
    if (!props.connId) return
    loading.value = true
    searchQuery.value = ''
    try {
        const res: any = await getTablesByConnId(props.connId)
        tables.value = Array.isArray(res) ? res : (res.data || [])
        filteredTables.value = tables.value
    } catch (error) {
        console.error('Failed to fetch tables', error)
    } finally {
        loading.value = false
    }
}

const handleSearch = () => {
    const query = searchQuery.value.toLowerCase()
    if (!query) {
        filteredTables.value = tables.value
        return
    }
    filteredTables.value = tables.value.filter(t =>
        t.table_name.toLowerCase().includes(query) ||
        (t.table_comment && t.table_comment.toLowerCase().includes(query))
    )
}

const handleSelectionChange = (val: MdTable[]) => {
    selectedTables.value = val
}

const handleRowClick = (row: MdTable) => {
    tableRef.value?.toggleRowSelection(row)
}

const handleConfirm = () => {
    emit('confirm', [...selectedTables.value])
    visible.value = false
}

const handleClosed = () => {
    searchQuery.value = ''
    selectedTables.value = []
}

defineExpose({
    show
})
</script>
<style scoped>
.dialog-content {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.search-bar {
    width: 100%;
}
</style>
