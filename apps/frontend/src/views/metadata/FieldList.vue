<template>
    <div class="field-list">
        <div class="page-header">
            <h1 class="text-primary">字段列表</h1>
        </div>
        <el-card>
            <div class="filter-bar m-b-lg">
                <div class="filter-left">
                    <el-select v-model="selectedTable" placeholder="选择表" style="width: 240px" @change="handleTableChange" clearable>
                        <el-option v-for="table in tables" :key="table.id" :label="table.tableName" :value="table.id" />
                    </el-select>
                    <el-input v-model="searchQuery" placeholder="搜索字段名称" clearable :prefix-icon="Search" style="width: 300px; margin-left: 10px" @clear="fetchFields" @keyup.enter="fetchFields" />
                </div>
                <div class="filter-right">
                    <el-button type="primary" :icon="Refresh" @click="fetchFields">刷新</el-button>
                </div>
            </div>
            <el-table v-loading="loading" :data="filteredFields" border stripe style="width: 100%">
                <el-table-column prop="columnName" label="名称" width="200" sortable />
                <el-table-column prop="columnType" label="类型" width="150" />
                <el-table-column prop="columnLength" label="长度" width="100" />
                <el-table-column prop="isNullable" label="可为空" width="100">
                    <template #default="scope">
                        <el-tag :type="scope.row.isNullable ? 'info' : 'danger'" size="small">
                            {{ scope.row.isNullable ? 'Yes' : 'No' }}
                        </el-tag>
                    </template>
                </el-table-column>
                <el-table-column prop="isPrimaryKey" label="主键" width="80" align="center">
                    <template #default="scope">
                        <el-icon v-if="scope.row.isPrimaryKey" color="#E6A23C">
                            <Key />
                        </el-icon>
                    </template>
                </el-table-column>
                <el-table-column prop="columnComment" label="备注" min-width="200" />
            </el-table>
        </el-card>
    </div>
</template>
<script setup lang="ts">
import { getFieldsByTableId, getTables } from '@/api/metadata'
import type { MdTable, MdTableField } from '@/types/metadata'
import { Key, Refresh, Search } from '@element-plus/icons-vue'
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// 响应式数据
const loading = ref(false)
const selectedTable = ref('')
const tables = ref<MdTable[]>([])
const allFields = ref<MdTableField[]>([])
const searchQuery = ref('')

// 计算属性
const filteredFields = computed(() => {
    if (!searchQuery.value) return allFields.value
    const query = searchQuery.value.toLowerCase()
    return allFields.value.filter(f =>
        f.columnName.toLowerCase().includes(query) ||
        (f.columnComment && f.columnComment.toLowerCase().includes(query))
    )
})

// 生命周期
onMounted(async () => {
    await fetchTables()

    // 如果 URL 中有 tableId 参数，则默认选中并加载
    if (route.query.tableId) {
        selectedTable.value = route.query.tableId as string
        fetchFields()
    }
})

// 获取表列表
const fetchTables = async () => {
    try {
        const res: any = await getTables('')
        tables.value = res?.data || []
    } catch (error) {
        console.error('获取表列表失败:', error)
    }
}

// 获取字段列表
const fetchFields = async () => {
    if (!selectedTable.value) {
        allFields.value = []
        return
    }
    loading.value = true
    try {
        const res: any = await getFieldsByTableId(selectedTable.value)
        allFields.value = res?.data || []
    } catch (error) {
        console.error('获取字段列表失败:', error)
    } finally {
        loading.value = false
    }
}

const handleTableChange = () => {
    fetchFields()
}
</script>
<style scoped>
.field-list {
    padding: 10px;
}

.page-header {
    margin-bottom: 20px;
}

.filter-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.filter-left {
    display: flex;
    align-items: center;
}
</style>
